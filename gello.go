package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/adlio/trello"
	"github.com/fatih/color"
	t "github.com/tgrhavoc/gello/translations"
)

type Config struct {
	TrelloAPIKey string        `json:"apiKey"`
	TrelloToken  string        `json:"token"`
	Me           t.Me          `json:"me"`
	Translations t.Translation `json:"translations,omitempty"`
}

func main() {

	// Get the directory we want to use for the config
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	homedir := usr.HomeDir
	configDir := path.Join(homedir, ".gello", "config.json")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		fmt.Println(configDir + " doesn't exist. Making it with defaults")
		def := getDefaultConfig()

		writeJson(configDir, def)

		os.Exit(0)
	}

	file, _ := os.Open(configDir)
	defer file.Close()
	decoder := json.NewDecoder(file)

	config := Config{}
	decoder.Decode(&config)

	if config.TrelloAPIKey == "GET FROM https://trello.com/app-key/" {
		println("Please put your API key from https://trello.com/app-key into the config file and re-run.")
		println("The config file can be found at: " + configDir)
		os.Exit(1)
	}

	if config.TrelloToken == "" {
		url := fmt.Sprintf("https://trello.com/1/connect?key=%s&name=gello&response_type=token&scope=account,read,write&expiration=never", config.TrelloAPIKey)
		fmt.Printf("Please put the token from %s into the config file.", url)
		println("The config file can be found at: " + configDir)
		os.Exit(1)
	}

	client := trello.NewClient(config.TrelloAPIKey, config.TrelloToken)

	member, te := client.GetMember("me", trello.Defaults())
	if te != nil {
		log.Fatal(te)
	}
	config.Me.ID = member.ID
	config.Me.Initials = member.Initials
	config.Me.Name = member.FullName
	config.Me.Username = member.Username

	boards, te := member.GetBoards(trello.Defaults())
	if te != nil {
		log.Fatal(te)
	}

	putBoardsIntoConfig(&config, boards)

	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println("Trello CLI")

	// register a function for "greet" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "greet",
		Help: "greet user",
		Func: func(c *ishell.Context) {
			d := color.New(color.FgCyan).SprintFunc()

			c.Println("Hello, ", strings.Join(c.Args, " "), d("world!"))
		},
	})

	shell.Interrupt(func(c *ishell.Context, count int, someString string) {
		if count >= 2 {
			c.Println("Interrupted")
			writeJson(configDir, config)
			os.Exit(1)
		}
		c.Println("Input Ctrl-c once more to exit")
	})

	// run shell
	shell.Run()

	// Make sure if we changed anything, it gets saved
	writeJson(configDir, config)
}

func writeJson(filename string, data interface{}) {
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
}

func getDefaultConfig() Config {
	fmt.Println(t.Defaults())
	return Config{
		TrelloAPIKey: "GET FROM https://trello.com/app-key/",
		Translations: t.Defaults(),
	}
}

func putBoardsIntoConfig(config *Config, boards []*trello.Board) {
	for _, b := range boards {
		config.Translations.Boards[b.ID] = t.BoardTranslation{
			Name:         b.Name,
			Organisation: b.IDOrganization,
			Closed:       b.Closed,
			URL:          b.URL,
		}
	}
}
