package commands

import (
	"fmt"
	"strings"

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

	board.AddCmd(&ishell.Cmd{
		Name: "select",
		Help: "Selects a board. If no ID/Name is provided then a menu will appear",
		Func: selectBoard(translations, shell),
	})

	shell.AddCmd(board)
}

func selectBoard(translations *translations.Translation, shell *ishell.Shell) func(c *ishell.Context) {
	return func(c *ishell.Context) {

		var args = c.Args

		boardIDSorted := translations.SortedBoardKeys()
		shell.Set("selected_board", nil)

		if len(args) == 0 {
			var options []string
			for _, id := range boardIDSorted {
				board := translations.Boards[id]
				options = append(options, fmt.Sprintf("%s\t%-4s", id, board.Name))
			}

			selectedOption := c.MultiChoice(options, "Select a board")

			shell.Set("selected_board", translations.Boards[boardIDSorted[selectedOption]])
			println("You have selected " + boardIDSorted[selectedOption])
			return
		}

		if len(args) > 1 {
			var newArgs = make([]string, 0)
			newArgs = append(newArgs, strings.Join(c.Args, " "))
			args = newArgs
		}

		// It's either an ID or a name...
		if board, ok := translations.Boards[args[0]]; ok {
			// It was an ID
			shell.Set("selected_board", board)
			c.Println("Selected " + board.Name)
		} else {
			// Check for names
			for _, ID := range boardIDSorted {
				board := translations.Boards[ID]
				if strings.Compare(strings.ToLower(board.Name), strings.ToLower(args[0])) == 0 {
					shell.Set("selected_board", board)
					c.Println("You have selected " + ID)
				}
			}

			if shell.Get("selected_board") == nil {
				shell.Println("No board found!")
			}
		}

	}
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
