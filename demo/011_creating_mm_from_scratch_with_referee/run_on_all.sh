if [ "$#" -lt 1 ]; then
    echo "Usage: <command-to-run-all-instances>"
    exit
fi

source ./conf.env

for port in $MM_PORT1 $MM_PORT2
do
    psql -U $MM_USER -p $port -h localhost -d $MM_DB -c "$1"
done
