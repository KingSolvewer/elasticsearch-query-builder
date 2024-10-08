package esearch

type BoolClauseType int

const (
	Must BoolClauseType = iota
	MustNot
	Should
	FilterClause
)

type RangeType int

const (
	Gt RangeType = iota
	Gte
	Lt
	Lte
)

type MatchType int

const (
	Match MatchType = iota
	MatchPhrase
	MatchPhrasePrefix
)

type FieldType int

const (
	BestFields FieldType = iota
	MostFields
	CrossFields
	Phrase
	PhrasePrefix
	BoolPrefix
)

var FieldTypeSet = map[FieldType]string{
	BestFields:   "best_fields",
	MostFields:   "most_fields",
	CrossFields:  "cross_fields",
	Phrase:       "phrase",
	PhrasePrefix: "phrase_prefix",
	BoolPrefix:   "bool_prefix",
}

type OrderType string

const (
	Asc  OrderType = "asc"
	Desc OrderType = "desc"
)

const (
	Terms         = "terms"
	Histogram     = "histogram"
	Range         = "range"
	DateHistogram = "date_histogram"
	Avg           = "avg"
	Max           = "max"
	Min           = "min"
	Sum           = "sum"
	Count         = "count"
	Stats         = "stats"
	ExtendedStats = "extended_stats"
	TopHits       = "top_hits"
)

type QueryBuilder interface {
	QueryBuild() string
}

type BoolBuilder interface {
	BoolBuild() string
}

type Sort map[string]Order

type Order struct {
	Order OrderType `json:"order"`
}

type Paginator interface {
	Page() uint
}

type Uint uint

func (i Uint) Page() uint {
	return uint(i)
}

type ElasticQuery struct {
	Source []string  `json:"_source,omitempty"`
	Size   Paginator `json:"size,omitempty"`
	From   Paginator `json:"from,omitempty"`
	Sort   []Sort    `json:"sort,omitempty"`
	Query  `json:"query,omitempty"`
	Aggs   map[string]Aggregator `json:"aggs,omitempty"`
}

type Query map[string]QueryBuilder

type BoolQuery struct {
	Must               []BoolBuilder `json:"must,omitempty"`
	MustNot            []BoolBuilder `json:"must_not,omitempty"`
	Should             []BoolBuilder `json:"should,omitempty"`
	Filter             []BoolBuilder `json:"filter,omitempty"`
	MinimumShouldMatch int           `json:"minimum_should_match,omitempty"`
}

func (b BoolQuery) QueryBuild() string {
	return ""
}

func (query Query) QueryBuild() string {
	return ""
}

func (b BoolQuery) BoolBuild() string {
	return ""
}

func (query Query) BoolBuild() string {
	return ""
}

type Aggregator interface {
	Aggregate(args string) Aggregator
}

type Request interface {
	Query() ([]byte, error)
	ScrollQuery() ([]byte, error)
	Decorate()
}
