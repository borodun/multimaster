if [[ -z "${1}" ]]; then
  echo "No argument, use start or stop"
  exit
fi


pg_ctl -D ./mm1/node11 -o "-p 5411" -l ./mm1/node11/logfile $1
pg_ctl -D ./mm1/node12 -o "-p 5412" -l ./mm1/node12/logfile $1
pg_ctl -D ./mm1/node13 -o "-p 5413" -l ./mm1/node13/logfile $1
