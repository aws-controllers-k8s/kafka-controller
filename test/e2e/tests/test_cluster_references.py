# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# 	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for MSK Cluster resource references (securityGroupRefs)."""

import time

import pytest

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name

from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common.types import CLUSTER_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import cluster

# The controller should report the reference as unresolved well within this
# window. No MSK cluster is ever created (creation is blocked on the
# unresolved reference), so the test is fast and incurs no AWS cost.
RESOLVE_WAIT_PERIODS = 18
RESOLVE_PERIOD_LENGTH = 10
DELETE_WAIT_SECONDS = 300


@pytest.fixture(scope="module")
def cluster_with_unresolvable_sg_ref():
    cluster_name = random_suffix_name("ack-sg-ref", 24)
    # Reference an ec2 SecurityGroup custom resource that is never created so
    # that the reference can not be resolved by the controller.
    sg_ref_name = random_suffix_name("nonexistent-sg", 24)

    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]

    replacements = REPLACEMENT_VALUES.copy()
    replacements["CLUSTER_NAME"] = cluster_name
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2
    replacements["SECURITY_GROUP_REF_NAME"] = sg_ref_name

    resource_data = load_resource(
        "cluster_security_group_ref",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP,
        CRD_VERSION,
        CLUSTER_RESOURCE_PLURAL,
        cluster_name,
        namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr, cluster_name)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_SECONDS,
    )
    assert deleted


@service_marker
@pytest.mark.canary
class TestClusterSecurityGroupRef:
    def test_unresolved_security_group_ref_blocks_create(
        self, cluster_with_unresolvable_sg_ref
    ):
        ref, _, cluster_name = cluster_with_unresolvable_sg_ref

        # Poll until the controller reports on reference resolution.
        cond = None
        for _ in range(RESOLVE_WAIT_PERIODS):
            cond = k8s.get_resource_condition(
                ref, condition.CONDITION_TYPE_REFERENCES_RESOLVED
            )
            if cond is not None:
                break
            time.sleep(RESOLVE_PERIOD_LENGTH)

        assert cond is not None, (
            "expected an ACK.ReferencesResolved condition because the resource "
            "uses securityGroupRefs"
        )

        # The referenced SecurityGroup can not be resolved, so the reference
        # must NOT be resolved successfully. Depending on whether the ec2
        # SecurityGroup CRD is installed in the test cluster, ACK reports:
        #   - status="False"   (CRD present, referenced CR not found), or
        #   - status="Unknown" (CRD absent -> "no matches for kind
        #                        SecurityGroup"); this is the case in this
        #                        suite, which does not install the ec2
        #                        controller/CRDs.
        # In both cases the reference is unresolved and creation must be blocked.
        assert cond.get("status") != "True", (
            f"expected securityGroupRefs to be unresolved, got condition {cond}"
        )

        # Because the reference is unresolved, the resource must not have been
        # assigned an ARN (i.e. never created in MSK) ...
        cr = k8s.get_resource(ref)
        arn = cr.get("status", {}).get("ackResourceMetadata", {}).get("arn")
        assert arn is None

        # ... and no MSK cluster with this name should exist.
        assert cluster.get(cluster_name) is None



# ---------------------------------------------------------------------------
# Positive test: an MSK Cluster references an ACK-managed ec2 SecurityGroup via
# securityGroupRefs, and the resulting MSK cluster is actually created using
# that security group.
#
# This requires the ec2-controller (and its SecurityGroup CRD) to be running in
# the test cluster. Add it to test_config.yaml under
# cluster.configuration.additional_controllers, e.g.:
#
#   cluster:
#     configuration:
#       additional_controllers:
#       - ec2-controller@v1.13.1
#
# If the CRD is not present the test is skipped rather than failed.
# ---------------------------------------------------------------------------

EC2_CRD_GROUP = "ec2.services.k8s.aws"
EC2_CRD_VERSION = "v1alpha1"
SECURITY_GROUP_RESOURCE_PLURAL = "securitygroups"
SECURITY_GROUP_CRD_NAME = "securitygroups.ec2.services.k8s.aws"

SG_SYNC_WAIT_PERIODS = 30
SG_SYNC_PERIOD_LENGTH = 10
CREATE_WAIT_AFTER_SECONDS = 180
REFS_RESOLVED_WAIT_PERIODS = 18
REFS_RESOLVED_PERIOD_LENGTH = 10


def _require_ec2_securitygroup_crd():
    """Skips the calling test unless the ec2 SecurityGroup CRD is installed."""
    import kubernetes

    api = kubernetes.client.ApiextensionsV1Api(k8s._get_k8s_api_client())
    try:
        api.read_custom_resource_definition(SECURITY_GROUP_CRD_NAME)
    except Exception:
        pytest.skip(
            "ec2 SecurityGroup CRD not installed; install the ec2-controller "
            "(add 'ec2-controller@<version>' to cluster.configuration."
            "additional_controllers in test_config.yaml) to run this test."
        )


@pytest.fixture(scope="module")
def referenced_security_group():
    _require_ec2_securitygroup_crd()

    resources = get_bootstrap_resources()
    vpc_id = resources.ClusterVPC.vpc_id

    sg_name = random_suffix_name("ack-msk-sg", 24)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["SECURITY_GROUP_NAME"] = sg_name
    replacements["VPC_ID"] = vpc_id

    sg_data = load_resource(
        "security_group",
        additional_replacements=replacements,
    )

    sg_ref = k8s.CustomResourceReference(
        EC2_CRD_GROUP,
        EC2_CRD_VERSION,
        SECURITY_GROUP_RESOURCE_PLURAL,
        sg_name,
        namespace="default",
    )
    k8s.create_custom_resource(sg_ref, sg_data)
    k8s.wait_resource_consumed_by_controller(sg_ref)

    # Wait until the security group is created and its ID is populated.
    assert k8s.wait_on_condition(
        sg_ref,
        condition.CONDITION_TYPE_RESOURCE_SYNCED,
        "True",
        wait_periods=SG_SYNC_WAIT_PERIODS,
        period_length=SG_SYNC_PERIOD_LENGTH,
    )
    sg_cr = k8s.get_resource(sg_ref)
    sg_id = sg_cr.get("status", {}).get("id")
    assert sg_id is not None and sg_id.startswith("sg-")

    yield (sg_ref, sg_name, sg_id)

    # Deleted after the cluster fixture (reverse setup order) so the MSK cluster
    # has released its ENIs before we remove the security group.
    k8s.delete_custom_resource(sg_ref, period_length=DELETE_WAIT_SECONDS)


@pytest.fixture(scope="module")
def cluster_using_managed_sg(referenced_security_group):
    _, sg_name, sg_id = referenced_security_group

    cluster_name = random_suffix_name("ack-msk-sgref", 24)

    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]

    replacements = REPLACEMENT_VALUES.copy()
    replacements["CLUSTER_NAME"] = cluster_name
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2
    # Point securityGroupRefs at the real, ACK-managed SecurityGroup.
    replacements["SECURITY_GROUP_REF_NAME"] = sg_name

    resource_data = load_resource(
        "cluster_security_group_ref",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP,
        CRD_VERSION,
        CLUSTER_RESOURCE_PLURAL,
        cluster_name,
        namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr, cluster_name, sg_id)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_SECONDS,
    )
    assert deleted
    cluster.wait_until_deleted(cluster_name)


@service_marker
@pytest.mark.canary
class TestClusterUsingManagedSecurityGroup:
    def test_cluster_uses_ack_managed_security_group(
        self, cluster_using_managed_sg
    ):
        ref, _, cluster_name, sg_id = cluster_using_managed_sg

        # The securityGroupRef must resolve to the ACK-managed SecurityGroup.
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_REFERENCES_RESOLVED,
            "True",
            wait_periods=REFS_RESOLVED_WAIT_PERIODS,
            period_length=REFS_RESOLVED_PERIOD_LENGTH,
        )

        # The resolved reference allows the cluster to be created in MSK.
        time.sleep(CREATE_WAIT_AFTER_SECONDS)

        cr = k8s.get_resource(ref)
        assert "status" in cr
        assert "ackResourceMetadata" in cr["status"]
        assert "arn" in cr["status"]["ackResourceMetadata"]
        cluster_arn = cr["status"]["ackResourceMetadata"]["arn"]

        cluster.wait_until(
            cluster_arn,
            cluster.state_matches("ACTIVE"),
        )

        # Verify MSK actually attached the ACK-managed security group resolved
        # from securityGroupRefs.
        latest = cluster.get_by_arn(cluster_arn)
        assert latest is not None
        attached_sgs = latest["BrokerNodeGroupInfo"]["SecurityGroups"]
        assert sg_id in attached_sgs, (
            f"expected resolved security group {sg_id} to be attached to the "
            f"cluster, found {attached_sgs}"
        )
