package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_printEnv(t *testing.T) {
	type args struct {
		envNames []string
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			args:  args{envNames: []string{"G_AUTHOR"}},
			wantW: "G_AUTHOR=\"voidint\"\n",
		},
	}

	os.Setenv("G_AUTHOR", "voidint")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			printEnv(w, tt.args.envNames)

			assert.Equal(t, tt.wantW, w.String())
		})
	}
}
