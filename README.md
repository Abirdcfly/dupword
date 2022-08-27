# dupword

A linter that checks for duplicate words in the source code (usually miswritten)

## example

1. Repeated words appear on two adjacent lines [commit](https://github.com/golang/go/commit/d8f90ce0f8119bf593efb6fb91825de5b61fcda7)

```diff
--- a/src/cmd/compile/internal/ssa/schedule.go
+++ b/src/cmd/compile/internal/ssa/schedule.go
@@ -179,7 +179,7 @@ func schedule(f *Func) {
 					// scored CarryChainTail (and prove w is not a tail).
 					score[w.ID] = ScoreFlags
 				}
-				// Verify v has not been scored. If v has not been visited, v may be the
+				// Verify v has not been scored. If v has not been visited, v may be
 				// the final (tail) operation in a carry chain. If v is not, v will be
 				// rescored above when v's carry-using op is scored. When scoring is done,
 				// only tail operations will retain the CarryChainTail score.
```

2. Repeated words appearing on the same line [commit](https://github.com/golang/go/commit/48da729e8468b630ee003ac51cbaac595d53bec8)

```diff
--- a/src/net/http/cookiejar/jar.go
+++ b/src/net/http/cookiejar/jar.go
@@ -465,7 +465,7 @@ func (j *Jar) domainAndType(host, domain string) (string, bool, error) {
 		// dot in the domain-attribute before processing the cookie.
 		//
 		// Most browsers don't do that for IP addresses, only curl
-		// version 7.54) and and IE (version 11) do not reject a
+		// version 7.54) and IE (version 11) do not reject a
 		//     Set-Cookie: a=1; domain=.127.0.0.1
 		// This leading dot is optional and serves only as hint for
 		// humans to indicate that a cookie with "domain=.bbc.co.uk"
```

## Install

```bash
go install github.com/Abirdcfly/dupword/cmd/dupword@latest
```

**Or** install the main branch (including the last commit) with:

```bash
go install github.com/Abirdcfly/dupword/cmd/dupword@main
```

## Usage

Run with default settings(include test file, only check keywords `the`,`a` and `and`):

```bash
dupword ./...
```

Skip detection test file(`*_test.go`):

```bash
dupword -test=false ./...
```

All options:

```bash
$ dupword --help
dupword: checks for duplicate words in the source code (usually miswritten)

Usage: dupword [-flag] [package]

This analyzer checks miswritten duplicate words in comments or package doc or string declaration

Flags:
  -V    print version and exit
  -all
        no effect (deprecated)
  -c int
        display offending line with this many lines of context (default -1)
  -cpuprofile string
        write CPU profile to this file
  -debug string
        debug flags, any subset of "fpstv"
  -fix
        apply all suggested fixes
  -flags
        print analyzer flags in JSON
  -json
        emit JSON output
  -memprofile string
        write memory profile to this file
  -source
        no effect (deprecated)
  -tags string
        no effect (deprecated)
  -test
        indicates whether test files should be analyzed, too (default true)
  -trace string
        write trace log to this file
  -v    no effect (deprecated)
```

## TODO

- [ ] make keyword optional
- [ ] add this linter to golangci-lint
- [ ] rewrite the detection logic to make it more efficient

## License

MIT
