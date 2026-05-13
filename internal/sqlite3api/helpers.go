package sqlite3api

import (
	"io"
	"os"
	"path/filepath"
)

// Sqlite3DbFileIsExist проверяет наличие файла базы данных Sqlite3
// и при необходимости создает его из резервного файла
func Sqlite3DbFileIsExist(rootPath, pathFileDb string) (newPathToDb string, err error) {
	var (
		fr, fw *os.File
	)

	backupFile := filepath.Join(rootPath, "/backupdb/sqlite3_backup.db")
	pathFileDb = filepath.Join(rootPath, pathFileDb)

	newPathToDb = pathFileDb

	// наличие файла backup
	if _, err = os.Stat(backupFile); err != nil {
		return
	}

	//файл с основной БД
	_, err = os.Stat(pathFileDb)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}

		fr, err = os.OpenFile(backupFile, os.O_RDONLY, 0666)
		if err != nil {
			return
		}
		defer fr.Close()

		fw, err = os.OpenFile(pathFileDb, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer fw.Close()

		if _, err = io.Copy(fw, fr); err != nil {
			return
		}
	}

	return
}
