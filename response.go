package twiml

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"go.opencensus.io/trace"
)

// MethodType is an enum for the http method
type MethodType string

const (
	// Post represents an HTTP POST method
	Post MethodType = "POST"

	// Get represents an HTTP GETmethod
	Get MethodType = "GET"
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
func (r *Response) Render(ctx context.Context) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, "twiml.Response.Render()")
	defer span.End()

	buff := new(bytes.Buffer)
	buff.WriteString(xml.Header)
	enc := xml.NewEncoder(buff)
	enc.Indent("", "  ")
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	span.AddAttributes(trace.StringAttribute("twiml", buff.String()))

	return buff.Bytes(), nil
}

// RenderTo writes the Rendered TwiML to the writer
func (r *Response) RenderTo(ctx context.Context, w io.Writer) error {
	ctx, span := trace.StartSpan(ctx, "twiml.Response.RenderTo()")
	defer span.End()

	res, err := r.Render(ctx)
	if err != nil {
		return err
	}
	_, err = w.Write(res)
	return err
}

// Dial represents the TwiML Dial Verb
type Dial struct {
	XMLName xml.Name   `xml:"Dial"`
	Action  string     `xml:"action,attr,omitempty"`
	Method  MethodType `xml:"method,attr,omitempty"`
	Timeout uint       `xml:"timeout,attr,omitempty"`
	Verbs   []interface{}
}

// NewDial returns a Dial verb
func NewDial() *Dial {
	return &Dial{}
}

// Number appends a Number verb to Dial
func (d *Dial) Number(number *Number) *Dial {
	d.Verbs = append(d.Verbs, number)
	return d
}

// Conference appends a Conference verb to Dial
func (d *Dial) Conference(conference *Conference) *Dial {
	d.Verbs = append(d.Verbs, conference)
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

	// PollyMatthew is used in Say to select Amazon Poly voice Matthew
	PollyMatthew VoiceType = "Polly.Matthew"
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
	Method                      MethodType `xml:"method,attr,omitempty"`
	Timeout                     uint       `xml:"timeout,attr,omitempty"`
	FinishOnKey                 *string    `xml:"finishOnKey,attr"`
	NumDigits                   uint       `xml:"numDigits,attr,omitempty"`
	PartialResultCallback       string     `xml:"partialResultCallback,attr,omitempty"`
	PartialResultCallbackMethod MethodType `xml:"partialResultCallbackMethod,attr,omitempty"`
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
func (g *Gather) SetMethod(method MethodType) *Gather {
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
	Method  MethodType `xml:"method,attr,omitempty"`
	Value   string     `xml:",chardata"`
}

// NewRedirect returns a Redirect verb
func NewRedirect(redirect string) *Redirect {
	return &Redirect{Value: redirect}
}

// SetMethod sets the method attribute
func (r *Redirect) SetMethod(method MethodType) *Redirect {
	r.Method = method
	return r
}

// BeepType is an enum type for Beep
type BeepType string

const (
	// BeepsOn beeps on enter and exit
	BeepsOn BeepType = "true"

	// BeepsOff disables beeps
	BeepsOff BeepType = "false"

	// BeepOnEnter beeps on enter only
	BeepOnEnter BeepType = "onEnter"

	// BeepOnExit beeps on exit only
	BeepOnExit BeepType = "onExit"
)

// Conference represents the twiml Conference verb
type Conference struct {
	XMLName                       xml.Name   `xml:"Conference"`
	Muted                         bool       `xml:"muted,attr,omitempty"`
	Beep                          BeepType   `xml:"beep,attr,omitempty"`
	StartConferenceOnEnter        *bool      `xml:"startConferenceOnEnter,attr"`
	EndConferenceOnExit           bool       `xml:"endConferenceOnExit,attr,omitempty"`
	WaitURL                       *string    `xml:"waitUrl,attr"`
	WaitMethod                    MethodType `xml:"waitMethod,attr,omitempty"`
	MaxParticipants               int        `xml:"maxParticipants,attr,omitempty"`
	Record                        string     `xml:"record,attr,omitempty"`
	Region                        string     `xml:"region,attr,omitempty"`
	Trim                          string     `xml:"trim,attr,omitempty"`
	Coach                         string     `xml:"coach,attr,omitempty"`
	StatusCallbackEvent           string     `xml:"statusCallbackEvent,attr,omitempty"`
	StatusCallback                string     `xml:"statusCallback,attr,omitempty"`
	StatusCallbackMethod          MethodType `xml:"statusCallbackMethod,attr,omitempty"`
	RecordingStatusCallback       string     `xml:"recordingStatusCallback,attr,omitempty"`
	RecordingStatusCallbackMethod MethodType `xml:"recordingStatusCallbackMethod,attr,omitempty"`
	RecordingStatusCallbackEvent  string     `xml:"recordingStatusCallbackEvent,attr,omitempty"`
	EventCallbackURL              string     `xml:"eventCallbackUrl,attr,omitempty"`
	Value                         string     `xml:",chardata"`
}

// NewConference returns a Conference verb
func NewConference(conferenceName string) *Conference {
	return &Conference{Value: conferenceName}
}

// SetMuted sets the muted attribute
func (c *Conference) SetMuted(muted bool) *Conference {
	c.Muted = muted
	return c
}

// SetBeep sets the beep attribute
func (c *Conference) SetBeep(beep BeepType) *Conference {
	c.Beep = beep
	return c
}

// SetStartConferenceOnEnter sets the startConferenceOnEnter attribute
func (c *Conference) SetStartConferenceOnEnter(startConferenceOnEnter bool) *Conference {
	c.StartConferenceOnEnter = &startConferenceOnEnter
	return c
}

// SetEndConferenceOnExit sets the endConferenceOnExit attribute
func (c *Conference) SetEndConferenceOnExit(endConferenceOnExit bool) *Conference {
	c.EndConferenceOnExit = endConferenceOnExit
	return c
}

// SetWaitURL sets the waitURL attribute
func (c *Conference) SetWaitURL(waitURL string) *Conference {
	c.WaitURL = &waitURL
	return c
}

// SetWaitMethod sets the waitMethod attribute
func (c *Conference) SetWaitMethod(waitMethod MethodType) *Conference {
	c.WaitMethod = waitMethod
	return c
}

// SetMaxParticipants sets the maxParticipants attribute
func (c *Conference) SetMaxParticipants(maxParticipants int) *Conference {
	c.MaxParticipants = maxParticipants
	return c
}

// SetRecord sets the record attribute
func (c *Conference) SetRecord(record string) *Conference {
	c.Record = record
	return c
}

// SetRegion sets the region attribute
func (c *Conference) SetRegion(region string) *Conference {
	c.Region = region
	return c
}

// SetTrim sets the trim attribute
func (c *Conference) SetTrim(trim string) *Conference {
	c.Trim = trim
	return c
}

// SetCoach sets the coach attribute
func (c *Conference) SetCoach(coach string) *Conference {
	c.Coach = coach
	return c
}

// SetStatusCallbackEvent sets the statusCallbackEvent attribute
func (c *Conference) SetStatusCallbackEvent(statusCallbackEvent conferenceCallbackEvents) *Conference {
	c.StatusCallbackEvent = string(statusCallbackEvent)
	return c
}

// SetStatusCallback sets the statusCallback attribute
func (c *Conference) SetStatusCallback(statusCallback string) *Conference {
	c.StatusCallback = statusCallback
	return c
}

// SetStatusCallbackMethod sets the statusCallbackMethod attribute
func (c *Conference) SetStatusCallbackMethod(statusCallbackMethod MethodType) *Conference {
	c.StatusCallbackMethod = statusCallbackMethod
	return c
}

// SetRecordingStatusCallback sets the recordingStatusCallback attribute
func (c *Conference) SetRecordingStatusCallback(recordingStatusCallback string) *Conference {
	c.RecordingStatusCallback = recordingStatusCallback
	return c
}

// SetRecordingStatusCallbackMethod sets the recordingStatusCallbackMethod attribute
func (c *Conference) SetRecordingStatusCallbackMethod(recordingStatusCallbackMethod MethodType) *Conference {
	c.RecordingStatusCallbackMethod = recordingStatusCallbackMethod
	return c
}

// SetRecordingStatusCallbackEvent sets the recordingStatusCallbackEvent attribute
func (c *Conference) SetRecordingStatusCallbackEvent(recordingStatusCallbackEvent string) *Conference {
	c.RecordingStatusCallbackEvent = recordingStatusCallbackEvent
	return c
}

// SetEventCallbackURL sets the eventCallbackURL attribute
func (c *Conference) SetEventCallbackURL(eventCallbackURL string) *Conference {
	c.EventCallbackURL = eventCallbackURL
	return c
}

// start end join leave mute hold speaker
type conferenceCallbackEvents string

// ConferenceCallbackEvents enables specific Callback Events
func ConferenceCallbackEvents() conferenceCallbackEvents {
	return conferenceCallbackEvents("")
}

// Start enables the Callback Event to indicate Conference has Started
func (c conferenceCallbackEvents) Start() conferenceCallbackEvents {
	return conferenceCallbackEvents(strings.TrimLeft(fmt.Sprintf("%s start", c), " "))
}

// End enables the Callback Event to indicate Conference has Ended
func (c conferenceCallbackEvents) End() conferenceCallbackEvents {
	return conferenceCallbackEvents(strings.TrimLeft(fmt.Sprintf("%s end", c), " "))
}

// Join enables the Callback Event to indicate Participant has joined
func (c conferenceCallbackEvents) Join() conferenceCallbackEvents {
	return conferenceCallbackEvents(strings.TrimLeft(fmt.Sprintf("%s join", c), " "))
}

// Leave enables the Callback Event to indicate Participant has left
func (c conferenceCallbackEvents) Leave() conferenceCallbackEvents {
	return conferenceCallbackEvents(strings.TrimLeft(fmt.Sprintf("%s leave", c), " "))
}

// Mute enables the Callback Event to indicate Participant has been muted/unmuted
func (c conferenceCallbackEvents) Mute() conferenceCallbackEvents {
	return conferenceCallbackEvents(strings.TrimLeft(fmt.Sprintf("%s mute", c), " "))
}

// Hold enables the Callback Event to indicate Participant has been held
func (c conferenceCallbackEvents) Hold() conferenceCallbackEvents {
	return conferenceCallbackEvents(strings.TrimLeft(fmt.Sprintf("%s hold", c), " "))
}

// Speaker enables the Callback Event to indicate Participant has started/stoped speaking
func (c conferenceCallbackEvents) Speaker() conferenceCallbackEvents {
	return conferenceCallbackEvents(strings.TrimLeft(fmt.Sprintf("%s speaker", c), " "))
}
