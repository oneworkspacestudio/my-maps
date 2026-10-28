# Map capabilities and install depth

| Feature | Local implementation |
| --- | --- |
| Basemap | Bundled Leaflet + Protomaps Canvas renderer, local SQLite MVT tiles, no WebGL required |
| Map appearance | Pale land, mint parks, cyan water, outlined blue-grey main roads, white local streets, buildings and zoom-aware street/park/place labels |
| Depth | Overview z0–6; regional z0–10; streets z0–14 with overzoom to19 |
| Coverage | World overview; UK/Europe overview or regional; user-chosen local area at all depths |
| Search | Local SQLite FTS index of mapped names/roads, GeoNames towns, optional UK postcodes |
| Nearby | Named POIs in the installed maximum-detail tiles; truthful empty results when absent |
| Place sheet | Compact after search; manual pointer drag or keyboard/tap handle reveals details |
| Private places | Save-to-list dropdown on places/results, Favourites/Home/Work/custom lists, private notes, deletion and GeoJSON export |
| Search history | Account-isolated recent queries/places, bounded20, Clear all |
| Share | Place name/coordinates or bounded route stops and mode; private notes excluded |
| Directions | Integrated local BRouter downloads; drive/walk/cycle; car/EV/hybrid and bicycle/hiking profiles; avoid tolls/motorways/unpaved; up to8 stops, alternatives/turns/GPX |
| Location | Explicit device request, foreground next-turn and remaining-distance navigation, explicit off-route update; no background navigation or voice guidance |
| Removed | Public imagery/terrain/3D, weather/AQI, traffic, Street View, voice/Lens, reviews/photos/ratings; no transit, battery range, charging stops or heavy-vehicle restrictions |

## Installing or changing data

Only the server owner can prepare/start/cancel a download or configure routing.
All signed-in users can browse the installed pack. Choose a download at first run
or from the download icon/account menu. Choose coverage, depth and optional UK
postcodes; check the plan, then explicitly Install. Changing a choice invalidates
the previous plan. Plans expire after20 minutes; restarting the app requires a new
plan. The displayed transfer size covers selected compressed tiles and GeoNames
ZIPs; free space is a conservative allowance for SQLite/indexing and safe staging.
A one-tile margin at each zoom avoids blank edges around the requested bounds.
Local-area radius defines a surrounding bounding box. Beyond the downloaded
coverage, additional streets/places are unavailable; choose a larger area to expand.
World street detail and an entire-country street graph are not bundled.

The installer reads bounded PMTiles v3 ranges from a recent Protomaps daily build
and downloads fixed GeoNames ZIP sources. It checks exact byte ranges, lengths,
stable ETags, archive/expansion limits and staging storage. It records a local
pack SHA256. This subset is not a verification of the planet archive's full-file
BLAKE3 digest. ZIPs are streamed without extracting filesystem paths. The new pack
is staged, indexed and activated atomically; failure/cancellation keeps the old
pack. A completed replacement removes only its previous active map pack. Settings,
saved places and recents stay in the app database. Native app backups preserve accounts, saved places and app objects. Downloaded
map packs and directions data are separate, replaceable files; reinstall coverage
after restoring into a fresh folder. Existing coverage remains during app updates.

No coverage download runs automatically during a software update. An installation
without local data shows the chooser rather than a blank map or public fallback.
Download progress/cancel is available in Map downloads. An interrupted staging file
is cleaned on restart. Refresh the plan and retry after a failed source download.
Data is a snapshot; install a fresh selection to update it. There is no automatic
planet sync, transit timetable, street-photo source or arbitrary-download URL.

## Optional local services

The map chooser includes local directions by default for UK/local coverage. It
shows the combined download/space estimate, installs the map first and then routing
while this page remains open. Reopening Map downloads offers a separate directions
installer if the page was closed before the second download started. Two installers
cannot run concurrently. A plan is refreshed after a long map install; a larger
routing download requires a new size check. Failed/cancelled downloads retain the
previous active data.

Directions downloads SHA-pinned BRouter and Temurin runtime archives plus bounded
OSM-derived segments. Strong ETags/lengths, archive path/expansion checks and graph
CRC checks protect staging; activation switches a local pointer only after private
engine startup passes. Successful replacement removes the previous directions pack.
The local bridge accepts only authenticated loopback route/cancel/health calls;
no upstream profile uploads, suspects/recent-waypoint routes, route logs or public
routing service are exposed. Process tokens pass through stdin and Windows job
supervision stops the child with Maps. Car options affect actual routing profiles.
Generic turn instructions do not announce road names; times have no live traffic.

Legacy literal loopback OSRM URLs remain supported by the existing routing API and
MYMAPS_ROUTER setting for standard driving; other modes/options require the bundled
local engine. Directions and Start remain visible, with in-app setup when needed.

Advanced administrators may set `MYMAPS_GEOCODER` to a local Photon base,
`MYMAPS_ROUTER` to a local OSRM base or `MYMAPS_TILES` to a local XYZ raster URL.
Only literal loopback IP URLs are accepted (no credentials). The attribution env
`MYMAPS_TILE_ATTRIBUTION` accompanies custom tiles. Public/default external
providers and `MYMAPS_OVERPASS` are removed. Legacy remote provider configuration
must be changed to a local service or unset before startup. Redirects cannot leave
the configured local origin. Browsers only contact the app origin for map assets.
A browser still needs access to this server; this is server-hosted data, not a
service-worker map download to a disconnected phone.

Database schema3 upgrades saved places/history from schemas1/2 without replacing
owner records. An older executable refuses a newer schema; use the normal compatible
backup recovery rather than downgrading the executable in place. The managed app
keeps four registered program/license packages and a minimum Workspace92 contract.

## Saving and Home/Work shortcuts

Save opens a native select picker without changing the place sheet or writing a
record. Favourites is the default for new places; existing places keep their list
and notes. Nearby/search/recent locations have independent save controls. An unset
Home or Work shortcut starts search and preselects that list after choosing a place.
No free-text list editing is needed for either shortcut. New list names are the
only list text entry. Confirm writes the list/note together; cancel keeps all data.

Home and Work each resolve to one location per account. Replacement is transactional:
the previous place moves to Favourites with its notes intact. Other accounts are
unaffected. Existing custom lists remain selectable. This remains one list per
saved place, not multiple simultaneous memberships. No new data schema is needed.
