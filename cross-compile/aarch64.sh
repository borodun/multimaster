#!/bin/bash

# Install Android NDK
curl https://dl.google.com/android/repository/android-ndk-r25-linux.zip -O
unzip android-ndk-r25-linux.zip
echo "export PATH=$PATH:$HOME/android-ndk-r25/toolchains/llvm/prebuilt/linux-x86_64/bin" >> ~/.bashrc
source ~/.bashrc

# Install compiler for aarch64 (not needed)
sudo apt install gcc-aarch64-linux-gnu binutils-aarch64-linux-gnu # for armv8/aarch64

# compile postgres for aarch64 (not working yet)
mkdir -p aarch64-postgres
cd aarch64-postgres
../postgres_cluster/configure --host=aarch64-linux-android28 --without-readline --without-zlib --prefix=$PWD/postgres-aarch64 CC=aarch64-linux-android28-clang USE_DEV_URANDOM=1
make
make install
tar -czf postgres-aarch64.tar.gz postgres-aarch64/ --hard-dereferenc
cd ../

# Errors:



