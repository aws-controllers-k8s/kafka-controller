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

"""Integration tests for the MSK Cluster resource"""

import logging
import time

import pytest

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest import tags
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common.types import CLUSTER_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import cluster

# Creating clusters and getting them into an ACTIVE state takes a LONG time --
# often 25+ minutes even for small clusters. We wait this amount of time before
# even trying to fetch a cluster's state.
CREATE_WAIT_AFTER_SECONDS = 180
DELETE_WAIT_SECONDS = 30
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait after the cluster has changed status, for the CR to update
CHECK_STATUS_WAIT_SECONDS = 60


@pytest.fixture(scope="module")
def simple_cluster():
    cluster_name = random_suffix_name("my-simple-cluster", 24)

    resources = get_bootstrap_resources()
    vpc = resources.ClusterVPC
    subnet_id_1 = vpc.public_subnets.subnet_ids[0]
    subnet_id_2 = vpc.public_subnets.subnet_ids[1]
    global secret_1, secret_2
    secret_1 = resources.SCRAMSecret1.arn
    secret_2 = resources.SCRAMSecret2.arn

    replacements = REPLACEMENT_VALUES.copy()
    replacements["CLUSTER_NAME"] = cluster_name
    replacements["SUBNET_ID_1"] = subnet_id_1
    replacements["SUBNET_ID_2"] = subnet_id_2
    replacements["SECRET_ARN"] = secret_1

    resource_data = load_resource(
        "cluster_simple",
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

    # NOTE(jaypipes): We wait until the cluster can no longer be found in MSK,
    # otherwise, trying to delete subnets referenced in the cluster
    # configuration will result in a failure...
    cluster.wait_until_deleted(cluster_name)


@service_marker
@pytest.mark.canary
class TestCluster:
    def test_crud(self, simple_cluster):
        ref, res = simple_cluster

        time.sleep(CREATE_WAIT_AFTER_SECONDS)

        # Before we update the Topic CR below, let's check to see that the
        # DisplayName field in the CR is still what we set in the original
        # Create call.
        cr = k8s.get_resource(ref)

        assert "status" in cr
        assert "ackResourceMetadata" in cr["status"]
        assert "arn" in cr["status"]["ackResourceMetadata"]
        cluster_arn = cr["status"]["ackResourceMetadata"]["arn"]

        latest = cluster.get_by_arn(cluster_arn)
        assert latest is not None
        assert "State" in latest

        cluster.wait_until(
            cluster_arn,
            cluster.state_matches("ACTIVE"),
        )

        # ensure status is updated properly once it has become active
        time.sleep(CHECK_STATUS_WAIT_SECONDS)
        condition.assert_synced(ref)

        latest_secrets = cluster.get_associated_scram_secrets(cluster_arn)
        assert len(latest_secrets) == 1
        assert secret_1 in latest_secrets

        # associate one more secret to the cluster
        updates = {
            "spec": {"associatedSCRAMSecrets": [secret_1, secret_2]},
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        latest_secrets = cluster.get_associated_scram_secrets(cluster_arn)
        assert len(latest_secrets) == 2
        assert secret_1 in latest_secrets and secret_2 in latest_secrets

        # remove all associated secrets
        updates = {
            "spec": {"associatedSCRAMSecrets": []},
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        latest_secrets = cluster.get_associated_scram_secrets(cluster_arn)
        assert len(latest_secrets) == 0


#
#        # update code path check for tags...
#        latest_tags = cluster.get_tags(cluster_name)
#        before_update_expected_tags = [
#            {
#                "Key": "tag1",
#                "Value": "val1"
#            }
#        ]
#        tags.assert_equal_without_ack_tags(
#           latest_tags,
#           before_update_expected_tags,
#        )
#        new_tags = [
#            {
#                "key": "tag2",
#                "value": "val2",
#            }
#        ]
#        updates = {
#            "spec": {"tags": new_tags},
#        }
#        k8s.patch_custom_resource(ref, updates)
#        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
#
#        after_update_expected_tags = [
#            {
#                "Key": "tag2",
#                "Value": "val2",
#            }
#        ]
#        latest_tags = role.get_tags(role_name)
#        tags.assert_equal_without_ack_tags(
#           latest_tags,
#           after_update_expected_tags,
#        )
