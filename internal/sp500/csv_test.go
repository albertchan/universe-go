package sp500

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteCSV(t *testing.T) {
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

	path := filepath.Join(t.TempDir(), "sp500.csv")
	if err := WriteCSV(path, companies); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open csv: %v", err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("read csv: %v", err)
	}

	if got, want := len(rows), 3; got != want {
		t.Fatalf("rows = %d, want %d (header + 2 data rows)", got, want)
	}

	wantHeader := []string{
		"symbol", "name", "gics_sector", "gics_sub_industry",
		"location", "cik", "added_at", "founded_at",
	}
	if !reflect.DeepEqual(rows[0], wantHeader) {
		t.Errorf("header = %v, want %v", rows[0], wantHeader)
	}

	wantRow1 := []string{
		"MMM", "3M", "Industrials", "Industrial Conglomerates",
		"Saint Paul, Minnesota", "0000066740", "1957-03-04", "1902",
	}
	if !reflect.DeepEqual(rows[1], wantRow1) {
		t.Errorf("row 1 = %v, want %v", rows[1], wantRow1)
	}
}

func TestWriteCSVEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.csv")
	if err := WriteCSV(path, nil); err != nil {
		t.Fatalf("WriteCSV with no companies: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open csv: %v", err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("read csv: %v", err)
	}
	if got, want := len(rows), 1; got != want {
		t.Fatalf("rows = %d, want %d (header only)", got, want)
	}
}

func TestWriteCSVBadPath(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "does-not-exist", "out.csv")
	if err := WriteCSV(bad, nil); err == nil {
		t.Fatal("expected error writing to nonexistent directory")
	}
}
