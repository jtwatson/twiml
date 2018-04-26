package twiml

import (
	"bytes"
	"encoding/xml"
	"io"
)

type methodType string

const (
	// Post represents an HTTP POST method
	Post methodType = "POST"

	// Get represents an HTTP GETmethod
	Get methodType = "GET"
)

// Response represents the TwiML Response Verb
type Response struct {
	Verbs []interface{}
}

// NewResponse returns a Response
func NewResponse() *Response {
	return &Response{}
}

// Gather adds the Gather verb to the response
func (r *Response) Gather(gather *Gather) *Response {
	r.Verbs = append(r.Verbs, gather)
	return r
}

// Dial adds the dial verb to the response
func (r *Response) Dial(dial *Dial) *Response {
	r.Verbs = append(r.Verbs, dial)
	return r
}

// Say adds the say verb to the Response
func (r *Response) Say(say *Say) *Response {
	r.Verbs = append(r.Verbs, say)
	return r
}

// Pause appends a Pause verb to Dial
func (r *Response) Pause(length uint) *Response {
	r.Verbs = append(r.Verbs, NewPause(length))
	return r
}

// Redirect appends a Redirect verb to Response
func (r *Response) Redirect(redirect *Redirect) *Response {
	r.Verbs = append(r.Verbs, redirect)
	return r
}

// Render returns a reader which returns the rendered twiml response
func (r *Response) Render() ([]byte, error) {
	buff := new(bytes.Buffer)
	buff.WriteString(xml.Header)
	enc := xml.NewEncoder(buff)
	enc.Indent("", "  ")
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// WriteTo writes the Rendered TwiML to the writer
func (r *Response) WriteTo(w io.Writer) (int64, error) {
	res, err := r.Render()
	if err != nil {
		return 0, err
	}
	n, err := w.Write(res)
	return int64(n), err
}

// Dial represents the TwiML Dial Verb
type Dial struct {
	XMLName xml.Name   `xml:"Dial"`
	Action  string     `xml:"action,attr,omitempty"`
	Method  methodType `xml:"method,attr,omitempty"`
	Timeout uint       `xml:"timeout,attr,omitempty"`
	Verbs   []interface{}
}

// NewDial returns a Dial verb
func NewDial() *Dial {
	return &Dial{}
}

// Say appends a Say verb to Dial
func (d *Dial) Say(say *Say) *Dial {
	d.Verbs = append(d.Verbs, say)
	return d
}

// Number appends a Number verb to Dial
func (d *Dial) Number(number *Number) *Dial {
	d.Verbs = append(d.Verbs, number)
	return d
}

// VoiceType is enum type for voice
type VoiceType string

// Say returns a Say verb with voice set to current value
func (v VoiceType) Say(msg string) *Say {
	return &Say{
		Voice: v,
		Value: msg,
	}
}

const (
	// ManVoice is used in Say to select man voice
	ManVoice VoiceType = "man"

	// WomenVoice is used in Say to select woman voice
	WomenVoice VoiceType = "women"

	// AliceVoice is used in Say to select alice voice
	AliceVoice VoiceType = "alice"
)

// Say represents the TwiML Say verb
type Say struct {
	XMLName xml.Name  `xml:"Say"`
	Voice   VoiceType `xml:"voice,attr,omitempty"`
	Loop    uint      `xml:"loop,attr,omitempty"`
	Value   string    `xml:",chardata"`
}

// NewSay returns a Say verb
func NewSay(msg string) *Say {
	return &Say{Value: msg}
}

// SetVoice sets the voice value
func (s *Say) SetVoice(voice VoiceType) *Say {
	s.Voice = voice
	return s
}

// Number represents a phone number to call
type Number struct {
	XMLName xml.Name `xml:"Number"`
	Value   string   `xml:",chardata"`
}

// NewNumber returns a Number verb
func NewNumber(number string) *Number {
	return &Number{Value: number}
}

// Gather represents the TwiML Gather verb
type Gather struct {
	XMLName                     xml.Name   `xml:"Gather"`
	Input                       string     `xml:"input,attr,omitempty"`
	Action                      string     `xml:"action,attr,omitempty"`
	Method                      methodType `xml:"method,attr,omitempty"`
	Timeout                     uint       `xml:"timeout,attr,omitempty"`
	FinishOnKey                 *string    `xml:"finishOnKey,attr"`
	NumDigits                   uint       `xml:"numDigits,attr,omitempty"`
	PartialResultCallback       string     `xml:"partialResultCallback,attr,omitempty"`
	PartialResultCallbackMethod methodType `xml:"partialResultCallbackMethod,attr,omitempty"`
	Language                    string     `xml:"language,attr,omitempty"`
	Hints                       string     `xml:"hints,attr,omitempty"`
	ProfanityFilter             bool       `xml:"profanityFilter,attr,omitempty"`
	SpeechTimeout               uint       `xml:"speechTimeout,attr,omitempty"`
	Verbs                       []interface{}
}

// NewGather returns a Gather verb
func NewGather() *Gather {
	return &Gather{}
}

// Say appends a Say verb to Gather
func (g *Gather) Say(say *Say) *Gather {
	g.Verbs = append(g.Verbs, say)
	return g
}

// Pause appends a Pause verb to Gather
func (g *Gather) Pause(length uint) *Gather {
	g.Verbs = append(g.Verbs, NewPause(length))
	return g
}

// SetInput sets the input attribute
func (g *Gather) SetInput(input string) *Gather {
	g.Input = input
	return g
}

// SetAction sets the action attribute
func (g *Gather) SetAction(action string) *Gather {
	g.Action = action
	return g
}

// SetMethod sets the method attribute
func (g *Gather) SetMethod(method methodType) *Gather {
	g.Method = method
	return g
}

// SetTimeout sets the timeout attribute
func (g *Gather) SetTimeout(timeout uint) *Gather {
	g.Timeout = timeout
	return g
}

// Pause represents the TwiML Pause verb
type Pause struct {
	XMLName xml.Name `xml:"Pause"`
	Length  uint     `xml:"length,attr,omitempty"`
}

// NewPause returns a Pause verb
func NewPause(length uint) *Pause {
	return &Pause{Length: length}
}

// Redirect represents the TwiML Redirect verb
type Redirect struct {
	XMLName xml.Name   `xml:"Redirect"`
	Method  methodType `xml:"method,attr,omitempty"`
	Value   string     `xml:",chardata"`
}

// NewRedirect returns a Redirect verb
func NewRedirect(redirect string) *Redirect {
	return &Redirect{Value: redirect}
}

// SetMethod sets the method attribute
func (r *Redirect) SetMethod(method methodType) *Redirect {
	r.Method = method
	return r
}
