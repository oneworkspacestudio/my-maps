package app

import (
	"errors"
	"github.com/oneworkspacestudio/my-maps/internal/directions"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"net/http"
)

func (a *App) directionsStatus(w http.ResponseWriter, r *http.Request) {
	apphost.JSON(w, 200, a.directions.Status())
}
func (a *App) directionsPlan(w http.ResponseWriter, r *http.Request) {
	var choice directions.Choice
	if !apphost.Decode(w, r, &choice) {
		return
	}
	plan, err := a.directions.Plan(r.Context(), choice)
	if err != nil {
		apphost.Problem(w, 409, err.Error())
		return
	}
	apphost.JSON(w, 200, plan)
}
func (a *App) directionsInstall(w http.ResponseWriter, r *http.Request) {
	a.installMu.Lock()
	defer a.installMu.Unlock()
	a.maps.mu.RLock()
	busy := a.maps.progress.Busy
	a.maps.mu.RUnlock()
	if busy {
		apphost.Problem(w, 409, "Wait for the map download to finish first")
		return
	}
	var in struct {
		Plan string `json:"plan"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	if err := a.directions.Install(in.Plan); err != nil {
		apphost.Problem(w, 409, err.Error())
		return
	}
	apphost.JSON(w, 202, a.directions.Status())
}
func (a *App) directionsCancel(w http.ResponseWriter, r *http.Request) {
	a.directions.Cancel()
	apphost.JSON(w, 200, a.directions.Status())
}
func (a *App) localRoute(w http.ResponseWriter, r *http.Request, in directions.Request) {
	result, err := a.directions.Route(r.Context(), in)
	if err != nil {
		status, code := 422, "route_unavailable"
		if errors.Is(err, directions.ErrNoData) {
			status, code = 503, "directions_download"
		}
		if errors.Is(err, directions.ErrCoverage) {
			code = "directions_coverage"
		}
		apphost.JSON(w, status, map[string]string{"error": err.Error(), "code": code})
		return
	}
	apphost.JSON(w, 200, result)
}

func (a *App) mapInstall(w http.ResponseWriter, r *http.Request) {
	a.installMu.Lock()
	defer a.installMu.Unlock()
	if a.directions.Status().Progress.Busy {
		apphost.Problem(w, 409, "Wait for the directions download to finish first")
		return
	}
	a.maps.install(w, r)
}
