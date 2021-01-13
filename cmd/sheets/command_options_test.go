package sheets

import "testing"

func Test_spreadIDFromURL(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "reg",
			args: args{s: "https://docs.google.com/spreadsheets/d/abcd1234EFGH5678xx-xyz012B/edit#gid=12345678"},
			want: "abcd1234EFGH5678xx-xyz012B",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := spreadIDFromURL(tt.args.s); got != tt.want {
				t.Errorf("spreadIDFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
