package twiml

import (
	"reflect"
	"testing"
	"time"
)

func TestRequestValues_Duration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       RequestValues
		want    time.Duration
		wantErr bool
	}{
		{name: "Test 1", r: RequestValues{"CallDuration": "12"}, want: 12 * time.Second, wantErr: false},
		{name: "Test 2", r: RequestValues{"CallDuration": "11 sec"}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.r.CallDuration()
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestValues.CallDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RequestValues.CallDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNumber(t *testing.T) {
	t.Parallel()

	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want *ParsedNumber
	}{
		{name: "Valid Number", args: args{v: "+18005642365"}, want: &ParsedNumber{Valid: true, Number: "+18005642365", Raw: "+18005642365"}},
		{name: "Valid SIP us1", args: args{v: "sips:8005642365@domain.sip.us1.twilio.com:5061"}, want: &ParsedNumber{Valid: true, Number: "+18005642365", SIP: true, SIPDomain: "domain", Region: "sip", Raw: "sips:8005642365@domain.sip.us1.twilio.com:5061"}},
		{name: "Valid SIP us2", args: args{v: "sips:8005642365@domain.sip.us2.twilio.com:5061"}, want: &ParsedNumber{Valid: true, Number: "+18005642365", SIP: true, SIPDomain: "domain", Region: "sip", Raw: "sips:8005642365@domain.sip.us2.twilio.com:5061"}},
		{name: "Invalid SIP sip.domain.com", args: args{v: "sips:8005642365@domain.sip.us1.domain.com:5061"}, want: &ParsedNumber{Number: "sips:8005642365@domain.sip.us1.domain.com:5061", Raw: "sips:8005642365@domain.sip.us1.domain.com:5061"}},
		{name: "Invalid SIP sip2.twilio.com", args: args{v: "sips:8005642365@domain.sip2.us1.twilio.com:5061"}, want: &ParsedNumber{Number: "sips:8005642365@domain.sip2.us1.twilio.com:5061", Raw: "sips:8005642365@domain.sip2.us1.twilio.com:5061"}},
		{name: "Invalid SIP twilio.com", args: args{v: "sips:8005642365@sip.us1.twilio.com:5061"}, want: &ParsedNumber{Number: "sips:8005642365@sip.us1.twilio.com:5061", Raw: "sips:8005642365@sip.us1.twilio.com:5061"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ParseNumber(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
