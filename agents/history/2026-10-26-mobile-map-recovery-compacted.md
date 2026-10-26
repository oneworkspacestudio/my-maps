# Mobile blank-map correction

Continues the authorized Maps push/release/install on both PCs after the user
reported a blank Android map on workspace-2. Load the reference-layout handoff
for existing installation/publication constraints. Released5's vector-only path
reproduces a blank screen and uncaught GPUInitializationError with WebGL disabled;
the reporting phone's exact graphics/network state is not independently known.

Phones/touch devices now start with Leaflet1.9.4 DOM/image tiles, preserving the
configured XYZ endpoint and attribution through the existing authenticated,
allowlisted resource proxy. A shared map controller retains search, private pins,
route SVG/GPX, measurements and foreground GPS. Desktop/explicit 3D retains
MapLibre6.15. Constructor failure, context loss or a 12s blank/stalled map switches
to standard tiles and a clear notice. Failed 3D resets layer selection; camera
survives renderer changes. Existing custom providers and private schemas retained.

Synthetic protected NTFS fixture: C:/Projects/.cache/maps-mobile-render-20261026,
port19315; tracked receipt and build/proof artifacts under
I:/Projects/tooling/maps-redesign-20261026. Never use owner data as fixtures.
Candidate checks pass real raster tiles with GPU disabled, Photon search,
private notes/pins, OSRM SVG/GPX, synthetic GPS, touch default with GPU enabled,
360/390/576 layout, failed 3D recovery, normal desktop vector features,
context-loss/camera recovery and deliberately aborted vector requests. No uncaught
page errors; physical Android/GPS/mic/reboot acceptance is not claimed.

Go1.27.3 tests/vet pass. Final immutable6 package, exact-binary repeat, canonical
registered four-package/min92 parser check, public release/channel promotion and
normal Maps-only backed-up installation on both PCs are pending. Revalidate
dynamic manager/app state; preserve every other app and existing recovery. Never
replace historical assets/tags, weaken updater checks or run GitHub Actions.
