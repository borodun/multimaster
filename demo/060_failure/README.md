## 6. Эмуляция отказа

### Поднятие кластера

0. убеждаемся, что собран контейнер demo-mtm-postgres
1. редактируем ./env.sh по вкусу:
  - connstring для каждой ноды
  - порты для каждой ноды
2.
3. docker compose -f mtm-compose.yml up -d
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

### 2. Полный обрыв связи
л. docker network disconnect mtm-netw mtm-1
ч. docker network connect mtm-new mtm-1

### 3. Частичный обрыв связи (A-B-C топология)

Оставим на самый конец
