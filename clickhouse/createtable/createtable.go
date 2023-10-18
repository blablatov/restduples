// Создание таблицы conn_log в СУБД ClickHouse YC

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

// Файл сертификата. Структура таблицы
var (
	crtFile  = filepath.Join("..", "certs", "YandexInternalRootCA.crt")
	createDB = `CREATE TABLE conn_log
	(
    user_id UInt32,
    ip_addr IPv4,
    timestamp DateTime,
	)
	ENGINE = MergeTree
	PRIMARY KEY (user_id)`
)

func main() {
	// DSN для подключения к СУБД ClickHouse
	const DB_HOST = "rc1a-u620db3mp7svl13i.mdb.yandexcloud.net"
	const DB_NAME = "duples"
	const DB_USER = "gorest"
	const DB_PASS = "rootroot"

	// Формирование метаданных структуры запроса
	caCert, err := ioutil.ReadFile(crtFile)
	if err != nil {
		p := recover()
		log.Fatalln(err)
		panic(p)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	conn := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	// Форматирование запроса
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s:8443/", DB_HOST), nil)
	query := req.URL.Query()
	query.Add("database", DB_NAME)
	query.Add("query", createDB)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", DB_USER)
	req.Header.Add("X-ClickHouse-Key", DB_PASS)

	resp, err := conn.Do(req) // Выполнение запроса
	if err != nil {
		p := recover()
		log.Fatalln(err)
		panic(p)
	}

	// Отложеное выполнение закрытия запроса, до получения ответа
	defer resp.Body.Close()

	// Чтение данных сервера, обработка ошибок
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}
