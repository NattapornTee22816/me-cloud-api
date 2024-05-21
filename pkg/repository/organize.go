package repository

import (
	clickhouse2 "github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Organize struct {
	UUID         clickhouse2.UUID
	Name         string
	Contact      string
	Website      string
	Status       string
	VerifyStatus string
	CreatedTime  clickhouse2.DateTime
	CreatedBy    clickhouse2.UUID
	UpdatedTime  clickhouse2.DateTime
	UpdatedBy    clickhouse2.UUID
}

type OrganizeRepository struct {
	db *clickhouse.Conn
}

func NewOrganizeRepository(db *clickhouse.Conn) *OrganizeRepository {
	return &OrganizeRepository{
		db,
	}
}
