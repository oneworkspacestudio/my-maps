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

## Published12 and software delivered; secondary data setup pending

The earlier implementation/build-pending notes above are superseded by this section.
Immutable12 source2cc3cf1eb1d22fd579caf55a3d7b12ff842710c6, tagv0.1.0-preview.12.
EXE04d5af2b7c3cf7aeed00112e65448eb712629fa3f48ccfbbc0b979d067020d95;
ZIP15e67b6e62a849e44277bce8736b4aec996836e5107f5f891521ca8bd3412c2d.
All17 asset sizes/digests and anonymous EXE verified. Exact packaged native UI7
interaction groups pass, including a new visible-mobile-Exit bounding check. Actual
served JS/CSS/title/hash match. Public92–96 parser tests pass from byte-verified
local NTFS source subsets (I source reads stalled). No Actions or asset/tag changes.
Official channel promoted with c9f6b50/241398a after retaining concurrent97 channel
7700484 and Shield4 pins; anonymous Maps12 entry/all17assets verified. Final public
release: https://github.com/oneworkspacestudio/my-maps/releases/tag/v0.1.0-preview.12.

Primary normal Maps-only11→12 update complete; matching public96 controller used,
manager96/Shield4/all other30app contracts and process bytes/settings/secrets retained,
29services running, exact accepted app/JS/CSS bytes served. Latest verified recovery
update-20261028T215940-c617688861d6. Owner signed in privately in Chrome3 Maps tab449.
Primary previously had no map data. Normal owner UI now installed UK regional map+
postcodes (2853tiles,1852091 searchable places,575.8MiB) and UK local directions
(612.1MiB), both ready. Open map verified with local Canvas/attribution and visible
Directions control. No owner searches/records/GPS were development tests.

Secondary own updater tab444 selected ONLY Maps12, submitted once. All9stages
completed2m57s, all apps healthy; latest recovery update-20261028T215928-b3834bc283d6.
Manager95/unselected app versions retained by selection, but a follow-up release
check shows an empty choices table. Do not repeat Apply. New Maps tab453 redirects
to accounts-2 sign-in. Async question asks owner to sign in there for UK routing
setup; awaiting reply. Primary sign-in is separate and already completed. Keep tab453
handoff. After sign-in, verify Maps12 and existing coverage, open Map downloads,
install ONLY UK local directions (keep existing map pack), verify Ready, capture
proof and29running services via authenticated UI/normal health page. Gateway protects
/maps/healthz from anonymous access; do not bypass auth/extract cookies or claim a
remote disk/credential audit. Secondary97/Shield4 updates belong to another task;
leave them unchecked. Preserve current95/AI14/Chat10/Shield3 actual versions.

Browser AX/native automation started timing out; browser-only DOM APIs work:
getBrowser{id:'3'}, browser.user.claimTab/user.openTabs, tab.playwright.domSnapshot.
Primary handle tab449, secondary handoff453, own updater444; old tab445 closed.
Do not use categories441 or Shield tabs. Token-auth URLs are not history material.
Both protected fixture apps/Job engines19315/19316 stopped and ports verified closed.
Recursive deletion of C:/Projects/.cache/maps-routing-20261028 was rejected by
approval review without detailed reason. Do not retry via another deletion/move
mechanism. Synthetic data remains protected there; cleanup.json records it. Public
source parser subsets C:/Projects/.cache/maps-parsers-20261028 are also retained.
No preview is running. Tooling/public release artifacts remain on I.

Next: finish secondary routing download once owner signs in; record sanitized
both-PC receipts/health/directions-ready evidence and final histories, push doc-only
handoffs, retain exact immutable12 source/assets and latest verified recoveries.
