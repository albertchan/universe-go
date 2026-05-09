package sp500

import (
	"encoding/csv"
	"fmt"
	"os"
)

func WriteCSV(path string, companies []Company) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create csv %s: %w", path, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	header := []string{
		"symbol", "name", "gics_sector", "gics_sub_industry",
		"location", "cik", "added_at", "founded_at",
	}
	if err := w.Write(header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, c := range companies {
		row := []string{
			c.Symbol, c.Name, c.GICSSector, c.GICSSubIndustry,
			c.Location, c.CIK, c.AddedAt, c.FoundedAt,
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("write row %s: %w", c.Symbol, err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("flush csv: %w", err)
	}
	return nil
}
