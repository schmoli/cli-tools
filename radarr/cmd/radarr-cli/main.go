package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/radarr/pkg/radarr"
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
	Use:     "radarr-cli",
	Short:   "CLI for Radarr",
	Version: version,
}

// Movies commands
var moviesCmd = &cobra.Command{
	Use:   "movies",
	Short: "Manage movies",
}

var moviesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all movies",
	RunE:  runMoviesList,
}

var moviesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show movie details",
	Args:  cobra.ExactArgs(1),
	RunE:  runMoviesShow,
}

// Other commands
var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Show upcoming releases",
	RunE:  runCalendar,
}

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Show download queue",
	RunE:  runQueue,
}

var wantedCmd = &cobra.Command{
	Use:   "wanted",
	Short: "Show missing movies",
	RunE:  runWanted,
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for new movies",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Radarr URL (or set RADARR_URL)")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "apikey", "", "API key (or set RADARR_API_KEY)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS certificate verification")

	calendarCmd.Flags().IntVar(&flagDays, "days", 30, "Number of days to show")
	wantedCmd.Flags().IntVar(&flagLimit, "limit", 20, "Maximum items to return")

	moviesCmd.AddCommand(moviesListCmd, moviesShowCmd)

	rootCmd.AddCommand(moviesCmd)
	rootCmd.AddCommand(calendarCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(wantedCmd)
	rootCmd.AddCommand(searchCmd)
}

func getConfig() (url, apiKey string, err error) {
	url = flagURL
	if url == "" {
		url = os.Getenv("RADARR_URL")
	}
	if url == "" {
		return "", "", radarr.ConfigError("missing URL. Use --url or set RADARR_URL")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", "", radarr.ConfigError("URL must start with http:// or https://")
	}

	apiKey = flagAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("RADARR_API_KEY")
	}
	if apiKey == "" {
		return "", "", radarr.ConfigError("missing API key. Use --apikey or set RADARR_API_KEY")
	}

	return url, apiKey, nil
}

func getClient() (*radarr.Client, error) {
	url, apiKey, err := getConfig()
	if err != nil {
		return nil, err
	}
	return radarr.NewClient(url, apiKey, flagInsecure), nil
}

func handleError(err error) {
	radarr.PrintError(err)
	if re, ok := err.(*radarr.RadarrError); ok {
		os.Exit(re.ExitCode())
	}
	os.Exit(1)
}

func runMoviesList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	movies, err := client.ListMovies()
	if err != nil {
		handleError(err)
		return nil
	}

	var items []radarr.MovieListItem
	for _, m := range movies {
		items = append(items, m.ToListItem())
	}

	if err := radarr.PrintYAML(radarr.MovieList{Movies: items}); err != nil {
		handleError(err)
	}
	return nil
}

func runMoviesShow(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		handleError(radarr.ConfigError("invalid movie ID"))
		return nil
	}

	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	movie, err := client.GetMovie(id)
	if err != nil {
		handleError(err)
		return nil
	}

	if err := radarr.PrintYAML(radarr.MovieDetail{Movie: movie.ToDetail()}); err != nil {
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

	var items []radarr.CalendarItem
	for _, e := range entries {
		items = append(items, e.ToListItem())
	}

	if err := radarr.PrintYAML(radarr.CalendarList{Movies: items}); err != nil {
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

	var items []radarr.QueueItem
	for _, q := range queue.Records {
		items = append(items, q.ToListItem())
	}

	if err := radarr.PrintYAML(radarr.QueueList{Queue: items, Total: queue.TotalRecords}); err != nil {
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

	var items []radarr.WantedItem
	for _, m := range wanted.Records {
		items = append(items, m.ToWantedItem())
	}

	if err := radarr.PrintYAML(radarr.WantedList{Movies: items, Total: wanted.TotalRecords}); err != nil {
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

	var items []radarr.SearchResultItem
	for _, r := range results {
		items = append(items, r.ToListItem())
	}

	if err := radarr.PrintYAML(radarr.SearchResultList{Results: items}); err != nil {
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
