# Info
Deploy multimaster cluster from master with 2 streaming replicas

Logs are in _mm/node*/logfile_

## Usage
You need to change **LOCAL_IP** in _conf.env_

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

1. Init master instance:
```bash
mkdir -p mm/node1
initdb -D ./mm/node1
```

2. Configure master:
```bash
echo -e $PG_CONF_LINES >> mm/node1/postgresql.conf
echo -e $PG_HBA_LINES >> mm/node1/pg_hba.conf
```

3. Start master:
```bash
pg_ctl -D mm/node1 -o "-p $MM_PORT1" -l mm/node1/logfile start
```

4. Create database for mm:
```bash
psql -h localhost -p $MM_PORT1 -d postgres -a -c "$CREATE_USER"
psql -h localhost -p $MM_PORT1 -d postgres -a -c "$CREATE_DB"
```

#### Start streaming replication

1. Create standby nodes:
```bash
mkdir -p mm/node2 -m 0750
mkdir -p mm/node3 -m 0750
pg_basebackup -D mm/node2 -h localhost -p $MM_PORT1 -U $MM_USER -P -Xs -R
pg_basebackup -D mm/node3 -h localhost -p $MM_PORT1 -U $MM_USER -P -Xs -R
```

2. Start standby nodes:
```bash
pg_ctl -D ./mm/node2 -o "-p $MM_PORT2" -l ./mm/node2/logfile start
pg_ctl -D ./mm/node3 -o "-p $MM_PORT3" -l ./mm/node3/logfile start
```

3. Fill db with some data for testing:
```bash
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$CREATE_DATA"
psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -a -c "$FILL_DATA"
```

#### Stop streaming replication

1. Stop standby nodes:
```bash
pg_ctl -D ./mm/node2 -o "-p $MM_PORT2" -l ./mm/node2/logfile stop
pg_ctl -D ./mm/node3 -o "-p $MM_PORT3" -l ./mm/node3/logfile stop
```

2. Remove stanby signal:
```bash
rm mm/node2/standby.signal
rm mm/node3/standby.signal
```

3. Start nodes:
```bash
pg_ctl -D ./mm/node2 -o "-p $MM_PORT2" -l ./mm/node2/logfile start
pg_ctl -D ./mm/node3 -o "-p $MM_PORT3" -l ./mm/node3/logfile start
```

#### Start multimaster cluster from master and replicas

1. Init mm:
```bash
psql -h localhost -p $MM_PORT1 -d $MM_DB -a -c "$INIT_MM"
```
