package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/config"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage"
)

type Clickhouse struct {
	Db *sql.DB
}

func New(conf *config.Config) (*Clickhouse, error) {
	clickhouse := &Clickhouse{}

	var err error
	clickhouse.Db, err = connect(conf)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to clickhouse: %s", err.Error())
	}

	return clickhouse, nil
}

func (c *Clickhouse) Test(metric storage.Metric) error {
	tx, err := c.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error occured while starting transaction: %s", err.Error())
	}
	timestamp, err := time.Parse(time.DateTime, metric.TimeStamp)
	if err != nil {
		return fmt.Errorf("error occured while parsing timestamp: %s", err.Error())
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

func connect(conf *config.Config) (*sql.DB, error) {
	dbInfo := url.Values{}
	dbInfo.Add("database", conf.Clickhouse.Name)
	dbInfo.Add("username", conf.Clickhouse.User)
	dbInfo.Add("password", conf.Clickhouse.Password)
	urlConnectionString := url.URL{
		Scheme:   "tcp",
		Host:     fmt.Sprintf("%s:%s", conf.Clickhouse.Host, conf.Clickhouse.Port),
		RawQuery: dbInfo.Encode(),
	}

	Db, err := sql.Open("clickhouse", urlConnectionString.String())
	if err != nil {
		return nil, fmt.Errorf("error connecting to ClickHouse: %s", err.Error())
	}

	// Ping the database to verify the connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = Db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while pinging clickhouse: %s", err.Error())
	}

	return Db, nil
}
