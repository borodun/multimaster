# multimaster
Postgres multimaster \
Deploy modified postgres and configure everything
```shell
ansible-playbook site.yml -i hosts  --ask-become-pass -t prepare,build,install
```
BECOME password is password for superuser on hosts

On one node
```shell
psql
```
```sql
CREATE EXTENSION multimaster;
SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=node0',
    '{"dbname=mydb user=mtmuser host=node1", "dbname=mydb user=mtmuser host=node2"}');
```
Check
```sql
SELECT * FROM mtm.status();
SELECT * FROM mtm.nodes();
```
