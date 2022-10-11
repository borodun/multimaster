CONF=./resources/postgresql.conf
HBA=./resources/pg_hba.conf

mkdir -p mm1
cd mm1
mkdir -p node11
mkdir -p node12
mkdir -p node13
cd ..

initdb -D ./mm1/node11
initdb -D ./mm1/node12
initdb -D ./mm1/node13

cp $CONF ./mm1/node11/
cp $CONF ./mm1/node12/
cp $CONF ./mm1/node13/

cp $HBA ./mm1/node11/
cp $HBA ./mm1/node12/
cp $HBA ./mm1/node13/
