package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"strings"
)

func canonicalList(name string) string {
	name = strings.TrimSpace(name)
	switch strings.ToLower(name) {
	case "", "favourites", "favorites":
		return "Favourites"
	case "home":
		return "Home"
	case "work":
		return "Work"
	}
	return name
}

// A shortcut has one location per account. Replacement keeps the old place and
// its notes in Favourites; the entire change commits with the new saved place.
func persistPlace(ctx context.Context, db *sql.DB, owner string, p Place) (Place, int, error) {
	p.List = canonicalList(p.List)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return p, 500, err
	}
	defer tx.Rollback()
	if p.ID == "" {
		var count int
		if err = tx.QueryRowContext(ctx, "SELECT count(*) FROM places WHERE owner=?", owner).Scan(&count); err != nil {
			return p, 500, err
		}
		if count >= 10000 {
			return p, 409, errors.New("Export or remove some saved places before adding more")
		}
		p.ID = apphost.ID()
		_, err = tx.ExecContext(ctx, "INSERT INTO places VALUES(?,?,?,?,?,?,?,?)", p.ID, owner, p.Name, p.Address, p.Lat, p.Lon, p.List, p.Notes)
	} else {
		var result sql.Result
		result, err = tx.ExecContext(ctx, "UPDATE places SET name=?,address=?,lat=?,lon=?,list=?,notes=? WHERE id=? AND owner=?", p.Name, p.Address, p.Lat, p.Lon, p.List, p.Notes, p.ID, owner)
		if err == nil {
			n, e := result.RowsAffected()
			if e != nil {
				return p, 500, e
			}
			if n != 1 {
				return p, 404, errors.New("Saved place not found")
			}
		}
	}
	if err != nil {
		return p, 500, err
	}
	if p.List == "Home" || p.List == "Work" {
		if _, err = tx.ExecContext(ctx, "UPDATE places SET list='Favourites' WHERE owner=? AND id<>? AND lower(trim(list))=lower(?)", owner, p.ID, p.List); err != nil {
			return p, 500, err
		}
	}
	if err = tx.Commit(); err != nil {
		return p, 500, err
	}
	return p, 200, nil
}
