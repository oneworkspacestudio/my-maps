# Managed release publishing

The managed release for Workspace92 contains Maps.exe and the three registered
Maps licence/notice files. The standalone archive also contains compatibility.json,
guides, source receipts and screenshots.

Start from the hash-verified published Workspace92 catalog. For the Maps group,
apply catalog.InstallPackageIDs for canonical selected apps: maps-contract at
docs/MyMaps/compatibility.json is optional and is not registered by that install
path. Do not include it in this managed manifest or move another installed path.
Use schema2/integration3 with minimumManager0.3.0-preview.92. Keep launch arguments,
callback, health route, storage and identities unchanged. Never use an owner's
runtime catalog as publisher metadata.

Before publication, pass the released92 manifest parser/selection against this
canonical registered layout, including a synthetic custom port and retained
unselected app packages. After promotion, require the deployed92 normal read-only
preflight to select the exact expected executable hash before applying.

Preview3's product UI passed but its five-file managed manifest was rejected by
92 before installation. It remains immutable and is superseded by4. Do not replace
published assets or weaken the updater's package-boundary checks. Build4 uses
Go1.27.3 and the accepted3 JavaScript/CSS/worker bytes.
