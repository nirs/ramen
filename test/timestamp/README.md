# timestamp

Prefix each line of `stdout` and `stderr` with a timestamp (`[HH:MM:SS.mmm]`).
Useful for profiling commands that do not log timing themselves, such as
`podman build` or long CI steps.

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
$ ./timestamp bash -c 'echo started; echo quick; sleep 1.234; echo slow'
[17:53:43.953] started
[17:53:43.954] quick
[17:53:45.200] slow
```

Exit status matches the wrapped command.

## Notes

Timestamps are per line (`\n`). Output without a newline is stamped when
the process exits. Times are local wall clock with millisecond precision.

The child uses pipes, not a TTY, so programs may buffer output. Timestamps
reflect when this tool reads a line, not when the child wrote it.

Stdin is not passed through. Use for non-interactive commands only.
