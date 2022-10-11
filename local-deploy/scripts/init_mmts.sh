SQL1='./resources/init_mm1.sql'
MASTERPORT1=5411

psql -U mtmuser -p $MASTERPORT1 -h localhost -d mydb -a -f $SQL1
