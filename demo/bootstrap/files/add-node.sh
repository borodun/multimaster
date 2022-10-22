#if [[ -z "${LOCAL_IP}" ]]; then
#  echo "LOCAL_IP env is not set, use ifconfig to get it"
#  exit
#fi

export CONNECTOR_URL="http://192.168.31.144:8080"

./mtm-joiner -u "$CONNECTOR_URL" -a "$LOCAL_IP" $1
