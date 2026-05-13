---
name: go-test-coverage
description: Validate and expand unit-test coverage for Go packages in this repo. Use when the user asks to "validate and add testcases for <package>", "increase coverage", "add unit tests", or similar. Defines the full workflow: build check → coverage audit → table-driven tests → doc comments → bug-flagging protocol.
---

# Go Test Coverage Workflow (go-lib)

Apply this workflow whenever the user asks to add / expand tests for a package in this repo.

## 1. Validate the package builds first

```bash
go build ./<pkg>/
go test -cover ./<pkg>/
```

If the build fails:

- Check `git status <pkg>/` — if the broken files are **untracked**, they are local scratch work. Do NOT delete them. Surface the problem to the user with the exact import errors and ask which route to take (fix imports, remove files, skip the package).
- If the broken files are tracked and the fix is obvious + minimal (e.g. swap a private import for `comparable`, drop a dead import), apply it — but state clearly that you're modifying source, not just tests.

Never silently rewrite source while the user expects "just add tests."

## 2. Map uncovered lines before writing anything

```bash
go test -coverprofile=/tmp/<pkg>_cover.out ./<pkg>/
go tool cover -func=/tmp/<pkg>_cover.out        # per-function %
go tool cover -html=/tmp/<pkg>_cover.out -o /tmp/<pkg>_cover.html
grep -B 1 -A 3 'cov0" title="0"' /tmp/<pkg>_cover.html | head -60
```

For each uncovered line decide:

- **Reachable via the public API** → write a test that hits it.
- **Reachable only by calling an unexported helper** (because JSON unmarshal can't produce the required dVal type, etc.) → write a same-package test that calls the helper directly with a synthetic input. This is legitimate — e.g. `TestProcessStruct_DirectDValTypes` in `generators/default.tags_test.go`.
- **Unreachable / defensive dead code** (post-bounds-check `IsValid` guard, regex-guarded `strconv.Atoi` error, library `New()` that can't fail with the passed options) → leave uncovered. In the PR / summary, list which lines are intentionally skipped and why.

## 3. Prefer table-driven tests

This repo's existing style is a slice of structs + one loop with `t.Run`. Match it. Don't mix subtest blocks and table cases in the same test function.

```go
tests := []struct {
    name    string
    source  any
    path    string
    value   any
    wantErr error
    wantX   any // nil → skip this check
}{
    {name: "empty path is no-op", source: map[string]any{"a": 1}, path: "", wantX: map[string]any{"a": 1}},
    // ...
}

for id, test := range tests {
    t.Run(fmt.Sprintf("%v_%s", id, test.name), func(t *testing.T) {
        // ...
    })
}
```

For functions that return `(value, error)`, this repo uses a `wantT` struct in `algo/extract_test.go` — reuse the same pattern inside a package rather than inventing a new one.

## 4. Comment discipline for tests

- **One short doc comment above each test function** explaining what behavior / branches it covers. Not a restatement of the name.
- **One inline comment per table case** when the case's purpose isn't obvious from its field values (e.g. "missing top-level key is a silent no-op", "duration-style string goes through ParseDuration").
- No trailing "// Test case 1" style comments. No emoji.

Look at `datatype/inspect_test.go` and `algo/extract_test.go` for the tone to match.

## 5. Bug-flagging protocol

If a test you add exposes a latent bug in the source:

- Do **not** silently fix the source unless the user asked for a fix.
- Pin the current (broken) behavior with `assert.Panics` / `assert.Error` and a test named `TestX_YBug` — see `TestGenerateDefaults_UintPanicsBug` in `generators/default.tags_test.go` for the template.
- Leave a comment on the pinned test saying: which source line is wrong, what the fix would be, and "delete this test once fixed."
- In your summary to the user, call out the bug explicitly and offer to patch it as a separate step.

## 6. Coverage target

- Aim for ≥95% per-function.
- 100% is only worth chasing if it doesn't require contrived tests. Clearly-defensive branches (sqids `New()` failing, `reflect.Value.IsValid()` after a length guard) can stay uncovered.
- In the summary: report `before% → after%` and list any functions still below 100% with a one-line reason.

## 7. Output format for the user

End with a terse summary:

- **Fix applied** (if any source was changed, with file:line)
- **Tests added** (file names + what they cover)
- **Coverage** (before → after, and residual gaps)
- **Bugs found** (if any, with offer to patch)

No emojis. No multi-paragraph explanations of what the tests do — the test comments already explain that.

## Useful conventions already in this repo

- Same-package tests live in `<file>_test.go` alongside the source.
- `github.com/stretchr/testify/assert` is the assertion library.
- Test names: `TestFuncName` for the main table, `TestFuncName_Scenario` for targeted single-case tests (e.g. `TestDecodeHash_Empty`, `TestGenerateDefaults_UintPanicsBug`).
- `wantT` struct pattern for `(value, error)` returns — reuse when present.
