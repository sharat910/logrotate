package main

import (
	"github.com/sharat910/logrotator"
	"log"
)

func main() {
	lr, err := logrotator.New("./logs/example.log",
		logrotator.StartHour(3),
		logrotator.Header([]byte("header\n")),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = lr.Write([]byte("hello world\n"))
	if err != nil {
		log.Fatal(err)
	}

	err = lr.Close()
	if err != nil {
		log.Fatal(err)
	}
}
