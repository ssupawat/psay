# Audio Announcement Protocol

Announce key moments out loud with `psay` (macOS, must be in `PATH`):

```bash
psay '[Action/Discovery] on [Target] because [Context]. [My Take / Advice]'
```

Example: `psay 'Found deadlock in auth_service.go, holding manual tests while I refactor to atomic pointers.'`

* 1–2 sentences (~20 words) — you condense before calling.
* No code blocks, stack traces, or logs.
* Starts, pivots, discoveries, recommendations only — never routine completion; report results silently.
