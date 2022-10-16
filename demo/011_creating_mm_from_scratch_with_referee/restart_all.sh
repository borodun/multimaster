# Get env vars
source ./conf.env

# Stop and remove existing instances if they exist
./poke_all.sh stop
rm -rf mm

# Init instances
mkdir -p mm
cd mm
mkdir -p node1
mkdir -p node2
mkdir -p node3
cd ..

initdb -D ./mm/node1
initdb -D ./mm/node2
initdb -D ./mm/node3

cp $PG_CONF ./mm/node1/
cp $PG_CONF ./mm/node2/
cp $PG_CONF ./mm/node3/

cp $PG_HBA ./mm/node1/
cp $PG_HBA ./mm/node2/
cp $PG_HBA ./mm/node3/

# Make a referee
pg_ctl -D mm/node3 -o "-p $REF_PORT" -l mm/node3/logfile start
psql -h localhost -p $REF_PORT -d postgres -a -c "$CREATE_USER"
psql -h localhost -p $REF_PORT -d postgres -a -c "$CREATE_DB"
psql -h localhost -p $REF_PORT -d $MM_DB -a -c "$INIT_REF"

# Configure other instances for referee
echo $REF_CONNCONF >> mm/node1/postgresql.conf
echo $REF_CONNCONF >> mm/node2/postgresql.conf

# Start instances
pg_ctl -D mm/node1 -o "-p $MM_PORT1" -l mm/node1/logfile start
pg_ctl -D mm/node2 -o "-p $MM_PORT2" -l mm/node2/logfile start

# Create user and database for multimaster
for port in $MM_PORT1 $MM_PORT2
do
    psql -U $(whoami) -p $port -h localhost -d postgres -a -c "$CREATE_USER"
    psql -U $(whoami) -p $port -h localhost -d postgres -a -c "$CREATE_DB"
done


# Init cluster
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$INIT_MM"
