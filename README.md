<h1 align="center">使用go语言构建DSL查询语句</h1>

## 特点
1. 支持链式调用
2. 学习简单，使用方便
3. 支持ES常用查询和聚合语法
4. 支持原生dsl语句查询

## 要求
1. 熟悉go语言基础知识，熟悉es查询语句
2. 简单校验，不符合规范的参数，会进行过滤和跳过组合。使用时，请严格按照es语法输入参数
3. go 运行版本 1.18及以上

## 使用
```go
go get github.com/KingSolvewer/elasticsearch-query-builder
```

#### 默认全局初始化Builder，如果要构建新的查询语句，请重置:
```go
elastic.Where...

...

// 重置后继续使用
elastic.Reset()
elastic.Where...

```

#### 手动初始化Builder，，如果要构建新的查询语句，请重置:
```go
builer := elastic.NewBuilder()
builer.Where...

builer.Reset().Where

```

#### Where,对应must, OrWhere,对应should, WhereNot,对应must_not, Filter,对应filter 都是并列关系。构建后，在bool查询条件下同级。
```go
import (
    elastic "github.com/KingSolvewer/elasticsearch-query-builder"
)

elastic.Where("status", 1).Where("title", "中国").OrWhere("status", 1).WhereNot("country", "日本").Filter("city", "合肥")


DSL:
{
    "size": 0,
    "query": {
        "bool": {
            "must": [
                {
                    "term": {
                        "status": 1
                    }
                },
                {
                "term": {
                        "title": "中国"
                    }
                }
            ],
            "must_not": [
                {
                    "term": {
                        "country": "日本"
                    }
                }
            ],
            "should": [
                {
                    "term": {
                        "status": 1
                    }
                }
            ],
            "filter": [
                {
                    "term": {
                        "city": "合肥"
                    }
                }
            ]
        }
    }
}
```

#### 使用闭包的方式嵌套子句，例如：must 语句中嵌套 should 语句，可以使用，WhereNested。同理还有WhereNotNested，OrWhereNested，FilterNested
```go
import (
    elastic "github.com/KingSolvewer/elasticsearch-query-builder"
)

elastic.Where("title", "中国").WhereNested(func(b *elastic.Builder) *elastic.Builder {
    return b.OrWhere("title", "美国").OrWhere("title", "日本")
})

DSL:
{
    "size": 0,
    "query": {
        "bool": {
            "must": [
                {
                    "term": {
                        "title": "中国"
                    }
                },
                {
                "bool": {
                     "should": [
                           {
                                "term": {
                                    "title": "美国"
                                }
                            },
                            {
                                "term": {
                                    "title": "日本"
                                }
                            }
                        ]
                    }
                }
            ]
        }
    }
}
```

## 方法
##### Select(fields ...string)
```go
    elastic.Select("title", "content")
```

##### AppendField(fields ...string)
```go
    elastic.AppendField("post_time", "author")
```

##### From(value uint)
```go
    elastic.From(0)
```

##### Size(value uint)
```go
    elastic.Size(0)
```

##### OrderBy(field string, orderType esearch.OrderType)
```go
	esearch.Asc
    esearch.Desc

    elastic.OrderBy("post_time", esearch.Desc)
```

##### Raw(raw string)
```go
    elastic.Raw(`{"query":{"bool":{"must":[{"term":{"news_uuid":"*********"}}]}}}`)
```

##### Where(field string, value any)
```go
    elastic.Where("post_time", "2024-11-01 00:00:23")
```

##### WherePrefix(field string, value any)
```go
    elastic.WherePrefix("username", "ki")
```

##### WhereExists(field string)
```go
    elastic.WhereExists("username")
```

##### WhereRegexp(field string, value string, fn termlevel.RegexpParamFunc)
```go
    elastic.WhereRegexp("username", "ki.*y", nil)
    elastic.WhereRegexp("username", "ki.*y", func() termlevel.RegexpParam {
        return termlevel.RegexpParam{}
    })
```

##### WhereWildcard(field string, value string, fn termlevel.WildcardParamFunc)
```go
    elastic.WhereWildcard("username", "ki*y", nil)
    elastic.WhereWildcard("username", "ki*y", func() termlevel.WildcardParam {
        return termlevel.WildcardParam{}
    })
```

##### WhereIn(field string, value []any)
```go
    any 只允许 go的基础类型 int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, string, bool
    
    elastic.WhereIn("news_emotion", elastic.SliceToAny([]string{"中性", ""})
```

##### WhereBetween(field string, value1, value2 any)
```go
    any 只允许 go的基础类型 int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, string, bool
    
    elastic.WhereBetween("publish_time", "2024-11-01 00:00:00", "2024-12-01 00:00:00")
```

##### WhereRange(field string, value any, rangeType esearch.RangeType)
```go
    any 只允许 go的基础类型 int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, string, bool
    esearch.Gte
    esearch.Gt
    esearch.Lt
    esearch.Lte    

    elastic.WhereRange("age", 18, esearch.Gt)
```

##### WhereMatch(field string, value string, matchType esearch.MatchType, fn fulltext.AppendParamsFunc)
```go
    esearch.BestFields
    esearch.MostFields
    esearch.CrossFields
    esearch.Phrase
    esearch.PhrasePrefix
    esearch.BoolPrefix

    elastic.WhereMatch("title", "中国电信", esearch.MatchPhrasePrefix, nil)
    elastic.WhereMatch("title", "中国电信", esearch.MatchPhrasePrefix, func() fulltext.AppendParams {
        return fulltext.AppendParams{}
    })
```

##### WhereMultiMatch(field []string, value string, fieldType esearch.FieldType, fn fulltext.AppendParamsFunc)
```go
    esearch.Match,esearch.MatchPhrase,esearch.MatchPhrasePrefix

    elastic.WhereMultiMatch([]string{"title", "content"}, "中国电信", esearch.BestFields, nil)
    elastic.WhereMultiMatch([]string{"title", "content"}, "中国电信", esearch.BestFields, func() fulltext.AppendParams {
        return fulltext.AppendParams{}
    })
```

##### WhereNested(fn NestWhereFunc)
```go
    elastic.WhereNested(func(b *elastic.Builder) {
        b.Where().WhereIn().OrWhere().WhereNot().Filter()...
    })
```

## WhereNot,OrWhere,Filter 都有以上对应的方法，使用方式相同

## 其他方法
##### Scroll(scroll string) 使用游标查询
```go
    elastic.Scroll("2m")
```

##### GetScroll() string
```go
    elastic.GetScroll()
```

##### ScrollId(scrollId string) 设置游标值
```go
    elastic.ScrollId("**************")
```

##### GetScrollId() string
```go
    elastic.GetScrollId()
```

##### Collapse(field string)  去重
```go
    elastic.Collapse("simhash")
```

##### GetCollapse()
```go
    elastic.GetScrollId()
```

##### Dsl() string 获取DSL语句
```go
    elastic.Dsl()
```

##### Marshal() (string, error) 获取DSL语句和构建语句时的错误
```go
    elastic.Marshal()
```

##### GetQuery() *esearch.ElasticQuery 获取构建DSL语句的结构体实例
```go
    elastic.GetQuery()
```

##### Clone() *Builder 复制Builder实例的字段和方法，生成一个新的Builder实例
```go
    elastic.Clone()
```

##### Reset() *Builder 重置同一个Builder实例，重复使用
```go
    elastic.Reset()
```





#### 常用 Bucket Aggregations
| 名称      | ES语法           | 方法                    | 参数                    | 说明             |
|---------|----------------|-----------------------|-----------------------|----------------|
| 分组聚合    | terms          | elastic.GroupBy()     | aggs.TermsParam{}     | 对某个字段进行分组聚合    |
| 直方图聚合   | histogram      | elastic.Histogram()   | aggs.HistogramParam{} |                |
| 日期直方图聚合 | date_histogram | elastic.DateGroupBy() | aggs.HistogramParam{} | 严格遵照日期直方图聚合的写法 |
| 数值范围聚合  | range          | elastic.Range()       | aggs.RangeParam{}     | 必须是数值类型的数据     |
| 日期范围聚合  | date_range     | elastic.DateRange()   | aggs.RangeParam{}     | 严格遵照日期范围聚合的写法  |

#### 常用 Metrics Aggregations
| 名称      | ES语法           | 方法                      | 参数                    | 说明                                    |
|---------|----------------|-------------------------|-----------------------|---------------------------------------|
| 平均数计算   | avg            | elastic.Avg()           | aggs.MetricParam      |                                       |
| 最大值计算   | max            | elastic.Mx()            | aggs.MetricParam      |                                       |
| 最小值计算   | min            | elastic.Min()           | aggs.MetricParam      |                                       |
| 总和计算    | sum            | elastic.Sum()           | aggs.MetricParam      |                                       |
| 数量计算    | value_count    | elastic.ValueCount()    | 无                     | 类似于 total                             |
| 统计      | stats          | elastic.Stats()         | aggs.MetricParam      |                                       |
| 扩展统计    | extended_stats | elastic.ExtendedStats() | aggs.CardinalityParam |                                       |
| 分组聚合的数据 | top_hits       | elastic.TopHits()       | aggs.TopHitsParam     |                                       |
| 分组聚合的数据 | top_hits       | elastic.TopHitsFunc()   | 闭包函数                  | 支持 b.From(0).Size(10).Select().Sort() |

## 函数
##### elastic.SliceToAny[T SliceInterface](sets []T) 将满足约束的任意类型转换成any类型
```go
elastic.SliceToAny([]string{"a", "b", "c"})
```

