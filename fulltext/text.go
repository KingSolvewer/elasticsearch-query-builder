package fulltext

type MultiMatcher interface {
	MultiMatch() string
}

type TextQuery struct {
	Match             map[string]MatchQuery `json:"match,omitempty"`
	MatchPhrase       map[string]MatchQuery `json:"match_phrase,omitempty"`
	MatchPhrasePrefix map[string]MatchQuery `json:"match_phrase_fix,omitempty"`
	MultiMatch        *MultiMatchQuery      `json:"multi_match,omitempty"`
}

type MatchQuery struct {
	Query string `json:"query,omitempty"`
	AppendParams
}

type MultiMatchQuery struct {
	Query  string   `json:"query,omitempty"`
	Type   string   `json:"type,omitempty"`
	Fields []string `json:"fields,omitempty"`
	AppendParams
}

type AppendParamsFunc func() AppendParams

type AppendParams struct {
	Analyzer           string  `json:"analyzer,omitempty"`
	Boost              float32 `json:"boost,omitempty"`
	Operator           string  `json:"operator,omitempty"`
	Fuzziness          string  `json:"fuzziness,omitempty"`
	MaxExpansions      int     `json:"max_expansions,omitempty"`
	PrefixLength       int     `json:"prefix_length,omitempty"`
	Lenient            bool    `json:"lenient,omitempty"`
	MinimumShouldMatch string  `json:"minimum_should_match,omitempty"`
}

func (m MultiMatchQuery) MultiMatch() string {
	return ""
}

func (text TextQuery) BoolBuild() string {
	return ""
}
