# Get env vars
source ./conf.env

# Stop and remove existing instances if they exist
./poke_all.sh stop
rm -rf mm

# Init instances
mkdir -p mm/node1 -m 0750
mkdir -p mm/node2 -m 0750
mkdir -p mm/node3 -m 0750

cp -r $BACKUP_FOLDER/* mm/node1/
cp -r $BACKUP_FOLDER/* mm/node2/
cp -r $BACKUP_FOLDER/* mm/node3/

echo -e $PG_CONF_LINES >> mm/node1/postgresql.conf
echo -e $PG_CONF_LINES >> mm/node2/postgresql.conf
echo -e $PG_CONF_LINES >> mm/node3/postgresql.conf

echo -e $PG_HBA_LINES >> mm/node1/pg_hba.conf
echo -e $PG_HBA_LINES >> mm/node2/pg_hba.conf
echo -e $PG_HBA_LINES >> mm/node3/pg_hba.conf

# Start instances
./poke_all.sh start

# Init cluster
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$INIT_MM"

