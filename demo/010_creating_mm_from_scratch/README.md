# Info
Deploy multimaster cluster from scratch

Logs are in _mm/node*/logfile_

## Usage
You need to have **LOCAL_IP** environment variable that will store IP accessible to future nodes:
```bash
export LOCAL_IP=192.168.31.144
```

Start scenario:
```bash
./start_scenario.sh
```

End scenrio:
```bash
./end_scenario.sh
```

Connecting to instance:
```bash
source conf.env
psql -h localhost -p $MM_PORT1 -d $MM_DB -U $MM_USER
```

Other utils:

```bash
# To start or restart everything
./restart_all.sh
# Stop instances
./poke_all.sh stop
# Start instances
./poke_all.sh start
# Clean up after dropping node
./clean_up.sh <node-id>
```

## By hand

0. Set environment vars:
```bash
source conf.env
```

1. Init 3 postgres instances:
```bash
mkdir -p mm/node1
mkdir -p mm/node2
mkdir -p mm/node3

initdb -D mm/node1
initdb -D mm/node2
initdb -D mm/node3
```

2. Configure instances:
```bash
echo -e $PG_CONF_LINES >> mm/node1/postgresql.conf
echo -e $PG_CONF_LINES >> mm/node2/postgresql.conf
echo -e $PG_CONF_LINES >> mm/node3/postgresql.conf

echo -e $PG_HBA_LINES >> mm/node1/pg_hba.conf
echo -e $PG_HBA_LINES >> mm/node2/pg_hba.conf
echo -e $PG_HBA_LINES >> mm/node3/pg_hba.conf
```

3. Start all instances:
```bash
pg_ctl -D mm/node1 -o "-p $MM_PORT1" -l mm/node1/logfile start
pg_ctl -D mm/node2 -o "-p $MM_PORT2" -l mm/node2/logfile start
pg_ctl -D mm/node3 -o "-p $MM_PORT3" -l mm/node3/logfile start
```

4. Create database for mm:
```bash
psql -h localhost -p $MM_PORT1 -d postgres -a -c "$CREATE_USER"
psql -h localhost -p $MM_PORT2 -d postgres -a -c "$CREATE_USER"
psql -h localhost -p $MM_PORT3 -d postgres -a -c "$CREATE_USER"

psql -h localhost -p $MM_PORT1 -d postgres -a -c "$CREATE_DB"
psql -h localhost -p $MM_PORT2 -d postgres -a -c "$CREATE_DB"
psql -h localhost -p $MM_PORT3 -d postgres -a -c "$CREATE_DB"
```

5. Init mm:
```bash
psql -h localhost -p $MM_PORT1 -d $MM_DB -a -c "$INIT_MM"
```
