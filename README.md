# universe-go

This project is a Go script that scrapes the current S&P 500 list from Wikipedia.

For each constituent it stores: `symbol`, `name`, `gics_sector`, `gics_sub_industry`, `location`, `cik`, `added_at`, `founded_at`.

## Quick Start

Requires Go 1.23+.

Install dependencies:

```sh
go mod tidy
```

Write the list to a CSV file:

```sh
go run ./cmd/sp500 csv -o sp500.csv
```

Write the list to a local SQLite database (table `sp500_companies`, upsert on `symbol`):

```sh
go run ./cmd/sp500 sqlite -o sp500.db
```

The `-o` / `--out` flag is optional; defaults are `sp500.csv` and `sp500.db` respectively.

Build a standalone binary:

```sh
go build -o sp500 ./cmd/sp500
./sp500 csv -o sp500.csv
```

Run the tests:

```sh
go test ./...
```

Run tests with verbose output:

```sh
go -C /Users/albert/Developer/github.com/bishop-bot/universe-go test -v ./internal/sp500 2>&1 | tail -50
```

## Sources

- [List of S&P 500 companies](https://en.wikipedia.org/wiki/List_of_S%26P_500_companies)