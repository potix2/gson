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
			name: "string with whitespace",
			args: args{text: ` "a"`},
			want: "a",
		},
		{
			name: "trim string",
			args: args{text: `"a" `},
			want: "a",
		},
		{
			name: "null",
			args: args{text: `null`},
			want: nil,
		},
		{
			name:    "invalid token",
			args:    args{text: `nil`},
			wantErr: true,
		},
		{
			name: "true",
			args: args{text: `true`},
			want: true,
		},
		{
			name: "false",
			args: args{text: `false`},
			want: false,
		},
		{
			name: "empty array",
			args: args{text: `[]`},
			want: []interface{}{},
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

func Test_parseString(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		want1   int
		wantErr bool
	}{
		{
			name:  "single character",
			args:  args{text: `"a"`},
			want:  "a",
			want1: 3,
		},
		{
			name:  "trim string",
			args:  args{text: `"a" `},
			want:  "a",
			want1: 3,
		},
		{
			name:  "empty string",
			args:  args{text: `""`},
			want:  "",
			want1: 2,
		},
		{
			name:  "multiple character",
			args:  args{text: `"ab1234"`},
			want:  "ab1234",
			want1: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseString([]byte(tt.args.text), 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
