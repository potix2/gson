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
			name: "last whitespace",
			args: args{text: ` true `},
			want: true,
		},
		{
			name: "false",
			args: args{text: `false`},
			want: false,
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
	tests := []struct {
		name    string
		input   string
		want    interface{}
		want1   int
		wantErr bool
	}{
		{
			name:  "single character",
			input: `"a"`,
			want:  "a",
			want1: 3,
		},
		{
			name:  "trim string",
			input: `"a" `,
			want:  "a",
			want1: 3,
		},
		{
			name:  "empty string",
			input: `""`,
			want:  "",
			want1: 2,
		},
		{
			name:  "multiple character",
			input: `"ab1234"`,
			want:  "ab1234",
			want1: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseString([]byte(tt.input), 0)
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

func Test_parseArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		want1   int
		wantErr bool
	}{
		{
			name:  "[]",
			input: `[]`,
			want:  []interface{}{},
			want1: 2,
		},
		{
			name:  "[ ]",
			input: `[ ]`,
			want:  []interface{}{},
			want1: 3,
		},
		{
			name:  `["a"]`,
			input: `["a"]`,
			want:  []interface{}{"a"},
			want1: 5,
		},
		{
			name:  `[true]`,
			input: `[ true ]`,
			want:  []interface{}{true},
			want1: 8,
		},
		{
			name:  `[true,false]`,
			input: `[ true, false ]`,
			want:  []interface{}{true, false},
			want1: 15,
		},
		{
			name:  `["a", "b"]`,
			input: `[ "a" , "b" ]`,
			want:  []interface{}{"a", "b"},
			want1: 13,
		},
		{
			name:    `[ `,
			input:   `[ `,
			wantErr: true,
			want1:   2,
		},
		{
			name:    `["a"`,
			input:   `["a"`,
			wantErr: true,
			want1:   4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseArray([]byte(tt.input), 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArray() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseArray() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_parseValue(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		want1   int
		wantErr bool
	}{
		{
			name:  "string",
			input: `"a"`,
			want:  "a",
			want1: 3,
		},
		{
			name:  "trim",
			input: ` "a" `,
			want:  "a",
			want1: 5,
		},
		{
			name:  "true",
			input: `true`,
			want:  true,
			want1: 4,
		},
		{
			name:  "null",
			input: `null`,
			want:  nil,
			want1: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseValue([]byte(tt.input), 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
