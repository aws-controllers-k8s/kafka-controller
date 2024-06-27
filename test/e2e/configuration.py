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

"""Utilities for working with Configuration resources"""

import datetime
import time
import typing

import boto3
import pytest

DELETE_WAIT_TIME_TIMEOUT_SECONDS = 60*10

ConfigurationMatchFunc = typing.NewType(
    'ConfigurationMatchFunc',
    typing.Callable[[dict], bool],
)

class StateMatcher:
    def __init__(self, state):
        self.match_on = state

    def __call__(self, record: dict) -> bool:
        return ('State' in record
                and record['State'] == self.match_on)


def state_matches(state: str) -> ConfigurationMatchFunc:
    return StateMatcher(state)


def wait_until_deleted(
        configuration_arn: str,
        timeout_seconds: int = DELETE_WAIT_TIME_TIMEOUT_SECONDS,
        interval_seconds: int = 15,
    ) -> None:
    """Waits until a Configuration with a supplied ID is no longer returned from
    the MSK API.

    Usage:
        from e2e.configuration import wait_until_deleted

        wait_until_deleted(configuration_name)

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while True:
        if datetime.datetime.now() >= timeout:
            pytest.fail(
                "Timed out waiting for Configuration to be "
                "deleted in MSK API"
            )
        time.sleep(interval_seconds)

        latest = get_by_arn(configuration_arn)
        if latest is None:
            break

def get_by_arn(config_arn):
    c = boto3.client("kafka")
    try:
        resp = c.describe_configuration(Arn=config_arn)
        return resp
    except c.exceptions.BadRequestException:
        return None