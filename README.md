<p align="center"><img src="web/public/logo.svg" alt="My Maps" width="112"></p>

# My Maps

Explore the world from One Workspace. A portable Windows **Maps.exe** with the
Workspace design system, MyAccounts sign-in, saved places and open-data directions.
Built on [Leaflet](https://leafletjs.com/) and OpenStreetMap, without a paid API key.

![My Maps desktop](docs/assets/desktop.png)

## Use it

Download the Windows ZIP from [Releases](https://github.com/oneworkspacestudio/my-maps/releases),
extract it into a dedicated NTFS folder and run **Maps.exe**. Its browser opens a
one-time local owner setup. No Docker, Node or Go installation is needed.

One Workspace installs the app automatically at `/maps/`, configures MyAccounts,
and stores its private state in the Workspace data folder under `maps`.
An existing Workspace login starts the normal authenticated SSO flow automatically.
Standalone accounts are disabled in managed mode.

## Included

- Interactive street map, zoom, location and dropped pins.
- Place/address search and reverse lookup, only when requested.
- Private saved lists, notes, favourites and deletion.
- Share links containing a place name and coordinates, without private notes.
- Driving, cycling and walking directions, up to eight stops, alternatives,
  distance/time estimates and written route steps.
- Follow your location and measure distances.
- GeoJSON saved-place and GPX route downloads.
- Light/dark theme, mobile layout, Workspace launcher, encrypted backups.

## Open services and limits

The default map uses [OSM tiles](https://operations.osmfoundation.org/policies/tiles/):
attribution stays visible, browser caching is honoured, and there is no offline
prefetch or bulk download. [Photon](https://github.com/komoot/photon#demo-server)
provides community geocoding; [FOSSGIS OSRM](https://routing.openstreetmap.de/about.html)
provides routes. Search and route requests are cached for 24 hours and share a
server-wide maximum of one upstream request per second. These community services
have no availability guarantee and can restrict heavy usage.

Search text and route coordinates go to the chosen providers. Location is requested
only after pressing My location / Follow my location; saved places stay on your
server. Do not submit confidential addresses to public providers.

This is an initial functional preview, **not full Google Maps parity**. It does not
include live traffic, Street View, satellite imagery, transit timetables, proprietary
reviews, background navigation, voice guidance, offline map packages or Google data.
GPS following is not turn-by-turn navigation. Routes can be incomplete or outdated.

Administrators can select compatible HTTPS endpoints using `MYMAPS_TILES`
(XYZ URL), `MYMAPS_GEOCODER` (Photon base URL), and `MYMAPS_ROUTER`
(FOSSGIS-compatible base with `routed-car`, `routed-bike`, `routed-foot`).
Use a provider with suitable capacity or host the services yourself for larger use.

## Build

Go 1.27.2, Node 26, PowerShell 7:

```powershell
./scripts/build.ps1 -Go go -Npm npm.cmd -Output ./dist
```

Screens are authored in `web/src/screens/*.aui`; `npm run build` regenerates React.
Release binaries are built locally, with dependency licences and source receipts.
No GitHub Actions. Application code is Apache-2.0; see `NOTICE`, `SECURITY.md` and the
third-party licence bundle. This project is not affiliated with Google.
