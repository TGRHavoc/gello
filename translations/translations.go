package translations

type Me struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Initials string `json:"initials,omitempty"`
}

type Translation struct {
	Organisations map[string]OrgTranslation   `json:"orgs"`
	Boards        map[string]BoardTranslation `json:"boards"`
	Lists         map[string]ListTranslation  `json:"lists"`
	Users         map[string]UserTranslation  `json:"users"`
}

type OrgTranslation struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayname"`
	URL         string `json:"url"`
}

type BoardTranslation struct {
	Organisation string `json:"organization"`
	Name         string `json:"name"`
	Closed       bool   `json:"closed"`
	URL          string `json:"url"`
}

type ListTranslation struct {
	Board string `json:"board"`
	Name  string `json:"name"`
}

type UserTranslation struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Initals  string `json:"initials"`
}

type LabelTranslations struct {
}

func Defaults() Translation {
	return Translation{
		Boards:        make(map[string]BoardTranslation),
		Organisations: make(map[string]OrgTranslation),
		Lists:         make(map[string]ListTranslation),
		Users:         make(map[string]UserTranslation),
	}
}
