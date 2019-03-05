package simplechatbot_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	SqlSchemaFile = "./schema/schema.sql"
)

func generateFileName() string {
	now := time.Now().UnixNano()
	return strconv.FormatInt(now, 16)
}

func createTestDatabase() (string, *sql.DB) {
	testDbFile := generateFileName() + ".db"
	db, err := sql.Open("sqlite3", testDbFile)
	if err != nil {
		log.Fatal(err)
	}
	return testDbFile, db
}

func createTables(db *sql.DB) {
	fileContent, _ := ioutil.ReadFile(SqlSchemaFile)
	db.Exec(string(fileContent))
}

func dropTables(db *sql.DB) {
	tableNames := []string{"Bots", "Users", "UserBots", "Responses", "Commands"}
	for _, tableName := range tableNames {
		db.Exec("DROP TABLE $1;", tableName)
	}
}

func populateTables(db *sql.DB) {
	// TODO: aaaaaaaaaaahhhhhhhh
}

func test_ttt() {

}

func TestDatabaseCreation(t *testing.T) {

	var test_test_test int
	test_test_test = 1
	_ = test_test_test
	dbFilename, db := createTestDatabase()
	createTables(db)

	sqlText := "SELECT * FROM UserBots"
	rows, _ := db.Query(sqlText)
	hasRow := rows.Next()
	assert.False(t, hasRow)
	// Clean up DB file
	db.Close()
	os.Remove(dbFilename)
}

func prepareMySQLDB(t *testing.T) (db *sql.DB, cleanup func() error) {
	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, err := sql.Open("mysqltx", cName)

	if err != nil {
		t.Fatalf("open mysqltx connection: %s", err)
	}

	return db, db.Close
}
