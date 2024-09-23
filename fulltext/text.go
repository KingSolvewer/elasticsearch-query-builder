package fulltext

type MultiMatcher interface {
	MultiMatch() string
}

type TextQuery struct {
	Match       map[string]MatchQuery       `json:"match,omitempty"`
	MatchPhrase map[string]MatchPhraseQuery `json:"match_phrase,omitempty"`
	MultiMatch  MultiMatcher                `json:"multi_match,omitempty"`
}

type MatchQuery struct {
	Query    string `json:"query,omitempty"`
	Operator string `json:"operator,omitempty"`
}

type MatchPhraseQuery struct {
	Query string `json:"query,omitempty"`
}

type MultiMatchQuery struct {
	Query    string   `json:"query,omitempty"`
	Type     string   `json:"type,omitempty"`
	Fields   []string `json:"fields,omitempty"`
	Operator string   `json:"operator,omitempty"`
}

func (m MultiMatchQuery) MultiMatch() string {
	return ""
}

func (text TextQuery) BoolBuild() string {
	return ""
}
