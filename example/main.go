package main

import (
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
	result, err := es.Select(zyzx.NewsUuid, zyzx.NewsTitle, zyzx.NewsContent, zyzx.PostTime).Get(&esResult)

	log.Println(es, es.QueryDsl)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(result, es.GetResult(), esResult)
}
