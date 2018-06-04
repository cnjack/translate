package translate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var (
	BAutoLang Language = "auto"
	BZhLang   Language = "zh"
	BEnLang   Language = "en"
	BJpLang   Language = "jp"
)

var (
	BaiduHttpApiGateway  = "http://api.fanyi.baidu.com/api/trans/vip/translate"
	BaiduHttpsApiGateway = "https://fanyi-api.baidu.com/api/trans/vip/translate"
)

type baiduTranslate struct {
	appID, appSecret string
	apiGateway       string
	httpClient       *http.Client
}

func NewBaiduTranslate(appID, appSecret, apiGateway string) ITranslate {
	return &baiduTranslate{
		appID:      appID,
		appSecret:  appSecret,
		apiGateway: apiGateway,
		httpClient: GetHttpClient(),
	}
}

func (t *baiduTranslate) Translate(query string, op Option) (string, error) {
	if op.From == "" {
		op.From = BAutoLang
	}
	if op.To == "" {
		op.To = BAutoLang
	}
	salt := t.randInt()
	param := &url.Values{}
	param.Set("q", query)
	param.Set("from", string(op.From))
	param.Set("to", string(op.To))
	param.Set("appid", t.appID)
	param.Set("salt", salt)
	param.Set("sign", t.sign(query, salt))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", t.apiGateway, param.Encode()), nil)
	if err != nil {
		return "", err
	}
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respParam, err := t.decode(body)
	if err != nil {
		return "", err
	}
	if respParam.ErrorMsg != "" {
		return "", respParam.BaiduRespErr
	}
	if len(respParam.TransResult) > 0 {
		return respParam.TransResult[0].DST, nil
	}
	return "", errors.New("baidu translate service error")
}

type BaiduRespErr struct {
	ErrorCode int    `json:"error_code,string,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

func (b BaiduRespErr) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", b.ErrorCode, b.ErrorMsg)
}

type BaiduResp struct {
	BaiduRespErr
	From        Language `json:"from,omitempty"`
	To          Language `json:"to,omitempty"`
	TransResult []struct {
		SRC string `json:"src"`
		DST string `json:"dst"`
	} `json:"trans_result,omitempty"`
}

func (t *baiduTranslate) decode(body []byte) (*BaiduResp, error) {
	resp := &BaiduResp{}
	err := json.Unmarshal(body, resp)
	return resp, err
}

func (t *baiduTranslate) sign(query, salt string) string {
	s := t.appID + query + salt + t.appSecret
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (t *baiduTranslate) randInt() string {
	return strconv.Itoa(rand.Intn(100000000))
}
