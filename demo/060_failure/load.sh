#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: <host> <port>"
    exit 1
fi

i=0

while true
do
    echo $i;
    COMMAND="begin; insert into category (name) values ('load_$i'); commit;"
    psql -U mtmuser -p $2 -h $1 -d $DBNAME -c "$COMMAND"
    ((i++))
    sleep 0.1
done
