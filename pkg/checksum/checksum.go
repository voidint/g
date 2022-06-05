package checksum

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"

	"github.com/voidint/g/pkg/errs"
)

// Algorithm 校验和算法
type Algorithm string

const (
	// SHA256 校验和算法-sha256
	SHA256 Algorithm = "SHA256"
	// SHA1 校验和算法-sha1
	SHA1 Algorithm = "SHA1"
)

// VerifyFile 检查目标文件校验和
func VerifyFile(algo Algorithm, expectedChecksum, filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var h hash.Hash
	switch algo {
	case SHA256:
		h = sha256.New()
	case SHA1:
		h = sha1.New()
	default:
		return errs.ErrUnsupportedChecksumAlgorithm
	}

	if _, err = io.Copy(h, f); err != nil {
		return err
	}

	if expectedChecksum != hex.EncodeToString(h.Sum(nil)) {
		return errs.ErrChecksumNotMatched
	}
	return nil
}
