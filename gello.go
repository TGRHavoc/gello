package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	c "github.com/TGRHavoc/gello/commands"
	t "github.com/TGRHavoc/gello/translations"
	"github.com/abiosoft/ishell"
	"github.com/adlio/trello"
	"github.com/faith/color"
)

// Config represents a config object that is stored in a JSON file in the user's home directory
type Config struct {
	TrelloAPIKey string        `json:"apiKey"`
	TrelloToken  string        `json:"token"`
	Me           t.Me          `json:"me,omitempty"`
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

		writeJSON(configDir, def)

		os.Exit(0)
	}

	file, _ := os.Open(configDir)
	defer file.Close()
	decoder := json.NewDecoder(file)

	config := Config{}
	confErr := decoder.Decode(&config)
	if confErr != nil {
		log.Fatal(confErr)
	}

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

	go config.syncMe(client)
	go config.syncBoardsAndLists(client)
	config.syncOrganizations(client)
	config.syncUsers(client)

	done := make(chan bool)
	go realStart(&config, done)

	<-done

	println("Goodbye!")
	// Make sure if we changed anything, it gets saved
	writeJSON(configDir, config)
}

func (config *Config) syncMe(client *trello.Client) {
	if config.Me.Name == "" {
		println("Syncronising self values")
		member, te := client.GetMember("me", trello.Defaults())
		if te != nil {
			log.Fatal(te)
		}
		config.Me.ID = member.ID
		config.Me.Initials = member.Initials
		config.Me.Name = member.FullName
		config.Me.Username = member.Username
	}
}

func (config *Config) syncBoardsAndLists(client *trello.Client) {
	if len(config.Translations.Boards) == 0 || len(config.Translations.Lists) == 0 {
		println("Syncronizing boards and lists")

		// Make sure they exist in the config
		config.Translations.Lists = make(map[string]t.ListTranslation)
		config.Translations.Boards = make(map[string]t.BoardTranslation)

		boards, te := client.GetMyBoards(trello.Defaults())
		if te != nil {
			log.Fatal(te)
		}

		putBoardsIntoConfig(config, boards)
	}
}

func (config *Config) syncOrganizations(client *trello.Client) {
	if len(config.Translations.Organisations) == 0 {
		println("Syncronizing orgs")

		config.Translations.Organisations = make(map[string]t.OrgTranslation)

		var orgs []*trello.Organization
		trelloError := client.Get("members/me/organizations", trello.Defaults(), &orgs)
		if trelloError != nil {
			log.Fatal(trelloError)
		}

		putOrgsIntoConfig(config, orgs)
	}
}

func (config *Config) syncUsers(client *trello.Client) {
	if len(config.Translations.Users) == 0 {
		println("Syncronizing users")

		config.Translations.Users = make(map[string]t.UserTranslation)

		for id := range config.Translations.Organisations {
			org, _ := client.GetOrganization(id, trello.Defaults())
			members, _ := org.GetMembers(trello.Defaults())
			for _, m := range members {
				if m.ID == config.Me.ID || config.Translations.Users[m.ID].Name != "" {
					continue
				}
				config.Translations.Users[m.ID] = t.UserTranslation{
					Name:     m.FullName,
					Username: m.Username,
					Initials: m.Initials,
				}
			}
		}
	}
}

func realStart(config *Config, amIDone chan bool) {
	// create new shell.
	shell := ishell.New()

	// display welcome info.
	shell.Println(color.BlueString("Trello ") + color.YellowString("CLI"))
	c.InitCard(shell, &config.Translations)
	c.InitOrgs(shell, &config.Translations)

	shell.Interrupt(func(c *ishell.Context, count int, someString string) {
		if count >= 2 {
			c.Println("Interrupted")
			shell.Stop()
			amIDone <- true
			return // make sure we don't print the message below again
		}
		c.Println("Input Ctrl-c once more to exit")
	})

	// run shell
	shell.Run()

	amIDone <- true
}

func writeJSON(filename string, data interface{}) {
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
}

func getDefaultConfig() Config {
	return Config{
		TrelloAPIKey: "GET FROM https://trello.com/app-key/",
		Translations: t.Defaults(),
	}
}

func putBoardsIntoConfig(config *Config, boards []*trello.Board) {
	ourLists := make([]*trello.List, 0)

	for _, b := range boards {
		config.Translations.Boards[b.ID] = t.BoardTranslation{
			Name:         b.Name,
			Organisation: b.IDOrganization,
			Closed:       b.Closed,
			URL:          b.URL,
		}

		lists, err := b.GetLists(trello.Defaults())
		if err != nil {
			log.Fatal(err)
		}
		ourLists = append(ourLists, lists...)
	}

	for _, l := range ourLists {
		config.Translations.Lists[l.ID] = t.ListTranslation{
			Name:  l.Name,
			Board: l.IDBoard,
		}
	}

}

func putOrgsIntoConfig(config *Config, orgs []*trello.Organization) {
	for _, o := range orgs {
		config.Translations.Organisations[o.ID] = t.OrgTranslation{
			Name:        o.Name,
			DisplayName: o.DisplayName,
			URL:         o.URL,
		}
	}
}
