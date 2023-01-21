package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel
// kind of enveloping
type Models struct {
	Movies    MovieModel
	Actors    ActorModel
	Directors DirectorModel
	Users     UserModel
	Tokens    TokenModel
}

// method which returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies:    MovieModel{DB: db},
		Actors:    ActorModel{DB: db},
		Directors: DirectorModel{DB: db},
		Users:     UserModel{DB: db},
		Tokens:    TokenModel{DB: db},
	}
}
