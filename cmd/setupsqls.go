package main

import (
	"gitlab.com/pos_malaysia/golib/database"
)

var sqlStatements = map[string]string{
	// SQL Name : SQL statement
	// ok to have duplicate SQL statements, but not ok to have duplicate SQL names
	// Using map ensure that no duplicate keys. The compiler will stop if there's any duplicate in map literal

	//"GetAllFromTable": "SELECT * FROM SOME_TABLE",

}

func init() {
	database.InitSQLStatements(sqlStatements)
}
