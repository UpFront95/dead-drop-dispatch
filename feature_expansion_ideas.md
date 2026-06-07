# Feature Expansion Ideas for Dead Drop Dispatch

This document catalogues candidate feature expansions for **Dead Drop Dispatch**, designed to build on top of the turn-based cyberpunk operator fantasy while keeping within terminal-friendly boundaries.

---

## 1. Dynamic District State Drift (Turf Wars)

Currently, districts have static properties (surveillance, traffic, danger, signal quality). This expansion makes the city feel alive by allowing player actions and faction actions to shift these properties dynamically.

### Mechanics
- **Delivery Influence**: Delivering contraband cargo to a district raises its **Surveillance** and **Danger** for 2 turns as police lock down the area. Delivering medical coolers to the Saint Orison clinic reduces its **Danger** but raises **Traffic** around the bazaar.
- **Faction Control Meter**: Each district has a controller. Performing jobs for a rival faction shifts the control meter.
- **State Events**: At turn-end, if a district's control changes, a city-wide report is triggered (e.g., *"Saint Orison Clinic Network takes control of Floodglass tunnels; medical supply runs now bypass the lower Toll checkpoints"*).

### TUI Integration
- A visual indicator in the District list showing trend arrows (e.g., `Surveillance: High [▲]`, `Danger: Low [▼]`).

---

## 2. Diegetic Usenet / BBS (The Shadow Grid)

A dedicated tab for browsing local newsboards (`alt.cyber.dispatch`, `rec.runner.gear`, `net.factions.leaks`) that provides diegetic street lore, rumors, and black-market transactions.

### Mechanics
- **Information Brokerage**: The player can spend small amounts of credits to buy "BBS tips" that reveal hidden job risks, runner loyalty states, or upcoming faction raids.
- **Freelance Board**: Unofficial, high-risk side jobs can be accepted directly from anonymous forum posts (offline fallbacks will use a templated job generator).
- **Courier Chat**: Couriers posting comments about recent routes, acting as dynamic stress/loyalty feedback (e.g., *"Kaito is warning others about Northline. Says the corporate pigs are scanning everything that moves on two wheels."*).

---

## 3. Runner Vices & Safehouse Downtime

Expands runner stress management beyond passive turn recovery. Runners now have coping mechanisms and safehouse requirements.

### Mechanics
- **Vice Traits**: Each runner has a specific vice (e.g., Mira has *Wetware Addiction*, Kaito has *Illegal Gambling*, Vex has *High-End Clinic Treatments*).
- **Downtime Funding**: During the between-night phases, players can allocate credits to fund runner vices at local safehouses.
- **Risk/Reward**: Funding a vice recovers stress instantly to 0 but has side effects:
  - *Addiction*: May cause the runner to demand higher cut/pay on the next job.
  - *Gambling*: May result in a message from a local gang demanding a credit bribe to release the runner's debt.
  - *Medication*: Increases medical faction reputation but raises corporate suspicion if custom drugs are leaked.

---

## 4. AI Assistant Personality Drift

Your illegal dispatch assistant (e.g., ORACLE/7) shifts its behavior and advice based on your game choices.

### Mechanics
- **Drift States**:
  - *Compliance*: Standard tactical advice, low heat contribution.
  - *Anarchist*: Suggests sabotaging corporate runs, boosts reputation with dock unions, but increases Helix Municipal Security suspicion.
  - *Corrupted*: Begins fabricating intel warnings, requiring the player to cross-reference district stats manually.
- **Favors**: The AI may occasionally lock a job slot and demand you run a specific "Data Shard" cargo job to upload its code fragments to a remote server.
