package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

type Repo struct {
	FullName    string   `json:"full_name"`
	Description string   `json:"description"`
	Stars       int      `json:"stargazers_count"`
	Language    string   `json:"language"`
	HTMLURL     string   `json:"html_url"`
	CreatedAt   string   `json:"created_at"`
	Fork        bool     `json:"fork"`
	Archived    bool     `json:"archived"`
	Topics      []string `json:"topics"`
}

type Response struct {
	Items []Repo `json:"items"`
}

var (
	duration    = pflag.StringP("duration", "d", "week", "Time range: day, week, month, year")
	limit       = pflag.IntP("limit", "l", 10, "Number of repositories (1-100)")
	language    = pflag.String("language", "", "Filter by programming language")
	langAlias   = pflag.String("lang", "", "Alias for --language")
	jsonOutput  = pflag.Bool("json", false, "Output as JSON")
	saveFile    = pflag.String("save", "", "Save output to file")
	openBrowser = pflag.Bool("open", false, "Open first repository")
	spoken      = pflag.String("spoken", "", "Filter by spoken language")
	proxy       = pflag.String("proxy", "", "HTTP/HTTPs proxy")

	noColor = pflag.Bool("no-color", false, "Disable Colors")
	watch   = pflag.DurationP("watch", "w", 0, "Auto refresh (e.g. 5m)")
	today   = pflag.Bool("today", false, "Shortcut for --duration day")
	weekly  = pflag.Bool("week", false, "Shortcut for --duration week")
	monthly = pflag.Bool("monthly", false, "Shortcut for --duration month")
)

func main() {
	pflag.Parse()

	if *noColor {
		color.NoColor = true
	}

	if *today {
		*duration = "day"
	}

	if *weekly {
		*duration = "week"
	}

	if *monthly {
		*duration = "month"
	}

	if *langAlias != "" && *language == "" {
		*language = *langAlias
	}

	if *limit < 1 || *limit > 100 {
		*limit = 10
	}

	if *watch > 0 {
		for {
			fetchAndShow()
			color.New(color.Faint).Printf("\n Refreshing in %s... (Ctrl + C to stop)\n", *watch)
			time.Sleep(*watch)
		}
	} else {
		fetchAndShow()

		color.New(color.FgCyan).Println("\nPress Enter to exit...")
		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func fetchAndShow() ([]Repo, bool) {
	query := buildQuery()
	apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=%d", url.QueryEscape(query), *limit)

	client := createHTTPClient()
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		color.Yellow("Failed to fetch live data ----> using cache")
		return loadCache(), true
	}

	defer resp.Body.Close()

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || len(data.Items) == 0 {
		return loadCache(), true
	}

	sort.Slice(data.Items, func(i, j int) bool {
		return dailyStars(&data.Items[i]) > dailyStars(&data.Items[j])
	})

	return data.Items, false
}

func buildQuery() string {
	daysMap := map[string]int{
		"day":   1,
		"week":  7,
		"month": 30,
		"year":  365,
	}
	days := daysMap[strings.ToLower(*duration)]

	if days == 0 {
		days = 7
	}

	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	parts := []string{fmt.Sprintf("created:>%s", since)}

	if *language != "" {
		parts = append(parts, "language:"+strings.ToLower(*language))
	}
	if *spoken != "" {
		parts = append(parts, "language:"+strings.ToLower(*spoken))
	}

	parts = append(parts, "stars:>50", "fork:false")
	return strings.Join(parts, " ")
}

func dailyStars(r *Repo) int {

}
