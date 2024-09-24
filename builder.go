package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/fulltext"
	"github.com/KingSolvewer/elasticsearch-query-builder/termlevel"
)

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

type NestedFunc func(c *Condition) *Condition

type Condition struct {
	where              map[BoolClauseType][]BoolBuilder
	nested             map[BoolClauseType][]NestedFunc
	minimumShouldMatch int
	*Dsl
}

var condition = NewCondition()

func NewCondition() *Condition {
	return &Condition{
		where:  make(map[BoolClauseType][]BoolBuilder),
		nested: make(map[BoolClauseType][]NestedFunc),
	}
}

func GetCondition() *Condition {
	condition.compile()
	return condition
}

// Where Must term 查询语句
func Where(field string, value string) *Condition {
	return condition.Where(field, value)
}

// Where Must term 查询语句
func (c *Condition) Where(field string, value string) *Condition {
	c.termQuery(Must, field, value)

	return c
}

// WhereIn Must terms 查询语句
func WhereIn(field string, value []string) *Condition {
	return condition.WhereIn(field, value)
}

// WhereIn Must terms 查询语句
func (c *Condition) WhereIn(field string, value []string) *Condition {
	c.termsQuery(Must, field, value)

	return c
}

// WhereBetween Must Range ( gte:>=a and lte:<=b) 查询语句
func WhereBetween(field string, value1, value2 int) *Condition {
	return condition.WhereBetween(field, value1, value2)
}

// WhereBetween Must Range ( gte:>=a and lte:<=b) 查询语句
func (c *Condition) WhereBetween(field string, value1, value2 int) *Condition {

	c.whereBetween(Must, field, value1, value2)

	return c
}

// WhereRange Must Range 单向范围查询, gt:>a, lt:<a等等
func WhereRange(field string, value int, rangeType RangeType) *Condition {
	return condition.WhereRange(field, value, rangeType)
}

// WhereRange Must Range 单向范围查询, gt:>a, lt:<a等等
func (c *Condition) WhereRange(field string, value int, rangeType RangeType) *Condition {

	c.whereRange(Must, field, value, rangeType)

	return c
}

// WhereMatch Must match 匹配
func WhereMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Condition {
	return condition.WhereMatch(field, value, matchType, params)
}

// WhereMatch Must match 匹配
func (c *Condition) WhereMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Condition {
	c.whereMatch(Must, field, value, matchType, params)
	return c
}

// WhereMultiMatch Must multi_match 匹配
func WhereMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Condition {
	return condition.WhereMultiMatch(field, value, fieldType, params)
}

// WhereMultiMatch Must multi_match 匹配
func (c *Condition) WhereMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Condition {
	c.whereMultiMatch(Must, field, value, fieldType, params)

	return c
}

// WhereNested Must 嵌套查询, 例如嵌套 should 语句
func WhereNested(fn NestedFunc) *Condition {
	return condition.WhereNested(fn)
}

// WhereNested Must 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (c *Condition) WhereNested(fn NestedFunc) *Condition {
	c.nested[Must] = append(c.nested[Must], fn)

	return c
}

// WhereNot MustNot term 查询
func WhereNot(field string, value string) *Condition {
	return condition.WhereNot(field, value)
}

// WhereNot MustNot term 查询
func (c *Condition) WhereNot(field string, value string) *Condition {
	c.termQuery(MustNot, field, value)

	return c
}

// WhereNotIn  MustNot terms 查询语句
func WhereNotIn(field string, value []string) *Condition {
	return condition.WhereIn(field, value)
}

// WhereNotIn Must terms 查询语句
func (c *Condition) WhereNotIn(field string, value []string) *Condition {
	c.termsQuery(MustNot, field, value)

	return c
}

// WhereNotBetween MustNot Range ( gte:>=a and lte:<=b) 查询语句
func WhereNotBetween(field string, value1, value2 int) *Condition {
	return condition.WhereNotBetween(field, value1, value2)
}

// WhereNotBetween MustNot Range ( gte:>=a and lte:<=b) 查询语句
func (c *Condition) WhereNotBetween(field string, value1, value2 int) *Condition {
	c.whereBetween(MustNot, field, value1, value2)

	return c
}

// WhereNotRange MustNot Range 单向范围查询, gt:>a, lt:<a等等
func WhereNotRange(field string, value int, rangeType RangeType) *Condition {
	return condition.WhereNotRange(field, value, rangeType)
}

// WhereNotRange MustNot Range 单向范围查询, gt:>a, lt:<a等等
func (c *Condition) WhereNotRange(field string, value int, rangeType RangeType) *Condition {
	c.whereRange(MustNot, field, value, rangeType)

	return c
}

// WhereNotMatch MustNot match 匹配
func WhereNotMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Condition {
	return condition.WhereNotMatch(field, value, matchType, params)
}

// WhereNotMatch MustNot match 匹配
func (c *Condition) WhereNotMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Condition {
	c.whereMatch(MustNot, field, value, matchType, params)
	return c
}

// WhereNotMultiMatch MustNot multi_match 匹配
func WhereNotMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Condition {
	return condition.WhereNotMultiMatch(field, value, fieldType, params)
}

// WhereNotMultiMatch MustNot multi_match 匹配
func (c *Condition) WhereNotMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Condition {
	c.whereMultiMatch(MustNot, field, value, fieldType, params)

	return c
}

// WhereNotNested MustNot 嵌套查询, 例如嵌套 should 语句
func WhereNotNested(fn NestedFunc) *Condition {
	return condition.WhereNotNested(fn)
}

// WhereNotNested Must 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (c *Condition) WhereNotNested(fn NestedFunc) *Condition {
	c.nested[MustNot] = append(c.nested[MustNot], fn)

	return c
}

// OrWhere Should term 查询语句
func OrWhere(field string, value string) *Condition {
	return condition.OrWhere(field, value)
}

// OrWhere Should term 查询语句
func (c *Condition) OrWhere(field string, value string) *Condition {
	c.termQuery(Should, field, value)

	return c
}

// OrWhereIn Should terms 查询语句
func OrWhereIn(field string, value []string) *Condition {
	return condition.OrWhereIn(field, value)
}

// OrWhereIn Should terms 查询语句
func (c *Condition) OrWhereIn(field string, value []string) *Condition {
	c.termsQuery(Should, field, value)

	return c
}

// OrWhereBetween Should Range ( gte:>=a and lte:<=b) 查询语句
func OrWhereBetween(field string, value1, value2 int) *Condition {
	return condition.OrWhereBetween(field, value1, value2)
}

// OrWhereBetween Should Range ( gte:>=a and lte:<=b) 查询语句
func (c *Condition) OrWhereBetween(field string, value1, value2 int) *Condition {

	c.whereBetween(Should, field, value1, value2)

	return c
}

// OrWhereRange Should Range 单向范围查询, gt:>a, lt:<a等等
func OrWhereRange(field string, value int, rangeType RangeType) *Condition {
	return condition.OrWhereRange(field, value, rangeType)
}

// OrWhereRange Should Range 单向范围查询, gt:>a, lt:<a等等
func (c *Condition) OrWhereRange(field string, value int, rangeType RangeType) *Condition {

	c.whereRange(Should, field, value, rangeType)

	return c
}

// OrWhereMatch Should match 匹配
func OrWhereMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Condition {
	return condition.OrWhereMatch(field, value, matchType, params)
}

// OrWhereMatch Should match 匹配
func (c *Condition) OrWhereMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Condition {
	c.whereMatch(Should, field, value, matchType, params)
	return c
}

// OrWhereMultiMatch Should multi_match 匹配
func OrWhereMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Condition {
	return condition.OrWhereMultiMatch(field, value, fieldType, params)
}

// OrWhereMultiMatch Should multi_match 匹配
func (c *Condition) OrWhereMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Condition {
	c.whereMultiMatch(Should, field, value, fieldType, params)

	return c
}

// OrWhereNested Should 嵌套查询, 例如嵌套 should 语句
func OrWhereNested(fn NestedFunc) *Condition {
	return condition.OrWhereNested(fn)
}

// OrWhereNested Should 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (c *Condition) OrWhereNested(fn NestedFunc) *Condition {
	c.nested[Should] = append(c.nested[Should], fn)

	return c
}

func MinimumShouldMatch(value int) *Condition {
	return condition.MinimumShouldMatch(value)
}

func (c *Condition) MinimumShouldMatch(value int) *Condition {
	c.minimumShouldMatch = value
	return c
}

// Filter Filter term 查询语句
func Filter(field string, value string) *Condition {
	return condition.Filter(field, value)
}

// Filter Filter term 查询语句
func (c *Condition) Filter(field string, value string) *Condition {
	c.termQuery(FilterClause, field, value)

	return c
}

// FilterIn Filter terms 查询语句
func FilterIn(field string, value []string) *Condition {
	return condition.WhereIn(field, value)
}

// FilterIn Filter terms 查询语句
func (c *Condition) FilterIn(field string, value []string) *Condition {
	c.termsQuery(FilterClause, field, value)

	return c
}

// FilterBetween Filter Range ( gte:>=a and lte:<=b) 查询语句
func FilterBetween(field string, value1, value2 int) *Condition {
	return condition.WhereBetween(field, value1, value2)
}

// FilterBetween Filter Range ( gte:>=a and lte:<=b) 查询语句
func (c *Condition) FilterBetween(field string, value1, value2 int) *Condition {

	c.whereBetween(FilterClause, field, value1, value2)

	return c
}

// FilterRange Filter Range 单向范围查询, gt:>a, lt:<a等等
func FilterRange(field string, value int, rangeType RangeType) *Condition {
	return condition.FilterRange(field, value, rangeType)
}

// FilterRange Filter Range 单向范围查询, gt:>a, lt:<a等等
func (c *Condition) FilterRange(field string, value int, rangeType RangeType) *Condition {

	c.whereRange(FilterClause, field, value, rangeType)

	return c
}

// FilterMatch Filter match 匹配
func FilterMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Condition {
	return condition.FilterMatch(field, value, matchType, params)
}

// FilterMatch Filter match 匹配
func (c *Condition) FilterMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Condition {
	c.whereMatch(FilterClause, field, value, matchType, params)
	return c
}

// FilterMultiMatch Filter multi_match 匹配
func FilterMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Condition {
	return condition.FilterMultiMatch(field, value, fieldType, params)
}

// FilterMultiMatch Filter multi_match 匹配
func (c *Condition) FilterMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Condition {
	c.whereMultiMatch(FilterClause, field, value, fieldType, params)

	return c
}

// FilterNested Filter 嵌套查询, 例如嵌套 should 语句
func FilterNested(fn NestedFunc) *Condition {
	return condition.FilterNested(fn)
}

// FilterNested Filter 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (c *Condition) FilterNested(fn NestedFunc) *Condition {
	c.nested[FilterClause] = append(c.nested[FilterClause], fn)

	return c
}

func (c *Condition) termQuery(clauseTyp BoolClauseType, field string, value string) {
	term := termlevel.TermQuery{
		Term: make(map[string]string),
	}

	term.Term[field] = value

	c.append(clauseTyp, term)
}

func (c *Condition) termsQuery(clauseTyp BoolClauseType, field string, value []string) {
	terms := termlevel.TermQuery{
		Terms: make(map[string][]string),
	}

	terms.Terms[field] = value

	c.append(clauseTyp, terms)
}

func (c *Condition) whereRange(clauseType BoolClauseType, field string, value int, rangeType RangeType) {
	rangeQuery := termlevel.RangeQuery{}

	switch rangeType {
	case Gt:
		rangeQuery.Gt = value
	case Gte:
		rangeQuery.Gte = value
	case Lt:
		rangeQuery.Lt = value
	case Lte:
		rangeQuery.Lte = value
	default:
		rangeQuery.Gt = value
	}

	c.rangeQuery(clauseType, field, rangeQuery)
}

func (c *Condition) whereBetween(clauseType BoolClauseType, field string, value1, value2 int) {
	rangeQuery := termlevel.RangeQuery{
		Gte: value1,
		Lte: value2,
	}

	c.rangeQuery(clauseType, field, rangeQuery)
}

func (c *Condition) rangeQuery(clauseType BoolClauseType, field string, rangeQuery termlevel.RangeQuery) {
	termQuery := termlevel.TermQuery{
		Range: make(map[string]termlevel.RangeQuery),
	}

	termQuery.Range[field] = rangeQuery

	c.append(clauseType, termQuery)
}

func (c *Condition) whereMatch(clauseType BoolClauseType, field string, value string, matchType MatchType, params fulltext.AppendParams) {
	textQuery := fulltext.TextQuery{
		Match:       make(map[string]fulltext.MatchQuery),
		MatchPhrase: make(map[string]fulltext.MatchQuery),
	}

	matchQuery := fulltext.MatchQuery{Query: value, AppendParams: params}

	switch matchType {
	case Match:
		textQuery.Match[field] = matchQuery
	case MatchPhrase:
		textQuery.MatchPhrase[field] = matchQuery
	case MatchPhrasePrefix:
	default:
		textQuery.Match[field] = matchQuery
	}

	c.append(clauseType, textQuery)
}

func (c *Condition) whereMultiMatch(clauseType BoolClauseType, field []string, value string, fieldType FieldType, params fulltext.AppendParams) {
	if field == nil || len(field) == 0 {
		return
	}

	typ, ok := FieldTypeSet[fieldType]
	if !ok {
		typ = FieldTypeSet[BestFields]
	}

	textQuery := fulltext.TextQuery{}

	multiMatch := &fulltext.MultiMatchQuery{
		Query:        value,
		Type:         typ,
		Fields:       field,
		AppendParams: params,
	}

	textQuery.MultiMatch = multiMatch
	c.append(clauseType, textQuery)
}

func (c *Condition) append(clauseTyp BoolClauseType, clause BoolBuilder) {
	c.where[clauseTyp] = append(c.where[clauseTyp], clause)
}

//func WhereBetween(field string, value []int, rangeType ...RangeType) *Condition {
//	return condition.WhereBetween(field, value, rangeType...)
//}
//
//func (c *Condition) WhereBetween(field string, value []int, rangeType ...RangeType) *Condition {
//	if value == nil {
//		panic("范围查询，参数value必须传值")
//	}
//
//	if len(value) == 1 && rangeType == nil {
//		panic("单范围查询，参数rangeType必须传值")
//	}
//
//	whereRange := termlevel.RangeQuery{}
//
//	if len(value) == 1 {
//
//		switch rangeType[0] {
//		case Gt:
//			whereRange.Gt = value[0]
//		case Gte:
//			whereRange.Gte = value[0]
//		case Lt:
//			whereRange.Lt = value[0]
//		case Lte:
//			whereRange.Lte = value[0]
//		}
//	} else {
//		whereRange.Gte = value[0]
//		whereRange.Lte = value[1]
//	}
//
//	termQuery := termlevel.TermQuery{
//		Range: make(map[string]termlevel.RangeQuery),
//	}
//
//	termQuery.Range[field] = whereRange
//	c.where[Must] = append(c.where[Must], termQuery)
//
//	return c
//}
