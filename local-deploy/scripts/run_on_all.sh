if [ "$#" -lt 1 ]; then
    echo "Usage: <command-to-run-all-instances>"
    exit 1
fi

echo "${1}"

for i in 5411 5412 5413
do
    psql -U mtmuser -p $i -h localhost -d mydb -c "${1}"
done
