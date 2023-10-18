// Вставка данных лога через SQL запрос.
// TODO вставка данных посредством сгенеренного .cvs файла

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

// Файл сертификата. Структура запроса
var (
	crtFile = filepath.Join("..", "certs", "YandexInternalRootCA.crt")
	//insertLog = `INSERT INTO conn_log FROM INFILE 'con_log.csv' FORMAT CSV`
	insertLog = `INSERT INTO conn_log (user_id, ip_addr, timestamp) VALUES (1, '127.0.0.1', '2023-10-15 00:00:00'), (2, '127.0.0.1', '2023-10-15 00:00:00'), (2, '127.0.0.3', '2023-10-15 00:00:00'), (1, '127.0.0.4', '2023-10-15 00:00:00'), (2, '127.0.0.2', '2023-10-15 00:00:00'), (2, '127.0.0.3', '2023-10-15 00:00:00'), (3, '127.0.0.3', '2023-10-15 00:00:00'), (3, '127.0.0.1', '2023-10-15 00:00:00'), (4, '127.0.0.1', '2023-10-15 00:00:00'), (1, '127.0.0.3', '2023-10-15 00:00:00')`
)

func main() {
	// DSN для подключения к СУБД ClickHouse.
	const DB_HOST = "rc1a-u620db3mp7svl13i.mdb.yandexcloud.net"
	const DB_NAME = "duples"
	const DB_USER = "gorest"
	const DB_PASS = "rootroot"

	// Формирование метаданных структуры запроса. Struct of request
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
	query.Add("query", insertLog)

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
