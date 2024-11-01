package main

import (
	"encoding/json"
	"fmt"
	elastic "github.com/KingSolvewer/elasticsearch-query-builder"
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"log"
	"main/ela"
)

type Dsl struct {
	Term string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	type Result struct {
		CategoryId  int64  `json:"category_id_"`
		ClaimUserId string `json:"claim_user_id_"`
	}

	m := make(map[string]any)
	var r Result
	var rs []Result

	es := ela.NewYingyanEs()

	//result, err := es.Select(ela.Title, ela.CategoryId, ela.PublishTime, ela.CreateTime).Where(ela.Stat, ela.StatWaitFilter).Size(1).Get()
	es.Select(ela.Title, ela.CategoryId, ela.PublishTime, ela.CreateTime, ela.ClaimUserId)
	es.Where(ela.Stat, ela.StatArchived).Size(10)
	//es.SearchTime("2024-02-01 00:00:00", "2024-02-05 00:00:00")
	es.WhereExists(ela.ClaimUserId)
	es.Collapse(ela.NewsSimHash)
	es.AggsFilter(ela.NewsSimHash, func(b *elastic.Builder) {
		b.Where(ela.Stat, ela.StatArchived).Where(ela.IsDelete, 0)
	}, func(b *elastic.Builder) {
		b.GroupBy(ela.CategoryId, aggs.TermsParam{})
	})
	//es.GroupBy(ela.NewsSimHash, aggs.TermsParam{Size: 10, Order: esearch.SortMap{"_key": esearch.Asc}}, func(b *elastic.Builder) {
	//	b.GroupBy(ela.Stat, aggs.TermsParam{}, func(b *elastic.Builder) {
	//		b.Histogram(ela.CreateTime, aggs.HistogramParam{})
	//	}, func(b *elastic.Builder) {
	//		b.TopHitsFunc(func(b *elastic.Builder) {
	//			b.From(2).Size(10).OrderBy("_count", esearch.Desc)
	//		})
	//	})
	//}, func(b *elastic.Builder) {
	//	b.GroupBy(ela.CategoryId, aggs.TermsParam{})
	//}, func(b *elastic.Builder) {
	//	b.TopHits(aggs.TopHitsParam{From: 1, Size: 10})
	//})
	//es.TopHits(aggs.TopHitsParam{})
	//es.TopHitsFunc(func(b *elastic.Builder) {
	//	b.From(0).Size(10).OrderBy(ela.NewsSimHash, esearch.Desc)
	//})
	//result, err := es.Get()
	//json.RawMessage{}

	log.Fatalln(es.Dsl())
	result, err := es.Paginator(1, 10, m)
	//result, err := es.Get(nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(es.Dsl(), m, r, rs)
	log.Println(result, err)
	jsonData, err := json.Marshal(result)
	log.Println(string(jsonData), err)
	return
	//dsl.SetSize(0)
	//dsl.SetFrom(10)
	//sort := make(dsl.Sort)
	//sort["posttime"] = dsl.Order{Order: "asc"}
	//dsl.SetSort(sort)
	//dsl.SetSource([]string{"posttime", "title"})
	//m := make(map[string]elastic.Builder)
	//m["asc"] = String("fsdff")
	//elastic.SetQuery(m)

	//elastic.Where("title", "1231312").Where("posttime", "12212").WhereIn("status", []string{"a", "b", "c"})
	//elastic.WhereIn("state", []string{"1", "2"}).WhereNot("create_time", "12313")
	//elastic.Where("title", "1231312").Where("posttime", "12212")
	//elastic.Where("title", "1231312").Where("posttime", "12212").WhereIn("status", []string{"a", "b", "c"})
	//elastic.WhereNot("title", "safjsdf").WhereNot("posttime", "sfasdfasf").WhereGt("stat_", 1)
	//elastic.Where("add_type", true)
	//elastic.Select("sdjf", "fsajlfas").Size(10).Page(2)
	//elastic.OrWhere("title", "saffdasdf").OrWhere("title", "fsdafasfd").WhereIn("posttime", []any{"1213", "fasfa", 42134})
	//elastic.WhereExists("title").WhereRegexp("title", "asfjsda", termlevel.RegexpParam{})
	//elastic.WhereWildcard("title", "safdsal", termlevel.WildcardParam{})
	//
	//elastic.WhereNested(func(c *elastic.Builder) *elastic.Builder {
	//	return c.OrWhere("create_time", "12312312").OrWhere("create_time", "456354").Where("create_time", "345253")
	//}).WhereNested(func(c *elastic.Builder) *elastic.Builder {
	//	return c.WhereNested(func(c *elastic.Builder) *elastic.Builder {
	//		return c.OrWhere("posttime", "123123123").OrWhere("posttime", "34432341")
	//	})
	//}).WhereRange("posttime", "9789978", elastic.Gt).WhereBetween("fsdafd", 1, 2).WhereMatch("title", "sfdasf", elastic.MatchPhrase, fulltext.AppendParams{})
	//elastic.WhereMultiMatch([]string{"create_time", "posttime"}, "fasdfasdf", elastic.BestFields, fulltext.AppendParams{})
	//elastic.OrWhereNested(func(c *elastic.Builder) *elastic.Builder {
	//	return c.Where("posttime", "fdajlsdf").Where("create_time", "sdfsadjf")
	//}).OrWhere("author", "fdsafjl").MinimumShouldMatch(1)
	//
	//elastic.Filter("comment", "dsaffd").FilterBetween("comment_time", 1, 2)
	//elastic.FilterNested(func(c *elastic.Builder) *elastic.Builder {
	//	return c.OrWhere("zan", "safdsa").OrWhere("zan2", "fsafdsa")
	//}).Order("create_time", elastic.Desc).Order("posttime", elastic.Asc)

	//[]string{"create_time", "posttime"}
	//elastic.OrWhereNested(func(c *elastic.Condition) *elastic.Condition {
	//
	//})

	//elastic.Where("title", "中国").WhereNested(func(b *elastic.Builder) *elastic.Builder {
	//	return b.OrWhere("title", "美国").OrWhere("title", "日本")
	//})

	builder := elastic.NewBuilder()
	//
	//builder.Where("status", 1).Where("title", "中国").OrWhere("status", 1).WhereNot("country", "日本").Filter("city", "合肥")
	//builder.OrderBy("status", esearch.Asc).GroupBy("status", aggs.TermsParam{Size: 20, Order: map[string]esearch.OrderType{"_count": esearch.Asc}}, func() aggs.TopHitsParam {
	//	return aggs.TopHitsParam{From: 0, Size: 100}
	//}).GroupBy("modify_date", aggs.TermsParam{}, func() aggs.TopHitsParam {
	//	return aggs.TopHitsParam{Size: 43}
	//}).Sum("count", aggs.MetricParam{}).Stats("state", aggs.MetricParam{}).TopHitsFunc(func(b *elastic.Builder) {
	//	b.Size(100).Select("state,title").OrderBy("news_posttime", esearch.Desc)
	//}).WhereNested(func(b *elastic.Builder) {
	//	b.OrWhere("fsaf", "fsadfsa").OrWhere("abc", "abc")
	//}).WhereMultiMatch([]string{"fsadf", "fsaf"}, "fsadfas", esearch.Phrase, func() fulltext.AppendParams {
	//	return fulltext.AppendParams{
	//		Operator:           "and",
	//		MinimumShouldMatch: "100%",
	//	}
	//})
	//elastic.DateGroupBy("posttime", aggs.HistogramParam{Interval: "1day", Format: "yyyy-MM-dd"})
	//elastic.Range("create_time", aggs.RangeParam{Format: "yyyy-MM-dd", Ranges: []aggs.Ranges{{To: 50}, {From: 50, To: 100}, {From: 100}}})
	//elastic.TopHits(aggs.TopHits{From: 0, Size: 10, Sort: map[string]es.Order{"posttime": {Order: es.Asc}}})

	dsl := builder.Dsl()

	fmt.Println(dsl)

	return
	builder = elastic.NewBuilder().Select("status").Where("status", 1000).Where("title", "中国").OrWhere("status", 13).WhereNot("country", "日本").Filter("city", "合肥")
	//elastic.OrderBy("status", es.Asc).GroupBy("status", aggs.TermsParam{Size: 20, Order: map[string]es.OrderType{"_count": es.Asc}}, func() aggs.TopHitsParam {
	//	return aggs.TopHitsParam{From: 0, Size: 100}
	//}).GroupBy("modify_date", aggs.TermsParam{}, func() aggs.TopHitsParam {
	//	return aggs.TopHitsParam{Size: 43}
	//}).Sum("count", aggs.MetricParam{}).Stats("state", aggs.MetricParam{}).TopHitsFunc(func(b *elastic.Builder) *elastic.Builder {
	//	return b.Size(100).Select("state,title").OrderBy("news_posttime", es.Desc)
	//})

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()

	dsl1 := builder.Dsl()

	fmt.Println(dsl1)

	builder1 := builder.Clone()
	builder1.AppendField("content").Where("fsafas", "fsafas").Where("opop", 12132)
	fmt.Println(builder1.Dsl())
	fmt.Println(builder.Dsl())

}
