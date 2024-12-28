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

	svcapitypes "github.com/aws-controllers-k8s/kafka-controller/apis/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	ackutil "github.com/aws-controllers-k8s/runtime/pkg/util"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/kafka"
	corev1 "k8s.io/api/core/v1"
)

var (
	// TerminalStatuses are the status strings that are terminal states for a
	// cluster.
	TerminalStatuses = []string{
		svcsdk.ClusterStateDeleting,
		svcsdk.ClusterStateFailed,
	}
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		fmt.Errorf("cluster in '%s' state, cannot be modified or deleted", svcsdk.ClusterStateDeleting),
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueWaitWhileCreating = ackrequeue.NeededAfter(
		fmt.Errorf("cluster in '%s' state, cannot be modified or deleted", svcsdk.ClusterStateCreating),
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
		state, svcsdk.ClusterStateActive,
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
	return cs == svcsdk.ClusterStateActive
}

// clusterCreating returns true if the supplied cluster is in the process
// of being created
func clusterCreating(r *resource) bool {
	if r.ko.Status.State == nil {
		return false
	}
	cs := *r.ko.Status.State
	return cs == svcsdk.ClusterStateCreating
}

// clusterDeleting returns true if the supplied cluster is in the process
// of being deleted
func clusterDeleting(r *resource) bool {
	if r.ko.Status.State == nil {
		return false
	}
	cs := *r.ko.Status.State
	return cs == svcsdk.ClusterStateDeleting
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

	if clusterDeleting(latest) {
		msg := "Cluster is currently being deleted"
		ackcondition.SetSynced(updatedRes, corev1.ConditionFalse, &msg, nil)
		return updatedRes, requeueWaitWhileDeleting
	}

	if !clusterActive(latest) {
		msg := "Cluster is in '" + *latest.ko.Status.State + "' state"
		ackcondition.SetSynced(updatedRes, corev1.ConditionFalse, &msg, nil)
		if clusterHasTerminalStatus(latest) {
			ackcondition.SetTerminal(updatedRes, corev1.ConditionTrue, &msg, nil)
			return updatedRes, nil
		}
		return updatedRes, requeueWaitUntilCanModify(latest)
	}

	if delta.DifferentAt("Spec.AssociatedSCRAMSecrets") {
		err = rm.syncAssociatedScramSecrets(ctx, updatedRes, latest)
		if err != nil {
			return nil, err
		}
	}

	return updatedRes, nil
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

	err = rm.sdkapi.ListScramSecretsPagesWithContext(
		ctx, input, func(page *svcsdk.ListScramSecretsOutput, _ bool) bool {
			if page == nil {
				return true
			}
			res = append(res, page.SecretArnList...)
			return page.NextToken != nil
		},
	)
	rm.metrics.RecordAPICall("READ_MANY", "ListScramSecrets", err)
	return res, err
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
	input.SecretArnList = secretARNs
	_, err = rm.sdkapi.BatchAssociateScramSecretWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "BatchAssociateScramSecret", err)
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
	input.SecretArnList = secretARNs
	_, err = rm.sdkapi.BatchDisassociateScramSecretWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "BatchDisassociateScramSecret", err)
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
	output, err = rm.sdkapi.GetBootstrapBrokersWithContext(
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
		a.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution = aws.String(svcsdk.BrokerAZDistributionDefault)
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
		a.ko.Spec.EnhancedMonitoring = aws.String(svcsdk.EnhancedMonitoringDefault)
	}
	if a.ko.Spec.OpenMonitoring == nil {
		a.ko.Spec.OpenMonitoring = b.ko.Spec.OpenMonitoring
	}
	if a.ko.Spec.StorageMode == nil {
		a.ko.Spec.StorageMode = b.ko.Spec.StorageMode
	}
}
