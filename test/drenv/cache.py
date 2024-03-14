# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

import os
import shutil
import subprocess
import time

from . import commands

HOUR = 3600

# Refresh threshold in seconds.
REFRESH_SECONDS = 12 * HOUR

# Fetch threshold in seconds.
FETCH_SECONDS = 48 * HOUR


def clear(key=""):
    """
    Clear cached key. If key is not set clear the entire cache.
    """
    filename = _path(key)
    try:
        shutil.rmtree(filename)
    except FileNotFoundError:
        pass


def refresh(kustomization_dir, key, log=print):
    """
    Rebuild the cache if older than REFRESH_SECONDS.
    """
    filename = _path(key)
    if _age(filename) > REFRESH_SECONDS:
        _fetch(kustomization_dir, filename, log=log)


def get(kustomization_dir, key, log=print):
    """
    Return path to cached kustomization output. Rebuild the cache if older than
    FETCH_SECONDS.
    """
    filename = _path(key)
    if _age(filename) > FETCH_SECONDS:
        _fetch(kustomization_dir, filename, log=log)
    return filename


def _path(key):
    cache_home = os.environ.get("XDG_CACHE_HOME", ".cache")
    return os.path.expanduser(f"~/{cache_home}/drenv/{key}")


def _age(filename):
    try:
        mtime = os.path.getmtime(filename)
    except FileNotFoundError:
        mtime = 0
    return time.time() - mtime


def _fetch(kustomization_dir, dest, log=print):
    """
    Build kustomization and store the output yaml in dest.

    TODO: retry on errors.
    """
    log(f"Fetching {dest}")
    dest_dir = os.path.dirname(dest)
    os.makedirs(dest_dir, exist_ok=True)
    tmp = dest + ".tmp"
    try:
        _build_kustomization(kustomization_dir, tmp)
        os.rename(tmp, dest)
    finally:
        _silent_remove(tmp)


def _build_kustomization(kustomization_dir, dest):
    with open(dest, "w") as f:
        args = ["kustomize", "build", kustomization_dir]
        try:
            cp = subprocess.run(
                args,
                stdout=f,
                stderr=subprocess.PIPE,
            )
        except OSError as e:
            os.unlink(dest)
            raise commands.Error(args, f"Could not execute: {e}").with_exception(e)

        if cp.returncode != 0:
            error = cp.stderr.decode(errors="replace")
            raise commands.Error(args, error, exitcode=cp.returncode)

        os.fsync(f.fileno())


def _silent_remove(path):
    try:
        os.remove(path)
    except FileNotFoundError:
        pass
