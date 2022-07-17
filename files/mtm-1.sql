CREATE EXTENSION multimaster;

SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=node3',
    '{"dbname=mydb user=mtmuser host=node4", "dbname=mydb user=mtmuser host=node5"}');
