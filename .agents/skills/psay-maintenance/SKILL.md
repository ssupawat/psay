---
name: psay-maintenance
description: Build, test, release, and evolve psay — the kokoro-TTS announce tool in this repo. Use whenever working in the psay repository — changing code, lexicon, or voice config, cutting releases, debugging what psay speaks, rewriting history, or weighing design changes — even if the user just says "psay".
---

# psay maintenance

psay turns developer jargon into speakable text and plays it through local TTS (kokoro + afplay). It is deliberately a thin wrapper: config + lexicon + exec into a Python venv + afplay. Preserve that shape; new features need evidence, not speculation.

## Architecture

- `main.go` — CLI (`-v` voice override), temp-wav plumbing, kokoro (via `~/.psay/venv` python + embedded `speak.py`) → afplay chain
- `speak.py` — embedded (`go:embed`) 7-line kokoro-onnx script; changing it requires a new release
- `settings.go` — `~/.psay/settings.json` load/auto-init (voice, `speed`, lexicon); migrates old piper files (`length_scale` key) on load
- `lexicon.go` — single-pass, word-boundary lexicon replacer
- `init.go` — installs the embedded PSAY.md protocol into harness instruction files (appends to existing files only, marker-wrapped, updates stale blocks in place)
- `PSAY.md` — agent-facing protocol, embedded via `go:embed`
- `install.sh` / `uninstall.sh` — one-shot user setup (binary + kokoro venv + model) and its reverse

## Verify loop

Run before every push:

```bash
gofmt -l . && go vet ./... && go test ./... && go build -o psay .
```

## Releases

- A tag push triggers goreleaser → GitHub Release with darwin amd64+arm64 tarballs: `git tag vX.Y.Z && git push origin vX.Y.Z`
- Tag whenever the binary changes — including PSAY.md or speak.py edits, because both are `go:embed`ded and the released binary serves them. Docs-only changes don't need a tag.
- Then reinstall locally and verify: `curl -fsSL https://raw.githubusercontent.com/ssupawat/psay/main/install.sh | sh`

## Never force-push tags

Force-pushing a tag re-fires the release workflow; goreleaser then fails with `422 already_exists` on every asset. The releases survive, but the Actions history shows red duplicates. If history must be rewritten (e.g. removing a word from commit messages), rewrite messages only via `git filter-branch --msg-filter` with `--tag-name-filter cat`, force-push main AND tags together, then delete the failed duplicate runs (`gh run delete`). Better: never put anything in a commit message you wouldn't want public forever.

## Design decisions — do not undo without evidence

- **No phonetic engine.** Listening tests proved the G2P front-end (espeak lineage in both piper and kokoro) already splits camelCase and reads symbols, extensions, and versions: `getUserById` → "get user by I-D", `auth_service.go` → "auth service. go". Don't build preprocessing the engine already does.
- **No LLM summarization.** The calling agent condenses its own announce; PSAY.md enforces brevity (~15 words).
- **Lexicon-first fixing.** A misread word is a lexicon entry in `~/.psay/settings.json`, not code. Change code only when a whole token class proves broken.
- **Word-boundary replacement.** Naive `strings.ReplaceAll` turns "the database is slow" into "the datadatabasease is slow" when mapping "db". Never simplify `lexicon.go` back to it.

## Diagnosing pronunciation

Phonemes are ground truth — kokoro's pipeline is text → G2P tokenizer → ONNX model, nothing else. Dump what will be spoken:

```bash
~/.psay/venv/bin/python - <<'PYEOF'
from kokoro_onnx import Kokoro
k = Kokoro("/Users/art/.psay/kokoro/kokoro-v1.0.onnx", "/Users/art/.psay/kokoro/voices-v1.0.bin")
for t in ["auth_service.go", "getUserById"]:
    print(repr(t), "->", k.tokenizer.phonemize(t, lang="en-us"))
PYEOF
```

Then confirm by ear: synthesize a wav and `afplay` it. Correct phonemes + wrong audio = model rendering (lexicon can't fix). Wrong phonemes = G2P level (lexicon entry or upstream fix). Don't trust your decoding of IPA over the user's ears — play the wav.

## Kokoro gotchas

- No CLI exists — psay execs `~/.psay/venv/bin/python -c <embedded speak.py>`; the venv and `~/.psay/kokoro/{kokoro-v1.0.onnx,voices-v1.0.bin}` come from install.sh
- Model download (~330MB): `https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files-v1.1/…`
- `speed` in settings.json is kokoro's pace: higher = FASTER (opposite of piper's old length_scale)
- Voices: `af_heart` is the only grade-A default; list at hexgrad/Kokoro-82M VOICES.md on Hugging Face
- Known misreads (identical to piper, lexicon material): `aws` → "awz", `db` → "dee bee", `psay` → "say"
- The model loads on every invocation (~1s) — fine for short announces, wrong for streaming
- `afplay` is macOS-only — that's why releases are darwin-only
