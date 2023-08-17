package shura

import "testing"

func Test_fetchContent(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{"https://example.com/"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fetchContent(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
