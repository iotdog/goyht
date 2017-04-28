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
)

// M is a convenient alias for a map[string]interface{} map.
type M map[string]interface{}

// UserType is a type for user types.
type UserType int

// constants about user types.
const (
	UserTypePersonal   UserType = 1 // 个人
	UserTypeEnterprise UserType = 2 // 企业
	UserTypePlatform   UserType = 4 // 平台
)

// CertType is a type for certification types.
type CertType int

// constants about certification types.
const (
	CertTypeIDCard   CertType = 1 // 身份证
	CertTypePassport CertType = 2 // 护照
	CertTypeOfficer  CertType = 3 // 军官证
	CertTypeLicence  CertType = 4 // 营业执照
	CertTypeOrgan    CertType = 5 // 组织机构代码
	CertTypeSocial   CertType = 6 // 社会代码
)

// constants for url and keys
const (
	YHTAPIGateway = "https://sdk.yunhetong.com/sdk"
	AppIDKey      = "appId"
	PasswordKey   = "passWord"
)

// Config contains configurations about YunHeTong service.
type Config struct {
	AppID      string
	Password   string
	APIGateway string
}

// Client handles all APIs for YunHeTong service.
type Client struct {
	config    Config
	tlsClient http.Client
}

// NewClient returns a *Client.
func NewClient(cfg Config) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := http.Client{Transport: tr}
	return &Client{
		config:    cfg,
		tlsClient: client,
	}
}

func httpRequest(c *Client, uri string, paramMap map[string]string, fileData []byte, factory func() interface{}) (interface{}, error) {
	var data []byte
	var err error
	if fileData != nil {
		data, err = c.doMultipartRequest(uri, paramMap, fileData)
	} else {
		data, err = c.doHTTPRequest(uri, paramMap)
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

// AddUser imports user into YunHeTong service.
func (c *Client) AddUser(userID, phone, name, certNum string, userType UserType, certType CertType, autoSign bool) (*AddUserResponse, error) {
	var createSign string
	if autoSign {
		createSign = "1"
	}
	p := addUserParams{
		AppUserID:       userID,
		CellNum:         phone,
		UserType:        fmt.Sprintf("%d", userType),
		UserName:        name,
		CertifyType:     fmt.Sprintf("%d", certType),
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
func (c *Client) AddPartner(contractID, token string, partners ...Partner) (*AddPartnerResponse, error) {
	data, err := json.Marshal(partners)
	if err != nil {
		return nil, err
	}

	p := addPartnerParams{
		ContractID: contractID,
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
func (c *Client) InvalidateContract() {}

// ListContracts returns a list of contracts finished or invalidated.
func (c *Client) ListContracts() {}

// LookupContractDetail returns the detail of a contract.
func (c *Client) LookupContractDetail() {}

// DownloadContract downloads a contract.
func (c *Client) DownloadContract() {}

// AsyncNotifyResult represents the result returned from YunHeTong service.
type AsyncNotifyResult struct{}

// AsyncNotify returns asynchronous notification from YunHeTong service.
func (c *Client) AsyncNotify(req *http.Request) (*AsyncNotifyResult, error) {
	return nil, errors.New("not defined")
}

func (c *Client) doMultipartRequest(uri string, paramMap map[string]string, fileData []byte) ([]byte, error) {
	if token, ok := paramMap["token"]; ok {
		delete(paramMap, "token")
		uri = fmt.Sprintf("%s?token=%s", uri, token)
	}
	apiURL := fmt.Sprintf("%s%s", c.config.APIGateway, uri)

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
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

func (c *Client) doHTTPRequest(uri string, paramMap map[string]string) ([]byte, error) {
	if token, ok := paramMap["token"]; ok {
		delete(paramMap, "token")
		uri = fmt.Sprintf("%s?token=%s", uri, token)
	}
	apiURL := fmt.Sprintf("%s%s", c.config.APIGateway, uri)

	formData := url.Values{}
	for k, v := range paramMap {
		formData.Add(k, v)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
