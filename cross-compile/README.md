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
./build-package.sh -a <architecture> -I postgresql-mm
```
Built package will be in _output_ folder \
Move it to termux device 
```shell
scp termux-packages/output/postgresql-mm_13.2-4_aarch64.deb termux:
```
Run in _termux_:
```shell
pkg install openssl libcrypt readline libandroid-shmem libuuid libxml2 libicu zlib
dpkg -i postgresql-mm_13.2-4_aarch64.deb
```

### Configuring 
You can configre build process by changing _[postgresql-mm/build.sh](postgresql-mm/build.sh)_. See [guide](https://github.com/termux/termux-packages/wiki/Building-packages).

## Creating your own repository
Clone _termux-apt-repo_ tool:
```shell
git clone https://github.com/termux/termux-apt-repo.git
```
Create input folder and put your packages in it:
```shell
cd termux-apt-repo
mkdir input
cp ../termux-packages/output/postgresql-mm* input/
termux-apt-repo input output
```
### Serve _output_ folder with nginx
```shell
sudo apt install nginx
```
Edit _/etc/nginx/sites-available/default_ to serve folder with your packages:
```conf
server {
	listen 80 default_server;
	listen [::]:80 default_server;
	root /home/user/termux-apt-repo/output;
	autoindex on;
	server_name example.com;

	location / {
		try_files $uri $uri/ =404;
	}
}
```
### Add your repository to termux
```shell
mkdir -p $PREFIX/etc/apt/sources.list.d
echo "deb [trusted=yes] http://example.com termux extras" > $PREFIX/etc/apt/sources.list.d/termux-extras.list
pkg update
```
Now you can install your custom packagases using _pkg install_
### Add my repository
```shell
mkdir -p $PREFIX/etc/apt/sources.list.d
echo "deb [trusted=yes] http://termux.borodun.works termux extras" > $PREFIX/etc/apt/sources.list.d/termux-extras-borodun.list
pkg update
```
Install Postgres with multimaster:
```shell
pkg install postgresql-mm
```