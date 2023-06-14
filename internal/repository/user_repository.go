package repository

import (
	"database/sql"
	"fmt"
	"log"
	"ricknmorty/internal/domain/model"
	"strconv"

	"github.com/lib/pq"
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
	query := `INSERT INTO users (username, email, password, avatar) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Avatar)
	if err != nil {
		log.Printf("error saving user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	query := `SELECT id, username, email, password, favourite_characters FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	user := &model.User{}
	var favouriteCharacters pq.Int64Array
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &favouriteCharacters)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		log.Printf("error scanning user row: %v", err)
		return nil, err
	}

	intSlice := make([]int, len(favouriteCharacters))
	for i, v := range favouriteCharacters {
		intSlice[i] = int(v)
	}

	user.FavouriteCharacters = intSlice
	return user, nil
}

func (r *UserRepository) AddFavouriteCharacter(userId int, characterId int) error {
	query := `UPDATE users 
	SET favourite_characters = array_append(favourite_characters, $1)
	WHERE id = $2 AND NOT EXISTS (
		SELECT 1 FROM UNNEST(favourite_characters) AS elem WHERE elem = $1
	)`
	_, err := r.db.Exec(query, characterId, userId)
	return err
}

func (r *UserRepository) RemoveFavouriteCharacter(userId int, characterId int) error {
	query := `UPDATE users SET favourite_characters = array_remove(favourite_characters, $1) WHERE id = $2`
	_, err := r.db.Exec(query, characterId, userId)
	return err
}

func (r *UserRepository) FindById(id int) (*model.User, error) {
	query := `SELECT id, username, email, password, favourite_characters FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	user := &model.User{}
	var favouriteCharacters []string
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, pq.Array(&favouriteCharacters))

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		log.Printf("error scanning user row: %v", err)
		return nil, err
	}

	intSlice := make([]int, len(favouriteCharacters))
	for i, v := range favouriteCharacters {
		intValue, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("error converting string to int: %v", err)
			return nil, err
		}
		intSlice[i] = intValue
	}

	user.FavouriteCharacters = intSlice
	return user, nil
}

func (r *UserRepository) initializeTable() {
	createTableStatement := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			favourite_characters INTEGER[] DEFAULT ARRAY[]::INTEGER[],
			avatar VARCHAR(255) NULL
		);`

	_, err := r.db.Exec(createTableStatement)
	if err != nil {
		log.Printf("error creating users table: %v", err)
	}
}
