if [[ -z "${LOCAL_IP}" ]]; then
  echo "LOCAL_IP is not set, read Usage in README"
  exit
fi

# Get env vars
source ./conf.env

# Stop and remove existing instances if they exist
./poke_all.sh stop
rm -rf mm

# Init instances
mkdir -p mm/node1
mkdir -p mm/node2
mkdir -p mm/node3

initdb -D ./mm/node1
initdb -D ./mm/node2
initdb -D ./mm/node3

#echo "wal_level = logical" >> mm/node1/postgresql.conf

echo -e $PG_CONF_LINES >> mm/node1/postgresql.conf
echo -e $PG_CONF_LINES >> mm/node2/postgresql.conf
echo -e $PG_CONF_LINES >> mm/node3/postgresql.conf

echo -e $PG_HBA_LINES >> mm/node1/pg_hba.conf
echo -e $PG_HBA_LINES >> mm/node2/pg_hba.conf
echo -e $PG_HBA_LINES >> mm/node3/pg_hba.conf

# Start instances
./poke_all.sh start

# Create user and database for multimaster
for port in $MM_PORT1 $MM_PORT2 $MM_PORT3
do
    psql -U $(whoami) -p $port -h localhost -d postgres -a -c "$CREATE_USER"
    psql -U $(whoami) -p $port -h localhost -d postgres -a -c "$CREATE_DB"
done

# Start logical replication
for port in $MM_PORT1 $MM_PORT2 $MM_PORT3
do
    psql -U $MM_USER -p $port -h localhost -d $MM_DB -a -c "$CREATE_DATA"
done

psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$CREATE_PUB"
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$FILL_DATA"

psql -U $MM_USER -p $MM_PORT2 -h localhost -d $MM_DB -a -c "$CREATE_SUB1"
psql -U $MM_USER -p $MM_PORT3 -h localhost -d $MM_DB -a -c "$CREATE_SUB2"

sleep 1

# Drop replication
psql -U $MM_USER -p $MM_PORT2 -h localhost -d $MM_DB -a -c "$DROP_SUB1"
psql -U $MM_USER -p $MM_PORT3 -h localhost -d $MM_DB -a -c "$DROP_SUB2"
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$DROP_PUB"

# Init cluster
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$INIT_MM"
