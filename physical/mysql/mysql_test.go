package mysql

import (
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLBackend(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("MYSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MYSQL_TABLE")
	if table == "" {
		table = "test"
	}

	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewMySQLBackend(map[string]string{
		"address":  address,
		"database": database,
		"table":    table,
		"username": username,
		"password": password,
	}, logger)

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mysql := b.(*MySQLBackend)
		_, err := mysql.client.Exec("DROP TABLE IF EXISTS " + mysql.dbTable + " ," + mysql.dbLockTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestMySQLHABackend(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("MYSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MYSQL_TABLE")
	if table == "" {
		table = "test"
	}

	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewMySQLBackend(map[string]string{
		"address":    address,
		"database":   database,
		"table":      table,
		"username":   username,
		"password":   password,
		"ha_enabled": "true",
	}, logger)

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mysql := b.(*MySQLBackend)
		_, err := mysql.client.Exec("DROP TABLE IF EXISTS " + mysql.dbTable + " ," + mysql.dbLockTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	ha, ok := b.(physical.HABackend)
	if !ok {
		t.Fatalf("MySQL does not implement HABackend")
	}
	physical.ExerciseHABackend(t, ha, ha)
}
