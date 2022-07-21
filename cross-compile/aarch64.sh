#!/bin/bash

sudo apt install gcc-aarch64-linux-gnu binutils-aarch64-linux-gnu # for armv8/aarch64
sudo apt install gcc-arm-linux-gnueabi binutils-arm-linux-gnueabi # for armv7/aarch32


# compile postgres for aarch64 (not working yet)
mkdir -p aarch64-postgres
cd aarch64-postgres
../postgres_cluster/configure --host=aarch64-linux-gnu --without-readline --without-zlib --prefix=$PWD/packaged
make
make install
cd ../

# Errors:
# /postgres_cluster/src/port/pg_crc32c_armv8_choose.c:88:20: error: ‘pg_comp_crc32c_armv8’ undeclared (first use in this function); did you mean ‘pg_comp_crc32c_sb8’?



# compile postgres for arm 32 bit (not working)

mkdir -p arm-postgres
cd arm-postgres
../postgres_cluster/configure CC=arm-linux-gnueabi-gcc CPP=arm-linux-gnueabi-cpp USE_DEV_URANDOM=1 --host=arm-linux --without-readline --without-zlib --disable-spinlocks
make
cd ../

# Errors:
# /postgres_cluster/src/include/pg_config.h:772:24: error: ‘__int128’ is not supported on this target
