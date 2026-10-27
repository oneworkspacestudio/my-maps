package app

import (
	"context"
	"database/sql"
	"testing"
)

func TestShortcutReplacementIsPrivateAndPreservesPreviousPlace(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(schemas[0]); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	save := func(owner, name, list, notes string) Place {
		p, status, e := persistPlace(ctx, db, owner, Place{Name: name, Lat: 51.5, Lon: -.1, List: list, Notes: notes})
		if status != 200 || e != nil {
			t.Fatal(status, e)
		}
		return p
	}
	old := save("a", "Previous home", " home ", "keep this note")
	other := save("b", "Other account home", "Home", "other note")
	work := save("a", "Office", "Work", "")
	next := save("a", "New home", "HOME", "new note")
	var list, note string
	if err = db.QueryRow("SELECT list,notes FROM places WHERE id=?", old.ID).Scan(&list, &note); err != nil || list != "Favourites" || note != "keep this note" {
		t.Fatal("old place or note lost", list, note, err)
	}
	if err = db.QueryRow("SELECT list FROM places WHERE id=?", other.ID).Scan(&list); err != nil || list != "Home" {
		t.Fatal("changed another account")
	}
	var homes int
	db.QueryRow("SELECT count(*) FROM places WHERE owner='a' AND list='Home'").Scan(&homes)
	if homes != 1 {
		t.Fatal("ambiguous home", homes)
	}
	next.List = "Work"
	if _, status, err := persistPlace(ctx, db, "a", next); status != 200 || err != nil {
		t.Fatal(status, err)
	}
	db.QueryRow("SELECT list FROM places WHERE id=?", work.ID).Scan(&list)
	if list != "Favourites" {
		t.Fatal("work not replaced safely")
	}
	foreign := other
	foreign.List = "Work"
	if _, status, err := persistPlace(ctx, db, "a", foreign); status != 404 || err == nil {
		t.Fatal("foreign update accepted")
	}
	db.QueryRow("SELECT list FROM places WHERE id=?", next.ID).Scan(&list)
	if list != "Work" {
		t.Fatal("failed update changed shortcut")
	}
}
func TestListUpdateRollbackAndCustomNotes(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.Exec(schemas[0])
	ctx := context.Background()
	old, status, err := persistPlace(ctx, db, "a", Place{Name: "Saved", Lat: 0, Lon: 0, List: "Home", Notes: "keep"})
	if status != 200 || err != nil {
		t.Fatal(err)
	}
	db.Exec("CREATE TRIGGER deny_new_place BEFORE INSERT ON places WHEN NEW.name='Failed' BEGIN SELECT RAISE(ABORT,'synthetic failure'); END")
	if _, status, err = persistPlace(ctx, db, "a", Place{Name: "Failed", Lat: 1, Lon: 1, List: "Home"}); status != 500 || err == nil {
		t.Fatal("failure ignored")
	}
	var list, note string
	db.QueryRow("SELECT list,notes FROM places WHERE id=?", old.ID).Scan(&list, &note)
	if list != "Home" || note != "keep" {
		t.Fatal("failed save displaced existing Home")
	}
	old.List = " Weekend trips "
	old.Notes = "updated"
	if _, status, err = persistPlace(ctx, db, "a", old); status != 200 || err != nil {
		t.Fatal(err)
	}
	db.QueryRow("SELECT list,notes FROM places WHERE id=?", old.ID).Scan(&list, &note)
	if list != "Weekend trips" || note != "updated" {
		t.Fatal(list, note)
	}
}
