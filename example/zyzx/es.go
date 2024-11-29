package zyzx

import (
	"bytes"
	"encoding/json"
	"errors"
	elastic "github.com/KingSolvewer/elasticsearch-query-builder"
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

const PostTime = "news_posttime"
const PostDate = "news_postdate"
const PostHour = "news_posthour"
const FetchTime = "news_fetch_time"
const EnterTime = "solr_create_time"
const Emotion = "news_emotion"
const Platform = "platform"
const PlatformDomainPri = "platform_domain_pri"
const PlatformDomainSec = "platform_domain_sec"
const SimHash = "news_sim_hash"
const ReadCount = "news_read_count"
const ReportCount = "news_reposts_count"
const CommentCount = "news_comment_count"
const LikeCount = "news_like_count"
const KeywordsList = "news_keywords_list"
const PlatformProvince = "platform_province"
const PlatformCity = "platform_city"
const MediaProvince = "media_province"
const MediaCity = "media_city"
const MediaCounty = "media_county"
const RefProvince = "news_content_province" //提及省份
const PlatformName = "platform_name"
const MediaName = "media_name"
const MediaLevel = "media_level"
const ContentCate = "news_content_field"
const MoodPri = "news_mood_pri"
const NewsLocalUrl = "news_local_url"
const NewsDigest = "news_digest"
const NewsTitle = "news_title"
const NewsUuid = "news_uuid"
const MediaCi = "media_CI"
const NewsContent = "news_content"
const NewsEmotion = "news_emotion"
const NewsContentCity = "news_content_city"
const NewsContentCounty = "news_content_county"
const NewsUrl = "news_url"
const NewsEmotionScore = "news_emotion_score"
const NewsOriginContent = "news_origin_content"
const NewsOriginAuthorName = "news_origin_author_name"
const NewsAuthor = "news_author"
const NewsImgUrls = "news_img_urls"
const PubCountry = "media_county" //发布区县
const MediaGrade = "media_grade"
const NewsIsOrigin = "news_is_origin"
const NewsKeywords = "news_keywords"
const MediaRankLocal = "media_rank_local"                         //发布人排名
const NewsContentIndustry = "news_content_industry"               //行业分类
const NewsContentIndustrySec = "news_content_industry_sec"        //二级行业分类
const NewsContentIndustrySecHit = "news_content_industry_sec_hit" //二级行业命中关键词
const MediaOrganization = "media_organization"                    //
const ChannelName = "channel_name"                                //频道
const MediaSpare2 = "media_spare2"                                //项目标签，表示有哪些项目关注这条数据，数组
const MediaIdentity = "media_identity"
const NewsOrigin = "news_origin"
const NewsRepostsCount = "news_reposts_count"       //转发数
const NewsLikeCount = "news_like_count"             //点赞
const NewsReadCount = "news_read_count"             //阅读数
const MediaFollowersCount = "media_followers_count" //粉丝数
const PlatformId = "platform_id"
const ChannelNav = "channel_nav"                         //文章导航栏
const NewsContentIpLocation = "news_content_ip_location" //IP所属地
const NewsContentField = "news_content_field"            // 分类
const MediaVerifiendtype = "media_verifiedtype"          //认证类型
const NewsSensitiveBy = "news_sensitive_by"              //表示文章是否存在敏感
const NewsOcr = "news_ocr"
const NewsSimHashHash = "news_sim_hash" //相似文章数聚合查询
const IsRumor = "is_rumor"              //是否是谣言
const GroupUuid = "group_uuid"          //所属话题ID
const NewsPosition = "news_position"

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

	dsl, err := es.Marshal()
	if err != nil {
		return nil, err
	}

	params := Params{
		Index:      es.index,
		Statement:  dsl,
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
		es.Cardinality(es.GetCollapse().Field, nil)
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
	if es.GetScroll() == "" {
		data, err = es.Query()
	} else {
		data, err = es.ScrollQuery()
	}

	return data, err
}

func (es *Es) GetRawData() ([]byte, error) {
	data, err := es.runQuery()
	return data, err
}

func (es *Es) GetResult() (string, error) {
	data, err := es.runQuery()
	return string(data), err
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
