# Info

Dropping node from multimaster cluster

# Using _mtm-connector_ and _mtm-joiner_

To drop node with that _mtm-connector_ node should be added _mtm-joiner_

To add node, see **scenario 040**.

## On Phone

1. Get your ip addr in local network:
```bash
ifconfig
export LOCAL_IP="192.168.31.166"
```

2. Export URL of _mtm-connector_
```bash
export CONNECTOR_URL="http://192.168.31.144:8080"
```
Note: protocol is mandatory

3. To remove node from cluster run _mtm-joiner_ with **--drop** flag:
```bash
./mtm-joiner -u $CONNECTOR_URL -a $LOCAL_IP --drop
```

### Clean up

In all scenarios (except 011) ther is a script that will clean up after dropping node. 
Just pass it and dropped node _id_ as 1st agument:
```bash
cd <scenario>
./clean_up.sh <id>
```

# By hand

[Mtm docs](https://postgrespro.github.io/mmts/#multimaster-removing-nodes-from-the-cluster)

1. On one online node:
```sql
SELECT mtm.drop_node(<id>);
```
It will return _id_ for new node

### Clean up

For now after dropping node you need clean two tables: _mtm.nodes_init_done_ and _mtm.syncpoints_

1. On all online nodes:
```sql
DELETE FROM mtm.nodes_init_done WHERE id = <id>;
```

2. On one online node:
```sql
DELETE FROM mtm.syncpoints WHERE receiver_node_id = <id> OR origin_node_id = <id>;
```

