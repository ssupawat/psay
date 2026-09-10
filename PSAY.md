# Audio Announcement Protocol

Announce key moments out loud with `psay` — a local-TTS command (piper + afplay, no cloud). Requires the `psay` binary in `PATH` (macOS).

```bash
psay '[Action/Discovery] on [Target] because [Context]. [My Take / Advice]'
```

## Examples

* `psay 'Found critical deadlock in auth_service.go due to mutex contention. Hold off on manual testing while I refactor it to atomic pointers.'`
* `psay 'Pivoting payment service to mock gateway because Stripe API returned rate limit errors. No action needed on your side, I will handle the retry logic.'`
* `psay 'PostgreSQL cluster migration failed on edge nodes due to memory limits. Approve increasing the node specs before I re-run the job.'`

## Rules

* **Short by design:** 1–2 sentences (~20 words). You condense before calling — psay does not rewrite.
* **Include context & take:** Always specify **Action/Discovery**, **Target**, **Reason**, and **Your Take** (what you think or what to do next).
* **No raw syntax:** Do not pass raw code blocks, stack traces, or massive terminal logs directly.
* **When:** Announce starting tasks, pivots, discoveries, and recommendations. Never announce routine completion — report results silently.
