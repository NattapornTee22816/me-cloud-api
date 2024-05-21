package repository

import (
	clickhouse2 "github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type ApiKey struct {
	AppUUID              clickhouse2.UUID
	KeyUUID              clickhouse2.UUID
	AccessKey            string
	AccessKeyExpireTime  clickhouse2.DateTime
	RefreshKey           string
	RefreshKeyExpireTime clickhouse2.DateTime
	Permissions          []string
	Status               string
	CreatedTime          clickhouse2.DateTime
	CreatedBy            string
	UpdatedTime          clickhouse2.DateTime
	UpdatedBy            string
}

type ApiKeyRepository struct {
	db *clickhouse.Conn
}

func NewApiKeyRepository(db *clickhouse.Conn) *ApiKeyRepository {
	return &ApiKeyRepository{
		db,
	}
}

func (r *ApiKeyRepository) Save(row ApiKey) *ApiKey {

	return nil
}

func (r *ApiKeyRepository) Delete(row ApiKey) {

}
