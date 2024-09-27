package es

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
)
