package repository

import (
	clickhouse2 "github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Application struct {
	OrganizeUUID   clickhouse2.UUID
	AppUUID        clickhouse2.UUID
	AppName        string
	AppDescription string
	DefaultApp     string
	Status         string
	CreatedTime    clickhouse2.DateTime
	CreatedBy      string
	UpdatedTime    clickhouse2.DateTime
	UpdatedBy      string
}

type ApplicationRepository struct {
	db *clickhouse.Conn
}

func NewApplicationRepository(db *clickhouse.Conn) *ApplicationRepository {
	return &ApplicationRepository{
		db,
	}
}

func (r *ApplicationRepository) Save(row Application) *Application {
	return nil
}

func (r *ApplicationRepository) Delete(row Application) {

}
