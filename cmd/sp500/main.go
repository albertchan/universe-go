package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bishop-bot/universe-go/internal/sp500"
)

func main() {
	root := &cobra.Command{
		Use:   "sp500",
		Short: "Scrape the current S&P 500 list from Wikipedia",
		Long: "sp500 scrapes the current S&P 500 constituents table from Wikipedia and " +
			"writes the result to a CSV file or a local SQLite database.",
		SilenceUsage: true,
	}

	root.AddCommand(newCSVCmd())
	root.AddCommand(newSQLiteCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func newCSVCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "csv",
		Short: "Scrape and write the S&P 500 list to a CSV file",
		RunE: func(_ *cobra.Command, _ []string) error {
			companies, err := sp500.Fetch()
			if err != nil {
				return err
			}
			if err := sp500.WriteCSV(out, companies); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "wrote %d companies to %s\n", len(companies), out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&out, "out", "o", "sp500.csv", "output CSV file path")
	return cmd
}

func newSQLiteCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "sqlite",
		Short: "Scrape and store the S&P 500 list in a local SQLite database",
		RunE: func(_ *cobra.Command, _ []string) error {
			companies, err := sp500.Fetch()
			if err != nil {
				return err
			}
			if err := sp500.WriteSQLite(out, companies); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "upserted %d companies into %s\n", len(companies), out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&out, "out", "o", "sp500.db", "output SQLite database path")
	return cmd
}
