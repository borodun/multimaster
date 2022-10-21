export INIT_SQL="CREATE EXTENSION multimaster; SELECT mtm.init_cluster('dbname=$DBNAME user=mtmuser host=$HOST1 port=$PORT1', '{\"dbname=$DBNAME user=mtmuser host=$HOST2 port=$PORT2\", \"dbname=$DBNAME user=mtmuser host=$HOST3 port=$PORT3\"}');"

echo $INIT_SQL

psql -U mtmuser -p $PORT1 -h $HOST1 -d $DBNAME -a -c "$INIT_SQL"

echo "mmts cluster started, use mtm.status() to check"
