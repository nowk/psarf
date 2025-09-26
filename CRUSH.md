# CRUSH.md - Guidelines for coding agents

## Build, Lint & Test Commands

- Build the project: `go build ./...`
- Run all tests: `go test ./...`
- Run a single test function: `go test -run ^TestFunctionName$ ./...`
- Run tests with verbose output: `go test -v ./...`
- Format code (go fmt): `go fmt ./...`

## Environment & Dependencies

- Go version is set in go.mod (1.15) - maintain compatibility
- Environment variables managed with asdf version loader via `.envrc` file

## Code Style Guidelines

- Use standard Go idioms for imports and error handling
- Group imports as standard library and third-party separately
- Use explicit error checks and return errors upwards
- Use clear and descriptive variable names with camelCase (e.g. psarSeries, pipOffset)
- Use iota for constant enums (e.g. Direction type)
- Use pointers when modifying state or avoiding copies (e.g. *time.Time, *PsarPeriod)
- Keep code idiomatic Go style formatting (use `go fmt` regularly)
- Comment exported types and functions with capitalized TODO sentences
- Use initialisms properly capitalized (e.g. EP, AF, Sar in structs)

## Testing

- Tests are in *_test.go files in the same package
- Use table-driven tests where applicable
- Use time.Now() or fixed timestamps for test date values
- Run tests with `go test -v` for detailed logging

## Misc Notes

- This project uses daily bars with truncated time operations for date comparison
- The main logic is in iterative style structure allowing stepwise Psar calculation
- Project uses third-party packages such as `github.com/piquette/finance-go` for charting (in cmd)

## Cursor & Copilot Rules

- No .cursor/rules or copilot instructions files were found in the repository

## Adding Go Language Server Protocol (LSP)

- Install the official Go language server `gopls` with:
  ```
  go install golang.org/x/tools/gopls@latest
  ```
- Ensure your PATH includes the directory where `gopls` is installed (e.g., `$HOME/go/bin`)
- Configure your editor or IDE to use `gopls` as the Go language server
  - VS Code: Install the Go extension by Microsoft (auto-uses `gopls`)
  - Vim/Neovim: Use plugins like `coc.nvim` or `nvim-lspconfig` to connect `gopls`
  - Other editors: Follow specific instructions for Go LSP integration
- Restart your editor after installation to enable features like auto-completion, go-to definition, and inline error checking

## Summary

Follow standard Go development best practices.
Use provided test files as reference for idiomatic usage and test coverage.
Maintain cleanliness and clarity in code and commit messages.

💘 Generated with Crush
Co-Authored-By: Crush <crush@charm.land>
