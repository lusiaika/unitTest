package database

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

type Database struct {
	ConnectionString string
	SqlDb            *sql.DB
}

func NewSqlConnection(connectionString string) *Database {
	s := Database{
		ConnectionString: connectionString,
	}

	db, err := sql.Open("sqlserver", s.ConnectionString)
	if err != nil {
		fmt.Printf("[mssql] Error connecting to SQL Server: %v", err)
	}

	s.SqlDb = db
	s.SqlDb.SetMaxIdleConns(255)
	s.SqlDb.SetMaxOpenConns(255)

	return &s
}
