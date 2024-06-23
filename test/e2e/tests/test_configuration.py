# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the MSK Cluster Configuration resource"""


import pytest
import time
import boto3


from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.common.types import CONFIG_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES

CREATE_WAIT_SECONDS = 10
MODIFY_WAIT_SECONDS = 10
DELETE_WAIT_SECONDS = 10


def get_by_arn(config_arn):
    c = boto3.client("kafka")
    try:
        resp = c.describe_configuration(Arn=config_arn)
        return resp
    except c.exceptions.BadRequestException:
        return None


@pytest.fixture(scope="module")
def simple_config():
    resource_name = random_suffix_name("my-config", 24)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["CONFIG_NAME"] = resource_name

    resource_data = load_resource(
        "configuration_simple", additional_replacements=replacements
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP,
        CRD_VERSION,
        CONFIG_RESOURCE_PLURAL,
        resource_name,
        namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    try:
        _, deleted = k8s.delete_custom_resource(ref, DELETE_WAIT_SECONDS)
        assert deleted
    except:
        pass


@service_marker
@pytest.mark.canary
class TestConfiguration:
    def test_crud(self, simple_config):
        # create the CR
        ref, cr = simple_config
        time.sleep(CREATE_WAIT_SECONDS)

        assert "status" in cr
        assert "ackResourceMetadata" in cr["status"]
        assert "arn" in cr["status"]["ackResourceMetadata"]
        config_arn = cr["status"]["ackResourceMetadata"]["arn"]

        latest = get_by_arn(config_arn)
        assert latest is not None
        assert "State" in latest
        assert latest["State"] == "ACTIVE"

        # update the CR
        description = "my description"
        updates = {
            "spec": {"description": description},
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_SECONDS)

        latest = get_by_arn(config_arn)
        assert "LatestRevision" in latest
        assert "Description" in latest["LatestRevision"]
        assert latest["LatestRevision"]["Description"] == description

        # delete the CR
        _, deleted = k8s.delete_custom_resource(ref, 2, 5)
        assert deleted
        time.sleep(DELETE_WAIT_SECONDS)

        latest = get_by_arn(config_arn)
        assert latest is None
