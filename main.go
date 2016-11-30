package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/google/gops/agent"
	log "github.com/sirupsen/logrus"
	"github.com/szuecs/go-cli/client"
	"github.com/szuecs/go-cli/conf"
)

//Buildstamp and Githash are used to set information at build time regarding
//the version of the build.
//Buildstamp is used for storing the timestamp of the build
var Buildstamp = "Not set"

//Githash is used for storing the commit hash of the build
var Githash = "Not set"

// Version is used to store the tagged version of the build
var Version = "Not set"

func init() {
	if err := agent.Start(); err != nil {
		log.Fatal(err)
	}

}

func main() {
	var (
		debug       = kingpin.Flag("debug", "enable debug mode").Default("false").Bool()
		username    = kingpin.Flag("username", "Set username to authenticate with.").Default("").String()
		oauth2Token = kingpin.Flag("oauth2-token", "Set OAuth2 Access Token.").Default("").String()
		baseURL     = kingpin.Flag("url", "Set Base URL.").Default("").String()
		_           = kingpin.Command("example", "Handle example subcmd.")
		_           = kingpin.Command("version", "show version")
	)

	switch kingpin.Parse() {
	case "example":
		fmt.Println("Example subcommand will create a client object and GET somethhing from passed URL")
		cli := createClient(*baseURL, *oauth2Token, *username, *debug)
		cli.Get(cli.Config.RealURL)
		time.Sleep(time.Second * 600)
	case "version":
		fmt.Printf(`%s Version: %s
================================
    Buildtime: %s
    GitHash: %s
`, path.Base(os.Args[0]), Version, Buildstamp, Githash)
	}
}

func createClient(url, token, username string, debug bool) client.Client {
	//loading cfg from file. it is overridden by the command line parameters
	cfg, err := conf.New()
	if err != nil {
		fmt.Printf("Could not parse config, caused by: %s\n", err)
		os.Exit(2)
	}

	if debug {
		cfg.DebugEnabled = debug
		fmt.Println("Enabled debug mode")
	}

	// URL, cli parameter overwrites config
	if url != "" {
		cfg.URL = url
	}
	cfg.RealURL = parseURL(cfg.URL)

	// client
	cli := client.Client{
		Config:      cfg,
		AccessToken: token,
	}

	// username
	cli.GetUsername(username)

	// AccessToken
	if cli.AccessToken == "" {
		cli.GetAccessToken()
	}

	return cli
}

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		fmt.Printf("Failed to parse url %s, caused by: %s", s, err)
		os.Exit(2)
	}
	return u
}
