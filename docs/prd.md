# NodeCleaner Product Requirements Document

**Version:** 1.0  
**Last Updated:** January 13, 2026  
**Status:** Draft  
**Owner:** Engineering & Product Team

---

## Executive Summary

**NodeCleaner** is a CLI tool designed to help developers reclaim disk space by identifying and safely removing stale dependency directories (initially `node_modules`, expandable to other ecosystems). Unlike general system cleaners, NodeCleaner understands developer workflows and provides transparency, control, and safety in cleaning operations.

---

## 1. Problem Statement

### Why This Matters

Developers working across multiple projects accumulate dependency directories that consume significant disk space—often hundreds of megabytes to several gigabytes per project. Over time, these directories accumulate from:

- Active projects that install dependencies regularly
- Abandoned or archived projects never cleaned up
- Learning repositories and tutorial code
- Hackathon experiments and proof-of-concepts

**Pain Points:**

1. **Silent Storage Drain:** Dependency folders grow invisibly, consuming disk space without user awareness
2. **Poor Tool Support:** Existing system cleaners (MacCleaner, CleanMyMac, DaisyDisk) focus on media files and binaries, not developer-specific directories
3. **Manual Hunting:** Developers waste time manually searching for large folders or using generic disk analyzers
4. **Risk of Breaking Projects:** Deleting the wrong folders can break active projects, making developers hesitant to clean
5. **Performance Impact:** Full disks slow down development machines and interrupt workflows

### What Success Looks Like

Developers can:

- Quickly scan their system for stale dependency directories
- See transparent information about what will be removed
- Safely reclaim gigabytes of disk space with confidence
- Maintain this cleanliness without manual effort

---

## 2. Goals & Non-Goals

### Goals

**Primary Goal:**  
Provide developers with a safe, transparent tool to identify and remove stale dependency directories, reclaiming significant disk space without disrupting active projects.

**Secondary Goals:**

- Build trust through transparency (dry-run, confirmations, detailed reporting)
- Minimize scan time through intelligent caching
- Create a foundation extensible to multiple language ecosystems
- Establish a CLI-first, automation-friendly workflow

### Non-Goals (MVP)

- **GUI application** — CLI-first approach for MVP
- **Automatic scheduling** — Manual invocation only
- **Multi-language support** — JavaScript/TypeScript ecosystem only for MVP
- **Cloud integrations** — Local-only operation
- **Analytics/telemetry** — Privacy-first approach with no data collection
- **Project health analysis** — Pure disk space focus, not code quality

---

## 3. Target Users

### Primary Personas

**1. The Multi-Project Developer**

- Maintains 10-50 repositories locally
- Works on multiple codebases simultaneously
- Frequently clones new projects for reference or contribution
- **Pain:** Disk constantly filling up, unclear where space goes

**2. The Learning Developer**

- Students, bootcamp graduates, self-taught programmers
- Clones tutorial repos, follows courses, experiments frequently
- **Pain:** Limited disk space on laptop, many abandoned learning projects

**3. The Hackathon Contributor**

- Rapid prototyping across multiple projects
- Creates many short-lived experimental repos
- **Pain:** No time to clean up after events, projects pile up

**4. The Context-Switcher**

- Engineers working on client projects, consulting, freelancing
- Switches between codebases frequently
- **Pain:** Each project installs full dependency trees, multiplying storage needs

### User Environment

- **Primary OS:** macOS and Linux
- **Storage:** 256GB-512GB SSD typical (expensive to upgrade)
- **Workflow:** Terminal-heavy, comfortable with CLI tools
- **Values:** Control, transparency, safety, speed

---

## 4. User Stories

### Core User Stories

**As a developer, I want to:**

1. **Discover stale dependencies** so I can understand what's consuming my disk space

   - Acceptance: See list of all `node_modules` folders with size and last access time

2. **Preview deletions before they happen** so I can avoid accidentally breaking projects

   - Acceptance: Dry-run mode shows exactly what will be deleted and space reclaimed

3. **Selectively choose what to remove** so I can keep active projects untouched

   - Acceptance: Interactive selection UI with checkboxes or numbered list

4. **Reclaim space safely** so I don't corrupt my development environment

   - Acceptance: Explicit confirmation required before any deletion occurs

5. **Avoid re-scanning unchanged directories** so subsequent runs are fast
   - Acceptance: Cache previous scan results, only rescan changed folders

### Future User Stories (Post-MVP)

6. **Clean multiple dependency types** so I can manage Python, Go, Java projects too
7. **Schedule automated cleanups** so my disk stays clean without manual intervention
8. **See recommendations** so the tool suggests what's safe to remove based on usage patterns

---

## 5. Features & Requirements

### 5.1 Core Features (MVP)

#### Feature 1: Filesystem Scanning

**What:** Recursively scan filesystem from a specified root path to detect all `node_modules` directories

**Why:** Developers need complete visibility into dependency locations and sizes to make informed cleanup decisions

**Requirements:**

- **REQ-1.1:** Accept user-defined root path as input (default: `$HOME`)
- **REQ-1.2:** Perform recursive directory traversal to find all `node_modules` folders
- **REQ-1.3:** Skip system and protected directories via ignore list (e.g., `/System`, `/Library`, `/Applications`)
- **REQ-1.4:** Collect metadata for each folder:
  - Absolute path
  - Size on disk (human-readable format)
  - Last modified timestamp
  - Last accessed timestamp
- **REQ-1.5:** Display progress indicator during scan (folder count, elapsed time)
- **REQ-1.6:** Handle permission errors gracefully (skip and log, don't crash)

**Acceptance Criteria:**

- Scan completes on typical developer machine (10-20 projects) in under 60 seconds
- All `node_modules` folders are detected regardless of nesting depth
- Metadata is accurate within ±1MB and ±1 hour for timestamps

---

#### Feature 2: Results Presentation

**What:** Display scan results in an organized, sortable CLI interface

**Why:** Developers need to quickly identify the largest and oldest folders to prioritize cleanup

**Requirements:**

- **REQ-2.1:** Present results in a table format with columns:
  - Size (sorted largest to smallest by default)
  - Path (relative to scan root if possible)
  - Last accessed date
  - Selection checkbox/indicator
- **REQ-2.2:** Support sorting by size, path, or last accessed time
- **REQ-2.3:** Show summary statistics:
  - Total folders found
  - Total space consumed
  - Potential space reclaimable
- **REQ-2.4:** Paginate results if more than 50 folders found
- **REQ-2.5:** Highlight folders not accessed in 30+ days (suggested for cleanup)

**Acceptance Criteria:**

- Results are readable in standard 80-column terminal
- User can navigate and sort results without confusion
- Summary accurately reflects totals

---

#### Feature 3: Interactive Selection

**What:** Allow users to select which folders to delete through interactive CLI prompts

**Why:** Gives developers full control over what gets deleted, preventing accidental removal of active projects

**Requirements:**

- **REQ-3.1:** Provide checkbox-style selection interface (spacebar to toggle, enter to confirm)
- **REQ-3.2:** Support "select all" and "deselect all" options
- **REQ-3.3:** Support "select by age" (e.g., all folders not accessed in 60+ days)
- **REQ-3.4:** Show running total of space to be reclaimed as selections change
- **REQ-3.5:** Allow filtering by path pattern (e.g., select all in `/archived-projects/`)

**Acceptance Criteria:**

- Selection interface is intuitive for CLI users
- Users can select/deselect individual or bulk items easily
- Running total updates accurately in real-time

---

#### Feature 4: Dry-Run Mode

**What:** Preview deletion operations without actually removing files

**Why:** Builds user trust by showing exactly what will happen before making irreversible changes

**Requirements:**

- **REQ-4.1:** Enabled by default or via `--dry-run` flag
- **REQ-4.2:** Display list of folders that would be deleted
- **REQ-4.3:** Show total space that would be reclaimed
- **REQ-4.4:** Clearly indicate "DRY RUN" mode in output
- **REQ-4.5:** Exit without making any filesystem changes
- **REQ-4.6:** Provide command to execute actual deletion

**Acceptance Criteria:**

- No files or directories are modified in dry-run mode
- Output clearly differentiates between dry-run and actual execution
- User understands exactly what will happen before proceeding

---

#### Feature 5: Safe Deletion

**What:** Remove selected folders with explicit confirmation and safety checks

**Why:** Prevents accidental data loss through explicit consent and validation

**Requirements:**

- **REQ-5.1:** Require explicit confirmation prompt before deletion (Yes/No)
- **REQ-5.2:** Show summary of what will be deleted in confirmation prompt
- **REQ-5.3:** Validate that selected folders still exist before deletion
- **REQ-5.4:** Delete folders recursively and report success/failure for each
- **REQ-5.5:** Display progress during deletion (e.g., "Deleting 3 of 10...")
- **REQ-5.6:** Generate deletion report:
  - Successful deletions
  - Failed deletions (with reasons)
  - Total space reclaimed
- **REQ-5.7:** Log operations to file (`~/.nodecleaner/deletions.log`)

**Acceptance Criteria:**

- User must explicitly type "yes" or confirm via prompt
- Deletion only proceeds after confirmation
- Report accurately reflects space reclaimed
- Failed deletions don't crash the program

---

#### Feature 6: Scan Caching

**What:** Store scan results locally to speed up subsequent scans

**Why:** Reduces scan time from minutes to seconds for unchanged directories, improving user experience

**Requirements:**

- **REQ-6.1:** Store scan results in `~/.nodecleaner/cache.json`
- **REQ-6.2:** Cache includes:
  - Folder path
  - Size
  - Last modified timestamp
  - Last scanned timestamp
- **REQ-6.3:** On subsequent scans, compare folder modified time against cached time
- **REQ-6.4:** Only re-scan folders where modified time changed
- **REQ-6.5:** Provide `--no-cache` flag to force full rescan
- **REQ-6.6:** Provide `--clear-cache` command to reset cache

**Acceptance Criteria:**

- Second scan completes in <10 seconds if no changes detected
- Cache invalidation correctly identifies changed folders
- Stale cache entries are removed for deleted folders

---

### 5.2 Non-Functional Requirements

#### Performance

- **NFR-1:** Initial scan completes in <2 minutes for filesystem with 100 projects
- **NFR-2:** Cached scan completes in <10 seconds
- **NFR-3:** Memory usage stays under 200MB during scan

#### Reliability

- **NFR-4:** Tool handles permission errors without crashing
- **NFR-5:** Tool handles symbolic links safely (follows or ignores based on config)
- **NFR-6:** Deletion failures for individual folders don't halt entire operation

#### Usability

- **NFR-7:** CLI output is readable on 80-column terminal
- **NFR-8:** Help documentation accessible via `--help` flag
- **NFR-9:** Error messages are clear and actionable

#### Security

- **NFR-10:** Tool never writes outside `$HOME/.nodecleaner/` for its own data
- **NFR-11:** Tool prompts for confirmation before any destructive operation
- **NFR-12:** Tool respects filesystem permissions (no privilege escalation)

#### Compatibility

- **NFR-13:** Works on macOS 10.15+ (Catalina and later)
- **NFR-14:** Works on Linux distributions with Node.js 18+ or Go 1.20+
- **NFR-15:** Installation via npm or standalone binary

---

## 6. Future Enhancements (Post-MVP)

### Phase 2: Multi-Ecosystem Support

**What:** Extend beyond JavaScript to other language dependency directories

**Target Directories:**

- Python: `.venv`, `venv`, `__pycache__`, `~/.cache/pip`
- Go: `vendor`, `~/go/pkg`
- Java: `.m2`, `.gradle/caches`
- Rust: `target/` (in Cargo projects)
- PHP: `vendor`

**Why:** Developers work across multiple languages; single tool is more valuable

---

### Phase 3: Intelligence Layer

**What:** Provide smart recommendations based on usage patterns

**Features:**

- Detect Git repositories with no commits in 90+ days
- Identify projects with no recent file access
- Score projects by "abandonment likelihood"
- Suggest safe-to-remove candidates automatically

**Why:** Reduces cognitive load on developers, automates decision-making

---

### Phase 4: GUI Application

**What:** Desktop application with visual interface (Electron or Tauri)

**Features:**

- Visual tree map of disk usage
- Drag-and-drop folder selection
- System tray integration
- macOS Finder extension

**Why:** Broader appeal to developers less comfortable with CLI

---

### Phase 5: Automation & Scheduling

**What:** Automated cleanup on schedule or triggers

**Features:**

- Cron job integration
- Git hooks (pre-commit, post-merge)
- Scheduled scans with email/notification reports
- "Auto-clean" mode with safety rules

**Why:** Zero-maintenance disk management

---

### Phase 6: Team & Enterprise Features

**What:** Multi-user and organizational capabilities

**Features:**

- Scan reports aggregated across team
- Shared ignore lists and policies
- CI/CD integration for cleanup in pipelines
- Storage reclamation dashboards

**Why:** Organizations can optimize disk usage across development fleets

---

## 7. Milestones & Release Plan

### Milestone 1: Foundation (Weeks 1-2)

**Goal:** Core scanning functionality

**Deliverables:**

- Filesystem traversal engine
- `node_modules` detection
- Metadata collection (size, timestamps)
- Basic CLI output (list of folders)

**Success Criteria:**

- Successfully scans 100+ projects in under 2 minutes
- Accurately reports folder sizes and timestamps

---

### Milestone 2: User Interaction (Weeks 3-4)

**Goal:** Selection and preview capabilities

**Deliverables:**

- Interactive selection UI (checkboxes)
- Sorting and filtering options
- Dry-run mode implementation
- Summary statistics display

**Success Criteria:**

- Users can select folders interactively
- Dry-run mode previews deletions accurately
- Zero false positives in reporting

---

### Milestone 3: Safe Deletion (Weeks 5-6)

**Goal:** Deletion functionality with safety guarantees

**Deliverables:**

- Confirmation prompts
- Deletion engine with error handling
- Operation logging
- Deletion reports

**Success Criteria:**

- Deletion requires explicit user confirmation
- Failed deletions don't crash tool
- Logs provide audit trail of operations

---

### Milestone 4: Performance Optimization (Weeks 7-8)

**Goal:** Caching and speed improvements

**Deliverables:**

- Scan caching implementation
- Cache invalidation logic
- Ignore list for system directories
- Performance benchmarking

**Success Criteria:**

- Cached scans complete in <10 seconds
- Cache correctly identifies changed folders
- No performance regressions from baseline

---

### Milestone 5: Polish & Release (Weeks 9-10)

**Goal:** Production-ready MVP

**Deliverables:**

- Comprehensive documentation
- Installation packages (npm/binary)
- Error message improvements
- Unit and integration tests (>80% coverage)
- README with usage examples

**Success Criteria:**

- Tool installable via package manager
- Documentation covers all features
- Beta testers successfully use tool
- Zero critical bugs in issue tracker

---

### Release Schedule

**Week 10:** v0.1.0 Beta Release (internal testing)  
**Week 12:** v1.0.0 Public Release  
**Week 16:** v1.1.0 (bug fixes, minor improvements)  
**Week 24:** v2.0.0 (multi-ecosystem support)

---

## 8. Success Metrics

### MVP Success Metrics

**Adoption Metrics:**

- 1,000+ npm/package downloads in first month
- 100+ GitHub stars in first month
- 10+ community-submitted issues or PRs

**Usage Metrics:**

- Average disk space reclaimed per user: >5GB
- Time saved vs manual hunting: >15 minutes per session
- Repeat usage rate: >30% of users run tool monthly

**Quality Metrics:**

- Zero critical bugs reported in first month
- User satisfaction: >4.0/5.0 stars
- False positive rate for deletion: <1%

### Long-Term Success Metrics (6-12 months)

**Growth:**

- 10,000+ active users
- 500+ GitHub stars
- Featured in developer newsletters/blogs

**Engagement:**

- 50%+ monthly active user rate
- 20+ community contributors
- Active Discord/Slack community

**Impact:**

- Average 10GB+ reclaimed per user lifetime
- Tool cited in "essential dev tools" lists
- Enterprise adoption inquiries

---

## 9. Risks & Mitigations

### Risk 1: Accidental Deletion of Active Projects

**Impact:** High | **Likelihood:** Medium

**Mitigation:**

- Mandatory dry-run by default
- Explicit confirmation required
- Highlight recently-accessed folders
- Provide undo/recovery documentation

---

### Risk 2: Slow Performance on Large Filesystems

**Impact:** Medium | **Likelihood:** High

**Mitigation:**

- Implement robust caching strategy
- Use system ignore lists
- Allow user-defined scan boundaries
- Optimize traversal algorithm (consider parallel scanning)

---

### Risk 3: Permission Errors Crashing Tool

**Impact:** High | **Likelihood:** Medium

**Mitigation:**

- Wrap all filesystem operations in try-catch
- Gracefully skip inaccessible directories
- Log permission errors for debugging
- Never require sudo/admin privileges

---

### Risk 4: Low Adoption Due to CLI-Only Interface

**Impact:** Medium | **Likelihood:** Low

**Mitigation:**

- Target early adopters comfortable with CLI
- Invest in clear documentation and examples
- Plan GUI for Phase 4 based on feedback
- Integrate with popular dev tools (VS Code)

---

### Risk 5: Competing Tools Emerge

**Impact:** Medium | **Likelihood:** Medium

**Mitigation:**

- Focus on developer-specific features competitors lack
- Build extensible architecture for quick feature additions
- Foster open-source community
- Maintain high code quality and reliability

---

## 10. Open Questions

**To be resolved before development:**

1. **Language Choice:** Node.js (faster to ship, familiar ecosystem) vs Go (better performance, standalone binaries)?

   - **Recommendation:** Start with Node.js for MVP speed, consider Go rewrite in v2.0

2. **Default Behavior:** Should dry-run be the default, or require a flag?

   - **Recommendation:** Dry-run by default for safety; explicit `--execute` flag for actual deletion

3. **Cache Location:** Should cache live in project directories or centralized location?

   - **Recommendation:** Centralized in `~/.nodecleaner/cache.json` to avoid polluting repos

4. **Symlink Handling:** Follow symlinks or ignore them?

   - **Recommendation:** Ignore symlinks by default; add `--follow-symlinks` flag for advanced users

5. **Deletion Method:** Move to trash/recycle bin vs permanent deletion?
   - **Recommendation:** Permanent deletion by default (dependency folders recoverable via package manager); consider `--trash` flag for v1.1

---

## 11. Appendix

### A. Competitive Analysis

| Tool       | Focus                  | Pros              | Cons                    | Differentiation                       |
| ---------- | ---------------------- | ----------------- | ----------------------- | ------------------------------------- |
| npkill     | `node_modules` cleanup | CLI, interactive  | JS-only, no caching     | We add caching + multi-ecosystem      |
| CleanMyMac | System cleaning        | User-friendly GUI | Expensive, no dev focus | We're free + dev-specific             |
| DaisyDisk  | Disk visualization     | Beautiful UI      | Generic, no automation  | We understand dev directories         |
| DevCleaner | Xcode cleanup          | Xcode-specific    | Apple ecosystem only    | We're cross-platform + multi-language |

### B. Technical Stack (Recommended)

**MVP (Node.js approach):**

- **Language:** Node.js 18+
- **CLI Framework:** Commander.js for argument parsing
- **Interactive UI:** Inquirer.js for prompts and selection
- **Filesystem:** Native `fs/promises` module
- **Testing:** Jest or Vitest
- **Distribution:** npm package + optional bundled binary (pkg or nexe)

**Alternative (Go approach):**

- **Language:** Go 1.20+
- **CLI Framework:** Cobra
- **Interactive UI:** bubbletea or survey
- **Distribution:** Single binary via goreleaser

### C. Installation & Usage (Planned)

```bash
# Install
npm install -g nodecleaner
# or
brew install nodecleaner

# Basic usage
nodecleaner scan ~/projects

# With options
nodecleaner scan --path ~/projects --no-cache --dry-run

# Interactive mode
nodecleaner clean

# Clear cache
nodecleaner cache --clear
```

---

## Document History

| Version | Date       | Author       | Changes              |
| ------- | ---------- | ------------ | -------------------- |
| 1.0     | 2026-01-13 | Product Team | Initial PRD creation |

---

## Approvals

| Role             | Name | Signature | Date |
| ---------------- | ---- | --------- | ---- |
| Product Manager  |      |           |      |
| Engineering Lead |      |           |      |
| UX Lead          |      |           |      |

---

**End of Document**
