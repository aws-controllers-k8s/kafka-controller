// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package cluster

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	svcapitypes "github.com/aws-controllers-k8s/kafka-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/kafka-controller/pkg/sync"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	ackutil "github.com/aws-controllers-k8s/runtime/pkg/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/kafka"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/kafka/types"
	corev1 "k8s.io/api/core/v1"
)

var (
	// TerminalStatuses are the status strings that are terminal states for a
	// cluster.
	TerminalStatuses = []string{
		string(svcsdktypes.ClusterStateDeleting),
		string(svcsdktypes.ClusterStateFailed),
	}
	RequeueAfterUpdateDuration = 15 * time.Second
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		fmt.Errorf("cluster in '%s' state, cannot be modified or deleted", string(svcsdktypes.ClusterStateDeleting)),
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueWaitWhileCreating = ackrequeue.NeededAfter(
		fmt.Errorf("cluster in '%s' state, cannot be modified or deleted", string(svcsdktypes.ClusterStateCreating)),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// requeueWaitUntilCanModify returns a `ackrequeue.RequeueNeededAfter` struct
// explaining the cluster cannot be modified until it reaches an ACTIVE status.
func requeueWaitUntilCanModify(r *resource) *ackrequeue.RequeueNeededAfter {
	if r.ko.Status.State == nil {
		return nil
	}
	state := *r.ko.Status.State
	msg := fmt.Sprintf(
		"Cluster in '%s' state, cannot be modified until '%s'.",
		state, string(svcsdktypes.ClusterStateActive),
	)
	return ackrequeue.NeededAfter(
		errors.New(msg),
		ackrequeue.DefaultRequeueAfterDuration,
	)
}

// clusterHasTerminalStatus returns whether the supplied cluster is in a
// terminal state
func clusterHasTerminalStatus(r *resource) bool {
	if r.ko.Status.State == nil {
		return false
	}
	cs := *r.ko.Status.State
	for _, s := range TerminalStatuses {
		if cs == s {
			return true
		}
	}
	return false
}

// clusterActive returns true if the supplied cluster is in an
// active status
func clusterActive(r *resource) bool {
	if r.ko.Status.State == nil {
		return false
	}
	cs := *r.ko.Status.State
	return cs == string(svcsdktypes.ClusterStateActive)
}

// clusterCreating returns true if the supplied cluster is in the process
// of being created
func clusterCreating(r *resource) bool {
	if r.ko.Status.State == nil {
		return false
	}
	cs := *r.ko.Status.State
	return cs == strings.ToLower(string(svcsdktypes.ClusterStateCreating))
}

// clusterDeleting returns true if the supplied cluster is in the process
// of being deleted
func clusterDeleting(r *resource) bool {
	if r.ko.Status.State == nil {
		return false
	}
	cs := *r.ko.Status.State
	return cs == strings.ToLower(string(svcsdktypes.ClusterStateDeleting))
}

// requeueAfterAsyncUpdate returns a `ackrequeue.RequeueNeededAfter` struct
// explaining the cluster cannot be modified until after the asynchronous update
// has (first, started and then) completed and the cluster reaches an active
// status.
func requeueAfterAsyncUpdate() *ackrequeue.RequeueNeededAfter {
	return ackrequeue.NeededAfter(
		fmt.Errorf("cluster has started asynchronously updating, cannot be modified until '%s'",
			"Active"),
		RequeueAfterUpdateDuration,
	)
}

func (rm *resourceManager) customUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdate")
	defer func() { exit(err) }()

	// For asynchronous updates, latest(from ReadOne) contains the
	// outdate values for Spec fields. However the status(Cluster status)
	// is correct inside latest.
	// So we construct the updatedRes object from the desired resource to
	// obtain correct spec fields and then copy the status from latest.
	updatedRes := rm.concreteResource(desired.DeepCopy())

	// Copy status from latest since it has the current cluster state
	updatedRes.ko.Status = latest.ko.Status

	if !clusterActive(latest) {
		msg := "Cluster is in '" + *latest.ko.Status.State + "' state"
		ackcondition.SetSynced(updatedRes, corev1.ConditionFalse, &msg, nil)
		if clusterHasTerminalStatus(latest) {
			ackcondition.SetTerminal(updatedRes, corev1.ConditionTrue, &msg, nil)
			return updatedRes, nil
		}
		return updatedRes, requeueWaitUntilCanModify(latest)
	}

	switch {
	case delta.DifferentAt("Spec.Tags"):
		if err = sync.Tags(
			ctx,
			desired.ko.Spec.Tags, latest.ko.Spec.Tags,
			(*string)(latest.ko.Status.ACKResourceMetadata.ARN),
			convertToOrderedACKTags, rm.sdkapi, rm.metrics,
		); err != nil {
			return updatedRes, nil
		}
		return updatedRes, requeueAfterAsyncUpdate()

	case delta.DifferentAt("Spec.ClientAuthentication"):
		return rm.updateClientAuthentication(ctx, updatedRes, latest)
	case delta.DifferentAt("Spec.AssociatedSCRAMSecrets"):
		err = rm.syncAssociatedScramSecrets(ctx, updatedRes, latest)
		if err != nil {
			return latest, err
		}
		return updatedRes, requeueAfterAsyncUpdate()

	case delta.DifferentAt("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize"):
		return rm.updateBrokerStorage(ctx, updatedRes, latest)

	case delta.DifferentAt("Spec.BrokerNodeGroupInfo.InstanceType"):
		return rm.updateBrokerType(ctx, desired, latest)

	case delta.DifferentAt("Spec.NumberOfBrokerNodes"):
		return rm.updateNumberOfBrokerNodes(ctx, desired, latest)
	}

	return updatedRes, nil
}

// updateNumberOfBrokerNodes updates the number of broker
// nodes for the kafka cluster
func (rm *resourceManager) updateNumberOfBrokerNodes(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (updatedRes *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateNumberOfBrokerNodes")
	defer func() { exit(err) }()

	_, err = rm.sdkapi.UpdateBrokerCount(ctx, &svcsdk.UpdateBrokerCountInput{
		ClusterArn:                (*string)(latest.ko.Status.ACKResourceMetadata.ARN),
		CurrentVersion:            latest.ko.Status.CurrentVersion,
		TargetNumberOfBrokerNodes: int32OrNil(desired.ko.Spec.NumberOfBrokerNodes),
	})
	rm.metrics.RecordAPICall("UPDATE", "UpdateBrokerCount", err)
	if err != nil {
		return latest, err
	}
	message := "kafka is updating broker number of broker nodes"
	ackcondition.SetSynced(updatedRes, corev1.ConditionFalse, &message, nil)

	return desired, requeueAfterAsyncUpdate()
}

// updateBrokerType updates the broker type of the
// kafka cluster
func (rm *resourceManager) updateBrokerType(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (updatedRes *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateBrokerType")
	defer func() { exit(err) }()
	_, err = rm.sdkapi.UpdateBrokerType(ctx, &svcsdk.UpdateBrokerTypeInput{
		ClusterArn:         (*string)(latest.ko.Status.ACKResourceMetadata.ARN),
		CurrentVersion:     latest.ko.Status.CurrentVersion,
		TargetInstanceType: desired.ko.Spec.BrokerNodeGroupInfo.InstanceType,
	})
	rm.metrics.RecordAPICall("UPDATE", "UpdateBrokerType", err)
	if err != nil {
		return nil, err
	}
	message := "kafka is updating broker instanceType"
	ackcondition.SetSynced(updatedRes, corev1.ConditionFalse, &message, nil)

	return desired, requeueAfterAsyncUpdate()
}

// updateBrokerStorate updates the volumeSize of the
// kafka cluster broker storage
func (rm *resourceManager) updateBrokerStorage(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (updatedRes *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateBrokerStorage")
	defer func() { exit(err) }()

	_, err = rm.sdkapi.UpdateBrokerStorage(ctx, &svcsdk.UpdateBrokerStorageInput{
		ClusterArn:     (*string)(latest.ko.Status.ACKResourceMetadata.ARN),
		CurrentVersion: latest.ko.Status.CurrentVersion,
		TargetBrokerEBSVolumeInfo: []svcsdktypes.BrokerEBSVolumeInfo{
			{
				KafkaBrokerNodeId: aws.String("ALL"),
				VolumeSizeGB:      int32OrNil(desired.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize),
			},
		},
	})
	rm.metrics.RecordAPICall("UPDATE", "UpdateBrokerStorage", err)
	if err != nil {
		return nil, err
	}
	message := "kafka is updating broker storage"
	ackcondition.SetSynced(desired, corev1.ConditionFalse, &message, nil)
	return desired, requeueAfterAsyncUpdate()
}

// updateClientAuthentication updates the kafka cluster
// authentication settings
func (rm *resourceManager) updateClientAuthentication(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (updatedRes *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateClientAuthentication")
	defer func() { exit(err) }()

	input := &svcsdk.UpdateSecurityInput{}
	if latest.ko.Status.CurrentVersion != nil {
		input.CurrentVersion = desired.ko.Status.CurrentVersion
	}
	if latest.ko.Status.ACKResourceMetadata.ARN != nil {
		input.ClusterArn = (*string)(desired.ko.Status.ACKResourceMetadata.ARN)
	}
	if desired.ko.Spec.ClientAuthentication != nil {
		f0 := &svcsdktypes.ClientAuthentication{}
		if desired.ko.Spec.ClientAuthentication.SASL != nil {
			f0f0 := &svcsdktypes.Sasl{}
			if desired.ko.Spec.ClientAuthentication.SASL.IAM != nil &&
				desired.ko.Spec.ClientAuthentication.SASL.IAM.Enabled != nil {
				f0f0f0 := &svcsdktypes.Iam{
					Enabled: desired.ko.Spec.ClientAuthentication.SASL.IAM.Enabled,
				}
				f0f0.Iam = f0f0f0
			}
			if desired.ko.Spec.ClientAuthentication.SASL.SCRAM != nil &&
				desired.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled != nil {
				f0f0f1 := &svcsdktypes.Scram{
					Enabled: desired.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled,
				}
				f0f0.Scram = f0f0f1
			}
			f0.Sasl = f0f0
		}
		if desired.ko.Spec.ClientAuthentication.TLS != nil {
			f0f1 := &svcsdktypes.Tls{}
			if desired.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList != nil {
				f0f1.CertificateAuthorityArnList = aws.ToStringSlice(desired.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList)
			}
			if desired.ko.Spec.ClientAuthentication.TLS.Enabled != nil {
				f0f1.Enabled = desired.ko.Spec.ClientAuthentication.TLS.Enabled
			}
			f0.Tls = f0f1
		}
		if desired.ko.Spec.ClientAuthentication.Unauthenticated != nil &&
			desired.ko.Spec.ClientAuthentication.Unauthenticated.Enabled != nil {
			f0.Unauthenticated = &svcsdktypes.Unauthenticated{
				Enabled: desired.ko.Spec.ClientAuthentication.Unauthenticated.Enabled,
			}
		}
		input.ClientAuthentication = f0
	}

	_, err = rm.sdkapi.UpdateSecurity(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateSecurity", err)
	if err != nil {
		return nil, err
	}
	message := "kafka is updating the client authentication"
	ackcondition.SetSynced(desired, corev1.ConditionFalse, &message, nil)

	return desired, err
}

// syncAssociatedScramSecrets examines the Secret ARNs in the supplied Cluster
// and calls the ListScramSecrets, BatchAssociateScramSecrets and
// BatchDisassciateScramSecret APIs to ensure that the set of assciacted secrets stays in
// sync with the Cluster.Spec.AssociatedScramSecrets field, which is a list of strings
// containing Secret ARNs.
func (rm *resourceManager) syncAssociatedScramSecrets(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncAssociatedScramSecrets")
	defer func() { exit(err) }()
	toAdd := []*string{}
	toDelete := []*string{}

	existingPolicies := latest.ko.Spec.AssociatedSCRAMSecrets

	for _, p := range desired.ko.Spec.AssociatedSCRAMSecrets {
		if !ackutil.InStringPs(*p, existingPolicies) {
			toAdd = append(toAdd, p)
		}
	}

	for _, p := range existingPolicies {
		if !ackutil.InStringPs(*p, desired.ko.Spec.AssociatedSCRAMSecrets) {
			toDelete = append(toDelete, p)
		}
	}

	if len(toAdd) > 0 {
		rlog.Debug("associate scram secrets to cluster", "secret_arn", toAdd)
		if err = rm.batchAssociateScramSecret(ctx, desired, toAdd); err != nil {
			return err
		}
	}

	if len(toDelete) > 0 {
		rlog.Debug("disassociate scram secrets from cluster", "secret_arn", toDelete)
		if err = rm.batchDisassociateScramSecret(ctx, desired, toDelete); err != nil {
			return err
		}
	}

	// Set synced condition to True after successful update
	ackcondition.SetSynced(desired, corev1.ConditionFalse, nil, nil)

	return nil
}

// getAssociatedScramSecrets returns the list of scram secrets currently
// associated with the Cluster
func (rm *resourceManager) getAssociatedScramSecrets(
	ctx context.Context,
	r *resource,
) ([]*string, error) {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getAssociatedScramSecrets")
	defer func() { exit(err) }()

	input := &svcsdk.ListScramSecretsInput{}
	input.ClusterArn = (*string)(r.ko.Status.ACKResourceMetadata.ARN)
	res := []*string{}

	paginator := svcsdk.NewListScramSecretsPaginator(rm.sdkapi, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if page == nil {
			continue
		}
		// Convert []string to []*string
		for _, arn := range page.SecretArnList {
			arnCopy := arn
			res = append(res, &arnCopy)
		}
	}
	rm.metrics.RecordAPICall("READ_MANY", "ListScramSecrets", err)
	return res, err
}

// unprocessedSecrets is an error returned by the
// BatchAssociateScramSecret or Disassociate. It represents the
// secretArns that could not be associated and the reason
type unprocessedSecrets struct {
	errorCodes    []string
	errorMessages []string
	secretArns    []string
}

// Error implementation of unprocessedSecrets loops over the errorCodes
// errorMessages, and failedSecretArns
func (us unprocessedSecrets) Error() string {
	// I don't see a case where the lengths will differ
	// getting the minimum just in case, so we can avoid
	// an index out of bounds
	lenErrs := min(len(us.errorCodes), len(us.errorMessages), len(us.secretArns))
	errorMessage := ""
	for i := range lenErrs {
		errorMessage += fmt.Sprintf("ErrorCode: %s, ErrorMessage %s, SecretArn: %s\n", us.errorCodes[i], us.errorMessages[i], us.secretArns[i])
	}
	return errorMessage
}

// batchAssociateScramSecret associates the supplied scram secrets to the supplied Cluster
// resource
func (rm *resourceManager) batchAssociateScramSecret(
	ctx context.Context,
	r *resource,
	secretARNs []*string,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.associateScramSecret")
	defer func() { exit(err) }()

	input := &svcsdk.BatchAssociateScramSecretInput{}
	input.ClusterArn = (*string)(r.ko.Status.ACKResourceMetadata.ARN)
	input.SecretArnList = aws.ToStringSlice(secretARNs)
	resp, err := rm.sdkapi.BatchAssociateScramSecret(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "BatchAssociateScramSecret", err)
	if err != nil {
		return err
	}

	if len(resp.UnprocessedScramSecrets) > 0 {
		unprocessedSecrets := unprocessedSecrets{}
		for _, uss := range resp.UnprocessedScramSecrets {
			unprocessedSecrets.errorCodes = append(unprocessedSecrets.errorCodes, aws.ToString(uss.ErrorCode))
			unprocessedSecrets.errorMessages = append(unprocessedSecrets.errorMessages, aws.ToString(uss.ErrorMessage))
			unprocessedSecrets.secretArns = append(unprocessedSecrets.secretArns, aws.ToString(uss.SecretArn))
		}

		return ackerr.NewTerminalError(unprocessedSecrets)
	}

	return err
}

// batchDisassociateScramSecret disassociates the supplied scram secrets from the supplied
// Cluster resource
func (rm *resourceManager) batchDisassociateScramSecret(
	ctx context.Context,
	r *resource,
	secretARNs []*string,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.disassociateScramSecret")
	defer func() { exit(err) }()

	input := &svcsdk.BatchDisassociateScramSecretInput{}
	input.ClusterArn = (*string)(r.ko.Status.ACKResourceMetadata.ARN)
	// Convert []*string to []string
	unrefSecrets := make([]string, len(secretARNs))
	for i, s := range secretARNs {
		unrefSecrets[i] = *s
	}
	input.SecretArnList = unrefSecrets
	resp, err := rm.sdkapi.BatchDisassociateScramSecret(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "BatchDisassociateScramSecret", err)

	if len(resp.UnprocessedScramSecrets) > 0 {
		unprocessedSecrets := unprocessedSecrets{}
		for _, uss := range resp.UnprocessedScramSecrets {
			unprocessedSecrets.errorCodes = append(unprocessedSecrets.errorCodes, aws.ToString(uss.ErrorCode))
			unprocessedSecrets.errorMessages = append(unprocessedSecrets.errorCodes, aws.ToString(uss.ErrorMessage))
			unprocessedSecrets.secretArns = append(unprocessedSecrets.errorCodes, aws.ToString(uss.SecretArn))
		}

		return ackerr.NewTerminalError(unprocessedSecrets)
	}

	return err
}

// setResourceDefaults queries the MSK Cluster for the current state of the
// fields that are not returned by the ReadOne or List APIs.
func (rm *resourceManager) setResourceAdditionalFields(ctx context.Context, r *svcapitypes.Cluster) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setResourceAdditionalFields")
	defer func() { exit(err) }()

	err = rm.setBootstrapBrokerStringInformations(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

// setBootstrapBrokerStringInformations sets the bootstrapBrokerString
// information status fields.
func (rm *resourceManager) setBootstrapBrokerStringInformations(ctx context.Context, r *svcapitypes.Cluster) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setBootstrapBrokerStringInformations")
	defer func() { exit(err) }()

	var output *svcsdk.GetBootstrapBrokersOutput
	output, err = rm.sdkapi.GetBootstrapBrokers(
		ctx,
		&svcsdk.GetBootstrapBrokersInput{
			ClusterArn: (*string)(r.Status.ACKResourceMetadata.ARN),
		},
	)
	rm.metrics.RecordAPICall("GET", "GetBootstrapBrokers", err)
	if err != nil {
		return err
	}

	r.Status.BootstrapBrokerString = output.BootstrapBrokerString
	r.Status.BootstrapBrokerStringPublicSASLIAM = output.BootstrapBrokerStringPublicSaslIam
	r.Status.BootstrapBrokerStringPublicSASLSCRAM = output.BootstrapBrokerStringPublicSaslScram
	r.Status.BootstrapBrokerStringPublicTLS = output.BootstrapBrokerStringPublicTls
	r.Status.BootstrapBrokerStringSASLIAM = output.BootstrapBrokerStringSaslIam
	r.Status.BootstrapBrokerStringSASLSCRAM = output.BootstrapBrokerStringSaslScram
	r.Status.BootstrapBrokerStringTLS = output.BootstrapBrokerStringTls
	r.Status.BootstrapBrokerStringVPCConnectivitySASLIAM = output.BootstrapBrokerStringVpcConnectivitySaslIam
	r.Status.BootstrapBrokerStringVPCConnectivitySASLSCRAM = output.BootstrapBrokerStringVpcConnectivitySaslScram
	r.Status.BootstrapBrokerStringVPCConnectivityTLS = output.BootstrapBrokerStringVpcConnectivityTls
	return nil
}

func customPreCompare(_ *ackcompare.Delta, a, b *resource) {
	// Set MSK defaults
	if a.ko.Spec.BrokerNodeGroupInfo == nil {
		a.ko.Spec.BrokerNodeGroupInfo = b.ko.Spec.BrokerNodeGroupInfo
	}
	if a.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution == nil {
		a.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution = aws.String(string(svcsdktypes.BrokerAZDistributionDefault))
	}
	if a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo == nil && b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo != nil {
		a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo = b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo
	}
	if a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo != nil {
		if a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess == nil && b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo != nil {
			a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess = b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess
		}
		if a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type == nil {
			a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type = aws.String("DISABLED")
		}
	}
	if a.ko.Spec.BrokerNodeGroupInfo.SecurityGroups == nil {
		a.ko.Spec.BrokerNodeGroupInfo.SecurityGroups = b.ko.Spec.BrokerNodeGroupInfo.SecurityGroups
	}
	if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo == nil {
		a.ko.Spec.BrokerNodeGroupInfo.StorageInfo = b.ko.Spec.BrokerNodeGroupInfo.StorageInfo
	}
	if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo != nil {
		if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo == nil {
			a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo = b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo
		}
		if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize == nil {
			a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize = b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize
		}
	}

	if a.ko.Spec.ClientAuthentication == nil {
		a.ko.Spec.ClientAuthentication = b.ko.Spec.ClientAuthentication
	}
	if a.ko.Spec.EncryptionInfo == nil {
		a.ko.Spec.EncryptionInfo = b.ko.Spec.EncryptionInfo
	}
	if a.ko.Spec.EnhancedMonitoring == nil {
		a.ko.Spec.EnhancedMonitoring = aws.String(string(svcsdktypes.EnhancedMonitoringDefault))
	}
	if a.ko.Spec.OpenMonitoring == nil {
		a.ko.Spec.OpenMonitoring = b.ko.Spec.OpenMonitoring
	}
	if a.ko.Spec.StorageMode == nil {
		a.ko.Spec.StorageMode = aws.String(string(svcsdktypes.StorageModeLocal))
	}
}

func int32OrNil(num *int64) *int32 {
	if num == nil {
		return nil
	}

	return aws.Int32(int32(*num))
}
