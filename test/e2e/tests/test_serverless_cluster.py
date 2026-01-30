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

"""Integration tests for the MSK ServerlessCluster resource"""

import logging
import time

import pytest

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest import tags
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common.types import SERVERLESSCLUSTER_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import serverlesscluster

CREATE_WAIT_AFTER_SECONDS = 180
DELETE_WAIT_SECONDS = 300
MODIFY_WAIT_AFTER_SECONDS = 10
CHECK_STATUS_WAIT_SECONDS = 60

@pytest.fixture(scope="module")
def simple_provisioned_cluster():
    cluster_name = random_suffix_name("my-provisioned-cluster", 24)

    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]
    global secret_1, secret_2
    secret_1 = resources.SCRAMSecret1.arn
    secret_2 = resources.SCRAMSecret2.arn

    replacements = REPLACEMENT_VALUES.copy()
    replacements["PROVISIONED_CLUSTER_NAME"] = cluster_name
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2
    replacements["SECRET_ARN"] = secret_1

    resource_data = load_resource(
        "provisioned_serverlesscluster",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP,
        CRD_VERSION,
        SERVERLESSCLUSTER_RESOURCE_PLURAL,
        cluster_name,
        namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_SECONDS,
    )
    assert deleted

    serverlesscluster.wait_until_deleted(cluster_name)

@pytest.fixture(scope="module")
def simple_serverless_cluster():
    cluster_name = random_suffix_name("my-serverless-cluster", 24)

    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]

    replacements = REPLACEMENT_VALUES.copy()
    replacements["SERVERLESS_CLUSTER_NAME"] = cluster_name
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2

    resource_data = load_resource(
        "serverlesscluster_simple",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP,
        CRD_VERSION,
        SERVERLESSCLUSTER_RESOURCE_PLURAL,
        cluster_name,
        namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_SECONDS,
    )
    assert deleted

    serverlesscluster.wait_until_deleted(cluster_name)

@service_marker
@pytest.mark.canary
class TestServerlessCluster:
    def test_provisioned_crud(self, simple_provisioned_cluster):
        ref, _ = simple_provisioned_cluster

        time.sleep(CREATE_WAIT_AFTER_SECONDS)

        cr = k8s.get_resource(ref)

        assert "status" in cr
        assert "ackResourceMetadata" in cr["status"]
        assert "arn" in cr["status"]["ackResourceMetadata"]
        cluster_arn = cr["status"]["ackResourceMetadata"]["arn"]

        latest = serverlesscluster.get_by_arn(cluster_arn)
        assert latest is not None
        assert "State" in latest

        serverlesscluster.wait_until(
            cluster_arn,
            serverlesscluster.state_matches("ACTIVE"),
        )

        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        condition.assert_synced(ref)

        # Test SCRAM secrets - first verify initial secret is associated
        latest_secrets = serverlesscluster.get_associated_scram_secrets(cluster_arn)
        assert len(latest_secrets) == 1
        assert secret_1 in latest_secrets

        # Test associating an additional secret using direct ARN references
        updates = {
            "spec": {
                "associatedSCRAMSecrets": [secret_1, secret_2],
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "True",
            wait_periods=MODIFY_WAIT_AFTER_SECONDS,
        )

        serverlesscluster.wait_until(
            cluster_arn,
            serverlesscluster.state_matches("ACTIVE"),
        )

        latest_secrets = serverlesscluster.get_associated_scram_secrets(cluster_arn)
        assert len(latest_secrets) == 2
        assert secret_1 in latest_secrets and secret_2 in latest_secrets

        # Test using resource references for SCRAM secrets
        updated_volume_size = cr['spec']['provisioned']['brokerNodeGroupInfo']['storageInfo']['ebsStorageInfo']['volumeSize'] + 10
        updates = {
            "spec": {
                'provisioned': {
                    'brokerNodeGroupInfo': {
                        "storageInfo": {
                            "ebsStorageInfo": {
                                "volumeSize": updated_volume_size
                            }
                        }
                    }
                }
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "True",
            wait_periods=MODIFY_WAIT_AFTER_SECONDS,
        )

        serverlesscluster.wait_until(
            cluster_arn,
            serverlesscluster.state_matches("ACTIVE"),
        )

        latest_cluster = serverlesscluster.get_by_arn(cluster_arn)
        assert latest_cluster is not None
        
        cr = k8s.get_resource(ref)

        latest_volume = latest_cluster['Provisioned']['BrokerNodeGroupInfo']["StorageInfo"]["EbsStorageInfo"]["VolumeSize"]
        desired_volume = cr['spec']['provisioned']['brokerNodeGroupInfo']['storageInfo']['ebsStorageInfo']['volumeSize']
        assert latest_volume == desired_volume

        # Test removing all associated secrets
        updates = {
            "spec": {
                "associatedSCRAMSecrets": [],
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        latest_secrets = serverlesscluster.get_associated_scram_secrets(cluster_arn)
        assert len(latest_secrets) == 0

        # Test adding tags to the cluster
        updates = {
            "spec": {
                "tags": {
                    "tag1": "val1",
                    "tag2": "val2"
                }
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "True",
            wait_periods=MODIFY_WAIT_AFTER_SECONDS,
        )

        cr = k8s.get_resource(ref)
        latest_tags = serverlesscluster.get_tags(cluster_arn)
        desired_tags = cr['spec']['tags']
        tags.assert_ack_system_tags(
            tags=latest_tags,
        )
        tags.assert_equal_without_ack_tags(
            expected=desired_tags,
            actual=latest_tags,
        )

        
    def test_serverless_crud(self, simple_provisioned_cluster):
        ref, _ = simple_provisioned_cluster

        time.sleep(CREATE_WAIT_AFTER_SECONDS)

        cr = k8s.get_resource(ref)

        assert "status" in cr
        assert "ackResourceMetadata" in cr["status"]
        assert "arn" in cr["status"]["ackResourceMetadata"]
        cluster_arn = cr["status"]["ackResourceMetadata"]["arn"]

        latest = serverlesscluster.get_by_arn(cluster_arn)
        assert latest is not None
        assert "State" in latest

        serverlesscluster.wait_until(
            cluster_arn,
            serverlesscluster.state_matches("ACTIVE"),
        )

        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        condition.assert_synced(ref)

        # Test adding tags to the cluster
        updates = {
            "spec": {
                "tags": {
                    "tag1": "val1",
                    "tag2": "val2"
                }
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "True",
            wait_periods=MODIFY_WAIT_AFTER_SECONDS,
        )

        cr = k8s.get_resource(ref)
        latest_tags = serverlesscluster.get_tags(cluster_arn)
        desired_tags = cr['spec']['tags']
        tags.assert_ack_system_tags(
            tags=latest_tags,
        )
        tags.assert_equal_without_ack_tags(
            expected=desired_tags,
            actual=latest_tags,
        )

        