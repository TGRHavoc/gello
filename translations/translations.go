package translations

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
	Organisations map[string]OrgTranslation   `json:"orgs"`
	Boards        map[string]BoardTranslation `json:"boards"`
	Lists         map[string]ListTranslation  `json:"lists"`
	Users         map[string]UserTranslation  `json:"users"`
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
	}
}
