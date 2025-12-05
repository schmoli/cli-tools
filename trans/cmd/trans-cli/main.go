package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/trans/pkg/trans"
)

var version = "dev"

var (
	flagURL      string
	flagUser     string
	flagPass     string
	flagInsecure bool
)

var rootCmd = &cobra.Command{
	Use:     "trans-cli",
	Short:   "CLI for Transmission RPC",
	Version: version,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all torrents",
	RunE:  runList(nil),
}

var downloadingCmd = &cobra.Command{
	Use:   "downloading",
	Short: "List downloading torrents",
	RunE:  runList(func(t *trans.APITorrent) bool { return t.IsDownloading() }),
}

var seedingCmd = &cobra.Command{
	Use:   "seeding",
	Short: "List seeding torrents",
	RunE:  runList(func(t *trans.APITorrent) bool { return t.IsSeeding() }),
}

var stoppedCmd = &cobra.Command{
	Use:   "stopped",
	Short: "List stopped torrents",
	RunE:  runList(func(t *trans.APITorrent) bool { return t.IsStopped() }),
}

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show torrent details",
	Args:  cobra.ExactArgs(1),
	RunE:  runShow,
}

var addCmd = &cobra.Command{
	Use:   "add <magnet|file>",
	Short: "Add torrent (magnet URI or .torrent file)",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdd,
}

var startCmd = &cobra.Command{
	Use:   "start <id>",
	Short: "Start torrent",
	Args:  cobra.ExactArgs(1),
	RunE:  runStart,
}

var stopCmd = &cobra.Command{
	Use:   "stop <id>",
	Short: "Stop torrent",
	Args:  cobra.ExactArgs(1),
	RunE:  runStop,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Transmission URL (or set TRANSMISSION_URL)")
	rootCmd.PersistentFlags().StringVar(&flagUser, "user", "", "Username (or set TRANSMISSION_USER)")
	rootCmd.PersistentFlags().StringVar(&flagPass, "pass", "", "Password (or set TRANSMISSION_PASS)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS certificate verification")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(downloadingCmd)
	rootCmd.AddCommand(seedingCmd)
	rootCmd.AddCommand(stoppedCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}

func getConfig() (url, user, pass string, err error) {
	url = flagURL
	if url == "" {
		url = os.Getenv("TRANSMISSION_URL")
	}
	if url == "" {
		return "", "", "", trans.ConfigError("missing URL. Use --url or set TRANSMISSION_URL")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", "", "", trans.ConfigError("URL must start with http:// or https://")
	}

	user = flagUser
	if user == "" {
		user = os.Getenv("TRANSMISSION_USER")
	}

	pass = flagPass
	if pass == "" {
		pass = os.Getenv("TRANSMISSION_PASS")
	}

	return url, user, pass, nil
}

func getClient() (*trans.Client, error) {
	url, user, pass, err := getConfig()
	if err != nil {
		return nil, err
	}
	return trans.NewClient(url, user, pass, flagInsecure), nil
}

func parseID(arg string) (int64, error) {
	var id int64
	if _, err := fmt.Sscanf(arg, "%d", &id); err != nil {
		return 0, trans.ConfigError(fmt.Sprintf("invalid ID: %s", arg))
	}
	if id <= 0 {
		return 0, trans.ConfigError("ID must be positive")
	}
	return id, nil
}

func handleError(err error) {
	trans.PrintError(err)
	if te, ok := err.(*trans.TransError); ok {
		os.Exit(te.ExitCode())
	}
	os.Exit(1)
}

func runList(filter func(*trans.APITorrent) bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return nil
		}

		torrents, err := client.ListTorrents()
		if err != nil {
			handleError(err)
			return nil
		}

		var items []trans.TorrentListItem
		for _, t := range torrents {
			if filter == nil || filter(&t) {
				items = append(items, t.ToListItem())
			}
		}

		if err := trans.PrintYAML(trans.TorrentList{Torrents: items}); err != nil {
			handleError(err)
		}
		return nil
	}
}

func runShow(cmd *cobra.Command, args []string) error {
	id, err := parseID(args[0])
	if err != nil {
		handleError(err)
		return nil
	}

	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	torrent, err := client.GetTorrent(id)
	if err != nil {
		handleError(err)
		return nil
	}

	if err := trans.PrintYAML(trans.TorrentDetail{Torrent: torrent.ToDetail()}); err != nil {
		handleError(err)
	}
	return nil
}

func runAdd(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	input := args[0]
	var info *trans.TorrentAddedInfo

	if strings.HasPrefix(input, "magnet:") {
		info, err = client.AddTorrentMagnet(input)
	} else {
		info, err = client.AddTorrentFile(input)
	}

	if err != nil {
		handleError(err)
		return nil
	}

	output := struct {
		Added struct {
			ID   int64  `yaml:"id"`
			Name string `yaml:"name"`
		} `yaml:"added"`
	}{}
	output.Added.ID = info.ID
	output.Added.Name = info.Name

	if err := trans.PrintYAML(output); err != nil {
		handleError(err)
	}
	return nil
}

func runStart(cmd *cobra.Command, args []string) error {
	id, err := parseID(args[0])
	if err != nil {
		handleError(err)
		return nil
	}

	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	if err := client.StartTorrent(id); err != nil {
		handleError(err)
		return nil
	}

	fmt.Printf("started torrent %d\n", id)
	return nil
}

func runStop(cmd *cobra.Command, args []string) error {
	id, err := parseID(args[0])
	if err != nil {
		handleError(err)
		return nil
	}

	client, err := getClient()
	if err != nil {
		handleError(err)
		return nil
	}

	if err := client.StopTorrent(id); err != nil {
		handleError(err)
		return nil
	}

	fmt.Printf("stopped torrent %d\n", id)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
