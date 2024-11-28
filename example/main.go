package main

import (
	elastic "github.com/KingSolvewer/elasticsearch-query-builder"
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"log"
	"main/zyzx"
)

type Dsl struct {
	Term string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	es := elastic.NewBuilder()
	//es.GroupBy(zyzx.Platform, aggs.TermsParam{}, func(b *elastic.Builder) {
	//	b.GroupBy(zyzx.PostDate, aggs.TermsParam{})
	//	b.GroupBy(zyzx.SimHash, aggs.TermsParam{})
	//	//b.TopHitsFunc(func(b *elastic.Builder) {
	//	//	b.Size(1).OrderBy(zyzx.PostTime, esearch.Asc).Select(zyzx.NewsUuid, zyzx.NewsTitle, zyzx.NewsContent, zyzx.PostTime)
	//	//})
	//})
	//es.DateGroupBy(zyzx.PostTime, aggs.HistogramParam{Interval: "day", Format: "yyyy-MM-dd"}, func(b *elastic.Builder) {
	//	b.GroupBy(zyzx.SimHash, aggs.TermsParam{})
	//})

	//es.Range(zyzx.NewsPosition, aggs.RangeParam{Ranges: []aggs.Ranges{{To: 1}, {From: 1, To: 2}, {From: 2}}})
	//es.Cardinality(zyzx.SimHash, nil)
	es.TopHits(aggs.TopHitsParam{Size: 10})

	dsl := es.Size(1).Select(zyzx.NewsUuid, zyzx.NewsTitle, zyzx.NewsContent, zyzx.PostTime).Dsl()
	log.Fatalln(dsl)
}
