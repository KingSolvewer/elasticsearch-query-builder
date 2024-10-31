package elastic

import (
	"errors"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
	"reflect"
)

func (b *Builder) Get(dest any) (*Result, error) {
	if err := checkDestType(dest); err != nil {
		return nil, err
	}

	if dest == nil {
		b.Size(0)
	}

	data, err := b.runQuery()
	if err != nil {
		return nil, err
	}
	b.byteData = data

	result := &Result{}

	err = b.Parser(result, data, dest)

	// 查询完毕之后，重置查询语句
	b.Reset()

	return result, err
}

func (b *Builder) Paginator(page, size uint, dest any) (*Result, error) {
	if err := checkDestType(dest); err != nil {
		return nil, err
	}

	var from uint = 0
	if page > 0 {
		from = page - 1
	}
	b.from = from

	if dest == nil {
		b.size = 0
	} else {
		b.size = size
	}

	//if b.collapse.Field != "" {
	//	b.Cardinality(b.collapse.Field, nil)
	//}

	data, err := b.runQuery()
	if err != nil {
		return nil, err
	}

	b.byteData = data

	result := &Result{
		CurrentPage: esearch.Uint(page),
		PerPage:     esearch.Uint(size),
	}

	err = b.Parser(result, data, dest)
	result.LastPage = esearch.Uint(uint(result.Total)/size + 1)

	// 查询完毕之后，重置查询语句
	b.Reset()
	return result, err
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

func (b *Builder) GetRawData() []byte {
	return b.byteData
}

func checkDestType(dest any) error {
	switch dest.(type) {
	case nil, map[string]any, *map[string]any, *[]map[string]any:
	default:
		reflectValue := reflect.ValueOf(dest)
		if reflectValue.Kind() != reflect.Ptr {
			return errors.New("expected Pointer dest")
		}
	}
	return nil
}
