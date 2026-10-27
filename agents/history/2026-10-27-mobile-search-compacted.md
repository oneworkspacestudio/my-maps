# Mobile search and place views

User supplied new Google Maps phone references: more map space, dedicated search
with recents/Home/Work, immediate place Directions/Start/Save/Share. Earlier no-
Google-data/fully-self-hosted constraint still applies. Read self-hosting and
search-nearby handoffs first. Working base5e670fd on codex/google-maps-layout;
prior cloud-control cleanup was uncommitted and is included in this revision.
No new release, official channel/pin update or owner installation in this turn.
No Google data/API/photo source added. Geographic coverage request remains open;
public map/geocoder/router defaults remain, so full self-hosting is NOT complete.

AUI Maps screen now has full-screen search, autofocus editor, Home/Work/exact-
coordinate shortcuts, saved matches and abortable/debounced provider suggestions,
private recents/Clear all and native browser Back. Submitted addresses choose the
first result in a compact place sheet; pending manually expanded sheets stay open.
Phone home peek112px+64px tabs, category136px without tabs, selected210px without
tabs; horizontal immediate Directions/Start/Share/Save row, icons and address.
Drag reveals notes/details. Compact result rows replace oversized decorative art.
Visual viewport sizing handles reduced keyboard space; manifest/browser chrome
uses neutral theme. Browser address bar itself is outside app layout control.

Start explicitly requests GPS, computes a route, starts foreground following;
denial preserves the place and explains manual Directions. Cleared searches stop
following/reset route points, route pins only appear in route mode. Saved places
merge private fields with current provider metadata rather than discarding phone,
hours/website. Source uses existing renderer, routes, private exports/backups.
No microphone/cloud speech, external weather/AQI, traffic placeholder, Mapillary
or FIRMS controls. No Google photos/reviews, native voice navigation or full Google
feature-parity claim. Existing satellite/terrain/public default providers remain
unresolved in the separate fully-local data provisioning work.

Recent API GET/POST/DELETE /api/recents is native account/Origin/CSRF protected;
20 bounded per-account query/place entries, no private notes/passive GPS. Writes
queue before Clear all; stale loads discarded. Schema2 adds recent_places while
preserving places/cache. migrate() applies only missing versions, supports restart
and refuses future versions; downgrade requires compatible protected backup.
Memory-SQLite upgrade/restart/future-data test and recent kind validation pass.
Explicit HTTPS endpoints or literal HTTP loopback accepted for native self-hosted
services; insecure hostname/public HTTP/credentials rejected. No new services/data
provisioned on owner machines. MyAI's Google-enabled SearXNG was not integrated.

Artifacts I:/Projects/tooling/maps-redesign-20261026, Maps-mobile.exe candidate.
Only protected C:/Projects/.cache/maps-mobile-20261027 used; local synthetic tile/
Photon/route provider at19316, app19315, tracked mobile-process.json. No owner
records or sessions used as fixtures. UI tests deny external browser requests;
local providers are synthetic contracts, NOT real geographic-data acceptance.
AUI/TS/Vite, Go tests/vet/native build passed. Functional checks caught saved-place
selection missing history; fixed. Screenshot review caught dark text specificity;
fixed with contrast assertion. Synthetic provider initially parsed semicolons as
URL params (one-point route); test provider fixed/cache cleared. Assertions now
wait for sheet height transitions rather than mistaking mid-animation bounds.

Final native candidate df5a5b6ad4bab6ca0363a664db5ee609b621eb5af25644e50e43b5c3afb7ae21
repeated18 groups with zero browser errors/external browser requests. Saved-place
metadata retention, dark contrast and manual pending expansion are included.
Sanitized evidence and synthetic mobile search/place screenshots are under docs/.
Guarded named test root and both tracked processes/ports removed; owner data was
not touched. Source-only branch publication is recorded by Git tracking state;
no binary release, official channel update or owner installation in this turn.
Next: provision real local map/index/router datasets once coverage known; no public
fallback, real denied-egress checks, new immutable Maps9+ with four-package/min92
contracts, official channel and normal backed-up Maps-only installs on both.
Revalidate other tasks' manager/app versions first. Do not use Actions, owner
fixtures, replace immutable tags/assets or ship a blank map under a self-host claim.
