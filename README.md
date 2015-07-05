# jsonproc

[![Build Status](https://travis-ci.org/hfinucane/jsonproc.svg?branch=master)](https://travis-ci.org/hfinucane/jsonproc)

A read-only `/proc` to json bridge. In general, the URL scheme looks like:

    /           # Everything in /proc
    /loadavg    # The contents of /proc/loadavg
    /proc/1     # All About Init

When hitting a directory, you should expect a blob that looks like this:

    { 
      "path": "/proc/sys/",
      "files": [],
      "dirs": ["abi", "debug", "dev", "fs", "kernel", "net", "vm"]
    }

When hitting a file, you should expect a blob that looks like this:

    {
      "path": "/proc/loadavg",
      "contents": "0.09 0.12 0.17 1/613 4319\n"
    }

Errors are signaled both by the appearance of the "err" field, and with a 500 code:

    $ curl -v localhost:9234/x; echo
    * Connected to localhost (127.0.0.1) port 9234 (#0)
    > GET /x HTTP/1.1
    > Host: localhost:9234
    > Accept: */*
    > 
    < HTTP/1.1 500 Internal Server Error
    < Date: Sun, 05 Jul 2015 06:27:06 GMT
    < Content-Length: 66
    < Content-Type: text/plain; charset=utf-8

    {"path":"/proc/x","err":"stat /proc/x: no such file or directory"}

The contents of "err" are not guaranteed to be stable.

# Installing

Checking out the source code and running 'go build' should be sufficient to get
you a binary. There should also be a linux/amd64 binary courtesy of Travis CI
attached to each [release on Github](https://github.com/hfinucane/jsonproc/releases).

# Usage

    ./jsonproc -listen 10.9.8.7:9234

will get you a /proc-to-json gateway running on port 9234, listening on a local
address. In general, you should prefer to explicitly bind `jsonproc` to a
non-routable address- `/proc` leaks all sorts of information.

# Notes around the design

The goal is to quit writing one-off exfiltrations for things in `/proc`. So
`jsonproc` needs to be lightweight, safe, and reasonably performant.

By default, you can only read the first 4MB of files, and the first 1024
entries of a directory. These limitations are designed to make it more
difficult to DOS yourself.

`jsonproc` does nothing special around permissions. If you would like to read
files that are `root`-readable only, running this as `root` should work.
Because it's written in a memory-safe language, and never calls anything other
that `stat`, `open`, `read`, and `readdir`, it's possible this isn't even that
terrible an idea, but, uh, this should not be taken as an endorsement. Run it
as an unprivileged user if possible.

Relative paths- `../` & company- are unsupported. You should be limited to `/proc`.
