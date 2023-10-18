// Тестовый rest-запрос к серверу go test -v client_test.go

package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestRestClient(t *testing.T) {
	var tests = []struct {
		user_id string
	}{
		{"1/2"},
		{"1/3"},
		{"2/1"},
		{"2/3"},
		{"3/2"},
		{"1/4"},
		{"3/4"},
		{"4/2"},
	}

	var prev_user_id string
	for _, test := range tests {
		if test.user_id != prev_user_id {
			fmt.Printf("\n%s\n", test.user_id)
			prev_user_id = test.user_id
		}

		// URL тестового сервера локально. Для облака указать внешний IP ВМ.
		apiUrl := "https://localhost:12345/" + test.user_id

		// Подгрузка сертификата и ключа. Loads the certs
		cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			log.Fatalf("Сертификат и ключ не получены: %v\n", err)
		}
		// Logs CLIENT_SERVER_HANDSHAKE_TRAFFIC_SECRETS
		var w io.Writer
		w = os.Stdout

		// Форматирование запроса. Formatting of the request
		req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
		// Формирование заголовков запроса. Headers of request
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Формирование метаданных структуры запроса. Struct of request
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					KeyLogWriter:       w,
					Certificates:       []tls.Certificate{cert},
					InsecureSkipVerify: true,
				},
			},
		}

		resp, err := client.Do(req) // Выполнение запроса. Send of request
		if err != nil {
			// Восстанавливается для анализа, после вывода err, завершается
			p := recover()
			log.Fatalln(err)
			panic(p)
		}

		// Отложеное выполнение закрытия запроса, до выполнения метода и получения ответа
		// Defer to finished the method and got response
		defer resp.Body.Close()

		fmt.Printf("Status = %v ", resp.Status) // Статус ответа сервера. Status of response

		// Чтение данных сервера, обработка ошибок. Reads data from server, check errors
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}
		log.Println("\nResponse of server: \n", string([]byte(body)))
	}
}
