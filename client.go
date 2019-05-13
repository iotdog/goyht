package goyht

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/leesper/holmes"
)

// M is a convenient alias for a map[string]interface{} map.
type M map[string]interface{}

// constants about user types.
const (
	UserTypePersonal   = "1" // 个人
	UserTypeEnterprise = "2" // 企业
	UserTypePlatform   = "4" // 平台
)

// constants about certification types.
const (
	CertTypeIDCard   = "1" // 身份证
	CertTypePassport = "2" // 护照
	CertTypeOfficer  = "3" // 军官证
	CertTypeLicence  = "4" // 营业执照
	CertTypeOrgan    = "5" // 组织机构代码
	CertTypeSocial   = "6" // 社会代码
)

// constants for url and keys
const (
	YHTAPIGateway   = "https://sdk.yunhetong.com/sdk"
	YHTAuthGateway  = "https://authentic.yunhetong.com"
	AppIDKey        = "appId"
	PasswordKey     = "password"
	YHTAPIGatewayV4 = "https://api.yunhetong.com/api" // V4版本使用
)

// Config contains configurations about YunHeTong service.
type Config struct {
	AppID       string // API V4使用
	AppKey      string // API V4使用
	Password    string
	APIGateway  string // API V4使用
	AuthID      string
	AuthPWD     string
	AuthGateway string // API V4使用
}

var (
	// YHTClient 云合同客户端单件对象
	yhtClient *Client
)

// updateLTTRoutine 每隔14分钟更新平台长效令牌
func updateLTTRoutine(yhtClient *Client) {
	for true {
		yhtClient.updateLongTimeToken()
		time.Sleep(14 * time.Minute)
	}
}

// InitYHTClient 初始化云合同客户端，该方法只可调用一次
func InitYHTClient(appID, appKey string) {
	yhtClient = newClient(Config{
		AppID:       appID,
		AppKey:      appKey,
		APIGateway:  YHTAPIGatewayV4,
		AuthGateway: YHTAuthGateway,
	})
	// 开启一个goroutine更新平台长效令牌
	go updateLTTRoutine(yhtClient)
	// yhtClient.updateLongTimeToken() // 测试
}

// GetClient 获取云合同客户端
func GetClient() *Client {
	if yhtClient == nil {
		fmt.Println("YHT client is not initialized, please invoke InitYHTClient() first!")
	}
	return yhtClient
}

// Client handles all APIs for YunHeTong service.
type Client struct {
	config    Config
	tlsClient http.Client
	ltt       string // 平台的长效令牌（Long Time Token），有效期15分钟
}

// newClient returns a *Client.
func newClient(cfg Config) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := http.Client{Transport: tr}
	return &Client{
		config:    cfg,
		tlsClient: client,
		ltt:       "",
	}
}

// updateLongTimeToken 更新长效令牌
func (c *Client) updateLongTimeToken() {
	req := yhtAuthLoginReq{
		AppID:  c.config.AppID,
		AppKey: c.config.AppKey,
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		holmes.Debugln(err)
		return
	}
	ret, ltt, err := httpRequestV4(c, "", req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtBaseResp{}
	})
	if err != nil {
		holmes.Debugln(err)
	} else {
		resp := ret.(*YhtBaseResp)
		if 200 == resp.Code {
			c.ltt = ltt // 保存token
		} else {
			holmes.Debugln(resp)
		}
	}
}

// CreatePersonV4 创建个人用户
func (c *Client) CreatePersonV4(req *YhtCreatePersonReq) (*YhtCreateUserResp, error) {
	if nil == req {
		return nil, errors.New("invalid parameter")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret, _, err := httpRequestV4(c, c.ltt, req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtCreateUserResp{}
	})
	if err != nil {
		return nil, err
	}
	return ret.(*YhtCreateUserResp), nil
}

// CreateCompanyV4 创建企业用户
func (c *Client) CreateCompanyV4(req *YhtCreateCompanyReq) (*YhtCreateUserResp, error) {
	if nil == req {
		return nil, errors.New("invalid parameter")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret, _, err := httpRequestV4(c, c.ltt, req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtCreateUserResp{}
	})
	if err != nil {
		return nil, err
	}
	return ret.(*YhtCreateUserResp), nil
}

// QuerySignerID 查询与合同平台用户ID
func (c *Client) QuerySignerID(req *YhtQuerySignerIDReq) (*YhtQuerySignerIDResp, error) {
	if nil == req {
		return nil, errors.New("invalid parameter")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret, _, err := httpRequestV4(c, c.ltt, req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtQuerySignerIDResp{}
	})
	if err != nil {
		return nil, err
	}
	return ret.(*YhtQuerySignerIDResp), nil
}

// CreatePersonMoulageV4 创建个人印章
func (c *Client) CreatePersonMoulageV4(req *YhtCreatePersonMoulageReq) (*YhtCreateMoulageResp, error) {
	if nil == req {
		return nil, errors.New("invalid parameter")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret, _, err := httpRequestV4(c, c.ltt, req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtCreateMoulageResp{}
	})
	if err != nil {
		return nil, err
	}
	return ret.(*YhtCreateMoulageResp), nil
}

// CreateCompanyMoulageV4 创建企业印章
func (c *Client) CreateCompanyMoulageV4(req *YhtCreateCompanyMoulageReq) (*YhtCreateMoulageResp, error) {
	if nil == req {
		return nil, errors.New("invalid parameter")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret, _, err := httpRequestV4(c, c.ltt, req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtCreateMoulageResp{}
	})
	if err != nil {
		return nil, err
	}
	return ret.(*YhtCreateMoulageResp), nil
}

// CreateContractFromTemplateV4 根据模板创建合同
func (c *Client) CreateContractFromTemplateV4(req *YhtCreateTemplateContractReq) (*YhtCreateTemplateContractResp, error) {
	if nil == req {
		return nil, errors.New("invalid parameter")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret, _, err := httpRequestV4(c, c.ltt, req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtCreateTemplateContractResp{}
	})
	if err != nil {
		return nil, err
	}
	return ret.(*YhtCreateTemplateContractResp), nil
}

// AddSignerV4 添加签署者
func (c *Client) AddSignerV4(req *YhtAddSignerReq) (*YhtBaseResp, error) {
	if nil == req {
		return nil, errors.New("invalid parameter")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret, _, err := httpRequestV4(c, c.ltt, req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtBaseResp{}
	})
	if err != nil {
		return nil, err
	}
	return ret.(*YhtBaseResp), nil
}

// SignContractV4 签署合同（V4版本）
func (c *Client) SignContractV4(req *YhtSignContractReq) (*YhtBaseResp, error) {
	if nil == req {
		return nil, errors.New("invalid parameter")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ret, _, err := httpRequestV4(c, c.ltt, req.URI(), req.Method(), jsonData, func() interface{} {
		return &YhtBaseResp{}
	})
	if err != nil {
		return nil, err
	}
	return ret.(*YhtBaseResp), nil
}

// AuthRealNameMobileV4 运营商三要素认证，认证成功返回nil，否则返回error
func (c *Client) AuthRealNameMobileV4(idNo, idName, phone string) error {
	uri := "/authentic/personal/mobile/realName"
	req := map[string]string{
		"appId":  c.config.AppID,
		"appKey": c.config.AppKey,
		"idNo":   idNo,
		"idName": idName,
		"mobile": phone,
	}
	ret, err := httpRequest(c, uri, req, nil, func() interface{} {
		return &AuthRealNameResp{}
	})
	if err != nil {
		return err
	}
	resp := ret.(*AuthRealNameResp)
	if 200 != resp.Code {
		holmes.Debugln(resp)
		return errors.New(resp.Message())
	}

	return nil
}

// AuthRealNameBankV4 银行四要素认证
func (c *Client) AuthRealNameBankV4(idNo, idName, phone, bankCardNo string) error {
	uri := "/authentic/personal/bankFour"
	req := map[string]string{
		"appId":      c.config.AppID,
		"appKey":     c.config.AppKey,
		"idNo":       idNo,
		"idName":     idName,
		"mobile":     phone,
		"bankCardNo": bankCardNo,
	}
	ret, err := httpRequest(c, uri, req, nil, func() interface{} {
		return &AuthRealNameResp{}
	})
	if err != nil {
		return err
	}
	resp := ret.(*AuthRealNameResp)
	if 200 != resp.Code {
		holmes.Debugln(resp)
		return errors.New(resp.Message())
	}

	return nil
}

// AuthRealName authenticates ID number and name via YunHeTong service.
func (c *Client) AuthRealName(idNum, idName string, portrait bool) (*AuthResponse, error) {
	reqType := "1"
	if portrait {
		reqType = "2"
	}
	p := authParams{
		IDNo:   idNum,
		IDName: idName,
	}

	paramMap, err := toMap(p, map[string]string{
		"key":            c.config.AuthID,
		"value":          c.config.AuthPWD,
		"rcaRequestType": reqType,
	})

	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &AuthResponse{}
	})

	if err != nil {
		return nil, err
	}

	rsp := ret.(*AuthResponse)

	if err = checkAuthErr(rsp.Code, rsp.Msg, rsp.Success); err != nil {
		return nil, err
	}

	info := struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{}

	if err = json.Unmarshal([]byte(rsp.Data), &info); err != nil {
		return nil, err
	}

	rsp.Message = info.Message
	rsp.Status = info.Status

	return rsp, nil
}

func (c *Client) AuthRealNameBank(idNum, idName, bankCard, mobile string) (*AuthResponse, error) {
	reqType := "3"
	p := authParams{
		IDNo:       idNum,
		IDName:     idName,
		BankCardNo: bankCard,
	}

	if mobile != "" {
		reqType = "4"
		p.Mobile = mobile
	}

	paramMap, err := toMap(p, map[string]string{
		"key":            c.config.AuthID,
		"value":          c.config.AuthPWD,
		"rcaRequestType": reqType,
	})

	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &AuthResponse{}
	})

	if err != nil {
		return nil, err
	}

	rsp := ret.(*AuthResponse)

	if err = checkAuthErr(rsp.Code, rsp.Msg, rsp.Success); err != nil {
		return nil, err
	}

	info := struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{}

	if err = json.Unmarshal([]byte(rsp.Data), &info); err != nil {
		return nil, err
	}

	rsp.Message = info.Message
	rsp.Status = info.Status

	return rsp, nil
}

// AddUser imports user into YunHeTong service.
func (c *Client) AddUser(userID, phone, name, certNum string, userType string, certType string, autoSign bool) (*AddUserResponse, error) {
	createSign := "0"
	if autoSign {
		createSign = "1"
	}
	p := addUserParams{
		AppUserID:       userID,
		CellNum:         phone,
		UserType:        userType,
		UserName:        name,
		CertifyType:     certType,
		CertifyNumber:   certNum,
		CreateSignature: createSign,
	}

	paramMap, err := toMap(p, map[string]string{
		AppIDKey:    c.config.AppID,
		PasswordKey: c.config.Password,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &AddUserResponse{}
	})

	if err != nil {
		return nil, err
	}

	rsp := ret.(*AddUserResponse)

	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// ModifyPhoneNumber modifies user's cell phone number.
func (c *Client) ModifyPhoneNumber(phone, token string) (*ModifyPhoneNumberResponse, error) {
	p := modifyPhoneNumberParams{
		CellNum: phone,
	}
	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &ModifyPhoneNumberResponse{}
	})

	if err != nil {
		return nil, err
	}

	rsp := ret.(*ModifyPhoneNumberResponse)

	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// ModifyUserName modifies user's name.
func (c *Client) ModifyUserName(name, token string, autoSign bool) (*ModifyUserNameResponse, error) {
	var createSign string
	if autoSign {
		createSign = "1"
	}
	p := modifyUserNameParams{
		UserName:        name,
		CreateSignature: createSign,
	}
	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &ModifyUserNameResponse{}
	})

	if err != nil {
		return nil, err
	}

	rsp := ret.(*ModifyUserNameResponse)

	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// UserToken gets user's token string.
func (c *Client) UserToken(userID string) (*UserTokenResponse, error) {
	p := userTokenParams{
		AppUserID: userID,
	}
	paramMap, err := toMap(p, map[string]string{
		AppIDKey:    c.config.AppID,
		PasswordKey: c.config.Password,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &UserTokenResponse{}
	})

	if err != nil {
		return nil, err
	}

	rsp := ret.(*UserTokenResponse)

	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// CreateTemplateContract creates contract based on template.
func (c *Client) CreateTemplateContract(title, contractNo, templateID, token string, useCer bool, placeHolders M) (*CreateTemplateContractResponse, error) {
	var cer string
	if useCer {
		cer = "1"
	}

	data, err := json.Marshal(placeHolders)
	if err != nil {
		return nil, err
	}

	p := createTemplateContractParams{
		Title:         title,
		DefContractNo: contractNo,
		TemplateID:    templateID,
		UseCer:        cer,
		Param:         string(data),
	}

	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &CreateTemplateContractResponse{}
	})

	if err != nil {
		return nil, err
	}

	rsp := ret.(*CreateTemplateContractResponse)

	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// CreateFileContract creates contract by uploading file.
func (c *Client) CreateFileContract(title, contractNo, token string, useCer bool, data []byte) (*CreateFileContractResponse, error) {
	var cer string
	if useCer {
		cer = "1"
	}
	p := createFileContractParams{
		Title:         title,
		DefContractNo: contractNo,
		UseCer:        cer,
	}
	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, data, func() interface{} {
		return &CreateFileContractResponse{}
	})

	if err != nil {
		return nil, err
	}

	rsp := ret.(*CreateFileContractResponse)
	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// Partner represents a participant in a contract.
type Partner struct {
	AppUserID    string `json:"appUserId"`
	LocationName string `json:"locationName,omitempty"` // 模板签名占位符名称(与keyWord必填其一)
	Keyword      string `json:"keyWord,omitempty"`
}

// AddPartner adds partners of contract.
func (c *Client) AddPartner(contractID int64, token string, partners ...Partner) (*AddPartnerResponse, error) {
	data, err := json.Marshal(partners)
	if err != nil {
		return nil, err
	}

	p := addPartnerParams{
		ContractID: fmt.Sprintf("%d", contractID),
		Partners:   string(data),
	}

	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &AddPartnerResponse{}
	})
	if err != nil {
		return nil, err
	}

	rsp := ret.(*AddPartnerResponse)
	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// SignContract signs contract automatically.
func (c *Client) SignContract(contractID, token string, signers ...string) (*SignContractResponse, error) {
	data, err := json.Marshal(signers)
	if err != nil {
		return nil, err
	}

	p := signContractParams{
		ContractID: contractID,
		Signer:     string(data),
	}

	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &SignContractResponse{}
	})
	if err != nil {
		return nil, err
	}

	rsp := ret.(*SignContractResponse)
	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// InvalidateContract invalidates contract.
func (c *Client) InvalidateContract(contractID, token string) (*InvalidateContractResponse, error) {
	p := invalidateContractParams{
		ContractID: contractID,
	}

	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &InvalidateContractResponse{}
	})
	if err != nil {
		return nil, err
	}

	rsp := ret.(*InvalidateContractResponse)
	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// ListContracts returns a list of contracts finished or invalidated.
func (c *Client) ListContracts(pageNum, pageSize int, token string) (*ListContractsResponse, error) {
	p := listContractsParams{
		PageNum:  fmt.Sprintf("%d", pageNum),
		PageSize: fmt.Sprintf("%d", pageSize),
	}

	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &ListContractsResponse{}
	})
	if err != nil {
		return nil, err
	}

	rsp := ret.(*ListContractsResponse)
	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// LookupContractDetail returns the detail of a contract.
func (c *Client) LookupContractDetail(contractID, token string) (*LookupContractDetailResponse, error) {
	p := lookupContractDetailParams{
		ContractID: contractID,
	}

	paramMap, err := toMap(p, map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	ret, err := httpRequest(c, p.URI(), paramMap, nil, func() interface{} {
		return &LookupContractDetailResponse{}
	})
	if err != nil {
		return nil, err
	}

	rsp := ret.(*LookupContractDetailResponse)
	if err = checkErr(rsp.Code, rsp.SubCode, rsp.Message); err != nil {
		return nil, err
	}

	return rsp, nil
}

// DownloadContract downloads a contract.
func (c *Client) DownloadContract(contractID, token string) ([]byte, error) {
	p := downloadContractParams{
		ContractID: contractID,
	}

	paramMap, err := toMap(p, nil)
	if err != nil {
		return nil, err
	}

	vals := url.Values{}
	for k, v := range paramMap {
		vals.Add(k, v)
	}

	uri := fmt.Sprintf("%s?token=%s&contractId=%s", p.URI(), token, contractID)
	apiURL := fmt.Sprintf("%s%s", c.config.APIGateway, uri)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	rsp, err := c.tlsClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// AsyncNotifyResult represents the result returned from YunHeTong service.
type AsyncNotifyResult struct {
	Content      string                 `json:"content"`
	NoticeType   int                    `json:"noticeType"`
	NoticeParams string                 `json:"noticeParams"`
	InfoMap      map[string]interface{} `json:"map"`
}

// AsyncNotify returns asynchronous notification from YunHeTong service.
func (c *Client) AsyncNotify(req *http.Request) (*AsyncNotifyResult, error) {
	defer req.Body.Close()
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	tmp, err := url.QueryUnescape(string(bodyBytes)) // url decode as YHT notification is url encoded
	if err != nil {
		return nil, err
	}
	jsonStr := strings.Replace(tmp, "notice=", "", -1) // remove "notice=" segment as it is not json format
	result := &AsyncNotifyResult{}
	err = json.Unmarshal([]byte(jsonStr), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AnswerAsyncNotify returns a json string answering async notification.
func (c *Client) AnswerAsyncNotify(rsp bool, msg string) string {
	ret := map[string]interface{}{
		"response": rsp,
		"msg":      msg,
	}
	data, err := json.Marshal(ret)
	if err != nil {
		return ""
	}
	return string(data)
}

// httpRequestV4 云合同V4版本接口请求
func httpRequestV4(c *Client, token, uri, method string, jsonData []byte, factory func() interface{}) (interface{}, string, error) {
	apiURL := c.config.APIGateway + uri
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, "", err
	}
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	if "" != token {
		req.Header.Add("token", token)
	}

	yhtResp, err := c.tlsClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	llt := ""
	if uri == "/auth/login" {
		if v, ok := yhtResp.Header["Token"]; ok { // 取token
			if ok {
				llt = v[0]
			}
		}
	}

	defer yhtResp.Body.Close()

	data, err := ioutil.ReadAll(yhtResp.Body)
	if err != nil {
		return nil, "", err
	}
	rsp := factory()
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, "", err
	}

	return rsp, llt, nil
}

func httpRequest(c *Client, uri string, paramMap map[string]string, fileData []byte, factory func() interface{}) (interface{}, error) {
	if token, ok := paramMap["token"]; ok {
		delete(paramMap, "token")
		uri = fmt.Sprintf("%s?token=%s", uri, token)
	}
	apiURL := fmt.Sprintf("%s%s", c.config.APIGateway, uri)
	if strings.Contains(uri, "authentic") {
		apiURL = fmt.Sprintf("%s%s", c.config.AuthGateway, uri)
	}

	var data []byte
	var err error
	if fileData != nil {
		data, err = c.doMultipartRequest(apiURL, paramMap, fileData)
	} else {
		data, err = c.doHTTPRequest(apiURL, paramMap)
	}

	if err != nil {
		return nil, err
	}

	rsp := factory()
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *Client) doMultipartRequest(apiURL string, paramMap map[string]string, fileData []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	for k, v := range paramMap {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	fw, err := writer.CreateFormField("file")
	if err != nil {
		return nil, err
	}
	if _, err = fw.Write(fileData); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	rsp, err := c.tlsClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) doHTTPRequest(apiURL string, paramMap map[string]string) ([]byte, error) {
	formData := url.Values{}
	for k, v := range paramMap {
		formData.Add(k, v)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	rsp, err := c.tlsClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
