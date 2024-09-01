# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

import logging
import time

from . import commands
from . import kubectl


def wait_until_ready(profile, attempts=20, delay=3):
    """
    Wait until API server /readyz endpoint report ready status.
    https://kubernetes.io/docs/reference/using-api/health-checks/
    """
    logging.debug("[%s] Waiting until API server is ready", profile["name"])
    start = time.monotonic()

    for i in range(1, attempts + 1):
        try:
            # Returns "ok" if ready, raises if not.
            kubectl.get("--raw=/readyz", context=profile["name"])
        except commands.Error as e:
            if i == attempts:
                raise
            logging.debug(
                "[%s] API server not ready: %s",
                profile["name"],
                e.error.rstrip(),
            )
            time.sleep(delay)
        else:
            logging.debug(
                "[%s] API server is ready in %.2f seconds",
                profile["name"],
                time.monotonic() - start,
            )
            break
