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

// Code generated by ack-generate. DO NOT EDIT.

package cluster

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/kafka"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/kafka-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.Kafka{}
	_ = &svcapitypes.Cluster{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.DescribeClusterOutput
	resp, err = rm.sdkapi.DescribeClusterWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeCluster", err)
	if err != nil {
		if reqErr, ok := ackerr.AWSRequestFailure(err); ok && reqErr.StatusCode() == 404 {
			return nil, ackerr.NotFound
		}
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "NotFoundException" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.ClusterInfo.BrokerNodeGroupInfo != nil {
		f1 := &svcapitypes.BrokerNodeGroupInfo{}
		if resp.ClusterInfo.BrokerNodeGroupInfo.BrokerAZDistribution != nil {
			f1.BrokerAZDistribution = resp.ClusterInfo.BrokerNodeGroupInfo.BrokerAZDistribution
		}
		if resp.ClusterInfo.BrokerNodeGroupInfo.ClientSubnets != nil {
			f1f1 := []*string{}
			for _, f1f1iter := range resp.ClusterInfo.BrokerNodeGroupInfo.ClientSubnets {
				var f1f1elem string
				f1f1elem = *f1f1iter
				f1f1 = append(f1f1, &f1f1elem)
			}
			f1.ClientSubnets = f1f1
		}
		if resp.ClusterInfo.BrokerNodeGroupInfo.ConnectivityInfo != nil {
			f1f2 := &svcapitypes.ConnectivityInfo{}
			if resp.ClusterInfo.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess != nil {
				f1f2f0 := &svcapitypes.PublicAccess{}
				if resp.ClusterInfo.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type != nil {
					f1f2f0.Type = resp.ClusterInfo.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type
				}
				f1f2.PublicAccess = f1f2f0
			}
			f1.ConnectivityInfo = f1f2
		}
		if resp.ClusterInfo.BrokerNodeGroupInfo.InstanceType != nil {
			f1.InstanceType = resp.ClusterInfo.BrokerNodeGroupInfo.InstanceType
		}
		if resp.ClusterInfo.BrokerNodeGroupInfo.SecurityGroups != nil {
			f1f4 := []*string{}
			for _, f1f4iter := range resp.ClusterInfo.BrokerNodeGroupInfo.SecurityGroups {
				var f1f4elem string
				f1f4elem = *f1f4iter
				f1f4 = append(f1f4, &f1f4elem)
			}
			f1.SecurityGroups = f1f4
		}
		if resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo != nil {
			f1f5 := &svcapitypes.StorageInfo{}
			if resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo != nil {
				f1f5f0 := &svcapitypes.EBSStorageInfo{}
				if resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.ProvisionedThroughput != nil {
					f1f5f0f0 := &svcapitypes.ProvisionedThroughput{}
					if resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.ProvisionedThroughput.Enabled != nil {
						f1f5f0f0.Enabled = resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.ProvisionedThroughput.Enabled
					}
					if resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.ProvisionedThroughput.VolumeThroughput != nil {
						f1f5f0f0.VolumeThroughput = resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.ProvisionedThroughput.VolumeThroughput
					}
					f1f5f0.ProvisionedThroughput = f1f5f0f0
				}
				if resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.VolumeSize != nil {
					f1f5f0.VolumeSize = resp.ClusterInfo.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.VolumeSize
				}
				f1f5.EBSStorageInfo = f1f5f0
			}
			f1.StorageInfo = f1f5
		}
		ko.Spec.BrokerNodeGroupInfo = f1
	} else {
		ko.Spec.BrokerNodeGroupInfo = nil
	}
	if resp.ClusterInfo.ClientAuthentication != nil {
		f2 := &svcapitypes.ClientAuthentication{}
		if resp.ClusterInfo.ClientAuthentication.Sasl != nil {
			f2f0 := &svcapitypes.SASL{}
			if resp.ClusterInfo.ClientAuthentication.Sasl.Iam != nil {
				f2f0f0 := &svcapitypes.IAM{}
				if resp.ClusterInfo.ClientAuthentication.Sasl.Iam.Enabled != nil {
					f2f0f0.Enabled = resp.ClusterInfo.ClientAuthentication.Sasl.Iam.Enabled
				}
				f2f0.IAM = f2f0f0
			}
			if resp.ClusterInfo.ClientAuthentication.Sasl.Scram != nil {
				f2f0f1 := &svcapitypes.SCRAM{}
				if resp.ClusterInfo.ClientAuthentication.Sasl.Scram.Enabled != nil {
					f2f0f1.Enabled = resp.ClusterInfo.ClientAuthentication.Sasl.Scram.Enabled
				}
				f2f0.SCRAM = f2f0f1
			}
			f2.SASL = f2f0
		}
		if resp.ClusterInfo.ClientAuthentication.Tls != nil {
			f2f1 := &svcapitypes.TLS{}
			if resp.ClusterInfo.ClientAuthentication.Tls.CertificateAuthorityArnList != nil {
				f2f1f0 := []*string{}
				for _, f2f1f0iter := range resp.ClusterInfo.ClientAuthentication.Tls.CertificateAuthorityArnList {
					var f2f1f0elem string
					f2f1f0elem = *f2f1f0iter
					f2f1f0 = append(f2f1f0, &f2f1f0elem)
				}
				f2f1.CertificateAuthorityARNList = f2f1f0
			}
			if resp.ClusterInfo.ClientAuthentication.Tls.Enabled != nil {
				f2f1.Enabled = resp.ClusterInfo.ClientAuthentication.Tls.Enabled
			}
			f2.TLS = f2f1
		}
		if resp.ClusterInfo.ClientAuthentication.Unauthenticated != nil {
			f2f2 := &svcapitypes.Unauthenticated{}
			if resp.ClusterInfo.ClientAuthentication.Unauthenticated.Enabled != nil {
				f2f2.Enabled = resp.ClusterInfo.ClientAuthentication.Unauthenticated.Enabled
			}
			f2.Unauthenticated = f2f2
		}
		ko.Spec.ClientAuthentication = f2
	} else {
		ko.Spec.ClientAuthentication = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ClusterInfo.ClusterArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ClusterInfo.ClusterArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ClusterInfo.EncryptionInfo != nil {
		f8 := &svcapitypes.EncryptionInfo{}
		if resp.ClusterInfo.EncryptionInfo.EncryptionAtRest != nil {
			f8f0 := &svcapitypes.EncryptionAtRest{}
			if resp.ClusterInfo.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyId != nil {
				f8f0.DataVolumeKMSKeyID = resp.ClusterInfo.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyId
			}
			f8.EncryptionAtRest = f8f0
		}
		if resp.ClusterInfo.EncryptionInfo.EncryptionInTransit != nil {
			f8f1 := &svcapitypes.EncryptionInTransit{}
			if resp.ClusterInfo.EncryptionInfo.EncryptionInTransit.ClientBroker != nil {
				f8f1.ClientBroker = resp.ClusterInfo.EncryptionInfo.EncryptionInTransit.ClientBroker
			}
			if resp.ClusterInfo.EncryptionInfo.EncryptionInTransit.InCluster != nil {
				f8f1.InCluster = resp.ClusterInfo.EncryptionInfo.EncryptionInTransit.InCluster
			}
			f8.EncryptionInTransit = f8f1
		}
		ko.Spec.EncryptionInfo = f8
	} else {
		ko.Spec.EncryptionInfo = nil
	}
	if resp.ClusterInfo.EnhancedMonitoring != nil {
		ko.Spec.EnhancedMonitoring = resp.ClusterInfo.EnhancedMonitoring
	} else {
		ko.Spec.EnhancedMonitoring = nil
	}
	if resp.ClusterInfo.LoggingInfo != nil {
		f10 := &svcapitypes.LoggingInfo{}
		if resp.ClusterInfo.LoggingInfo.BrokerLogs != nil {
			f10f0 := &svcapitypes.BrokerLogs{}
			if resp.ClusterInfo.LoggingInfo.BrokerLogs.CloudWatchLogs != nil {
				f10f0f0 := &svcapitypes.CloudWatchLogs{}
				if resp.ClusterInfo.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled != nil {
					f10f0f0.Enabled = resp.ClusterInfo.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled
				}
				if resp.ClusterInfo.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup != nil {
					f10f0f0.LogGroup = resp.ClusterInfo.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup
				}
				f10f0.CloudWatchLogs = f10f0f0
			}
			if resp.ClusterInfo.LoggingInfo.BrokerLogs.Firehose != nil {
				f10f0f1 := &svcapitypes.Firehose{}
				if resp.ClusterInfo.LoggingInfo.BrokerLogs.Firehose.DeliveryStream != nil {
					f10f0f1.DeliveryStream = resp.ClusterInfo.LoggingInfo.BrokerLogs.Firehose.DeliveryStream
				}
				if resp.ClusterInfo.LoggingInfo.BrokerLogs.Firehose.Enabled != nil {
					f10f0f1.Enabled = resp.ClusterInfo.LoggingInfo.BrokerLogs.Firehose.Enabled
				}
				f10f0.Firehose = f10f0f1
			}
			if resp.ClusterInfo.LoggingInfo.BrokerLogs.S3 != nil {
				f10f0f2 := &svcapitypes.S3{}
				if resp.ClusterInfo.LoggingInfo.BrokerLogs.S3.Bucket != nil {
					f10f0f2.Bucket = resp.ClusterInfo.LoggingInfo.BrokerLogs.S3.Bucket
				}
				if resp.ClusterInfo.LoggingInfo.BrokerLogs.S3.Enabled != nil {
					f10f0f2.Enabled = resp.ClusterInfo.LoggingInfo.BrokerLogs.S3.Enabled
				}
				if resp.ClusterInfo.LoggingInfo.BrokerLogs.S3.Prefix != nil {
					f10f0f2.Prefix = resp.ClusterInfo.LoggingInfo.BrokerLogs.S3.Prefix
				}
				f10f0.S3 = f10f0f2
			}
			f10.BrokerLogs = f10f0
		}
		ko.Spec.LoggingInfo = f10
	} else {
		ko.Spec.LoggingInfo = nil
	}
	if resp.ClusterInfo.NumberOfBrokerNodes != nil {
		ko.Spec.NumberOfBrokerNodes = resp.ClusterInfo.NumberOfBrokerNodes
	} else {
		ko.Spec.NumberOfBrokerNodes = nil
	}
	if resp.ClusterInfo.OpenMonitoring != nil {
		f12 := &svcapitypes.OpenMonitoringInfo{}
		if resp.ClusterInfo.OpenMonitoring.Prometheus != nil {
			f12f0 := &svcapitypes.PrometheusInfo{}
			if resp.ClusterInfo.OpenMonitoring.Prometheus.JmxExporter != nil {
				f12f0f0 := &svcapitypes.JmxExporterInfo{}
				if resp.ClusterInfo.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker != nil {
					f12f0f0.EnabledInBroker = resp.ClusterInfo.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker
				}
				f12f0.JmxExporter = f12f0f0
			}
			if resp.ClusterInfo.OpenMonitoring.Prometheus.NodeExporter != nil {
				f12f0f1 := &svcapitypes.NodeExporterInfo{}
				if resp.ClusterInfo.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker != nil {
					f12f0f1.EnabledInBroker = resp.ClusterInfo.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker
				}
				f12f0.NodeExporter = f12f0f1
			}
			f12.Prometheus = f12f0
		}
		ko.Spec.OpenMonitoring = f12
	} else {
		ko.Spec.OpenMonitoring = nil
	}
	if resp.ClusterInfo.State != nil {
		ko.Status.State = resp.ClusterInfo.State
	} else {
		ko.Status.State = nil
	}
	if resp.ClusterInfo.StorageMode != nil {
		ko.Spec.StorageMode = resp.ClusterInfo.StorageMode
	} else {
		ko.Spec.StorageMode = nil
	}
	if resp.ClusterInfo.Tags != nil {
		f16 := map[string]*string{}
		for f16key, f16valiter := range resp.ClusterInfo.Tags {
			var f16val string
			f16val = *f16valiter
			f16[f16key] = &f16val
		}
		ko.Spec.Tags = f16
	} else {
		ko.Spec.Tags = nil
	}
	if resp.ClusterInfo.ZookeeperConnectString != nil {
		ko.Status.ZookeeperConnectString = resp.ClusterInfo.ZookeeperConnectString
	} else {
		ko.Status.ZookeeperConnectString = nil
	}
	if resp.ClusterInfo.ZookeeperConnectStringTls != nil {
		ko.Status.ZookeeperConnectStringTLS = resp.ClusterInfo.ZookeeperConnectStringTls
	} else {
		ko.Status.ZookeeperConnectStringTLS = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return (r.ko.Status.ACKResourceMetadata == nil || r.ko.Status.ACKResourceMetadata.ARN == nil)

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeClusterInput, error) {
	res := &svcsdk.DescribeClusterInput{}

	if r.ko.Status.ACKResourceMetadata != nil && r.ko.Status.ACKResourceMetadata.ARN != nil {
		res.SetClusterArn(string(*r.ko.Status.ACKResourceMetadata.ARN))
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateClusterOutput
	_ = resp
	resp, err = rm.sdkapi.CreateClusterWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateCluster", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ClusterArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ClusterArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ClusterName != nil {
		ko.Spec.Name = resp.ClusterName
	} else {
		ko.Spec.Name = nil
	}
	if resp.State != nil {
		ko.Status.State = resp.State
	} else {
		ko.Status.State = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateClusterInput, error) {
	res := &svcsdk.CreateClusterInput{}

	if r.ko.Spec.BrokerNodeGroupInfo != nil {
		f0 := &svcsdk.BrokerNodeGroupInfo{}
		if r.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution != nil {
			f0.SetBrokerAZDistribution(*r.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution)
		}
		if r.ko.Spec.BrokerNodeGroupInfo.ClientSubnets != nil {
			f0f1 := []*string{}
			for _, f0f1iter := range r.ko.Spec.BrokerNodeGroupInfo.ClientSubnets {
				var f0f1elem string
				f0f1elem = *f0f1iter
				f0f1 = append(f0f1, &f0f1elem)
			}
			f0.SetClientSubnets(f0f1)
		}
		if r.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo != nil {
			f0f2 := &svcsdk.ConnectivityInfo{}
			if r.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess != nil {
				f0f2f0 := &svcsdk.PublicAccess{}
				if r.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type != nil {
					f0f2f0.SetType(*r.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type)
				}
				f0f2.SetPublicAccess(f0f2f0)
			}
			f0.SetConnectivityInfo(f0f2)
		}
		if r.ko.Spec.BrokerNodeGroupInfo.InstanceType != nil {
			f0.SetInstanceType(*r.ko.Spec.BrokerNodeGroupInfo.InstanceType)
		}
		if r.ko.Spec.BrokerNodeGroupInfo.SecurityGroups != nil {
			f0f4 := []*string{}
			for _, f0f4iter := range r.ko.Spec.BrokerNodeGroupInfo.SecurityGroups {
				var f0f4elem string
				f0f4elem = *f0f4iter
				f0f4 = append(f0f4, &f0f4elem)
			}
			f0.SetSecurityGroups(f0f4)
		}
		if r.ko.Spec.BrokerNodeGroupInfo.StorageInfo != nil {
			f0f5 := &svcsdk.StorageInfo{}
			if r.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo != nil {
				f0f5f0 := &svcsdk.EBSStorageInfo{}
				if r.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput != nil {
					f0f5f0f0 := &svcsdk.ProvisionedThroughput{}
					if r.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled != nil {
						f0f5f0f0.SetEnabled(*r.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled)
					}
					if r.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput != nil {
						f0f5f0f0.SetVolumeThroughput(*r.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput)
					}
					f0f5f0.SetProvisionedThroughput(f0f5f0f0)
				}
				if r.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize != nil {
					f0f5f0.SetVolumeSize(*r.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize)
				}
				f0f5.SetEbsStorageInfo(f0f5f0)
			}
			f0.SetStorageInfo(f0f5)
		}
		res.SetBrokerNodeGroupInfo(f0)
	}
	if r.ko.Spec.ClientAuthentication != nil {
		f1 := &svcsdk.ClientAuthentication{}
		if r.ko.Spec.ClientAuthentication.SASL != nil {
			f1f0 := &svcsdk.Sasl{}
			if r.ko.Spec.ClientAuthentication.SASL.IAM != nil {
				f1f0f0 := &svcsdk.Iam{}
				if r.ko.Spec.ClientAuthentication.SASL.IAM.Enabled != nil {
					f1f0f0.SetEnabled(*r.ko.Spec.ClientAuthentication.SASL.IAM.Enabled)
				}
				f1f0.SetIam(f1f0f0)
			}
			if r.ko.Spec.ClientAuthentication.SASL.SCRAM != nil {
				f1f0f1 := &svcsdk.Scram{}
				if r.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled != nil {
					f1f0f1.SetEnabled(*r.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled)
				}
				f1f0.SetScram(f1f0f1)
			}
			f1.SetSasl(f1f0)
		}
		if r.ko.Spec.ClientAuthentication.TLS != nil {
			f1f1 := &svcsdk.Tls{}
			if r.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList != nil {
				f1f1f0 := []*string{}
				for _, f1f1f0iter := range r.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList {
					var f1f1f0elem string
					f1f1f0elem = *f1f1f0iter
					f1f1f0 = append(f1f1f0, &f1f1f0elem)
				}
				f1f1.SetCertificateAuthorityArnList(f1f1f0)
			}
			if r.ko.Spec.ClientAuthentication.TLS.Enabled != nil {
				f1f1.SetEnabled(*r.ko.Spec.ClientAuthentication.TLS.Enabled)
			}
			f1.SetTls(f1f1)
		}
		if r.ko.Spec.ClientAuthentication.Unauthenticated != nil {
			f1f2 := &svcsdk.Unauthenticated{}
			if r.ko.Spec.ClientAuthentication.Unauthenticated.Enabled != nil {
				f1f2.SetEnabled(*r.ko.Spec.ClientAuthentication.Unauthenticated.Enabled)
			}
			f1.SetUnauthenticated(f1f2)
		}
		res.SetClientAuthentication(f1)
	}
	if r.ko.Spec.Name != nil {
		res.SetClusterName(*r.ko.Spec.Name)
	}
	if r.ko.Spec.ConfigurationInfo != nil {
		f3 := &svcsdk.ConfigurationInfo{}
		if r.ko.Spec.ConfigurationInfo.ARN != nil {
			f3.SetArn(*r.ko.Spec.ConfigurationInfo.ARN)
		}
		if r.ko.Spec.ConfigurationInfo.Revision != nil {
			f3.SetRevision(*r.ko.Spec.ConfigurationInfo.Revision)
		}
		res.SetConfigurationInfo(f3)
	}
	if r.ko.Spec.EncryptionInfo != nil {
		f4 := &svcsdk.EncryptionInfo{}
		if r.ko.Spec.EncryptionInfo.EncryptionAtRest != nil {
			f4f0 := &svcsdk.EncryptionAtRest{}
			if r.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID != nil {
				f4f0.SetDataVolumeKMSKeyId(*r.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID)
			}
			f4.SetEncryptionAtRest(f4f0)
		}
		if r.ko.Spec.EncryptionInfo.EncryptionInTransit != nil {
			f4f1 := &svcsdk.EncryptionInTransit{}
			if r.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker != nil {
				f4f1.SetClientBroker(*r.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker)
			}
			if r.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster != nil {
				f4f1.SetInCluster(*r.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster)
			}
			f4.SetEncryptionInTransit(f4f1)
		}
		res.SetEncryptionInfo(f4)
	}
	if r.ko.Spec.EnhancedMonitoring != nil {
		res.SetEnhancedMonitoring(*r.ko.Spec.EnhancedMonitoring)
	}
	if r.ko.Spec.KafkaVersion != nil {
		res.SetKafkaVersion(*r.ko.Spec.KafkaVersion)
	}
	if r.ko.Spec.LoggingInfo != nil {
		f7 := &svcsdk.LoggingInfo{}
		if r.ko.Spec.LoggingInfo.BrokerLogs != nil {
			f7f0 := &svcsdk.BrokerLogs{}
			if r.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs != nil {
				f7f0f0 := &svcsdk.CloudWatchLogs{}
				if r.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled != nil {
					f7f0f0.SetEnabled(*r.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled)
				}
				if r.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup != nil {
					f7f0f0.SetLogGroup(*r.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup)
				}
				f7f0.SetCloudWatchLogs(f7f0f0)
			}
			if r.ko.Spec.LoggingInfo.BrokerLogs.Firehose != nil {
				f7f0f1 := &svcsdk.Firehose{}
				if r.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream != nil {
					f7f0f1.SetDeliveryStream(*r.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream)
				}
				if r.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled != nil {
					f7f0f1.SetEnabled(*r.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled)
				}
				f7f0.SetFirehose(f7f0f1)
			}
			if r.ko.Spec.LoggingInfo.BrokerLogs.S3 != nil {
				f7f0f2 := &svcsdk.S3{}
				if r.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket != nil {
					f7f0f2.SetBucket(*r.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket)
				}
				if r.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled != nil {
					f7f0f2.SetEnabled(*r.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled)
				}
				if r.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix != nil {
					f7f0f2.SetPrefix(*r.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix)
				}
				f7f0.SetS3(f7f0f2)
			}
			f7.SetBrokerLogs(f7f0)
		}
		res.SetLoggingInfo(f7)
	}
	if r.ko.Spec.NumberOfBrokerNodes != nil {
		res.SetNumberOfBrokerNodes(*r.ko.Spec.NumberOfBrokerNodes)
	}
	if r.ko.Spec.OpenMonitoring != nil {
		f9 := &svcsdk.OpenMonitoringInfo{}
		if r.ko.Spec.OpenMonitoring.Prometheus != nil {
			f9f0 := &svcsdk.PrometheusInfo{}
			if r.ko.Spec.OpenMonitoring.Prometheus.JmxExporter != nil {
				f9f0f0 := &svcsdk.JmxExporterInfo{}
				if r.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker != nil {
					f9f0f0.SetEnabledInBroker(*r.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker)
				}
				f9f0.SetJmxExporter(f9f0f0)
			}
			if r.ko.Spec.OpenMonitoring.Prometheus.NodeExporter != nil {
				f9f0f1 := &svcsdk.NodeExporterInfo{}
				if r.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker != nil {
					f9f0f1.SetEnabledInBroker(*r.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker)
				}
				f9f0.SetNodeExporter(f9f0f1)
			}
			f9.SetPrometheus(f9f0)
		}
		res.SetOpenMonitoring(f9)
	}
	if r.ko.Spec.StorageMode != nil {
		res.SetStorageMode(*r.ko.Spec.StorageMode)
	}
	if r.ko.Spec.Tags != nil {
		f11 := map[string]*string{}
		for f11key, f11valiter := range r.ko.Spec.Tags {
			var f11val string
			f11val = *f11valiter
			f11[f11key] = &f11val
		}
		res.SetTags(f11)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	return nil, ackerr.NewTerminalError(ackerr.NotImplemented)
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	if clusterDeleting(r) {
		return nil, requeueWaitWhileDeleting
	}
	if clusterCreating(r) {
		return nil, requeueWaitWhileCreating
	}

	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteClusterOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteClusterWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteCluster", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteClusterInput, error) {
	res := &svcsdk.DeleteClusterInput{}

	if r.ko.Status.ACKResourceMetadata != nil && r.ko.Status.ACKResourceMetadata.ARN != nil {
		res.SetClusterArn(string(*r.ko.Status.ACKResourceMetadata.ARN))
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Cluster,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case "BadRequestException":
		return true
	default:
		return false
	}
}
