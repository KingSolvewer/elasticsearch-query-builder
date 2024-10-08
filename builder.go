package elastic

import (
	"encoding/json"
	"github.com/KingSolvewer/elasticsearch-query-builder/collapse"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
	"github.com/KingSolvewer/elasticsearch-query-builder/fulltext"
	"github.com/KingSolvewer/elasticsearch-query-builder/termlevel"
)

type NestedFunc func(b *Builder)

type Builder struct {
	fields             []string
	size               uint
	from               uint
	sort               []esearch.Sort
	where              map[esearch.BoolClauseType][]esearch.BoolBuilder
	nested             map[esearch.BoolClauseType][]NestedFunc
	minimumShouldMatch int
	aggs               map[string]esearch.Aggregator
	query              *esearch.ElasticQuery
	scroll             string
	scrollId           string
	collapse           collapse.Collapse
	Request            esearch.Request
	*Result
	*Response
}

func NewBuilder() *Builder {
	return &Builder{
		fields: make([]string, 0),
		sort:   make([]esearch.Sort, 0),
		where:  make(map[esearch.BoolClauseType][]esearch.BoolBuilder),
		nested: make(map[esearch.BoolClauseType][]NestedFunc),
		aggs:   make(map[string]esearch.Aggregator),
	}
}

func (b *Builder) Dsl() string {
	bytes, _ := b.Marshal()

	return string(bytes)
}

func (b *Builder) Marshal() ([]byte, error) {
	b.compile()

	return json.Marshal(b.query)
}

func (b *Builder) Reset() *Builder {
	b.fields = make([]string, 0)
	b.size = 0
	b.from = 0
	b.sort = make([]esearch.Sort, 0)
	b.where = make(map[esearch.BoolClauseType][]esearch.BoolBuilder)
	b.nested = make(map[esearch.BoolClauseType][]NestedFunc)
	b.minimumShouldMatch = 0
	b.query = nil

	return b
}

func (b *Builder) Clone() *Builder {
	where := make(map[esearch.BoolClauseType][]esearch.BoolBuilder)
	for key, boolBuilderSet := range b.where {
		where[key] = append(where[key], boolBuilderSet...)
	}

	nested := make(map[esearch.BoolClauseType][]NestedFunc)
	for key, nestedFuncSet := range b.nested {
		nested[key] = append(nested[key], nestedFuncSet...)
	}

	return &Builder{
		fields:             b.fields,
		where:              where,
		nested:             nested,
		minimumShouldMatch: b.minimumShouldMatch,
	}
}

func (b *Builder) Scroll(scroll string) *Builder {
	b.scroll = scroll
	return b
}

func (b *Builder) GetScroll() string {
	return b.scroll
}

func (b *Builder) ScrollId(scrollId string) *Builder {
	b.scrollId = scrollId
	return b
}

func (b *Builder) GetScrollId() string {
	return b.scrollId
}

func (b *Builder) Collapse(field string) *Builder {
	b.collapse = collapse.Collapse{
		Field: field,
	}
	return b
}

func (b *Builder) Select(fields ...string) *Builder {
	b.fields = fields
	return b
}

func (b *Builder) AppendField(fields ...string) *Builder {
	if len(fields) > 0 {
		b.fields = append(b.fields, fields...)
	} else {
		b.fields = fields
	}
	return b
}

func (b *Builder) Size(value uint) *Builder {
	b.size = value

	return b
}

func (b *Builder) From(value uint) *Builder {
	b.from = value
	return b
}

func (b *Builder) OrderBy(field string, orderType esearch.OrderType) *Builder {

	sort := make(map[string]esearch.Order)

	switch orderType {
	case esearch.Asc, esearch.Desc:
		sort[field] = esearch.Order{Order: orderType}
	default:
		sort[field] = esearch.Order{Order: esearch.Asc}
	}

	b.sort = append(b.sort, sort)

	return b
}

// Where Must term 查询语句
func (b *Builder) Where(field string, value any) *Builder {
	b.termQuery(esearch.Must, field, value)

	return b
}

// WhereExists Must exists 查询语句
func (b *Builder) WhereExists(field string) *Builder {
	b.exists(esearch.Must, field)

	return b
}

// WhereRegexp Must regexp 查询语句
func (b *Builder) WhereRegexp(field string, value string, fn termlevel.RegexpParamFunc) *Builder {
	b.regexp(esearch.Must, field, value, fn)

	return b
}

// WhereWildcard Must wildcard 查询语句
func (b *Builder) WhereWildcard(field string, value string, fn termlevel.WildcardParamFunc) *Builder {
	b.wildcard(esearch.Must, field, value, fn)

	return b
}

// WhereIn Must terms 查询语句
func (b *Builder) WhereIn(field string, value []any) *Builder {
	b.termsQuery(esearch.Must, field, value)

	return b
}

// WhereBetween Must Range ( gte:>=a and lte:<=b) 查询语句
func (b *Builder) WhereBetween(field string, value1, value2 any) *Builder {
	b.whereBetween(esearch.Must, field, value1, value2)

	return b
}

// WhereRange Must Range 单向范围查询, gt:>a, lt:<a等等
func (b *Builder) WhereRange(field string, value any, rangeType esearch.RangeType) *Builder {

	b.whereRange(esearch.Must, field, value, rangeType)

	return b
}

// WhereMatch Must match 匹配
func (b *Builder) WhereMatch(field string, value string, matchType esearch.MatchType, fn fulltext.AppendParamsFunc) *Builder {
	b.whereMatch(esearch.Must, field, value, matchType, fn)
	return b
}

// WhereMultiMatch Must multi_match 匹配
func (b *Builder) WhereMultiMatch(field []string, value string, fieldType esearch.FieldType, fn fulltext.AppendParamsFunc) *Builder {
	b.whereMultiMatch(esearch.Must, field, value, fieldType, fn)

	return b
}

// WhereNested Must 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (b *Builder) WhereNested(fn NestedFunc) *Builder {
	b.nested[esearch.Must] = append(b.nested[esearch.Must], fn)

	return b
}

// WhereNot MustNot term 查询
func (b *Builder) WhereNot(field string, value string) *Builder {
	b.termQuery(esearch.MustNot, field, value)

	return b
}

// WhereNotExistsNot MustNot exists 查询语句
func (b *Builder) WhereNotExistsNot(field string) *Builder {
	b.exists(esearch.MustNot, field)

	return b
}

// WhereNotRegexp MustNot regexp 查询语句
func (b *Builder) WhereNotRegexp(field string, value string, fn termlevel.RegexpParamFunc) *Builder {
	b.regexp(esearch.MustNot, field, value, fn)

	return b
}

// WhereNotWildcard MustNot wildcard 查询语句
func (b *Builder) WhereNotWildcard(field string, value string, fn termlevel.WildcardParamFunc) *Builder {
	b.wildcard(esearch.MustNot, field, value, fn)

	return b
}

// WhereNotIn Must terms 查询语句
func (b *Builder) WhereNotIn(field string, value []any) *Builder {
	b.termsQuery(esearch.MustNot, field, value)

	return b
}

// WhereNotBetween MustNot Range ( gte:>=a and lte:<=b) 查询语句
func (b *Builder) WhereNotBetween(field string, value1, value2 int) *Builder {
	b.whereBetween(esearch.MustNot, field, value1, value2)

	return b
}

// WhereNotRange MustNot Range 单向范围查询, gt:>a, lt:<a等等
func (b *Builder) WhereNotRange(field string, value int, rangeType esearch.RangeType) *Builder {
	b.whereRange(esearch.MustNot, field, value, rangeType)

	return b
}

// WhereNotMatch MustNot match 匹配
func (b *Builder) WhereNotMatch(field string, value string, matchType esearch.MatchType, fn fulltext.AppendParamsFunc) *Builder {
	b.whereMatch(esearch.MustNot, field, value, matchType, fn)
	return b
}

// WhereNotMultiMatch MustNot multi_match 匹配
func (b *Builder) WhereNotMultiMatch(field []string, value string, fieldType esearch.FieldType, fn fulltext.AppendParamsFunc) *Builder {
	b.whereMultiMatch(esearch.MustNot, field, value, fieldType, fn)

	return b
}

// WhereNotNested Must 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (b *Builder) WhereNotNested(fn NestedFunc) *Builder {
	b.nested[esearch.MustNot] = append(b.nested[esearch.MustNot], fn)

	return b
}

// OrWhere Should term 查询语句
func (b *Builder) OrWhere(field string, value any) *Builder {
	b.termQuery(esearch.Should, field, value)

	return b
}

// OrWhereExists Should exists 查询语句
func (b *Builder) OrWhereExists(field string) *Builder {
	b.exists(esearch.Should, field)

	return b
}

// OrWhereRegexp Should regexp 查询语句
func (b *Builder) OrWhereRegexp(field string, value string, fn termlevel.RegexpParamFunc) *Builder {
	b.regexp(esearch.Should, field, value, fn)

	return b
}

// OrWhereWildcard Should wildcard 查询语句
func (b *Builder) OrWhereWildcard(field string, value string, fn termlevel.WildcardParamFunc) *Builder {
	b.wildcard(esearch.Should, field, value, fn)

	return b
}

// OrWhereIn Should terms 查询语句
func (b *Builder) OrWhereIn(field string, value []any) *Builder {
	b.termsQuery(esearch.Should, field, value)

	return b
}

// OrWhereBetween Should Range ( gte:>=a and lte:<=b) 查询语句
func (b *Builder) OrWhereBetween(field string, value1, value2 int) *Builder {

	b.whereBetween(esearch.Should, field, value1, value2)

	return b
}

// OrWhereRange Should Range 单向范围查询, gt:>a, lt:<a等等
func (b *Builder) OrWhereRange(field string, value int, rangeType esearch.RangeType) *Builder {

	b.whereRange(esearch.Should, field, value, rangeType)

	return b
}

// OrWhereMatch Should match 匹配
func (b *Builder) OrWhereMatch(field string, value string, matchType esearch.MatchType, fn fulltext.AppendParamsFunc) *Builder {
	b.whereMatch(esearch.Should, field, value, matchType, fn)
	return b
}

// OrWhereMultiMatch Should multi_match 匹配
func (b *Builder) OrWhereMultiMatch(field []string, value string, fieldType esearch.FieldType, fn fulltext.AppendParamsFunc) *Builder {
	b.whereMultiMatch(esearch.Should, field, value, fieldType, fn)

	return b
}

// OrWhereNested Should 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (b *Builder) OrWhereNested(fn NestedFunc) *Builder {
	b.nested[esearch.Should] = append(b.nested[esearch.Should], fn)

	return b
}

func (b *Builder) MinimumShouldMatch(value int) *Builder {
	b.minimumShouldMatch = value
	return b
}

// Filter Filter term 查询语句
func (b *Builder) Filter(field string, value any) *Builder {
	b.termQuery(esearch.FilterClause, field, value)

	return b
}

// FilterExists Filter exists 查询语句
func (b *Builder) FilterExists(field string) *Builder {
	b.exists(esearch.FilterClause, field)

	return b
}

// FilterRegexp Filter exists 查询语句
func (b *Builder) FilterRegexp(field string, value string, fn termlevel.RegexpParamFunc) *Builder {
	b.regexp(esearch.FilterClause, field, value, fn)

	return b
}

// FilterWildcard Filter wildcard 查询语句
func (b *Builder) FilterWildcard(field string, value string, fn termlevel.WildcardParamFunc) *Builder {
	b.wildcard(esearch.FilterClause, field, value, fn)

	return b
}

// FilterIn Filter terms 查询语句
func (b *Builder) FilterIn(field string, value []any) *Builder {
	b.termsQuery(esearch.FilterClause, field, value)

	return b
}

// FilterBetween Filter Range ( gte:>=a and lte:<=b) 查询语句
func (b *Builder) FilterBetween(field string, value1, value2 int) *Builder {

	b.whereBetween(esearch.FilterClause, field, value1, value2)

	return b
}

// FilterRange Filter Range 单向范围查询, gt:>a, lt:<a等等
func (b *Builder) FilterRange(field string, value int, rangeType esearch.RangeType) *Builder {

	b.whereRange(esearch.FilterClause, field, value, rangeType)

	return b
}

// FilterMatch Filter match 匹配
func (b *Builder) FilterMatch(field string, value string, matchType esearch.MatchType, fn fulltext.AppendParamsFunc) *Builder {
	b.whereMatch(esearch.FilterClause, field, value, matchType, fn)
	return b
}

// FilterMultiMatch Filter multi_match 匹配
func (b *Builder) FilterMultiMatch(field []string, value string, fieldType esearch.FieldType, fn fulltext.AppendParamsFunc) *Builder {
	b.whereMultiMatch(esearch.FilterClause, field, value, fieldType, fn)

	return b
}

// FilterNested Filter 嵌套查询, 例如在 must 语句中嵌套 should 语句
func (b *Builder) FilterNested(fn NestedFunc) *Builder {
	b.nested[esearch.FilterClause] = append(b.nested[esearch.FilterClause], fn)

	return b
}

func (b *Builder) termQuery(clauseTyp esearch.BoolClauseType, field string, value any) {
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

func (b *Builder) termsQuery(clauseTyp esearch.BoolClauseType, field string, value []any) {
	terms := termlevel.TermQuery{
		Terms: make(map[string][]any),
	}

	terms.Terms[field] = value

	b.append(clauseTyp, terms)
}

func (b *Builder) exists(clauseTyp esearch.BoolClauseType, field string) {
	term := termlevel.TermQuery{
		Exists: make(map[string]string),
	}

	term.Exists["field"] = field
	b.append(clauseTyp, term)
}

func (b *Builder) regexp(clauseTyp esearch.BoolClauseType, field string, value string, fn termlevel.RegexpParamFunc) {
	term := termlevel.TermQuery{
		Regexp: make(map[string]termlevel.Regexp),
	}

	regexp := termlevel.Regexp{Value: value}
	if fn != nil {
		regexp.RegexpParam = fn()
	}

	term.Regexp[field] = regexp

	b.append(clauseTyp, term)
}

func (b *Builder) wildcard(clauseTyp esearch.BoolClauseType, field string, value string, fn termlevel.WildcardParamFunc) {
	term := termlevel.TermQuery{
		Wildcard: make(map[string]termlevel.Wildcard),
	}

	wildcard := termlevel.Wildcard{Value: value}
	if fn != nil {
		wildcard.WildcardParam = fn()
	}
	term.Wildcard[field] = wildcard

	b.append(clauseTyp, term)
}

func (b *Builder) whereRange(clauseTyp esearch.BoolClauseType, field string, value any, rangeType esearch.RangeType) {
	if !checkType(value) {
		return
	}

	rangeQuery := termlevel.RangeQuery{}

	switch rangeType {
	case esearch.Gt:
		rangeQuery.Gt = value
	case esearch.Gte:
		rangeQuery.Gte = value
	case esearch.Lt:
		rangeQuery.Lt = value
	case esearch.Lte:
		rangeQuery.Lte = value
	default:
		rangeQuery.Gt = value
	}

	b.rangeQuery(clauseTyp, field, rangeQuery)
}

func (b *Builder) whereBetween(clauseTyp esearch.BoolClauseType, field string, value1, value2 any) {
	if !checkType(value1) || !checkType(value2) {
		return
	}

	rangeQuery := termlevel.RangeQuery{
		Gte: value1,
		Lte: value2,
	}

	b.rangeQuery(clauseTyp, field, rangeQuery)
}

func (b *Builder) rangeQuery(clauseTyp esearch.BoolClauseType, field string, rangeQuery termlevel.RangeQuery) {
	termQuery := termlevel.TermQuery{
		Range: make(map[string]termlevel.RangeQuery),
	}

	termQuery.Range[field] = rangeQuery

	b.append(clauseTyp, termQuery)
}

func (b *Builder) whereMatch(clauseTyp esearch.BoolClauseType, field string, value string, matchType esearch.MatchType, fn func() fulltext.AppendParams) {
	textQuery := fulltext.TextQuery{
		Match:       make(map[string]fulltext.MatchQuery),
		MatchPhrase: make(map[string]fulltext.MatchQuery),
	}

	matchQuery := fulltext.MatchQuery{Query: value}
	if fn != nil {
		matchQuery.AppendParams = fn()
	}

	switch matchType {
	case esearch.Match:
		textQuery.Match[field] = matchQuery
	case esearch.MatchPhrase:
		textQuery.MatchPhrase[field] = matchQuery
	case esearch.MatchPhrasePrefix:
	default:
		textQuery.Match[field] = matchQuery
	}

	b.append(clauseTyp, textQuery)
}

func (b *Builder) whereMultiMatch(clauseTyp esearch.BoolClauseType, field []string, value string, fieldType esearch.FieldType, fn fulltext.AppendParamsFunc) {
	if field == nil || len(field) == 0 {
		return
	}

	typ, ok := esearch.FieldTypeSet[fieldType]
	if !ok {
		typ = esearch.FieldTypeSet[esearch.BestFields]
	}

	textQuery := fulltext.TextQuery{}

	multiMatch := &fulltext.MultiMatchQuery{
		Query:  value,
		Type:   typ,
		Fields: field,
	}

	if fn != nil {
		multiMatch.AppendParams = fn()
	}

	textQuery.MultiMatch = multiMatch
	b.append(clauseTyp, textQuery)
}

func (b *Builder) append(clauseTyp esearch.BoolClauseType, clause esearch.BoolBuilder) {
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
