package app

import (
	"errors"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

var schemas = []string{
	`CREATE TABLE places(id TEXT PRIMARY KEY,owner TEXT NOT NULL,name TEXT NOT NULL,address TEXT NOT NULL,lat REAL NOT NULL,lon REAL NOT NULL,list TEXT NOT NULL,notes TEXT NOT NULL);
 CREATE INDEX places_owner ON places(owner,list);
 CREATE TABLE provider_cache(key TEXT PRIMARY KEY,body BLOB NOT NULL,expires INTEGER NOT NULL);`,
	`CREATE TABLE recent_places(owner TEXT NOT NULL,entry_key TEXT NOT NULL,body BLOB NOT NULL,updated INTEGER NOT NULL,PRIMARY KEY(owner,entry_key));
 CREATE INDEX recent_places_owner ON recent_places(owner,updated);`,
}

func migrate(h *apphost.Host) error {
	var current int
	if err := h.DB.QueryRow("SELECT coalesce(max(version),0) FROM host_migrations").Scan(&current); err != nil {
		return err
	}
	if current > len(schemas) {
		return errors.New("map data needs a newer executable; restore a compatible backup to downgrade")
	}
	for version := current + 1; version <= len(schemas); version++ {
		if err := h.Migrate(version, schemas[version-1]); err != nil {
			return err
		}
	}
	return nil
}
