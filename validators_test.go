package twiml

import "testing"

func Test_validSipURI(t *testing.T) {
	type args struct {
		v     interface{}
		param string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Valid SIPS URI", args: args{v: "sips:user@yourdomain.sip.us1.twilio.com:5061", param: ""}, wantErr: false},
		{name: "Valid SIP  URI", args: args{v: "sip:user@yourdomain.sip.us1.twilio.com:5061", param: ""}, wantErr: false},
		{name: "Valid", args: args{v: "", param: "allowempty"}, wantErr: false},
		{name: "Invaid", args: args{v: "", param: ""}, wantErr: true},
		{name: "Invalid SIP URI", args: args{v: "+18002368945", param: ""}, wantErr: true},
		{name: "Invalid URL", args: args{v: "https://user@yourdomain.sip.us1.twilio.com:5061", param: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validSipURI(tt.args.v, tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("validSipURI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
