// Генератор строк записей для таблицы conn_log, СУБД ClickHouse.
// Ключ - это число итераций основного цикла с 4 подциклами для дополнительных итераций.
// Общее количество строк, определяется простым подбором ключа.
// $gendtables 1000000

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	log.SetPrefix("Event main: ")
	log.SetFlags(log.Lshortfile)

	var num, sep string

	for i := 1; i < len(os.Args); i++ {
		num += sep + os.Args[i]
		sep = " "
	}

	sl := len("")
	if len(num) == sl {
		log.Fatalf("\nВведите количество итераций генератора\n")
	}

	n, err := strconv.Atoi(num)
	if err != nil {
		log.Fatalf("\nВведите число\n")
	}

	start := time.Now()
	var buf bytes.Buffer

	buf.WriteString("user_id, ip_addr, ts")
	for ; n > 0; n-- {
		for i := 0; i < 3; i++ {
			fmt.Fprintf(&buf, "%v\n", i)
			buf.WriteString("1, 127.0.0.1, 17:51:1")
		}
		for i := 0; i < 2; i++ {
			fmt.Fprintf(&buf, "%v\n", i)
			buf.WriteString("3, 127.0.0.2, 17:51:2")
		}
		for i := 0; i < 3; i++ {
			fmt.Fprintf(&buf, "%v\n", i)
			buf.WriteString("2, 127.0.0.4, 17:51:2")
		}
		for i := 0; i < 2; i++ {
			fmt.Fprintf(&buf, "%v\n", i)
			buf.WriteString("4, 127.0.0.3, 17:51:2")
		}
	}
	//fmt.Println(buf.String())
	fb := buf.Bytes()
	rb := bytes.TrimPrefix(fb, []byte("0\n"))
	err = ioutil.WriteFile("conn_log.csv", rb, 0644)
	if err != nil {
		log.Fatal(err)
	}
	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs Time of create file\n", secs)
}
