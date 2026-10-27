# Map capabilities and install depth

| Feature | Local implementation |
| --- | --- |
| Basemap | Bundled Leaflet + Protomaps Canvas renderer, local SQLite MVT tiles, no WebGL required |
| Depth | Overview z0–6; regional z0–10; streets z0–14 with overzoom to19 |
| Coverage | World overview; UK/Europe overview or regional; user-chosen local area at all depths |
| Search | Local SQLite FTS index of mapped names/roads, GeoNames towns, optional UK postcodes |
| Nearby | Named POIs in the installed maximum-detail tiles; truthful empty results when absent |
| Place sheet | Compact after search; manual pointer drag or keyboard/tap handle reveals details |
| Private places | Save-to-list dropdown on places/results, Favourites/Home/Work/custom lists, private notes, deletion and GeoJSON export |
| Search history | Account-isolated recent queries/places, bounded20, Clear all |
| Share | Place name and coordinates only; private notes excluded |
| Driving | Optional loopback OSRM driving service, up to8 stops, alternatives/steps/GPX |
| Location | Explicit device request, foreground following; no background navigation or voice guidance |
| Removed | Public imagery/terrain/3D, weather/AQI, traffic, Street View, voice/Lens, reviews/photos/ratings, walking/cycling modes without local graphs |

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
saved places and recents stay in the app database. Encrypted/native backups include
installed map packs, so larger selections also increase backup size.

No coverage download runs automatically during a software update. An installation
without local data shows the chooser rather than a blank map or public fallback.
Download progress/cancel is available in Map downloads. An interrupted staging file
is cleaned on restart. Refresh the plan and retry after a failed source download.
Data is a snapshot; install a fresh selection to update it. There is no automatic
planet sync, transit timetable, street-photo source or arbitrary-download URL.

## Optional local services

In Map downloads, connect a literal loopback OSRM URL, e.g.
`http://127.0.0.1:5000`. A compatible nearest-driving response is required before
saving it. A driving graph must already exist on this server; installing and
building arbitrary routing graphs is outside this map-pack installer. Route API
is standard `/route/v1/driving/…`. Directions/Start and the drive mode appear only
when configured; no public routing fallback or synthetic journey is supplied.
Disconnect hides them. Location following runs only while the page is open.

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
