# Info
Deploy multimaster cluster from existing database backup

Logs are in _mm/node*/logfile_

## Usage
You need to have some database backup to start from. You can make it using [_scenario 020_](../020_creating_mm_from_existing_db_dump/README.md).

```bash
cp -r ../020_creating_mm_from_existing_db_dump/backup/* ./databases/demo-small/
```

You need to have **LOCAL_IP** environment variable that will store IP accessible to future nodes:
```bash
export LOCAL_IP=192.168.31.144
```

Also change **BACKUP_FOLDER** in _conf.env_ if you want to use your own database.

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
mkdir -p mm/node1 -m 0750
mkdir -p mm/node2 -m 0750
mkdir -p mm/node3 -m 0750

cp -r $BACKUP_FOLDER/* mm/node1/
cp -r $BACKUP_FOLDER/* mm/node2/
cp -r $BACKUP_FOLDER/* mm/node3/
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

4. Init mm:
```bash
psql -h localhost -p $MM_PORT1 -d $MM_DB -a -c "$INIT_MM"
```
