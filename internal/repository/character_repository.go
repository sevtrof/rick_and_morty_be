package repository

import (
	"context"
	"database/sql"
	"fmt"
	"ricknmorty/internal/domain/model"
	"strings"
	"time"

	"github.com/lib/pq"
)

type CharacterRepository struct {
	db *sql.DB
}

func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

func (repo *CharacterRepository) FetchCharacters(filters map[string]string, page int) (model.CharactersWithInfo, error) {
	const pageSize = 20
	const timeout = 3 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	page, offset := normalizePageNumber(page, pageSize)
	whereClause, params := buildWhereClause(filters)

	totalCount, err := repo.fetchTotalCount(ctx, whereClause, params)
	if err != nil {
		return model.CharactersWithInfo{}, fmt.Errorf("error fetching total count: %w", err)
	}

	characters, err := repo.fetchCharacters(ctx, whereClause, params, pageSize, offset)
	if err != nil {
		return model.CharactersWithInfo{}, fmt.Errorf("error fetching characters: %w", err)
	}

	info := buildInfo(totalCount, pageSize, page)
	return model.CharactersWithInfo{
		Characters: characters,
		Info:       info,
	}, nil
}

func normalizePageNumber(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	return page, (page - 1) * pageSize
}

func buildWhereClause(filters map[string]string) (string, []interface{}) {
	whereFiltersClauses := []string{}
	params := []interface{}{}
	placeholderIndex := 1

	for key, value := range filters {
		if value != "" {
			whereFiltersClauses = append(whereFiltersClauses, fmt.Sprintf("%s ILIKE $%d", key, placeholderIndex))
			if key != "gender" && key != "species" {
				value = "%" + value + "%"
			}
			params = append(params, value)
			placeholderIndex++
		}
	}

	whereClause := ""
	if len(whereFiltersClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereFiltersClauses, " AND ")
	}
	return whereClause, params
}

func (repo *CharacterRepository) fetchTotalCount(ctx context.Context, whereClause string, params []interface{}) (int, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM characters %s", whereClause)
	row := repo.db.QueryRowContext(ctx, countQuery, params...)
	var totalCount int
	if err := row.Scan(&totalCount); err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (repo *CharacterRepository) fetchCharacters(ctx context.Context, whereClause string, params []interface{}, pageSize, offset int) ([]model.Character, error) {
	placeholderIndex := len(params) + 1
	charactersQuery := fmt.Sprintf("SELECT * FROM characters %s LIMIT $%d OFFSET $%d", whereClause, placeholderIndex, placeholderIndex+1)
	params = append(params, pageSize, offset)

	rows, err := repo.db.QueryContext(ctx, charactersQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []model.Character
	for rows.Next() {
		var character model.Character
		var episode pq.StringArray

		if err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.Status,
			&character.Species,
			&character.Type,
			&character.Gender,
			&character.Image,
			&character.Url,
			&character.Created,
			&character.Location.Name,
			&character.Location.Url,
			&character.Origin.Name,
			&character.Origin.Url,
			&episode,
		); err != nil {
			return nil, err
		}

		character.Episode = episode
		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}

func buildInfo(totalCount, pageSize, page int) model.Info {
	info := model.Info{
		Count: totalCount,
		Pages: (totalCount + pageSize - 1) / pageSize,
	}

	if page > 1 {
		info.Prev = fmt.Sprintf("/characters?page=%d", page-1)
	}

	if (page-1)*pageSize+pageSize < totalCount {
		info.Next = fmt.Sprintf("/characters?page=%d", page+1)
	}

	return info
}
