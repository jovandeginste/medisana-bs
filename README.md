This project is heavily inspired by [this Python project](https://github.com/keptenkurk/BS440).
Since I am not a Python developer and can't get used to the language, and Python and all the
dependencies for that project are quite heavy for a Pi, I decided to give a Go-port a try.

I use their csv format, so the output should be compatible between the code bases.

# Current status (Update: 2017/07/14)

The code compiles (:+1:) and the runs stable. I haven't run it long enough to decide that the
code is completely stable.

Earlier attempts resulted in the daemon locking up (probably trying
to scan using the device, but hanging because BlueZ was still around and locking). I disabled
BlueZ completely on my pi.

Last run I missed some data, probably because I set the scan duration too long. I'm now trying
with short scans (10s). Which should result in fast retries.

# Setup

I'm developing this for my Raspberry Pi Zero W running DietPi, connecting to my
BS410 scale. You should need nothing after compiling this (don't compile this
on the Pi!)

# Compilation

Copy `config.go.example` to `config.go` and change the parameters.

All dependencies are vendored in, so you should be able to just `go build` for
your local platform. The `Makefile` supports these platforms:

> arm6 arm7 linux32 linux64

# TODO

* Add more plugins
* Test it for a significant period
* Make some noise about it
