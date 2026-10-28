# Local directions and reference mobile route screens

User wants visible Directions/Start, vehicle/travel modes and Google-like screens,
with the persistent constraint: self-hosted, no Google data, remove unsupported
features. User chose "do what is needed out of the box experience for users" when
asked about coverage. UK routing downloads on both existing PCs are authorized;
new UK/local map setup enables routing by default, with size/depth review.

Active source C:/Projects/.cache/maps-source-20261028, base87cbb01. All work so far
uncommitted; source now version12. I:/Projects/my-maps is stale; do not use its Git.
Intentional web/node_modules junction points to I:/Projects/my-maps/web/node_modules;
do not follow/remove recursively. Caches/builds/artifacts on I. No subagents.
Skills repository-history, write-aui-ui, browser-screenshot applied. AUI screens
and generated React updated, no handwritten JSX.

Implemented: integrated bounded BRouter1.7.10/Temurin21.0.12+8 download plans,
fixed archive hashes, stable graph ETags/lengths/CRC, safe ZIP extraction/staging,
atomic active pointer, cancellation retaining prior data, successful replacement
cleans only prior validated pack. Owners install; members view/use routes. Shared
map/directions installer exclusion. Defaults UK regional + optional postcodes +
routing; combined size/free-space checks, sequential map then routing install.
Closing the setup page before stage2 requires reopening its directions installer.

Own Java bridge exposes only token-authenticated loopback POST route/cancel and
health. No upstream RouteServer, profile uploads, suspects/recent-waypoint APIs,
precise URLs/logs. Token through stdin, no CLI/environment secret, suppressed engine
stdout/stderr, hidden supervised child with Windows kill-on-close Job. Source/class
pair hashes verified at package build. Runtime/legal and MIT/OSM notices preserved.
No third-party runtime bundled in Maps.exe; explicit data install downloads it.

Drive/walk/cycle profiles, car/EV/hybrid energy model, touring/road/gravel/MTB and
walking/hiking, actual avoidance options, alternatives, up to8 stops, turn cues,
GPX and bounded route links. Place green actions always visible; search stays
compact. Foreground GPS next-turn/remaining-distance/exit and explicit off-route
update, no background/voice guarantee. No transit/live traffic/charging/height/
weight/Google photos or full parity claims. Moped excluded: upstream experimental.
Error codes survive frontend for setup/coverage recovery; stale route cancellation,
compatible mode defaults and privacy/account guards retained. First-pack env-router
seed persistence corrected. Native backups do not include separate downloaded
coverage; docs corrected to explain reinstallation after restore.

Go tests/vet, AUI/TypeScript/Vite/road-style checks pass. Real engine car/bike/foot
and longer public-city alternative routes passed. Normal owner UI installed graph.
Seven native mobile groups pass; actual combined first-run map+directions install
passed with no overlapping installers; cancellation/owner-member/CSRF/UK plan passed.
Private process loopback/auth/no-CLI-token checks passed; forced parent exit killed
its engine. GPS is simulated only at public fixture coordinates, not a real phone.
UK plan snapshot543546009 download/2167508939 required bytes. No owner test data.
Evidence in docs/evidence/directions-*.json and docs/assets/directions-*.png.

Protected synthetic C:/Projects/.cache/maps-routing-20261028 retains only disposable
fixtures; data2 is active native fixture on19315, first-run on19316; another public
graph is in root/directions from initial CLI acceptance. Process receipts under
I:/Projects/tooling/maps-routing-20261028. Use guarded restart.ps1; omit bootstrap
flags after owner creation. Initial bootstrap-password launch was rejected by
automatic approval; normal first-run setup succeeded. Do not retry rejected broad
fixture moves/deletes. Stop exact owned processes and their Job children at finish.

Release12 not yet committed/published/installed. Next: commit current intended files,
local scripts/build.ps1 -Release into fresh I build-release12, exact packaged native
QA and released manager92/93/94/95/96 parser checks; publish immutable17assets from
committed source with Actions disabled, official publish-update-index promotion,
anonymous bytes/channel verification; normal Maps-only updates on both PCs then UK
routing install in normal owner UI on both. No repeat of Maps11 delivery.

Current channel C:/Projects/.cache/maps-channel-20261028 fast-forwarded to3fa0535.
Primary now Workspace96, SHA bc11fb6bb68b93844ce0b47a627816f42f30f472a91719ef8769829e7661a10a,
Maps11 and29services. Secondary95/Maps11 and29services. Preserve actual installed
versions, data/identities/other app bytes. Public96 source downloaded and digest
verified70a825b53d1c946eb4968b059fea193bb1c09a1a52dea10d04e9a061fa3b2825,
under task/sdk96-source/one-workspace. Old95 helper must not be reused on primary96.
Release templates/canonical maps92 four-package layout under tooling/maps-redesign-20261026.
Own secondary updater Chrome3 tab585734444 authenticated; categories tab441 belongs
to other task, do not claim. Browser use is available again. No external writes for12.
