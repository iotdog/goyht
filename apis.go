package goyht

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
