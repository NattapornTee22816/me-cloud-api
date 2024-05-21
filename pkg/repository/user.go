package repository

import (
	clickhouse2 "github.com/ClickHouse/clickhouse-go"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type User struct {
	UUID         clickhouse2.UUID
	Email        string
	Secret       string
	VitalStatus  string
	VerifyStatus string
	CreatedTime  clickhouse2.DateTime
	CreatedBy    string
	UpdatedTime  clickhouse2.DateTime
	UpdatedBy    string
}

type UserRepository struct {
	db *clickhouse.Conn
}

func NewUserRepository(db *clickhouse.Conn) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (r *UserRepository) Save(row User) *User {
	return nil
}

func (r *UserRepository) Delete(row User) {

}
