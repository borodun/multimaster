source ./conf.env
./end_scenario.sh
mkdir -p mm/node1
initdb -D mm/node1
#cp $PG_CONF mm/node1
#cp $PG_HBA mm/node1
pg_ctl -D mm/node1 -o "-p $MM_PORT1" -l mm/node1/logfile start
psql -h localhost -p $MM_PORT1 -d postgres -a -c "$CREATE_USER"
psql -h localhost -p $MM_PORT1 -d postgres -a -c "$CREATE_DB"
psql -U $MM_USER -p $MM_PORT1 -h localhost -d postgres -b -f $DUMP_FILE

rm -r backup
pg_basebackup -D ./backup -h localhost -p $MM_PORT1 -U $MM_USER -P -c fast

./end_scenario.sh
