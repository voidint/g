package version

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
