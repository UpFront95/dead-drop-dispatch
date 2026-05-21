# Dead Drop Dispatch Implementation Plan

This plan converts `dead-drop-dispatch-spec.md` into a buildable task list for the first playable prototype and MVP.

## Locked Design Decisions

- [x] Build the deterministic offline game first.
- [x] Use the LLM only as a fiction layer: job dressing, messages, and after-action reports.
- [x] Start with generated route options per job rather than full city pathfinding.
- [x] Use a discrete turn loop.
- [x] Allow each runner to carry up to two bundled jobs at once.
- [x] Resolve active work at turn end.
- [x] Present complications as immediate fixed-choice events during resolution.
- [x] Hide exact risk math from the player.
- [x] Show contributing risk factors, route traits, cargo warnings, runner stress, district pressure, and faction context.
- [x] Make offline fallback content mandatory.
- [x] Treat the first playable prototype as the initial target before full MVP polish.

## Remaining Design Questions

- [ ] Define bundling constraints.
  - Can any two cargo types be bundled?
  - Should witnesses block bundling?
  - Should prototype, contraband, or medical cargo add extra bundle risk?
  - Can bundled jobs have different destinations, or must one route plausibly serve both?

- [ ] Define route compatibility for bundled jobs.
  - One shared route for both jobs?
  - A route sequence with two stops?
  - A simple compatibility score between job districts?

- [ ] Define how much hidden risk feedback the UI provides.
  - Show only contributing factors?
  - Show vague warnings like `unsteady`, `watched`, or `bad route`?
  - Show no aggregate `Low / Medium / High` risk label?

- [ ] Define active job timing.
  - Do all accepted jobs resolve at the end of the same turn?
  - Can longer routes occupy a runner across multiple turns?
  - For the first prototype, prefer same-turn resolution unless a delay occurs.

- [ ] Define first playable content size.
  - Spec says five jobs per night.
  - Turn structure suggests fresh jobs each turn.
  - Recommendation: generate three available jobs per turn, with roughly five meaningful accepts per night.

## Milestone 0: Project Foundation

- [x] Initialize Go module.
- [x] Add Bubble Tea v2 dependency.
- [x] Add Lip Gloss v2 dependency.
- [x] Add Bubbles v2 dependency.
- [x] Create `cmd/ddd/main.go`.
- [ ] Create initial package layout:
  - [x] `internal/app`
  - [x] `internal/game`
  - [x] `internal/tui`
  - [x] `internal/content`
  - `internal/save`
  - `internal/llm`
- [ ] Add `README.md` with build and run commands.
- [x] Add `.gitignore` for Go build artifacts and local config.
- [x] Add basic `go test ./...` validation.

## Milestone 1: Terminal Skeleton

- [x] Implement Bubble Tea app model.
- [x] Implement `Init`, `Update`, and `View`.
- [x] Handle quit keys: `q`, `ctrl+c`.
- [x] Handle window resize messages.
- [x] Add dashboard screen state.
- [x] Add focus model for panels.
- [ ] Add keyboard navigation:
  - [x] `tab` / `shift+tab` cycles panels
  - arrows and `hjkl` move selection
  - `enter` confirms
  - `esc` backs out
  - `space` advances or resolves
  - [x] `?` opens help
- [x] Render static dashboard panels:
  - [x] status bar
  - [x] city panel
  - [x] job board
  - [x] runner roster
  - [x] message feed
  - [x] detail pane
- [x] Add initial visual style:
  - [x] thin borders
  - [x] restrained terminal palette
  - [x] compact tables
  - [x] clear focused panel styling
- [x] Add smoke test for app initialization.

## Milestone 2: Core Game State

- [x] Define `GameState`.
- [x] Define `District`.
- [x] Define `Runner`.
- [x] Define `Faction`.
- [x] Define `CargoType`.
- [x] Define `Job`.
- [x] Define `Route`.
- [x] Define `ActiveJob`.
- [x] Define `Bundle`.
- [x] Define `Message`.
- [x] Define `LogEntry`.
- [x] Define `GamePhase`.
- [ ] Define victory and failure states.
- [x] Add seeded RNG to game state.
- [x] Add deterministic initial state creation.
- [x] Add tests for initial state validity.

## Milestone 3: Static Content

- [x] Add five districts:
  - Northline
  - Floodglass
  - Saint Orison Market
  - Port Kestrel
  - Crown Verge
- [x] Add district stats:
  - surveillance
  - traffic
  - faction control
  - danger
  - signal quality
- [x] Add three runners:
  - Mira Vale
  - Kaito Senn
  - Vex Calder
- [x] Add runner stats:
  - speed
  - stealth
  - nerve
  - talk
  - loyalty
  - stress
- [x] Add runner states:
  - ready
  - on job
  - injured
  - burned
  - missing
- [x] Add four factions:
  - Helix Municipal Security
  - Kestrel Dock Union
  - Saint Orison Clinic Network
  - Asterion Systems
- [x] Add faction reputation and suspicion.
- [x] Add five cargo types:
  - data shard
  - medical cooler
  - witness
  - contraband package
  - corporate prototype
- [x] Add five route types:
  - main artery
  - service tunnels
  - market weave
  - drone corridor
  - floodline
- [x] Add tests for content counts and references.

## Milestone 4: Job Generation

- [ ] Add ten job templates:
  - clinic rush
  - quiet data hand-off
  - witness transfer
  - prototype lift
  - union favor
  - dirty evidence
  - decoy package
  - split route
  - curfew run
  - bad client
- [ ] Generate jobs from templates.
- [ ] Assign cargo type, origin, destination, deadline, payout, faction source, and modifiers.
- [ ] Generate fallback title and client message.
- [ ] Generate contributing risk factors without exposing exact odds.
- [ ] Generate two to four route options per job.
- [ ] Ensure origin and destination are different.
- [ ] Ensure deadlines fit the seven-night run.
- [ ] Ensure payouts scale with cargo, distance, urgency, faction pressure, and hidden risk.
- [ ] Add tests for generated job validity.

## Milestone 5: Assignment And Bundling

- [ ] Add job acceptance flow.
- [ ] Add runner assignment flow.
- [ ] Add route selection flow.
- [ ] Let a ready runner take one job.
- [ ] Let a runner take a second bundled job when compatible.
- [ ] Add bundle compatibility checks.
- [ ] Add bundle penalties:
  - delay pressure
  - stress gain
  - cargo conflict
  - detection exposure
  - destination complexity
- [ ] Prevent assignment to injured, missing, or unavailable runners.
- [ ] Show bundled jobs in runner roster.
- [ ] Show bundled job details in detail pane.
- [ ] Add tests for valid and invalid assignments.

## Milestone 6: Resolver

- [ ] Implement deterministic route resolver independent of TUI and LLM.
- [ ] Compute hidden detection risk.
- [ ] Compute hidden delay risk.
- [ ] Compute hidden injury risk.
- [ ] Compute hidden cargo damage risk.
- [ ] Compute hidden betrayal or interception risk for eligible jobs.
- [ ] Apply runner stat modifiers.
- [ ] Apply district stat modifiers.
- [ ] Apply route type modifiers.
- [ ] Apply cargo type modifiers.
- [ ] Apply faction reputation and suspicion modifiers.
- [ ] Apply bundle modifiers.
- [ ] Resolve success, delay, detection, complication, injury, cargo damage, failure, and rare interception.
- [ ] Produce mechanical `JobResult`.
- [ ] Produce player-facing explanation using contributing factors rather than exact math.
- [ ] Update credits, heat, stress, loyalty, reputation, suspicion, cargo integrity, and dispatch integrity.
- [ ] Add deterministic logs for every meaningful result.
- [ ] Add fixed-seed resolver tests.

## Milestone 7: Complications

- [ ] Define complication model.
- [ ] Add first five complication types:
  - checkpoint
  - scanner sweep
  - runner panic
  - cargo leak
  - signal loss
- [ ] Add choices for each complication.
- [ ] Resolve each choice mechanically.
- [ ] Apply costs such as bribes, stress, heat, cargo damage, delay, or reputation shifts.
- [ ] Add complication messages.
- [ ] Add after-action log entries.
- [ ] Add tests for complication resolution.

## Milestone 8: Turn And Run Loop

- [ ] Implement turn phases:
  - messages arrive
  - jobs shown
  - assignments made
  - jobs resolve
  - complications resolve
  - rewards and damage apply
  - after-action reports appear
  - city state updates
- [ ] Add night counter.
- [ ] Add fixed turns per night.
- [ ] Generate new jobs at the appropriate phase.
- [ ] Recover runner stress between turns or nights.
- [ ] Tick injury recovery.
- [ ] Apply end-of-night operating costs.
- [ ] Apply city state changes.
- [ ] Add victory condition.
- [ ] Add failure conditions:
  - bankrupt
  - burned
  - collapse
  - roster loss
  - faction lockout
- [ ] Add tests for turn and night transitions.

## Milestone 9: TUI Gameplay Screens

- [ ] Implement dashboard as the primary command center.
- [ ] Implement jobs screen.
- [ ] Implement routing screen.
- [ ] Implement runners screen.
- [ ] Implement factions screen.
- [ ] Implement messages screen.
- [ ] Implement log screen.
- [ ] Implement help screen.
- [ ] Add screen-specific key maps.
- [ ] Add empty states for no jobs, no messages, and no active assignments.
- [ ] Add confirmation prompts for risky actions.
- [ ] Add compact risk-factor display.
- [ ] Add clear current-action prompt.
- [ ] Add visual indication of bundled jobs.

## Milestone 10: Save System

- [ ] Define save file schema.
- [ ] Add save version field.
- [ ] Save current run as JSON.
- [ ] Load current run from JSON.
- [ ] Autosave after turn.
- [ ] Handle missing or corrupt save files gracefully.
- [ ] Add seeded run support.
- [ ] Add tests for save/load round trip.

## Milestone 11: Offline Fiction Layer

- [ ] Add fallback job titles.
- [ ] Add fallback client messages.
- [ ] Add fallback runner messages.
- [ ] Add fallback faction messages.
- [ ] Add fallback after-action reports.
- [ ] Add fallback message selection by cargo, faction, result, and tone.
- [ ] Keep fallback text short enough for TUI panels.
- [ ] Add tests for fallback availability.

## Milestone 12: Optional LLM Layer

- [ ] Define provider-agnostic LLM client interface.
- [ ] Add config loading from file and environment variables.
- [ ] Add offline mode switch.
- [ ] Add job dressing prompt.
- [ ] Add NPC and runner message prompt.
- [ ] Add after-action report prompt.
- [ ] Require structured JSON responses.
- [ ] Validate JSON fields.
- [ ] Validate max text lengths.
- [ ] Reject unauthorized mechanics, rewards, choices, districts, and outcomes.
- [ ] Retry failed calls once.
- [ ] Fall back to offline text after validation or network failure.
- [ ] Add tests with fake LLM client.

## Milestone 13: Balance And Polish

- [ ] Tune runner stats.
- [ ] Tune district stats.
- [ ] Tune cargo risks.
- [ ] Tune payout ranges.
- [ ] Tune heat gain and decay.
- [ ] Tune stress gain and recovery.
- [ ] Tune bundle penalties.
- [ ] Tune operating costs.
- [ ] Ensure no single runner dominates every route.
- [ ] Ensure no single route dominates every job.
- [ ] Improve panel sizing across narrow and wide terminals.
- [ ] Improve help text.
- [ ] Add release build instructions.
- [ ] Add sample config.

## First Playable Acceptance Criteria

- [ ] Player can start a new run.
- [ ] Player can see the city, jobs, runners, messages, and status.
- [ ] Player can accept jobs.
- [ ] Player can assign runners.
- [ ] Player can bundle up to two compatible jobs on one runner.
- [ ] Player can choose routes.
- [ ] Player can advance turns.
- [ ] Jobs resolve through deterministic rules.
- [ ] Outcomes update credits, heat, stress, factions, cargo, and dispatch integrity.
- [ ] Logs explain outcomes using contributing factors.
- [ ] Player can win a seven-night run.
- [ ] Player can lose through defined failure conditions.
- [ ] Game is playable without an LLM.
- [ ] `go test ./...` passes.

## MVP Acceptance Criteria

- [ ] First playable acceptance criteria are complete.
- [ ] Save and load work.
- [ ] Autosave works after each turn.
- [ ] At least ten job templates are present.
- [ ] At least twelve complication types are present.
- [ ] At least thirty fallback messages are present.
- [ ] At least twenty fallback after-action reports are present.
- [ ] Optional LLM job dressing works when configured.
- [ ] Optional LLM after-action reports work when configured.
- [ ] LLM failures never block play.
- [ ] README explains install, run, config, save files, and offline mode.
