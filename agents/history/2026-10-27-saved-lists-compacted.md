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
