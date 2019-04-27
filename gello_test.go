package main

import (
	"log"
	"os"
	"testing"

	"github.com/adlio/trello"
)

var (
	TrelloKey   = os.Getenv("TRELLO_KEY")
	TrelloToken = os.Getenv("TRELLO_TOKEN")
)

// Make sure TRELLO_KEY and TRELLO_TOKEN is set in environment before running tests
func TestTrelloIntergration(t *testing.T) {
	if TrelloKey == "" {
		t.Error("No trello_key provided in env")
	}
	if TrelloToken == "" {
		t.Error("No trello_token provided in env")
	}

	trelloClient := trello.NewClient(TrelloKey, TrelloToken)

	t.Run("get self from tello", func(t *testing.T) {
		myself, err := trelloClient.GetMember("me", trello.Defaults())

		if err != nil {
			t.Error(err)
		}

		t.Logf("Got ID: %s, and Username: %s", myself.ID, myself.Username)
	})

	t.Run("get boards from self", func(t *testing.T) {
		boards, err := trelloClient.GetMyBoards(trello.Defaults())
		if err != nil {
			t.Error(err)
		}

		for _, board := range boards {
			board := board
			t.Run(board.Name, func(t *testing.T) {
				lists, err := board.GetLists(trello.Defaults())
				if err != nil {
					log.Fatal(err)
				}
				t.Logf("Found %d lists for %s", len(lists), board.Name)
			})
		}

		t.Logf("Board size: %d", len(boards))
	})

	var orgs []*trello.Organization
	t.Run("get organizations from self", func(t *testing.T) {
		trelloError := trelloClient.Get("member/me/organizations", trello.Defaults(), &orgs)
		if trelloError != nil {
			t.Error(trelloError)
		}

		t.Logf("Found %d orgs", len(orgs))
	})

	t.Run("get members from organizations", func(t *testing.T) {
		for _, o := range orgs {
			o := o
			t.Run(o.DisplayName, func(t *testing.T) {
				t.Parallel()
				org, _ := trelloClient.GetOrganization(o.ID, trello.Defaults())
				members, mErr := org.GetMembers(trello.Defaults())
				if mErr != nil {
					t.Error(mErr)
				}

				t.Logf("Found %d members for %s", len(members), org.DisplayName)
			})
		}
	})
}
