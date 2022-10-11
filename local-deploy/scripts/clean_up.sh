INIT_DONE="DELETE FROM mtm.nodes_init_done WHERE id = 4"
SYNCPOINTS="DELETE FROM mtm.syncpoints WHERE receiver_node_id = 4 OR origin_node_id = 4"

for i in 5411 5412 5413
do
    psql -U mtmuser -p $i -h localhost -d mydb -c "${INIT_DONE}"
done

psql -U mtmuser -p 5411 -h localhost -d mydb -c "${SYNCPOINTS}"
