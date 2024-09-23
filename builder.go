package elastic

import (
	"github.com/KingSolvewer/elasticsearch-query-builder/termlevel"
)

type BoolClauseType int

const (
	Must BoolClauseType = iota
	MustNot
	Should
	Filter
)

type NestedFunc func(c *Condition) *Condition

type Condition struct {
	where  map[BoolClauseType][]BoolBuilder
	nested map[BoolClauseType][]NestedFunc
	*Dsl
}

type Criteria struct {
	typ      string
	field    string
	value    any
	operator string
	boolean  string
	not      int
	filter   bool
}

var condition = NewCondition()

func NewCondition() *Condition {
	return &Condition{
		where:  make(map[BoolClauseType][]BoolBuilder),
		nested: make(map[BoolClauseType][]NestedFunc),
	}
}

func GetCondition() *Condition {
	condition.compile()
	return condition
}

// Where Must term 查询语句
func Where(field string, value string) *Condition {
	return condition.Where(field, value)
}

// Where Must term 查询语句
func (c *Condition) Where(field string, value string) *Condition {
	c.termQuery(Must, field, value)

	return c
}

// WhereIn Must terms 查询语句
func WhereIn(field string, value []string) *Condition {
	return condition.WhereIn(field, value)
}

// WhereIn Must terms 查询语句
func (c *Condition) WhereIn(field string, value []string) *Condition {
	c.termsQuery(Must, field, value)

	return c
}

// WhereNot MustNot term 查询
func WhereNot(field string, value string) *Condition {
	return condition.WhereNot(field, value)
}

// WhereNot MustNot term 查询
func (c *Condition) WhereNot(field string, value string) *Condition {
	c.termQuery(MustNot, field, value)

	return c
}

func WhereGt(field string, value int) *Condition {
	return condition.WhereGt(field, value)
}

func (c *Condition) WhereGt(field string, value int) *Condition {
	termQuery := termlevel.TermQuery{
		Range: make(map[string]termlevel.RangeQuery),
	}

	termQuery.Range[field] = termlevel.RangeQuery{
		Gt: value,
	}

	c.where[Must] = append(c.where[Must], termQuery)

	return c
}

//func WhereGte(field string, value int) *Condition {
//	return condition.WhereGt(field, value)
//}
//
//func (c *Condition) WhereGte(field string, value int) *Condition {
//	val := make(map[string]int)
//	val[">="] = value
//
//	c.appendRange(field, val)
//
//	return c
//}
//
//func WhereLt(field string, value int) *Condition {
//	return condition.WhereGt(field, value)
//}
//
//func (c *Condition) WhereLt(field string, value int) *Condition {
//	val := make(map[string]int)
//	val["<"] = value
//
//	c.appendRange(field, val)
//	return c
//}
//
//func WhereLte(field string, value int) *Condition {
//	return condition.WhereGt(field, value)
//}
//
//func (c *Condition) WhereLte(field string, value int) *Condition {
//	val := make(map[string]int)
//	val["<="] = value
//
//	c.appendRange(field, val)
//
//	return c
//}
//
//func WhereRange(field string, values [2]int) *Condition {
//	return condition.WhereRange(field, values)
//}
//
//func (c *Condition) WhereRange(field string, values [2]int) *Condition {
//	val := make(map[string]int)
//	val[">="] = values[0]
//	val["<="] = values[1]
//
//	c.appendRange(field, val)
//
//	return c
//}
//
//func (c *Condition) appendRange(field string, value any) {
//	criteria := Criteria{
//		typ:     "between",
//		field:   field,
//		value:   value,
//		boolean: "and",
//	}
//
//	c.where = append(c.where, criteria)
//}

func WhereNested(fn NestedFunc) *Condition {
	return condition.WhereNested(fn)
}

func (c *Condition) WhereNested(fn NestedFunc) *Condition {
	c.nested[Must] = append(c.nested[Must], fn)

	return c
}

func OrWhere(field string, value string) *Condition {
	return condition.OrWhere(field, value)
}

func (c *Condition) OrWhere(field string, value string) *Condition {
	c.termQuery(Should, field, value)

	return c
}

func (c *Condition) termQuery(clauseTyp BoolClauseType, field string, value string) {
	term := termlevel.TermQuery{
		Term: make(map[string]string),
	}

	term.Term[field] = value

	c.append(clauseTyp, term)
}

func (c *Condition) termsQuery(clauseTyp BoolClauseType, field string, value []string) {
	terms := termlevel.TermQuery{
		Terms: make(map[string][]string),
	}

	terms.Terms[field] = value

	c.append(clauseTyp, terms)
}

func (c *Condition) append(clauseTyp BoolClauseType, clause BoolBuilder) {
	c.where[clauseTyp] = append(c.where[clauseTyp], clause)
}
