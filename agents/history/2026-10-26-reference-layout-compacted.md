# Maps reference layout — 26 October 2026

User requests matching supplied Google Maps screens/features, push, release and
installation on both PCs. Read repository-history, Maps/Workspace AGENTS and
prior Maps2 delivery. Source I:/Projects/my-maps was clean at b35e327; branch
codex/google-maps-layout. Current owner installation is C:/tmp/workspace, not the
older other-host I:/Workspace handoff. Both PCs currently manager92; second known
installation at workspace-2.m8e.co.uk has an authenticated browser updater.
SSH m8e-pc failed auth and is not required for this established second endpoint.

AUI screen redesign: full map, floating search/chips/profile, responsive place
sheet with pointer/keyboard control, Explore/You/Contribute, layers dialog and
place/route UI. MapLibre6.15 replaces Leaflet; vector colours, raised buildings,
Mapterhorn terrain and NASA regional imagery. Nearby Overpass query is bounded,
rate limited/cached; loaded vector POIs provide a notice-labelled fallback.
Fixed-provider authenticated resource proxy retains same-origin CSP, bounds bytes
and redirects, and forwards no cookies/credentials. Existing XYZ provider preserved.
Weather/AQI use Open-Meteo, 15min cache and explicit user actions. Private places
schema/identity ownership retained. No Actions, delegation, paid services or Google
scraping. docs/FEATURES.md distinguishes all capabilities from missing Google data.

Development Go tests/vet, AUI/TS/Vite and Windows build pass. Native synthetic QA
passed actual Photon/OSRM, private place save, precise pin, GPX, tabs/profile,
Escape/focus, bounded layers and 360/390/576/1440-width rendering. Overpass failed
network/TLS/timeout in this environment; fallback added and final QA pending.
No full Google parity, real-device GPS/mic, background navigation, clean-host or
reboot acceptance. Satellite is NASA500m; Street View/fire maps are external links;
traffic/timetables/reviews/photos/Lens remain unavailable. Never claim otherwise.

Artifacts I:/Projects/tooling/maps-redesign-20261026. Only disposable protected
fixture C:/Projects/.cache/maps-redesign-20261026; process receipts track preview
programs/PIDs. Final release3/package/public-channel promotion, normal Maps-only
backed-up updates and verified receipts are pending. Revalidate dynamic installed
versions before applying because Tasks/Drive/AI tasks are active independently.

Final native QA confirms real loaded-map nearby discovery after zooming to POI
level, duplicate removal/cap60 and separated pins. The last mobile Settings check
found an inherited hidden Back link; final CSS restores that explicit link and
scrolling. Exact packaged binary/asset and complete final QA are required before
publication; do not claim installation from this source checkpoint.
