package termlevel

type TermQuery struct {
	Term     map[string]any        `json:"term,omitempty"`
	Range    map[string]RangeQuery `json:"range,omitempty"`
	Terms    map[string][]any      `json:"terms,omitempty"`
	Exists   map[string]string     `json:"exists,omitempty"`
	Regexp   map[string]Regexp     `json:"regexp,omitempty"`
	Wildcard map[string]Wildcard   `json:"wildcard,omitempty"`
}

type RangeQuery struct {
	Gte any `json:"gte,omitempty"`
	Lte any `json:"lte,omitempty"`
	Gt  any `json:"gt,omitempty"`
	Lt  any `json:"lt,omitempty"`
}

type Regexp struct {
	Value string `json:"value"`
	RegexpParam
}

type RegexpParam struct {
	Flags                 string `json:"flags,omitempty"`
	CaseInsensitive       bool   `json:"case_insensitive,omitempty"`
	MaxDeterminizedStates int32  `json:"max_determinized_states,omitempty"` // 相当于 es中的integer类型
	Rewrite               string `json:"rewrite,omitempty"`
}

type Wildcard struct {
	Value string `json:"value"`
	WildcardParam
}

type WildcardParam struct {
	Boost           float32 `json:"boost,omitempty"` // 相当于es中的float类型
	CaseInsensitive bool    `json:"case_insensitive,omitempty"`
	Rewrite         string  `json:"rewrite,omitempty"`
	Wildcard        string  `json:"wildcard,omitempty"`
}

func (term TermQuery) BoolBuild() string {
	return ""
}
