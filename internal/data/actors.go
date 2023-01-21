package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Actor struct {
	ID      int64  `json:"id"`      // Unique integer ID for the movie
	Name    string `json:"name"`    // Movie title
	Surname string `json:"surname"` // Movie title
	Version int32  `json:"version"` // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type ActorModel struct {
	DB *sql.DB
}

func (a ActorModel) Get(id int64) (*Actor, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, first_name, last_name, version
		FROM actors
		WHERE id = $1;`

	var actor Actor

	err := a.DB.QueryRow(query, id).Scan(
		&actor.ID,
		&actor.Surname,
		&actor.Name,
		&actor.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &actor, nil
}

func (a ActorModel) GetAll(name string, filters Filters) ([]*Actor, Metadata, error) {
	// Update the SQL query to include the window function which counts the total
	// (filtered) records.
	// (filtered) records.
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, first_name, last_name, version
		FROM actors
		WHERE (to_tsvector('simple', first_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{name, filters.limit(), filters.offset()}

	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}

	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0

	actors := []*Actor{}

	for rows.Next() {
		var actor Actor
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&actor.ID,
			&actor.Surname,
			&actor.Name,
			&actor.Version,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		actors = append(actors, &actor)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return actors, metadata, nil
}
