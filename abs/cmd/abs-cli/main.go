package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/abs/pkg/abs"
)

var version = "dev"

var (
	flagURL      string
	flagToken    string
	flagInsecure bool
	flagLibrary  string
	flagLimit    int
)

var rootCmd = &cobra.Command{
	Use:     "abs-cli",
	Short:   "CLI for Audiobookshelf",
	Version: version,
}

// Libraries commands
var librariesCmd = &cobra.Command{
	Use:   "libraries",
	Short: "Manage libraries",
}

var librariesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all libraries",
	RunE:  runLibrariesList,
}

// Books commands
var booksCmd = &cobra.Command{
	Use:   "books",
	Short: "Manage audiobooks",
}

var booksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List audiobooks",
	RunE:  runBooksList,
}

var booksShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show audiobook details",
	Args:  cobra.ExactArgs(1),
	RunE:  runBooksShow,
}

// Other commands
var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Show listening progress",
	RunE:  runProgress,
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search audiobooks",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Trigger library scan",
	RunE:  runScan,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Audiobookshelf URL (or set ABS_URL)")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token (or set ABS_TOKEN)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS certificate verification")

	booksListCmd.Flags().StringVar(&flagLibrary, "library", "", "Library ID (uses first library if not specified)")
	booksListCmd.Flags().IntVar(&flagLimit, "limit", 50, "Maximum items to return")
	
	searchCmd.Flags().StringVar(&flagLibrary, "library", "", "Library ID (uses first library if not specified)")
	
	scanCmd.Flags().StringVar(&flagLibrary, "library", "", "Library ID (scans all if not specified)")

	librariesCmd.AddCommand(librariesListCmd)
	booksCmd.AddCommand(booksListCmd, booksShowCmd)

	rootCmd.AddCommand(librariesCmd)
	rootCmd.AddCommand(booksCmd)
	rootCmd.AddCommand(progressCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(scanCmd)
}

func getConfig() (url, token string, err error) {
	url = flagURL
	if url == "" {
		url = os.Getenv("ABS_URL")
	}
	if url == "" {
		return "", "", abs.ConfigError("missing URL. Use --url or set ABS_URL")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", "", abs.ConfigError("URL must start with http:// or https://")
	}

	token = flagToken
	if token == "" {
		token = os.Getenv("ABS_TOKEN")
	}
	if token == "" {
		return "", "", abs.ConfigError("missing token. Use --token or set ABS_TOKEN")
	}

	return url, token, nil
}

func getClient() (*abs.Client, error) {
	url, token, err := getConfig()
	if err != nil {
		return nil, err
	}
	return abs.NewClient(url, token, flagInsecure), nil
}

func handleError(err error) {
	abs.PrintError(err)
	if ae, ok := err.(*abs.AbsError); ok {
		os.Exit(ae.ExitCode())
	}
	os.Exit(1)
}

func getDefaultLibrary(client *abs.Client) (string, error) {
	if flagLibrary != "" {
		return flagLibrary, nil
	}
	
	libs, err := client.ListLibraries()
	if err != nil {
		return "", err
	}
	if len(libs) == 0 {
		return "", abs.NotFoundError("no libraries found")
	}
	return libs[0].ID, nil
}

func runLibrariesList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	libraries, err := client.ListLibraries()
	if err != nil {
		handleError(err)
		return nil
	}

	var items []abs.LibraryListItem
	for _, lib := range libraries {
		items = append(items, lib.ToListItem())
	}

	if err := abs.PrintYAML(abs.LibraryList{Libraries: items}); err != nil {
		handleError(err)
	}
	return nil
}

func runBooksList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	libraryID, err := getDefaultLibrary(client)
	if err != nil {
		handleError(err)
		return nil
	}

	items, total, err := client.ListLibraryItems(libraryID, flagLimit)
	if err != nil {
		handleError(err)
		return nil
	}

	var books []abs.BookListItem
	for _, item := range items {
		books = append(books, item.ToListItem())
	}

	if err := abs.PrintYAML(abs.BookList{Books: books, Total: total}); err != nil {
		handleError(err)
	}
	return nil
}

func runBooksShow(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	item, err := client.GetItem(args[0])
	if err != nil {
		handleError(err)
		return nil
	}

	if err := abs.PrintYAML(abs.BookDetail{Book: item.ToDetail()}); err != nil {
		handleError(err)
	}
	return nil
}

func runProgress(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	progress, err := client.GetProgress()
	if err != nil {
		handleError(err)
		return nil
	}

	// Filter to in-progress items only
	var items []abs.ProgressListItem
	for _, p := range progress {
		if !p.IsFinished && p.Progress > 0 {
			items = append(items, p.ToListItem())
		}
	}

	if err := abs.PrintYAML(abs.ProgressList{Progress: items}); err != nil {
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

	libraryID, err := getDefaultLibrary(client)
	if err != nil {
		handleError(err)
		return nil
	}

	query := strings.Join(args, " ")
	results, err := client.Search(libraryID, query)
	if err != nil {
		handleError(err)
		return nil
	}

	var output abs.SearchResults
	for _, r := range results.Book {
		output.Books = append(output.Books, r.LibraryItem.ToListItem())
	}
	for _, a := range results.Authors {
		output.Authors = append(output.Authors, a.Name)
	}
	for _, s := range results.Series {
		output.Series = append(output.Series, s.Series.Name)
	}

	if err := abs.PrintYAML(output); err != nil {
		handleError(err)
	}
	return nil
}

func runScan(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	if flagLibrary != "" {
		if err := client.ScanLibrary(flagLibrary); err != nil {
			handleError(err)
			return nil
		}
		fmt.Printf("scan started for library %s\n", flagLibrary)
	} else {
		// Scan all libraries
		libs, err := client.ListLibraries()
		if err != nil {
			handleError(err)
			return nil
		}
		for _, lib := range libs {
			if err := client.ScanLibrary(lib.ID); err != nil {
				handleError(err)
				return nil
			}
			fmt.Printf("scan started for library %s (%s)\n", lib.Name, lib.ID)
		}
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
