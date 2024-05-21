package repository

import (
	clickhouse2 "github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type UserSession struct {
	UserUUID    clickhouse2.UUID
	MacAddress  string
	UserAgent   string
	CreatedTime clickhouse2.UUID
	CreatedBy   string
}

type UserSessionRepository struct {
	db *clickhouse.Conn
}

func NewUserSessionRepository(db *clickhouse.Conn) *UserSessionRepository {
	return &UserSessionRepository{
		db,
	}
}
