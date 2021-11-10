#!/bin/bash

# Download zeromq
# Ref http://zeromq.org/intro:get-the-software
wget https://github.com/zeromq/libzmq/releases/download/v4.3.4/zeromq-4.3.4.tar.gz
# Unpack tarball package
tar xvzf zeromq-4.3.4.tar.gz

# Install dependency
apt-get update && \
apt-get install -y libtool pkg-config build-essential autoconf automake uuid-dev apt-get libzmq3-dev mingw-w64



########
update-alternatives --set x86_64-w64-mingw32-gcc /usr/bin/x86_64-w64-mingw32-gcc-posix
update-alternatives --set x86_64-w64-mingw32-g++ /usr/bin/x86_64-w64-mingw32-g++-posix


# Install libsodium
# git clone git://github.com/jedisct1/libsodium.git
# cd libsodium
# ./autogen.sh
# ./configure && make check
# make install
# ldconfig
# cd ..

# cd zeromq-4.3.4/builds/mingw32
# wget https://github.com/oneclick/rubyinstaller/archive/refs/tags/devkit-4.7.2.tar.gz
# tar xvzf devkit-4.7.2.tar.gz
# cp ./rubyinstaller-devkit-4.7.2/resources/devkit/* .
# rm -rf devkit-*
# chmod +x ./devkitvars.bat
# ./devkitvars.bat

# make all -f Makefile.mingw32

# cd ../../..


# Create make file
cd zeromq-4.3.4
./autogen.sh
./configure


# CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CFLAGS="-I/usr/local/include" LDFLAGS="-L/usr/local/lib -lgcc -lgcc_s" ./configure --host=x86_64-w64-mingw32
 ## CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CFLAGS="-I/usr/local/include" LDFLAGS="-L/usr/local/lib -lgcc -lgcc_s" ./configure --host=x86_64-w64-mingw32


# ./configure CC=x86_64-w64-mingw32-gcc --disable-dependency-tracking --host=arm-apple-darwin10



make

# Build and install(root permission only)
make install

# Install zeromq driver on linux
ldconfig

# Check installed
ldconfig -p | grep zmq

#go get -v -x github.com/pebbe/zmq4



# Expected
############################################################
# libzmq.so.5 (libc6,x86-64) => /usr/local/lib/libzmq.so.5
# libzmq.so (libc6,x86-64) => /usr/local/lib/libzmq.so
############################################################
