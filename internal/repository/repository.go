package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nickolation/grpc-task-hezzl/model"
)

type GrpcUsersRepository interface {
	NewUser(u *model.User) (int, error)
	DeleteUser(username string) error
	GetUserList() ([]model.User, error)
}

type Repository struct {
	Db *sqlx.DB
}

func NewGrpsUserRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Db: db,
	}
}

const _usersTable = "users"


// NewUser-method part 
var _insertUserTableQuery = "INSERT INTO %s (username, password_hash, gender, age, description, user_hash) values ($1, $2, $3, $4, $5, $6) RETURNING id"
var _makeInsertTableQuery = func(t string) string {
	return fmt.Sprintf(_insertUserTableQuery, t)
}

const _failCRUDRowIndex = 0

func (repo *Repository) NewUser(u *model.User) (int, error) {
	query := _makeInsertTableQuery(_usersTable)
	row := repo.Db.QueryRow(
		query,
		u.Username,
		u.Password,
		u.Gender,
		u.Age,
		u.Description,
		u.Hash,
	)

	var idx int
	if err := row.Scan(&idx); err != nil {
		return _failCRUDRowIndex, err
	}

	return idx, nil
}

// DeleteUser-method part 
var _deleteFromUsersTableQuery = "DELETE FROM %s u WHERE u.username = $1"
var _makeDeleteTableQuery = func() string {
	return fmt.Sprintf(_deleteFromUsersTableQuery, _usersTable)
}

func (repo *Repository) DeleteUser(username string) error {
	query := _makeDeleteTableQuery()
	_, err := repo.Db.Exec(query, username)
	return err
}


// GetUserList-method part 
var _getAllUsersFromTableQuery = "SELECT * FROM %s u"
var _makeGetTableQuery = func() string {
	return fmt.Sprintf(_getAllUsersFromTableQuery, _usersTable)
}

func (repo *Repository) GetUserList() ([]model.User, error) {
	res := make([]model.User, 0)

	query := _makeGetTableQuery()
	err := repo.Db.Select(&res, query)

	return res, err
}
