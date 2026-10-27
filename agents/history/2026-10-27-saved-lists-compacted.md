# Save-to-list and direct Home/Work UX

User authorizes implementation, push/release/install: all locations should expose
Save, followed by a list dropdown. Existing manual typing of Home is too complicated.
No Google/public-data regression; retain Maps9 local-depth provisioning.

Changes: AUI native-dialog Save picker, default Favourites / Home / Work / private
existing lists / optional new-list name and private note. Explicit confirm only;
cancel writes nothing. Selected-place quick/header/detail Save, independent nearby
and search/recent save controls, saved-item change-list action. Unset Home/Work
opens focused search with intent then preselected picker; existing shortcuts open
the saved location. Programmatic search Back no longer closes a newly opened save
picker. Keyboard-sized viewport, modal validation/progress, dark Cancel contrast
and focus return handled. Provider metadata/selected notes merge retained.

Backend serializes saved-place writes and uses a transaction. Home/Work canonical
names, one shortcut per account; replacing demotes old entries to Favourites while
retaining their records/notes. Foreign IDs refuse404 before shortcut updates. No
schema bump or owner data fixture. In-memory tests cover replacement/privacy/failure
rollback/custom list notes. UI fixture tested11 groups; final12/exact binary pending.

Next release10, source version bumped. Build locally with Go1.27.3/AUI, no Actions.
Keep immutable9 and withheld3/7, four managed packages/min92, official promotion.
Follow Maps-depth handoff for release tooling and normal protected installation.
Primary last9/manager94/MyAI15 with28services; fresh runtime matches Maps9/manager94.
Secondary still requires MyAccounts owner sign-in (fresh workspace-2/updates tab
585734409/Chrome3, marked handoff; user asked asynchronously). Never bypass sign-in
or claim both installed without evidence. Preserve other versions/pins/identities.

Active synthetic fixture C:/Projects/.cache/maps-save-20261027, private NTFS ACL,
bootstrap credentials only inside that fixture. Loopback synthetic map/Photon on
19316, app19315; both PIDs/paths in tooling/save-process.json. Safe restart-save.ps1
tracks app; do not lose provider PID when restarting (original provider stays).
Stop exact task-owned services and remove only this guarded fixture after QA.
Screenshots contain only synthetic names. Source/release/receipts stay under
I:/Projects/tooling/maps-redesign-20261026. Source9 receipt/tag/ZIP immutable.

## Preview10 released and primary installed

Immutable source5ed357b2d5c4a59397703d48fe2498f1f536c8f7, all17 public assets sizes/
SHA256 and anonymous executable download verified. Accepted/public/installed EXE
SHAa6caa386520de445cf52cabd62ca4dbbfb0991f8b2712d4ec5407ec98c6078ec. Official
channel/pin/default promotion3699323, anonymously checked; all other pins/prior
rows/withholdings retained. Exact packaged12 UI groups + Go tests/vet + AUI/TS/Vite
and public92/93/94 parser checks pass. Modal dark Cancel contrast fixed and checked;
recent Save control added. Selected metadata refresh handles a demoted shortcut.
No Actions, schema bump, public provider or owner functional fixture.

Primary normal Maps9->10 protected update completed with exact public94 controller,
29catalog apps/28healthy services, manager94 retained. Other app contracts/packages/
EXE bytes, owner configuration and decrypted identities compared in memory/matched.
Live Maps health10, HTML title and JS/CSS/executable hashes verified. Newest verified
recovery update-20261027T195656-d8d8f547eb4a retained by normal updater cleanup.
Do not repeat primary install. No country/map data automatically downloaded.

Secondary normal MyAccounts sign-in still required in Chrome3/tab585734409; fresh
updates login marked handoff, user asked asynchronously and no reply received.
No secondary apply submitted. Last93/Maps8/MyAI14 is historical, not a fresh read.
After user sign-in, inspect fresh updater versions/state and install ONLY Maps10
via normal backed-up flow; preserve all other pending updates. Never bypass auth
or reuse an OAuth callback/state. Current user phone screenshots predate10; do not
claim its secondary deployment until verified. Primary is complete.

Protected owned C fixture maps-save-20261027, exact app and Python provider PIDs,
ports19315/19316 stopped/removed with executable/command/path/reparse guards. Safe
restart script now preserves provider tracking. Source/release/receipt artifacts
remain in task tooling on I; docs/evidence/saved-lists-20261027.json is sanitized.

## Secondary delivery completed after owner sign-in

Owner signed in; normal workspace-2 updater freshly showed manager93/Maps8/MyAI14.
Checked public channel, selected ONLY Maps10 (cleared manager/MyAI/Mail/Chat/new
Push), submitted once. Normal download/package validation, scoped33.8MB recovery
copy and checksum verification completed; restart took longer than earlier runs,
with a temporary public502/disconnection. Automatic reconnect confirmed all9 stages
complete in4m2s, all apps healthy, newest verified recovery
update-20261027T211311-247ea6e5ea1b retained. No repeated apply or auth bypass.

Fresh release check shows Maps10 up to date; every other visible row matches the
pre-update baseline, retaining manager93/MyAI14/Chat10 and Push not installed.
Live /maps/healthz reports10/ok, portal /healthz shows all27services running.
Secondary proof is updater completion/version rows/live health, not an independent
remote executable/disk/credential audit. No owner records as development fixtures,
map dataset download, physical phone test or new release/asset replacement.
Both PCs now have Maps10; primary delivery remains as recorded above. Tooling
secondary-delivery10.json/completion10.txt/completion10.jpg retain local receipts;
sanitized current evidence is docs/evidence/saved-lists-20261027.json in the app.
Normal updater tab585734409/Chrome3 retained as deliverable. Sign-in/install handoff
is resolved; no repeat installation needed.
