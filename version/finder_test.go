package version

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/pkg/errs"
)

func genVersions() ([]*Version, error) {
	data, err := os.ReadFile("./testdata/versions.json")
	if err != nil {
		return nil, err
	}
	var items []*struct {
		Name     string     `json:"version"`
		Packages []*Package `json:"packages"`
	}

	if err = json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	vs := make([]*Version, 0, len(items))
	for _, item := range items {
		vs = append(vs, MustNew(item.Name, WithPackages(item.Packages)))
	}
	return vs, nil
}

func TestFinder_Find(t *testing.T) {
	vs, err := genVersions()
	if err != nil {
		assert.Nil(t, err)
	}

	type fields struct {
		goos   string
		goarch string
		items  []*Version
	}
	type args struct {
		vname string
	}

	f := fields{
		goos:   "darwin",
		goarch: "arm64",
		items:  vs,
	}

	fdr := NewFinder(f.items, WithFinderGoos(f.goos), WithFinderGoarch(f.goarch))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Version
		wantErr bool
	}{
		{
			name:    "1、具体版本号",
			fields:  f,
			args:    args{vname: "1.21.4"},
			want:    vs[len(vs)-1],
			wantErr: false,
		},
		{
			name:    "2、最新版本",
			fields:  f,
			args:    args{vname: Latest},
			want:    vs[len(vs)-1],
			wantErr: false,
		},
		{
			name:    "3.1、通配符x",
			fields:  f,
			args:    args{vname: "1.20.x"},
			want:    fdr.MustFind("1.20.11"),
			wantErr: false,
		},
		{
			name:    "3.2、通配符X",
			fields:  f,
			args:    args{vname: "1.20.X"},
			want:    fdr.MustFind("1.20.11"),
			wantErr: false,
		},
		{
			name:    "3.3、通配符*",
			fields:  f,
			args:    args{vname: "1.20.*"},
			want:    fdr.MustFind("1.20.11"),
			wantErr: false,
		},
		{
			name:    "3.4、通配符*",
			fields:  f,
			args:    args{vname: "1.*"},
			want:    fdr.MustFind("1.21.4"),
			wantErr: false,
		},
		{
			name:    "4.1、^匹配最新的次版本号",
			fields:  f,
			args:    args{vname: "^1"},
			want:    fdr.MustFind("1.21.4"),
			wantErr: false,
		},
		{
			name:    "4.2、^匹配最新的次版本号",
			fields:  f,
			args:    args{vname: "^1.18"},
			want:    fdr.MustFind("1.21.4"),
			wantErr: false,
		},
		{
			name:    "4.3、^匹配最新的次版本号",
			fields:  f,
			args:    args{vname: "^1.18.10"},
			want:    fdr.MustFind("1.21.4"),
			wantErr: false,
		},
		{
			name:    "5.1、~匹配某个次版本号的最新修订号",
			fields:  f,
			args:    args{vname: "~1.18"},
			want:    fdr.MustFind("1.18.10"),
			wantErr: false,
		},
		{
			name:    "5.2、~匹配某个次版本号的最新修订版",
			fields:  f,
			args:    args{vname: "~1.18.2"},
			want:    fdr.MustFind("1.18.10"),
			wantErr: false,
		},
		{
			name:    "6、>大于某个版本",
			fields:  f,
			args:    args{vname: ">1.18.2"},
			want:    fdr.MustFind("1.21.4"),
			wantErr: false,
		},
		{
			name:    "7.1、<小于某个版本",
			fields:  f,
			args:    args{vname: "<1.18.10"},
			want:    fdr.MustFind("1.18.9"),
			wantErr: false,
		},
		{
			name:    "7.2、<小于某个版本",
			fields:  f,
			args:    args{vname: "<1.18"},
			want:    fdr.MustFind("1.17.13"),
			wantErr: false,
		},
		{
			name:    "7.3、<小于某个版本但没有合适的软件包",
			fields:  f,
			args:    args{vname: "<1.16"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "8.1、匹配目标版本区间内的最新版本",
			fields:  f,
			args:    args{vname: "1.16 - 1.21"},
			want:    fdr.MustFind("1.21.4"),
			wantErr: false,
		},
		{
			name:    "8.2、匹配目标版本区间内的最新版本",
			fields:  f,
			args:    args{vname: "1.16 - 1.20.7"},
			want:    fdr.MustFind("1.20.7"),
			wantErr: false,
		},
		{
			name:    "9、非法版本号",
			fields:  f,
			args:    args{vname: "voidint"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "10、不存在的版本号",
			fields:  f,
			args:    args{vname: "1.11.111"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fdr.Find(tt.args.vname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Finder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Finder.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFinder_MustFind(t *testing.T) {
	vs, err := genVersions()
	if err != nil {
		assert.Nil(t, err)
	}

	fdr := NewFinder(vs,
		WithFinderPackageKind(ArchiveKind),
		WithFinderGoos("darwin"),
		WithFinderGoarch("arm64"),
	)

	t.Run("查找到版本", func(t *testing.T) {
		v := fdr.MustFind("1.21.4")
		assert.NotNil(t, v)
		assert.Equal(t, v.Name(), "1.21.4")
	})

	t.Run("查找不到版本", func(t *testing.T) {
		assert.Panics(t, func() {
			fdr.MustFind("~1.15")
		})
	})
}

func TestFinder_findLatest(t *testing.T) {
	vs := []*Version{
		MustNew("1.16.1", WithPackages([]*Package{
			{
				FileName:  "go1.16.1.darwin-amd64.tar.gz",
				URL:       "https://golang.google.cn/dl/go1.16.1.darwin-amd64.tar.gz",
				Kind:      ArchiveKind,
				OS:        "macOS",
				Arch:      "x86-64",
				Size:      "124MB",
				Checksum:  "a760929667253cdaa5b10117f536a912be2b0be1006215ff86e957f98f76fd58",
				Algorithm: "SHA256",
			},
			{
				FileName:  "go1.16.1.darwin-arm64.tar.gz",
				URL:       "https://golang.google.cn/dl/go1.16.1.darwin-arm64.tar.gz",
				Kind:      ArchiveKind,
				OS:        "macOS",
				Arch:      "ARM64",
				Size:      "120MB",
				Checksum:  "de2847f49faac2d0608b4afc324cbb3029a496c946db616c294d26082e45f32d",
				Algorithm: "SHA256",
			},
		})),
	}

	tests := []struct {
		name    string
		fdr     *Finder
		wantV   *Version
		wantErr error
	}{
		{
			name:    "查找器中版本列表为空",
			fdr:     NewFinder(nil, WithFinderPackageKind(ArchiveKind), WithFinderGoos("darwin"), WithFinderGoarch("arm64")),
			wantV:   nil,
			wantErr: errs.NewVersionNotFoundError(Latest, "darwin", "arm64"),
		},
		{
			name:    "查找器中版本列表非空且软件包亦匹配",
			fdr:     NewFinder(vs, WithFinderPackageKind(ArchiveKind), WithFinderGoos("darwin"), WithFinderGoarch("arm64")),
			wantV:   vs[0],
			wantErr: nil,
		},
		{
			name:    "查找器中版本列表非空但未找到匹配的软件包",
			fdr:     NewFinder(vs, WithFinderPackageKind(InstallerKind), WithFinderGoos("windows"), WithFinderGoarch("arm64")),
			wantV:   nil,
			wantErr: errs.NewPackageNotFoundError(string(InstallerKind), "windows", "arm64"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := tt.fdr.findLatest()
			if err != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.wantV.Name(), v.Name())
			}
		})
	}
}
