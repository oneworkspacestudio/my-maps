# Map-depth installer and local geographic data

User authorizes per-user choice of installation depth, push/release and both-PC
install. Earlier UI references are layout only: no Google data, fully self-hosted,
remove unprovisioned features. Shared server owner chooses coverage for all users;
private places/recents remain account-isolated. This supersedes the prior unanswered
coverage question in self-hosting/mobile-search handoffs.

Candidate implementation: Overviewz6 (world/UK/Europe/area), regionalz10
(UK/Europe/area), streetsz14 (area1–100km). Owner plans exact compressed transfer
and conservative free-space allowance, then explicitly installs/cancels. Real
bounded Protomaps PMTiles ranges + GeoNames cities15000/optional full UK postcode
ZIP; stable ETag/length/range/expansion/metadata limits, protected staging SQLite,
disk reserve/max page budget, atomic activation and previous-pack retention on
cancel/failure. Named staging files cleaned after interruption. No automatic owner
data download on app update. Fresh unconfigured servers show the chooser.

Tiles/search/nearby are local, bundled CPU Canvas Leaflet + Protomaps; no WebGL or
public browsing fallback, no CDN/fonts. Local MVT names/roads/POIs + GeoNames FTS,
not full house-number geocoding. Named POIs depend on depth/source coverage. Index
bounds filter while tile selection includes one-tile margin. Snapshot updates are
explicit reinstalls. Attribution includes OSM/GeoNames/RoyalMail where applicable.
No arbitrary source URL/path. Map packs increase normal encrypted backup size.

Default public Photon/FOSSGIS/Overpass/OSM/OpenFreeMap/NASA/terrain removed.
Advanced providers require literal loopback URLs. Optional local OSRM driving
connection is probed; Directions/Start hidden until configured, no unprovisioned
walking/cycling modes. No routing graph installer. Standard OSRM route endpoint.
Map downloads/Settings/Back to map navigation fixed. AUI source/screens and generated
React retained; no MapLibre dependency. Earlier compact-sheet/search/schema2 work
retained; app schema3 adds active-pack/routing metadata, old private data preserved.

Verification: Go tests/vet + AUI/TS/Vite/native pass. Protected disposable C NTFS
fixture downloaded real2km London street pack+UK postcodes. Actual E46NW search,
restaurants/shopping and all six main categories valid, no invented tags/photos,
compact sheet and share/save, restart, CSRF/member privacy/concurrent start/cancel,
atomic replacement203tiles/2471features checked. Browser WebGL off, external
browser traffic denied/none attempted, no exceptions. World overview5461tiles
~48.5MB and UK regional2853tiles~74.9MB plans succeeded (not installed). Separate
synthetic OSRM contract test. No owner records as fixtures/physical-phone claim.
See docs/evidence/map-depth-20261027.json. Artifact root is task tooling directory.

Release9 source version bumped; immutable Maps8 is previous delivered build,
Maps7 remains withheld. Publication/channel and actual owner update are still
pending at this source commit. Use four registered packages/min92, local Go1.27.3,
Actions disabled, immutable17asset release + official publisher. Dedicated channel
checkout must retain latest pins/withholdings. Primary now manager94/MyAI15/Maps8,
28services, manager bytes public94 SHA1936c1e8b96e2e0cb51fe713c687dd961399fd8da1eae6cb9df1ff48278ad76e.
Public94 source ZIP verifiedbecb7fc1ebe783e554da09564657e5a42b559ae368993cac3bce84f3e21b88e3;
source under tooling/sdk94-source/one-workspace. Adapted normal protected Maps-only
updater helper; do not use old93 controller on primary. Preserve all other apps,
owner config/identities/paths/models/latest recovery. Secondary browser sign-in
expired; fresh MyAccounts tab is marked handoff and user asked asynchronously to
sign in. No owner update submitted yet. Never bypass sign-in or use owner as QA.

Only active owned fixture: C:/Projects/.cache/maps-depth-20261027, receipt
I:/Projects/tooling/maps-redesign-20261026/depth-process.json tracks exact binary,
PID/data/port19315; guarded restart script. Remove this fixture after exact-package
QA; no older dead PIDs reused. Keep source/release/evidence on I.

Packaging guard correctly refused missing licences in upstream npm tarballs for
@protomaps/basemaps5.7.2 and pmtiles3.2.1. Exact npm gitHead licence texts are
vendored under docs/licenses with source URLs/SHA256, including style/schema
attribution; version/digest-checked build fallback includes them in notices.
