# System Design — Roadmap & Future Work

This document captures planned features and design decisions intended to improve reliability, scalability, and observability of the LLM integration and the broader data-collection system.

## Current approach (short-term)

- **Round Robin Key Rotation (implemented / prioritized)**
  - Rationale: During active data collection we need a lightweight, reliable mechanism to reduce the risk of hitting provider rate limits or single-key outages. A round-robin pool rotates through multiple API keys for a provider (e.g., GROQ, OpenRouter), distributing requests more evenly across free-tier keys.
  - Benefits: Simple, low-maintenance, easy to reason about, and sufficient for early-stage experiments where throughput is moderate.
  - Limitations: Does not solve global rate limits, no true pooling/circuit-breaker logic, and limited observability per-key without metrics.

## Planned improvements (short-term / medium priority)

1. **Throttling (per-key and per-tenant)**
   - What: Rate-limit outgoing requests per API key and optionally per tenant/session to avoid provider rate-limit errors.
   - Acceptance criteria: Reproduces reliable request pacing under high load; prevents 429s in normal operation.

2. **Redis Caching Layer**
   - What: Cache repeatable LLM responses (where applicable), precomputed templates, and short-lived artifacts to reduce LLM calls and latency.
   - Acceptance criteria: Demonstrable reduction in LLM calls for repeated/similar requests and improved P95 latency.

3. **Retry & Exponential Backoff**
   - What: Implement retries with jitter for transient failures (e.g., 5xx and network errors) and allow immediate fail-over to the next key when a key consistently returns errors.
   - Acceptance criteria: Reduced error rates for transient faults and no request storms during retry storms.

4. **Monitoring, Metrics & Telemetry**
   - What: Instrument LLM calls, key usage, latency, error rates (4xx/5xx), queued requests, and cache hit/miss ratios.
   - Acceptance criteria: Dashboards for LLM health and alerts for rising error rates or key exhaustion.

5. **Circuit Breaker and Fallback Model**
   - What: Protect the system from cascading failures by stopping calls to a failing key/provider and either switching to another provider or returning a graceful degraded response.
   - Acceptance criteria: Automatic isolation of unhealthy keys and successful fallback to alternative provider or cached result.

6. **Key Management & Rotation Automation**
   - What: Store keys in a secure store (e.g., Vault / Secrets Manager) and provide tooling for adding/removing keys without server restart.
   - Acceptance criteria: Keys can be added via an admin endpoint or configuration update and become available to the pool dynamically.

7. **Quota Tracking & Billing Awareness**
   - What: Track per-key usage to detect exhaustion or billing spikes and optionally throttle or pause keys approaching limits.
   - Acceptance criteria: Notifications or automatic pausing for keys exceeding configured thresholds.

8. **Graceful Degradation & UX Handling**
   - What: Ensure frontend shows deterministic fallback messaging when feedback generation is delayed or fails (e.g., "Feedback generation delayed. We'll email you when ready" or show a cached summary).
   - Acceptance criteria: Clear UX states for degraded LLM service and preserve privacy / consistency.

9. **Integration & Load Testing**
   - What: Simulate production traffic to verify the behavior of round robin rotation, throttling, retries, and caching under load.
   - Acceptance criteria: No significant 5xx spike and acceptable latency under expected maxima.

10. **Security & Secret Handling**
    - What: Use secure secret manager for keys, rotate keys periodically, and ensure `.env` files are not committed.
    - Acceptance criteria: No secrets in the repo and documented key rotation policy.

## Long-term / Nice-to-have

- **Adaptive Scheduling**: Dynamically adjust throttling based on observed error rates and provider feedback headers.
- **Per-experiment Cost Accounting**: Attribute LLM usage to experiments/sessions for research accounting and budgetary controls.
- **ML-based Retry/Selection**: Use historical latency/error data to predict best key or provider for a given request.

## Prioritization suggestion

1. Throttling + Key-level rate limiting (High)
2. Retry & Backoff + Circuit Breaker (High)
3. Monitoring & Alerts (High)
4. Redis Caching Layer (Medium)
5. Key Management Automation (Medium)
6. Integration & Load Testing (Medium)
7. Graceful Degradation & UX (Medium)
8. Long-term enhancements (Low)

---

## LLM Safety & Guardrails (Security & Ethical Constraints)

To reduce the chance of harmful, biased, or sensitive output and to protect user privacy and research integrity, the following guardrails should be implemented around any LLM generation tasks (both scenario generation and feedback generation):

1. **Input Sanitization & Prompt Hardening**
   - Sanitize all user-provided data before including it in prompts (escape JSON, remove control characters).
   - Use structured system prompts and strict template prompts to minimize prompt injection.
   - Avoid sending PII; filter or obfuscate any user-provided fields.
   - Acceptance criteria: No direct user-supplied text appears unescaped in LLM prompts; prompts are generated from templates and validated.

2. **Output Schema Validation & Parsing**
   - Require LLM to return responses in a strict JSON schema when used for structured outputs (e.g., Entities, DilemmaOptions, Feedback summary). Use JSON schema validation and reject or sanitize non-conforming responses.
   - Acceptance criteria: Any non-conforming response is retried or replaced with a safe fallback and logged for review.

3. **Content Filtering & Safety Layers**
   - Run LLM outputs through safety filters for violence, hate, sexual content, PII leakage, and other banned categories before returning results to users or storing in DB.
   - Use provider-side safety endpoints where available and local heuristics (regex PII detectors, toxicity classifiers) as a secondary check.
   - Acceptance criteria: Blocked outputs are not shown to users; a safe fallback or human review process is triggered.

4. **Prompt Injection & Malicious Input Protection**
   - Treat LLM as untrusted output; avoid concatenating raw user inputs into system-level context. Use clearly delimited sections and system messages.
   - Include instructions in system prompts to ignore or not execute embedded instructions from user content.
   - Acceptance criteria: Unit tests simulate prompt injection attempts and the system resists them.

5. **Rate Limit & Abuse Detection**
   - Monitor for unusual request patterns that may indicate abuse (high frequency from a single session/fingerprint or repeating similar requests).
   - Auto-throttle or flag suspicious sessions for review and possible blocking.
   - Acceptance criteria: Abusive traffic patterns are detected and mitigated without affecting normal traffic.

6. **Human-in-the-loop & Escalation Paths**
   - Provide an easy path for manual review of flagged outputs (researcher dashboard to approve/override outputs) and a mechanism to re-run or correct outputs.
   - Acceptance criteria: Flagged outputs are surfaced to reviewers with links to raw prompts and responses.

7. **Auditing & Logging (Privacy-aware)**
   - Log prompts, responses (or hashes of them), and key-metadata (timestamps, key used, latencies, error codes) for debugging and safety audits.
   - Avoid storing raw user PII in logs; redact sensitive fields and keep logs access-controlled.
   - Acceptance criteria: Sufficient data for post-mortem without leaking PII; logs are retained for an appropriate retention policy.

8. **Model Selection & Conservative Defaults**
   - Prefer models/providers with safety features. Use smaller or more conservative models for feedback if that suffices, and reserve stronger models only where needed.
   - Tune temperature and max tokens to minimize hallucination and unexpected behaviors for structured outputs.
   - Acceptance criteria: Lower hallucination rates and higher schema conformance for structured tasks.

9. **Fallback & Graceful Degradation**
   - When an LLM response fails safety checks or the model is unavailable, return a deterministic fallback (e.g., generic feedback message or cached response) and surface an explanatory note to the user.
   - Acceptance criteria: Users never see unvalidated LLM content; failures are handled gracefully.

10. **Testing & Red Teaming**
    - Create unit/integration tests simulating adversarial inputs (toxic, PII, injection) and run periodic red-team exercises on prompts and outputs.
    - Acceptance criteria: Passing score on red-team test suite before releasing new prompt changes.

11. **Privacy Controls & Data Minimization**
    - Where possible, anonymize or aggregate data before sending it to LLM providers.
    - Consider differential privacy or local preprocessing for sensitive analytics.
    - Acceptance criteria: No identifiable PII is sent to third-party services; documented privacy guarantees for feedback generation.

12. **Operational Safety (Monitoring & Alerts)**
    - Add alerts for sudden spikes in unsafe outputs or error rates per provider/key and dashboards tracking safety metrics.
    - Acceptance criteria: On-call/team receives alerts for anomalous behavior within defined SLAs.

---

If you'd like, I can:
- Add a `system design/diagrams/llm-guardrails.puml` illustrating the guardrail flow (prompt templating, safety checks, cache, human review), or
- Add a `system design/llm-safety-checklist.md` with a checklist for PR review and release for LLM-related changes.

Which would you prefer I do next?