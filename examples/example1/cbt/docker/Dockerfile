# Start from an ubuntu image and add third-party deps as needed
#
FROM ubuntu

# Install build neccessities
# The image will only contain the build tool environment, we will use
# different types of volumes to store data of more persistent nature:
# /build is where we build in.
# /src is where the source code is mapped into.

RUN apt-get update && apt-get install -y wget autoconf automake \
    libtool make g++ unzip build-essential \
    pkg-config gdbserver mingw-w64

WORKDIR /build

COPY ./clean.sh /
CMD ['ls']