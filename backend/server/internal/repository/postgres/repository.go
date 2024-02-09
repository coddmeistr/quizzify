package postgres

import "github.com/maxik12233/quizzify-online-tests/server/internal/repository"

type PostgresRepository struct {
}

func NewPostgresRepository() repository.RootRepository {
	return &PostgresRepository{}
}

func (pg *PostgresRepository) CreateTest() {

}
