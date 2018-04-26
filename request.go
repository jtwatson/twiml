package twiml

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"

	"github.com/pkg/errors"
)

// Request is a twillio request expecting a TwiML response
type Request struct {
	r *http.Request
}

// NewRequest returns Request
func NewRequest(r *http.Request) *Request {
	return &Request{r: r}
}

// ValidatePost validates the Twilio Signature, requiring that the request is a POST
func (req *Request) ValidatePost(authToken string) error {
	if req.r.Method != "POST" {
		return fmt.Errorf("twiml.Request.ValidatePost(): Expected a POST request, received %s", req.r.Method)
	}

	if err := req.r.ParseForm(); err != nil {
		return errors.WithMessage(err, "twiml.Request.ValidatePost(): http.Request.ParseForm():")
	}

	params := make([]string, 0, len(req.r.PostForm))
	for p := range req.r.PostForm {
		params = append(params, p)
	}
	sort.Sort(sort.StringSlice(params))

	message := req.r.URL.String()
	for _, p := range params {
		message += p
		message += req.r.PostForm[p][0]
	}

	hash := hmac.New(sha1.New, []byte(authToken))
	if n, err := hash.Write([]byte(message)); err != nil || n != len(message) {
		return err
	}
	sig := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if xTwilioSigHdr := req.r.Header[http.CanonicalHeaderKey("X-Twilio-Signature")]; len(xTwilioSigHdr) != 1 || sig != xTwilioSigHdr[0] {
		var xTwilioSig string
		if len(xTwilioSigHdr) == 1 {
			xTwilioSig = xTwilioSigHdr[0]
		}
		return fmt.Errorf("twiml.Request.ValidatePost(): Calculated Signature: %s, failed to match X-Twilio-Signature: %s", sig, xTwilioSig)
	}
	return nil
}

func (req *Request) form(key string) string {
	value, exists := req.r.PostForm[key]
	if exists {
		return value[0]
	}

	return ""
}

func (req *Request) CallSid() string {
	return req.form("CallSid")
}

func (req *Request) AccountSid() string {
	return req.form("AccountSid")
}

func (req *Request) From() string {
	return req.form("From")
}

func (req *Request) To() string {
	return req.form("To")
}

func (req *Request) CallStatus() string {
	return req.form("CallStatus")
}

func (req *Request) ApiVersion() string {
	return req.form("ApiVersion")
}

func (req *Request) ForwardedFrom() string {
	return req.form("ForwardedFrom")
}

func (req *Request) CallerName() string {
	return req.form("CallerName")
}

func (req *Request) ParentCallSid() string {
	return req.form("ParentCallSid")
}

func (req *Request) FromCity() string {
	return req.form("FromCity")
}

func (req *Request) FromState() string {
	return req.form("FromState")
}

func (req *Request) FromZip() string {
	return req.form("FromZip")
}

func (req *Request) FromCountry() string {
	return req.form("FromCountry")
}

func (req *Request) ToCity() string {
	return req.form("ToCity")
}

func (req *Request) ToZip() string {
	return req.form("ToZip")
}

func (req *Request) ToCountry() string {
	return req.form("ToCountry")
}

func (req *Request) SipDomain() string {
	return req.form("SipDomain")
}

func (req *Request) SipUsername() string {
	return req.form("SipUsername")
}

func (req *Request) SipCallId() string {
	return req.form("SipCallId")
}

func (req *Request) SipSourceIp() string {
	return req.form("SipSourceIp")
}

func (req *Request) Digits() string {
	return req.form("Digits")
}
