# Setup

I'm developing this for my Raspberry Pi Zero W running DietPi, connecting to my
BS410 scale. You should need nothing after compiling this (don't compile this
on the Pi!)

# Compilation

All dependencies are vendored in, so you should be able to just `go build` for
your local platform. The `Makefile` supports these platforms:
  arm6 arm7 linux32 linux64

# TODO

Add incoming data to the csvs
Detect when to write the new csvs to disk
