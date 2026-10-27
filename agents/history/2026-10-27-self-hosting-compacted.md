# Self-hosting correction - awaiting coverage

User overrides the Google-clone goal: no Google data; fully self-hosted de-Googled
solution, reuse existing Workspace search where suitable, remove impossible features.
Prior Maps8 push/release/both protected installations are complete; this is a new
correction, not evidence that Maps8 was fully local. Base clean5e670fd, branch
codex/google-maps-layout. See search-nearby handoff for immutable releases, both
installed Maps8/manager93 and protected official update constraints.

Audit: Maps8 uses public non-Google OSM/Photon/FOSSGIS/OpenFreeMap/Mapterhorn/NASA,
plus optional Open-Meteo/Mapillary/FIRMS. Browser SpeechRecognition may be cloud-
backed. NASA's GoogleMapsCompatible matrix name is not Google map data. App/private
places self-hosted does not mean map data/services self-hosted. Default public
providers still remain in current working source; DO NOT claim fully local.

Read current my-search AGENTS/latest history/README/API: it federates private
Workspace apps, not geographic/web search. My AI has bundled SearXNG; read-only
inspection of internal/searchruntime/settings.yml found Google explicitly enabled
alongside other public engines. Maps does not use it. No MyAI/Search files changed;
do not wire Maps through Google-enabled web search or pretend web search supplies
geocoding, POIs, tiles or routing. Local geographic indexes/data are needed.

Working source (uncommitted, unpublished): removed browser speech button/handler,
weather/AQI frontend and backend, traffic availability-only placeholder, external
Mapillary/FIRMS controls. Renamed transit/cycle entries to Transport stops/Cycle
facilities to describe actual nearby POIs; replaced scan icon with LocateFixed
for coordinates. AUI regenerated; docs/NOTICE updated; historical acceptance kept
with explicit current correction appendix. Map-first sheet changes retained.

AUI generation/TypeScript/Vite pass. Go tests and vet pass with synthetic transports.
No new fixture/owner-data tests, commit/tag/release/channel change or installation.
Owner installations still8, including removed controls until a future update.

Asked asynchronously for required geographic coverage: UK / Europe / worldwide.
No answer yet. Need choose/provision local datasets before delivering self-hosting;
do not silently ship a blank map or label public-provider defaults self-hosted.
Photon supports local embedded OpenSearch + Java21 and downloadable dumps; global
index official docs ~95GB and64GB recommended RAM. Current App.New only permits
HTTPS, so local HTTP loopback support must be designed/tested for native providers.
PMTiles offers downloadable local tile archives; server rendering still must work
on GPU-disabled phones (keep current actual raster recovery requirement). OSRM
Windows upstream build is experimental. Evaluate actual Windows/native stack and
storage/update lifecycle before installing it; no provider/data downloaded yet.

Resume: get coverage response, implement local providers/capability gating with
no public fallback, audit endpoints/redirects/CSP/browser cloud calls, test real
local data under denied external egress and exact package. Remove unprovisioned
optional features, retain private places/exports/manual sheet. Then new immutable
Maps9+release if ready, canonical four-package/min92 manifests, official channel,
normal backed-up Maps-only update on both with current manager source/guards.
Revalidate mutable versions and all other app pins first. Never Actions or owner
functional fixtures. Current partial cleanup is not finished delivery.

Follow-up UI revision is in 2026-10-27-mobile-search-compacted.md. It includes
this cleanup, private schema2 recent history, literal HTTP-loopback provider
support and mobile reference flows. Full local geographic data is still pending.
