package repository

import "github.com/ClickHouse/clickhouse-go/v2"

type AbstractRepository[T any] interface {
	Save(row T) T
	Delete(row T)
}

type Repository struct {
	ApiKey          *ApiKeyRepository
	Application     *ApplicationRepository
	Organize        *OrganizeRepository
	Permission      *PermissionRepository
	User            *UserRepository
	UserInformation *UserInformationRepository
	UserSession     *UserSessionRepository
}

func NewRepository(db *clickhouse.Conn) *Repository {
	return &Repository{
		ApiKey:          NewApiKeyRepository(db),
		Application:     NewApplicationRepository(db),
		Organize:        NewOrganizeRepository(db),
		Permission:      NewPermissionRepository(db),
		User:            NewUserRepository(db),
		UserInformation: NewUserInformationRepository(db),
		UserSession:     NewUserSessionRepository(db),
	}
}
