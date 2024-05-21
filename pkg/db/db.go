package db

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"os"
	"strings"
)

func NewDatabase() (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: strings.Split(os.Getenv("DB_ADDR"), ","),
		Auth: clickhouse.Auth{
			Database: os.Getenv("DB_NAME"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
		},
		Protocol:         clickhouse.Native,
		MaxOpenConns:     10,
		MaxIdleConns:     10,
		ConnOpenStrategy: clickhouse.ConnOpenRoundRobin,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "me-cloud-go-client", Version: "0.0.1"},
			},
		},
		Settings: map[string]any{
			"max_execution_time": 60,
		},
		//TLS: &tls.Config{
		//	InsecureSkipVerify: true,
		//},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug: os.Getenv("DB_DEBUG") == "true",
		Debugf: func(format string, v ...any) {
			fmt.Printf(format, v)
		},
	})
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}
