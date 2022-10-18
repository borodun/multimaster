## 6. Эмуляция отказа

### Поднятие кластера

0. убеждаемся, что собран контейнер demo-mtm-postgres
1. редактируем ./env.sh по вкусу:
  - connstring для каждой ноды
  - порты для каждой ноды
2. docker compose -f mtm-compose.yml up -d
3. ./restore.sh && ./init_mmts.sh
...
4. docker compose -f mtm-compose.yml down

### Общий сценарий проверки

0. Жертва - всегда узел с тем портом, на который идет нагрузка
1. `cat ./load.sh`
Пару слов о характере нагрузки
2. `./load.sh <port>`
3. *ломаем (см. ниже)*
4. Показываем, что работа не нарушилась - подключимся  
`psql <port-reserve>`  
сделаем несколько запросов на запись и на чтение
`select`
`insert`
`select`
5. `./load.sh <port-reserve>`
6. *чиним (см. ниже)* 
7. `./load.sh <port>`

### 1. Выключение ноды
л. `docker stop 060_failure-mtm-1-1`
ч. `docker start 060_failure-mtm-1-1`

### 2. Полный обрыв связи
л. `docker network disconnect mtmnet 060_failure-mtm-1-1`
ч. `docker network connect mtmnet 060_failure-mtm-1-1`

### 3. Частичный обрыв связи (A-B-C топология)

` $ iptables -A OUTPUT -s 10.11.0.12 -j REJECT `

*здесь было убита куча времени на docker networks, прежде чем до меня дошло что ip один а не 2*
можно попробовать что-то с фаерволом, но highly unlikely
