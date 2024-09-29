<h1 align="center">使用go语言构建DSL查询语句</h1>

## 特点
1. 支持链式调用
2. 学习简单，使用方便

## 要求
1. 熟悉go语言基础知识，熟悉es查询语句
2. 简单校验，不符合规范的参数，会进行过滤和跳过组合。使用时，请严格按照es语法输入参数

## 使用
```go
go get github.com/KingSolvewer/elasticsearch-query-builder
```

#### Where, OrWhere, WhereNot, Filter 都是并列关系。构建后，在bool查询条件下同级。
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


#### 默认全局初始化Builder，如果要构建新的查询语句，请重新初始化:
```
builder := elastic.NewBuilder()
```

#### 重置后继续使用
```go
elastic.Reset().Where...

```
