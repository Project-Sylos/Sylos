## Sylos Roadmap / TODO

### Next Immediate Building

- [x] More end to end testing of various things
- [ ] Migrate off of DigitalOcean DNS to CloudFlare so we can redirect the chat.sylos.io site to the permanent discord invite link.
- [ ] Switch back to GitHub

### Before Broader Alpha (IT firms / local business testing)

- [x] **FS providers**: add SharePoint, OneDrive, and Box at minimum 
- [ ] **Autoscaler tuning (overall)**: validate large-scale behavior on all FS classes, ensure performance is optimized and not misbehaving. Add support for lower and upper bounds to prevent autoscaler from exceeding safe/constrained limits, especially for session-based operations (e.g., rate of uploaded/downloaded chunks).
  - [x] Autoscaler: Local - Physical Disk
    - Verify that scaling works correctly on physical disk access.
    - Confirm lower and upper bounds are enforced.
  - [x] Autoscaler: Local - Network Share
    - Validate scaling with network-mounted drives.
    - Ensure bounds logic prevents over/under-scaling.
  - [x] Autoscaler: Local - FUSE
    - Test scaling logic for FUSE (Linux/macOS).
  - [ ] Autoscaler: Local - WinFSP
    - Test scaling logic for WinFSP (Windows).
    - Confirm system respects set lower/upper bounds during extended use.
  - [ ] Autoscaler: SFTP: Test autoscaler logic with SFTP transfers under load, confirm proper scaling.
  - [ ] Autoscaler: OneDrive: Validate autoscaler under heavy operations, adjust bounds to prevent stalls or overload.
  - [ ] Autoscaler: SharePoint: Large-scale transfer validation, confirm limits are honored.
  - [x] Autoscaler: Dropbox: Ensure autoscaler is behaving optimally, no under/over-scaling, especially for batch uploads/downloads.
  - [ ] Autoscaler: Box: Confirm scaling honors limits for upload/download, prevents timeout/stall scenarios.
  - [x] Autoscaler: Google Drive: Test at scale, dial in lower/upper bounds for stable and performant operation.
- [x] **Size Comparisons**: Path Review shows Size (Src) / Selected / Size (Dst) / Free (Dest); selected size updates via exclude deltas; optional per-provider free space (Sylos-FS `FSStorageInfo`) with soft over-capacity warnings in review and root selection.
- [ ] **Symlink / junction policy**: exclude by default in V1; ensure traversal inspects item *type* only (not target) to avoid loops
- [ ] **Installer / uninstaller**: not strictly required (zip/tar works), but wanted for polish
- [ ] **Build script**: makefile-style script to build + run tests per-repo, plus a broader integration test runner for pre-release regression checks


### After Alpha, Before/Around Beta

- [ ] **Path Verifier Library (in progress via Go-Path-Linter, needs more test coverage)**
  - Initial implementation exists as Go-Path-Linter (GPL); core logic building out
  - Cross-platform + cross-provider path verification and cleaning
  - Verifies paths against destination service rules (Windows/Mac/Linux/Dropbox/SharePoint/OneDrive/Box/etc.)
  - Structured error output: raises or associates `.errors` to paths as needed
  - `.clean()` method returns auto-suggested clean paths
  - Integrated into path review phase: flags problem paths pre-copy, suggests fixes
  - Some issues are treated as warnings only and do not always require user action
  - [x] Everything above
  - [ ] More real-world scenario testing is still needed for full confidence (I'd say this is about 80-90% done)
- [ ] **Unified review page / cross-phase navigation**
  - Collapse traversal/copy/delete review into one state-aware page (or shared state layer across 3 pages — TBD)
  - New nav component/UX for flipping between phases (tabs? stepper? TBD — needs a design pass)
  - Navigation gated on **run state**, not phase state: free navigation when idle, locked to progress monitor when a phase is actively running (`{action} in progress...` messaging)
  - Retraversal after deletion: treat delete events in the status events table as authoritative — traversal-on-retry respects them rather than re-discovering/erroring on deleted items
- [ ] **Filter rules engine**: composable rules (include/exclude by pattern, age, metadata, etc.) with combinations/negations
  - Applied during traversal; excluded items don't get traversed further
  - Review UI lets users unmark auto-filtered items, triggering retraversal for just those items
- [ ] **Offline docs embedding**: pull `sylos.wiki` repo into build pipeline, render MD → HTML/CSS, bundle into UI for offline/locked-down environments (piggybacks on build script work)
- [ ] **Excel report outputs and final dashboard**
  - Build an in-app dashboard showing run metrics: API/OS call counts, item modifications, copies, total bytes moved, etc.
  - Make metrics exportable to Excel (async job or similar if needed)
- [ ] **Implement scan only mode** 
  - Build a feature where you can scan a singular source against no destination, purely for scanning and logging purposes, maybe for future automations.
- [ ] Build auomation hooks and allow for recurring migration / scanning scheduling
