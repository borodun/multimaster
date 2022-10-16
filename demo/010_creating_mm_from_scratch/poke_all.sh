if [[ -z "${1}" ]]; then
  echo "No argument, use start, stop or status"
  exit
fi

source ./conf.env

pg_ctl -D ./mm/node1 -o "-p $MM_PORT1" -l ./mm/node1/logfile $1
pg_ctl -D ./mm/node2 -o "-p $MM_PORT2" -l ./mm/node2/logfile $1
pg_ctl -D ./mm/node3 -o "-p $MM_PORT3" -l ./mm/node3/logfile $1
