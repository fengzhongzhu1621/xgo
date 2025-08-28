package ip

import "testing"

func TestGetDailAddress(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"", args{"http://localhost:80/path?q=a"}, "localhost:80", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDailAddress(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDailAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDailAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
