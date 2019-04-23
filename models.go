package bankid

// Request - A basic BankID request contain one or more of the variables below
import "fmt"

// Collect response statuses
const (
	OrderPending  = "pending"
	OrderFailed   = "failed"
	OrderComplete = "complete"
)

// HintCodes - for Pending and Failed statuses
const (
	PendOutstandingTransaction = "outstandingTransaction"
	PendNoClient               = "noClient"
	PendStarted                = "started"
	PendUserSign               = "userSign"
	// PendUnknown

	FailExpiredTransaction = "expiredTransaction"
	FailCertificateErr     = "certificateErr"
	FailUserCancel         = "userCancel"
	FailCancelled          = "cancelled"
	FailStartFailed        = "startFailed"
	// FailUnknown
)

type Request struct {
	OrderRef           string `json:"orderRef,omitempty"`
	EndUserIP          string `json:"endUserIp,omitempty"`
	PersonalNumber     string `json:"personalNumber,omitempty"`
	UserVisibleData    string `json:"userVisibleData,omitempty"`
	UserNonVisibleData string `json:"userNonVisibleData,omitempty"`
}

// Response - for Auth and Sign requests
type Response struct {
	AutoStartToken string `json:"autoStartToken"` // UUID, e.g "dbbee61c-357b-4fd8-b103-392eed10be7a"
	OrderRef       string `json:"orderRef"`       // UUID, e.g "131daac9-16c6-4618-beb0-365768f37288"
}

// ErrorResponse - when anything goes bad
type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Details   string `json:"details"`
}

// Error -
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("failed with code: %s. '%s'", e.ErrorCode, e.Details)
}

type CollectResponse struct {
	OrderRef       string      `json:"orderRef"`
	Status         string      `json:"status"`
	HintCode       string      `json:"hintCode,omitempty"`       // Pending and Failed orders only
	CompletionData *Completion `json:"completionData,omitempty"` // Complete orders only
}

type Completion struct {
	User         User   `json:"user"`
	Device       Device `json:"device"`
	Cert         Cert   `json:"cert"`
	Signature    string `json:"signature"` // base64 encoded signature, see https://www.bankid.com/bankid-i-dina-tjanster/rp-info
	OCSPResponse string `json:"ocspResponse"`
}

type User struct {
	PersonalNumber string `json:"personalNumber"` // e.g "197001010000"
	Name           string `json:"name"`
	GivenName      string `json:"giveName"`
	Surname        string `json:"surname"`
}

type Device struct {
	IPAddress string `json:"ipAddress"` // e.g "192.168.0.1"
}

type Cert struct {
	NotBefore string `json:"notBefore"` // e.g "1502983274000" UNIX Epoch in ms
	NotAfter  string `json:"notAfter"`
}
