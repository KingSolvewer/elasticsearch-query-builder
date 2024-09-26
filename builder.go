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

type OrderType string

const (
	Asc  OrderType = "asc"
	Desc OrderType = "desc"
)

type NestedFunc func(b *Builder) *Builder

type Builder struct {
	fields             []string
	size               uint
	page               uint
	sort               []Sort
	where              map[BoolClauseType][]BoolBuilder
	nested             map[BoolClauseType][]NestedFunc
	minimumShouldMatch int
	*Dsl
}

var builder = NewBuilder()

func NewBuilder() *Builder {
	return &Builder{
		fields: make([]string, 0),
		sort:   make([]Sort, 0),
		where:  make(map[BoolClauseType][]BoolBuilder),
		nested: make(map[BoolClauseType][]NestedFunc),
	}
}

func GetCondition() *Builder {
	builder.compile()
	return builder
}

func Select(fields ...string) *Builder {
	return builder.Select(fields...)
}

func (b *Builder) Select(fields ...string) *Builder {
	if len(fields) > 0 {
		b.fields = append(b.fields, fields...)
	} else {
		b.fields = fields
	}
	return b
}

func AppendField(fields ...string) *Builder {
	return builder.Select(fields...)
}

func (b *Builder) AppendField(fields ...string) *Builder {
	return b.Select(fields...)
}

func Size(value uint) *Builder {
	return builder.Size(value)
}

func (b *Builder) Size(value uint) *Builder {
	b.size = value

	return b
}

func Page(value uint) *Builder {
	return builder.Page(value)
}

func (b *Builder) Page(value uint) *Builder {
	if value < 1 {
		value = 1
	}

	b.page = value
	return b
}

func Order(field string, order OrderType) *Builder {
	return builder.Order(field, order)
}

func (b *Builder) Order(field string, orderType OrderType) *Builder {

	sort := make(map[string]OrderBy)

	switch orderType {
	case Asc, Desc:
		sort[field] = OrderBy{Order: orderType}
	default:
		sort[field] = OrderBy{Order: Asc}
	}

	b.sort = append(b.sort, sort)

	return b
}

// Where Must term 查询语句
func Where(field string, value any) *Builder {
	return builder.Where(field, value)
}

// Where Must term 查询语句
func (b *Builder) Where(field string, value any) *Builder {
	b.termQuery(Must, field, value)

	return b
}

// WhereExists Must exists 查询语句
func WhereExists(field string) *Builder {
	return builder.WhereExists(field)
}

// WhereExists Must exists 查询语句
func (b *Builder) WhereExists(field string) *Builder {
	b.exists(Must, field)

	return b
}

// WhereRegexp Must regexp 查询语句
func WhereRegexp(field string, value string, params termlevel.RegexpParam) *Builder {
	return builder.WhereRegexp(field, value, params)
}

// WhereRegexp Must regexp 查询语句
func (b *Builder) WhereRegexp(field string, value string, params termlevel.RegexpParam) *Builder {
	b.regexp(Must, field, value, params)

	return b
}

// WhereWildcard Must wildcard 查询语句
func WhereWildcard(field string, value string, params termlevel.WildcardParam) *Builder {
	return builder.WhereWildcard(field, value, params)
}

// WhereWildcard Must wildcard 查询语句
func (b *Builder) WhereWildcard(field string, value string, params termlevel.WildcardParam) *Builder {
	b.wildcard(Must, field, value, params)

	return b
}

// WhereIn Must terms 查询语句
func WhereIn(field string, value []any) *Builder {
	return builder.WhereIn(field, value)
}

// WhereIn Must terms 查询语句
func (b *Builder) WhereIn(field string, value []any) *Builder {
	b.termsQuery(Must, field, value)

	return b
}

// WhereBetween Must Range ( gte:>=a and lte:<=b) 查询语句
func WhereBetween(field string, value1, value2 int) *Builder {
	return builder.WhereBetween(field, value1, value2)
}

// WhereBetween Must Range ( gte:>=a and lte:<=b) 查询语句
func (b *Builder) WhereBetween(field string, value1, value2 any) *Builder {
	b.whereBetween(Must, field, value1, value2)

	return b
}

// WhereRange Must Range 单向范围查询, gt:>a, lt:<a等等
func WhereRange(field string, value int, rangeType RangeType) *Builder {
	return builder.WhereRange(field, value, rangeType)
}

// WhereRange Must Range 单向范围查询, gt:>a, lt:<a等等
func (b *Builder) WhereRange(field string, value any, rangeType RangeType) *Builder {

	b.whereRange(Must, field, value, rangeType)

	return b
}

// WhereMatch Must match 匹配
func WhereMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Builder {
	return builder.WhereMatch(field, value, matchType, params)
}

// WhereMatch Must match 匹配
func (b *Builder) WhereMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Builder {
	b.whereMatch(Must, field, value, matchType, params)
	return b
}

// WhereMultiMatch Must multi_match 匹配
func WhereMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Builder {
	return builder.WhereMultiMatch(field, value, fieldType, params)
}

// WhereMultiMatch Must multi_match 匹配
func (b *Builder) WhereMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Builder {
	b.whereMultiMatch(Must, field, value, fieldType, params)

	return b
}

// WhereNested Must 嵌套查询, 例如嵌套 should 语句
func WhereNested(fn NestedFunc) *Builder {
	return builder.WhereNested(fn)
}

// WhereNested Must 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (b *Builder) WhereNested(fn NestedFunc) *Builder {
	b.nested[Must] = append(b.nested[Must], fn)

	return b
}

// WhereNot MustNot term 查询
func WhereNot(field string, value string) *Builder {
	return builder.WhereNot(field, value)
}

// WhereNot MustNot term 查询
func (b *Builder) WhereNot(field string, value string) *Builder {
	b.termQuery(MustNot, field, value)

	return b
}

// WhereNotExistsNot MustNot exists 查询语句
func WhereNotExistsNot(field string) *Builder {
	return builder.WhereNotExistsNot(field)
}

// WhereNotExistsNot MustNot exists 查询语句
func (b *Builder) WhereNotExistsNot(field string) *Builder {
	b.exists(MustNot, field)

	return b
}

// WhereNotRegexp MustNot regexp 查询语句
func WhereNotRegexp(field string, value string, params termlevel.RegexpParam) *Builder {
	return builder.WhereNotRegexp(field, value, params)
}

// WhereNotRegexp MustNot regexp 查询语句
func (b *Builder) WhereNotRegexp(field string, value string, params termlevel.RegexpParam) *Builder {
	b.regexp(MustNot, field, value, params)

	return b
}

// WhereNotWildcard MustNot wildcard 查询语句
func WhereNotWildcard(field string, value string, params termlevel.WildcardParam) *Builder {
	return builder.WhereNotWildcard(field, value, params)
}

// WhereNotWildcard MustNot wildcard 查询语句
func (b *Builder) WhereNotWildcard(field string, value string, params termlevel.WildcardParam) *Builder {
	b.wildcard(MustNot, field, value, params)

	return b
}

// WhereNotIn  MustNot terms 查询语句
func WhereNotIn(field string, value []any) *Builder {
	return builder.WhereIn(field, value)
}

// WhereNotIn Must terms 查询语句
func (b *Builder) WhereNotIn(field string, value []any) *Builder {
	b.termsQuery(MustNot, field, value)

	return b
}

// WhereNotBetween MustNot Range ( gte:>=a and lte:<=b) 查询语句
func WhereNotBetween(field string, value1, value2 int) *Builder {
	return builder.WhereNotBetween(field, value1, value2)
}

// WhereNotBetween MustNot Range ( gte:>=a and lte:<=b) 查询语句
func (b *Builder) WhereNotBetween(field string, value1, value2 int) *Builder {
	b.whereBetween(MustNot, field, value1, value2)

	return b
}

// WhereNotRange MustNot Range 单向范围查询, gt:>a, lt:<a等等
func WhereNotRange(field string, value int, rangeType RangeType) *Builder {
	return builder.WhereNotRange(field, value, rangeType)
}

// WhereNotRange MustNot Range 单向范围查询, gt:>a, lt:<a等等
func (b *Builder) WhereNotRange(field string, value int, rangeType RangeType) *Builder {
	b.whereRange(MustNot, field, value, rangeType)

	return b
}

// WhereNotMatch MustNot match 匹配
func WhereNotMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Builder {
	return builder.WhereNotMatch(field, value, matchType, params)
}

// WhereNotMatch MustNot match 匹配
func (b *Builder) WhereNotMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Builder {
	b.whereMatch(MustNot, field, value, matchType, params)
	return b
}

// WhereNotMultiMatch MustNot multi_match 匹配
func WhereNotMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Builder {
	return builder.WhereNotMultiMatch(field, value, fieldType, params)
}

// WhereNotMultiMatch MustNot multi_match 匹配
func (b *Builder) WhereNotMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Builder {
	b.whereMultiMatch(MustNot, field, value, fieldType, params)

	return b
}

// WhereNotNested MustNot 嵌套查询, 例如嵌套 should 语句
func WhereNotNested(fn NestedFunc) *Builder {
	return builder.WhereNotNested(fn)
}

// WhereNotNested Must 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (b *Builder) WhereNotNested(fn NestedFunc) *Builder {
	b.nested[MustNot] = append(b.nested[MustNot], fn)

	return b
}

// OrWhere Should term 查询语句
func OrWhere(field string, value any) *Builder {
	return builder.OrWhere(field, value)
}

// OrWhere Should term 查询语句
func (b *Builder) OrWhere(field string, value any) *Builder {
	b.termQuery(Should, field, value)

	return b
}

// OrWhereExists Should exists 查询语句
func OrWhereExists(field string) *Builder {
	return builder.OrWhereExists(field)
}

// OrWhereExists Should exists 查询语句
func (b *Builder) OrWhereExists(field string) *Builder {
	b.exists(Should, field)

	return b
}

// OrWhereRegexp Should regexp 查询语句
func OrWhereRegexp(field string, value string, params termlevel.RegexpParam) *Builder {
	return builder.OrWhereRegexp(field, value, params)
}

// OrWhereRegexp Should regexp 查询语句
func (b *Builder) OrWhereRegexp(field string, value string, params termlevel.RegexpParam) *Builder {
	b.regexp(Should, field, value, params)

	return b
}

// OrWhereWildcard Should wildcard 查询语句
func OrWhereWildcard(field string, value string, params termlevel.WildcardParam) *Builder {
	return builder.OrWhereWildcard(field, value, params)
}

// OrWhereWildcard Should wildcard 查询语句
func (b *Builder) OrWhereWildcard(field string, value string, params termlevel.WildcardParam) *Builder {
	b.wildcard(Should, field, value, params)

	return b
}

// OrWhereIn Should terms 查询语句
func OrWhereIn(field string, value []any) *Builder {
	return builder.OrWhereIn(field, value)
}

// OrWhereIn Should terms 查询语句
func (b *Builder) OrWhereIn(field string, value []any) *Builder {
	b.termsQuery(Should, field, value)

	return b
}

// OrWhereBetween Should Range ( gte:>=a and lte:<=b) 查询语句
func OrWhereBetween(field string, value1, value2 int) *Builder {
	return builder.OrWhereBetween(field, value1, value2)
}

// OrWhereBetween Should Range ( gte:>=a and lte:<=b) 查询语句
func (b *Builder) OrWhereBetween(field string, value1, value2 int) *Builder {

	b.whereBetween(Should, field, value1, value2)

	return b
}

// OrWhereRange Should Range 单向范围查询, gt:>a, lt:<a等等
func OrWhereRange(field string, value int, rangeType RangeType) *Builder {
	return builder.OrWhereRange(field, value, rangeType)
}

// OrWhereRange Should Range 单向范围查询, gt:>a, lt:<a等等
func (b *Builder) OrWhereRange(field string, value int, rangeType RangeType) *Builder {

	b.whereRange(Should, field, value, rangeType)

	return b
}

// OrWhereMatch Should match 匹配
func OrWhereMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Builder {
	return builder.OrWhereMatch(field, value, matchType, params)
}

// OrWhereMatch Should match 匹配
func (b *Builder) OrWhereMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Builder {
	b.whereMatch(Should, field, value, matchType, params)
	return b
}

// OrWhereMultiMatch Should multi_match 匹配
func OrWhereMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Builder {
	return builder.OrWhereMultiMatch(field, value, fieldType, params)
}

// OrWhereMultiMatch Should multi_match 匹配
func (b *Builder) OrWhereMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Builder {
	b.whereMultiMatch(Should, field, value, fieldType, params)

	return b
}

// OrWhereNested Should 嵌套查询, 例如嵌套 should 语句
func OrWhereNested(fn NestedFunc) *Builder {
	return builder.OrWhereNested(fn)
}

// OrWhereNested Should 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (b *Builder) OrWhereNested(fn NestedFunc) *Builder {
	b.nested[Should] = append(b.nested[Should], fn)

	return b
}

func MinimumShouldMatch(value int) *Builder {
	return builder.MinimumShouldMatch(value)
}

func (b *Builder) MinimumShouldMatch(value int) *Builder {
	b.minimumShouldMatch = value
	return b
}

// Filter Filter term 查询语句
func Filter(field string, value any) *Builder {
	return builder.Filter(field, value)
}

// Filter Filter term 查询语句
func (b *Builder) Filter(field string, value any) *Builder {
	b.termQuery(FilterClause, field, value)

	return b
}

// FilterExists Filter exists 查询语句
func FilterExists(field string) *Builder {
	return builder.FilterExists(field)
}

// FilterExists Filter exists 查询语句
func (b *Builder) FilterExists(field string) *Builder {
	b.exists(FilterClause, field)

	return b
}

// FilterRegexp Filter exists 查询语句
func FilterRegexp(field string, value string, params termlevel.RegexpParam) *Builder {
	return builder.FilterRegexp(field, value, params)
}

// FilterRegexp Filter exists 查询语句
func (b *Builder) FilterRegexp(field string, value string, params termlevel.RegexpParam) *Builder {
	b.regexp(FilterClause, field, value, params)

	return b
}

// FilterWildcard Filter wildcard 查询语句
func FilterWildcard(field string, value string, params termlevel.WildcardParam) *Builder {
	return builder.FilterWildcard(field, value, params)
}

// FilterWildcard Filter wildcard 查询语句
func (b *Builder) FilterWildcard(field string, value string, params termlevel.WildcardParam) *Builder {
	b.wildcard(FilterClause, field, value, params)

	return b
}

// FilterIn Filter terms 查询语句
func FilterIn(field string, value []any) *Builder {
	return builder.WhereIn(field, value)
}

// FilterIn Filter terms 查询语句
func (b *Builder) FilterIn(field string, value []any) *Builder {
	b.termsQuery(FilterClause, field, value)

	return b
}

// FilterBetween Filter Range ( gte:>=a and lte:<=b) 查询语句
func FilterBetween(field string, value1, value2 int) *Builder {
	return builder.WhereBetween(field, value1, value2)
}

// FilterBetween Filter Range ( gte:>=a and lte:<=b) 查询语句
func (b *Builder) FilterBetween(field string, value1, value2 int) *Builder {

	b.whereBetween(FilterClause, field, value1, value2)

	return b
}

// FilterRange Filter Range 单向范围查询, gt:>a, lt:<a等等
func FilterRange(field string, value int, rangeType RangeType) *Builder {
	return builder.FilterRange(field, value, rangeType)
}

// FilterRange Filter Range 单向范围查询, gt:>a, lt:<a等等
func (b *Builder) FilterRange(field string, value int, rangeType RangeType) *Builder {

	b.whereRange(FilterClause, field, value, rangeType)

	return b
}

// FilterMatch Filter match 匹配
func FilterMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Builder {
	return builder.FilterMatch(field, value, matchType, params)
}

// FilterMatch Filter match 匹配
func (b *Builder) FilterMatch(field string, value string, matchType MatchType, params fulltext.AppendParams) *Builder {
	b.whereMatch(FilterClause, field, value, matchType, params)
	return b
}

// FilterMultiMatch Filter multi_match 匹配
func FilterMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Builder {
	return builder.FilterMultiMatch(field, value, fieldType, params)
}

// FilterMultiMatch Filter multi_match 匹配
func (b *Builder) FilterMultiMatch(field []string, value string, fieldType FieldType, params fulltext.AppendParams) *Builder {
	b.whereMultiMatch(FilterClause, field, value, fieldType, params)

	return b
}

// FilterNested Filter 嵌套查询, 例如嵌套 should 语句
func FilterNested(fn NestedFunc) *Builder {
	return builder.FilterNested(fn)
}

// FilterNested Filter 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (b *Builder) FilterNested(fn NestedFunc) *Builder {
	b.nested[FilterClause] = append(b.nested[FilterClause], fn)

	return b
}

func (b *Builder) termQuery(clauseTyp BoolClauseType, field string, value any) {
	ok := checkType(value)
	if !ok {
		return
	}

	//if _, ok := value.(int);

	term := termlevel.TermQuery{
		Term: make(map[string]any),
	}

	term.Term[field] = value

	b.append(clauseTyp, term)
}

func (b *Builder) termsQuery(clauseTyp BoolClauseType, field string, value []any) {
	terms := termlevel.TermQuery{
		Terms: make(map[string][]any),
	}

	terms.Terms[field] = value

	b.append(clauseTyp, terms)
}

func (b *Builder) exists(clauseTyp BoolClauseType, field string) {
	term := termlevel.TermQuery{
		Exists: make(map[string]string),
	}

	term.Exists["field"] = field
	b.append(clauseTyp, term)
}

func (b *Builder) regexp(clauseType BoolClauseType, field string, value string, params termlevel.RegexpParam) {
	term := termlevel.TermQuery{
		Regexp: make(map[string]termlevel.Regexp),
	}

	term.Regexp[field] = termlevel.Regexp{Value: value, RegexpParam: params}

	b.append(clauseType, term)
}

func (b *Builder) wildcard(clauseType BoolClauseType, field string, value string, params termlevel.WildcardParam) {
	term := termlevel.TermQuery{
		Wildcard: make(map[string]termlevel.Wildcard),
	}

	term.Wildcard[field] = termlevel.Wildcard{Value: value, WildcardParam: params}

	b.append(clauseType, term)
}

func (b *Builder) whereRange(clauseType BoolClauseType, field string, value any, rangeType RangeType) {
	if !checkType(value) {
		return
	}

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

	b.rangeQuery(clauseType, field, rangeQuery)
}

func (b *Builder) whereBetween(clauseType BoolClauseType, field string, value1, value2 any) {
	if !checkType(value1) || !checkType(value2) {
		return
	}

	rangeQuery := termlevel.RangeQuery{
		Gte: value1,
		Lte: value2,
	}

	b.rangeQuery(clauseType, field, rangeQuery)
}

func (b *Builder) rangeQuery(clauseType BoolClauseType, field string, rangeQuery termlevel.RangeQuery) {
	termQuery := termlevel.TermQuery{
		Range: make(map[string]termlevel.RangeQuery),
	}

	termQuery.Range[field] = rangeQuery

	b.append(clauseType, termQuery)
}

func (b *Builder) whereMatch(clauseType BoolClauseType, field string, value string, matchType MatchType, params fulltext.AppendParams) {
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

	b.append(clauseType, textQuery)
}

func (b *Builder) whereMultiMatch(clauseType BoolClauseType, field []string, value string, fieldType FieldType, params fulltext.AppendParams) {
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
	b.append(clauseType, textQuery)
}

func (b *Builder) append(clauseTyp BoolClauseType, clause BoolBuilder) {
	b.where[clauseTyp] = append(b.where[clauseTyp], clause)
}

func checkType(value any) bool {
	switch value.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, string, bool:
		return true
	default:
		return false
	}
}

//func WhereBetween(field string, value []int, rangeType ...RangeType) *Condition {
//	return builder.WhereBetween(field, value, rangeType...)
//}
//
//func (b *Condition) WhereBetween(field string, value []int, rangeType ...RangeType) *Condition {
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
//	b.where[Must] = append(b.where[Must], termQuery)
//
//	return b
//}
