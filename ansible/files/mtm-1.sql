CREATE EXTENSION multimaster;

SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=node1',
    '{"dbname=mydb user=mtmuser host=node2", "dbname=mydb user=mtmuser host=node3"}');
