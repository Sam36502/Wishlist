package model

//
//		DB Schema Functions
//
//	To ensure the schema matches expectations
//	and to provide meta functions
//

import (
	"database/sql"
	"fmt"
)

const (
	SCHEMA_VERSION = "1.1.0"
	SCHEMA_META    = "META"
)

func EnsureSchema(db *sql.DB) {
	if db == nil || db.Ping() != nil {
		fmt.Println("Error: tried to ensure schema of invalid database connection")
		return
	}

	//EnsureMetaTables(db)
}

func EnsureMetaTables(db *sql.DB) {

	// Meta Info
	row := db.QueryRow("SELECT VALUE FROM META_INFO WHERE FIELD = 'schema_version';")
	schema_version := ""
	err := row.Scan(&schema_version)
	if err != nil || schema_version != SCHEMA_VERSION {
		fmt.Println("--> Table META_INFO is invalid; rebuilding...")
		_, err = db.Exec("CREATE TABLE META_INFO (FIELD TEXT PRIMARY KEY, VALUE TEXT);")
		if err != nil {
			fmt.Printf("     Reconstruction failed! (%e)\n", err)
			return
		}
	}

}

func EnsureUserTables(db *sql.DB, prefix string) {

	// User
	row := db.QueryRow("SELECT * FROM " + prefix + "_USER LIMIT 1;")
	schema_version := ""
	err := row.Scan(&schema_version)
	if err != nil || schema_version != SCHEMA_VERSION {
		fmt.Println("--> Table META_INFO is invalid; rebuilding...")
		_, err = db.Exec("CREATE TABLE META_INFO (FIELD TEXT PRIMARY KEY, VALUE TEXT);")
		if err != nil {
			fmt.Printf("     Reconstruction failed! (%e)\n", err)
			return
		}
	}

}
