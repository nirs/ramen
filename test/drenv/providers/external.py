# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

import logging

# Provider scope


def setup():
    logging.info("[external] Skipping setup for external provider")


def cleanup():
    logging.info("[external] Skipping cleanup for external provider")


# Cluster scope


def exists(profile):
    return True


def start(profile, verbose=False):
    logging.info("[%s] Skipping start for external cluster", profile["name"])


def configure(profile, existing=False):
    logging.info("[%s] Skipping configure for external cluster", profile["name"])


def stop(profile):
    logging.info("[%s] Skipping stop for external cluster", profile["name"])


def delete(profile):
    logging.info("[%s] Skipping delete for external cluster", profile["name"])


def suspend(profile):
    logging.info("[%s] Skipping suspend for external cluster", profile["name"])


def resume(profile):
    logging.info("[%s] Skipping resume for external cluster", profile["name"])


def cp(name, src, dst):
    logging.warning("[%s] cp not implemented yet for external cluster", name)


def ssh(name, command):
    logging.warning("[%s] ssh not implemented yet for external cluster", name)
