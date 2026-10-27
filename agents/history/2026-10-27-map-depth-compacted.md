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

## Published and primary installed; secondary sign-in pending

Release9 immutable source35cca4a136821a7210ac9b9a948d770087aa42a4, all17 public
asset sizes/digests and anonymous EXE verified. Accepted/public/primary EXE SHA
b2352e87b92a75af677c33f0a1024eae1bf94901cc56f04e983b4192c17d9ab4.
Official channel promotionefca08a is anonymously visible on main; other pins,
historical rows and withholdings preserved. Exact native packaged bytes pass12
real-data groups, actual download/cancel/replacement plus browser/local assertions.
Public92/93/94 manifest parsers pass. No Actions. Missing upstream npm licence
texts are pinned to exact npm gitHead with verified bytes and -text attributes;
packager keeps licence metadata separate from commit identity. Final build path
is build-release9-ready; older failed/preliminary staging must not be published.

Primary normal Maps-only8->9 protected update completed with public94 controller.
29catalog apps/28running services, manager94/MyAI15 and every other app contract,
package and executable byte preserved; owner configuration and decrypted identity
settings compared only in memory and matched. Newest verified recovery:
C:/tmp/workspace/Private/backups/update-20261027T164148-425bd09170cb.
Runtime health/version9, exact EXE and served JS/CSS hashes/title verified. No
owner functional fixtures or arbitrary country-data download performed. Users
choose map coverage/depth after opening Maps. No local OSRM graph is bundled.

Secondary still presents normal MyAccounts sign-in on the freshly navigated
workspace-2/updates tab585734409 in Chrome3. It is marked handoff. User asked to
sign in asynchronously; no answer received/no secondary apply submitted. Last
known secondary manager93/Maps8/MyAI14 is historical, not fresh installed evidence.
After user signs in, inspect fresh updater versions/state, check published updates,
select only Maps9 and use normal backed-up install. Preserve other pending app
updates, do not upgrade manager/MyAI implicitly or bypass sign-in. If state expired,
start at fresh workspace-2 root/updates, never reuse OAuth callbacks/tokens.
Primary installation is complete; do not repeat it. Additional owner live Maps UI
checks need authorized sign-in; runtime/accepted package proofs already pass.

Exact owned protected C fixture maps-depth-20261027/PID/port19315 was stopped and
removed with path/executable/reparse guards. Sanitized receipts remain in tooling
and docs/evidence/map-depth-20261027.json. No preview process is left. Continue
secondary installation only when normal authentication becomes available.

Advanced follow-up outside default/UI connection path: a routing service supplied
only via MYMAPS_ROUTER before the first pack can be superseded by the empty routing
field inserted with that first pack on restart. Persisted UI router connections are
retained. Keep this edge case explicit if extending administrator env seeding.
