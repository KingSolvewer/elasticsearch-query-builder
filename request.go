package elastic

import (
	"encoding/json"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
)

func (b *Builder) Get() (*Result, error) {
	data, err := b.runQuery()
	if err != nil {
		return nil, err
	}

	b.Response = &Response{}

	err = json.Unmarshal(data, b.Response)
	if err != nil {
		return nil, err
	}

	b.Result = &Result{
		Total: b.Response.Hits.Total,
		List:  make([]map[string]any, 0),
	}

	if b.Response.ScrollId != "" {
		b.Result.ScrollId = b.Response.ScrollId
	}
	if b.Response.Aggregations != nil {
		b.Result.Aggs = b.Response.Aggregations
	}

	b.Request.Decorate()

	return b.Result, err
}

func (b *Builder) Paginator(page, size uint) (*Result, error) {
	var from uint = 0
	if page > 0 {
		from = page - 1
	}
	b.size = size
	b.from = from

	if b.collapse.Field != "" {
		b.Cardinality(b.collapse.Field, nil)
	}

	data, err := b.runQuery()
	if err != nil {
		return nil, err
	}

	b.Response = &Response{}

	err = json.Unmarshal(data, b.Response)
	if err != nil {
		return nil, err
	}

	b.Result = &Result{
		Total:         b.Response.Hits.Total,
		List:          make([]map[string]any, 0),
		OriginalTotal: esearch.Uint(b.Response.Hits.Total),
		PerPage:       esearch.Uint(size),
		CurrentPage:   esearch.Uint(page),
	}

	if b.Response.Aggregations != nil {
		b.Result.Aggs = b.Response.Aggregations
	}

	if b.collapse.Field != "" {
		b.Result.Total = b.Result.Aggs[b.collapse.Field+"_"+esearch.Cardinality].(*CardinalityAggResult).Value
	}

	b.Result.LastPage = esearch.Uint(uint(b.Result.Total)/size + 1)

	b.Request.Decorate()

	return b.Result, err
}

func (b *Builder) runQuery() (data []byte, err error) {
	if b.Request == nil {
		panic("请先实现elastic.Request接口")
	}

	if b.scroll == "" {
		data, err = b.Request.Query()
	} else {
		data, err = b.Request.ScrollQuery()
	}

	return data, err
}

func (b *Builder) Decorate() {
	for _, item := range b.Response.Hits.Hits {
		b.Result.List = append(b.Result.List, item.Source)
	}
}
