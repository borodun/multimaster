## 6. Эмуляция отказа

### Поднятие кластера

0. убеждаемся, что собран контейнер demo-mtm-postgres
1. редактируем ./env.sh по вкусу:
  - connstring для каждой ноды
  - порты для каждой ноды
2. docker compose up -d
3. ./scripts/restore.sh && ./scripts/init_mmts.sh
...
4. docker compose down

### Общий сценарий проверки

0. Жертва - всегда узел с тем портом, на который идет нагрузка
1. `cat ./load.sh`
Пару слов о характере нагрузки
2. `./load.sh <port>`
3. *ломаем (см. ниже)*
4. Показываем, что работа не нарушилась - подключимся  `$ psql -U mtmuser -p $PORT1 -h $HOST1 -d $DBNAME -a -c "select * from mtm.status();"`
сделаем несколько запросов на запись с включенным `watch`. 
5. `./load.sh <port-reserve>`  
6. *чиним (см. ниже)*   
7. `./load.sh <port>`  

### 1. Выключение ноды

Подготовка  

    $ export NODE=060_failure-mtm-1-1

Ломаем 

    $ docker stop $NODE

Чиним  
    
    $ docker start $NODE

### 2. Полный обрыв связи

Подготовка    
    
    $ export NODE1=060_failure-mtm-1-1

Ломаем   
    
    $ docker exec -u 0 -it $NODE1 iptables -A OUTPUT -d $HOST2 -j REJECT
    $ docker exec -u 0 -it $NODE1 iptables -A OUTPUT -d $HOST3 -j REJECT

Чиним  
    
    $ docker exec -u 0 -it $NODE1 iptables -D OUTPUT -d $HOST2 -j REJECT
    $ docker exec -u 0 -it $NODE1 iptables -D OUTPUT -d $HOST3 -j REJECT

### 3. Частичный обрыв связи (A-B-C топология)

Подготовка    

    $ export NODE1=060_failure-mtm-1-1
    $ export NODE2=060_failure-mtm-2-1

Ломаем   

    $ docker exec -u 0 -it $NODE1 iptables -A OUTPUT -d $HOST2 -j REJECT
    $ docker exec -u 0 -it $NODE2 iptables -A OUTPUT -d $HOST1 -j REJECT
Чиним  

    $ docker exec -u 0 -it $NODE1 iptables -D OUTPUT -d $HOST2 -j REJECT
    $ docker exec -u 0 -it $NODE2 iptables -D OUTPUT -d $HOST1 -j REJECT

