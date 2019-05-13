package goyht

import (
	"encoding/json"
	"net/http"
)

// YhtBaseResp 云合同基础应答模型
type YhtBaseResp struct {
	Code   int             `json:"code"`
	RawMsg json.RawMessage `json:"msg"`
}

// Message .
func (p YhtBaseResp) Message() string {
	if 200 != p.Code {
		return string(p.RawMsg)
	}
	return "请求成功"
}

type yhtAuthLoginReq struct {
	AppID    string `json:"appId"`    // 应用ID
	AppKey   string `json:"appKey"`   // 应用密钥
	SignerID string `json:"signerId"` // 用户ID，可选参数，不传则获取平台的长效令牌，否则获取指定用户的长效令牌
}

func (p yhtAuthLoginReq) URI() string {
	return "/auth/login"
}

func (p yhtAuthLoginReq) Method() string {
	return http.MethodPost
}

// 个人用户身份地区类型（V4版本）
const (
	YHTIdentityRegionMainland = "0" // 大陆
	YHTIdentityRegionHK       = "1" // 香港
	YHTIdentityRegionTaiwan   = "2" // 台湾
	YHTIdentityRegionMacao    = "3" // 澳门
	YHTIdentityRegionForeign  = "4" // 海外
)

// 个人用户证件类型（V4版本）
const (
	YHTPersonCertTypeIDCard   = "a" // 身份证
	YHTPersonCertTypePassport = "b" // 护照
	YHTPersonCertTypeEEP      = "d" // 港澳通行证
	YHTPersonCertTypeMTPForTW = "e" // 台胞证
	YHTPersonCertTypeMTPForHM = "f" // 港澳居民来往内地通行证
	YHTPersonCertTypeOther    = "z" // 其它
)

// 个人用户手机号地区类型（V4版本）
const (
	YHTPhoneRegionMainland = "0" // 大陆
	YHTPhoneRegionHKMacao  = "1" // 香港澳门
	YHTPhoneRegionTaiwan   = "2" // 台湾
)

// YhtCreatePersonReq 云合同创建个人用户请求模型
type YhtCreatePersonReq struct {
	Username       string `json:"userName"`
	IdentityRegion string `json:"identityRegion"`
	CertType       string `json:"certifyType"`
	CertNum        string `json:"certifyNum"`
	PhoneRegion    string `json:"phoneRegion"`
	Phone          string `json:"phoneNo"`
	CAType         string `json:"caType"` // 固定传B2
}

// SignerIDResp 用户ID应答模型
type SignerIDResp struct {
	SignerID int `json:"signerId"`
}

// YhtCreateUserResp 云合同创建个人用户应答
type YhtCreateUserResp struct {
	YhtBaseResp
	Data SignerIDResp `json:"data"`
}

// URI .
func (p YhtCreatePersonReq) URI() string {
	return "/user/person"
}

// Method .
func (p YhtCreatePersonReq) Method() string {
	return http.MethodPost
}

// YhtCreateCompanyReq 云合同创建企业用户请求
type YhtCreateCompanyReq struct {
	Username string `json:"userName"`
	CertType string `json:"certifyType"` // 固定为1， 社会统一信用代码
	CertNum  string `json:"certifyNum"`
	Phone    string `json:"phoneNo"`
	CAType   string `json:"caType"` // 固定传B2
}

// URI .
func (p YhtCreateCompanyReq) URI() string {
	return "/user/company"
}

// Method .
func (p YhtCreateCompanyReq) Method() string {
	return http.MethodPost
}

// YhtQuerySignerIDReq 云合同查询用户ID请求
type YhtQuerySignerIDReq struct {
	CertifyNumList []string `json:"certifyNumList"`
}

// URI .
func (p YhtQuerySignerIDReq) URI() string {
	return "/user/signerId/certifyNums"
}

// Method .
func (p YhtQuerySignerIDReq) Method() string {
	return http.MethodPost
}

// YhtQuerySignerIDResp 云合同查询用户ID应答
type YhtQuerySignerIDResp struct {
	YhtBaseResp
	Data []map[string]int `json:"data"`
}

// 云合同个人印章边框类型
const (
	YHTPMWithBorder    = "B1" // 有边框
	YHTPMWithoutBorder = "B2" // 无边框
)

// 云合同个人印章字体类型
const (
	YHTPMFontKaiti = "F1" // 楷体
	YHTPMFontHWFS  = "F2" // 华文仿宋
	YHTPMFontHWKT  = "F3" // 华文楷体
	YHTPMFontMSYH  = "F4" // 微软雅黑
)

// 云合同个人印章字体颜色类型
const (
	YHTMFontColorRed   = "C1" // 红
	YHTMFontColorBlue  = "C2" // 蓝
	YHTMFontColorBlack = "C3" // 黑
)

// 云合同个人印章模式
const (
	YHTMModeNormal      = "0" // 常规
	YHTMModeTransparent = "1" // 透明
	YHTMModeSafe        = "2" // 脱敏
)

// 云合同个人印章缩放类型
const (
	YHTPMZoomCodeLarge  = "0" // 大
	YHTPMZoomCodeNormal = "1" // 中
	YHTPMZoomCodeSmall  = "2" // 小
)

// YhtCreatePersonMoulageReq 云合同创建个人印章请求
type YhtCreatePersonMoulageReq struct {
	SignerID   string `json:"signerId"`
	BorderType string `json:"borderType"`
	FontFamily string `json:"fontFamily"`
	FontColor  string `json:"color"`
	Mode       string `json:"mode"`
	ZoomCode   string `json:"zoomCode"`
}

// URI .
func (p YhtCreatePersonMoulageReq) URI() string {
	return "/user/personMoulage"
}

// Method .
func (p YhtCreatePersonMoulageReq) Method() string {
	return http.MethodPost
}

// 云合同企业印章形状
const (
	YHTCMStyleTypeCircle  = "1" // 圆形
	YHTCMStyleTypeEllipse = "2" // 椭圆
)

// YhtCreateCompanyMoulageReq 云合同创建企业印章请求
type YhtCreateCompanyMoulageReq struct {
	SignerID    string `json:"signerId"`
	StyleType   string `json:"styleType"`
	TextContent string `json:"textContent"` // 横向文案
	KeyContent  string `json:"keyContent"`  // 防伪码，13位数字
	FontColor   string `json:"color"`
	Mode        string `json:"mode"`
}

// URI .
func (p YhtCreateCompanyMoulageReq) URI() string {
	return "/user/companyMoulage"
}

// Method .
func (p YhtCreateCompanyMoulageReq) Method() string {
	return http.MethodPost
}

// MoulageIDResp .
type MoulageIDResp struct {
	MoulageID int `json:"moulageId"`
}

// YhtCreateMoulageResp 云合同创建印章应答
type YhtCreateMoulageResp struct {
	YhtBaseResp
	Data MoulageIDResp `json:"data"`
}

// YhtCreateTemplateContractReq 云合同根据模板生成合同请求
type YhtCreateTemplateContractReq struct {
	Title        string      `json:"contractTitle"` // 合同标题
	ContractNo   string      `json:"contractNo"`    // 自定义合同编号
	TemplateID   string      `json:"templateId"`    // 模板ID
	ContractData interface{} `json:"contractData"`  // 合同参数
}

// URI .
func (p YhtCreateTemplateContractReq) URI() string {
	return "/contract/templateContract"
}

// Method .
func (p YhtCreateTemplateContractReq) Method() string {
	return http.MethodPost
}

// ContractIDResp .
type ContractIDResp struct {
	ContractID int `json:"contractId"`
}

// YhtCreateTemplateContractResp 根据模板生成合同应答
type YhtCreateTemplateContractResp struct {
	YhtBaseResp
	Data ContractIDResp `json:"data"`
}

// 合同ID类型
const (
	YHTIDTypeSystem = "0" // 云合同平台合同ID
	YHTIDTypeCustom = "1" // 第三方平台自定义合同ID
)

// 签署定位方式
const (
	YHTSignPositionTypeKeyWord     = "0" // 关键字定位
	YHTSignPositionTypePlaceHolder = "1" // 占位符定位
	YHTSignPositionTypeCoord       = "2" // 坐标定位
)

// 签署验证方式
const (
	YHTSignValidateTypeIgnore = "0" // 不校验
	YHTSignValidateTypeSMS    = "1" // 短信验证
)

// 印章使用类型
const (
	YHTSignModeSpecify = "0" // 指定印章
	YHTSignModeRender  = "1" // 每次绘制
)

// 合同签署形态
const (
	YHTSignFormJS = "0" // JS集成页面
	YHTSignFormH5 = "1" // 独立H5页面
)

// YhtSigner .
type YhtSigner struct {
	SignerID         string `json:"signerId"`
	SignPositionType string `json:"signPositionType"` // 签署的定位方式：0=关键字定位，1=签名占位符定位，2=签署坐标
	PositionContent  string `json:"positionContent"`
	SignValidateType string `json:"signValidateType"` // 签署验证方式：0=不校验，1=短信验证
	SignMode         string `json:"signMode"`         // 印章使用类型（针对页面签署）：0=指定印章，1=每次绘制
	SignForm         string `json:"signForm"`         // 签署形态，0=JS集成页面(默认)，1=独立H5页面
}

// YhtAddSignerReq 云合同添加签署者请求
type YhtAddSignerReq struct {
	IDType    string      `json:"idType"`    // 合同ID类型，0 合同ID，1 合同自定义编号
	IDContent string      `json:"idContent"` // 合同ID
	Signers   []YhtSigner `json:"signers"`   // 签署者列表
}

// URI .
func (p YhtAddSignerReq) URI() string {
	return "/contract/signer"
}

// Method .
func (p YhtAddSignerReq) Method() string {
	return http.MethodPost
}

// 签章样式
const (
	YHTSealClassNormal       = "0" // 常规
	YHTSealClassQF           = "1" // 骑缝
	YHTSealClassWithAbstract = "2" // 含摘要
	YHTSealClassWithSignTime = "3" // 含签署时间
	YHTSealClassNormalWithQF = "4" // 常规+骑缝
)

// YhtSignContractReq 云合同签署合同请求
type YhtSignContractReq struct {
	IDType    string `json:"idType"`
	IDContent string `json:"idContent"`
	SignerID  string `json:"signerId"`
	MoulageID string `json:"moulageId"`
	SealClass string `json:"sealClass"` // 签章样式，0=常规样式，1=骑缝章，2=含摘要样式，3=含签署时间样式，4=常规样式+骑缝章，可选参数，不传时使用常规样式
}

// URI .
func (p YhtSignContractReq) URI() string {
	return "/contract/sign"
}

// Method .
func (p YhtSignContractReq) Method() string {
	return http.MethodPost
}

// AuthSerialNumResp 实名认证流水号
type AuthSerialNumResp struct {
	ID string `json:"id"`
}

// AuthRealNameResp 实名认证应答
type AuthRealNameResp struct {
	YhtBaseResp
	Data AuthSerialNumResp `json:"data"`
}

type authParams struct {
	IDNo       string `param:"idNo"`
	IDName     string `param:"idName"`
	BankCardNo string `param:"bankCardNo"`
	Mobile     string `param:"mobile"`
}

func (p authParams) URI() string {
	return "/authentic/authentication"
}

// AuthResponse represents the response returned.
type AuthResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    string `json:"data"`
	Message string
	Status  string
}

type addUserParams struct {
	AppUserID       string `param:"appUserId"`
	CellNum         string `param:"cellNum"`
	UserType        string `param:"userType"`
	UserName        string `param:"userName"`
	CertifyType     string `param:"certifyType"`
	CertifyNumber   string `param:"certifyNumber"`
	CreateSignature string `param:"createSignature"`
}

func (p addUserParams) URI() string {
	return "/userInfo/addUser"
}

// AddUserResponse represents the reponse returned.
type AddUserResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
}

type modifyPhoneNumberParams struct {
	CellNum string `param:"cellNum"`
}

// URI returns the URL of API.
func (p modifyPhoneNumberParams) URI() string {
	return "/userInfo/modifyCellNum"
}

// ModifyPhoneNumberResponse represents the reponse returned.
type ModifyPhoneNumberResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
}

type modifyUserNameParams struct {
	UserName        string `json:"userName"`
	CreateSignature string `json:"createSignature"`
}

func (p modifyUserNameParams) URI() string {
	return "/userInfo/modifyUserName"
}

// ModifyUserNameResponse represents the reponse returned.
type ModifyUserNameResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
}

type userTokenParams struct {
	AppUserID string `param:"appUserId"`
}

func (p userTokenParams) URI() string {
	return "/token/getToken"
}

// UserTokenResponse represents the reponse returned.
type UserTokenResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
	Value   struct {
		Token string `json:"token"`
	} `json:"value"`
}

type createTemplateContractParams struct {
	Title         string `param:"title"`
	DefContractNo string `param:"defContractNo"`
	TemplateID    string `param:"templateId"`
	UseCer        string `param:"useCer"`
	Param         string `param:"param"`
}

// URI returns the URL of API.
func (p createTemplateContractParams) URI() string {
	return "/contract/templateContract"
}

// CreateTemplateContractResponse represents the reponse returned.
type CreateTemplateContractResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
	Value   struct {
		ContractID int64 `json:"contractId"`
	} `json:"value"`
}

type createFileContractParams struct {
	Title         string `param:"title"`
	DefContractNo string `param:"defContractNo"`
	UseCer        string `param:"useCer"`
}

func (p createFileContractParams) URI() string {
	return "/contract/fileContract"
}

// CreateFileContractResponse represents the reponse returned.
type CreateFileContractResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
	Value   struct {
		ContractID string `json:"contractId"`
	} `json:"value"`
}

type addPartnerParams struct {
	ContractID string `param:"contractId"`
	Partners   string `param:"partners"`
}

func (p addPartnerParams) URI() string {
	return "/contract/addPartner"
}

// AddPartnerResponse represents the reponse returned.
type AddPartnerResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
}

type signContractParams struct {
	ContractID string `param:"contractId"`
	Signer     string `param:"signer"`
}

func (p signContractParams) URI() string {
	return "/contract/signContract"
}

// SignContractResponse represents the reponse returned.
type SignContractResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
}

type invalidateContractParams struct {
	ContractID string `param:"contractId"`
}

func (p invalidateContractParams) URI() string {
	return "/contract/invalid"
}

// InvalidateContractResponse represents the reponse returned.
type InvalidateContractResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
}

type listContractsParams struct {
	PageNum  string `param:"pageNum"`
	PageSize string `param:"pageSize"`
}

func (p listContractsParams) URI() string {
	return "/contract/list"
}

// ListContractsResponse represents the reponse returned.
type ListContractsResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
	Value   struct {
		ContractList []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Status      string `json:"status"`
			AppName     string `json:"appName"`
			GmtModify   string `json:"gmtModify"`
			PartnerList string `json:"partnerList"`
		} `json:"contractList"`
	} `json:"value"`
}

type lookupContractDetailParams struct {
	ContractID string `param:"contractId"`
}

func (p lookupContractDetailParams) URI() string {
	return "/contract/detail"
}

// LookupContractDetailResponse represents the reponse returned.
type LookupContractDetailResponse struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subCode"`
	Message string `json:"message"`
	Value   struct {
		PartnerList []struct {
			SignStatus string `json:"signStatus"`
			UserID     string `json:"userId"`
		} `json:"partnerList"`
		Title  string `param:"title"`
		Status string `json:"status"`
	} `json:"value"`
}

type downloadContractParams struct {
	ContractID string `param:"contractId"`
}

func (p downloadContractParams) URI() string {
	return "/contract/download"
}

// DownloadContractResponse represents the reponse returned.
type DownloadContractResponse struct {
	File []byte
}
