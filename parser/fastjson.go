package parser

import (
	"errors"
	"fmt"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
	"github.com/valyala/fastjson"
	"github.com/valyala/fastjson/fastfloat"
	"reflect"
	"strconv"
	"strings"
)

func AggValueParser(obj *fastjson.Object, dest any) (aggsResult *esearch.AggsResult, errorSet []error) {
	aggsResult = &esearch.AggsResult{}
	errorSet = make([]error, 0)

	obj.Visit(func(k []byte, v *fastjson.Value) {
		key := string(k)

		lastIndex := strings.LastIndex(key, "_")

		if lastIndex != -1 && lastIndex+1 < len(key) {
			lastString := key[lastIndex:]
			switch lastString {
			case esearch.Terms:
				bucketArr := v.GetArray("buckets")

				buckets := make([]esearch.Bucket, len(bucketArr))
				for i, item := range bucketArr {
					bucketObj := item.GetObject()

					var bucket esearch.Bucket
					if bucketObj != nil {
						bucket, errorSet = subAggParser(bucketObj, dest)
					}
					buckets[i] = bucket
				}
				aggsResult.Terms = make(map[string]*esearch.TermsResult)
				aggsResult.Terms[key] = &esearch.TermsResult{
					DocCountErrorUpperBound: v.GetInt("doc_count_error_upper_bound"),
					SumOtherDocCount:        v.GetInt("sum_other_doc_count"),
					Buckets:                 buckets,
				}
			case esearch.Histogram, esearch.DateHistogram:
				bucketsArr := v.GetArray("buckets")

				buckets := make([]esearch.Bucket, len(bucketsArr))
				for i, item := range bucketsArr {
					bucketObj := item.GetObject()

					var bucket esearch.Bucket
					if bucketObj != nil {
						bucket, errorSet = subAggParser(bucketObj, dest)
					}

					buckets[i] = bucket

				}
				aggsResult.Histogram = make(map[string]*esearch.HistogramResult)
				aggsResult.Histogram[key] = &esearch.HistogramResult{
					Buckets: buckets,
				}
			case esearch.Range, esearch.DateRange:
				bucketsArr := v.GetArray("buckets")

				buckets := make([]esearch.Bucket, len(bucketsArr))
				for i, item := range bucketsArr {
					bucketObj := item.GetObject()

					var bucket esearch.Bucket
					if bucketObj != nil {
						bucket, errorSet = subAggParser(bucketObj, dest)
					}

					buckets[i] = bucket
				}

				aggsResult.Range = make(map[string]*esearch.RangeResult)
				aggsResult.Range[key] = &esearch.RangeResult{
					Buckets: buckets,
				}
			case esearch.Cardinality, esearch.ValueCount:
				aggsResult.Count = make(map[string]*esearch.CountResult)
				aggsResult.Count[key] = &esearch.CountResult{Value: v.GetInt("value")}
			case esearch.Avg, esearch.Max, esearch.Min, esearch.Sum:
				aggsResult.Arithmetic = make(map[string]*esearch.ArithmeticResult)
				aggsResult.Arithmetic[key] = &esearch.ArithmeticResult{Value: v.GetFloat64("value")}
			case esearch.Stats:
				aggsResult.Stats = make(map[string]*esearch.StatsResult)
				aggsResult.Stats[key] = &esearch.StatsResult{
					Count: v.GetInt("count"),
					Max:   v.GetFloat64("max"),
					Min:   v.GetFloat64("min"),
					Sum:   v.GetFloat64("sum"),
					Avg:   v.GetFloat64("avg"),
				}
			case esearch.ExtendedStats:
				aggsResult.ExtendedStats = make(map[string]*esearch.ExtendStatsResult)
				stdDeviationBoundsV := v.Get("std_deviation_bounds")
				stdDeviationBounds := esearch.StdDeviationBounds{
					Upper: stdDeviationBoundsV.GetFloat64("upper"),
					Lower: stdDeviationBoundsV.GetFloat64("lower"),
				}

				aggsResult.ExtendedStats[key] = &esearch.ExtendStatsResult{
					StatsResult: &esearch.StatsResult{
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
			case esearch.TopHits:
				topHitsV := v.Get("hits")

				hitsArr := topHitsV.GetArray("hits")
				hitsBuckets := make([]*esearch.HitsBucket, len(hitsArr))
				for i, hitsV := range hitsArr {
					newDest, err := topHitsParser(hitsV, dest)
					if err != nil {
						errorSet = append(errorSet, err)
					}

					hitsBuckets[i] = &esearch.HitsBucket{
						Source: newDest,
					}
				}
				aggsResult.TopHits = &esearch.HitsResult{
					Total: topHitsV.GetInt("total"),
					Hits:  hitsBuckets,
				}
			}
		}
	})

	return aggsResult, nil
}

func subAggParser(obj *fastjson.Object, dest any) (rootBucket esearch.Bucket, errorSet []error) {
	errorSet = make([]error, 0)
	rootBucket = esearch.Bucket{
		Aggs: esearch.AggsResult{
			Terms:         make(map[string]*esearch.TermsResult),
			Histogram:     make(map[string]*esearch.HistogramResult),
			Range:         make(map[string]*esearch.RangeResult),
			Count:         make(map[string]*esearch.CountResult),
			Arithmetic:    make(map[string]*esearch.ArithmeticResult),
			Stats:         make(map[string]*esearch.StatsResult),
			ExtendedStats: make(map[string]*esearch.ExtendStatsResult),
			TopHits:       &esearch.HitsResult{},
		},
	}
	obj.Visit(func(k []byte, v *fastjson.Value) {
		key := string(k)
		if key == "key" {
			rootBucket.Key = string(v.GetStringBytes())
		} else if key == "doc_count" {
			rootBucket.DocCount = v.GetInt()
		} else if key == "key_as_string" {
			rootBucket.KeyAsString = string(v.GetStringBytes())
		} else if key == "to" {
			rootBucket.To = v.GetFloat64()
		} else if key == "from" {
			rootBucket.From = v.GetFloat64()
		} else {
			lastIndex := strings.LastIndex(key, "_")
			if lastIndex != -1 && lastIndex+1 < len(key) {
				lastString := key[lastIndex:]
				switch lastString {
				case esearch.Terms:
					bucketsArr := v.GetArray("buckets")

					subBuckets := make([]esearch.Bucket, len(bucketsArr))
					for i, bucket := range bucketsArr {
						bucketObj := bucket.GetObject()

						var subBucket esearch.Bucket
						if bucketObj != nil {
							subBucket, errorSet = subAggParser(bucketObj, dest)
						}
						subBuckets[i] = subBucket
					}
					rootBucket.Aggs.Terms[key] = &esearch.TermsResult{
						DocCountErrorUpperBound: v.GetInt("doc_count_error_upper_bound"),
						SumOtherDocCount:        v.GetInt("sum_other_doc_count"),
						Buckets:                 subBuckets,
					}
				case esearch.Histogram, esearch.DateHistogram:
					bucketsArr := v.GetArray("buckets")

					buckets := make([]esearch.Bucket, len(bucketsArr))
					for i, item := range bucketsArr {
						bucketObj := item.GetObject()

						var bucket esearch.Bucket
						if bucketObj != nil {
							bucket, errorSet = subAggParser(bucketObj, dest)
						}

						buckets[i] = bucket
					}

					rootBucket.Aggs.Histogram[key] = &esearch.HistogramResult{
						Buckets: buckets,
					}
				case esearch.Range, esearch.DateRange:
					bucketsArr := v.GetArray("buckets")

					buckets := make([]esearch.Bucket, len(bucketsArr))
					for i, item := range bucketsArr {
						bucketObj := item.GetObject()

						var bucket esearch.Bucket
						if bucketObj != nil {
							bucket, errorSet = subAggParser(bucketObj, dest)
						}

						buckets[i] = bucket
					}

					rootBucket.Aggs.Range[key] = &esearch.RangeResult{
						Buckets: buckets,
					}
				case esearch.Cardinality, esearch.ValueCount:
					rootBucket.Aggs.Count[key] = &esearch.CountResult{Value: v.GetInt("value")}
				case esearch.Avg, esearch.Max, esearch.Min, esearch.Sum:
					rootBucket.Aggs.Arithmetic[key] = &esearch.ArithmeticResult{Value: v.GetFloat64("value")}
				case esearch.Stats:
					rootBucket.Aggs.Stats[key] = &esearch.StatsResult{
						Count: v.GetInt("count"),
						Max:   v.GetFloat64("max"),
						Min:   v.GetFloat64("min"),
						Sum:   v.GetFloat64("sum"),
						Avg:   v.GetFloat64("avg"),
					}
				case esearch.ExtendedStats:
					stdDeviationBoundsV := v.Get("std_deviation_bounds")
					stdDeviationBounds := esearch.StdDeviationBounds{
						Upper: stdDeviationBoundsV.GetFloat64("upper"),
						Lower: stdDeviationBoundsV.GetFloat64("lower"),
					}

					rootBucket.Aggs.ExtendedStats[key] = &esearch.ExtendStatsResult{
						StatsResult: &esearch.StatsResult{
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
				case esearch.TopHits:
					topHitsV := v.Get("hits")

					hitsArr := topHitsV.GetArray("hits")
					hitsBuckets := make([]*esearch.HitsBucket, len(hitsArr))
					for i, hitsV := range hitsArr {

						newDest, err := topHitsParser(hitsV, dest)
						if err != nil {
							errorSet = append(errorSet, err)
						}
						hitsBuckets[i] = &esearch.HitsBucket{
							Source: newDest,
						}
					}
					rootBucket.Aggs.TopHits = &esearch.HitsResult{
						Total: topHitsV.GetInt("total"),
						Hits:  hitsBuckets,
					}
				}
			}
		}
	})
	return
}

func topHitsParser(hitsV *fastjson.Value, dest any) (newDest any, err error) {

	switch dest.(type) {
	case nil, map[string]any, *map[string]any:
		mapValue := make(map[string]any)
		mapValue["id"] = string(hitsV.GetStringBytes("_id"))

		sourceV := hitsV.GetObject("_source")
		if sourceV != nil {
			sourceV.Visit(func(keyByte []byte, v *fastjson.Value) {
				key := string(keyByte)
				mapValue[key] = ConvertValue(v)
			})
		}
		_, ok := dest.(map[string]any)
		if dest == nil || ok {
			return mapValue, nil
		}

		return &mapValue, nil
	default:
		var val any

		reflectPtr := reflect.ValueOf(dest)
		reflectValue := reflectPtr.Elem()
		reflectValueType := reflectValue.Type()

		switch reflectValue.Kind() {
		case reflect.Struct:
			newStructPtr := reflect.New(reflect.TypeOf(dest).Elem())
			newStruct := newStructPtr.Elem()

			sourceV := hitsV.Get("_source")

			for i := 0; i < reflectValueType.NumField(); i++ {
				field := reflectValueType.Field(i)
				name := field.Tag.Get("json")
				if name == "" {
					name = field.Name
				}
				val, err = setFieldValue(field, sourceV.Get(name))
				if err != nil {
					return err, nil
				}
				newStruct.Field(i).Set(reflect.ValueOf(val))
			}
			newDest = newStruct.Interface()
		default:
			return errors.New("only supports parsing variables of the types map, *map, *struct"), nil
		}
	}
	return
}

func HitsValueParser(hitsListV []*fastjson.Value, dest any) (err error) {
	if len(hitsListV) == 0 {
		return nil
	}

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
				mapValue[key] = ConvertValue(v)
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
					m[key] = ConvertValue(v)
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
				return errors.New("only supports parsing variables of the types map, *map, *struct, *[]struct, *[]map")
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
			return errors.New("only supports parsing variables of the types map, *map, *struct, *[]struct, *[]map")
		}
	}
	return nil
}

func setFieldValue(field reflect.StructField, value *fastjson.Value) (val any, err error) {
	switch field.Type.Kind() {
	case reflect.Int:
		val, err = GetInt(value)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err = GetInt64(value)
	case reflect.Uint:
		val, err = GetUint(value)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err = GetUint64(value)
	case reflect.String:
		val, err = GetString(value)
	case reflect.Float64, reflect.Float32:
		val, err = GetFloat64(value)
	case reflect.Bool:
		val, err = GetBool(value)
	}

	return val, err
}

func GetString(value *fastjson.Value) (val string, err error) {
	if value == nil {
		return "", nil
	}

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

func GetInt(value *fastjson.Value) (val int, err error) {
	if value == nil {
		return 0, nil
	}

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

func GetInt64(value *fastjson.Value) (val int64, err error) {
	if value == nil {
		return 0, nil
	}
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

func GetUint(value *fastjson.Value) (val uint, err error) {
	if value == nil {
		return 0, nil
	}

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

func GetUint64(value *fastjson.Value) (val uint64, err error) {
	if value == nil {
		return 0, nil
	}
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

func GetFloat64(value *fastjson.Value) (val float64, err error) {
	if value == nil {
		return 0, nil
	}

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

func GetBool(value *fastjson.Value) (val bool, err error) {
	if value == nil {
		return false, nil
	}

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

// ConvertValue 将 fastjson.Value 转换为 Go 原生类型
func ConvertValue(v *fastjson.Value) any {
	switch v.Type() {
	case fastjson.TypeObject:
		m := make(map[string]any)
		v.GetObject().Visit(func(key []byte, val *fastjson.Value) {
			m[string(key)] = ConvertValue(val)
		})
		return m
	case fastjson.TypeArray:
		arr := v.GetArray()
		result := make([]any, len(arr))
		for i, val := range arr {
			result[i] = ConvertValue(val)
		}
		return result
	case fastjson.TypeString:
		return string(v.GetStringBytes())
	case fastjson.TypeNumber:
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
