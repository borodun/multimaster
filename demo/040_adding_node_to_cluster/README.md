# Info

Adding node to multimaster cluster

# Using _mtm-connector_ and _mtm-joiner_

## On PC

1. Run a multimaster cluster, see previous scenarios.

2. Build and run _[mtm-connector](../../mtm-connector)_
```bash
cd <scenario>
source conf.env
./mtm-connector -u "postgresql://${MM_USER}:${MM_PASSWORD}@${LOCAL_IP}:${MM_PORT1}/${MM_DB}?sslmode=disable"
```

Or run throught:
```bash
docker run -p 8080:8080 borodun/mtm-connector -u "postgresql://${MM_USER}:${MM_PASSWORD}@${LOCAL_IP}:${MM_PORT1}/${MM_DB}?sslmode=disable"
```

It will output on what address it is listening on. You'll need it for _mtm-joiner_.

## On Phone

1. Install Postgres with Multimaster

2. Build _[mtm-joiner](../../mtm-joiner)_ and copy it to Phone

3. Get your ip addr in local network:
```bash
ifconfig
export LOCAL_IP="192.168.31.166"
```

4. Export URL of _mtm-connector_
```bash
export CONNECTOR_URL="http://192.168.31.144:8080"
```
Note: protocol is mandatory

5. Add node with _mtm-joiner_:
```bash
./mtm-joiner -u $CONNECTOR_URL -a $LOCAL_IP
```
For additional arguments use **--help** flag

6. To remove node from cluster add **--drop** flag:
```bash
./mtm-joiner -u $CONNECTOR_URL -a $LOCAL_IP --drop
```

# By hand

[Mtm docs](https://postgrespro.github.io/mmts/#multimaster-adding-new-nodes-to-the-cluster)

1. Figure out the required connection string to access the new node.
 For example, for the database demo, user mtmuser, and the new node node4,
 the connection string can be "dbname=demo user=mtmuser host=node4". 

2. On one online node:
```sql
SELECT mtm.add_node('dbname=demo user=mtmuser host=node4');
```
It will return _id_ for new node

3. Go to the new node and clone all the data from one of the alive nodes to this node:
```bash
pg_basebackup -D datadir -h node1 -U mtmuser -c fast -v
```
It will output _write-ahead log end-point_

4. Configure postgres for recovery:
```bash
touch datadir/recovery.signal
export CONF_LINES="restore_command = 'false'
recovery_target = 'immediate'
recovery_target_action = 'promote'"
echo -e "$CONF_LINES" >> datadir/postgresql.conf
```

5. Start Postgres on new node:
```bash
pg_ctl -D datadir -l datadir/logfile start
```

6. On one online node:
```sql
SELECT mtm.join_node(<id>, '<end-point>');
```



