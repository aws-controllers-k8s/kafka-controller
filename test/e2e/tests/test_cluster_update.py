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

"""Integration tests for MSK Cluster update operations (monitoring, configuration)"""

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

CREATE_WAIT_AFTER_SECONDS = 180
DELETE_WAIT_SECONDS = 60 * 10
CHECK_STATUS_WAIT_SECONDS = 60
POLL_INTERVAL_SECONDS = 30
MODIFY_WAIT_AFTER_SECONDS = 60
LONG_UPDATE_WAIT = 60 * 40


@pytest.fixture(scope="module")
def update_cluster():
    cluster_name = random_suffix_name("ack-cluster-upd", 24)

    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]

    replacements = REPLACEMENT_VALUES.copy()
    replacements["CLUSTER_NAME"] = cluster_name
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2

    resource_data = load_resource(
        "cluster_update",
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


@service_marker
@pytest.mark.canary
class TestClusterUpdate:
    def test_update_monitoring(self, update_cluster):
        ref, _ = update_cluster

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

        resources = get_bootstrap_resources()
        bucket_name = resources.LogBucket.name

        updates = {
            "spec": {
                "enhancedMonitoring": "PER_BROKER",
                "loggingInfo": {
                    "brokerLogs": {
                        "s3": {
                            "enabled": True,
                            "bucket": bucket_name,
                            "prefix": "msk-logs",
                        },
                    },
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
            wait_periods=LONG_UPDATE_WAIT // POLL_INTERVAL_SECONDS,
            period_length=POLL_INTERVAL_SECONDS,
        )

        cluster.wait_until(
            cluster_arn,
            cluster.state_matches("ACTIVE"),
        )

        aws_cluster = cluster.get_by_arn(cluster_arn)
        assert aws_cluster is not None
        assert aws_cluster["EnhancedMonitoring"] == "PER_BROKER"
        assert aws_cluster["LoggingInfo"]["BrokerLogs"]["S3"]["Enabled"] is True
        assert aws_cluster["LoggingInfo"]["BrokerLogs"]["S3"]["Bucket"] == bucket_name
        assert aws_cluster["LoggingInfo"]["BrokerLogs"]["S3"]["Prefix"] == "msk-logs"

        cr = k8s.get_resource(ref)
        assert cr["spec"]["enhancedMonitoring"] == "PER_BROKER"
        assert cr["spec"]["loggingInfo"]["brokerLogs"]["s3"]["enabled"] is True
        assert cr["spec"]["loggingInfo"]["brokerLogs"]["s3"]["bucket"] == bucket_name

    def test_update_cluster_configuration(self, update_cluster, kafka_client):
        ref, _ = update_cluster

        cr = k8s.get_resource(ref)
        cluster_arn = cr["status"]["ackResourceMetadata"]["arn"]

        cluster.wait_until(
            cluster_arn,
            cluster.state_matches("ACTIVE"),
        )

        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        condition.assert_synced(ref)

        config_resp = kafka_client.create_configuration(
            Name=random_suffix_name("ack-e2e-config", 24),
            KafkaVersions=[cr["spec"]["kafkaVersion"]],
            ServerProperties=b"auto.create.topics.enable = true\n",
        )
        config_arn = config_resp["Arn"]
        config_revision = config_resp["LatestRevision"]["Revision"]

        try:
            updates = {
                "spec": {
                    "configurationInfo": {
                        "arn": config_arn,
                        "revision": config_revision,
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
                wait_periods=LONG_UPDATE_WAIT // POLL_INTERVAL_SECONDS,
                period_length=POLL_INTERVAL_SECONDS,
            )

            cluster.wait_until(
                cluster_arn,
                cluster.state_matches("ACTIVE"),
            )

            resp = kafka_client.describe_cluster_v2(ClusterArn=cluster_arn)
            aws_cluster = resp["ClusterInfo"]
            assert aws_cluster is not None
            assert aws_cluster["Provisioned"]["CurrentBrokerSoftwareInfo"]["ConfigurationArn"] == config_arn
            assert aws_cluster["Provisioned"]["CurrentBrokerSoftwareInfo"]["ConfigurationRevision"] == config_revision

            cr = k8s.get_resource(ref)
            assert cr["spec"]["configurationInfo"]["arn"] == config_arn
            assert cr["spec"]["configurationInfo"]["revision"] == config_revision
        finally:
            try:
                kafka_client.delete_configuration(Arn=config_arn)
            except Exception:
                pass
