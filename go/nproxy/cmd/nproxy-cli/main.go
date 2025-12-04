package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/go/nproxy/pkg/nproxy"
	"golang.org/x/term"
)

var (
	flagURL      string
	flagToken    string
	flagInsecure bool
)

var rootCmd = &cobra.Command{
	Use:     "nproxy-cli",
	Short:   "CLI for nginx-proxy-manager API",
	Version: "0.1.0",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and get a token",
	Run: func(cmd *cobra.Command, args []string) {
		url := flagURL
		if url == "" {
			url = os.Getenv("NPROXY_URL")
		}
		if url == "" {
			handleError(nproxy.ConfigError("missing URL. Use --url or set NPROXY_URL"))
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			handleError(nproxy.ConfigError("failed to read password"))
		}
		password := string(passwordBytes)

		token, err := nproxy.Login(url, email, password, flagInsecure)
		if err != nil {
			handleError(err)
		}

		fmt.Println(token)
	},
}

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage proxy hosts",
}

var hostsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all proxy hosts",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
		}

		hosts, err := client.ListProxyHosts()
		if err != nil {
			handleError(err)
		}

		output := nproxy.ProxyHostList{
			Hosts: make([]nproxy.ProxyHostListItem, len(hosts)),
		}
		for i, h := range hosts {
			output.Hosts[i] = h.ToListItem()
		}

		if err := nproxy.PrintYAML(output); err != nil {
			handleError(err)
		}
	},
}

var hostsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show proxy host details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			handleError(err)
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
		}

		host, err := client.GetProxyHost(id)
		if err != nil {
			handleError(err)
		}

		if err := nproxy.PrintYAML(host.ToProxyHost()); err != nil {
			handleError(err)
		}
	},
}

var certificatesCmd = &cobra.Command{
	Use:     "certificates",
	Aliases: []string{"certs"},
	Short:   "Manage certificates",
}

var certificatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all certificates",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
		}

		certs, err := client.ListCertificates()
		if err != nil {
			handleError(err)
		}

		output := nproxy.CertificateList{
			Certificates: make([]nproxy.CertificateListItem, len(certs)),
		}
		for i, c := range certs {
			output.Certificates[i] = c.ToListItem()
		}

		if err := nproxy.PrintYAML(output); err != nil {
			handleError(err)
		}
	},
}

var certificatesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show certificate details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := parseID(args[0])
		if err != nil {
			handleError(err)
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
		}

		cert, err := client.GetCertificate(id)
		if err != nil {
			handleError(err)
		}

		if err := nproxy.PrintYAML(cert.ToCertificate()); err != nil {
			handleError(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "nginx-proxy-manager URL (or set NPROXY_URL)")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token (or set NPROXY_TOKEN)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS certificate verification")

	hostsCmd.AddCommand(hostsListCmd)
	hostsCmd.AddCommand(hostsShowCmd)
	certificatesCmd.AddCommand(certificatesListCmd)
	certificatesCmd.AddCommand(certificatesShowCmd)

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(hostsCmd)
	rootCmd.AddCommand(certificatesCmd)
}

func parseID(arg string) (int64, error) {
	var id int64
	if _, err := fmt.Sscanf(arg, "%d", &id); err != nil {
		return 0, nproxy.ConfigError(fmt.Sprintf("invalid ID: %s", arg))
	}
	if id <= 0 {
		return 0, nproxy.ConfigError("ID must be positive")
	}
	return id, nil
}

func getConfig() (string, string, error) {
	url := flagURL
	if url == "" {
		url = os.Getenv("NPROXY_URL")
	}
	if url == "" {
		return "", "", nproxy.ConfigError("missing URL. Use --url or set NPROXY_URL")
	}

	token := flagToken
	if token == "" {
		token = os.Getenv("NPROXY_TOKEN")
	}
	if token == "" {
		return "", "", nproxy.ConfigError("missing token. Use --token or set NPROXY_TOKEN")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", "", nproxy.ConfigError("URL must start with http:// or https://")
	}

	return url, token, nil
}

func getClient() (*nproxy.Client, error) {
	url, token, err := getConfig()
	if err != nil {
		return nil, err
	}
	return nproxy.NewClient(url, token, flagInsecure), nil
}

func handleError(err error) {
	nproxy.PrintError(err)
	if ne, ok := err.(*nproxy.NproxyError); ok {
		os.Exit(ne.ExitCode())
	}
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
