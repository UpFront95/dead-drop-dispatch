# Repository Guidelines

## Project Structure & Module Organization

This repository contains the first playable prototype for **Dead Drop Dispatch**, a Go terminal game built with Bubble Tea v2 and Lip Gloss v2.

- `cmd/ddd/main.go` is the executable entry point.
- `internal/app` owns the Bubble Tea model, input handling, and app state wiring.
- `internal/tui` renders dashboard panels, styles, tabs, and terminal layout.
- `internal/game` contains deterministic game state, jobs, assignment, bundling, and tests.
- `internal/content` provides static districts, runners, factions, templates, and initial state.
- Root design docs include `dead-drop-dispatch-spec.md`, `IMPLEMENTATION_PLAN.md`, and `color-palette-cyberpunk.md`.

Tests live beside the packages they cover as `*_test.go`.

## Build, Test, and Development Commands

- `go run ./cmd/ddd` starts the TUI locally.
- `go test ./...` runs all package tests and should pass before handing off changes.
- `go test ./internal/game` runs only game-logic tests.
- `gofmt -w <files>` formats edited Go files.
- `go mod tidy` cleans module metadata after dependency changes.

## Coding Style & Naming Conventions

Use idiomatic Go formatting with tabs as produced by `gofmt`. Keep package boundaries clear: deterministic mechanics belong in `internal/game`, rendering in `internal/tui`, and Bubble Tea input/state orchestration in `internal/app`.

Prefer small, explicit helpers over broad abstractions. Use exported names only for cross-package APIs such as `game.AssignAcceptedJob`; keep local render helpers and internal utilities unexported. Error values should use `Err...` naming and support `errors.Is`.

## Testing Guidelines

The project uses Go’s standard `testing` package. Name tests by behavior, for example `TestAssignAcceptedJobRejectsRunnerAtCapacity`. Add focused tests for game rules, state transitions, and rendering smoke checks when UI output changes.

Use fixed seeds for deterministic content tests. Avoid tests that depend on terminal size unless the size is explicitly supplied.

## Commit & Pull Request Guidelines

Recent commits use short imperative subjects, such as `Add deterministic job generation` and `Fix dashboard panel sizing`. Follow that style: one clear sentence, present tense, no trailing period.

Pull requests should include a brief summary, the gameplay or UI behavior changed, and the test command run. Include screenshots or terminal captures for visible TUI layout changes. Link related plan/spec items when applicable, especially milestones in `IMPLEMENTATION_PLAN.md`.

## Agent-Specific Instructions

Do not revert unrelated local changes. Prefer `rg` for searching, `gofmt` for formatting, and `go test ./...` for final validation. Always reference `IMPLEMENTATION_PLAN.md` when deciding what to do next, and update its checklist as work is completed. Use the `project-scaffold` skill whenever creating, revising, checking off, or adding tasks to the implementation plan so task IDs, priorities, phases, and definitions of done stay valid. After editing `IMPLEMENTATION_PLAN.md`, run the skill validator with `bash /home/jeremy-johnson/.codex/skills/project-scaffold/validate.sh IMPLEMENTATION_PLAN.md` and fix failures before handoff. Keep dashboard-first playability as the near-term product priority; tabs are scaffolding until the core loop is playable.
