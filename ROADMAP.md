## Sylos Roadmap / TODO

### Before Broader Alpha (IT firms / local business testing)

- [x] Switch back to GitHub
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
  - [x] Autoscaler: SFTP: Test autoscaler logic with SFTP transfers under load, confirm proper scaling.
  - [x] Autoscaler: OneDrive: Validate autoscaler under heavy operations, adjust bounds to prevent stalls or overload.
  - [ ] Autoscaler: SharePoint: Large-scale transfer validation, confirm limits are honored.
  - [x] Autoscaler: Dropbox: Ensure autoscaler is behaving optimally, no under/over-scaling, especially for batch uploads/downloads.
  - [x] Autoscaler: Box: Confirm scaling honors limits for upload/download, prevents timeout/stall scenarios.
  - [x] Autoscaler: Google Drive: Test at scale, dial in lower/upper bounds for stable and performant operation.
- [x] **Size Comparisons**: Path Review shows Size (Src) / Selected / Size (Dst) / Free (Dest); selected size updates via exclude deltas; optional per-provider free space (Sylos-FS `FSStorageInfo`) with soft over-capacity warnings in review and root selection.
- [ ] **Symlink / junction policy**: exclude by default in V1; ensure traversal inspects item *type* only (not target) to avoid loops
- [ ] **Installer / uninstaller**: not strictly required (zip/tar works), but wanted for polish
- [ ] **Build script**: makefile-style script to build + run tests per-repo, plus a broader integration test runner for pre-release regression checks
- [ ] **DB Transfer Security Model**: build a good security model for transferring DB files across Sylos instances so sensitive credential info isn't stored and transferred when exported / uploaded. 
- [ ] **Sign out without breaking active migrations**: Make it so all users can sign out of sylos while an active migration is running, as long as the sylos backend process is not stopped this should not stop any active migrations from running. Separately, they should be able to log in again during an active migration without interrupting it
- [ ] **New Wiki Article**: Set up wiki article for how to import / export migration DBs
- [ ] **Set up Github Discussions page**: Set up Github's discussion page so that we have sections for the typical things people would want / ask for to get it ready for forum usage for the community before testing starts.

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
- [x] **Filter rules engine**: composable rules (include/exclude by pattern, age, metadata, etc.) with combinations/negations
  - [x] Preview then exclude-all-matching with spare list; provenance via `filter_applications` + `exclusion_source`
  - [x] Local FS omits `.DS_Store` / `Thumbs.db` at list time (not via filters)
- [ ] **Offline docs embedding**: pull `sylos.wiki` repo into build pipeline, render MD → HTML/CSS, bundle into UI for offline/locked-down environments (piggybacks on build script work)
- [ ] **Implement scan only mode**: Build a feature where you can scan a singular source against no destination, purely for scanning and logging purposes, maybe for future automations.
- [ ] **Automation hooks**: Build automation hooks and allow for recurring migration / scanning scheduling
- [ ] **Terminal UI**: Finish work on Terminal UI to catch it up to the main Web UI and make it a viable (and user friendly) option for those wanting to opt out of the Web UI instead.