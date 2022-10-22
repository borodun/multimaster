chmod +x ./bootstrap.sh

export SERVER_ADDR="192.168.31.144:80"

export DEB_PATH="/deb/postgresql_13.2-4_aarch64.deb"
curl -o postgres_with_mm.deb ${SERVER_ADDR}${DEB_PATH}
pkg install ./postgres_with_mm.deb

export JOINER_PATH="/mtm-joiner"
curl -o mtm-joiner ${SERVER_ADDR}${JOINER_PATH}
chmod +x ./mtm-joiner

export ADD_PATH="/add-node.sh"
curl -o add-node.sh ${SERVER_ADDR}${ADD_PATH}
chmod +x ./add-node.sh

export DROP_PATH="/drop-node.sh"
curl -o drop-node.sh ${SERVER_ADDR}${DROP_PATH}
chmod +x ./drop-node.sh

export UNINSTALL_PATH="/uninstall.sh"
curl -o uninstall.sh ${SERVER_ADDR}${UNINSTALL_PATH}
chmod +x ./uninstall.sh

export CONNECT_PATH="/connect.sh"
curl -o connect.sh ${SERVER_ADDR}${CONNECT_PATH}
chmod +x ./connect.sh
