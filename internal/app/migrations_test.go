package app

import (
	"database/sql"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"testing"
)

func TestMigrationUpgradeRestartAndFutureRefusal(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err = db.Exec("CREATE TABLE host_migrations(version INTEGER PRIMARY KEY,applied TEXT NOT NULL)"); err != nil {
		t.Fatal(err)
	}
	h := &apphost.Host{DB: db}
	if err = h.Migrate(1, schemas[0]); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO places VALUES('synthetic-id','synthetic-owner','Test place','',51.5,-0.1,'Home','Private synthetic note')"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err = migrate(h); err != nil {
			t.Fatalf("upgrade/restart %d: %v", i, err)
		}
	}
	var note string
	if err = db.QueryRow("SELECT notes FROM places WHERE owner='synthetic-owner'").Scan(&note); err != nil || note != "Private synthetic note" {
		t.Fatal("upgrade changed saved place", err)
	}
	if _, err = db.Exec("INSERT INTO recent_places VALUES('synthetic-owner','k','{}',1)"); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO host_migrations VALUES(4,'future')"); err != nil {
		t.Fatal(err)
	}
	if err = migrate(h); err == nil {
		t.Fatal("future data must not start on this binary")
	}
}
