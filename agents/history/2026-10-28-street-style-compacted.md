# Reference-based local map cartography

User requests actual roads/colours/labels matching supplied Google Maps screenshot,
continuing earlier release/install workflow. Google reference is appearance only:
retain self-hosted data, no Google/public runtime fallback or automatic owner pack
download. Both PCs were Maps10 before this change; primary manager94/28services,
secondary manager93/27services are last known, require fresh state before delivery.

Replaces flavor light with street-map-style.ts explicit Canvas rules: pale land,
cyan water, mint parks/woodland, pale buildings, blue-grey outlined main roads,
white local streets, bridge/tunnel layering, soft rail/paths/boundaries, named
park anchors from POI points, category-coloured labels with spacing/zoom gates.
RoadLabel follows gentle bends, keeps reversed lines upright and rejects short/
degenerate/hairpin labels; local Arial only. No new dependency/schema/provider.
Product AUI screens untouched. Version11 source prepared; no release yet.

Go tests/vet, AUI/TS/Vite and road-label geometry checks pass. Protected C fixture
maps-style-20261028 installed194 tiles/1250features for a3km Bognor Regis selection.
Candidate native11 passes7 real rendering/interaction groups: actual pixel palette,
road/park/POI text, mobile15/16/18, overview6, desktop resize, local search/map-first
sheet, nearby pins and Save/cancel, with WebGL disabled/no external browser requests
or exceptions. View screenshots in docs/assets/map-street-style.png and tooling.
Test harness initially targeted the old List label; corrected to Choose a list.
Restart helper now checks candidate existence before stopping and omits bootstrap
flags after the first owner exists; only disposable fixture was affected.

Artifact root I:/Projects/tooling/maps-style-20261028, tracked process.json and guarded
restart.ps1; current candidate Maps-style4.exe, loopback19315; private fixture path
C:/Projects/.cache/maps-style-20261028. Stop its exact tracked process and remove
only that guarded protected fixture when QA completes. Never reuse old fixture PIDs.
Release tooling/verified92/93/94 SDKs/canonical four-package layout/channel checkout
remain under tooling/maps-redesign-20261026. Build immutable17asset11 from committed
source, exact-package QA and public parser checks before publishing. No Actions,
minimum92, retain all other pins/history/withheld3/7. Promote official channel and
anonymously verify before normal app-only backed-up installation; preserve owner
identities/other apps/installed map data. Prior accepted10 assets/tag are immutable.
