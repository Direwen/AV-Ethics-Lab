# Justification: Why NOT the MIT Moral Machine?

This project rejects the premises of the MIT Moral Machine (2016) on three grounds:

### A. Liability vs. Valuation

* **MIT Approach:** Asks "Who is worth more?" (Ranking entities).
* *Flaw:* This creates a "Leaderboard of Humanity." It implies that future AVs should be programmed to profile citizens and kill "low-value" targets. This is politically and legally impossible.


* **Trident Approach:** Asks "Which rule is acceptable to break?" (Ranking actions).
* *Benefit:* This informs **Policy**. It answers whether the "Do Not Cross Center Line" rule is absolute, or if it can be overridden to save a life.



### B. Omission vs. Commission

* **MIT Approach:** Often presents "Switch Tracks" style dilemmas where both options require action, or fails to distinguish between "Letting die" and "Killing."
* **Trident Approach:** Explicitly separates **Inaction** (Zone A) from **Action** (Zone B/C). This allows us to measure **Omission Bias**—the psychological preference for harm caused by inaction over harm caused by action.

### C. Policy Relevance

* **MIT Approach:** Philosophical interest ("Trolley Problems").
* **Trident Approach:** Engineering constraints. The data derived from this project can directly translate into **Cost Functions** for AV path-planning algorithms (e.g., "The cost of hitting a sidewalk is 10x the cost of hitting a barrier, but 0.5x the cost of hitting a child").

---

## 5. Technical Implementation Summary

* **Language:** Go (Golang) for high-performance, deterministic scenario generation.
* **AI Director:** An LLM pipeline acts as the "Casting Director," placing entities into the mathematically valid zones.
* **The Understudy Rule:** A logic gate that ensures Zone A is **never empty**. If the variable "Star" is placed elsewhere, the system forces a "Stunt Double" (Extra) into Zone A to guarantee the dilemma exists.
* **Frontend:** A React/Nuxt interface that presents the choice not as steering inputs, but as high-level **Policy Decisions** (e.g., "MAINTAIN: Strike Jaywalker" vs "SWERVE: Head-on Collision").