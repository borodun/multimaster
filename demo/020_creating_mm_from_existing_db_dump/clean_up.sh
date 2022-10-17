source ./conf.env

for i in $MM_PORT1 $MM_PORT2 $MM_PORT3
do
    psql -U $MM_USER -p $i -h localhost -d $MM_DB -c "$CLEAR_INIT_DONE"
done

psql -U $MM_USER -p $MM_PORT1 -h localhost -d $MM_DB -c "$CLEAR_SYNCPOINTS"
