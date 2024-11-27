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

	type EsResult struct {
		NewsUuid     string `json:"news_uuid"`
		NewsTitle    string `json:"news_title"`
		NewsContent  string `json:"news_content"`
		NewsPostTime string `json:"news_posttime"`
	}

	var esResult []EsResult

	es := zyzx.NewEs()
	es.GroupBy(zyzx.Platform, aggs.TermsParam{}, func(b *elastic.Builder) {
		b.GroupBy(zyzx.SimHash, aggs.TermsParam{}, nil)
		//b.TopHitsFunc(func(b *elastic.Builder) {
		//	b.Size(1).OrderBy(zyzx.PostTime, esearch.Asc).Select(zyzx.NewsUuid, zyzx.NewsTitle, zyzx.NewsContent, zyzx.PostTime)
		//})
	})
	result, err := es.Size(0).Select(zyzx.NewsUuid, zyzx.NewsTitle, zyzx.NewsContent, zyzx.PostTime).Get(&esResult)

	log.Println(es, es.QueryDsl)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(result, es.GetResult(), esResult)
}
