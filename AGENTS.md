# wiki

`wiki` is a Go CLI that searches and reads English Wikipedia articles in the terminal, with Markdown output through Glamour. It ships through the Homebrew tap `aayush9029/tap`. Status: shipped.

## Commands

| Task | Command |
|------|---------|
| Build | `go build -o wiki .` |
| Search | `go run . search "quantum computing"` |
| Read | `go run . read "Quantum computing"` |
| Release | Load the `oss-cli-release` skill |

## Architecture

- `cmd/`: Cobra commands `search`, `read`, `random`. `cmd.Version` holds the version.
- `wikipedia/`: client for the MediaWiki API (`https://en.wikipedia.org/w/api.php`) and the REST API (`https://en.wikipedia.org/api/rest_v1`).

## Rules

- `read` shows the full article by default. `read -s` shows only the summary.
- The version is `var Version` in `cmd/root.go`, not `main.version`. The tap formula `../../homebrew-tap/Formula/wiki.rb` sets it with `-X github.com/Aayush9029/wiki/cmd.Version=#{version}`. If you move the variable, update the formula.
- The built `wiki` binary at the root is in `.gitignore`. Do not commit it.
