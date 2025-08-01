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
	"bytes"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
	_ = &acktags.Tags{}
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if (a == nil && b != nil) ||
		(a != nil && b == nil) {
		delta.Add("", a, b)
		return delta
	}
	customPreCompare(delta, a, b)

	if !reflect.DeepEqual(a.ko.Spec.AssociatedSCRAMSecretRefs, b.ko.Spec.AssociatedSCRAMSecretRefs) {
		delta.Add("Spec.AssociatedSCRAMSecretRefs", a.ko.Spec.AssociatedSCRAMSecretRefs, b.ko.Spec.AssociatedSCRAMSecretRefs)
	}
	if len(a.ko.Spec.AssociatedSCRAMSecrets) != len(b.ko.Spec.AssociatedSCRAMSecrets) {
		delta.Add("Spec.AssociatedSCRAMSecrets", a.ko.Spec.AssociatedSCRAMSecrets, b.ko.Spec.AssociatedSCRAMSecrets)
	} else if len(a.ko.Spec.AssociatedSCRAMSecrets) > 0 {
		if !ackcompare.SliceStringPEqual(a.ko.Spec.AssociatedSCRAMSecrets, b.ko.Spec.AssociatedSCRAMSecrets) {
			delta.Add("Spec.AssociatedSCRAMSecrets", a.ko.Spec.AssociatedSCRAMSecrets, b.ko.Spec.AssociatedSCRAMSecrets)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo, b.ko.Spec.BrokerNodeGroupInfo) {
		delta.Add("Spec.BrokerNodeGroupInfo", a.ko.Spec.BrokerNodeGroupInfo, b.ko.Spec.BrokerNodeGroupInfo)
	} else if a.ko.Spec.BrokerNodeGroupInfo != nil && b.ko.Spec.BrokerNodeGroupInfo != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution, b.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution) {
			delta.Add("Spec.BrokerNodeGroupInfo.BrokerAZDistribution", a.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution, b.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution)
		} else if a.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution != nil && b.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution != nil {
			if *a.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution != *b.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution {
				delta.Add("Spec.BrokerNodeGroupInfo.BrokerAZDistribution", a.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution, b.ko.Spec.BrokerNodeGroupInfo.BrokerAZDistribution)
			}
		}
		if len(a.ko.Spec.BrokerNodeGroupInfo.ClientSubnets) != len(b.ko.Spec.BrokerNodeGroupInfo.ClientSubnets) {
			delta.Add("Spec.BrokerNodeGroupInfo.ClientSubnets", a.ko.Spec.BrokerNodeGroupInfo.ClientSubnets, b.ko.Spec.BrokerNodeGroupInfo.ClientSubnets)
		} else if len(a.ko.Spec.BrokerNodeGroupInfo.ClientSubnets) > 0 {
			if !ackcompare.SliceStringPEqual(a.ko.Spec.BrokerNodeGroupInfo.ClientSubnets, b.ko.Spec.BrokerNodeGroupInfo.ClientSubnets) {
				delta.Add("Spec.BrokerNodeGroupInfo.ClientSubnets", a.ko.Spec.BrokerNodeGroupInfo.ClientSubnets, b.ko.Spec.BrokerNodeGroupInfo.ClientSubnets)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo, b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo) {
			delta.Add("Spec.BrokerNodeGroupInfo.ConnectivityInfo", a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo, b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo)
		} else if a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo != nil && b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess, b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess) {
				delta.Add("Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess", a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess, b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess)
			} else if a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess != nil && b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type, b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type) {
					delta.Add("Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type", a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type, b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type)
				} else if a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type != nil && b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type != nil {
					if *a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type != *b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type {
						delta.Add("Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type", a.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type, b.ko.Spec.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type)
					}
				}
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.InstanceType, b.ko.Spec.BrokerNodeGroupInfo.InstanceType) {
			delta.Add("Spec.BrokerNodeGroupInfo.InstanceType", a.ko.Spec.BrokerNodeGroupInfo.InstanceType, b.ko.Spec.BrokerNodeGroupInfo.InstanceType)
		} else if a.ko.Spec.BrokerNodeGroupInfo.InstanceType != nil && b.ko.Spec.BrokerNodeGroupInfo.InstanceType != nil {
			if *a.ko.Spec.BrokerNodeGroupInfo.InstanceType != *b.ko.Spec.BrokerNodeGroupInfo.InstanceType {
				delta.Add("Spec.BrokerNodeGroupInfo.InstanceType", a.ko.Spec.BrokerNodeGroupInfo.InstanceType, b.ko.Spec.BrokerNodeGroupInfo.InstanceType)
			}
		}
		if len(a.ko.Spec.BrokerNodeGroupInfo.SecurityGroups) != len(b.ko.Spec.BrokerNodeGroupInfo.SecurityGroups) {
			delta.Add("Spec.BrokerNodeGroupInfo.SecurityGroups", a.ko.Spec.BrokerNodeGroupInfo.SecurityGroups, b.ko.Spec.BrokerNodeGroupInfo.SecurityGroups)
		} else if len(a.ko.Spec.BrokerNodeGroupInfo.SecurityGroups) > 0 {
			if !ackcompare.SliceStringPEqual(a.ko.Spec.BrokerNodeGroupInfo.SecurityGroups, b.ko.Spec.BrokerNodeGroupInfo.SecurityGroups) {
				delta.Add("Spec.BrokerNodeGroupInfo.SecurityGroups", a.ko.Spec.BrokerNodeGroupInfo.SecurityGroups, b.ko.Spec.BrokerNodeGroupInfo.SecurityGroups)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.StorageInfo, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo) {
			delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo)
		} else if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo != nil && b.ko.Spec.BrokerNodeGroupInfo.StorageInfo != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo) {
				delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo)
			} else if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo != nil && b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput) {
					delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput)
				} else if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput != nil && b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput != nil {
					if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled) {
						delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled)
					} else if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled != nil && b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled != nil {
						if *a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled != *b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled {
							delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.Enabled)
						}
					}
					if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput) {
						delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput)
					} else if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput != nil && b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput != nil {
						if *a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput != *b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput {
							delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.ProvisionedThroughput.VolumeThroughput)
						}
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize) {
					delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize)
				} else if a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize != nil && b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize != nil {
					if *a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize != *b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize {
						delta.Add("Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize", a.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize, b.ko.Spec.BrokerNodeGroupInfo.StorageInfo.EBSStorageInfo.VolumeSize)
					}
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication, b.ko.Spec.ClientAuthentication) {
		delta.Add("Spec.ClientAuthentication", a.ko.Spec.ClientAuthentication, b.ko.Spec.ClientAuthentication)
	} else if a.ko.Spec.ClientAuthentication != nil && b.ko.Spec.ClientAuthentication != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.SASL, b.ko.Spec.ClientAuthentication.SASL) {
			delta.Add("Spec.ClientAuthentication.SASL", a.ko.Spec.ClientAuthentication.SASL, b.ko.Spec.ClientAuthentication.SASL)
		} else if a.ko.Spec.ClientAuthentication.SASL != nil && b.ko.Spec.ClientAuthentication.SASL != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.SASL.IAM, b.ko.Spec.ClientAuthentication.SASL.IAM) {
				delta.Add("Spec.ClientAuthentication.SASL.IAM", a.ko.Spec.ClientAuthentication.SASL.IAM, b.ko.Spec.ClientAuthentication.SASL.IAM)
			} else if a.ko.Spec.ClientAuthentication.SASL.IAM != nil && b.ko.Spec.ClientAuthentication.SASL.IAM != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.SASL.IAM.Enabled, b.ko.Spec.ClientAuthentication.SASL.IAM.Enabled) {
					delta.Add("Spec.ClientAuthentication.SASL.IAM.Enabled", a.ko.Spec.ClientAuthentication.SASL.IAM.Enabled, b.ko.Spec.ClientAuthentication.SASL.IAM.Enabled)
				} else if a.ko.Spec.ClientAuthentication.SASL.IAM.Enabled != nil && b.ko.Spec.ClientAuthentication.SASL.IAM.Enabled != nil {
					if *a.ko.Spec.ClientAuthentication.SASL.IAM.Enabled != *b.ko.Spec.ClientAuthentication.SASL.IAM.Enabled {
						delta.Add("Spec.ClientAuthentication.SASL.IAM.Enabled", a.ko.Spec.ClientAuthentication.SASL.IAM.Enabled, b.ko.Spec.ClientAuthentication.SASL.IAM.Enabled)
					}
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.SASL.SCRAM, b.ko.Spec.ClientAuthentication.SASL.SCRAM) {
				delta.Add("Spec.ClientAuthentication.SASL.SCRAM", a.ko.Spec.ClientAuthentication.SASL.SCRAM, b.ko.Spec.ClientAuthentication.SASL.SCRAM)
			} else if a.ko.Spec.ClientAuthentication.SASL.SCRAM != nil && b.ko.Spec.ClientAuthentication.SASL.SCRAM != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled, b.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled) {
					delta.Add("Spec.ClientAuthentication.SASL.SCRAM.Enabled", a.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled, b.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled)
				} else if a.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled != nil && b.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled != nil {
					if *a.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled != *b.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled {
						delta.Add("Spec.ClientAuthentication.SASL.SCRAM.Enabled", a.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled, b.ko.Spec.ClientAuthentication.SASL.SCRAM.Enabled)
					}
				}
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.TLS, b.ko.Spec.ClientAuthentication.TLS) {
			delta.Add("Spec.ClientAuthentication.TLS", a.ko.Spec.ClientAuthentication.TLS, b.ko.Spec.ClientAuthentication.TLS)
		} else if a.ko.Spec.ClientAuthentication.TLS != nil && b.ko.Spec.ClientAuthentication.TLS != nil {
			if len(a.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList) != len(b.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList) {
				delta.Add("Spec.ClientAuthentication.TLS.CertificateAuthorityARNList", a.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList, b.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList)
			} else if len(a.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList) > 0 {
				if !ackcompare.SliceStringPEqual(a.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList, b.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList) {
					delta.Add("Spec.ClientAuthentication.TLS.CertificateAuthorityARNList", a.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList, b.ko.Spec.ClientAuthentication.TLS.CertificateAuthorityARNList)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.TLS.Enabled, b.ko.Spec.ClientAuthentication.TLS.Enabled) {
				delta.Add("Spec.ClientAuthentication.TLS.Enabled", a.ko.Spec.ClientAuthentication.TLS.Enabled, b.ko.Spec.ClientAuthentication.TLS.Enabled)
			} else if a.ko.Spec.ClientAuthentication.TLS.Enabled != nil && b.ko.Spec.ClientAuthentication.TLS.Enabled != nil {
				if *a.ko.Spec.ClientAuthentication.TLS.Enabled != *b.ko.Spec.ClientAuthentication.TLS.Enabled {
					delta.Add("Spec.ClientAuthentication.TLS.Enabled", a.ko.Spec.ClientAuthentication.TLS.Enabled, b.ko.Spec.ClientAuthentication.TLS.Enabled)
				}
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.Unauthenticated, b.ko.Spec.ClientAuthentication.Unauthenticated) {
			delta.Add("Spec.ClientAuthentication.Unauthenticated", a.ko.Spec.ClientAuthentication.Unauthenticated, b.ko.Spec.ClientAuthentication.Unauthenticated)
		} else if a.ko.Spec.ClientAuthentication.Unauthenticated != nil && b.ko.Spec.ClientAuthentication.Unauthenticated != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.ClientAuthentication.Unauthenticated.Enabled, b.ko.Spec.ClientAuthentication.Unauthenticated.Enabled) {
				delta.Add("Spec.ClientAuthentication.Unauthenticated.Enabled", a.ko.Spec.ClientAuthentication.Unauthenticated.Enabled, b.ko.Spec.ClientAuthentication.Unauthenticated.Enabled)
			} else if a.ko.Spec.ClientAuthentication.Unauthenticated.Enabled != nil && b.ko.Spec.ClientAuthentication.Unauthenticated.Enabled != nil {
				if *a.ko.Spec.ClientAuthentication.Unauthenticated.Enabled != *b.ko.Spec.ClientAuthentication.Unauthenticated.Enabled {
					delta.Add("Spec.ClientAuthentication.Unauthenticated.Enabled", a.ko.Spec.ClientAuthentication.Unauthenticated.Enabled, b.ko.Spec.ClientAuthentication.Unauthenticated.Enabled)
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ConfigurationInfo, b.ko.Spec.ConfigurationInfo) {
		delta.Add("Spec.ConfigurationInfo", a.ko.Spec.ConfigurationInfo, b.ko.Spec.ConfigurationInfo)
	} else if a.ko.Spec.ConfigurationInfo != nil && b.ko.Spec.ConfigurationInfo != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ConfigurationInfo.ARN, b.ko.Spec.ConfigurationInfo.ARN) {
			delta.Add("Spec.ConfigurationInfo.ARN", a.ko.Spec.ConfigurationInfo.ARN, b.ko.Spec.ConfigurationInfo.ARN)
		} else if a.ko.Spec.ConfigurationInfo.ARN != nil && b.ko.Spec.ConfigurationInfo.ARN != nil {
			if *a.ko.Spec.ConfigurationInfo.ARN != *b.ko.Spec.ConfigurationInfo.ARN {
				delta.Add("Spec.ConfigurationInfo.ARN", a.ko.Spec.ConfigurationInfo.ARN, b.ko.Spec.ConfigurationInfo.ARN)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ConfigurationInfo.Revision, b.ko.Spec.ConfigurationInfo.Revision) {
			delta.Add("Spec.ConfigurationInfo.Revision", a.ko.Spec.ConfigurationInfo.Revision, b.ko.Spec.ConfigurationInfo.Revision)
		} else if a.ko.Spec.ConfigurationInfo.Revision != nil && b.ko.Spec.ConfigurationInfo.Revision != nil {
			if *a.ko.Spec.ConfigurationInfo.Revision != *b.ko.Spec.ConfigurationInfo.Revision {
				delta.Add("Spec.ConfigurationInfo.Revision", a.ko.Spec.ConfigurationInfo.Revision, b.ko.Spec.ConfigurationInfo.Revision)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.EncryptionInfo, b.ko.Spec.EncryptionInfo) {
		delta.Add("Spec.EncryptionInfo", a.ko.Spec.EncryptionInfo, b.ko.Spec.EncryptionInfo)
	} else if a.ko.Spec.EncryptionInfo != nil && b.ko.Spec.EncryptionInfo != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.EncryptionInfo.EncryptionAtRest, b.ko.Spec.EncryptionInfo.EncryptionAtRest) {
			delta.Add("Spec.EncryptionInfo.EncryptionAtRest", a.ko.Spec.EncryptionInfo.EncryptionAtRest, b.ko.Spec.EncryptionInfo.EncryptionAtRest)
		} else if a.ko.Spec.EncryptionInfo.EncryptionAtRest != nil && b.ko.Spec.EncryptionInfo.EncryptionAtRest != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID, b.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID) {
				delta.Add("Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID", a.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID, b.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID)
			} else if a.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID != nil && b.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID != nil {
				if *a.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID != *b.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID {
					delta.Add("Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID", a.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID, b.ko.Spec.EncryptionInfo.EncryptionAtRest.DataVolumeKMSKeyID)
				}
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.EncryptionInfo.EncryptionInTransit, b.ko.Spec.EncryptionInfo.EncryptionInTransit) {
			delta.Add("Spec.EncryptionInfo.EncryptionInTransit", a.ko.Spec.EncryptionInfo.EncryptionInTransit, b.ko.Spec.EncryptionInfo.EncryptionInTransit)
		} else if a.ko.Spec.EncryptionInfo.EncryptionInTransit != nil && b.ko.Spec.EncryptionInfo.EncryptionInTransit != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker, b.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker) {
				delta.Add("Spec.EncryptionInfo.EncryptionInTransit.ClientBroker", a.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker, b.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker)
			} else if a.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker != nil && b.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker != nil {
				if *a.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker != *b.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker {
					delta.Add("Spec.EncryptionInfo.EncryptionInTransit.ClientBroker", a.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker, b.ko.Spec.EncryptionInfo.EncryptionInTransit.ClientBroker)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster, b.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster) {
				delta.Add("Spec.EncryptionInfo.EncryptionInTransit.InCluster", a.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster, b.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster)
			} else if a.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster != nil && b.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster != nil {
				if *a.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster != *b.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster {
					delta.Add("Spec.EncryptionInfo.EncryptionInTransit.InCluster", a.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster, b.ko.Spec.EncryptionInfo.EncryptionInTransit.InCluster)
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.EnhancedMonitoring, b.ko.Spec.EnhancedMonitoring) {
		delta.Add("Spec.EnhancedMonitoring", a.ko.Spec.EnhancedMonitoring, b.ko.Spec.EnhancedMonitoring)
	} else if a.ko.Spec.EnhancedMonitoring != nil && b.ko.Spec.EnhancedMonitoring != nil {
		if *a.ko.Spec.EnhancedMonitoring != *b.ko.Spec.EnhancedMonitoring {
			delta.Add("Spec.EnhancedMonitoring", a.ko.Spec.EnhancedMonitoring, b.ko.Spec.EnhancedMonitoring)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.KafkaVersion, b.ko.Spec.KafkaVersion) {
		delta.Add("Spec.KafkaVersion", a.ko.Spec.KafkaVersion, b.ko.Spec.KafkaVersion)
	} else if a.ko.Spec.KafkaVersion != nil && b.ko.Spec.KafkaVersion != nil {
		if *a.ko.Spec.KafkaVersion != *b.ko.Spec.KafkaVersion {
			delta.Add("Spec.KafkaVersion", a.ko.Spec.KafkaVersion, b.ko.Spec.KafkaVersion)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo, b.ko.Spec.LoggingInfo) {
		delta.Add("Spec.LoggingInfo", a.ko.Spec.LoggingInfo, b.ko.Spec.LoggingInfo)
	} else if a.ko.Spec.LoggingInfo != nil && b.ko.Spec.LoggingInfo != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs, b.ko.Spec.LoggingInfo.BrokerLogs) {
			delta.Add("Spec.LoggingInfo.BrokerLogs", a.ko.Spec.LoggingInfo.BrokerLogs, b.ko.Spec.LoggingInfo.BrokerLogs)
		} else if a.ko.Spec.LoggingInfo.BrokerLogs != nil && b.ko.Spec.LoggingInfo.BrokerLogs != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs, b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs) {
				delta.Add("Spec.LoggingInfo.BrokerLogs.CloudWatchLogs", a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs, b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs)
			} else if a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs != nil && b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled) {
					delta.Add("Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled", a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled)
				} else if a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled != nil && b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled != nil {
					if *a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled != *b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled {
						delta.Add("Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled", a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.Enabled)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup, b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup) {
					delta.Add("Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup", a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup, b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup)
				} else if a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup != nil && b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup != nil {
					if *a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup != *b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup {
						delta.Add("Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup", a.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup, b.ko.Spec.LoggingInfo.BrokerLogs.CloudWatchLogs.LogGroup)
					}
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.Firehose, b.ko.Spec.LoggingInfo.BrokerLogs.Firehose) {
				delta.Add("Spec.LoggingInfo.BrokerLogs.Firehose", a.ko.Spec.LoggingInfo.BrokerLogs.Firehose, b.ko.Spec.LoggingInfo.BrokerLogs.Firehose)
			} else if a.ko.Spec.LoggingInfo.BrokerLogs.Firehose != nil && b.ko.Spec.LoggingInfo.BrokerLogs.Firehose != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream, b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream) {
					delta.Add("Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream", a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream, b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream)
				} else if a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream != nil && b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream != nil {
					if *a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream != *b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream {
						delta.Add("Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream", a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream, b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.DeliveryStream)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled) {
					delta.Add("Spec.LoggingInfo.BrokerLogs.Firehose.Enabled", a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled)
				} else if a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled != nil && b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled != nil {
					if *a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled != *b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled {
						delta.Add("Spec.LoggingInfo.BrokerLogs.Firehose.Enabled", a.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.Firehose.Enabled)
					}
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.S3, b.ko.Spec.LoggingInfo.BrokerLogs.S3) {
				delta.Add("Spec.LoggingInfo.BrokerLogs.S3", a.ko.Spec.LoggingInfo.BrokerLogs.S3, b.ko.Spec.LoggingInfo.BrokerLogs.S3)
			} else if a.ko.Spec.LoggingInfo.BrokerLogs.S3 != nil && b.ko.Spec.LoggingInfo.BrokerLogs.S3 != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket) {
					delta.Add("Spec.LoggingInfo.BrokerLogs.S3.Bucket", a.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket)
				} else if a.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket != nil && b.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket != nil {
					if *a.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket != *b.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket {
						delta.Add("Spec.LoggingInfo.BrokerLogs.S3.Bucket", a.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Bucket)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled) {
					delta.Add("Spec.LoggingInfo.BrokerLogs.S3.Enabled", a.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled)
				} else if a.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled != nil && b.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled != nil {
					if *a.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled != *b.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled {
						delta.Add("Spec.LoggingInfo.BrokerLogs.S3.Enabled", a.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Enabled)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix) {
					delta.Add("Spec.LoggingInfo.BrokerLogs.S3.Prefix", a.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix)
				} else if a.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix != nil && b.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix != nil {
					if *a.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix != *b.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix {
						delta.Add("Spec.LoggingInfo.BrokerLogs.S3.Prefix", a.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix, b.ko.Spec.LoggingInfo.BrokerLogs.S3.Prefix)
					}
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	} else if a.ko.Spec.Name != nil && b.ko.Spec.Name != nil {
		if *a.ko.Spec.Name != *b.ko.Spec.Name {
			delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.NumberOfBrokerNodes, b.ko.Spec.NumberOfBrokerNodes) {
		delta.Add("Spec.NumberOfBrokerNodes", a.ko.Spec.NumberOfBrokerNodes, b.ko.Spec.NumberOfBrokerNodes)
	} else if a.ko.Spec.NumberOfBrokerNodes != nil && b.ko.Spec.NumberOfBrokerNodes != nil {
		if *a.ko.Spec.NumberOfBrokerNodes != *b.ko.Spec.NumberOfBrokerNodes {
			delta.Add("Spec.NumberOfBrokerNodes", a.ko.Spec.NumberOfBrokerNodes, b.ko.Spec.NumberOfBrokerNodes)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.OpenMonitoring, b.ko.Spec.OpenMonitoring) {
		delta.Add("Spec.OpenMonitoring", a.ko.Spec.OpenMonitoring, b.ko.Spec.OpenMonitoring)
	} else if a.ko.Spec.OpenMonitoring != nil && b.ko.Spec.OpenMonitoring != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.OpenMonitoring.Prometheus, b.ko.Spec.OpenMonitoring.Prometheus) {
			delta.Add("Spec.OpenMonitoring.Prometheus", a.ko.Spec.OpenMonitoring.Prometheus, b.ko.Spec.OpenMonitoring.Prometheus)
		} else if a.ko.Spec.OpenMonitoring.Prometheus != nil && b.ko.Spec.OpenMonitoring.Prometheus != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.OpenMonitoring.Prometheus.JmxExporter, b.ko.Spec.OpenMonitoring.Prometheus.JmxExporter) {
				delta.Add("Spec.OpenMonitoring.Prometheus.JmxExporter", a.ko.Spec.OpenMonitoring.Prometheus.JmxExporter, b.ko.Spec.OpenMonitoring.Prometheus.JmxExporter)
			} else if a.ko.Spec.OpenMonitoring.Prometheus.JmxExporter != nil && b.ko.Spec.OpenMonitoring.Prometheus.JmxExporter != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker, b.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker) {
					delta.Add("Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker", a.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker, b.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker)
				} else if a.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker != nil && b.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker != nil {
					if *a.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker != *b.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker {
						delta.Add("Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker", a.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker, b.ko.Spec.OpenMonitoring.Prometheus.JmxExporter.EnabledInBroker)
					}
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.OpenMonitoring.Prometheus.NodeExporter, b.ko.Spec.OpenMonitoring.Prometheus.NodeExporter) {
				delta.Add("Spec.OpenMonitoring.Prometheus.NodeExporter", a.ko.Spec.OpenMonitoring.Prometheus.NodeExporter, b.ko.Spec.OpenMonitoring.Prometheus.NodeExporter)
			} else if a.ko.Spec.OpenMonitoring.Prometheus.NodeExporter != nil && b.ko.Spec.OpenMonitoring.Prometheus.NodeExporter != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker, b.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker) {
					delta.Add("Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker", a.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker, b.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker)
				} else if a.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker != nil && b.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker != nil {
					if *a.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker != *b.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker {
						delta.Add("Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker", a.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker, b.ko.Spec.OpenMonitoring.Prometheus.NodeExporter.EnabledInBroker)
					}
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.StorageMode, b.ko.Spec.StorageMode) {
		delta.Add("Spec.StorageMode", a.ko.Spec.StorageMode, b.ko.Spec.StorageMode)
	} else if a.ko.Spec.StorageMode != nil && b.ko.Spec.StorageMode != nil {
		if *a.ko.Spec.StorageMode != *b.ko.Spec.StorageMode {
			delta.Add("Spec.StorageMode", a.ko.Spec.StorageMode, b.ko.Spec.StorageMode)
		}
	}
	desiredACKTags, _ := convertToOrderedACKTags(a.ko.Spec.Tags)
	latestACKTags, _ := convertToOrderedACKTags(b.ko.Spec.Tags)
	if !ackcompare.MapStringStringEqual(desiredACKTags, latestACKTags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	}

	return delta
}
