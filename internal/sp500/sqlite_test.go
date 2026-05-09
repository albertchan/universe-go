package sp500

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestWriteSQLite(t *testing.T) {
	companies := []Company{
		{
			Symbol:          "MMM",
			Name:            "3M",
			GICSSector:      "Industrials",
			GICSSubIndustry: "Industrial Conglomerates",
			Location:        "Saint Paul, Minnesota",
			CIK:             "0000066740",
			AddedAt:         "1957-03-04",
			FoundedAt:       "1902",
		},
		{
			Symbol:          "AOS",
			Name:            "A. O. Smith",
			GICSSector:      "Industrials",
			GICSSubIndustry: "Building Products",
			Location:        "Milwaukee, Wisconsin",
			CIK:             "0000091142",
			AddedAt:         "2017-07-26",
			FoundedAt:       "1916",
		},
	}

	path := filepath.Join(t.TempDir(), "sp500.db")
	if err := WriteSQLite(path, companies); err != nil {
		t.Fatalf("WriteSQLite: %v", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM sp500_companies").Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Errorf("row count = %d, want 2", count)
	}

	var got Company
	err = db.QueryRow(`
		SELECT symbol, name, gics_sector, gics_sub_industry, location, cik, added_at, founded_at
		FROM sp500_companies WHERE symbol = ?`, "MMM").Scan(
		&got.Symbol, &got.Name, &got.GICSSector, &got.GICSSubIndustry,
		&got.Location, &got.CIK, &got.AddedAt, &got.FoundedAt,
	)
	if err != nil {
		t.Fatalf("query MMM: %v", err)
	}
	if got != companies[0] {
		t.Errorf("MMM row = %#v, want %#v", got, companies[0])
	}
}

func TestWriteSQLiteUpsert(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sp500.db")

	original := []Company{{
		Symbol: "MMM", Name: "3M", GICSSector: "Industrials",
		GICSSubIndustry: "Industrial Conglomerates",
		Location:        "Saint Paul, Minnesota",
		CIK:             "0000066740", AddedAt: "1957-03-04", FoundedAt: "1902",
	}}
	if err := WriteSQLite(path, original); err != nil {
		t.Fatalf("first write: %v", err)
	}

	updated := []Company{{
		Symbol: "MMM", Name: "3M Company", GICSSector: "Industrials",
		GICSSubIndustry: "Industrial Conglomerates",
		Location:        "St. Paul, Minnesota",
		CIK:             "0000066740", AddedAt: "1957-03-04", FoundedAt: "1902",
	}}
	if err := WriteSQLite(path, updated); err != nil {
		t.Fatalf("second write: %v", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM sp500_companies").Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("row count = %d, want 1 (upsert must not duplicate)", count)
	}

	var name, location string
	if err := db.QueryRow("SELECT name, location FROM sp500_companies WHERE symbol = ?", "MMM").Scan(&name, &location); err != nil {
		t.Fatalf("query: %v", err)
	}
	if name != "3M Company" || location != "St. Paul, Minnesota" {
		t.Errorf("after upsert: name=%q location=%q, want updated values", name, location)
	}
}

func TestWriteSQLiteSchemaIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sp500.db")
	if err := WriteSQLite(path, nil); err != nil {
		t.Fatalf("first write (empty): %v", err)
	}
	if err := WriteSQLite(path, nil); err != nil {
		t.Fatalf("second write (empty) on existing schema: %v", err)
	}
}
