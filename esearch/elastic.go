package esearch

import "encoding/json"

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
	Terms         = "_terms"
	Histogram     = "_histogram"
	Range         = "_range"
	DateRange     = "_dateRange"
	DateHistogram = "_dateHistogram"
	AggsFilter    = "_filter"

	Avg           = "_avg"
	Max           = "_max"
	Min           = "_min"
	Sum           = "_sum"
	ValueCount    = "_valueCount"
	Stats         = "_stats"
	ExtendedStats = "_extendedStats"
	TopHits       = "_topHits"
	Cardinality   = "_cardinality"
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
	Mode  string    `json:"mode,omitempty"`
}

type SortMap map[string]OrderType

type Sorter interface {
	Sort()
}

func (s Sort) Sort() {

}

func (s SortMap) Sort() {

}

type Paginator interface {
	Page() uint
}

type Uint uint

func (i Uint) Page() uint {
	return uint(i)
}

type ElasticQuery struct {
	Source     []string  `json:"_source,omitempty"`
	Size       Paginator `json:"size,omitempty"`
	From       Paginator `json:"from,omitempty"`
	Sort       []Sorter  `json:"sort,omitempty"`
	Query      `json:"query,omitempty"`
	PostFilter Query                 `json:"post_filter,omitempty"`
	Aggs       map[string]Aggregator `json:"aggs,omitempty"`
	Collapse   Collapsor             `json:"collapse,omitempty"`
}

type Query map[string]QueryBuilder

type BoolQuery struct {
	Must               []BoolBuilder `json:"must,omitempty"`
	MustNot            []BoolBuilder `json:"must_not,omitempty"`
	Should             []BoolBuilder `json:"should,omitempty"`
	Filter             []BoolBuilder `json:"filter,omitempty"`
	MinimumShouldMatch int           `json:"minimum_should_match,omitempty"`
}

func (b *BoolQuery) QueryBuild() string {
	return ""
}

func (query Query) QueryBuild() string {
	return ""
}

func (b *BoolQuery) BoolBuild() string {
	return ""
}

func (query Query) BoolBuild() string {
	return ""
}

type Aggregator interface {
	Aggregate(subAgg map[string]Aggregator)
}

type Request interface {
	Query() ([]byte, error)
	ScrollQuery() ([]byte, error)
}

type Collapsor interface {
	Collapse()
}

type ExpandInnerHits interface {
	ExpandHits()
}

type CountResult struct {
	Value int `json:"value"`
}

type ArithmeticResult struct {
	Value float64 `json:"value"`
}

type StatsResult struct {
	Count int     `json:"count"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Sum   float64 `json:"sum"`
}

type ExtendStatsResult struct {
	*StatsResult
	SumOfSquares       float64 `json:"sum_of_squares"`
	Variance           float64 `json:"variance"`
	StdDeviation       float64 `json:"std_deviation"`
	StdDeviationBounds `json:"std_deviation_bounds"`
}

type StdDeviationBounds struct {
	Upper float64 `json:"upper"`
	Lower float64 `json:"lower"`
}

type TermsResult struct {
	DocCountErrorUpperBound int      `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int      `json:"sum_other_doc_count"`
	Buckets                 []Bucket `json:"buckets"`
}

type Bucket struct {
	Key         string `json:"key"`
	DocCount    int    `json:"doc_count"`
	KeyAsString string `json:"key_as_string,omitempty"` // histogram, date_histogram 使用
	RangeBucket
	Aggs AggsResult `json:"aggs,omitempty"`
}

type HistogramResult struct {
	Buckets []Bucket `json:"buckets"`
}

type RangeResult struct {
	Buckets []Bucket
}

type RangeBucket struct {
	To           any    `json:"to,omitempty"`             // range, date_range 使用
	ToAsString   string `json:"to_as_string,omitempty"`   // range, date_range 使用
	From         any    `json:"from,omitempty"`           // range, date_range 使用
	FromAsString string `json:"from_as_string,omitempty"` // range, date_range 使用
}

type HitsResult struct {
	Total    int          `json:"total"`
	MaxScore float64      `json:"max_score"`
	Hits     []HitsBucket `json:"hits"`
}

type HitsBucket struct {
	Index string  `json:"_index"`
	Type  string  `json:"_type"`
	Id    string  `json:"_id"`
	Score float64 `json:"_score"`
	//Source    map[string]any       `json:"_source"`
	Source    json.RawMessage      `json:"_source"`
	Fields    map[string][]string  `json:"fields,omitempty"`
	InnerHits map[string]InnerHits `json:"inner_hits,omitempty"`
	Sort      []any                `json:"sort,omitempty"`
}

type InnerHits struct {
	HitsResult
}

type AggsResult struct {
	Terms         map[string]*TermsResult
	Histogram     map[string]*HistogramResult
	Range         map[string]*RangeResult
	Count         map[string]*CountResult
	Arithmetic    map[string]*ArithmeticResult
	Stats         map[string]*StatsResult
	ExtendedStats map[string]*ExtendStatsResult
	TopHits       *HitsResult
}
