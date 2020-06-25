package postgresdriver

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/KarinaBotova/Normalization/models"

	sq "github.com/Masterminds/squirrel"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var db *sql.DB

// Init database
func InitDatabaseConnection(dbUrl string) (err error) {
	// Open connection
	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		return fmt.Errorf("could not open database connection: %v", err)
	}
	// Test connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}
	return
}

// Init database structure
func InitDatabaseStructure() (err error) {
	// Get data from script
	path := "./script.sql"
	scriptFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	script := string(scriptFile)

	// Execute script
	_, err = db.Exec(script)
	if err != nil {
		return err
	}
	return nil
}

// Close db connection
func CloseConnection() (err error) {
	return db.Close()
}

// Send data to DB
func SaveStudents(students []models.Student) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare queries
	queries := make([]sq.InsertBuilder, 5)
	queries[0] = sq.Insert("faculties").
		Columns("name", "adress_fac").
		Suffix("ON CONFLICT DO NOTHING")
	queries[1] = sq.Insert("cafedries").
		Columns("id", "tel", "faculty").
		Suffix("ON CONFLICT DO NOTHING")
	queries[2] = sq.Insert("specialities").
		Columns("name", "cafedra").
		Suffix("ON CONFLICT DO NOTHING")
	queries[3] = sq.Insert("groups").
		Columns("id", "speciality").
		Suffix("ON CONFLICT DO NOTHING")
	queries[4] = sq.Insert("students").
		Columns("fio", "zachet_book", "\"group\"").
		Suffix("ON CONFLICT DO NOTHING")

	// Bind arguments to queries
	for _, c := range students {
		queries[0] = queries[0].
			Values(c.Faculty, c.Adress)
		queries[1] = queries[1].
			Values(c.Cafedry, c.Tel, c.Faculty)
		queries[2] = queries[2].
			Values(c.Speciality, c.Cafedry)
		queries[3] = queries[3].
			Values(c.Group, c.Speciality)
		queries[4] = queries[4].
			Values(c.FIO, c.ZachetBook, c.Group)
	}

	// Execute queries
	for _, query := range queries {
		q, a, err := query.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			log.Error().Err(err).Msg("Failed to change placeholder")
			return err
		}
		if _, err = tx.Exec(q, a...); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error().Err(rbErr).Msg("Failed to rollback transaction")
			}
			return err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Error().Err(rbErr).Msg("Failed to rollback transaction")
		}
		return err
	}

	return nil
}
