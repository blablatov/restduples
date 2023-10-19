// Тестовый вебсервер для отладки REST запросов.
// Метода выполнения запроса к СУБД ClickHouse Yandex Cloud
// go run restduples.go
// go test -v client_test.go

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	sl "github.com/blablatov/restduples/clickhouse/selectlog"
)

// Файлы сертификата и ключа
var (
	crtFile = filepath.Join(".", "certs", "server.crt")
	keyFile = filepath.Join(".", "certs", "server.key")
)

func main() {
	log.SetPrefix("Event main: ")
	log.SetFlags(log.Lshortfile)

	// TLS or simple connect. Подключение TLS или базовое
	http.HandleFunc("/", handle)
	http.ListenAndServeTLS("localhost:12345", crtFile, keyFile, nil)
	//http.ListenAndServe("localhost:12345", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nRequest of client:\n")

	// Параметры заголовков http-запроса. Parameters of headers
	fmt.Fprintf(w, "Method = %s\nURL(user_id) = %s\nProto = %s\n", r.Method, r.URL, r.Proto)
	fmt.Printf("Method = %s\nURL(user_id) = %s\nProto = %s\n", r.Method, r.URL, r.Proto)

	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		fmt.Printf("Header[%q] = %q\n", k, v)
	}

	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Printf("Host = %q\n", r.Host)

	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	fmt.Printf("RemoteAddr = %q\n", r.RemoteAddr)

	ps := strings.Split(r.URL.Path, "")

	// TODO буффер для дополнительной обработки данных, аналог strings.Builder
	var buf bytes.Buffer
	buf.WriteString(r.URL.Path)
	fmt.Println("buf: ", buf.String())

	// TODO Проверка повторного user_id через стек, при ошибке проверить client_test.go на дубли user_id
	var stack []string
	stack = append(stack, r.URL.Path)
	top := stack[len(stack)-1]
	fmt.Println("stack: ", top)

	// Формирование структуры с user_id
	var u sl.UserId
	u.Mu.Lock()
	sr := sl.UserId{
		Userid1: ps[1],
		Userid2: ps[3],
	}
	u.Mu.Unlock()
	fmt.Printf("Userid1 = %s\nUserid2 = %s\n", sr.Userid1, sr.Userid2)

	// Вызов метода выполнения запроса к СУБД ClickHouse
	var wg sync.WaitGroup
	chb := make(chan bool, 1)
	wg.Add(1)
	go func() {
		u.SelectLog(sr.Userid1, sr.Userid2, chb, wg)
		log.Println("This is duples:", <-chb)
	}()
	go func() {
		wg.Wait()
		close(chb)
	}()
}
