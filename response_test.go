package twiml

import (
	"context"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestResponse_Render(t *testing.T) {

	ctx := context.Background()

	header := xml.Header[:len(xml.Header)-1]

	response1 := &Response{
		Verbs: []interface{}{
			&Dial{
				Verbs: []interface{}{
					&Number{
						Value: "810-730-3842",
					},
				},
			},
			&Say{
				Value: "Failed to connect",
			},
		},
	}

	xml1 := header + `
<Response>
  <Dial>
    <Number>810-730-3842</Number>
  </Dial>
  <Say>Failed to connect</Say>
</Response>`

	response2 := NewResponse().
		Dial(NewDial().
			Number(NewNumber("810-730-3842"))).
		Say(NewSay("Failed to connect"))

	xml2 := header + `
<Response>
  <Dial>
    <Number>810-730-3842</Number>
  </Dial>
  <Say>Failed to connect</Say>
</Response>`

	response3 := NewResponse().
		Gather(NewGather().
			SetAction("%s").
			SetMethod(Post).
			SetTimeout(10).
			SetInput("dtmf").
			Say(AliceVoice.Say("Welcome to patch conferencing")).
			Pause(1).
			Say(AliceVoice.Say("Please enter your akkcess code, followed by the pound sign"))).
		Say(AliceVoice.Say("We did not hear a selection")).
		Pause(1).
		Redirect(NewRedirect("%s").
			SetMethod(Post))

	xml3 := header + `
<Response>
  <Gather input="dtmf" action="%s" method="POST" timeout="10">
    <Say voice="alice">Welcome to patch conferencing</Say>
    <Pause length="1"></Pause>
    <Say voice="alice">Please enter your akkcess code, followed by the pound sign</Say>
  </Gather>
  <Say voice="alice">We did not hear a selection</Say>
  <Pause length="1"></Pause>
  <Redirect method="POST">%s</Redirect>
</Response>`

	tests := []struct {
		name     string
		response *Response
		want     string
		wantErr  error
	}{
		{name: "Test1", response: response1, want: xml1},
		{name: "Test2", response: response2, want: xml2},
		{name: "Test3", response: response3, want: xml3},
		{name: "Test4", response: response1, want: xml1},
		{name: "Test5", response: response1, want: xml1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.response.Render(ctx)
			if err != nil {
				t.Errorf("Response.Render() = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Response.Render() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
