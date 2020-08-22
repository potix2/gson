package gson

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "single character",
			args: args{text: `"a"`},
			want: "a",
		},
		{
			name: "empty string",
			args: args{text: `""`},
			want: "",
		},
		{
			name: "multiple character",
			args: args{text: `"ab1234"`},
			want: "ab1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
