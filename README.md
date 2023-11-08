[![Go](https://github.com/blablatov/restduples/actions/workflows/server.yml/badge.svg)](https://github.com/blablatov/restduples/actions/workflows/server.yml)
### Тестовый REST API на Go  
Согласно [ТЗ:](https://gist.github.com/zemlya25/585ab3fb3b0704880f920728c7598beb)

### Описание   
Принимает два `user_id` и выдаёт ответ - являются ли они дублем или нет. Дублем считается пара `user_id`, у которых хотя бы два раза совпадает `ip-адрес` в логе соединений.
Лог соединений будет хранится в тестовой БД `duples` СУБД `ClickHouse` ОП `Yandex Cloud`, столбцовой `СУБД` для онлайн обработки аналитических запросов `OLAP`.    
Создается таблица, в нее загружаются тестовые данные. Структура таблицы:     
  
	`CREATE TABLE conn_log
	(
    user_id UInt32,
    ip_addr IPv4,
    timestamp DateTime,
	)
	ENGINE = MergeTree
	PRIMARY KEY (user_id)`
	

### Заполнение таблицы данными из SQL запроса:       
	go run insertlog.go  

### Генератор файла со строками записей таблицы (conn_log.csv).  
Ключ `1000000`, это число итераций основного цикла.  
Общее количество строк, определяется простым подбором ключа:  

	go build gendtables.go  
	gendtables 1000000
		
### Заполнение таблицы посредством сгенеренного .cvs файла:
	clickhouse-client -q "INSERT INTO conn_log FORMAT CSV" < conn_log.csv 
	

### Сборка локально и в Yandex Cloud:  
#### Локально. Local:  
	docker build -t restduples -f Dockerfile  
	
#### Облако. Cloud:    
	sudo docker build . -t cr.yandex/${REGISTRY_ID}/debian:restduples -f Dockerfile


### Тестирование локально и в Yandex Cloud:         
#### Локально. Local:    
	go test -v client_test.go    
	go test -v selectlog_test.go  

#### Облако. Cloud:   
	sudo docker run --name restduples -p 12345:12345 -d cr.yandex/${REGISTRY_ID}/debian:restduples 
	go test -v client_test.go  	

### Использование. How use:   
	go run restduples.go
	go run client.go    
	
### Ответ сервера, вывод модуля. Response of server and module:     
	Status = 200 OK 2023/10/17 10:49:28 
	This is duples:
 	{"true"} or {"false"}  



	


  




 
