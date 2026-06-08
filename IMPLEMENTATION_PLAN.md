# Dead Drop Dispatch Implementation Plan

## Purpose

This plan tracks the first playable prototype and MVP for **Dead Drop Dispatch**, a Go terminal dispatch game built with Bubble Tea v2, Bubbles v2, and Lip Gloss v2. The near-term end state is a dashboard-first playable loop where the player accepts jobs, assigns runners, chooses routes, advances turns, receives deterministic consequences, and can complete or lose a seven-night run without requiring an LLM.

Priority scale:
- `P2` = required for the next milestone
- `P1` = important hardening / enabler
- `P0` = optimization or enhancement

Status markers: `[ ]` = open, `[~]` = in progress, `[x]` = done.

## Context & Decisions

- Build the deterministic offline game first; the LLM is optional fiction dressing only.
- Keep deterministic mechanics in `internal/game`, static content in `internal/content`, Bubble Tea orchestration in `internal/app`, and rendering in `internal/tui`.
- Use generated route options per job for the prototype instead of full city pathfinding.
- Use a discrete turn loop. One turn is a dispatch cycle of roughly thirty to sixty in-world minutes.
- Each runner can carry up to two bundled jobs when compatibility rules allow it.
- Active work resolves at turn end for the first playable prototype unless a delay outcome keeps work active.
- Complications are fixed-choice events decided by deterministic mechanics.
- Hide exact risk math from the player, but show contributing factors, route traits, cargo warnings, runner stress, district pressure, and faction context.
- Offline fallback content is mandatory. LLM failures must never block play.
- Dashboard-first playability is the near-term product priority; secondary tabs are scaffolding until the core loop is playable.
- Settings screens are deferred until config or LLM work needs them; first playable should run with sensible offline defaults and no required setup.
- Complication choices resolve deterministically against game state with persistent cargo integrity, delay, economy, stress, heat, dispatch integrity, and faction deltas.
- Complication reporting uses deterministic after-action messages and effect-bearing event logs so UI and future save data can show what happened without re-running mechanics.
- Unreliable or inaccurate information should be an explicit intel mechanic with source/staleness/trust tags, not accidental LLM drift.
- `feature_expansion_ideas.md` is a candidate backlog source for post-MVP systems; keep those items after the first playable and MVP loop unless one becomes necessary for core playability.
- Open design choice: the exact hidden-risk presentation should settle on contributing factors plus terse warnings, without an aggregate numeric chance.
- Open design choice: first playable content cadence should likely generate three available jobs per turn, with roughly five meaningful accepts per night.
- Prototype city content now starts with six districts after adding Ashgate Yard; keep dashboard rendering compact before expanding further.
- Turn advancement now uses explicit game phases for messages, jobs, dispatch assignment, resolution, complications, reports, and city updates; `space` advances that phase coordinator from the dashboard.
- The run clock uses six turns per night and rolls over to the next night during city update; crossing past night seven leaves final win/loss interpretation to the existing run-status checks and later RUN-06/RUN-07 work.
- Job templates are copied into `GameState` at run creation so the game-layer turn coordinator can refresh available jobs without importing static content.
- Runner recovery happens during city update: ready runners shed stress, injured runners tick recovery, and recovered injured runners return to ready.
- After-action reporting should explain what physically happened during deliveries, especially partial and failed outcomes; injuries need concrete severity/cause text instead of only a binary injured state.

## ID Prefix Registry

| Prefix | Domain |
|---|---|
| DOC | Documentation, agent guidance, and release notes |
| FND | Project foundation and package structure |
| APP | Bubble Tea application model and input orchestration |
| TUI | Terminal rendering, dashboard panels, screens, and visual style |
| STA | Core game state, resources, and run-ending rules |
| CNT | Static districts, runners, factions, cargo, and route content |
| JOB | Job templates, generation, and route option generation |
| ASN | Job acceptance, runner assignment, and bundling |
| MSG | Player responses, message actions, and client/runner interaction rules |
| INT | Intel quality, source reliability, and imperfect information |
| DST | Dynamic district state, turf pressure, and control drift |
| BBS | Diegetic BBS, shadow grid, rumors, tips, and side boards |
| DOW | Runner downtime, vices, safehouse recovery, and side effects |
| AST | Dispatch assistant personality, advice, favors, and drift |
| RES | Deterministic resolver and outcome application |
| CMP | Complication events, choices, and consequences |
| RUN | Turn, night, and seven-night run loop |
| ECO | Credits, costs, bribes, treatment, and upgrades |
| SAV | Save, load, autosave, and seeded run persistence |
| FIC | Offline fiction fallback text and report content |
| LLM | Optional provider-agnostic LLM fiction layer |
| BAL | Balance, polish, accessibility, and release preparation |

## Phase 1 — Foundation And Shell

Establish the Go module, package boundaries, Bubble Tea shell, and initial dashboard surface that future gameplay can inhabit.

**Definition of done:** `go test ./...` passes, the app starts with `go run ./cmd/ddd`, and the dashboard renders the core panels with navigable focus.

| Status | ID | Priority | Task |
|---|---|---|---|
| [x] | FND-01 | P2 | Initialize the Go module and add Bubble Tea v2, Bubbles v2, and Lip Gloss v2 dependencies. |
| [x] | FND-02 | P2 | Create `cmd/ddd/main.go` and wire it to the Bubble Tea app. |
| [x] | FND-03 | P2 | Create initial package layout for `internal/app`, `internal/game`, `internal/tui`, and `internal/content`. |
| [ ] | FND-04 | P1 | Create planned `internal/save` and `internal/llm` package directories when their implementation begins. |
| [x] | FND-05 | P2 | Add `.gitignore` and basic `go test ./...` validation. |
| [ ] | DOC-01 | P1 | Add `README.md` with build, run, test, and prototype scope notes. |
| [x] | APP-01 | P2 | Implement Bubble Tea model, `Init`, `Update`, `View`, quit keys, resize handling, and dashboard state. |
| [x] | APP-02 | P2 | Add focus model and navigation for tabs, panel cycling, directional selection, confirm, and help. |
| [ ] | APP-03 | P2 | Add `esc` back/cancel behavior and `space` advance/resolve behavior for the playable loop. |
| [x] | TUI-01 | P2 | Render static dashboard panels: status bar, city, jobs, runners, messages, detail pane. |
| [x] | TUI-02 | P2 | Apply initial visual style: thin borders, restrained palette, compact tables, focused panel styling. |
| [x] | APP-04 | P2 | Add smoke tests for app initialization. |

## Phase 2 — Deterministic Content And Jobs

Build the static world model and generated jobs needed for repeatable prototype runs.

**Definition of done:** The initial state contains valid districts, runners, factions, cargo, routes, generated jobs, and fixed-seed tests cover content references and job validity.

| Status | ID | Priority | Task |
|---|---|---|---|
| [x] | STA-01 | P2 | Define `GameState`, `District`, `Runner`, `Faction`, `CargoType`, `Job`, `Route`, `ActiveJob`, `Bundle`, `Message`, `LogEntry`, and `GamePhase`. |
| [x] | STA-02 | P2 | **DONE (2026-06-06).** Added run status and end-reason types plus pure run evaluation for victory, bankrupt, burned, collapse, roster loss, and faction lockout endings. |
| [x] | STA-05 | P2 | **DONE (2026-06-06).** Set prototype thresholds for credit target, heat maximum, dispatch integrity floor, faction-hostile count, and all-runners-unavailable failure. |
| [x] | STA-03 | P2 | Add seeded RNG and deterministic initial state creation. |
| [x] | STA-04 | P2 | Add tests for initial state validity. |
| [x] | CNT-01 | P2 | Add five districts with surveillance, traffic, faction control, danger, and signal quality. |
| [x] | CNT-02 | P2 | Add three runners with speed, stealth, nerve, talk, loyalty, stress, traits, and availability states. |
| [x] | CNT-03 | P2 | Add four factions with reputation and suspicion. |
| [x] | CNT-04 | P2 | Add five cargo types and five route types. |
| [x] | CNT-05 | P2 | Add tests for content counts and references. |
| [x] | CNT-06 | P1 | **DONE (2026-06-07).** Added Ashgate Yard as a sixth static district using existing faction coverage, updated tests and MVP scope docs, and kept the current dashboard pane layout unchanged. |
| [x] | JOB-01 | P2 | Add ten job templates covering clinic rush, data hand-off, witness transfer, prototype lift, union favor, dirty evidence, decoy package, split route, curfew run, and bad client. |
| [x] | JOB-02 | P2 | Generate jobs from templates with cargo, origin, destination, deadline, payout, faction source, modifiers, fallback title, and client message. |
| [x] | JOB-03 | P2 | Generate contributing risk factors without exposing exact odds. |
| [x] | JOB-04 | P2 | Generate two to four route options per job. |
| [x] | JOB-05 | P2 | Validate generated jobs: different origin and destination, deadlines fit the seven-night run, and payouts scale with cargo, distance, urgency, faction pressure, and hidden risk. |
| [ ] | JOB-06 | P1 | Expand static job templates beyond the initial ten, covering each cargo type and faction with more varied objectives, risk factors, route preferences, deadlines, payout profiles, and short client messages; update content count/reference tests. |

## Phase 3 — Assignment And Resolution

Turn jobs into active work, support compatible bundles, and resolve route outcomes through deterministic game mechanics.

**Definition of done:** A player can accept jobs, assign runners, choose routes, bundle compatible jobs, resolve active work, and see mechanical logs and player-facing explanations backed by tests.

| Status | ID | Priority | Task |
|---|---|---|---|
| [x] | ASN-01 | P2 | Add job acceptance, runner assignment, and route selection flows. |
| [x] | ASN-02 | P2 | Let a ready runner take one job. |
| [x] | ASN-03 | P2 | Let a runner take a second bundled job when compatible. |
| [x] | ASN-04 | P2 | Add bundle compatibility checks: witnesses block bundling; bundled jobs require matching route type or shared route district; destination complexity can apply. |
| [x] | ASN-05 | P2 | Add bundle penalties for delay pressure, stress gain, cargo conflict, detection exposure, and destination complexity. |
| [x] | ASN-06 | P2 | Prevent assignment to injured, missing, burned, or otherwise unavailable runners. |
| [x] | ASN-07 | P2 | Show bundled jobs in runner roster and detail pane. |
| [x] | ASN-08 | P2 | Add tests for valid and invalid assignments. |
| [ ] | ASN-09 | P2 | Add cancel, abandon, and abort action hooks for accepted or active jobs where the current phase allows them. |
| [ ] | ASN-10 | P2 | Add reroute and delay action hooks for active jobs and complications; consumed by CMP-03 and RUN-01. |
| [x] | RES-01 | P2 | Implement deterministic route resolver independent of TUI and LLM. |
| [x] | RES-02 | P2 | Compute hidden detection, delay, injury, cargo damage, and eligible betrayal or interception risks. |
| [x] | RES-03 | P2 | Apply runner stat, district stat, route type, cargo type, faction, and bundle modifiers. |
| [x] | RES-04 | P2 | **DONE (2026-06-07).** Added first-class resolver representation for success, delay, detection, complication, injury, cargo damage, failure, and rare interception outcomes. |
| [x] | RES-05 | P2 | Produce mechanical `JobResult` values. |
| [x] | RES-06 | P2 | Produce player-facing explanations using contributing factors rather than exact math. |
| [x] | RES-07 | P2 | **DONE (2026-06-07).** Made resolver outcome application explicit for credits, heat, runner stress/state/loyalty, faction reputation/suspicion, cargo integrity, and dispatch integrity. |
| [x] | RES-08 | P2 | Add deterministic logs and fixed-seed resolver tests. |
| [x] | RES-09 | P1 | **DONE (2026-06-07).** Added structured injury detail on job results with severity, cause, recovery estimate, and summary text tied to route type, cargo, detection, complications, and clinic-favor recovery; delivery reports now include the injury detail. |

## Phase 4 — Complications, Economy, And Turn Loop

Make the prototype feel like a game loop instead of isolated mechanics by adding complications, player responses, turn phases, resource changes, and run endings.

**Definition of done:** The player can advance turns through job resolution, handle fixed-choice complications, make fixed message or route actions, see city and roster state change, and win or lose a seven-night run.

| Status | ID | Priority | Task |
|---|---|---|---|
| [x] | CMP-01 | P2 | **DONE (2026-06-07).** Added pending complication model, status, state storage, and resolver-result queueing with tests. |
| [x] | CMP-02 | P2 | **DONE (2026-06-07).** Added first five complication definitions: checkpoint, scanner sweep, runner panic, cargo leak, and signal loss. |
| [x] | CMP-03 | P2 | **DONE (2026-06-07).** Added fixed structured choices for each first-playable complication and snapshots on queued complications. |
| [x] | CMP-04 | P2 | **DONE (2026-06-07).** Added deterministic complication choice resolution with validation, resolved status, state effects, faction context, economy spending, logs, and tests. |
| [x] | CMP-05 | P2 | **DONE (2026-06-07).** Added structured complication costs for bribes, treatment, stress, heat, cargo integrity, delay turns, dispatch integrity, faction reputation, and faction suspicion. |
| [x] | CMP-06 | P1 | **DONE (2026-06-07).** Added complication opened/resolved after-action messages, effect-bearing resolution log entries, and reporting assertions around complication resolution. |
| [x] | CMP-07 | P1 | **DONE (2026-06-07).** Added the remaining MVP complication definitions and fixed choice sets: gang toll, client changes terms, drone tail, witness refuses, data trace, curfew drop, and rival courier. |
| [x] | CMP-08 | P1 | **DONE (2026-06-07).** Added invariant tests enforcing exactly twelve unique MVP complication definitions, lookup coverage for every MVP type, rejection of unknown types, non-empty choice metadata, unique choice IDs, fixed choice sets, and queued choice snapshots. |
| [x] | MSG-01 | P1 | **DONE (2026-06-07).** Defined fixed message response actions for clients, runners, and factions: refuse, ask more pay, ask more info, threaten, reassure, stall, deceive, cancel, and accept, with typed audiences, lookup helpers, allowed-action filtering, and invariant tests. |
| [x] | MSG-02 | P1 | **DONE (2026-06-07).** Added deterministic message response resolution by message ID with audience/action validation, resolved status storage, effect summaries, logs, response reports, and fixed state changes for credits, heat, dispatch integrity, runner stress/loyalty, and faction reputation/suspicion. |
| [x] | MSG-03 | P1 | **DONE (2026-06-07).** Added response validation and outcome tests covering missing snapshot choices, no-mutation rejection paths, every fixed response action, allowed audience choices, deterministic state deltas, stored effects, logs, and response reports. |
| [x] | INT-01 | P1 | **DONE (2026-06-07).** Added deterministic intel reports with source, timestamp, staleness, confidence, supported claims, omitted tags, false/incomplete/stale/biased risk tags, generated job intel snapshots, aging helpers, and claim validation so LLM-06 can reject unsupported claims instead of inventing misinformation. |
| [x] | ECO-01 | P2 | **DONE (2026-06-07).** Added economy transaction helpers for payouts, nonpayment, bribes, treatment, and operating costs, with deterministic logs and tests. |
| [x] | ECO-02 | P1 | **DONE (2026-06-07).** Added deterministic upgrade shop definitions, purchased-upgrade state, available-upgrade filtering, upgrade spending transactions, install logs, duplicate/unknown/insufficient-credit guards, and first upgrades: signal relay, safehouse, fake credential printer, clinic favor, dead-drop locker, and scrambler. |
| [ ] | ECO-02b | P1 | Add sub-upgrade trees for purchased upgrades, including safehouse equipment branches, nested costs, prerequisite checks, purchased sub-upgrade state, and tests for ownership and duplicate-purchase rules. |
| [x] | ECO-03 | P1 | **DONE (2026-06-07).** Applied upgrade effects to route intel claim depth, safehouse stress recovery, fake-credential checkpoint costs, clinic injury recovery, contraband risk, and data trace risk, with focused tests. |
| [ ] | ECO-04 | P1 | Add treatment and bribe spending paths used by complications and end-of-night recovery. |
| [ ] | ECO-05 | P1 | Add tests for economy transactions, upgrade effects, treatment, bribes, and operating costs. |
| [x] | RUN-01 | P2 | **DONE (2026-06-07).** Added explicit turn phases for messages, job review, dispatch assignment, active-job resolution, complication gating, reports, and city updates; dashboard `space` now advances the phase coordinator and focused tests cover phase transitions. |
| [x] | RUN-02 | P2 | **DONE (2026-06-07).** Added fixed turn-per-night clock advancement, night rollover during city update, night-start briefs/logs, final seven-night boundary progression, and focused tests for same-night, rollover, and completed-run cases. |
| [x] | RUN-03 | P2 | **DONE (2026-06-07).** Stored job templates on game state, added available-job refresh from the current night/turn clock, regenerated postings when messages advance into the job-board phase, preserved accepted/active work, and added focused tests. |
| [x] | RUN-04 | P2 | **DONE (2026-06-07).** Integrated runner recovery into city update so ready runners shed stress, safehouse bonuses still apply, injured runners tick recovery, recovered runners return to ready, and turn advancement reports/logs recovery effects with focused tests. |
| [x] | RUN-05 | P2 | **DONE (2026-06-07).** Applied nightly operating costs and heat decay (city update changes) on night transition, with focused unit tests. |
| [x] | RUN-06 | P2 | **DONE (2026-06-07).** Added terminal victory transition at the seven-night boundary with a resolved run-complete message and final event log entry when the credit target is met. |
| [x] | RUN-07 | P2 | **DONE (2026-06-07).** Added terminal failure transitions for bankrupt, burned, collapse, roster loss, faction lockout, and final credit shortfall; game-over reporting is duplicate-safe. |
| [x] | RUN-08 | P2 | **DONE (2026-06-07).** Added focused tests for final-night victory, final credit shortfall, immediate failure, duplicate game-over reporting, and existing turn/night transitions. |

## Phase 5 — Playable TUI Command Center

Promote the dashboard from a display surface into the primary command center for the playable loop, then fill out secondary screens where they improve repeated play.

**Definition of done:** The dashboard supports the first playable loop end to end, secondary screens have useful empty states and key maps, and visible TUI layout changes have smoke coverage.

| Status | ID | Priority | Task |
|---|---|---|---|
| [x] | TUI-03 | P2 | **DONE (2026-06-07).** Dashboard now accepts jobs, assigns runners, cycles selected routes, supports bundling through assignment, advances/resolves turns with `space`, renders pending assignment route detail, and has app-level workflow tests for accept/assign/resolve. |
| [ ] | TUI-04 | P2 | Add jobs, routing, runners, factions, messages, log, and help screens. |
| [ ] | TUI-05 | P2 | Add screen-specific key maps. |
| [x] | TUI-06 | P1 | **DONE (2026-06-07).** Added phase-aware empty states for no posted jobs and no messages plus clearer runner detail copy for no active assignments, with render smoke coverage. |
| [ ] | TUI-07 | P2 | Add confirmation prompts for risky actions. |
| [x] | TUI-08 | P2 | **DONE (2026-06-07).** Added a right-aligned header `NEXT` chip and a compact center spacer strip showing current action, selected job risk factors, selected route summary, and notices without adding a new panel. |
| [x] | TUI-09 | P2 | **DONE (2026-06-07).** Added compact `B1`/`B2` bundle markers in the runner roster, `Bundle n/2` detail headers with penalty display, and pending-assignment cues when the selected runner would bundle or is full, with render smoke coverage. |
| [ ] | TUI-10 | P1 | Improve panel sizing across narrow and wide terminals. |
| [ ] | TUI-11 | P1 | Improve help text once the loop input model is final. |
| [ ] | TUI-12 | P2 | Render an ASCII district map or route topology that makes district-to-district choices legible. |
| [ ] | TUI-13 | P2 | Add progress and status indicators for night, turn, deadlines, active jobs, runner state, cargo integrity, and heat. |
| [ ] | TUI-14 | P1 | Format event logs as timestamped diegetic terminal output with clear outcome and cause text. |
| [x] | TUI-15 | P1 | **DONE (2026-06-08).** Exposed fixed message response actions from the dashboard message feed with message selection, response cycling, enter-to-send resolution, action-strip reply cues, message detail rendering, and app/render tests. |
| [x] | TUI-16 | P1 | **DONE (2026-06-07).** Added in-panel city sector briefings: city focus can select districts, `enter` opens a same-panel briefing with district description, control, stats, pressure, signal, and job-touch count, and `esc` returns to the sector list. |

## Phase 6 — Persistence And Offline Fiction

Add save/load, autosave, seeded runs, and enough deterministic fiction to make the game playable without network access or an LLM.

**Definition of done:** Current runs persist safely, offline text covers every mechanical outcome used by the MVP, and LLM-disabled play remains complete.

| Status | ID | Priority | Task |
|---|---|---|---|
| [ ] | SAV-01 | P2 | Define save file schema with save version field. |
| [ ] | SAV-02 | P2 | Save and load the current run as JSON. |
| [ ] | SAV-03 | P2 | Autosave after each turn. |
| [ ] | SAV-04 | P1 | Handle missing or corrupt save files gracefully. |
| [ ] | SAV-05 | P1 | Add seeded run support. |
| [ ] | SAV-06 | P2 | Add tests for save/load round trip. |
| [ ] | FIC-01 | P2 | Add fallback job titles, client messages, runner messages, faction messages, and after-action reports. |
| [ ] | FIC-02 | P2 | Select fallback messages by cargo, faction, result, and tone. |
| [ ] | FIC-03 | P2 | Keep fallback text short enough for TUI panels. |
| [ ] | FIC-04 | P2 | Add tests for fallback availability. |
| [ ] | FIC-05 | P1 | Ensure MVP content includes at least thirty fallback messages and twenty fallback after-action reports. |
| [x] | FIC-06 | P2 | **DONE (2026-06-07).** Added deterministic delivery outcome text for success, partial, failed, and intercepted jobs, with cause clauses for delay, cargo damage, injury, detection, complications, payout cuts, and intercept fallout; after-action messages now use the richer result summaries. |

## Phase 7 — Optional LLM Layer

Layer in provider-agnostic fiction generation while keeping deterministic mechanics authoritative and offline fallback mandatory.

**Definition of done:** When configured, the LLM can dress jobs, messages, and after-action reports through validated structured JSON; when unavailable or invalid, offline fallback takes over without blocking play.

| Status | ID | Priority | Task |
|---|---|---|---|
| [ ] | LLM-01 | P1 | Define provider-agnostic LLM client interface. |
| [ ] | LLM-02 | P1 | Add config loading from file and environment variables plus an offline mode switch. |
| [ ] | LLM-03 | P1 | Add job dressing, NPC and runner message, and after-action report prompts. |
| [ ] | LLM-04 | P1 | Require structured JSON responses. |
| [ ] | LLM-05 | P1 | Validate JSON fields, max text lengths, and allowed claims. |
| [ ] | LLM-06 | P1 | Reject unauthorized mechanics, rewards, choices, districts, and outcomes. |
| [ ] | LLM-07 | P1 | Retry failed calls once and fall back to offline text after validation or network failure. |
| [ ] | LLM-08 | P1 | Add tests with a fake LLM client. |

## Phase 8 — Balance, Polish, And Release Readiness

Tune the MVP into a small, sharp, replayable terminal game and document how to run it.

**Definition of done:** The MVP can be played in a single terminal session, no single runner or route dominates, release instructions are present, and `go test ./...` passes.

| Status | ID | Priority | Task |
|---|---|---|---|
| [ ] | BAL-01 | P1 | Tune runner stats, district stats, cargo risks, payout ranges, heat gain and decay, stress gain and recovery, bundle penalties, and operating costs. |
| [ ] | BAL-02 | P1 | Ensure no single runner dominates every route. |
| [ ] | BAL-03 | P1 | Ensure no single route dominates every job. |
| [ ] | BAL-04 | P1 | Add release build instructions and sample config. |
| [ ] | BAL-05 | P0 | Add final visual and copy polish after MVP mechanics settle. |

## Phase 9 — Post-MVP Feature Expansions

Capture candidate systems from `feature_expansion_ideas.md` without letting them displace first-playable or MVP fundamentals.

**Definition of done:** Each expansion has a deterministic game model, offline content path, focused tests, and TUI exposure that fits the terminal dashboard without blocking the core dispatch loop.

| Status | ID | Priority | Task |
|---|---|---|---|
| [ ] | DST-01 | P0 | Add mutable district state overlays for surveillance, traffic, danger, signal quality, control pressure, active timers, and trend direction without losing static district baselines. |
| [ ] | DST-02 | P0 | Apply delivery influence rules: contraband raises surveillance and danger temporarily, medical cooler delivery can reduce danger and raise traffic near clinic districts, and rival faction jobs shift district control pressure. |
| [ ] | DST-03 | P0 | Resolve turn-end turf events and control changes with deterministic city-wide reports and mechanical route or complication modifiers. |
| [ ] | DST-04 | P0 | Show district state drift in TUI district lists with compact trend indicators and changed-stat emphasis. |
| [ ] | BBS-01 | P0 | Add diegetic BBS/shadow grid model with boards such as `alt.cyber.dispatch`, `rec.runner.gear`, and `net.factions.leaks`, backed by offline templated posts. |
| [ ] | BBS-02 | P0 | Add information brokerage actions that spend credits for BBS tips revealing hidden job risks, runner loyalty pressure, upcoming faction raids, or INT-01 intel report upgrades. |
| [ ] | BBS-03 | P0 | Add freelance board side jobs generated from anonymous posts with higher risk, offline fallback text, and ordinary assignment/resolution compatibility. |
| [ ] | BBS-04 | P0 | Add courier chat posts that reflect recent routes, runner stress, loyalty pressure, and district danger as dynamic diegetic feedback. |
| [ ] | DOW-01 | P0 | Add runner vice traits and safehouse downtime needs, including addiction, gambling, medication, and faction-linked treatment hooks. |
| [ ] | DOW-02 | P0 | Add between-night downtime funding actions that spend credits to recover runner stress immediately while recording deterministic side-effect risk. |
| [ ] | DOW-03 | P0 | Resolve vice side effects through messages and state changes, including higher runner cuts, gang debt/bribe pressure, medical faction reputation, and corporate suspicion. |
| [ ] | DOW-04 | P0 | Add tests and balance checks for downtime recovery, vice side effects, safehouse costs, and runner availability interactions. |
| [ ] | AST-01 | P0 | Add dispatch assistant drift states for compliance, anarchist, and corrupted modes with deterministic triggers from player choices. |
| [ ] | AST-02 | P0 | Add assistant advice and favor hooks, including occasional locked job slots and required data shard delivery jobs for code-fragment uploads. |
| [ ] | AST-03 | P0 | Model corrupted assistant fabricated warnings through INT-01 false-risk intel tags so inaccurate information is explicit and cross-checkable, not accidental text drift. |
| [ ] | AST-04 | P0 | Expose assistant personality state, advice, warnings, and favors in the TUI without crowding the dashboard-first loop. |

## Milestone 1 — First Playable Prototype

Exit criteria:
- Phase 1 complete.
- Phase 2 complete.
- Phase 3 complete.
- Threshold task STA-05 complete.
- Core loop tasks RUN-01, RUN-02, RUN-03, RUN-04, RUN-05, RUN-06, RUN-07, and RUN-08 complete.
- Dashboard play tasks TUI-03, TUI-07, TUI-08, TUI-09, TUI-12, and TUI-13 complete.
- Action hook tasks ASN-09, ASN-10, and ECO-01 complete.
- Player can start a new run; see city, jobs, runners, messages, and status; accept jobs; assign runners; bundle up to two compatible jobs; choose routes; advance turns; resolve deterministic outcomes; read outcome logs; win or lose a seven-night run; and play without an LLM.
- `go test ./...` passes.

## Milestone 2 — MVP

Exit criteria:
- Milestone 1 complete.
- Persistence tasks SAV-01, SAV-02, SAV-03, SAV-04, SAV-05, and SAV-06 complete.
- Offline fiction tasks FIC-01, FIC-02, FIC-03, FIC-04, and FIC-05 complete.
- Economy and interaction tasks ECO-02, ECO-03, ECO-04, ECO-05, MSG-01, MSG-02, MSG-03, TUI-14, and TUI-15 complete.
- Complication expansion tasks CMP-07 and CMP-08 complete.
- Optional LLM tasks LLM-01, LLM-02, LLM-03, LLM-04, LLM-05, LLM-06, LLM-07, and LLM-08 complete.
- At least ten job templates, twelve complication types, thirty fallback messages, and twenty fallback after-action reports are present.
- README explains install, run, config, save files, and offline mode.
- LLM failures never block play.
- `go test ./...` passes.
