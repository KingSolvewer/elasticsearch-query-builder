package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/es"
)

func (b *Builder) compile() {

	b.Dsl = &es.Dsl{
		Source: make([]string, 0),
		Sort:   make([]es.Sort, 0),
		Query:  make(map[string]es.QueryBuilder),
	}

	if b.fields != nil || len(b.fields) > 0 {
		b.Dsl.Source = b.fields
	}
	if b.sort != nil || len(b.sort) > 0 {
		b.Dsl.Sort = b.sort
	}

	if b.size >= 0 {
		b.Dsl.Size = es.Uint(b.size)
	} else {
		b.Dsl.Size = es.Uint(10)
	}

	if b.page > 0 {
		b.Dsl.From = es.Uint((b.page - 1) * b.size)
	}

	boolQuery := b.component()

	if len(boolQuery.Should) > 0 {
		boolQuery.MinimumShouldMatch = b.minimumShouldMatch
	}

	b.Dsl.Query["bool"] = boolQuery
	b.Dsl.Aggs = b.aggs
}

func (b *Builder) component() es.BoolQuery {
	boolQuery := es.BoolQuery{}

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

			newQuery := make(es.Query)
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
