package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func Connect_DB() {
	var err error

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	db_url := "postgres://postgres:12345@localhost:5432/dumbways41"

	Conn, err = pgx.Connect(context.Background(), db_url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unnable u connect to Database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Succes connect to Database")
}
