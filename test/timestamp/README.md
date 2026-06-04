# timestamp

Prefix each line of `stdout` and `stderr` with elapsed seconds since the
command started (`[seconds.mmm]`). Useful for profiling commands that do not log
timing themselves, such as `podman build`, `minikube start`, or long CI steps.

## Build

```console
make
```

The binary is `timestamp` in this directory (gitignored).

## Test

```console
make test
```

Builds `timestamp` first, then runs the tests.

## Usage

```console
./timestamp <command> [arguments...]
```

Example:

```console
$ ./timestamp bash -c 'echo start; echo quick; sleep 1.234; echo slow'
[    0.008] start
[    0.008] quick
[    1.255] slow
```

Exit status matches the wrapped command.

## Notes

Timestamps are per line (`\n`). Output without a newline is stamped when
the process exits. Elapsed time is measured from when the child process
starts, with millisecond precision. The first line is often a few
milliseconds in, not `0.000`, due to shell and pipe setup.

The child uses pipes, not a TTY, so programs may buffer output. Timestamps
reflect when this tool reads a line, not when the child wrote it.

Stdin is not passed through. Use for non-interactive commands only.
