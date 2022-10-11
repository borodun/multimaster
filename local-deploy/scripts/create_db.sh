SQL='resources/bootstrap.sql'

for i in 5411 5412 5413
do
    psql -U $(whoami) -p $i -h localhost -d postgres -a -f $SQL
done
