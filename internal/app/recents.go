package app

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

// Only explicit searches and chosen places enter this account's bounded history.
// Never store private place notes, provider responses or passive location fixes.
type Recent struct {
	Query   string   `json:"query"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Lat     *float64 `json:"lat,omitempty"`
	Lon     *float64 `json:"lon,omitempty"`
}

func (a *App) recents(w http.ResponseWriter, r *http.Request) {
	rows, err := a.H.DB.Query("SELECT body FROM recent_places WHERE owner=? ORDER BY updated DESC,entry_key LIMIT 20", apphost.Current(r).User.ID)
	if err != nil {
		a.H.DBError(w, err)
		return
	}
	defer rows.Close()
	result := []Recent{}
	for rows.Next() {
		var body []byte
		var item Recent
		if err = rows.Scan(&body); err != nil {
			a.H.DBError(w, err)
			return
		}
		if err = json.Unmarshal(body, &item); err != nil {
			a.H.DBError(w, err)
			return
		}
		result = append(result, item)
	}
	if err = rows.Err(); err != nil {
		a.H.DBError(w, err)
		return
	}
	apphost.JSON(w, http.StatusOK, result)
}

func (a *App) remember(w http.ResponseWriter, r *http.Request) {
	var item Recent
	if !apphost.Decode(w, r, &item) {
		return
	}
	item.Query, item.Name, item.Address = strings.TrimSpace(item.Query), strings.TrimSpace(item.Name), strings.TrimSpace(item.Address)
	if !validRecent(item) {
		apphost.Problem(w, 400, "Check the recent search or place")
		return
	}
	key := apphost.Hash(strings.ToLower(item.Query + "\n" + item.Name + "\n" + item.Address))
	body, _ := json.Marshal(item)
	owner := apphost.Current(r).User.ID
	tx, err := a.H.DB.Begin()
	if err != nil {
		a.H.DBError(w, err)
		return
	}
	defer tx.Rollback()
	if _, err = tx.Exec("INSERT OR REPLACE INTO recent_places VALUES(?,?,?,?)", owner, key, body, time.Now().UnixMilli()); err != nil {
		a.H.DBError(w, err)
		return
	}
	if _, err = tx.Exec("DELETE FROM recent_places WHERE owner=? AND entry_key NOT IN (SELECT entry_key FROM recent_places WHERE owner=? ORDER BY updated DESC,entry_key LIMIT 20)", owner, owner); err != nil {
		a.H.DBError(w, err)
		return
	}
	if err = tx.Commit(); err != nil {
		a.H.DBError(w, err)
		return
	}
	apphost.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func validRecent(item Recent) bool {
	if len(item.Query) > 200 || len(item.Name) > 200 || len(item.Address) > 1000 || (item.Lat == nil) != (item.Lon == nil) {
		return false
	}
	if item.Lat != nil {
		return item.Name != "" && item.Query == "" && coordinate(*item.Lat, *item.Lon)
	}
	return len(item.Query) >= 2 && item.Name == "" && item.Address == ""
}

func (a *App) clearRecents(w http.ResponseWriter, r *http.Request) {
	if _, err := a.H.DB.Exec("DELETE FROM recent_places WHERE owner=?", apphost.Current(r).User.ID); err != nil {
		a.H.DBError(w, err)
		return
	}
	apphost.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}
