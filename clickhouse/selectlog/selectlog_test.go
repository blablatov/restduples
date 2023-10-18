// Подключение к тестовой БД duples СУБД ClickHouse ОП Yandex Cloud
// go test -v selectlog_test.go

package selectlog

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestSelectLog(t *testing.T) {
	var tests = []struct {
		dbhost string
		dbname string
		dbuser string
		dbpass string
	}{
		{"rc1a-u620db3mp7svl13i.mdb.yandexcloud.net", "duples", "gorest", "rootroot"},
		{"rc1a-u620db3mp7svl13i.mdb.yandexcloud.net", "duples", "gorest", "rootroot"},
		{"rc1a-u620db3mp7svl13i.mdb.yandexcloud.net", "duples", "gorest", "rootroot"},
	}

	var prev_dbhost string
	for _, test := range tests {
		if test.dbhost != prev_dbhost {
			fmt.Printf("\n%s\n", test.dbhost)
			prev_dbhost = test.dbhost
		}
	}

	var prev_dbname string
	for _, test := range tests {
		if test.dbname != prev_dbname {
			fmt.Printf("\n%s\n", test.dbname)
			prev_dbname = test.dbname
		}
	}
	var prev_dbuser string
	for _, test := range tests {
		if test.dbuser != prev_dbuser {
			fmt.Printf("\n%s\n", test.dbuser)
			prev_dbuser = test.dbuser
		}
	}
	var prev_dbpass string
	for _, test := range tests {
		if test.dbpass != prev_dbpass {
			fmt.Printf("\n%s\n", test.dbpass)
			prev_dbpass = test.dbpass
		}
	}

	// Функциональный SQL запрос для получения дублей из БД
	duplesGet := `SELECT ip_addr, COUNT(*) FROM duples.conn_log WHERE user_id = 1 OR user_id = 2 GROUP BY ip_addr HAVING COUNT (*) > 1`

	// Чтение сертификата из файла. Формирование метаданных запроса
	caCert, err := ioutil.ReadFile(crtFile)
	if err != nil {
		// Восстанавливается для анализа, после вывода err, завершается
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

	// Формирование rest-запроса, его заголовков и тела
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s:8443/", prev_dbhost), nil)
	query := req.URL.Query()
	query.Add("database", prev_dbname)
	query.Add("query", duplesGet)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", prev_dbuser)
	req.Header.Add("X-ClickHouse-Key", prev_dbpass)

	// Выполнение запроса
	resp, err := conn.Do(req)
	if err != nil {
		// Восстанавливается для анализа, после вывода err, завершается
		p := recover()
		log.Fatalln(err)
		panic(p)
	}

	// Отложеное выполнение закрытия запроса, до выполнения метода и получения ответа
	defer resp.Body.Close()

	fmt.Printf("Status = %v ", resp.Status) // Статус ответа сервера

	// Чтение данных сервера, обработка ошибок
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("\nDuples:", false)
	}
	log.Println("\nResponse server:\n", string(data))
	fmt.Println("Duples:", true)
}
