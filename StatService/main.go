package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gorilla/mux"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	chClient, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"0.0.0.0:19400"},
		Auth: clickhouse.Auth{
			Database: "reactions",
			Username: "user",
			Password: "alphabet",
		},
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
		MaxIdleConns: 300,
	})
	if err != nil {
		panic(err)
	}

	chClient.Exec(context.Background(),
		"CREATE TABLE IF NOT EXISTS (reaction int) reactions ENGINE = Kafka('kafka:29092', 'my_topic', 'my_group', 'earliest');")

	router := mux.NewRouter()
	router.HandleFunc("/ping", PingHandler).Methods("Get")

	http.ListenAndServe(":8080", router)
	select {}
}
