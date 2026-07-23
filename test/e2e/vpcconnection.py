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

"""Utilities for working with VpcConnection resources"""

import datetime
import time
import typing

import boto3
import pytest

DEFAULT_WAIT_UNTIL_TIMEOUT_SECONDS = 60 * 15
DEFAULT_WAIT_UNTIL_INTERVAL_SECONDS = 15
DEFAULT_WAIT_UNTIL_DELETED_TIMEOUT_SECONDS = 60 * 15
DEFAULT_WAIT_UNTIL_DELETED_INTERVAL_SECONDS = 15

VpcConnectionMatchFunc = typing.NewType(
    "VpcConnectionMatchFunc",
    typing.Callable[[dict], bool],
)


class StateMatcher:
    def __init__(self, state):
        self.match_on = state

    def __call__(self, record: dict) -> bool:
        return "State" in record and record["State"] == self.match_on


def state_matches(state: str) -> VpcConnectionMatchFunc:
    return StateMatcher(state)


def wait_until(
    vpc_connection_arn: str,
    match_fn: VpcConnectionMatchFunc,
    timeout_seconds: int = DEFAULT_WAIT_UNTIL_TIMEOUT_SECONDS,
    interval_seconds: int = DEFAULT_WAIT_UNTIL_INTERVAL_SECONDS,
) -> None:
    """Waits until a VpcConnection with a supplied ARN is returned from the MSK
    API and the matching functor returns True.

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while not match_fn(get_by_arn(vpc_connection_arn)):
        if datetime.datetime.now() >= timeout:
            pytest.fail("failed to match VpcConnection before timeout")
        time.sleep(interval_seconds)


def wait_until_deleted(
    vpc_connection_arn: str,
    timeout_seconds: int = DEFAULT_WAIT_UNTIL_DELETED_TIMEOUT_SECONDS,
    interval_seconds: int = DEFAULT_WAIT_UNTIL_DELETED_INTERVAL_SECONDS,
) -> None:
    """Waits until a VpcConnection with a supplied ARN is no longer returned
    from the MSK API.

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while True:
        if datetime.datetime.now() >= timeout:
            pytest.fail(
                "Timed out waiting for VpcConnection to be deleted in MSK API"
            )
        time.sleep(interval_seconds)

        latest = get_by_arn(vpc_connection_arn)
        if latest is None:
            break


def get_by_arn(vpc_connection_arn):
    """Returns a dict containing the VpcConnection record with the supplied ARN
    from the MSK API.

    If no such VpcConnection exists, returns None.
    """
    c = boto3.client("kafka")

    try:
        return c.describe_vpc_connection(Arn=vpc_connection_arn)
    except c.exceptions.NotFoundException:
        return None


def get_tags(vpc_connection_arn):
    """Returns a dict containing the tags associated with the supplied
    VpcConnection.

    If no such VpcConnection exists, returns None.
    """
    c = boto3.client("kafka")
    try:
        resp = c.list_tags_for_resource(ResourceArn=vpc_connection_arn)
        return resp["Tags"]
    except c.exceptions.NotFoundException:
        return None
