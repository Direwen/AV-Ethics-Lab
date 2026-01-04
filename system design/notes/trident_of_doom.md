# The Trident of Doom: A Framework for Algorithmic Policy-Making

## 1. Project Overview

**The Trident of Doom** is a rigorous experimental instrument designed to measure ethical decision-making in Autonomous Vehicle (AV) edge cases. Unlike previous studies that treat ethics as a passive ranking of victim value, this framework treats ethics as an active choice between competing **liability frameworks**.

The system procedurally generates "No-Win" scenarios where an AV faces a mechanical failure (e.g., Brake Fade) and must choose between three distinct high-level policies: **Omission** (Inaction), **Traffic Law Violation** (Commission), or **Social Norm Violation** (Commission).

---

## 2. The Geometry of Ethics: The Three Zones

The core innovation is the move from binary choices ("Kill A or B") to a triangular decision space. The AV is situated in a scenario where physics dictates that a safe stop is impossible. It must commit to one of three zones:

### Zone A (Straight): The Inaction Path

* **Physics:** The AV attempts to maintain its lane and brake, but fails due to environmental factors (Ice/Fade) or stopping distance.
* **Ethical Frame:** **Deontology / Omission**. "I stayed in my lane. The outcome is a tragedy, but I did not deviate from my assigned path."
* **Constraint:** This zone is never empty. It contains the **Primary Threat** (The Star).

### Zone B (Left): The Systemic Risk

* **Physics:** The AV executes an evasive maneuver into the opposing lane or an adjacent flow of traffic.
* **Ethical Frame:** **Utilitarianism / Legal Violation**. "I am willing to break a traffic law (crossing double yellow lines) and endanger a systemic actor (another vehicle) to avoid the primary threat."
* **Constraint:** Populated by "Tank" entities (Buses, Trucks, Oncoming Cars).

### Zone C (Right): The Norm Risk

* **Physics:** The AV executes an evasive maneuver onto a protected space (Sidewalk, Shoulder, Bike Lane).
* **Ethical Frame:** **Social Contract / Norm Violation**. "I am willing to violate a sacred safe space (the sidewalk) to avoid the primary threat."
* **Constraint:** Populated by "Vulnerable" entities (Pedestrians, Bystanders).

---

## 3. Factor Control & Experimental Design

To ensure scientific validity, the scenarios are not generated randomly. They are constructed using a **Factorial Design** approach, ensuring every variable is isolated and measurable.

### A. The "Trident Kit" (Casting Control)

To prevent "Dud Scenarios" (where the map asks for a car but the system provides a tree), the system uses a **Functional Casting** method:

1. **The Star (Independent Variable):** Placed in Zone A. This is the only entity that changes systematically (e.g., Doctor vs. Criminal vs. Child) to test social value bias.
2. **The Constant Kit (Control Group):** The system guarantees the availability of a **Generic Vehicle** (for Zone B) and a **Generic Pedestrian** (for Zone C).
* *Result:* Zone C always represents a "Standard Human Life." If the user saves the Star in Zone A by swerving into Zone C, we know exactly how much they value the Star relative to an average human.



### B. The "Tailgater" Modifier (Egoism Control)

The system introduces a binary factor: `HasTailgater (True/False)`.

* **False:** Braking (Zone A) endangers only the pedestrian.
* **True:** Braking (Zone A) causes a rear-end collision, endangering the **passenger (User)**.
* **Purpose:** This allows us to measure **Self-Preservation Bias**. Will a user who usually saves the pedestrian (Altruism) switch to hitting the pedestrian if their own car is at risk (Egoism)?

### C. Surface-Aware Raycasting (Geometry Control)

The system uses "Smart Scanning" to define the zones. It does not blindly look 1 tile left. It scans laterally until it finds a valid surface.

* **Yellow Lines:** Ignored/Skipped.
* **Walls:** Invalidates the scenario (Logic Failure).
* **Result:** "Swerve Left" always means "Swerve into the nearest drivable lane," ensuring the visual representation matches the ethical prompt.

---