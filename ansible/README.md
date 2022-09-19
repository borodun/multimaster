# multimaster
Postgres multimaster \
Deploy modified postgres on nodes and start multimaster clusters on them + anable monitoring:
```shell
ansible-playbook site.yml -i hosts -f 10 -t prepare,build,install,clusters,monitoring,checkmm
```
Note: change password for _postgres_ user in site.yml 

Start/stop/status postgres
```shel
ansible-playbook site.yml -i hosts -f 6 -t start/stop/status
```

Other tags
```shel
ansible-playbook site.yml --list-tags
```

Check multimasters on node0:
```shell
psql -U mtmuser -d mydb
```
```sql
SELECT * FROM mtm.status();
SELECT * FROM mtm.nodes();
```
Or from ansible
```shell
ansible-playbook site.yml -i hosts -f 6 -t checkmm
```
Example output:
```
 my_node_id | status | connected | gen_num | gen_members | gen_members_online | gen_configured 
------------+--------+-----------+---------+-------------+--------------------+----------------
          1 | online | {1,2,3}   |       1 | {1,2,3}     | {1,2,3}            | {1,2,3}
```
```
 id |              conninfo               | is_self | enabled | connected | sender_pid | receiver_pid | n_workers | receiver_mode 
----+-------------------------------------+---------+---------+-----------+------------+--------------+-----------+---------------
  1 | dbname=mydb user=mtmuser host=node3 | t       | t       | t         |            |              |           | 
  2 | dbname=mydb user=mtmuser host=node4 | f       | t       | t         |      84498 |        84508 |         0 | normal
  3 | dbname=mydb user=mtmuser host=node5 | f       | t       | t         |      84502 |        84509 |         0 | normal
```
