package migrator

import (
	"bufio"
	"database/sql"
	"os"
	"strings"
)

func Migrate(filepath string, db *sql.DB) error {
	// Open the file
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// read line till ';' and execute
	statement := strings.Builder{}
	for scanner.Scan() {
		line := scanner.Text()
		statement.WriteString(line)
		if strings.Contains(line, ";") {
			_, err := db.Exec(statement.String())
			if err != nil {
				return err
			}
			statement.Reset()
		}
	}
	return nil
}
