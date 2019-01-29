package twiml

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// RequestValues hold form values from a validated Request
type RequestValues map[string]string

// CallDuration Parses the duration rom the string value
func (r RequestValues) CallDuration() (int, error) {
	var duration int
	if r["CallDuration"] != "" {
		d, err := strconv.Atoi(r["CallDuration"])
		if err != nil {
			return 0, errors.WithMessage(err, "RequestValues.CallDuration()")
		}
		duration = d
	}
	return duration, nil
}

// SequenceNumber Parses the duration rom the string value
func (r RequestValues) SequenceNumber() (int, error) {
	var seq int
	if r["SequenceNumber"] != "" {
		d, err := strconv.Atoi(r["SequenceNumber"])
		if err != nil {
			return 0, errors.WithMessage(err, "RequestValues.SequenceNumber()")
		}
		seq = d
	}
	return seq, nil
}

// TimestampOrNow parses the Timestamp from string. If Timestamp does not exist in the
// current request, time.Now() is returned instead.
func (r RequestValues) TimestampOrNow() time.Time {
	t, err := time.Parse(time.RFC1123Z, r["Timestamp"])
	if err != nil {
		t = time.Now()
	}
	return t
}

// Request is a twillio request expecting a TwiML response
type Request struct {
	host   string
	r      *http.Request
	Values RequestValues
}

// NewRequest returns Request
func NewRequest(host string, r *http.Request) *Request {
	return &Request{host: host, r: r, Values: RequestValues{}}
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

	message := req.host + req.r.URL.String()
	for _, p := range params {
		message += p
		if len(req.r.PostForm[p]) > 0 {
			message += req.r.PostForm[p][0]
		}
	}

	hash := hmac.New(sha1.New, []byte(authToken))
	if n, err := hash.Write([]byte(message)); err != nil {
		return errors.WithMessage(err, "twiml.Request.ValidatePost(): hash.Write()")
	} else if n != len(message) {
		err := fmt.Errorf("expected %d bytes, got %d bytes", len(message), n)
		return errors.WithMessage(err, "twiml.Request.ValidatePost(): hash.Write()")
	}
	sig := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if xTwilioSigHdr := req.r.Header[http.CanonicalHeaderKey("X-Twilio-Signature")]; len(xTwilioSigHdr) != 1 || sig != xTwilioSigHdr[0] {
		var xTwilioSig string
		if len(xTwilioSigHdr) == 1 {
			xTwilioSig = xTwilioSigHdr[0]
		}
		return fmt.Errorf("twiml.Request.ValidatePost(): Calculated Signature: %s, failed to match X-Twilio-Signature: %s", sig, xTwilioSig)
	}

	// Validate data
	for _, p := range params {
		var val string
		if len(req.r.PostForm[p]) > 0 {
			val = req.r.PostForm[p][0]
		}
		if valParam, ok := fieldValidators[p]; ok {
			if err := valParam.valFunc(val, valParam.valParam); err != nil {
				return fmt.Errorf("twiml.Request.ValidatePost(): Invalid form value: %s=%s, err: %s", p, val, err)
			}
		}
		req.Values[p] = val
	}

	return nil
}

type valCfg struct {
	valFunc  func(interface{}, string) error
	valParam string
}

var fieldValidators = map[string]valCfg{
	// "CallSid":       "CallSid",
	// "AccountSid":    "AccountSid",
	"From": valCfg{valFunc: validPhoneNumber},
	"To":   valCfg{valFunc: validPhoneNumber},
	// "CallStatus":    "CallStatus",
	// "ApiVersion":    "ApiVersion",
	// "ForwardedFrom": "ForwardedFrom",
	// "CallerName":    "CallerName",
	// "ParentCallSid": "ParentCallSid",
	// "FromCity":      "FromCity",
	// "FromState":     "FromState",
	// "FromZip":       "FromZip",
	// "FromCountry":   "FromCountry",
	// "ToCity":        "ToCity",
	// "ToZip":         "ToZip",
	// "ToCountry":     "ToCountry",
	// "SipDomain":     "SipDomain",
	// "SipUsername":   "SipUsername",
	// "SipCallId":     "SipCallId",
	// "SipSourceIp":   "SipSourceIp",
	"Digits": valCfg{valFunc: validNumeric},
}
