package elastic

import (
	"encoding/json"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
	"strings"
)

type Response struct {
	ResponseResult
	Aggregations map[string]esearch.AggResult `json:"aggregations"`
}

type ResponseResult struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits     `json:"hits"`
	ScrollId string `json:"_scroll_id,omitempty"`
}

func (r *Response) UnmarshalJSON(data []byte) error {
	var temp struct {
		ResponseResult
		Aggregations map[string]json.RawMessage `json:"aggregations"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	agg, err := CreateAggResult(temp.Aggregations)
	if err != nil {
		return err
	}
	r.Aggregations = agg
	r.ResponseResult = temp.ResponseResult

	return nil
}

type Hits struct {
	Total    int     `json:"total"`
	MaxScore float32 `json:"max_score"`
	Hits     []struct {
		Index     string               `json:"_index"`
		Type      string               `json:"_type"`
		Id        string               `json:"_id"`
		Score     float32              `json:"_score"`
		Source    map[string]any       `json:"_source"`
		Fields    map[string][]string  `json:"fields,omitempty"`
		InnerHits map[string]InnerHits `json:"inner_hits,omitempty"`
		Sort      []any                `json:"sort,omitempty"`
	} `json:"hits"`
}

type InnerHits struct {
	Hits
}

type Result struct {
	Total         int                          `json:"total"`
	List          []map[string]any             `json:"list"`
	Aggs          map[string]esearch.AggResult `json:"aggs,omitempty"`
	ScrollId      string                       `json:"scroll_id,omitempty"`
	OriginalTotal esearch.Paginator            `json:"original_total,omitempty"`
	PerPage       esearch.Paginator            `json:"per_page,omitempty"`
	CurrentPage   esearch.Paginator            `json:"current_page,omitempty"`
	LastPage      esearch.Paginator            `json:"last_page,omitempty"`
}

type CardinalityAggResult struct {
	Value int `json:"value"`
}

func (c *CardinalityAggResult) UnmarshalJSON(data []byte) error {
	type Alias CardinalityAggResult
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type TermsAggResult struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	} `json:"buckets"`
}

func (terms TermsAggResult) UnmarshalJSON(data []byte) error {
	type Alias TermsAggResult
	aux := &struct {
		Alias
	}{
		Alias: (Alias)(terms),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

//func CreateAggResult(data []byte) (esearch.AggResult, error) {
//	log.Println(string(data))
//	var temp map[string]json.RawMessage
//	if err := json.Unmarshal(data, &temp); err != nil {
//		return nil, err
//	}
//
//	log.Println(temp)
//
//	if _, ok := temp["buckets"]; ok {
//		var result TermsAggResult
//		err := json.Unmarshal(data, &result)
//		return &result, err
//	} else if _, ok := temp["value"]; ok {
//		var result CardinalityAggResult
//		err := json.Unmarshal(data, &result)
//		return &result, err
//	}
//
//	return nil, fmt.Errorf("unknown aggregation result type")
//}

func CreateAggResult(data map[string]json.RawMessage) (res map[string]esearch.AggResult, err error) {
	res = make(map[string]esearch.AggResult)

	for key, item := range data {
		lastIndex := strings.LastIndex(key, "_")
		if lastIndex != -1 && lastIndex+1 < len(key) {
			lastString := key[lastIndex+1:]
			switch lastString {
			case esearch.Cardinality:
				val := &CardinalityAggResult{}
				err = json.Unmarshal(item, val)
				if err == nil {
					res[key] = val
				}
			}
		}
	}
	return res, err

}
