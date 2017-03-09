#!/bin/bash
set -e
pushd /src
rm -rf missing Makefile.in install-sh hello.exe depcomp configure config.h.in compile aclocal.m4 src/Makefile.in src/.deps autom4te.cache
popd

