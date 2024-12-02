package yingyan

import (
	"bytes"
	"encoding/json"
	"errors"
	elastic "github.com/KingSolvewer/elasticsearch-query-builder"
	"github.com/KingSolvewer/elasticsearch-query-builder/aggs"
	"github.com/KingSolvewer/elasticsearch-query-builder/esearch"
	"github.com/KingSolvewer/elasticsearch-query-builder/parser"
	"github.com/valyala/fastjson"
	"io"
	"net/http"
	"time"
)

const (
	DateTime           = "2006-01-02 15:04:05"
	DataSize           = 10000
	EsGateWayUrl       = ""
	EsScrollGateWayUrl = ""
)

const (
	FeelingId      = "feeling_id"           // 舆情ID
	DataSource     = "data_source_"         // 数据来源("spider"—爬虫,yingyan"--鹰眼接口,"import"--导入
	NewsSimHash    = "news_sim_hash_"       // 新闻相似哈希
	NewsUuid       = "news_uuid"            // 新闻uuid
	YingYanId      = "feelings_yingyan_id_" // 鹰眼id
	Province       = "province_"            // 发布省
	City           = "city_"                // 发布市
	Title          = "title_"               // 新闻标题
	TitleId        = "title_id_"            // 新闻标题ID
	Digest         = "digest_"              // 新闻摘要
	SourceMedia    = "sourceMedia_"         // 来源媒体，类似于清博的 platform_name
	PublishMedia   = "publishMedia_"        // 发布媒体，类似于清博的 media_name
	NewsUrl        = "url_"                 // 新闻url
	PublishTime    = "publishTime_"         // 发布时间
	Author         = "author_"              // 作者
	Cks            = "cks_"                 // 新闻热词
	MediaType      = "mediaTypeId_"         // 媒体类型，类似于清博的 platform
	RepostsNum     = "reposts_num_"         // 转发数
	CommentsNum    = "comments_nem_"        // 评论数
	DocType        = "doc_type_"            // 类型("feeds":微博, "bbs": 论坛, "news": 新闻)
	OriginId       = "original_id_"         // 原创ID
	OriginUrl      = "original_url_"        // 原创URL
	OriginSource   = "original_source_"     // 转载媒体源，从发布页面上获取
	Subject        = "subject_"             // 所属企业
	Category       = "category_"            // 新闻分类
	CategoryId     = "category_id_"         // 新闻分类ID
	SubCategory    = "subcategory_"         // 新闻子分类
	SubCategoryId  = "subcategory_id_"      // 新闻子分类
	LevelWarning   = "level_warning_"       // 预警级别
	Keywords       = "keywords_"            // 关键词
	Label          = "label_"               // 舆情标签分类(1-上海舆情,2-全国舆情,3-高危舆情)
	DocScore       = "doc_score_"           // 感情打分(1:正面(绿色)，0：中性(橙色)，-1：负面(浅红色) -2:负面(红色)
	IsInvalid      = "is_invalid_"          // 是否有效(0:有效,1:无效)
	Tendency       = "tendency_"            // 倾向性(根据doc_score设定(1:正面，0：中性，-1：负面)
	Stat           = "state_"               // 状态(1:待认领，2:待处理、3:待审批、4:已审批、5:已归档、6:人工无效)
	CreateTime     = "create_time_"         // 创建时间
	CreateUserId   = "create_user_id_"      // 创建人ID
	ModifyTime     = "modify_time_"         // 修改时间
	ModifyUserId   = "modify_user_id_"      // 修改人ID
	IsDelete       = "is_delete_"           // 是否删除(0:未删除,1:已删除)
	ClaimUserId    = "claim_user_id_"       // 认领人ID
	ClaimTime      = "claim_time_"          // 认领时间
	ApprovalUserId = "approve_user_id_"     // 审批人ID
	ApprovalTime   = "approval_time_"       // 审批时间
	LabelGroupId   = "label_group_id_"      // 标签组ID
	LabelItemId    = "label_item_id_"       // 标签项ID
	CustomLabelId  = "custom_label_id_"     // 自定义标签ID
	Content        = "content_"             // 新闻内容
	EnterTime      = "data_create_time"     // 入库时间
	IsPushed       = "is_pushed_"           // 上报舆情 Y是 N否
	IsRecomplain   = "is_recomplain_"       // 是否是夜班认领和处理的数据 Y是 N否
	IsSpecial      = "is_special_"          // 归档类型 Y自动归档 N人工归档
)
const (
	// 状态常量
	StatFilterOcr    = -2 // 过滤链接后缀为#ocr的数据
	StatMissKeyword  = -1 // 黑名单过滤后的数据，以后再也不用过滤
	StatWaitFilter   = 0  // 待关键词黑白名单过滤
	StatWaitClaim    = 1  // 待认领
	StatWaitProcess  = 2  // 待处理
	StatWaitApproval = 3  // 待审核
	//    const STAT_APPROVAL = 4; // 已审核
	StatArchived = 5 // 已归档
	StatInvalid  = 6 // 无效
)

const (
	SourceYingyan = "yingyan"
	SourceSpider  = "spider"
	SourceImport  = "import"
	SourceHistory = "history"
)

const (
	SupSer       = "¤"
	SupSerRegexp = SupSer + ".*?" + SupSer
)

var (
	SelectValidationSet = map[string]string{
		"标题": Title,
		"正文": Content,
		"作者": Author,
	}
)

type Params struct {
	Index      string `json:"index"`
	Statement  string `json:"statement"`
	StartStamp int64  `json:"startStamp"`
	EndStamp   int64  `json:"endStamp"`
	Scroll     string `json:"scroll,omitempty"`
	ScrollId   string `json:"scrollId,omitempty"`
}

type Es struct {
	*elastic.Builder
	startTime string
	endTime   string
	index     string
	jsonValue *fastjson.Value
}

type Result struct {
	Total         int                 `json:"total"`
	List          any                 `json:"list"`
	Aggs          *esearch.AggsResult `json:"aggs,omitempty"`
	ScrollId      string              `json:"scroll_id,omitempty"`
	OriginalTotal esearch.Paginator   `json:"original_total,omitempty"`
	PerPage       esearch.Paginator   `json:"per_page,omitempty"`
	CurrentPage   esearch.Paginator   `json:"current_page,omitempty"`
	LastPage      esearch.Paginator   `json:"last_page,omitempty"`
	TopHits       any                 `json:"-"`
}

func NewEs() *Es {
	builder := elastic.NewBuilder()
	es := &Es{
		Builder:   builder,
		startTime: time.Now().AddDate(0, -3, 0).Format(DateTime),
		endTime:   time.Now().Format(DateTime),
		index:     "all",
	}

	return es
}

func (es *Es) Clone() *Es {
	newEs := &Es{
		Builder: es.Builder.Clone(),
	}

	return newEs
}

func (es *Es) SetIndex(index string) *Es {
	es.index = index
	return es
}

func (es *Es) SearchTime(startTime, endTime string) *Es {
	es.startTime = startTime
	es.endTime = endTime

	return es
}

func (es *Es) Copy() *Es {
	return &Es{
		Builder: es.Builder.Clone(),
	}
}

func (es *Es) parseTime(typ string) (int64, error) {
	dateTime := es.startTime
	if typ == "end" {
		dateTime = es.endTime
	}

	var (
		t   time.Time
		err error
	)
	if dateTime == "" {
		if typ == "start" {
			t = time.Now().AddDate(0, -3, 0)
		} else {
			t = time.Now()
		}
	} else {
		t, err = time.Parse(DateTime, dateTime)
		if err != nil {
			return 0, err
		}
	}

	return t.UnixMilli(), nil
}

func (es *Es) Query() ([]byte, error) {
	jsonData, err := es.getParams(true)
	if err != nil {
		return nil, err
	}

	return es.request(jsonData, EsGateWayUrl)
}

func (es *Es) ScrollQuery() ([]byte, error) {
	jsonData, err := es.getParams(true)
	if err != nil {
		return nil, err
	}

	return es.request(jsonData, EsScrollGateWayUrl)
}

func (es *Es) getParams(scroll bool) ([]byte, error) {

	startStamp, err := es.parseTime("start")
	if err != nil {
		return nil, err
	}
	endStamp, err := es.parseTime("end")
	if err != nil {
		return nil, err
	}

	params := Params{
		Index:      es.index,
		Statement:  es.Dsl(),
		StartStamp: startStamp,
		EndStamp:   endStamp,
	}

	if scroll {
		params.Scroll = es.GetScroll()
		params.ScrollId = es.GetScrollId()
	}

	jsonData, err := json.Marshal(params)
	return jsonData, err
}

func (es *Es) request(jsonData []byte, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "")
	req.Header.Set("Expect", "")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	return body, err
}

func (es *Es) Get(result *Result) error {
	data, err := es.runQuery()
	if err != nil {
		return err
	}

	// 解析成对应的json数据
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	es.jsonValue, err = fastjson.ParseBytes(dataCopy)
	if err != nil {
		return err
	}

	err = es.Parser(result)

	// 查询完毕之后，重置查询语句
	es.Reset()

	return err
}

func (es *Es) Paginator(result *Result, page, size uint) error {

	var from uint = 0
	if page > 0 {
		from = page - 1
	}
	es.From(from)

	if size == 0 {
		es.Size(10)
	} else {
		es.Size(size)
	}

	if es.GetCollapse() != nil {
		es.Cardinality(es.GetCollapse().Field, aggs.CardinalityParam{})
	}

	data, err := es.runQuery()
	if err != nil {
		return err
	}

	// 解析成对应的json数据
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	es.jsonValue, err = fastjson.ParseBytes(dataCopy)
	if err != nil {
		return err
	}

	err = es.Parser(result)

	// 查询完毕之后，重置查询语句
	es.Reset()

	return err
}

func (es *Es) runQuery() (data []byte, err error) {

	_, err = es.Marshal()
	if err != nil {
		return nil, err
	}

	if es.GetScroll() == "" {
		data, err = es.Query()
	} else {
		data, err = es.ScrollQuery()
	}

	return data, err
}

func (es *Es) GetRawData() []byte {
	return es.jsonValue.GetStringBytes()
}

func (es *Es) GetResult() string {
	return string(es.jsonValue.GetStringBytes())
}

func (es *Es) Parser(result *Result) error {
	err := elastic.CheckHitsDestType(result.List)
	if err != nil {
		return err
	}

	err = elastic.CheckTopHitsDestType(result.TopHits)
	if err != nil {
		return err
	}

	code := es.jsonValue.GetInt("code")
	if code != 0 {
		msgV := es.jsonValue.GetStringBytes("msg")
		return errors.New(string(msgV))
	}

	result.Total = es.jsonValue.GetInt("total")

	// 列表
	if result.List != nil {
		hitsV := es.jsonValue.Get("hits").GetArray("hits")

		err = parser.HitsValueParser(hitsV, result.List)
		if err != nil {
			return err
		}
	}

	// 聚合数据
	aggsObj := es.jsonValue.GetObject("aggregations")
	aggsResult, errSet := parser.AggValueParser(aggsObj, result.TopHits)
	if len(errSet) > 0 && errSet[0] != nil {
		return errSet[0]
	}

	result.Aggs = aggsResult

	if es.GetScroll() != "" {
		result.ScrollId = string(es.jsonValue.GetStringBytes("scroll_Id"))
	}

	return nil
}
