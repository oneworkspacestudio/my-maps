# Initial application — implementation and release preparation

User requested My Maps and My Meet as native Windows apps in One Workspace with shared design, SSO, standalone use, installation and public org releases. Free/open source foundations only; no GitHub Actions. UI uses AUI compiler with an adapter derived from Apache-2.0 My Tasks; shared apphost pinned to 0578d36ce297. Local Go/TypeScript/Vite builds passed. Isolated NTFS preview at C:/tmp/maps-meet-preview-20261020 only, no owner data used. Meet preview opens a real LiveKit room, publishes synthetic audio/video, admits a second account and displays meeting chat. Maps street tiles load; Photon search/results also work; route interaction has not received separate end-to-end acceptance. These checks are not physical-device, public-network or capacity acceptance.

Source is I:/Projects/my-maps; suite integration is I:/Projects/one-workspace-maps-meet branch codex/maps-meet. Build material lives in I:/Projects/tooling/maps-meet-20261020. No release or owner install yet. Preserve explicit README preview boundaries, provider policies, engine licence inventory and private-data boundaries. Continue by completing package/build receipts, org publication, combined Workspace release and normal backed-up owner update. Reconcile newer primary Workspace commits before final tag. Clean only named synthetic C fixture and owned child processes afterward.

Integrated navigation follow-up: explicit shared-picker trigger suppresses the floating launcher, and compact labelled icon controls keep the managed mobile header within the viewport. Meet also binds its signalling grants to the current authenticated app session. Maps2 and Meet3 supersede the initial previews in Workspace54; do not rewrite published assets/tags.

Final integrated minimum is Workspace54, incorporating the concurrent Clips4 restart repair. Initial Maps1 and Meet2 releases are public; Maps2/Meet3 are the final integrated candidates. Owner installation remains pending.

## Delivery checkpoint

Maps2 source7708ea6 and Meet3 sourceb39d099 are now public in their organization
repos, each with16 verified release assets including managed-update metadata.
The normal owner53-to54 update installed both;24 apps and23 services are healthy,
identity/storage fields were retained, and20 other app versions/packages were
preserved. Both public routes return the expected unauthenticated login redirect.
Managed browser SSO was not rechecked because Chrome control could not load its
request-header policy. Synthetic previews were removed from C after use.

The final combined suite is being superseded by55 to fix inherited
publication-time release ordering.54's source tag remains immutable and its
draft is held.55 chooses higher versions and prevents manager downgrades.
Drive3 and full suite55 binary publication are pending slow uploads; do not
claim those releases are public from this checkpoint. No paid services enabled.
