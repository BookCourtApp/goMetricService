package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage"
)

type Clickhouse struct {
	db *sql.DB
}

func New() (*Clickhouse, error) {
	db, err := connect()
	if err != nil {
		return nil, fmt.Errorf("error while connecting to clickhouse: %s", err.Error())
	}

	return &Clickhouse{
		db: db,
	}, nil
}

func (c *Clickhouse) Init() error {
	tx, err := c.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error occured while starting transaction: %s", err.Error())
	}
	query := `
		create table if not exists metrics(
			TimeStamp 	DateTime,
			IsApp 		UInt8,
			IsAuth 		UInt8,
			IsNew 		UInt8,
			ResWidth 	UInt16,
			ResHeight 	UInt16,
			UserAgent 	String,
			UserId 		String,
			SessionID 	String,
			DeviceType 	String,
			Reffer 		String,
			Stage 		LowCardinality(String),
			Action 		LowCardinality(String),
			ExtraKeys 	Array(String),
			ExtraValues Array(String)
		) 
		engine = MergeTree() 
		order by Action
		`
	_, err = tx.ExecContext(context.Background(), query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error occured while inserting in database: %s", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %s", err.Error())
	}

	return nil
}

func (c *Clickhouse) Test(metric storage.Metric) error {
	tx, err := c.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error occured while starting transaction: %s", err.Error())
	}
	timestamp, err := time.Parse("2006-01-02 15:04:05", metric.TimeStamp)
	if err != nil {
		return fmt.Errorf("error occured while parsing timestamp")
	}
	query := `
		insert into metrics(
			TimeStamp 	,
			IsApp 		,
			IsAuth 	 	,
			IsNew 		,
			ResWidth ,
			ResHeight ,
			UserAgent ,
			UserId 		,
			SessionID,
			DeviceType,
			Reffer 	,
			Stage 	,
			Action ,
			ExtraKeys,
			ExtraValues
			)
		values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,)
	`
	_, err = tx.ExecContext(context.Background(), query,
		timestamp,
		metric.IsApp,
		metric.IsAuth,
		metric.IsNew,
		metric.ResWidth,
		metric.ResHeight,
		metric.UserAgent,
		metric.UserID,
		metric.SessionID,
		metric.DeviceType,
		metric.Reffer,
		metric.Stage,
		metric.Action,
		metric.ExtraKeys,
		metric.ExtraValues,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error occured while inserting in database: %s", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %s", err.Error())
	}

	return nil
}

func connect() (*sql.DB, error) {
	//connectionString := "tcp://localhost:8123?username=your_username&password=your_password&database=test_db"
	connectionString := "tcp://localhost:9000?&database=test_db"

	// Create the connection to ClickHouse
	db, err := sql.Open("clickhouse", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to ClickHouse: %s", err.Error())
	}

	// Ping the database to verify the connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while pinging clickhouse: %s", err.Error())
	}

	return db, nil
}
