package repository

import (
	clickhouse2 "github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type UserInformation struct {
	UserUUID    clickhouse2.UUID
	Language    string
	TagName     string
	TagValue    string
	CreatedTime clickhouse2.DateTime
	CreatedBy   string
	UpdatedTime clickhouse2.DateTime
	UpdatedBy   string
}

type UserInformationRepository struct {
	db *clickhouse.Conn
}

func NewUserInformationRepository(db *clickhouse.Conn) *UserInformationRepository {
	return &UserInformationRepository{
		db,
	}
}
