# Info

Bootstrap for adding multimaster node on phone.

# Hosting .deb package
You need to build Postgres with Multimaster for Termux, see [how](../../cross-compile).

After building, put it in _files/deb_. Run nginx container and mount path to .deb package 
to /usr/share/nginx/html in container:
```bash
docker run -p 80:80 \
-v $(pwd)/nginx.conf:/etc/nginx/conf.d/default.conf \
-v $(pwd)/files:/files \
nginx:stable-alpine
```

# On phone

Save server addr:
```bash
export SERVER_ADDR="192.168.31.144:80"
```

Download bootstrap script.
```bash
curl -o bootstrap.sh ${SERVER_ADDR}/bootstrap.sh
```

Run bootstrap script:
```bash
chmod +x ./bootstrap.sh
./bootstrap.sh
```

If phone connected to several networks then get and save phone's ip address:
```bash
ifconfig
export LOCAL_IP=192.168.31.166
```

Add node:
```bash
./add-node.sh
```

Drop node:
```bash
./drop-node.sh
```

Connect to db:
```bash
./connect.sh
```

Uninstall everything:
```bash
./uninstall.sh
```

# For QR

```bash
curl -o bootstrap.sh 192.168.31.144/bootstrap.sh;
chmod +x bootstrap.sh;
./bootstrap.sh
```
