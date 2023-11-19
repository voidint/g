package version

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/pkg/errs"
)

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
