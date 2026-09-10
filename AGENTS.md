# AGENTS.md

## Build

- `go test ./...` — run tests before pushing
- `go build -o psay .` — build
- Zero Go dependencies — stdlib only
- Release: push a tag (`git tag vX.Y.Z && git push origin vX.Y.Z`), CI builds it via goreleaser
- Never force-push tags — re-fires the release workflow against existing releases (422 already_exists failures)

@PSAY.md
