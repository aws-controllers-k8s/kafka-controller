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

"""Integration tests for the MSK VpcConnection resource"""

import time

import pytest

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest import tags
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common.types import VPCCONNECTION_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import vpcconnection

CREATE_WAIT_AFTER_SECONDS = 30
DELETE_WAIT_SECONDS = 300
MODIFY_WAIT_AFTER_SECONDS = 20
CHECK_STATUS_WAIT_SECONDS = 30

# Target MSK cluster provisioning can take >30 minutes.
CLUSTER_CREATE_TIMEOUT_SECONDS = 60 * 45
# Enabling VpcConnectivity via UpdateConnectivity triggers a heavyweight
# cluster reconfiguration (PrivateLink/NLB setup, rolling broker updates) that
# commonly takes 30-60+ minutes while the cluster sits in UPDATING.
CLUSTER_UPDATE_TIMEOUT_SECONDS = 60 * 90
CLUSTER_POLL_INTERVAL_SECONDS = 30


def _wait_target_cluster_active(kafka_client, cluster_arn, timeout_seconds=CLUSTER_CREATE_TIMEOUT_SECONDS):
    deadline = time.time() + timeout_seconds
    while time.time() < deadline:
        resp = kafka_client.describe_cluster_v2(ClusterArn=cluster_arn)
        state = resp["ClusterInfo"]["State"]
        if state == "ACTIVE":
            return
        if state in ("FAILED",):
            pytest.fail(f"target cluster entered {state} state")
        time.sleep(CLUSTER_POLL_INTERVAL_SECONDS)
    pytest.fail("timed out waiting for target cluster to become ACTIVE")


@pytest.fixture(scope="module")
def target_cluster(request, kafka_client):
    """Provisions an MSK cluster with VPC connectivity (SASL/IAM) enabled so a
    VpcConnection can target it. The ACK Cluster CRD does not expose
    VpcConnectivity, so this is created directly via the MSK API.
    """
    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_ids = vpc.public_subnets.subnet_ids[:2]
    sg_id = vpc.security_group.group_id

    cluster_name = random_suffix_name("ack-vpcconn-target", 32)
    resp = kafka_client.create_cluster_v2(
        ClusterName=cluster_name,
        Provisioned={
            "BrokerNodeGroupInfo": {
                "InstanceType": "kafka.m5.large",
                "ClientSubnets": subnet_ids,
                "SecurityGroups": [sg_id],
                "ConnectivityInfo": {
                    "VpcConnectivity": {
                        "ClientAuthentication": {
                            # MSK does not allow enabling VpcConnectivity auth
                            # schemes at cluster CREATE time. This is enabled
                            # after the cluster is ACTIVE via UpdateConnectivity.
                            "Sasl": {"Iam": {"Enabled": False}},
                        },
                    },
                },
            },
            "ClientAuthentication": {
                "Sasl": {"Iam": {"Enabled": True}},
            },
            "NumberOfBrokerNodes": 2,
            "KafkaVersion": "3.9.x",
        },
    )
    cluster_arn = resp["ClusterArn"]

    # Register teardown immediately after the cluster ARN is known so the
    # cluster is always deleted even if a later step in this fixture (waiting
    # for ACTIVE, UpdateConnectivity, the second wait) raises before `yield`.
    def _delete_target_cluster():
        try:
            kafka_client.delete_cluster(ClusterArn=cluster_arn)
        except Exception:
            pass

    request.addfinalizer(_delete_target_cluster)

    _wait_target_cluster_active(kafka_client, cluster_arn)

    # MSK requires a two-step flow: the cluster is created with VpcConnectivity
    # auth schemes disabled, then enabled via UpdateConnectivity once the
    # cluster is ACTIVE. Enable SASL/IAM VpcConnectivity now.
    resp = kafka_client.describe_cluster_v2(ClusterArn=cluster_arn)
    current_version = resp["ClusterInfo"]["CurrentVersion"]
    kafka_client.update_connectivity(
        ClusterArn=cluster_arn,
        CurrentVersion=current_version,
        ConnectivityInfo={
            "VpcConnectivity": {
                "ClientAuthentication": {
                    "Sasl": {"Iam": {"Enabled": True}},
                },
            },
        },
    )

    # UpdateConnectivity is asynchronous: the cluster moves to UPDATING then
    # back to ACTIVE. Give the state a moment to flip before polling.
    time.sleep(CLUSTER_POLL_INTERVAL_SECONDS)
    _wait_target_cluster_active(
        kafka_client,
        cluster_arn,
        timeout_seconds=CLUSTER_UPDATE_TIMEOUT_SECONDS,
    )

    yield cluster_arn


@pytest.fixture(scope="module")
def simple_vpc_connection(target_cluster):
    resources = get_bootstrap_resources()
    # The VpcConnection lives in its own DNS-enabled VPC, separate from the
    # ClusterVPC that hosts the target MSK cluster.
    vpc = resources.VpcConnectionVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]
    sg_id = vpc.security_group.group_id

    vpc_connection_name = random_suffix_name("ack-vpc-conn", 24)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["VPC_CONNECTION_NAME"] = vpc_connection_name
    replacements["AUTHENTICATION"] = "SASL_IAM"
    replacements["TARGET_CLUSTER_ARN"] = target_cluster
    replacements["VPC_ID"] = vpc.vpc_id
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2
    replacements["SECURITY_GROUP_ID"] = sg_id

    resource_data = load_resource(
        "vpcconnection_simple",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP,
        CRD_VERSION,
        VPCCONNECTION_RESOURCE_PLURAL,
        vpc_connection_name,
        namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    # Capture the ARN before deletion so we can confirm the VpcConnection is
    # removed on the AWS side as well.
    latest_cr = k8s.get_resource(ref)
    vpc_connection_arn = None
    if latest_cr and "status" in latest_cr:
        vpc_connection_arn = (
            latest_cr["status"]
            .get("ackResourceMetadata", {})
            .get("arn")
        )

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_SECONDS,
    )
    assert deleted

    # Dual verification: confirm the VpcConnection is gone on the AWS side.
    if vpc_connection_arn is not None:
        vpcconnection.wait_until_deleted(vpc_connection_arn)


@service_marker
@pytest.mark.canary
class TestVpcConnection:
    def test_crud(self, simple_vpc_connection):
        ref, _ = simple_vpc_connection

        time.sleep(CREATE_WAIT_AFTER_SECONDS)

        cr = k8s.get_resource(ref)
        assert "status" in cr
        assert "ackResourceMetadata" in cr["status"]
        assert "arn" in cr["status"]["ackResourceMetadata"]
        vpc_connection_arn = cr["status"]["ackResourceMetadata"]["arn"]

        latest = vpcconnection.get_by_arn(vpc_connection_arn)
        assert latest is not None
        assert "State" in latest

        vpcconnection.wait_until(
            vpc_connection_arn,
            vpcconnection.state_matches("AVAILABLE"),
        )

        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        condition.assert_synced(ref)

        # Verify the initial tag is present on the resource
        latest_tags = vpcconnection.get_tags(vpc_connection_arn)
        tags.assert_ack_system_tags(tags=latest_tags)
        tags.assert_equal_without_ack_tags(
            expected={"initialKey": "initialValue"},
            actual=latest_tags,
        )

        # Update the tags and verify they sync
        updates = {
            "spec": {
                "tags": {
                    "tag1": "val1",
                    "tag2": "val2",
                },
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "True",
            wait_periods=MODIFY_WAIT_AFTER_SECONDS,
        )

        cr = k8s.get_resource(ref)
        latest_tags = vpcconnection.get_tags(vpc_connection_arn)
        desired_tags = cr["spec"]["tags"]
        tags.assert_ack_system_tags(tags=latest_tags)
        tags.assert_equal_without_ack_tags(
            expected=desired_tags,
            actual=latest_tags,
        )
