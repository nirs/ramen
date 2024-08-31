# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

import concurrent.futures
import os
import platform
import subprocess
import tempfile

from drenv import kubectl
from drenv import commands

from . import command

IMAGE = "quay.io/ramendr/ramen-operator:latest"


def register(commands):
    parser = commands.add_parser(
        "deploy",
        help="Deploy ramen on the hub and managed clusters",
    )
    parser.set_defaults(func=run)
    command.add_common_arguments(parser)
    command.add_source_arguments(parser)
    command.add_ramen_arguments(parser)
    parser.add_argument(
        "--image",
        default=IMAGE,
        help=f"The container image to deploy (default '{IMAGE}')",
    )


def run(args):
    env = command.env_info(args)

    command.info("Preparing resources")
    command.watch("make", "-C", args.source_dir, "resources")

    with tempfile.TemporaryDirectory(prefix="ramenctl-deploy-") as tmpdir:
        tar = os.path.join(tmpdir, "image.tar")
        command.info("Saving image '%s'", args.image)
        command.watch("podman", "save", args.image, "-o", tar)

        with concurrent.futures.ThreadPoolExecutor() as executor:
            futures = []

            if env["hub"]:
                f = executor.submit(deploy, args, env["hub"], tar, "hub", distro="k8s")
                futures.append(f)

            for cluster in env["clusters"]:
                f = executor.submit(deploy, args, cluster, tar, "dr-cluster")
                futures.append(f)

            for f in concurrent.futures.as_completed(futures):
                f.result()


def deploy(args, cluster, tar, deploy_type, distro="", timeout=120):
    command.info("Loading image in cluster '%s'", cluster)
    # TODO: move to new drenv command.
    system = platform.system().lower()
    if system == "linux":
        minikube_load(cluster, tar)
    elif system == "darwin":
        lima_load(cluster, tar)
    else:
        raise RuntimeError(f"Don't know how to load image on {system}")

    command.info("Deploying ramen operator in cluster '%s'", cluster)
    overlay = os.path.join(args.source_dir, f"config/{deploy_type}/default", distro)
    yaml = kubectl.kustomize(overlay, load_restrictor="LoadRestrictionsNone")
    kubectl.apply("--filename=-", input=yaml, context=cluster, log=command.debug)

    deploy = f"ramen-{deploy_type}-operator"
    command.info("Waiting until '%s' is rolled out in cluster '%s'", deploy, cluster)
    kubectl.rollout(
        "status",
        f"deploy/{deploy}",
        f"--namespace={args.ramen_namespace}",
        f"--timeout={timeout}s",
        context=cluster,
        log=command.debug,
    )


def minikube_load(cluster, tar):
    command.watch("minikube", "--profile", cluster, "image", "load", tar)


def lima_load(cluster, tar):
    cmd = [
        "limactl",
        "shell",
        cluster,
        "sudo",
        "nerdctl",
        "--namespace",
        "k8s.io",
        "load",
    ]
    command.debug("Running %s", cmd)
    with open(tar) as f:
        cp = subprocess.run(
            cmd,
            stdin=f,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        error = cp.stderr.decode(errors="replace")
        if cp.returncode != 0:
            raise commands.Error(
                cmd,
                error,
                exitcode=cp.returncode,
                output=cp.stdout.decode(errors="replace"),
            )
        command.debug("%s", error)
