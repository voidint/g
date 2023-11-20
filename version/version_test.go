package version

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/pkg/errs"
)

func TestSemantify(t *testing.T) {
	type args struct {
		vname string
	}
	tests := []struct {
		name       string
		args       args
		wantString string
		wantErr    bool
	}{
		{
			name:       "普通三位版本号",
			args:       args{vname: "1.21.4"},
			wantString: "1.21.4",
			wantErr:    false,
		},
		{
			name:       "普通两位版本号",
			args:       args{vname: "1.18"},
			wantString: "1.18.0",
			wantErr:    false,
		},
		{
			name:       "1.21.0版本号",
			args:       args{vname: "1.21.0"},
			wantString: "1.21.0",
			wantErr:    false,
		},
		{
			name:       "alpha",
			args:       args{vname: "1.19alpha1"},
			wantString: "1.19.0-alpha1",
			wantErr:    false,
		},
		{
			name:       "beta",
			args:       args{vname: "1.19beta1"},
			wantString: "1.19.0-beta1",
			wantErr:    false,
		},
		{
			name:       "rc",
			args:       args{vname: "1.21rc4"},
			wantString: "1.21.0-rc4",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Semantify(tt.args.vname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Semantify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.wantString) {
				t.Errorf("Semantify().String() = %v, want %v", got.String(), tt.wantString)
			}
		})
	}
}

func TestVerifyChecksum(t *testing.T) {
	t.Run("检查安装包校验和", func(t *testing.T) {
		filename := fmt.Sprintf("%d.txt", time.Now().Unix())
		f, err := os.Create(filename)
		assert.Nil(t, err)
		defer os.Remove(filename)
		defer f.Close()
		_, err = f.WriteString("hello 世界！")
		assert.Nil(t, err)

		t.Run("SHA256", func(t *testing.T) {
			_, _ = f.Seek(0, 0)
			h := sha256.New()
			_, err = io.Copy(h, f)
			assert.Nil(t, err)

			pkg := &Package{
				Algorithm: "SHA256",
				Checksum:  fmt.Sprintf("%x", h.Sum(nil)),
			}
			assert.Nil(t, pkg.VerifyChecksum(filename))
		})

		t.Run("校验和不匹配", func(t *testing.T) {
			f.Seek(0, 0)
			h := sha1.New()
			_, err = io.Copy(h, f)
			assert.Nil(t, err)

			pkg := &Package{
				Algorithm: "SHA1",
				Checksum:  fmt.Sprintf("%x", h.Sum(nil)),
			}
			assert.Nil(t, pkg.VerifyChecksum(filename))
		})

		t.Run("SHA1", func(t *testing.T) {
			f.Seek(0, 0)
			h := sha1.New()
			_, err = io.Copy(h, f)
			assert.Nil(t, err)

			pkg := &Package{
				Algorithm: "SHA1",
				Checksum:  "hello",
			}
			assert.Equal(t, errs.ErrChecksumNotMatched, pkg.VerifyChecksum(filename))
		})

		t.Run("SHA1024", func(t *testing.T) {
			pkg := &Package{
				Algorithm: "SHA1024",
			}
			assert.Equal(t, errs.ErrUnsupportedChecksumAlgorithm, pkg.VerifyChecksum(filename))
		})
	})
}
