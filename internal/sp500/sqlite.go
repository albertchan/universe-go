package sp500

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const sqliteSchema = `
CREATE TABLE IF NOT EXISTS sp500_companies (
    symbol            TEXT PRIMARY KEY,
    name              TEXT NOT NULL,
    gics_sector       TEXT,
    gics_sub_industry TEXT,
    location          TEXT,
    cik               TEXT,
    added_at          TEXT,
    founded_at        TEXT
);
`

const upsertSQL = `
INSERT INTO sp500_companies (
    symbol, name, gics_sector, gics_sub_industry,
    location, cik, added_at, founded_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(symbol) DO UPDATE SET
    name              = excluded.name,
    gics_sector       = excluded.gics_sector,
    gics_sub_industry = excluded.gics_sub_industry,
    location          = excluded.location,
    cik               = excluded.cik,
    added_at          = excluded.added_at,
    founded_at        = excluded.founded_at
`

func WriteSQLite(path string, companies []Company) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("open sqlite %s: %w", path, err)
	}
	defer db.Close()

	if _, err := db.Exec(sqliteSchema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(upsertSQL)
	if err != nil {
		return fmt.Errorf("prepare upsert: %w", err)
	}
	defer stmt.Close()

	for _, c := range companies {
		if _, err := stmt.Exec(
			c.Symbol, c.Name, c.GICSSector, c.GICSSubIndustry,
			c.Location, c.CIK, c.AddedAt, c.FoundedAt,
		); err != nil {
			return fmt.Errorf("upsert %s: %w", c.Symbol, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
