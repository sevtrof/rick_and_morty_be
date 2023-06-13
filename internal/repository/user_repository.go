package repository

import (
	"database/sql"
	"log"
	"ricknmorty/internal/domain/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	repo := &UserRepository{db: db}
	repo.initializeTable()
	return repo
}

func (r *UserRepository) Save(user *model.User) error {
	log.Printf("Saving user: %v", user)
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("error saving user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		log.Printf("user not found: %v", err)
		return nil, err
	} else if err != nil {
		log.Printf("error searching user: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) initializeTable() {
	createTableStatement := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		);`

	_, err := r.db.Exec(createTableStatement)
	if err != nil {
		log.Printf("error creating users table: %v", err)
	}
}
