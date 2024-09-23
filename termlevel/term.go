package termlevel

type Compare interface {
	String() string
}

type TermQueryMap map[string]string
type TermsQueryMap map[string][]string

type TermQuery struct {
	Term  map[string]string     `json:"term,omitempty"`
	Range map[string]RangeQuery `json:"range,omitempty"`
	Terms map[string][]string   `json:"terms,omitempty"`
}

type RangeQuery struct {
	Gte int `json:"gte,omitempty"`
	Lte int `json:"lte,omitempty"`
	Gt  int `json:"gt,omitempty"`
	Lt  int `json:"lt,omitempty"`
}

func (term TermQuery) BoolBuild() string {
	return ""
}

func (term TermQueryMap) BoolBuild() string {
	return ""
}

func (term TermsQueryMap) BoolBuild() string {
	return ""
}
