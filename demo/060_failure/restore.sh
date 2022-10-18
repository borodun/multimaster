pg_restore -U mtmuser -h $HOST1 -p $PORT1 -d $DBNAME ./dvdrental.tar
pg_restore -U mtmuser -h $HOST2 -p $PORT2 -d $DBNAME ./dvdrental.tar
pg_restore -U mtmuser -h $HOST3 -p $PORT3 -d $DBNAME ./dvdrental.tar
