package commands

import (
	"github.com/TGRHavoc/gello/translations"
	"github.com/abiosoft/ishell"
	"github.com/faith/color"
)

func InitOrgs(shell *ishell.Shell, translations *translations.Translation) {
	orgs := &ishell.Cmd{
		Name: "org",
		Help: "organization related commands",
		Func: func(c *ishell.Context) {
			c.Println("Orgs command")
		},
	}

	orgs.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "List the available organizations.",
		Func: listOrgs(translations),
	})

	shell.AddCmd(orgs)
}

func listOrgs(translations *translations.Translation) func(c *ishell.Context) {
	return func(c *ishell.Context) {
		// 5a 4ca1221676df30b716c4f7
		c.Printf("ID\t\t\t\t%-4s\n", "Name")
		for _, i := range translations.SortedOrgKeys() {
			b := translations.Organisations[i]
			c.Printf("%s\t%-4s\n", color.GreenString(i), color.YellowString(b.Name))
		}
	}
}
