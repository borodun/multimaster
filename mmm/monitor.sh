#!/bin/bash

USER=$1
DATABASE=$2
HOSTS=$(cat hosts)

MM_STATUS_QUERY="SELECT status, gen_num FROM mtm.status()"
PG_STAT_ACTIVITY_QUERY="SELECT count(*) FROM pg_stat_activity WHERE datname = '$DATABASE'"
XACT_ID_QUERY="SELECT pg_current_xact_id()"
DATABASE_SIZE_QUERY="SELECT pg_size_pretty(pg_database_size('$DATABASE'))"
DATABASE_STAT_QUERY="SELECT xact_commit, xact_rollback, tup_returned, tup_inserted, tup_updated, tup_deleted FROM pg_stat_database WHERE datname='$DATABASE'"

nodes_info_list="host,status,gen_num,db_processes,xact_id,database_size,xact_commit,xact_rollback,tup_returned,tup_inserted,tup_updated,tup_deleted\n"
online_nodes=""

run_query() {
    host=$(echo $1 | cut -d':' -f1)
    port=$(echo $1 | cut -d':' -f2)
    query="$2"
    out=$(psql -U $USER -d $DATABASE -h $host -p $port -c "$2" -t -A -F "," 2> log)
    retval=$?
    echo $out
    return $retval
}

# Get info about multimaster nodes 
for host in $HOSTS
do  
    node_status=$(run_query $host "$MM_STATUS_QUERY")
    retval=$?
    node_info=$host

    # If return code is 0 than OK
    if [ $retval -eq 0 ]; then
        node_info+=",$node_status"

        if [ $(echo $node_status | cut -d',' -f1) == "online" ]; then
            online_nodes+="$host \n"         
        fi
    # Return code 2 is connection error
    elif [ $retval -eq 2 ]; then
        node_info+=",unreachable"
    elif [[ $(tail -n 1 log) == *"current status isolated"* ]]; then
        node_info+=",isolated"
    else
        node_info+=",error"
    fi

    db_processes=$(run_query $host "$PG_STAT_ACTIVITY_QUERY")
    node_info+=",$db_processes"

    current_xact_id=$(run_query $host "$XACT_ID_QUERY")
    node_info+=",$current_xact_id"

    database_size=$(run_query $host "$DATABASE_SIZE_QUERY")
    node_info+=",$database_size"

    database_stat=$(run_query $host "$DATABASE_STAT_QUERY")
    node_info+=",$database_stat"

    nodes_info_list+="$node_info \n"
done

echo -e $nodes_info_list | column -t -s ',' 
echo -e $online_nodes > online_nodes

