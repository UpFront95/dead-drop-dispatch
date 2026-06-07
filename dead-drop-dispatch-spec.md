# Cyberpunk Courier Dispatch TUI Game Spec

Working title: **Dead Drop Dispatch**

## 1. Concept

**Dead Drop Dispatch** is a terminal user interface game built with Bubble Tea in Go. The player runs a small illegal courier dispatch desk in a cyberpunk city, assigning runners to jobs, choosing routes, managing heat, handling client messages, and surviving faction pressure.

The player is not a gunfighter or street samurai. The player is the operator behind the terminal: routing people, packages, lies, and consequences through a city that watches everything.

The terminal is not just the interface. It is the fiction.

## 2. MVP Goal

The MVP should prove the core loop:

1. The player receives courier jobs.
2. The player assigns runners.
3. The player chooses routes.
4. The city simulation resolves risk, delay, heat, faction impact, runner stress, and cargo status.
5. The player receives messages and after-action reports.
6. The next turn begins with a changed city state.

The MVP should feel playable in a single terminal session without requiring a large content set, procedural sprawl, or a fully open-ended LLM game master.

The target MVP should be small, sharp, and replayable.

## 3. MVP Design Pillars

### 3.1 Dispatch, Not Direct Control

The player makes decisions from the desk. They do not move a character tile-by-tile. They assign, route, negotiate, delay, abandon, split, or reroute jobs.

### 3.2 Consequence Over Combat

Conflict appears as risk, reputation damage, heat, injured runners, lost cargo, burned contacts, faction retaliation, and changed city conditions.

### 3.3 Legible Simulation

The player should understand enough of the system to make informed choices. The game can hide some intel, but the simulation should not feel arbitrary.

### 3.4 TUI-Native Fiction

The game should use panels, logs, tables, progress bars, status indicators, message feeds, and ASCII maps as part of the world.

### 3.5 LLM as Fiction Layer, Not Rules Engine

The LLM may generate job flavor, NPC messages, runner dialogue, and after-action reports. It should not decide mechanical outcomes directly. The game engine remains authoritative.

## 4. MVP Scope

The MVP includes:

### 4.1 One City Sector

The MVP contains one playable city sector with six districts.

Example districts:

| District | Description | Mechanical Traits |
|---|---|---|
| Northline | Corporate transit spine and checkpoint corridor | High surveillance, fast travel, high police risk |
| Floodglass | Low streets, tunnels, pump stations, illegal clinics | Medium risk, slow routes, useful for hiding |
| Saint Orison Market | Dense bazaar, shrines, noodle stalls, black-market kiosks | Low surveillance, high faction activity |
| Port Kestrel | Cargo yards, drone cranes, container stacks | Good for contraband, gang risk |
| Ashgate Yard | Rail spurs, furnace stacks, and union-controlled salvage lanes | High traffic, weak signal, useful industrial routing |
| Crown Verge | Luxury towers and private security zones | High payout, high corporate security |

Each district has the following MVP stats:

| Stat | Purpose |
|---|---|
| Surveillance | Increases chance of scans, trace events, and police heat |
| Traffic | Increases delay risk |
| Faction Control | Determines which faction has influence there |
| Danger | Increases runner injury or job failure risk |
| Signal Quality | Affects quality of intel and AI assistant reliability |

### 4.2 Three Runners

The MVP starts with three runners. Each runner has a name, movement style, strengths, weakness, stress, injury state, loyalty, and one personal trait.

Example runners:

| Runner | Style | Strength | Weakness |
|---|---|---|---|
| Mira Vale | Bike courier | Fast through crowded streets | High stress gain under surveillance |
| Kaito Senn | Tunnel runner | Safer in Floodglass and service corridors | Slow through corporate districts |
| Vex Calder | Social operator | Better at checkpoints and negotiations | Higher chance of betraying or side-dealing if loyalty drops |

Runner MVP stats:

| Stat | Purpose |
|---|---|
| Speed | Reduces delay chance |
| Stealth | Reduces detection chance |
| Nerve | Reduces stress from dangerous jobs |
| Talk | Improves checkpoint and negotiation events |
| Loyalty | Affects willingness to take risks and future complications |
| Stress | Increases chance of mistakes, refusal, or breakdown |

Runner states:

| State | Effect |
|---|---|
| Ready | Can take jobs |
| On Job | Currently assigned |
| Injured | Cannot take jobs for a number of turns |
| Burned | Temporarily unsafe in one or more districts |
| Missing | Unavailable until a recovery event resolves |

### 4.3 Four Factions

The MVP includes four factions with reputation and suspicion values.

Example factions:

| Faction | Description |
|---|---|
| Helix Municipal Security | Police-adjacent surveillance and enforcement body |
| Kestrel Dock Union | Smugglers, workers, cargo handlers, and strike captains |
| Saint Orison Clinic Network | Underground medics and biological cargo brokers |
| Asterion Systems | Corporate client, data broker, and security contractor |

Faction variables:

| Variable | Purpose |
|---|---|
| Reputation | Higher reputation unlocks better jobs, safer passage, and favors |
| Suspicion | Higher suspicion causes retaliation, traps, surveillance, or bad terms |

### 4.4 Five Cargo Types

The MVP has five cargo categories.

| Cargo Type | Gameplay Effect |
|---|---|
| Data Shard | High trace risk, low physical danger |
| Medical Cooler | Degrades over time, clinic reputation impact |
| Witness | Can panic, talk, or refuse route changes |
| Contraband Package | High police and faction risk |
| Corporate Prototype | High payout, high betrayal and interception risk |

Each job combines cargo type, origin district, destination district, deadline, payout, faction source, and special modifier.

### 4.5 Job Templates

The MVP should ship with around ten job templates. These can be mechanically fixed but cosmetically varied by the LLM.

Example templates:

| Template | Core Challenge |
|---|---|
| Clinic Rush | Deliver a medical cooler before integrity drops |
| Quiet Data Hand-Off | Move a data shard while trace risk rises |
| Witness Transfer | Move a person through hostile or surveilled territory |
| Prototype Lift | High-value package with corporate pursuit risk |
| Union Favor | Low payout but improves dock reputation |
| Dirty Evidence | Delivery hurts one faction while helping another |
| Decoy Package | The job may be bait or surveillance setup |
| Split Route | Two legs, possible handoff between runners |
| Curfew Run | Deadline before district lockdown |
| Bad Client | Strong payout but increased chance of nonpayment or betrayal |

### 4.6 Routes

For MVP, routes can be district-to-district paths rather than precise map movement.

Each job has two to four possible route options.

Example route types:

| Route Type | Effect |
|---|---|
| Main Artery | Fast, high surveillance |
| Service Tunnels | Slow, lower surveillance, higher danger |
| Market Weave | Medium speed, faction-dependent |
| Drone Corridor | Fast for certain cargo, high trace risk |
| Floodline | Low visibility, delay risk, medical access bonus |

The route resolver calculates:

| Outcome | Description |
|---|---|
| Success | Cargo delivered |
| Delay | Deadline pressure increases or job fails if too late |
| Detection | Heat rises, faction suspicion changes |
| Complication | A mid-job choice appears |
| Injury | Runner becomes injured or unavailable |
| Cargo Damage | Cargo integrity drops |
| Betrayal or Interception | Rare event for high-risk jobs |

### 4.7 Turn Structure

The MVP uses discrete turns. One turn represents a dispatch cycle, roughly thirty to sixty in-world minutes.

A turn proceeds as follows:

1. New messages arrive.
2. Available jobs are shown.
3. Player assigns runners and routes.
4. Player optionally sends short responses to clients or runners.
5. Active jobs resolve.
6. Complications are presented.
7. Rewards, damage, heat, reputation, and stress update.
8. After-action reports appear.
9. City state updates.

### 4.8 Core Resources

The MVP tracks:

| Resource | Meaning |
|---|---|
| Credits | Money used for upgrades, bribes, treatment, and survival |
| Heat | Global law-enforcement pressure |
| Dispatch Integrity | Operational health of the player’s outfit |
| Runner Stress | Individual runner pressure |
| Faction Reputation | Access and favor |
| Faction Suspicion | Retaliation and risk |
| Cargo Integrity | Job-specific condition |
| Deadline | Job-specific time pressure |

### 4.9 Failure and Victory Conditions

MVP victory condition:

Survive seven nights and reach a credit target while keeping dispatch integrity above zero.

MVP failure conditions:

| Failure | Trigger |
|---|---|
| Bankrupt | Credits below required operating threshold |
| Burned | Heat reaches maximum |
| Collapse | Dispatch integrity reaches zero |
| Roster Loss | All runners unavailable |
| Faction Lockout | Too many factions become hostile at once |

The MVP should support short runs of approximately thirty to sixty minutes.

## 5. TUI Interface Spec

The game should use Bubble Tea, Lip Gloss, and Bubbles components.

### 5.1 Main Dashboard Layout

The main screen should contain:

| Panel | Contents |
|---|---|
| City Map | District list or ASCII route map |
| Job Board | Available jobs, deadlines, payouts, faction source |
| Runner Roster | Runner states, stress, injury, loyalty |
| Message Feed | Client messages, runner updates, system alerts |
| Detail Pane | Selected job, route, runner, or faction details |
| Status Bar | Credits, heat, day/night, dispatch integrity |

### 5.2 Primary Screens

| Screen | Purpose |
|---|---|
| Dashboard | Main command center |
| Jobs | Browse and accept jobs |
| Routing | Pick runner and route |
| Runners | Inspect roster |
| Factions | Inspect reputation and suspicion |
| Messages | Read and respond to messages |
| Log | Review previous events |
| Settings | LLM/API configuration and display options |

### 5.3 Input Model

Suggested controls:

| Key | Action |
|---|---|
| Tab / Shift+Tab | Cycle panels |
| Arrow keys / hjkl | Move selection |
| Enter | Select / confirm |
| Esc | Back / cancel |
| j | Open jobs |
| r | Open runners |
| f | Open factions |
| m | Open messages |
| l | Open log |
| Space | Advance turn or resolve selection |
| ? | Help |
| q | Quit |

### 5.4 Visual Style

The interface should feel like an illegal dispatch console rather than a generic CLI.

Style goals:

| Element | Direction |
|---|---|
| Borders | Thin, practical, utilitarian |
| Colors | Muted terminal palette with limited accent colors |
| Text | Short, dense, readable |
| Alerts | Clear but not noisy |
| Animation | Minimal spinners, typing effects, signal recovery effects |
| Logs | Timestamped diegetic terminal output |

## 6. LLM Integration for MVP

### 6.1 MVP LLM Features

The MVP should include three LLM-powered systems.

#### 6.1.1 Job Dressing

The game engine creates a structured job. The LLM writes a short in-world description.

Input to LLM:

```json
{
  "job_type": "Medical Cooler",
  "origin": "Saint Orison Market",
  "destination": "Crown Verge",
  "deadline_turns": 3,
  "payout": 420,
  "faction": "Saint Orison Clinic Network",
  "risk_tags": ["perishable", "checkpoint risk", "high-value client"]
}
```

Output from LLM:

```json
{
  "title": "Blue Organ Run",
  "client_message": "Clinic courier went dark near the market shrine. Cooler is sealed, cold, and already late. Crown Verge buyer pays if the tissue still scans clean.",
  "public_summary": "Move a medical cooler from Saint Orison Market to Crown Verge before integrity drops."
}
```

#### 6.1.2 NPC and Runner Messages

The LLM writes messages from clients, runners, factions, and the in-game AI assistant.

The game engine provides:

| Input | Purpose |
|---|---|
| Speaker | Who is sending the message |
| Tone | Nervous, angry, cold, friendly, encrypted, etc. |
| Context | What happened mechanically |
| Allowed Claims | Facts the message may reference |
| Allowed Choices | Choices available to the player |

The LLM must not invent new actions, mechanics, districts, rewards, or outcomes.

#### 6.1.3 After-Action Reports

After job resolution, the LLM writes a short report in one of several formats:

| Format | Example |
|---|---|
| Runner voice memo | A tired courier recap |
| Police bulletin | Sanitized enforcement report |
| Client receipt | Cold transactional note |
| Street rumor | Unreliable public chatter |
| Dispatch ledger | Internal operational summary |

### 6.2 Optional Freeform Text Input in MVP

The MVP may include simple freeform negotiation, but only if it remains classification-based.

The player can type a response. The LLM classifies intent into a fixed action.

Allowed classified intents:

| Intent | Engine Action |
|---|---|
| Accept | Accept job or terms |
| Refuse | Decline |
| Ask More Pay | Run payment negotiation check |
| Ask More Info | Reveal extra job detail if available |
| Threaten | Reputation or suspicion check |
| Reassure | Stress or loyalty check |
| Stall | Delay decision |
| Deceive | Risky manipulation check |

The engine decides results. The LLM writes the reply after the result.

### 6.3 LLM Safety Rails

The LLM must return structured JSON for all gameplay-relevant calls.

The engine validates:

| Validation | Reason |
|---|---|
| No new mechanics | Prevent hallucinated game systems |
| No unauthorized rewards | Prevent model-created payouts |
| No unauthorized choices | Preserve engine control |
| Max length limits | Keep TUI readable |
| Tone constraints | Preserve setting consistency |
| Retry or fallback | Game remains playable offline or if the model fails |

### 6.4 Offline Mode

The MVP should be playable without an LLM.

Fallback content should include:

| Content | Fallback |
|---|---|
| Job titles | Template-based strings |
| Client messages | Prewritten snippets |
| Runner messages | Prewritten snippets |
| After-action reports | Deterministic summaries |

The LLM should enhance the experience but not be required for the game to function.

## 7. Technical Architecture

### 7.1 Language and Libraries

| Area | Suggested Tool |
|---|---|
| Language | Go |
| TUI Framework | Bubble Tea |
| Styling | Lip Gloss |
| Components | Bubbles |
| Config | YAML or TOML |
| Save Data | JSON or SQLite |
| LLM API | Provider-abstracted client |
| Testing | Go test |

### 7.2 Project Structure

Suggested structure:

```text
dead-drop-dispatch/
  cmd/
    ddd/
      main.go
  internal/
    app/
      model.go
      update.go
      view.go
    game/
      state.go
      jobs.go
      runners.go
      factions.go
      routes.go
      resolver.go
      events.go
    tui/
      layout.go
      styles.go
      screens.go
      components.go
    llm/
      client.go
      prompts.go
      schemas.go
      fallback.go
    content/
      districts.go
      templates.go
      names.go
    save/
      save.go
      load.go
  assets/
    sample-config.toml
  docs/
    design.md
  README.md
```

### 7.3 Game State Model

Core state:

```go
type GameState struct {
    Turn              int
    Night             int
    Credits           int
    Heat              int
    DispatchIntegrity int
    Districts         []District
    Runners           []Runner
    Factions          []Faction
    AvailableJobs     []Job
    ActiveJobs        []ActiveJob
    Messages          []Message
    EventLog          []LogEntry
    RandomSeed        int64
}
```

### 7.4 Deterministic Resolver

The resolver should be independent of the TUI and LLM.

Inputs:

| Input | Description |
|---|---|
| Job | Cargo, deadline, faction, payout |
| Runner | Stats, stress, state |
| Route | Speed, surveillance, danger |
| District State | Current city modifiers |
| Faction State | Reputation and suspicion |
| Random Seed | Reproducible outcomes |

Outputs:

| Output | Description |
|---|---|
| JobResult | Success, failure, partial, delayed |
| State Changes | Heat, credits, stress, reputation |
| Events | Complications and logs |
| LLM Context | Summary for message/report generation |

### 7.5 Save System

MVP save support:

| Feature | Requirement |
|---|---|
| Save current run | Required |
| Load current run | Required |
| Autosave after turn | Recommended |
| Export run log | Optional |
| Seeded runs | Recommended |

## 8. MVP Content Requirements

Minimum content counts:

| Content Type | Count |
|---|---|
| Districts | 5 |
| Runners | 3 |
| Factions | 4 |
| Cargo Types | 5 |
| Job Templates | 10 |
| Route Types | 5 |
| Complication Types | 12 |
| Message Fallbacks | 30 |
| After-Action Fallbacks | 20 |

### 8.1 MVP Complication Types

Suggested complications:

| Complication | Player Choice |
|---|---|
| Checkpoint | Talk through, reroute, bribe, abandon |
| Scanner Sweep | Hide cargo, rush through, spoof tag |
| Runner Panic | Reassure, order forward, abort |
| Cargo Leak | Continue, seek clinic, dump cargo |
| Gang Toll | Pay, negotiate, detour |
| Client Changes Terms | Accept, demand more, cancel |
| Drone Tail | Go dark, switch route, bait tail |
| Signal Loss | Trust route, wait, reroute |
| Witness Refuses | Calm, threaten, sedate if available |
| Data Trace | Scrub, hurry, handoff |
| Curfew Drop | Sprint, hide, postpone |
| Rival Courier | Race, cooperate, sabotage |

## 9. MVP Progression

### 9.1 Run Length

A run lasts seven nights.

Each night has a fixed number of turns, such as six turns per night.

### 9.2 Upgrades

The MVP should include a small upgrade shop between nights.

Example upgrades:

| Upgrade | Effect |
|---|---|
| Signal Relay | Improves signal quality in one district |
| Safehouse | Reduces runner stress after jobs through one district |
| Fake Credential Printer | Improves checkpoint outcomes |
| Clinic Favor | Reduces injury recovery time |
| Dead-Drop Locker | Reduces risk for contraband jobs |
| Scrambler | Reduces data trace risk |

### 9.3 Economy

Credits are gained from completed jobs and spent on upgrades, treatment, bribes, and operating costs.

Each night ends with an operating cost. This creates pressure to take riskier work.

## 10. Development Milestones

### Milestone 1: Terminal Skeleton

Deliverables:

| Deliverable | Description |
|---|---|
| Bubble Tea app shell | Main loop, model/update/view |
| Layout | Dashboard with panels |
| Navigation | Panel focus and keyboard commands |
| Static content | Districts, runners, factions shown |

### Milestone 2: Core Simulation

Deliverables:

| Deliverable | Description |
|---|---|
| Job generation | Template-based jobs |
| Runner assignment | Assign runner and route |
| Resolver | Resolve job outcomes |
| State updates | Credits, heat, stress, reputation |
| Event log | Deterministic logs |

### Milestone 3: MVP Run Loop

Deliverables:

| Deliverable | Description |
|---|---|
| Turn advancement | Full turn cycle |
| Night advancement | End-of-night state changes |
| Failure conditions | Bankrupt, burned, collapse, roster loss |
| Victory condition | Survive seven nights |
| Save/load | Basic persistence |

### Milestone 4: LLM Layer

Deliverables:

| Deliverable | Description |
|---|---|
| LLM client abstraction | Provider-agnostic interface |
| Job dressing | LLM-generated job text |
| NPC messages | LLM-generated messages from structured context |
| After-action reports | LLM-generated reports |
| Fallbacks | Offline template mode |

### Milestone 5: Polish Pass

Deliverables:

| Deliverable | Description |
|---|---|
| UI refinement | Better tables, colors, help screen |
| Balance pass | Risk, payout, heat, stress |
| Content pass | More fallbacks and event variants |
| Packaging | README, install instructions, release binary |

## 11. MVP Non-Goals

The MVP should not include:

| Non-Goal | Reason |
|---|---|
| Full open-world city simulation | Too large for first playable version |
| Direct character movement | Dispatch fantasy is stronger |
| Real-time action | Turn structure is easier to build and balance |
| LLM as unrestricted game master | Too unstable and hard to test |
| Multiplayer | Not needed for core loop |
| Complex combat system | Consequence-based conflict fits better |
| Huge procedural lore generator | Start with focused content |
| Networked persistent world | Premature infrastructure burden |

## 12. Optional Future Content

The following systems should be considered future expansions after the MVP proves the core loop.

### 12.1 Expanded City

Add more districts with deeper traits.

Future district ideas:

| District | Concept |
|---|---|
| Glass Cathedral | Corporate religious wellness complex |
| Blackrail | Abandoned maglev spine used by smugglers |
| Neon Marsh | Flooded outskirts with illegal bio-labs |
| Parliament of Wires | AI cult enclave and server monastery |
| Old Financial Core | Dead trading towers and automated security |
| Low Sun Estates | Luxury enclave with private armies |
| Underbridge Nine | Refugee market below suspended expressways |

Additional district variables:

| Variable | Effect |
|---|---|
| Curfew Status | Blocks or penalizes travel |
| Riot Level | Raises danger but lowers surveillance |
| Blackout Level | Lowers signal and camera risk |
| Medical Access | Helps injuries and biological cargo |
| Smuggling Density | Helps contraband but raises faction entanglement |
| Weather Exposure | Affects drones, bikes, and cargo integrity |

### 12.2 Larger Runner Roster

Add recruitable runners with personal arcs, rivalries, injuries, favors, vices, and permanent scars.

Possible runner traits:

| Trait | Effect |
|---|---|
| Ex-Cop | Better at checkpoints, hated by gangs |
| Ghostware Implant | Harder to track, unstable under signal noise |
| Debt-Bound | Takes bad jobs unless managed |
| Local Saint | Protected in one district |
| Chrome Allergy | Bad with prototype cargo |
| Union Blood | Better dock outcomes, faction obligations |
| Nerve Burn | Powerful under danger, stress recovery is poor |

### 12.3 Runner Personal Arcs

Each runner can develop storylines based on repeated decisions.

Example arcs:

| Arc | Trigger |
|---|---|
| The Debt Collector | Runner’s old debt resurfaces |
| Dead Sibling Signal | Runner receives messages from impossible source |
| Burned Identity | Runner needs new papers |
| Revenge Job | Runner asks for off-ledger route |
| Clinic Dependency | Injury treatment creates faction obligation |
| Loyalty Break | Runner questions player after repeated high-risk assignments |

LLM role:

The LLM can write messages and scenes based on compact memory facts. The engine controls actual consequences.

### 12.4 Seasonal Campaigns

Instead of a single seven-night run, add seasons with city-level crises.

Possible season arcs:

| Season Arc | Mechanical Impact |
|---|---|
| Corporate Merger | More prototype jobs, more private security |
| Plague Rumor | Medical cargo demand rises, checkpoints tighten |
| Dock Strike | Port routes destabilize, union jobs increase |
| Election Week | Police visibility rises, bribes become expensive |
| AI Blackout | Signal quality collapses, strange jobs appear |
| Monsoon Flood | Low districts become dangerous but stealthy |
| Celebrity Assassination | Witness and evidence jobs spike |

### 12.5 In-Game AI Assistant

Add an illegal dispatch assistant, such as ORACLE/7, DOGSTAR, or Little Saint.

Functions:

| Function | Gameplay Role |
|---|---|
| Route Advice | Suggests safest or fastest route |
| Suspicion Analysis | Flags likely bait jobs |
| Runner Warnings | Notes stress and loyalty risks |
| Message Drafting | Suggests client replies |
| Pattern Detection | Identifies faction strategy |
| Unreliable Inference | Can be wrong if signal is poor or compromised |

Possible twist:

The assistant may develop its own agenda, become compromised, or attract attention from AI cult factions.

### 12.6 Freeform Negotiation

Allow players to type natural language responses to clients, runners, and factions.

Important constraint:

The LLM classifies player intent. The engine resolves the outcome.

Expanded intent set:

| Intent | Use |
|---|---|
| Bluff | Pretend to have leverage |
| Appeal | Ask based on relationship |
| Delay | Buy time |
| Redirect | Push job to another party |
| Threaten | Increase pressure |
| Confess | Admit failure |
| Conceal | Hide partial failure |
| Bargain | Trade pay, time, route, or favor |
| Extract Intel | Ask probing questions |

### 12.7 Memory System

Add persistent NPC and faction memory.

Memory examples:

| Memory | Later Effect |
|---|---|
| Player abandoned Mira at checkpoint | Mira loyalty drops, future tone hardens |
| Player saved clinic cargo at personal cost | Clinic grants treatment discount |
| Player lied to Asterion | Corporate suspicion rises |
| Player paid runners extra after bad night | Stress recovery improves |
| Player repeatedly used tunnels | Tunnel faction starts charging tolls |

### 12.8 Expanded Cargo

Future cargo types:

| Cargo | Special Rules |
|---|---|
| Wetware Imprint | Requires cold signal isolation |
| Memory Pearl | Can corrupt runner dreams |
| Organ Case | Perishable and morally loaded |
| Synthetic Pet | Screams if scanned |
| Political Evidence | Changes faction balance |
| AI Fragment | May communicate with dispatch |
| Luxury Narcotic | High payout, high gang risk |
| Stolen Body | Requires route secrecy and clinic access |
| Diplomatic Bag | Legal immunity until tampered with |
| Debt Ledger | Can be sold, delivered, copied, or destroyed |

### 12.9 Infrastructure Layer

Add persistent dispatch upgrades and physical network assets.

Future infrastructure:

| Asset | Effect |
|---|---|
| Safehouses | Rest, hide cargo, reduce heat |
| Signal Relays | Better route intel |
| Bribe Channels | Lower checkpoint risk |
| Dead-Drop Lockers | Enable package handoffs |
| Drone Nests | Allow unmanned delivery |
| Med Bay | Heal runners |
| Credential Printer | Generate fake access passes |
| Packet Scrubber | Lower data trace |
| Chop Shop Contact | Modify vehicles |
| Rumor Broker | Reveal job hidden traits |

### 12.10 More Complex Routing

Future routing can include multi-leg jobs, runner handoffs, decoys, split cargo, parallel routes, and timed coordination.

Future route features:

| Feature | Description |
|---|---|
| Multi-leg Route | One job crosses several route nodes |
| Handoff | Transfer cargo between runners |
| Decoy | Send fake cargo to draw heat |
| Shadow Route | Runner follows another for protection |
| Dead Zone | No messages during route |
| Forced Stop | Delay or complication node |
| Hidden Shortcut | Unlocked by faction or infrastructure |
| Route Memory | Repeated use increases detection |

### 12.11 Rival Dispatchers

Add competing courier outfits.

Possible rival behaviors:

| Rival | Behavior |
|---|---|
| CleanLine | Corporate-friendly, steals premium contracts |
| Red Moth | Violent, sabotages runners |
| Saint Bicycle Choir | Idealistic, helps civilians |
| Null Cartel | AI-assisted and unpredictable |
| Paper Tiger | Cheap, unreliable, useful scapegoat |

Interactions:

| Interaction | Effect |
|---|---|
| Race | Compete for urgent delivery |
| Sabotage | Raise risk for rival |
| Partnership | Split jobs |
| Blame Shift | Move heat to rival |
| Buyout | Recruit rival runner |
| Feud | Repeated retaliation |

### 12.12 News and Rumor System

Add a city feed that reacts to player action.

Feed formats:

| Format | Tone |
|---|---|
| Police Bulletin | Official, sanitized |
| Street Board | Rumor-heavy |
| Corporate Memo | Cold and euphemistic |
| Pirate Radio | Dramatic and unreliable |
| Clinic Notice | Practical and urgent |
| Courier Forum | Bitter, funny, tactical |

### 12.13 Moral Pressure

Add jobs where success is not obviously good.

Examples:

| Job | Dilemma |
|---|---|
| Deliver Witness | Client may kill the witness |
| Move Medicine | Destination is rich buyer, not clinic |
| Carry Evidence | Delivery could start crackdown |
| Prototype Return | Corporate job suppresses whistleblower |
| Debt Ledger | Destroying cargo helps workers but angers client |
| AI Fragment | Delivering it may awaken dangerous system |

### 12.14 Endgame Arcs

Longer campaigns can culminate in major city changes.

Possible endings:

| Ending | Condition |
|---|---|
| Ghost Network | Dispatch becomes untraceable but isolated |
| Corporate Asset | Player survives by serving Asterion |
| Street Legend | High street trust and low betrayal |
| Burn Notice | Player escapes after citywide crackdown |
| Little Saint Ascendant | AI assistant becomes major power |
| Union City | Dock and street factions reshape routes |
| Quiet Exit | Player sells network and disappears |

## 13. Design Risks

### 13.1 LLM Scope Creep

Risk:

The LLM starts controlling too much of the game.

Mitigation:

Keep mechanics deterministic. Use structured JSON. Validate outputs. Preserve fallback content.

### 13.2 Interface Overload

Risk:

The TUI becomes dense and confusing.

Mitigation:

Start with five panels only. Add detail screens. Keep alerts clear. Make the current required action obvious.

### 13.3 Simulation Opacity

Risk:

The player cannot understand why outcomes happen.

Mitigation:

Show risk summaries, route traits, and post-result explanations.

### 13.4 Content Burden

Risk:

The game needs too many authored jobs and messages.

Mitigation:

Use templates, tags, and LLM dressing. Start with small content counts.

### 13.5 Balance Problems

Risk:

The player finds one dominant route or runner.

Mitigation:

Use changing district states, faction pressure, stress, route memory, and cargo-specific risks.

## 14. First Playable Prototype

The first playable prototype should include:

| Feature | Requirement |
|---|---|
| One dashboard | Shows jobs, runners, messages, status |
| Three runners | Assignable to jobs |
| Five districts | Used in route risk |
| Five jobs per night | Generated from templates |
| Basic resolver | Success, delay, detection, injury |
| Seven-night run | Win/loss state |
| Event log | Explains outcomes |
| Offline text | No LLM required |
| Optional LLM job text | If API key configured |

The prototype is successful if the player can complete a run, understand why things happened, and want to try a different routing strategy next time.

## 15. Tone Reference

The tone should be grounded, sharp, and restrained.

Avoid:

| Avoid | Reason |
|---|---|
| Wacky cyberpunk parody | Undercuts tension |
| Endless lore dumps | Slows play |
| Generic hacker jargon | Makes the world feel cheap |
| Overwritten NPCs | TUI needs compressed text |
| Player-as-chosen-one plot | Dispatch fantasy is stronger |

Prefer:

| Prefer | Reason |
|---|---|
| Short messages | Terminal-readable |
| Concrete details | Makes jobs memorable |
| Moral ambiguity | Adds pressure |
| Operational language | Fits dispatch role |
| Consequence-driven fiction | Supports replayability |

Example tone:

> “Cooler is still cold. Runner is not. Mira came back shaking, jacket wet to the elbows, asking whether Saint Orison pays extra for silence.”

## 16. Summary

The MVP should be a compact turn-based cyberpunk operations sim in a terminal. The player manages runners, routes, cargo, heat, factions, and limited resources over seven nights. The LLM provides voice, variation, and atmosphere, but the game engine owns the rules.

Build the desk first. Let the city crawl in through the wires later.
