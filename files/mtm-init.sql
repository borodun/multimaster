CREATE EXTENSION multimaster;

--- 1st try
SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=node0',
    '{"dbname=mydb user=mtmuser host=node1", "dbname=mydb user=mtmuser host=node2"}');

--- 2nd try
SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=localhost',
    '{"dbname=mydb user=mtmuser host=10.10.10.200", "dbname=mydb user=mtmuser host=10.10.10.167"}');
