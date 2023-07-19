package clickhouse

import (
	"context"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Clickhouse struct {
	db driver.Conn
}

func New() (*Clickhouse, error) {
	db, err := connect()
	if err != nil {
		return nil, fmt.Errorf("error while connecting to clickhouse: %s", err.Error()) //log.Fatalf("error while connecting: %s", err.Error())
	}

	return &Clickhouse{
		db: db,
	}, nil
}

func (c *Clickhouse) Test() error {
	rows, err := c.db.Query(context.Background(), "SELECT name, sex FROM human")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var (
			name, sex string
		)
		if err := rows.Scan(
			&name,
			&sex,
		); err != nil {
			log.Fatal(err)
		}
		log.Printf("name: %s, sex: %s", name, sex)
	}
	return nil
}

func connect() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			// Other options like TLS configuration can be added here
		})
	)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}
