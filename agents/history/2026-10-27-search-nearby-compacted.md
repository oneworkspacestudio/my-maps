# Map-first search and nearby reliability

User requests search results remain on map, manually dragged detail sheet, working
restaurants/shopping/etc, all-feature checks, push/release/install both PCs.
Source clean6a13457 on codex/google-maps-layout; installed Maps6/manager92 per
prior verified delivery. Load mobile-map-recovery/reference-layout handoffs for
immutable releases, canonical four-package/min92 scope and protected update flow.

Found search/category/choose explicitly expanding the sheet. Now collapse once
at request start or place selection, release keyboard focus, retain any manual
expansion while the request runs, expose result title/count/drag hint. AUI source
regenerated; handle remains usable during requests and pointer cancellation resets.
Empty/error/retry state is distinct from discovery intro, and address refresh
repeats the query. Old async fallback results are discarded after a new action.

Overpass was failing; vector fallback was unavailable on the phone raster renderer.
Photon bounded reverse/category search now supplies all8 categories on either
renderer, with Overpass backup (explicit configured Overpass retains priority).
Validate category/coordinates/type/radius; fixed queries, public caches only,
no owner credentials forwarded. Nearby cache1h; generic provider gate spaces
requests rather than429, cancellation cannot reserve future rate slots. Bound
response bytes, JSON shape and deadlines. Private places/schema retained.

Broader review also fixed Settings backup/member delegation, route selection/
swap/destination fields and repeat unsupported3D selection feedback. Candidate
15-group browser QA passes actual8category data, drive/walk/cycle/SVG/GPX,
exports/private notes, mocked sharing, syntheticGPS, dragging/pending/retry,
weather/AQI/layers, settings encrypted backup and second-account isolation.
Go tests/vet pass including failover/way-centre/invalid provider/rate cancellation.
No physical Android/GPS/mic/OS sharing/Google feature parity claim.

Artifacts I:/Projects/tooling/maps-redesign-20261026; only fixture
C:/Projects/.cache/maps-search-20261027, port19315, tracked search preview receipt.
Final7 source/build/exact-package repeat/public release/channel and both normal
Maps-only backed-up installs pending. Revalidate mutable owner versions. Never
use owner data as fixtures, replace public assets/tags, bypass updater or use
GitHub Actions. Preserve other apps/config/identities and newest recovery.
