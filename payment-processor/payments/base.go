package payments

import (
	"log"
	"os"
	"path/filepath"
)

var DB_ROOT string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: File path error")
	}

	DB_ROOT = filepath.Join(dir, "db")
}
