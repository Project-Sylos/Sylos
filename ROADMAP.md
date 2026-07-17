## Sylos Roadmap / TODO

### Next Immediate Building

- [x] Minor UI polish / tweaks
- [ ] More end to end testing of various things

### Before Broader Alpha (IT firms / local business testing)

- [ ] **FS providers**: add SharePoint, OneDrive, and Box at minimum 
- [ ] **Autoscaler tuning**: validate large-scale behavior on Drive/Dropbox, make sure it's not leaving performance on the table or misbehaving. (Also need to add in lower and upper bound support to never let it go below certain lower / upper bounds. Especially in regards to upload / download sessions where time between uploaded / downloaded chunks might create issues.
- [ ] **Symlink / junction policy**: exclude by default in V1; ensure traversal inspects item *type* only (not target) to avoid loops
- [ ] **Installer / uninstaller**: not strictly required (zip/tar works), but wanted for polish
- [ ] **Build script**: makefile-style script to build + run tests per-repo, plus a broader integration test runner for pre-release regression checks
- [x] **UI polish pass**: themes readable, all buttons functional, no obvious bugs



### After Alpha, Before/Around Beta

- [ ] **Path Verifier Library** — cross-platform + cross-provider path verification and cleaning library
  - Verifies paths against destination service rules (Windows/Mac/Linux/Dropbox/SharePoint/OneDrive/Box/etc.)
  - Structured error output (raise or silent + `.errors`)
  - `.clean()` method for auto-suggested clean paths
  - Integrates into path review phase: flag problem paths pre-copy, suggest fixes
- [ ] **Unified review page / cross-phase navigation**
  - Collapse traversal/copy/delete review into one state-aware page (or shared state layer across 3 pages — TBD)
  - New nav component/UX for flipping between phases (tabs? stepper? TBD — needs a design pass)
  - Navigation gated on **run state**, not phase state: free navigation when idle, locked to progress monitor when a phase is actively running (`{action} in progress...` messaging)
  - Retraversal after deletion: treat delete events in the status events table as authoritative — traversal-on-retry respects them rather than re-discovering/erroring on deleted items
- [ ] **Filter rules engine**: composable rules (include/exclude by pattern, age, metadata, etc.) with combinations/negations
  - Applied during traversal; excluded items don't get traversed further
  - Review UI lets users unmark auto-filtered items, triggering retraversal for just those items



### After V1 (or whenever)

- [ ] **Offline docs embedding**: pull `sylos.wiki` repo into build pipeline, render MD → HTML/CSS, bundle into UI for offline/locked-down environments (piggybacks on build script work)