package elastic

type QueryBuilder interface {
	QueryBuild() string
}

type BoolBuilder interface {
	BoolBuild() string
}

type Sort map[string]Order

type Order struct {
	Order string
}

type Paginator interface {
	Page() int
}

type Int int

func (i Int) Page() int {
	return int(i)
}

type Dsl struct {
	Source []string  `json:"_source,omitempty"`
	Size   Paginator `json:"size,omitempty"`
	From   Paginator `json:"from,omitempty"`
	Sort   []Sort    `json:"sort,omitempty"`
	Query  `json:"query,omitempty"`
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

func (c *Condition) compile() {

	c.Dsl = &Dsl{
		Query: make(map[string]QueryBuilder),
	}

	boolQuery := c.component()

	if len(boolQuery.Should) > 0 {
		boolQuery.MinimumShouldMatch = c.minimumShouldMatch
	}

	c.Dsl.Query["bool"] = boolQuery
}

func (c *Condition) component() BoolQuery {
	boolQuery := BoolQuery{}

	for key, items := range c.where {
		switch key {
		case Must:
			boolQuery.Must = append(boolQuery.Must, items...)
		case MustNot:
			boolQuery.MustNot = append(boolQuery.MustNot, items...)
		case Should:
			boolQuery.Should = append(boolQuery.Should, items...)
		case FilterClause:
			boolQuery.Filter = append(boolQuery.Filter, items...)
		}
	}

	for key, fns := range c.nested {
		for _, fn := range fns {
			newCondition := NewCondition()
			newCondition = fn(newCondition)
			newBoolQuery := newCondition.component()

			newQuery := make(Query)
			newQuery["bool"] = newBoolQuery

			switch key {
			case Must:
				boolQuery.Must = append(boolQuery.Must, newQuery)
			case MustNot:
				boolQuery.MustNot = append(boolQuery.MustNot, newQuery)
			case Should:
				boolQuery.Should = append(boolQuery.Should, newQuery)
			case FilterClause:
				boolQuery.Filter = append(boolQuery.Filter, newQuery)
			}
		}
	}

	return boolQuery
}
