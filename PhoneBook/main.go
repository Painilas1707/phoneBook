package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {

	urlExample := "postgres://postgres:root@localhost:5432/Phonebook"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе", err)
	}
	defer conn.Close(context.Background())
}
