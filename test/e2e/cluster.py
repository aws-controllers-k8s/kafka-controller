# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Utilities for working with Cluster resources"""

import datetime
import time
import typing

import boto3
import pytest

# Creating MSK Clusters often takes >25 minutes...
DEFAULT_WAIT_UNTIL_TIMEOUT_SECONDS = 60*35
DEFAULT_WAIT_UNTIL_INTERVAL_SECONDS = 15
DEFAULT_WAIT_UNTIL_EXISTS_TIMEOUT_SECONDS = 60*10
DEFAULT_WAIT_UNTIL_EXISTS_INTERVAL_SECONDS = 15
# Deleting MSK Clusters often takes >10 minutes if the cluster's dependencies
# have been deleted already and thus the cluster will transition into a FAILED
# state.
DEFAULT_WAIT_UNTIL_DELETED_TIMEOUT_SECONDS = 60*15
DEFAULT_WAIT_UNTIL_DELETED_INTERVAL_SECONDS = 15

ClusterMatchFunc = typing.NewType(
    'ClusterMatchFunc',
    typing.Callable[[dict], bool],
)

class StateMatcher:
    def __init__(self, state):
        self.match_on = state

    def __call__(self, record: dict) -> bool:
        return ('State' in record
                and record['State'] == self.match_on)


def state_matches(state: str) -> ClusterMatchFunc:
    return StateMatcher(state)


def wait_until(
        cluster_arn: str,
        match_fn: ClusterMatchFunc,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_INTERVAL_SECONDS,
    ) -> None:
    """Waits until a cluster with a supplied ARN is returned from the MSK API
    and the matching functor returns True.

    Usage:
        from e2e.cluster import wait_until, state_matches

        wait_until(
            cluster_arn,
            state_matches("ACTIVE"),
        )

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while not match_fn(get_by_arn(cluster_arn)):
        if datetime.datetime.now() >= timeout:
            pytest.fail("failed to match Cluster before timeout")
        time.sleep(interval_seconds)


def wait_until_exists(
        cluster_name: str,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_EXISTS_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_EXISTS_INTERVAL_SECONDS,
    ) -> None:
    """Waits until a Cluster with a supplied name is returned from MSK
    ListClusters API.

    Usage:
        from e2e.cluster import wait_until_exists

        wait_until_exists(cluster_name)

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while True:
        if datetime.datetime.now() >= timeout:
            pytest.fail(
                "Timed out waiting for Cluster to exist "
                "in MSK API"
            )
        time.sleep(interval_seconds)

        latest = get(cluster_name)
        if latest is not None:
            break


def wait_until_deleted(
        cluster_name: str,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_DELETED_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_DELETED_INTERVAL_SECONDS,
    ) -> None:
    """Waits until a Cluster with a supplied ID is no longer returned from
    the MSK API.

    Usage:
        from e2e.cluster import wait_until_deleted

        wait_until_deleted(cluster_name)

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while True:
        if datetime.datetime.now() >= timeout:
            pytest.fail(
                "Timed out waiting for Cluster to be "
                "deleted in MSK API"
            )
        time.sleep(interval_seconds)

        latest = get(cluster_name)
        if latest is None:
            break


def get(cluster_name):
    """Returns a dict containing the Cluster record with the supplied Cluster
    Name from the MSK API.

    :note: MSK doesn't have a DescribeCluster API call that uses Cluster Name.
    Instead, it only has a ListClusters API with a ClusterNameFilter which
    returns all clusters that are *prefixed* with the filter string. So, we
    first call ListClusters and then call DescribeCluster with the Cluster ARN
    we get from ListClusters...

    If no such Cluster exists, returns None.
    """
    c = boto3.client('kafka')

    # NOTE(jaypipes): We deliberately do not wrap this in a try/catch because
    # we want to bubble up the exceptions that ListClusters raises.
    resp = c.list_clusters(ClusterNameFilter=cluster_name)
    for c in resp['ClusterInfoList']:
        if c['ClusterName'] == cluster_name:
            cluster_arn = c['ClusterArn']
            return get_by_arn(cluster_arn)


def get_by_arn(cluster_arn):
    """Returns a dict containing the Cluster record with the supplied Cluster
    ARN from the MSK API.

    If no such Cluster exists, returns None.
    """
    c = boto3.client('kafka')

    try:
        resp = c.describe_cluster(ClusterArn=cluster_arn)
        return resp['ClusterInfo']
    except c.exceptions.NotFoundException:
        return None


def get_tags(cluster_arn):
    """Returns a list containing the tags that have been associated to the
    supplied Cluster.

    If no such Cluster exists, returns None.
    """
    c = boto3.client('kafka')
    try:
        resp = c.list_tags_for_resource(ResourceArn=cluster_arn)
        return resp['Tags']
    except c.exceptions.NotFoundException:
        return None
