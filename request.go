// package twiml implements tooling to create, parse, and validate Twilio Twiml requests and responses
package twiml

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/errors/v5"
	"github.com/ttacon/libphonenumber"
	"go.opencensus.io/trace"
)

// RequestValues hold form values from a validated Request
type RequestValues map[string]string

// CallDuration Parses the duration from the string value
func (r RequestValues) CallDuration() (time.Duration, error) {
	var duration int
	if r["CallDuration"] != "" {
		d, err := strconv.Atoi(r["CallDuration"])
		if err != nil {
			return 0, errors.Wrap(err, "RequestValues.CallDuration()")
		}
		duration = d
	}

	return time.Second * time.Duration(duration), nil
}

// SequenceNumber Parses the sequence number from the string value
func (r RequestValues) SequenceNumber() (int, error) {
	var seq int
	if r["SequenceNumber"] != "" {
		d, err := strconv.Atoi(r["SequenceNumber"])
		if err != nil {
			return 0, errors.Wrap(err, "RequestValues.SequenceNumber()")
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

// From returns a Number parsed from the raw From value
func (r RequestValues) From() *ParsedNumber {
	return ParseNumber(r["From"])
}

// To returns a Number parsed from the raw To value
func (r RequestValues) To() *ParsedNumber {
	return ParseNumber(r["To"])
}

// ParseNumber parses ether a E164 number or a SIP URI returning a ParsedNumber
func ParseNumber(v string) *ParsedNumber {
	number := &ParsedNumber{
		Number: v,
		Raw:    v,
	}

	if err := validPhoneNumber(v, ""); err == nil {
		number.Valid = true

		return number
	}

	u, err := parseSIPURI(v)
	if err != nil {
		return number
	}

	parts := strings.Split(u.Hostname(), ".")
	l := len(parts)

	if l < 5 || strings.Join(parts[l-2:], ".") != "twilio.com" || parts[l-4] != "sip" {
		return number
	}

	num, err := FormatNumber(u.User.Username())
	if err != nil {
		return number
	}

	number.Valid = true
	number.SIP = true
	number.Number = num
	number.SIPDomain = strings.Join(parts[:l-4], ".")
	number.Region = parts[l-4]

	return number
}

type ParsedNumber struct {
	Valid     bool
	Number    string
	SIP       bool
	SIPDomain string
	Region    string
	Raw       string
}

// FormatNumber formates a number to E164 format
func FormatNumber(number string) (string, error) {
	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return number, errors.Wrapf(errors.New("Invalid phone number"), "twiml.FormatNumber(): %s", err)
	}
	if !libphonenumber.IsValidNumber(num) {
		return number, errors.Wrap(errors.New("Invalid phone number"), "twiml.FormatNumber()")
	}

	return libphonenumber.Format(num, libphonenumber.E164), nil
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
func (req *Request) ValidatePost(ctx context.Context, authToken string) error {
	_, span := trace.StartSpan(ctx, "twiml.Request.ValidatePost()")
	defer span.End()

	url := req.host + req.r.URL.String()
	span.AddAttributes(trace.StringAttribute("url", url))

	if req.r.Method != "POST" {
		return errors.Wrap(fmt.Errorf("expected a POST request, received %s", req.r.Method), "twiml.Request.ValidatePost()")
	}

	if err := req.r.ParseForm(); err != nil {
		return errors.Wrap(err, "http.Request.ParseForm()")
	}

	params := make([]string, 0, len(req.r.PostForm))
	for p := range req.r.PostForm {
		params = append(params, p)
	}
	sort.Strings(params)

	message := url
	for _, p := range params {
		message += p
		if len(req.r.PostForm[p]) > 0 {
			message += req.r.PostForm[p][0]
		}
	}

	hash := hmac.New(sha1.New, []byte(authToken))
	if n, err := hash.Write([]byte(message)); err != nil {
		return errors.Wrap(err, "twiml.Request.ValidatePost(): hash.Write()")
	} else if n != len(message) {
		err := fmt.Errorf("expected %d bytes, got %d bytes", len(message), n)

		return errors.Wrap(err, "twiml.Request.ValidatePost(): hash.Write()")
	}
	sig := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if xTwilioSigHdr := req.r.Header[http.CanonicalHeaderKey("X-Twilio-Signature")]; len(xTwilioSigHdr) != 1 || sig != xTwilioSigHdr[0] {
		var xTwilioSig string
		if len(xTwilioSigHdr) == 1 {
			xTwilioSig = xTwilioSigHdr[0]
		}

		return errors.Wrap(fmt.Errorf("calculated Signature: %s, failed to match X-Twilio-Signature: %s", sig, xTwilioSig), "twiml.Request.ValidatePost()")
	}

	// Validate data
	for _, p := range params {
		var val string
		if len(req.r.PostForm[p]) > 0 {
			val = req.r.PostForm[p][0]
		}
		if valParam, ok := fieldValidators[p]; ok {
			if err := valParam.valFunc(val, valParam.valParam); err != nil {
				return errors.Wrapf(err, "Invalid form value: %s=%s", p, val)
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
	"From": {valFunc: validFromOrTo},
	"To":   {valFunc: validFromOrTo},
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
	"Digits": {valFunc: validateKeyPadEntry},
}
