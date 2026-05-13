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

"""Integration tests for MSK Kafka version upgrades"""

import time

import pytest

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common.types import CLUSTER_RESOURCE_PLURAL, SERVERLESSCLUSTER_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import cluster
from e2e import serverlesscluster

CREATE_WAIT_AFTER_SECONDS = 180
DELETE_WAIT_SECONDS = 60 * 10
CHECK_STATUS_WAIT_SECONDS = 60
POLL_INTERVAL_SECONDS = 30
MODIFY_WAIT_AFTER_SECONDS = 60
LONG_UPDATE_WAIT = 60 * 40
KAFKA_VERSION_UPGRADE_WAIT = 60 * 60

INITIAL_KAFKA_VERSION = "3.7.x"
TARGET_KAFKA_VERSION = "3.8.x"


@pytest.fixture(scope="module")
def kafka_version_cluster():
    cluster_name = random_suffix_name("ack-kafka-ver", 24)

    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]

    replacements = REPLACEMENT_VALUES.copy()
    replacements["CLUSTER_NAME"] = cluster_name
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2

    resource_data = load_resource(
        "cluster_kafka_version",
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

    yield (ref, cr)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_SECONDS,
    )
    assert deleted
    cluster.wait_until_deleted(cluster_name)


@pytest.fixture(scope="module")
def kafka_version_provisioned_cluster():
    cluster_name = random_suffix_name("ack-kafka-ver-prov", 24)

    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]

    replacements = REPLACEMENT_VALUES.copy()
    replacements["PROVISIONED_CLUSTER_NAME"] = cluster_name
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2

    resource_data = load_resource(
        "provisioned_serverlesscluster_kafka_version",
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
class TestKafkaVersionUpgrade:
    def test_upgrade_cluster(self, kafka_version_cluster):
        ref, _ = kafka_version_cluster

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

        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        condition.assert_synced(ref)

        cr = k8s.get_resource(ref)
        assert cr["spec"]["kafkaVersion"] == INITIAL_KAFKA_VERSION

        updates = {
            "spec": {
                "kafkaVersion": TARGET_KAFKA_VERSION,
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(CHECK_STATUS_WAIT_SECONDS)

        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "False",
            wait_periods=MODIFY_WAIT_AFTER_SECONDS // POLL_INTERVAL_SECONDS,
            period_length=POLL_INTERVAL_SECONDS,
        )

        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "True",
            wait_periods=KAFKA_VERSION_UPGRADE_WAIT // POLL_INTERVAL_SECONDS,
            period_length=POLL_INTERVAL_SECONDS,
        )

        cluster.wait_until(
            cluster_arn,
            cluster.state_matches("ACTIVE"),
        )

        aws_cluster = cluster.get_by_arn(cluster_arn)
        assert aws_cluster is not None
        assert aws_cluster["CurrentBrokerSoftwareInfo"]["KafkaVersion"] == TARGET_KAFKA_VERSION

        cr = k8s.get_resource(ref)
        assert cr["spec"]["kafkaVersion"] == TARGET_KAFKA_VERSION

    def test_upgrade_provisioned_cluster(self, kafka_version_provisioned_cluster):
        ref, _ = kafka_version_provisioned_cluster

        time.sleep(CREATE_WAIT_AFTER_SECONDS)

        cr = k8s.get_resource(ref)
        assert "status" in cr
        assert "ackResourceMetadata" in cr["status"]
        assert "arn" in cr["status"]["ackResourceMetadata"]
        cluster_arn = cr["status"]["ackResourceMetadata"]["arn"]

        serverlesscluster.wait_until(
            cluster_arn,
            serverlesscluster.state_matches("ACTIVE"),
        )

        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        condition.assert_synced(ref)

        cr = k8s.get_resource(ref)
        assert cr["spec"]["provisioned"]["kafkaVersion"] == INITIAL_KAFKA_VERSION

        updates = {
            "spec": {
                "provisioned": {
                    "kafkaVersion": TARGET_KAFKA_VERSION,
                },
            },
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(CHECK_STATUS_WAIT_SECONDS)

        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "False",
            wait_periods=MODIFY_WAIT_AFTER_SECONDS // POLL_INTERVAL_SECONDS,
            period_length=POLL_INTERVAL_SECONDS,
        )

        assert k8s.wait_on_condition(
            ref,
            "ACK.ResourceSynced",
            "True",
            wait_periods=KAFKA_VERSION_UPGRADE_WAIT // POLL_INTERVAL_SECONDS,
            period_length=POLL_INTERVAL_SECONDS,
        )

        serverlesscluster.wait_until(
            cluster_arn,
            serverlesscluster.state_matches("ACTIVE"),
        )

        aws_cluster = serverlesscluster.get_by_arn(cluster_arn)
        assert aws_cluster is not None
        assert aws_cluster["Provisioned"]["CurrentBrokerSoftwareInfo"]["KafkaVersion"] == TARGET_KAFKA_VERSION

        cr = k8s.get_resource(ref)
        assert cr["spec"]["provisioned"]["kafkaVersion"] == TARGET_KAFKA_VERSION
