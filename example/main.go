package main

import (
	"encoding/json"
	"fmt"
	elastic "github.com/KingSolvewer/elasticsearch-query-builder"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
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
	elastic.OrWhere("title", "saffdasdf").OrWhere("title", "fsdafasfd")

	elastic.WhereNested(func(c *elastic.Condition) *elastic.Condition {
		return c.OrWhere("create_time", "12312312").OrWhere("create_time", "456354").Where("create_time", "345253")
	}).WhereNested(func(c *elastic.Condition) *elastic.Condition {
		return c.WhereNested(func(c *elastic.Condition) *elastic.Condition {
			return c.OrWhere("posttime", "123123123").OrWhere("posttime", "34432341")
		})
	})

	//elastic.OrWhereNested(func(c *elastic.Condition) *elastic.Condition {
	//
	//})

	condition := elastic.GetCondition()
	fmt.Println(condition)

	dsl := condition.Dsl
	fmt.Println(dsl)

	dslJson, err := json.Marshal(dsl)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(dslJson))
}
