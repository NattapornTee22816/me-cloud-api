package repository

import (
	clickhouse2 "github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Permission struct {
	Group       string
	Code        string
	Description string
	CreatedTime clickhouse2.DateTime
	CreatedBy   clickhouse2.UUID
	UpdatedTime clickhouse2.DateTime
	UpdatedBy   clickhouse2.UUID
}

type PermissionRepository struct {
	db *clickhouse.Conn
}

func NewPermissionRepository(db *clickhouse.Conn) *PermissionRepository {
	return &PermissionRepository{
		db,
	}
}
