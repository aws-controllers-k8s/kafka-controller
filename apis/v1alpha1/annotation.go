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

package v1alpha1

import "fmt"

var (
	// AssociatedSCRAMSecretsManagedByAnnotation is the annotation key used to
	// set the management style for the SCRAM secret associations of a cluster.
	// This annotation can only be set on a Cluster custom resource.
	//
	// The value of this annotation must be one of the following:
	//
	// - 'external':            The SCRAM secret associations are managed by an
	//                          external entity, causing the controller to
	//                          completely ignore the `spec.associatedSCRAMSecrets`
	//                          field. The controller will neither associate nor
	//                          disassociate any SCRAM secret, leaving association
	//                          management entirely to the customer.
	//
	// - 'ack-kafka-controller': The SCRAM secret associations are managed by the
	//                          ACK controller, causing the controller to reconcile
	//                          the set of associated SCRAM secrets with the value
	//                          of the `spec.associatedSCRAMSecrets` field.
	//
	// By default the SCRAM secret associations are managed by the controller. If
	// the annotation is not set, or the value is not one of the above, the
	// controller will default to managing the associations as if the annotation
	// was set to "ack-kafka-controller". Only the "external" value changes the
	// controller's behavior.
	AssociatedSCRAMSecretsManagedByAnnotation = fmt.Sprintf("%s/scram-secrets-managed-by", GroupVersion.Group)
)

const (
	// AssociatedSCRAMSecretsManagedByExternal is the value of the
	// AssociatedSCRAMSecretsManagedByAnnotation annotation that indicates that
	// the SCRAM secret associations of a cluster are managed by an external
	// entity. When set, the controller fully ignores the
	// `spec.associatedSCRAMSecrets` field.
	AssociatedSCRAMSecretsManagedByExternal = "external"
	// AssociatedSCRAMSecretsManagedByACKController is the value of the
	// AssociatedSCRAMSecretsManagedByAnnotation annotation that indicates that
	// the SCRAM secret associations of a cluster are managed by the ACK
	// controller. This is the default behavior.
	AssociatedSCRAMSecretsManagedByACKController = "ack-kafka-controller"
)
