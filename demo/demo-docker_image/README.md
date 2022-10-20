## Usage:

`$ git clone -b rel13_mm_2 --single-branch https://github.com/postgrespro/postgres_cluster.git`  
`$ git clone -b master --single-branch https://github.com/postgrespro/mmts.git ./postgres_cluster/contrib/mmts/`
`$ docker build -t demo-mtm-postgres .`    
cry because buid failed on `error: 'IDLE_SESSION_TIMEOUT' undeclared`  
go fix and retry  

then run docker-compose of choice

Partially taken from https://github.com/kernogo/mtm-operator
