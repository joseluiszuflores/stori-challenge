package internal

import "testing"

func TestToIntFromFile(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success to get name",
			args: args{val: "1.csv"},
			want: "1",
		},
		{
			name: "Success to get name without extension",
			args: args{val: "1"},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToIntFromFile(tt.args.val); got != tt.want {
				t.Errorf("ToIntFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
