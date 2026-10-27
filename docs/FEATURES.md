# Map capabilities

The phone layout follows the supplied reference: floating search and category
chips, a full-screen map, right-side controls, a draggable/clickable place sheet,
and Explore, You and Contribute tabs. Desktop uses a floating place panel.
My Maps keeps its own name and branding.

| Reference feature | My Maps implementation |
| --- | --- |
| Search, pins and place details | Photon address search, exact coordinates, open map labels and nearby POIs |
| Nearby discovery | Photon category search on both renderers, Overpass backup and loaded-vector fallback; no invented ratings/photos |
| Your places | Private lists, Home/Work list shortcuts, notes and GeoJSON export |
| Directions | Drive/walk/cycle, up to eight stops, alternatives, steps, GPX, foreground GPS following |
| Default map | OpenStreetMap image tiles on phones; OpenFreeMap vector map on desktops; existing custom tile endpoint retained |
| Satellite | NASA Blue Marble regional 500 m imagery, not street-resolution satellite |
| Terrain | Mapterhorn elevation/hillshade with open vector cartography |
| Raised buildings | Mapped building heights and tilted view where available |
| Transport stops | Nearby stops/stations |
| Cycle facilities | Nearby cycle parking/hire/repair |
| Traffic | Removed: no self-hosted traffic service |
| Street View | Removed: external-service shortcut |
| Wildfires | Removed: external-service shortcut |
| Air quality | Removed: external-service dependency |
| Weather | Removed: external-service dependency |
| Voice search | Removed: browser speech can depend on a cloud provider |
| Coordinates | Exact latitude/longitude lookup with a location icon |
| Contribute | OpenStreetMap edit/report links and private saved-place notes |

Google photos, reviews, ratings, Lens and Google data services are not cloned or
scraped. No paid account, offline map download or background navigation is enabled.
MapLibre and open data providers retain their own licences and attribution.

Phones use a DOM/image-tile renderer that does not require WebGL. Desktop maps,
terrain and raised buildings use WebGL2, with automatic standard-map recovery
if graphics initialization fails, the context is lost or no map renders within
12 seconds. Search, pins, route lines, measurements and foreground GPS also work
on the standard map. Loaded-vector POI discovery is available only in vector mode.

Searches, nearby categories and selected places keep the map visible with the
sheet collapsed. Drag or tap its handle to open details; result arrival preserves
a sheet the user has already opened. The title and count identify current results.
Public-provider coverage varies; failures show a retry state and empty results
are distinguished from an unavailable service. Nearby responses are cached for
one hour and provider requests are spaced rather than rejected for quick clicks.

The source cleanup above is not released yet. Public map/geocoder/router defaults
still exist; local tiles, geographic index and routing data are required before
claiming a fully self-hosted installation. No Google data is used by the Maps
provider implementation, and Maps does not call My AI's Google-enabled SearXNG.

## Mobile search and place revision (working source)

The supplied phone references and the [Google Maps Android search/list flow](https://support.google.com/maps/answer/9948049?co=GENIE.Platform%3DAndroid&hl=en) inform the layout only; Maps does not call a Google data API.
Home uses a 112px peek above 64px navigation at phone sizes. Categories use a 136px
peek; selected places use 210px with navigation hidden. Directions, Start, Share
and Save are immediately available; dragging reveals private notes and details.
Address search selects its first matching place in the collapsed view; choosing
a suggestion, recent or saved place also opens that compact view. A manual
expansion during a pending search is retained. Giant decorative result art is
replaced by compact place rows. Search has an autofocus editor, Home/Work,
local saved matches, provider suggestions, private bounded recents and Clear all.
Native browser Back closes search. Visual viewport changes reserve keyboard space.
Start requests location explicitly, routes from it, then follows in the foreground;
permission failures explain manual Directions. This is not voice guidance.

Recent history is stored per immutable app account in schema 2, capped at 20 entries,
with validated query/place kinds and native Origin/CSRF authorization. Private notes
and passive device locations are excluded. Upgrade/restart preserve saved places;
a schema downgrade requires the normal compatible backup. Public-provider defaults
still remain: the separate fully-local map-data provisioning is not complete.
