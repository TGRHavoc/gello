package commands

import (
	"github.com/TGRHavoc/gello/translations"
	"github.com/abiosoft/ishell"
	"github.com/faith/color"
)

func InitCard(shell *ishell.Shell, translations *translations.Translation) {
	board := &ishell.Cmd{
		Name: "board",
		Help: "board related commands",
		Func: func(c *ishell.Context) {
			c.Println("Board command")
		},
	}

	board.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "List the available boards. Red represents closed boards.",
		Func: listBoards(translations),
	})

	shell.AddCmd(board)
}

func listBoards(translations *translations.Translation) func(c *ishell.Context) {
	return func(c *ishell.Context) {
		// 5a 4ca1221676df30b716c4f7
		c.Printf("ID\t\t\t\t%-4s\n", "Name")
		for _, i := range translations.SortedBoardKeys() {
			b := translations.Boards[i]

			var nameColor func(a ...interface{}) string
			if b.Closed {
				nameColor = color.New(color.FgRed).SprintFunc()
			} else {
				nameColor = color.New(color.FgYellow).SprintFunc()
			}

			c.Printf("%s\t%-4s\n", color.GreenString(i), nameColor(b.Name))
		}
	}
}
