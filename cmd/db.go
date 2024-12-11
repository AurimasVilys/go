package cmd

import (
	"github.com/stephenafamo/bob"
	"os"
)

func openDB() (*bob.DB, error) {
	db, err := bob.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		return nil, err
	}

	return &db, nil
}
