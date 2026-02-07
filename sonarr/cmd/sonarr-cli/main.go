package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/sonarr/pkg/sonarr"
)

var version = "dev"

var (
	flagURL      string
	flagAPIKey   string
	flagInsecure bool
	flagDays     int
	flagLimit    int
)

var rootCmd = &cobra.Command{
	Use:     "sonarr-cli",
	Short:   "CLI for Sonarr",
	Version: version,
}

// Series commands
var seriesCmd = &cobra.Command{
	Use:   "series",
	Short: "Manage TV series",
}

var seriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all series",
	RunE:  runSeriesList,
}

var seriesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show series details",
	Args:  cobra.ExactArgs(1),
	RunE:  runSeriesShow,
}

// Other commands
var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Show upcoming episodes",
	RunE:  runCalendar,
}

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Show download queue",
	RunE:  runQueue,
}

var wantedCmd = &cobra.Command{
	Use:   "wanted",
	Short: "Show missing episodes",
	RunE:  runWanted,
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for new series",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Sonarr URL (or set SONARR_URL)")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "apikey", "", "API key (or set SONARR_API_KEY)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS certificate verification")

	calendarCmd.Flags().IntVar(&flagDays, "days", 7, "Number of days to show")
	wantedCmd.Flags().IntVar(&flagLimit, "limit", 20, "Maximum items to return")

	seriesCmd.AddCommand(seriesListCmd, seriesShowCmd)

	rootCmd.AddCommand(seriesCmd)
	rootCmd.AddCommand(calendarCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(wantedCmd)
	rootCmd.AddCommand(searchCmd)
}

func getConfig() (url, apiKey string, err error) {
	url = flagURL
	if url == "" {
		url = os.Getenv("SONARR_URL")
	}
	if url == "" {
		return "", "", sonarr.ConfigError("missing URL. Use --url or set SONARR_URL")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", "", sonarr.ConfigError("URL must start with http:// or https://")
	}

	apiKey = flagAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("SONARR_API_KEY")
	}
	if apiKey == "" {
		return "", "", sonarr.ConfigError("missing API key. Use --apikey or set SONARR_API_KEY")
	}

	return url, apiKey, nil
}

func getClient() (*sonarr.Client, error) {
	url, apiKey, err := getConfig()
	if err != nil {
		return nil, err
	}
	return sonarr.NewClient(url, apiKey, flagInsecure), nil
}

func handleError(err error) {
	sonarr.PrintError(err)
	if se, ok := err.(*sonarr.SonarrError); ok {
		os.Exit(se.ExitCode())
	}
	os.Exit(1)
}

func runSeriesList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	series, err := client.ListSeries()
	if err != nil {
		handleError(err)
		return nil
	}

	var items []sonarr.SeriesListItem
	for _, s := range series {
		items = append(items, s.ToListItem())
	}

	if err := sonarr.PrintYAML(sonarr.SeriesList{Series: items}); err != nil {
		handleError(err)
	}
	return nil
}

func runSeriesShow(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		handleError(sonarr.ConfigError("invalid series ID"))
		return nil
	}

	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	series, err := client.GetSeries(id)
	if err != nil {
		handleError(err)
		return nil
	}

	if err := sonarr.PrintYAML(sonarr.SeriesDetail{Series: series.ToDetail()}); err != nil {
		handleError(err)
	}
	return nil
}

func runCalendar(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	entries, err := client.GetCalendar(flagDays)
	if err != nil {
		handleError(err)
		return nil
	}

	var items []sonarr.CalendarItem
	for _, e := range entries {
		items = append(items, e.ToListItem())
	}

	if err := sonarr.PrintYAML(sonarr.CalendarList{Episodes: items}); err != nil {
		handleError(err)
	}
	return nil
}

func runQueue(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	queue, err := client.GetQueue()
	if err != nil {
		handleError(err)
		return nil
	}

	var items []sonarr.QueueItem
	for _, q := range queue.Records {
		items = append(items, q.ToListItem())
	}

	if err := sonarr.PrintYAML(sonarr.QueueList{Queue: items, Total: queue.TotalRecords}); err != nil {
		handleError(err)
	}
	return nil
}

func runWanted(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	wanted, err := client.GetWanted(flagLimit)
	if err != nil {
		handleError(err)
		return nil
	}

	var items []sonarr.WantedItem
	for _, e := range wanted.Records {
		items = append(items, e.ToWantedItem())
	}

	if err := sonarr.PrintYAML(sonarr.WantedList{Episodes: items, Total: wanted.TotalRecords}); err != nil {
		handleError(err)
	}
	return nil
}

func runSearch(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	query := strings.Join(args, " ")
	results, err := client.Search(query)
	if err != nil {
		handleError(err)
		return nil
	}

	var items []sonarr.SearchResultItem
	for _, r := range results {
		items = append(items, r.ToListItem())
	}

	if err := sonarr.PrintYAML(sonarr.SearchResultList{Results: items}); err != nil {
		handleError(err)
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
