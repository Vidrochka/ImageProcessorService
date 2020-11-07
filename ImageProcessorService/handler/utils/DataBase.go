package utils

import (
	"database/sql"
	"log"

	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler/dto"

	_ "github.com/mattn/go-sqlite3"
)

//DataBase - implement db funtional
type DataBase struct {
	db     *sql.DB
	logger *log.Logger
	config *Configuration
}

//CreateDB - create and connect db
func CreateDB(logger *log.Logger, config *Configuration) *DataBase {
	db, err := sql.Open("sqlite3", config.DataBasePath)

	if err != nil {
		logger.Println(err)
		panic(err)
	}

	return &DataBase{db: db, logger: logger, config: config}
}

//Close - close connection
func (dataBase *DataBase) Close() error {
	return dataBase.db.Close()
}

//GetDBContext - return db handler
func (dataBase *DataBase) GetDBContext() *sql.DB {
	return dataBase.db
}

//CreateTable - create db table
func (dataBase DataBase) CreateTable() error {
	_, err := dataBase.db.Exec("CREATE TABLE IF NOT EXISTS [Images] (" +
		"[Id] INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT," +
		"[Extension] NVARCHAR(64) NOT NULL," +
		"[Name] NVARCHAR(128) NOT NULL," +
		"[Path] NVARCHAR(128) NOT NULL," +
		"[DateCreated] TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")

	if err != nil {
		dataBase.logger.Println(err)
		return err
	}

	dataBase.logger.Println("BB table created")

	return nil
}

//SaveImage - save imafe in db
func (dataBase DataBase) SaveImage(name, extension, path string) (int64, error) {
	result, err := dataBase.db.Exec("insert into Images (Name, Extension, Path) values (?, ?, ?)",
		name, extension, path)

	if err != nil {
		dataBase.logger.Println(err)
		return -1, err
	}

	return result.LastInsertId()
}

//RestoreImage - restore image by id
func (dataBase DataBase) RestoreImage(id int64) (*dto.Image, error) {
	row, err := dataBase.db.Query("SELECT Name, Extension, Path FROM Images WHERE Id = $1", id)

	if err != nil {
		dataBase.logger.Println(err)
		return nil, err
	}

	image := dto.Image{}

	row.Next()
	err = row.Scan(&image.Name, &image.Extension, &image.Data)

	dataBase.logger.Println(image.Data)

	if err != nil {
		dataBase.logger.Println(err)
		return nil, err
	}

	return &image, nil
}
