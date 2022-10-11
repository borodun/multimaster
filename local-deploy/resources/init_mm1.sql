CREATE EXTENSION multimaster;
SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=192.168.31.144 port=5411 password=1234', 
    '{"dbname=mydb user=mtmuser host=192.168.31.144 port=5412 password=1234", 
    "dbname=mydb user=mtmuser host=192.168.31.144 port=5413 password=1234"}');
