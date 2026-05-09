package sp500

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const fixtureHTML = `
<html><body>
<table id="constituents">
  <thead>
    <tr><th>Symbol</th><th>Security</th><th>GICS Sector</th><th>GICS Sub-Industry</th>
        <th>Headquarters Location</th><th>Date added</th><th>CIK</th><th>Founded</th></tr>
  </thead>
  <tbody>
    <tr>
      <td><a href="/wiki/3M">MMM</a></td>
      <td><a href="/wiki/3M">3M</a></td>
      <td>Industrials</td>
      <td>Industrial Conglomerates</td>
      <td>Saint Paul, Minnesota</td>
      <td>1957-03-04</td>
      <td>0000066740</td>
      <td>1902</td>
    </tr>
    <tr>
      <td>AOS<sup class="reference">[1]</sup></td>
      <td>A. O. Smith</td>
      <td>Industrials</td>
      <td>Building Products</td>
      <td>Milwaukee, Wisconsin</td>
      <td>2017-07-26</td>
      <td>0000091142</td>
      <td>1916</td>
    </tr>
    <tr>
      <td></td><td></td><td></td><td></td>
      <td></td><td></td><td></td><td></td>
    </tr>
    <tr>
      <td>SHORT</td><td>too few cells</td>
    </tr>
  </tbody>
</table>
</body></html>
`

func TestParseConstituents(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fixtureHTML))
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}

	companies, err := parseConstituents(doc)
	if err != nil {
		t.Fatalf("parseConstituents: %v", err)
	}

	if got, want := len(companies), 2; got != want {
		t.Fatalf("companies count = %d, want %d (empty + short rows must be skipped)", got, want)
	}

	want := Company{
		Symbol:          "MMM",
		Name:            "3M",
		GICSSector:      "Industrials",
		GICSSubIndustry: "Industrial Conglomerates",
		Location:        "Saint Paul, Minnesota",
		CIK:             "0000066740",
		AddedAt:         "1957-03-04",
		FoundedAt:       "1902",
	}
	if companies[0] != want {
		t.Errorf("row 0 = %#v, want %#v", companies[0], want)
	}

	if companies[1].Symbol != "AOS" {
		t.Errorf("row 1 symbol = %q, want %q (citation footnote should be stripped)", companies[1].Symbol, "AOS")
	}
	if companies[1].CIK != "0000091142" {
		t.Errorf("row 1 CIK = %q, want %q (leading zeros must be preserved)", companies[1].CIK, "0000091142")
	}
}

func TestParseConstituentsMissingTable(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<html><body><p>no table</p></body></html>"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, err := parseConstituents(doc); err == nil {
		t.Fatal("expected error when constituents table is missing")
	}
}

func TestParseConstituentsEmptyTable(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(
		`<table id="constituents"><tbody></tbody></table>`,
	))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, err := parseConstituents(doc); err == nil {
		t.Fatal("expected error when no rows are parsed")
	}
}

func TestClean(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"  hello  ", "hello"},
		{"AOS[1]", "AOS"},
		{"foo[note 2]bar", "foobar"},
		{"line break", "line break"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := clean(tc.in); got != tc.want {
			t.Errorf("clean(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
