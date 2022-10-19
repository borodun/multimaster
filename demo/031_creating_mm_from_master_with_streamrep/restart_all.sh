# Get env vars
source ./conf.env

# Stop and remove existing instances if they exist
./poke_all.sh stop
rm -rf mm

# Init master
mkdir -p mm/node1
initdb -D ./mm/node1
echo -e $PG_CONF_LINES >> mm/node1/postgresql.conf
echo -e $PG_HBA_LINES >> mm/node1/pg_hba.conf

# Start master
pg_ctl -D ./mm/node1 -o "-p $MM_PORT1" -l ./mm/node1/logfile start

# Create user and db on master
psql -U $(whoami) -p $MM_PORT1 -h localhost -d postgres -a -c "$CREATE_USER"
psql -U $(whoami) -p $MM_PORT1 -h localhost -d postgres -a -c "$CREATE_DB"

# Creating standby nodes
mkdir -p mm/node2 -m 0750
mkdir -p mm/node3 -m 0750
pg_basebackup -D mm/node2 -h localhost -p $MM_PORT1 -U $MM_USER -P -Xs -R
pg_basebackup -D mm/node3 -h localhost -p $MM_PORT1 -U $MM_USER -P -Xs -R

# Start standby nodes
pg_ctl -D ./mm/node2 -o "-p $MM_PORT2" -l ./mm/node2/logfile start
pg_ctl -D ./mm/node3 -o "-p $MM_PORT3" -l ./mm/node3/logfile start

# Create some data for testing
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$CREATE_DATA"
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$FILL_DATA"

sleep 1

# Stop streaming replication
./poke_all.sh stop
rm mm/node2/standby.signal
rm mm/node3/standby.signal
./poke_all.sh start

# Init cluster
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$INIT_MM"
