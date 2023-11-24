package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsset_IsCompressedFile(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{
			name:        "application/zip",
			contentType: "application/zip",
			want:        true,
		},
		{
			name:        "application/x-gzip",
			contentType: "application/x-gzip",
			want:        true,
		},
		{
			name:        "application/json",
			contentType: "application/json",
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Asset{
				ContentType: tt.contentType,
			}
			assert.Equal(t, tt.want, a.IsCompressedFile())
		})
	}
}
