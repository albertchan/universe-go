package sp500

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	wikipediaURL = "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies"
	userAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 14.7; rv:135.0) Gecko/20100101 Firefox/135.0"
)

type Company struct {
	Symbol          string
	Name            string
	GICSSector      string
	GICSSubIndustry string
	Location        string
	CIK             string
	AddedAt         string
	FoundedAt       string
}

func Fetch() ([]Company, error) {
	req, err := http.NewRequest(http.MethodGet, wikipediaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch wikipedia: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected wikipedia status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	return parseConstituents(doc)
}

func parseConstituents(doc *goquery.Document) ([]Company, error) {
	table := doc.Find("table#constituents").First()
	if table.Length() == 0 {
		return nil, fmt.Errorf("constituents table (id=constituents) not found")
	}

	var companies []Company
	table.Find("tbody > tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 8 {
			return
		}
		c := Company{
			Symbol:          clean(cells.Eq(0).Text()),
			Name:            clean(cells.Eq(1).Text()),
			GICSSector:      clean(cells.Eq(2).Text()),
			GICSSubIndustry: clean(cells.Eq(3).Text()),
			Location:        clean(cells.Eq(4).Text()),
			AddedAt:         clean(cells.Eq(5).Text()),
			CIK:             clean(cells.Eq(6).Text()),
			FoundedAt:       clean(cells.Eq(7).Text()),
		}
		if c.Symbol == "" {
			return
		}
		companies = append(companies, c)
	})

	if len(companies) == 0 {
		return nil, fmt.Errorf("no rows parsed from constituents table")
	}
	return companies, nil
}

var citationRE = regexp.MustCompile(`\[[^\]]*\]`)

func clean(s string) string {
	s = citationRE.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, " ", " ")
	return strings.TrimSpace(s)
}
