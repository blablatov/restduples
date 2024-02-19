// Подключение к тестовой БД duples СУБД ClickHouse ОП Yandex Cloud
// go run selectlog_test.go
package selectlog

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Структура user_id
type UserId struct {
	Userid1 string
	Userid2 string
	Mu      sync.Mutex
}

// Файл сертификата
var (
	crtFile = filepath.Join(".", "certs", "YandexInternalRootCA.crt")
)

// Метод подключения, аутентификации, выполнения запроса к БД duples
func (u *UserId) SelectLog(Userid1, Userid2 string, chb chan bool, wg sync.WaitGroup) error {

	defer wg.Done()

	// DSN для подключения к СУБД ClickHouse.
	const DB_HOST = "rc1a-u620db3mp7svl13i.mdb.yandexcloud.net"
	const DB_NAME = "duples"
	const DB_USER = "gorest"
	const DB_PASS = "rootroot"

	// Тип запроса select
	duplesGet := `SELECT ip_addr, COUNT(*) FROM duples.conn_log WHERE ` + Userid1 + ` OR ` + Userid2 + ` GROUP BY ip_addr HAVING COUNT (*) > 1`

	// Чтение сертификата из файла. Формирование метаданных запроса
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

	// Формирование rest-запроса, его заголовков и тела
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s:8443/", DB_HOST), nil)
	query := req.URL.Query()
	query.Add("database", DB_NAME)
	query.Add("query", duplesGet)

	req.URL.RawQuery = query.Encode()

	req.Header.Add("X-ClickHouse-User", DB_USER)
	req.Header.Add("X-ClickHouse-Key", DB_PASS)

	start := time.Now()
	// Выполнение запроса
	resp, err := conn.Do(req)
	if err != nil {
		//log.Fatalln(err)
		log.Println(err)
		return err
	}

	// Отложеное выполнение закрытия запроса, до выполнения метода и получения ответа
	defer resp.Body.Close()

	fmt.Printf("Status = %v\n", resp.Status) // Статус ответа сервера

	// Чтение данных сервера, обработка ошибок
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs Time of request\n", secs)

	log.Println("\nResponse ClickHouse:\n", string(data))

	if strings.Contains(string(data), "Exception") {
		fmt.Println("\nDuples:", false)
		chb <- false
	} else {
		fmt.Println("\nDuples:", true)
		chb <- true
	}
	return nil
}
