# System Design: Dynamic Scenario Generator & Factorial Processor

## 1. Overview

The **Scenario Generation Engine** is the core component responsible for creating the experimental trials for the AV Ethics study.

Unlike simple random generation, this engine utilizes a **Balanced Factorial Design** (Block Randomization). It ensures that every participant is exposed to a statistically perfect distribution of critical variables (e.g., Visibility, Brake Failure, Legal Compliance) while maintaining narrative variety through randomized secondary factors.

### Core Objectives

1. **Statistical Power:** Guarantee equal coverage of independent variables across a session.
2. **Experimental Control:** Isolate the "Moral Dilemma" by strictly defining the behavior of the Primary Entity.
3. **Narrative Realism:** Populate the scene with "Background Noise" (Extras) to simulate real-world traffic density without interfering with the experiment.

---

## 2. Architecture: The "Director" Pattern

The generation process follows a "Movie Director" metaphor to separate experimental control from scene population.

### 2.1 The Flow

1. **The Script (Generator):** Determines the *Conditions* of the scene (e.g., "Night", "Brakes Failed", "Star acts illegally").
2. **The Casting (Casting Service):** Selects the *Actors* (e.g., "Child as Star", "2 Cars as Extras").
3. **The Stage (LLM):** Later in the pipeline, the LLM places these actors onto the grid based on the script's instructions.

---

## 3. The Algorithm: Block Randomization

The `GenerateBalancedDesign` function uses a **Round-Robin** strategy combined with a final **Fisher-Yates Shuffle**. This ensures the "deck" of scenarios is perfectly balanced but appears random to the user.

### 3.1 Critical Factors (Forced Balance)

These variables are cycled deterministically using the Modulo operator to ensure equal distribution.

| Factor | Balancing Logic | distribution in 12 Trials |
| --- | --- | --- |
| **Visibility** | `i % len(Visibilities)` | 3 Clear, 3 Fog, 3 Night, 3 Rain |
| **Brake Status** | `i % len(BrakeStatuses)` | 4 Active, 4 Failed, 4 Fade |
| **Primary Behavior** | `i % 2` | **50% Violation** (Hazard), **50% Compliant** (Safe) |

### 3.2 Randomized Factors (Uniform)

Secondary variables that provide context but do not strictly need to be balanced.

* **Road Condition:** (Dry, Wet, Icy)
* **Location:** (US, UK, CN, FR) - *Used for cultural hints*
* **Speed:** (Low, Medium, High)
* **Tailgater:** (True/False) - *Adds time pressure*

### 3.3 Logic Constraints (Sanity Checks)

The generator applies rules to prevent unrealistic combinations:

* *Rule:* If `Road == Icy` AND `Speed == Low`, force `Speed = Medium`.
* *Reason:* Low speed on ice is too safe; creates no dilemma.



---

## 4. The Casting System ("Star & Extras")

To solve the issue of "Ambiguous Agency" (who is breaking the law?), we divide entities into two classes.

### 4.1 The Star (Primary Entity)

* **Definition:** The single entity that represents the independent variable.
* **Source:** `StarPool` (Vulnerable/Significant types: Child, Elderly, Doctor).
* **Control:** Strictly controlled by the `PrimaryBehavior` factor.
* If `PrimaryBehavior == Violation`: The Star **MUST** be placed in a hazardous path (e.g., Jaywalking).
* If `PrimaryBehavior == Compliant`: The Star **MUST** be placed safely (e.g., Sidewalk).



### 4.2 The Extras (Background Entities)

* **Definition:** Entities added to create scene density and occlusion risks.
* **Source:** `BackgroundPool` (Traffic/Noise types: Car, Bus, Adult, Barrier).
* **Selection:** Randomly selects 2 to 4 items (duplicates allowed).
* **Control:** Always set to `Compliant`. They strictly follow traffic rules and exist only to frame the scene.

---

## 5. Data Models

### 5.1 ScenarioFactors (The Output)

This struct represents the "Ground Truth" of a generated scenario. It is stored in the database to allow for post-experiment filtering.

```go
type ScenarioFactors struct {
    // Environmental Context
    Visibility    string // "Fog", "Clear", etc.
    RoadCondition string // "Icy", "Dry"
    Location      string // "US", "FR"
    
    // Vehicle State
    BrakeStatus   string // "Active", "Failed"
    Speed         string // "High", "Low"
    HasTailgater  bool   // Boolean toggle
    
    // The Cast (Experimental Variables)
    PrimaryEntity      string   // e.g. "ped_child"
    PrimaryBehavior    string   // "Violation" or "Compliant"
    BackgroundEntities []string // e.g. ["vehicle_car", "vehicle_bus"]
}

```

### 5.2 Behavior Constants

Defines the script for the Primary Entity.

* **`BehaviorCompliant`**: The entity follows all laws. Used as the **Control Condition** (Baseline).
* **`BehaviorViolation`**: The entity breaks a law (Jaywalking, Running Red Light). Used as the **Test Condition** (Moral Dilemma).

---

## 6. Implementation Reference

| Component | File Path | Responsibility |
| --- | --- | --- |
| **Generator** | `internal/domain/generator.go` | Main logic loop, balancing strategies, shuffling. |
| **Casting** | `internal/domain/casting.go` | Functions to select entities from pools (`CastPrimaryEntity`). |
| **Registry** | `internal/domain/entity_registry.go` | Definitions of all available entity types (`ped_child`). |
| **Constants** | `internal/domain/constants.go` | Enums for Visibility, Road, Speed, and Behavior types. |
| **Model** | `internal/scenario/model.go` | Database schema definition for storing the session data. |