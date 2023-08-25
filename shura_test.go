package shura

import (
	"os"
	"testing"
)

func Test_tryAccess(t *testing.T) {
	// TODO: assert whether a screenshot is taken
	expected, _ := os.ReadFile("testdata/example_com.txt")

	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"normal", args{"https://example.com/"}, string(expected), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tryAccess(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("tryAccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tryAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
