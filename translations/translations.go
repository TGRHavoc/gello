package translations

import (
	"sort"
	"strings"
)

// Me represents the currently logged in trello user.
// Used to speed up calls to get the current user data (rather than doing a HTTP call every time)
type Me struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Initials string `json:"initials,omitempty"`
}

// Translation represents various Trello objects for speeding up various calls (e.g. getting a org name)
type Translation struct {
	Organisations map[string]OrgTranslation    `json:"orgs"`
	Boards        map[string]BoardTranslation  `json:"boards"`
	Lists         map[string]ListTranslation   `json:"lists"`
	Users         map[string]UserTranslation   `json:"users"`
	Labels        map[string]LabelTranslations `json:"labels"`
}

// OrgTranslation represents various organizational data
type OrgTranslation struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayname"`
	URL         string `json:"url"`
}

// BoardTranslation represents various board data
type BoardTranslation struct {
	Organisation string `json:"organization"`
	Name         string `json:"name"`
	Closed       bool   `json:"closed"`
	URL          string `json:"url"`
}

// ListTranslation represents various lists data
type ListTranslation struct {
	Board string `json:"board"`
	Name  string `json:"name"`
}

// UserTranslation represents various user data
type UserTranslation struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Initials string `json:"initials"`
}

// LabelTranslations represents various label data
type LabelTranslations struct {
}

// Defaults returns a new translation structure with an empty maps
func Defaults() Translation {
	return Translation{
		Boards:        make(map[string]BoardTranslation),
		Organisations: make(map[string]OrgTranslation),
		Lists:         make(map[string]ListTranslation),
		Users:         make(map[string]UserTranslation),
		Labels:        make(map[string]LabelTranslations),
	}
}

func (translation *Translation) GetBoardMatches(name string) []BoardTranslation {
	return translation.FilterBoard(func(id string, board BoardTranslation) bool {
		return strings.Contains(board.Name, name)
	})
}

func (translsation *Translation) FilterBoard(f func(string, BoardTranslation) bool) []BoardTranslation {
	vsf := make([]BoardTranslation, 0)
	for i, v := range translsation.Boards {
		if f(i, v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (translsation *Translation) FilterOrg(f func(string, OrgTranslation) bool) []OrgTranslation {
	vsf := make([]OrgTranslation, 0)
	for i, v := range translsation.Organisations {
		if f(i, v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (translsation *Translation) FilterList(f func(string, ListTranslation) bool) []ListTranslation {
	vsf := make([]ListTranslation, 0)
	for i, v := range translsation.Lists {
		if f(i, v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (translsation *Translation) FilterUsers(f func(string, UserTranslation) bool) []UserTranslation {
	vsf := make([]UserTranslation, 0)
	for i, v := range translsation.Users {
		if f(i, v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (translsation *Translation) FilterLabel(f func(string, LabelTranslations) bool) []LabelTranslations {
	vsf := make([]LabelTranslations, 0)
	for i, v := range translsation.Labels {
		if f(i, v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (translsation *Translation) SortedBoardKeys() []string {
	sortedKeys := make([]string, 0)
	for id := range translsation.Boards {
		sortedKeys = append(sortedKeys, id)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
func (translsation *Translation) SortedOrgKeys() []string {
	sortedKeys := make([]string, 0)
	for id := range translsation.Organisations {
		sortedKeys = append(sortedKeys, id)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
func (translsation *Translation) SortedListKeys() []string {
	sortedKeys := make([]string, 0)
	for id := range translsation.Lists {
		sortedKeys = append(sortedKeys, id)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
func (translsation *Translation) SortedUsersKeys() []string {
	sortedKeys := make([]string, 0)
	for id := range translsation.Users {
		sortedKeys = append(sortedKeys, id)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
func (translsation *Translation) SortedLabelsKeys() []string {
	sortedKeys := make([]string, 0)
	for id := range translsation.Labels {
		sortedKeys = append(sortedKeys, id)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
