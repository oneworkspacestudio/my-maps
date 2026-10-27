# Preview evidence

The redesign uses the supplied phone screenshots for layout. Its map, floating
search and category chips, Explore/You/Contribute tabs, place sheet, profile menu
and map-type/detail dialog were reviewed at phone and desktop sizes.

Go tests/vet, AUI generation, TypeScript, Vite and native Windows builds pass.
Isolated protected NTFS data was used for address search, exact coordinates,
private saved notes, saved route endpoints, real driving routes and GPX download.
Photon and FOSSGIS returned real results. Terrain, NASA imagery, weather and
Open-Meteo/CAMS air-quality requests completed without provider errors.

The default Overpass endpoint returned network/server failures during testing.
Nearby discovery therefore also uses actual POIs from loaded vector tiles,
clearly labelled when fuller provider details are unavailable. No sample ratings,
photos, businesses or traffic are presented as live data.

Final release and installed-runtime receipts are recorded separately. These are
synthetic browser/Windows checks, not physical GPS/microphone, clean-host, reboot,
background-navigation or exhaustive security acceptance. Full Google parity is
not achieved; see FEATURES.md and README for each capability and remaining gap.

Preview6 reproduces the released5 blank-map failure with WebGL disabled and
verifies visible real raster tiles in phone/touch emulation and desktop GPU-off
mode. Checks cover actual Photon search and OSRM/SVG routes/GPX, synthetic GPS,
private pin/note save, 360/390/576-width layouts, failed 3D selection recovery,
preserved camera after simulated context loss and bounded stalled-vector recovery.
The normal desktop vector path also renders real features. No uncaught browser
errors were observed. These are browser fault-injection checks, not a physical
Android test; the exact cause on the reporting phone is not independently known.

Preview7 checks the map-first search/category flow, actual pointer drags, manual
expansion while search is pending, retry after a simulated nearby outage, query
refresh, exact-coordinate validation, private pins/notes and Home shortcut,
mocked sharing without notes, synthetic GPS, GeoJSON export/delete, actual
drive/walk/cycle routes and SVG/GPX, route controls, measurement, imagery/layers,
weather/AQI, contribution links, appearance, Settings and encrypted backup.
Standalone-user creation and a second account's private-place isolation are
checked using disposable protected data. All eight nearby categories return
real Photon results in the phone renderer without WebGL. Tests also cover
provider failover, malformed upstream data, rate spacing/cancellation and bounds.
Exact-package and deployed receipts are separate. Physical device permissions,
GPS/microphone, OS share sheets, reboot and missing Google services are unclaimed.

Final8 adds a visibility assertion: category pins must lie inside the unobscured
phone map, including sparse categories such as petrol stations. Bounds now fit
the returned places with space for search controls and the collapsed sheet.
The earlier7 check counted marker DOM nodes and missed off-screen pins. Public7
was never promoted/installed and remains immutable, explicitly withheld in favour
of8. Do not treat DOM marker count alone as map visibility acceptance.
