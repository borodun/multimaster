## Compiling Postgres with multimaster to termux
Get [termux-packages](https://github.com/termux/termux-packages) repository:
```shell
git clone https://github.com/termux/termux-packages.git
```
Move _postgresql-mm_ package to _pacakges_:
```shell
cd termux-packages
cp -r ../multimaster/cross-compile/postgresql-mm packages/
```
Run _termux-package-builder_ container:
```shell
sudo ./scripts/run-docker.sh
```
Build package(choose needed architecture: aarch64(default), arm, i686, x86_64 or all):
```shell
./build-package.sh -a architecture postgresql-mm
```
Built package will be in _output_ folder \
Run in _termux_:
```shell
pkg install openssl libcrypt readline libandroid-shmem libuuid libxml2 libicu zlib
dpkg -i postgresql-mm_13.2-4_aarch64.deb
```

### Configuring 
You can configre package by changing variables in _postgresql-mm/build.sh_