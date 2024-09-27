package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"github.com/KingSolvewer/elasticsearch-query-builder/es"
)

type QueryBuilder interface {
	QueryBuild() string
}

type BoolBuilder interface {
	BoolBuild() string
}

type Sort map[string]Order

type Order struct {
	Order es.OrderType `json:"order"`
}

type Paginator interface {
	Page() uint
}

type Uint uint

func (i Uint) Page() uint {
	return uint(i)
}

type Dsl struct {
	Source []string  `json:"_source,omitempty"`
	Size   Paginator `json:"size,omitempty"`
	From   Paginator `json:"from,omitempty"`
	Sort   []Sort    `json:"sort,omitempty"`
	Query  `json:"query,omitempty"`
	Aggs   map[string]map[string]aggs.Aggregator `json:"aggs,omitempty"`
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

func (b *Builder) compile() {

	b.Dsl = &Dsl{
		Source: make([]string, 0),
		Sort:   make([]Sort, 0),
		Query:  make(map[string]QueryBuilder),
	}

	if b.fields != nil || len(b.fields) > 0 {
		b.Dsl.Source = b.fields
	}
	if b.sort != nil || len(b.sort) > 0 {
		b.Dsl.Sort = b.sort
	}

	if b.size >= 0 {
		b.Dsl.Size = Uint(b.size)
	} else {
		b.Dsl.Size = Uint(10)
	}

	if b.page > 0 {
		b.Dsl.From = Uint((b.page - 1) * b.size)
	}

	boolQuery := b.component()

	if len(boolQuery.Should) > 0 {
		boolQuery.MinimumShouldMatch = b.minimumShouldMatch
	}

	b.Dsl.Query["bool"] = boolQuery
	b.Dsl.Aggs = b.aggs
}

func (b *Builder) component() BoolQuery {
	boolQuery := BoolQuery{}

	for key, items := range b.where {
		switch key {
		case es.Must:
			boolQuery.Must = append(boolQuery.Must, items...)
		case es.MustNot:
			boolQuery.MustNot = append(boolQuery.MustNot, items...)
		case es.Should:
			boolQuery.Should = append(boolQuery.Should, items...)
		case es.FilterClause:
			boolQuery.Filter = append(boolQuery.Filter, items...)
		}
	}

	for key, fns := range b.nested {
		for _, fn := range fns {
			newBuilder := NewBuilder()
			newBuilder = fn(newBuilder)
			newBoolQuery := newBuilder.component()

			newQuery := make(Query)
			newQuery["bool"] = newBoolQuery

			switch key {
			case es.Must:
				boolQuery.Must = append(boolQuery.Must, newQuery)
			case es.MustNot:
				boolQuery.MustNot = append(boolQuery.MustNot, newQuery)
			case es.Should:
				boolQuery.Should = append(boolQuery.Should, newQuery)
			case es.FilterClause:
				boolQuery.Filter = append(boolQuery.Filter, newQuery)
			}
		}
	}

	return boolQuery
}
