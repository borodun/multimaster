CREATE EXTENSION multimaster;

SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=localhost',
    '{"dbname=mydb user=mtmuser host=10.10.10.200", "dbname=mydb user=mtmuser host=10.10.10.167"}');
