package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sharat910/logrotator"
)

func main() {
	lr, err := logrotator.New("./logs/example.log",
		logrotator.StartHour(3),
		logrotator.PrependTimeFormat("2006-01-02", "_"),
		logrotator.WithHeaderWriter(func(w io.Writer) error {
			_, err := w.Write([]byte("header\n"))
			return err}),
		logrotator.WithImmediateFlush,
		logrotator.RotateCallback(func(t time.Time) {
			fmt.Println("Rotating logs", t)
		}),
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
