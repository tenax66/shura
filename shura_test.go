package shura

import (
	"os"
	"reflect"
	"regexp"
	"testing"
)

func Test_extractLinks(t *testing.T) {

	html, _ := os.ReadFile("testdata/example_com.txt")

	type args struct {
		html  string
		regex *regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"normal", args{string(html), regexp.MustCompile(`http.*\.org`)}, []string{"https://www.iana.org"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractLinks(tt.args.html, tt.args.regex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
