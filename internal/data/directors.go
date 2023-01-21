package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Director struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Awards  []string `json:"awards,omitempty"`
}

type DirectorModel struct {
	DB *sql.DB
}

func (d DirectorModel) Insert(director *Director) error {
	query := `
		INSERT INTO directors(name, surname, awards)
		VALUES ($1, $2, $3)
		RETURNING id`

	return d.DB.QueryRow(query, &director.Name, &director.Surname, pq.Array(&director.Awards)).Scan(&director.ID)
}

func (d DirectorModel) Get(id int64) (*Director, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM directors
		WHERE id = $1`

	var director Director

	err := d.DB.QueryRow(query, id).Scan(
		&director.ID,
		&director.Name,
		&director.Surname,
		pq.Array(&director.Awards),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &director, nil
}

func (d DirectorModel) GetAll(name string, awards []string, filters Filters) ([]*Director, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, surname, awards
		FROM directors
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (awards @> $2 OR $2 = '{}')
		ORDER BY %s
		LIMIT $3 OFFSET $4`, filters.sortDirection2())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	args := []interface{}{name, pq.Array(awards), filters.limit(), filters.offset()}
	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	totalRecords := 0
	directors := []*Director{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var director Director
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords,
			&director.ID,
			&director.Name,
			&director.Surname,
			pq.Array(&director.Awards),
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		directors = append(directors, &director)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// If everything went OK, then return the slice of movies.
	return directors, metadata, nil
}
