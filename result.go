package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
	"github.com/valyala/fastjson"
	"github.com/valyala/fastjson/fastfloat"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Response struct {
	Took         int  `json:"took"`
	TimedOut     bool `json:"timed_out"`
	Shards       `json:"_shards"`
	Hits         `json:"hits"`
	ScrollId     string         `json:"_scroll_id,omitempty"`
	Aggregations map[string]any `json:"aggregations"`
	Dest         any
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits struct {
	Total    int          `json:"total"`
	MaxScore float32      `json:"max_score"`
	Hits     []HitsResult `json:"hits"`
}

type HitsResult struct {
	Index string  `json:"_index"`
	Type  string  `json:"_type"`
	Id    string  `json:"_id"`
	Score float32 `json:"_score"`
	//Source    map[string]any       `json:"_source"`
	Source    json.RawMessage      `json:"_source"`
	Fields    map[string][]string  `json:"fields,omitempty"`
	InnerHits map[string]InnerHits `json:"inner_hits,omitempty"`
	Sort      []any                `json:"sort,omitempty"`
}

type InnerHits struct {
	Hits
}

type Result struct {
	Total         int               `json:"total"`
	List          any               `json:"list"`
	Aggs          map[string]any    `json:"aggs,omitempty"`
	ScrollId      string            `json:"scroll_id,omitempty"`
	OriginalTotal esearch.Paginator `json:"original_total,omitempty"`
	PerPage       esearch.Paginator `json:"per_page,omitempty"`
	CurrentPage   esearch.Paginator `json:"current_page,omitempty"`
	LastPage      esearch.Paginator `json:"last_page,omitempty"`
}

type CountResult struct {
	Value int `json:"value"`
}

type ArithmeticResult struct {
	Value float64 `json:"value"`
}

type StatsResult struct {
	Count int     `json:"count"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Sum   float64 `json:"sum"`
}

type ExtendStatsResult struct {
	*StatsResult
	SumOfSquares       float64 `json:"sum_of_squares"`
	Variance           float64 `json:"variance"`
	StdDeviation       float64 `json:"std_deviation"`
	StdDeviationBounds `json:"std_deviation_bounds"`
}

type StdDeviationBounds struct {
	Upper float64 `json:"upper"`
	Lower float64 `json:"lower"`
}

type TermsResult struct {
	DocCountErrorUpperBound int      `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int      `json:"sum_other_doc_count"`
	Buckets                 []Bucket `json:"buckets"`
}

type Bucket struct {
	Key         string `json:"key"`
	KeyAsString string `json:"key_as_string,omitempty"`
	DocCount    int    `json:"doc_count"`
}

type HistogramResult struct {
	Buckets []Bucket `json:"buckets"`
}

type RangeResult struct {
	Buckets []RangeBucket
}

type RangeBucket struct {
	Key          string `json:"key"`
	To           string `json:"to,omitempty"`
	ToAsString   string `json:"to_as_string,omitempty"`
	From         string `json:"from,omitempty"`
	FromAsString string `json:"from_as_string,omitempty"`
	DocCount     int    `json:"doc_count"`
}

func (b *Builder) Parser(result *Result, value []byte, dest any) error {
	data := make([]byte, len(value))
	copy(data, value)

	var parser fastjson.Parser
	v, err := parser.ParseBytes(data)
	if err != nil {
		return err
	}

	code := v.GetInt("code")
	if code != 0 {
		msgV := v.Get("msg")
		errorV := msgV.Get("error")
		log.Fatalln(errorV)
	}

	hitsV := v.Get("hits")
	result.Total = hitsV.GetInt("total")

	if dest != nil {
		hitsListV := hitsV.GetArray("hits")
		if len(hitsListV) > 0 {
			switch dest.(type) {
			case map[string]any, *map[string]any:
				mapValue, ok := dest.(map[string]interface{})
				if !ok {
					if m, ok := dest.(*map[string]interface{}); ok {
						if *m == nil {
							*m = map[string]interface{}{}
						}
						mapValue = *m
					}
				}
				mapValue["id"] = string(hitsListV[0].GetStringBytes("_id"))

				sourceV := hitsListV[0].GetObject("_source")
				if sourceV != nil {
					sourceV.Visit(func(keyByte []byte, v *fastjson.Value) {
						key := string(keyByte)
						mapValue[key] = convertValue(v)
					})
				}
			case *[]map[string]any:
				d := dest.(*[]map[string]any)
				for _, itemV := range hitsListV {
					sourceV := itemV.GetObject("_source")
					if sourceV != nil {
						sourceV.Visit(func(keyByte []byte, v *fastjson.Value) {
							m := make(map[string]any)
							key := string(keyByte)
							m[key] = convertValue(v)
							*d = append(*d, m)
						})
					}
				}
				dest = d
			default:
				var val any

				reflectPtr := reflect.ValueOf(dest)
				reflectValue := reflectPtr.Elem()
				reflectValueType := reflectValue.Type()

				switch reflectValue.Kind() {
				case reflect.Slice, reflect.Array:
					reflectElemType := reflectValueType.Elem()
					if reflectElemType.Kind() == reflect.Struct {
						elemValue := reflect.New(reflectElemType).Elem()
						for _, itemV := range hitsListV {
							sourceV := itemV.Get("_source")

							for i := 0; i < reflectElemType.NumField(); i++ {
								field := reflectElemType.Field(i)
								name := field.Tag.Get("json")
								if name == "" {
									name = field.Name
								}

								res := sourceV.Get(name)
								if res == nil {
									return fmt.Errorf("%v key is not existing in query result", name)
								}

								val, err = setFieldValue(field, res)
								if err != nil {
									return err
								}
								elemValue.Field(i).Set(reflect.ValueOf(val))
							}

							reflectValue.Set(reflect.Append(reflectValue, elemValue))
						}

					} else {
						return errors.New("仅支持解析map,*map,*struct,*[]struct,*[]map类型的变量")
					}
				case reflect.Struct:
					sourceV := hitsListV[0].Get("_source")

					for i := 0; i < reflectValueType.NumField(); i++ {
						field := reflectValueType.Field(i)
						name := field.Tag.Get("json")
						if name == "" {
							name = field.Name
						}

						val, err = setFieldValue(field, sourceV.Get(name))
						if err != nil {
							return err
						}
						reflectValue.Field(i).Set(reflect.ValueOf(val))
					}
				default:
					return errors.New("仅支持解析map,*map,*struct,*[]struct,*[]map类型的变量")
				}
			}
		}
		result.List = dest
	}

	aggsV := v.GetObject("aggregations")
	if aggsV != nil {
		result.Aggs = make(map[string]any)
		aggsV.Visit(func(k []byte, v *fastjson.Value) {
			key := string(k)

			lastIndex := strings.LastIndex(key, "_")
			if lastIndex != -1 && lastIndex+1 < len(key) {
				lastString := key[lastIndex+1:]
				switch lastString {
				case esearch.Terms:
					bucketsV := v.GetArray("buckets")
					buckets := make([]Bucket, len(bucketsV))
					for i, item := range bucketsV {
						buckets[i] = Bucket{
							Key:      string(item.GetStringBytes("key")),
							DocCount: item.GetInt("doc_count"),
						}
					}

					result.Aggs[key] = &TermsResult{
						DocCountErrorUpperBound: v.GetInt("doc_count_error_upper_bound"),
						SumOtherDocCount:        v.GetInt("sum_other_doc_count"),
						Buckets:                 buckets,
					}
				case esearch.Histogram, esearch.DateHistogram:
					bucketsV := v.GetArray("buckets")
					buckets := make([]Bucket, len(bucketsV))
					for i, item := range bucketsV {
						buckets[i] = Bucket{
							Key:      string(item.GetStringBytes("key")),
							DocCount: item.GetInt("doc_count"),
						}
						if lastString == esearch.DateHistogram {
							buckets[i].KeyAsString = string(item.GetStringBytes("key_as_string"))
						}
					}

					result.Aggs[key] = &HistogramResult{
						Buckets: buckets,
					}
				case esearch.Range, esearch.DateRange:
					bucketsV := v.GetArray("buckets")
					buckets := make([]RangeBucket, len(bucketsV))

					for i, item := range bucketsV {
						buckets[i] = RangeBucket{
							Key:      string(item.GetStringBytes("key")),
							To:       string(item.GetStringBytes("to")),
							From:     string(item.GetStringBytes("from")),
							DocCount: item.GetInt("doc_count"),
						}
						if lastString == esearch.DateRange {
							buckets[i].ToAsString = string(item.GetStringBytes("to_as_string"))
							buckets[i].FromAsString = string(item.GetStringBytes("from_as_string"))
						}
					}

					result.Aggs[key] = &RangeResult{
						Buckets: buckets,
					}
				case esearch.Cardinality, esearch.ValueCount:
					result.Aggs[key] = &CountResult{Value: v.GetInt("value")}
				case esearch.Avg, esearch.Max, esearch.Min, esearch.Sum:
					f := v.GetFloat64("value")
					result.Aggs[key] = &ArithmeticResult{Value: f}
				case esearch.Stats:
					result.Aggs[key] = &StatsResult{
						Count: v.GetInt("count"),
						Max:   v.GetFloat64("max"),
						Min:   v.GetFloat64("min"),
						Sum:   v.GetFloat64("sum"),
						Avg:   v.GetFloat64("avg"),
					}
				case esearch.ExtendedStats:
					stdDeviationBoundsV := v.Get("std_deviation_bounds")
					stdDeviationBounds := StdDeviationBounds{
						Upper: stdDeviationBoundsV.GetFloat64("upper"),
						Lower: stdDeviationBoundsV.GetFloat64("lower"),
					}

					result.Aggs[key] = &ExtendStatsResult{
						StatsResult: &StatsResult{
							Count: v.GetInt("count"),
							Max:   v.GetFloat64("max"),
							Min:   v.GetFloat64("min"),
							Sum:   v.GetFloat64("sum"),
							Avg:   v.GetFloat64("avg"),
						},
						SumOfSquares:       v.GetFloat64("sum_of_squares"),
						Variance:           v.GetFloat64("variance"),
						StdDeviation:       v.GetFloat64("std_deviation"),
						StdDeviationBounds: stdDeviationBounds,
					}
				}
			}
		})
	}

	if b.GetScroll() != "" {
		result.ScrollId = string(v.GetStringBytes("scroll_Id"))
	}

	return nil
}

func setFieldValue(field reflect.StructField, value *fastjson.Value) (val any, err error) {
	switch field.Type.Kind() {
	case reflect.Int:
		val, err = getInt(value)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err = getInt64(value)
	case reflect.Uint:
		val, err = getUint(value)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err = getUint64(value)
	case reflect.String:
		val, err = getString(value)
	case reflect.Float64, reflect.Float32:
		val, err = getFloat64(value)
	case reflect.Bool:
		val, err = getBool(value)
	}

	return val, err
}

func getString(value *fastjson.Value) (val string, err error) {
	switch value.Type() {
	case fastjson.TypeString:
		val = string(value.GetStringBytes())
	case fastjson.TypeNumber:
		val = value.String()
	case fastjson.TypeTrue:
		val = "true"
	case fastjson.TypeFalse:
		val = "false"
	case fastjson.TypeNull:
		val = ""
	default:
		err = fmt.Errorf("cannot convert %v to string", value.Type())
	}

	return val, err
}

func getInt(value *fastjson.Value) (val int, err error) {
	switch value.Type() {
	case fastjson.TypeNumber:
		val = value.GetInt()
	case fastjson.TypeString:
		// 尝试将字符串转换为整数
		val, err = strconv.Atoi(string(value.GetStringBytes()))
	case fastjson.TypeNull:
		val = 0
	default:
		err = fmt.Errorf("cannot convert %v to int", value.Type())
	}

	return val, err
}

func getInt64(value *fastjson.Value) (val int64, err error) {
	switch value.Type() {
	case fastjson.TypeNumber:
		val = value.GetInt64()
	case fastjson.TypeString:
		// 尝试将字符串转换为整数
		val, err = fastfloat.ParseInt64(string(value.GetStringBytes()))
	case fastjson.TypeNull:
		val = 0
	default:
		err = fmt.Errorf("cannot convert %v to int", value.Type())
	}

	return val, err
}

func getUint(value *fastjson.Value) (val uint, err error) {
	switch value.Type() {
	case fastjson.TypeNumber:
		val = value.GetUint()
	case fastjson.TypeString:
		// 尝试将字符串转换为整数
		var temp uint64
		temp, err = fastfloat.ParseUint64(string(value.GetStringBytes()))
		val = uint(temp)
	case fastjson.TypeNull:
		val = 0
	default:
		err = fmt.Errorf("cannot convert %v to int", value.Type())
	}

	return val, err
}

func getUint64(value *fastjson.Value) (val uint64, err error) {
	switch value.Type() {
	case fastjson.TypeNumber:
		val = value.GetUint64()
	case fastjson.TypeString:
		// 尝试将字符串转换为整数
		val, err = fastfloat.ParseUint64(string(value.GetStringBytes()))
	case fastjson.TypeNull:
		val = 0
	default:
		err = fmt.Errorf("cannot convert %v to int", value.Type())
	}

	return val, err
}

func getFloat64(value *fastjson.Value) (val float64, err error) {
	switch value.Type() {
	case fastjson.TypeNumber:
		val = value.GetFloat64()
	case fastjson.TypeString:
		// 尝试将字符串转换为整数
		val, err = strconv.ParseFloat(string(value.GetStringBytes()), 64)
	case fastjson.TypeNull:
		val = float64(0)
	default:
		err = fmt.Errorf("cannot convert %v to int", value.Type())
	}

	return val, err
}

func getBool(value *fastjson.Value) (val bool, err error) {
	switch value.Type() {
	case fastjson.TypeTrue, fastjson.TypeFalse:
		val = value.GetBool()
	case fastjson.TypeString:
		str := string(value.GetStringBytes())
		val = str == "true" || str == "1" || str == "yes" || str == "on"
	case fastjson.TypeNumber:
		val = value.GetInt64() != 0
	case fastjson.TypeNull:
		val = false
	default:
		err = fmt.Errorf("cannot convert %v to bool", value.Type())
	}
	return val, err
}

// 将 fastjson.Value 转换为 Go 原生类型
func convertValue(v *fastjson.Value) any {
	switch v.Type() {
	case fastjson.TypeObject:
		m := make(map[string]any)
		v.GetObject().Visit(func(key []byte, val *fastjson.Value) {
			m[string(key)] = convertValue(val)
		})
		return m
	case fastjson.TypeArray:
		arr := v.GetArray()
		result := make([]any, len(arr))
		for i, val := range arr {
			result[i] = convertValue(val)
		}
		return result
	case fastjson.TypeString:
		return string(v.GetStringBytes())
	case fastjson.TypeNumber:
		// 尝试先解析为整数，如果失败则解析为浮点数
		if v.GetFloat64() == float64(v.GetInt64()) {
			return v.GetInt64()
		}
		return v.GetFloat64()
	case fastjson.TypeTrue:
		return true
	case fastjson.TypeFalse:
		return false
	case fastjson.TypeNull:
		return nil
	default:
		return nil
	}
}
