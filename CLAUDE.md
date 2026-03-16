# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an unofficial mirror of the SQLite sqllogictest suite, maintained by DoltHub. It includes a Go-based parser and runner for executing SQL correctness tests against pluggable database engines. The test suite contains 5.9 million SQL queries with expected results regenerated for MySQL 8.x.

## Build & Test Commands

All Go code lives under `go/`. The module path is `github.com/dolthub/sqllogictest/go`.

```bash
# Build
cd go && go build ./...

# Run unit tests
cd go && go test ./logictest/parser/...

# Run the MySQL test runner (requires local MySQL setup)
cd go && go run logictest/mysql/main/main.go verify ../../test/evidence/in1.test

# Modes: verify, generate, filter, analyze
cd go && go run logictest/mysql/main/main.go analyze ../../test/evidence/
```

### MySQL Setup (required for running tests against MySQL)

```sql
CREATE DATABASE sqllogictest;
CREATE USER sqllogictest@localhost IDENTIFIED BY "password";
GRANT ALL ON sqllogictest.* TO sqllogictest@localhost;
```

## Architecture

The Go code follows a clean separation of concerns:

- **Parser** (`go/logictest/parser/`): Parses `.test` files into `Record` structs. Each record is either a `Statement` (DDL/DML expecting ok/error), a `Query` (SELECT with expected results), or `Halt`. Records support engine-specific conditions (`onlyif`/`skipif`).

- **Harness interface** (`go/logictest/harness.go`): Defines the `Harness` interface that database adapters must implement: `Init()`, `ExecuteStatement()`, `ExecuteQuery()`, `EngineStr()`, `GetTimeout()`. Column types use single-char schema strings: `I` (integer), `R` (real), `T` (text).

- **Runner** (`go/logictest/runner.go`): Orchestrates test execution. Takes a harness and test file paths, parses them, runs each record, and logs results. Supports verify (check results), generate (produce expected output), and filter (generate excluding failures) modes.

- **MySQL harness** (`go/logictest/mysql/`): Reference `Harness` implementation for MySQL. The `main/` subdirectory provides the CLI entry point.

- **Result parser** (`go/logictest/resultparser.go`): Parses runner log output back into structured results (Ok, NotOk, Skipped, Timeout, DidNotRun).

## Test File Format

```
statement ok
CREATE TABLE t1(a INTEGER, b INTEGER)

query II nosort
SELECT a, b FROM t1
----
1
2
```

- `statement ok|error` — DDL/DML that should succeed or fail
- `query [SCHEMA] [nosort|rowsort|valuesort] [label]` — SELECT with expected results after `----`
- Results can be inline values or MD5 hashes (for large result sets)
- `onlyif engine` / `skipif engine` — conditional execution

## Key Environment Variables

- `SQLLOGICTEST_TRUNCATE_QUERIES` — controls query truncation in runner output
