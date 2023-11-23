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

	"github.com/Masterminds/semver/v3"
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
		{
			name:    "无效版本号",
			args:    args{vname: "abcdef"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Semantify(tt.args.vname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Semantify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if !reflect.DeepEqual(got.String(), tt.wantString) {
					t.Errorf("Semantify().String() = %v, want %v", got.String(), tt.wantString)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("合法的版本号", func(t *testing.T) {
		vname := "1.21.4"
		opts := []func(v *Version){nil, nil, nil}
		v, err := New(vname, opts...)
		assert.Nil(t, err)
		assert.True(t, reflect.DeepEqual(v, &Version{name: "1.21.4", sv: semver.MustParse("1.21.4")}))
	})

	t.Run("非法的版本号", func(t *testing.T) {
		vname := "1.2.3.4"
		v, err := New(vname)
		assert.NotNil(t, err)
		assert.True(t, errs.IsMalformedVersion(err))
		assert.Nil(t, v)
	})
}

func TestMustNew(t *testing.T) {
	t.Run("合法的版本号", func(t *testing.T) {
		vname := "1.21.4"
		v := MustNew(vname)
		assert.True(t, reflect.DeepEqual(v, &Version{name: "1.21.4", sv: semver.MustParse("1.21.4")}))
	})

	t.Run("非法的版本号", func(t *testing.T) {
		vname := "1.2.3.4"
		assert.Panics(t, func() {
			MustNew(vname)
		})
	})
}

func TestVersion_FindPackages(t *testing.T) {
	vs, err := genVersions()
	if err != nil {
		assert.Nil(t, err)
	}

	t.Run("查找到唯一的软件包", func(t *testing.T) {
		v1214 := vs[len(vs)-1] // 1.21.4

		pkgs, err := v1214.FindPackages(ArchiveKind, "darwin", "arm64")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(pkgs))
		assert.Equal(t, pkgs[0].FileName, "go1.21.4.darwin-arm64.tar.gz")
		assert.Equal(t, pkgs[0].URL, "https://golang.google.cn/dl/go1.21.4.darwin-arm64.tar.gz")
		assert.Equal(t, pkgs[0].Kind, ArchiveKind)
		assert.Equal(t, pkgs[0].OS, "macOS")
		assert.Equal(t, pkgs[0].Arch, "ARM64")
		assert.Equal(t, pkgs[0].Size, "62MB")
		assert.Equal(t, pkgs[0].Checksum, "8b7caf2ac60bdff457dba7d4ff2a01def889592b834453431ae3caecf884f6a5")
		assert.Equal(t, pkgs[0].Algorithm, "SHA256")
	})

	t.Run("查找到多个软件包", func(t *testing.T) {
		v1214 := vs[len(vs)-1] // 1.21.4

		pkgs, err := v1214.FindPackages(ArchiveKind, "linux", "arm")
		assert.Nil(t, err)
		assert.Equal(t, 2, len(pkgs))
	})

	t.Run("未查找到软件包", func(t *testing.T) {
		v1214 := vs[len(vs)-1] // 1.21.4

		pkgs, err := v1214.FindPackages(ArchiveKind, "darwin", "ppc64")
		assert.True(t, errs.IsPackageNotFound(err))
		assert.Equal(t, 0, len(pkgs))
	})
}

func TestVersion_Packages(t *testing.T) {
	filename := "go1.21.4.darwin-arm64.tar.gz"
	url := "https://golang.google.cn/dl/go1.21.4.darwin-arm64.tar.gz"
	kind := ArchiveKind
	os := "macOS"
	arch := "ARM64"
	size := "62MB"
	checksum := "8b7caf2ac60bdff457dba7d4ff2a01def889592b834453431ae3caecf884f6a5"
	algorithm := "SHA256"

	tests := []struct {
		name string
		v    *Version
		want []Package
	}{
		{
			name: "软件包列表为空",
			v:    MustNew("1"),
			want: make([]Package, 0),
		},
		{
			name: "软件包列表非空",
			v: MustNew("1.21.4", WithPackages([]*Package{
				{
					FileName:  filename,
					URL:       url,
					Kind:      kind,
					OS:        os,
					Arch:      arch,
					Size:      size,
					Checksum:  checksum,
					Algorithm: algorithm,
				},
			})),
			want: []Package{
				{
					FileName:  filename,
					URL:       url,
					Kind:      kind,
					OS:        os,
					Arch:      arch,
					Size:      size,
					Checksum:  checksum,
					Algorithm: algorithm,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Packages(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Version.Packages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_MatchConstraint(t *testing.T) {
	cs, err := semver.NewConstraint("~1.18")
	assert.Nil(t, err)
	assert.NotNil(t, cs)

	tests := []struct {
		name string
		v    *Version
		cs   *semver.Constraints
		want bool
	}{
		{
			name: "不匹配约束",
			v:    MustNew("1.21.4"),
			cs:   cs,
			want: false,
		},
		{
			name: "匹配约束",
			v:    MustNew("1.18.4"),
			cs:   cs,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.MatchConstraint(tt.cs); got != tt.want {
				t.Errorf("Version.MatchConstraint() = %v, want %v", got, tt.want)
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
			_, _ = f.Seek(0, 0)
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
			_, _ = f.Seek(0, 0)
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
