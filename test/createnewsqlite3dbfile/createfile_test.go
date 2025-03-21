package createnewsqlite3dbfile_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// checkSqlite3DbFileExist проверяет наличие файла БД Sqlite3
// и при необходимости создает его из резервного файла
func checkSqlite3DbFileExist(pathFileDb string) error {
	backupFile := "../../backupdb/sqlite3_backup.db"

	// наличие файла backup
	if _, err := os.Stat(backupFile); err != nil {
		return err
	}

	//файл с основной БД
	_, err := os.Stat(pathFileDb)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		fr, err := os.OpenFile(backupFile, os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		defer fr.Close()

		fw, err := os.OpenFile(pathFileDb, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer fw.Close()

		if _, err := io.Copy(fw, fr); err != nil {
			return err
		}
	}

	return nil
}

func TestCreateFile(t *testing.T) {
	err := checkSqlite3DbFileExist("../../sqlite3/sqlite3.db")
	assert.NoError(t, err)
}
