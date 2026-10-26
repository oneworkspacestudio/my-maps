# Map capabilities

The phone layout follows the supplied reference: floating search and category
chips, a full-screen map, right-side controls, a draggable/clickable place sheet,
and Explore, You and Contribute tabs. Desktop uses a floating place panel.
My Maps keeps its own name and branding.

| Reference feature | My Maps implementation |
| --- | --- |
| Search, pins and place details | Photon address search, exact coordinates, open map labels and nearby POIs |
| Nearby discovery | Overpass details; loaded-map fallback during provider outages; no invented ratings/photos |
| Your places | Private lists, Home/Work list shortcuts, notes and GeoJSON export |
| Directions | Drive/walk/cycle, up to eight stops, alternatives, steps, GPX, foreground GPS following |
| Default map | OpenStreetMap image tiles on phones; OpenFreeMap vector map on desktops; existing custom tile endpoint retained |
| Satellite | NASA Blue Marble regional 500 m imagery, not street-resolution satellite |
| Terrain | Mapterhorn elevation/hillshade with open vector cartography |
| Raised buildings | Mapped building heights and tilted view where available |
| Public transport | Nearby stops/stations; no timetables or live departures |
| Cycling | Nearby cycle parking/hire/repair; no cycle-network overlay |
| Traffic | Availability explanation; no live congestion provider |
| Street View | Coordinate-specific link to Mapillary, opens separately |
| Wildfires | Coordinate-specific link to NASA FIRMS, opens separately |
| Air quality | Open-Meteo/CAMS model values with timestamp; no map-wide overlay |
| Weather | User-requested local temperature and timestamp |
| Voice search | Browser speech recognition where available; submit captured search text |
| Lens | Not implemented; the scan control opens exact-coordinate search |
| Contribute | OpenStreetMap edit/report links and private saved-place notes |

Google photos, reviews, ratings, Lens and Google data services are not cloned or
scraped. No paid account, offline map download or background navigation is enabled.
MapLibre and open data providers retain their own licences and attribution.

Phones use a DOM/image-tile renderer that does not require WebGL. Desktop maps,
terrain and raised buildings use WebGL2, with automatic standard-map recovery
if graphics initialization fails, the context is lost or no map renders within
12 seconds. Search, pins, route lines, measurements and foreground GPS also work
on the standard map. Loaded-vector POI discovery is available only in vector mode.
