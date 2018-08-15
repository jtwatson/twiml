package twiml

import "testing"

func TestRequestValues_Duration(t *testing.T) {
	tests := []struct {
		name    string
		r       RequestValues
		want    int
		wantErr bool
	}{
		{name: "Test 1", r: RequestValues{"CallDuration": "12"}, want: 12, wantErr: false},
		{name: "Test 2", r: RequestValues{"CallDuration": "11 sec"}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
