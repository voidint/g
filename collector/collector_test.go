package collector

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/collector/autoindex"
	"github.com/voidint/g/collector/fancyindex"
	"github.com/voidint/g/collector/official"
	"github.com/voidint/g/pkg/errs"
)

func TestNewCollector(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("hello world")),
	}

	patches := gomonkey.ApplyFuncReturn(http.Get, resp, nil)
	defer patches.Reset()

	type args struct {
		urls []string
	}
	tests := []struct {
		name              string
		args              args
		wantCollectorName string
		wantErr           error
	}{
		{
			name:              "nil parameter",
			args:              args{urls: nil},
			wantCollectorName: official.Name,
		},
		{
			name:              "A slice containing an empty string",
			args:              args{urls: []string{""}},
			wantCollectorName: official.Name,
		},
		{
			name:              "The parameter is a URL slice without a trailing backslash",
			args:              args{urls: []string{"https://mirrors.aliyun.com/golang"}},
			wantCollectorName: fancyindex.Name,
		},
		{
			name:              "A slice containing the name of the official collector",
			args:              args{urls: []string{"official|https://golang.google.cn/dl/"}},
			wantCollectorName: official.Name,
		},
		{
			name:              "A slice containing the name of the fancyindex collector",
			args:              args{urls: []string{"fancyindex|https://mirrors.hust.edu.cn/golang/"}},
			wantCollectorName: fancyindex.Name,
		},
		{
			name:              "A slice containing the name of the autoindex collector",
			args:              args{urls: []string{"autoindex|https://mirrors.ustc.edu.cn/golang/"}},
			wantCollectorName: autoindex.Name,
		},
		{
			name:              "A slice containing only official collector URLs",
			args:              args{urls: []string{OfficialDownloadPageURL}},
			wantCollectorName: official.Name,
		},
		{
			name:              "A slice containing only original official collector URLs",
			args:              args{urls: []string{OriginalOfficialDownloadPageURL}},
			wantCollectorName: official.Name,
		},
		{
			name:              "A slice containing only china official mirror site collector URLs",
			args:              args{urls: []string{CNDownloadPageURL}},
			wantCollectorName: official.Name,
		},
		{
			name:              "A slice containing only Alibaba cloud mirror site collector URLs",
			args:              args{urls: []string{AliYunDownloadPageURL}},
			wantCollectorName: fancyindex.Name,
		},
		{
			name:              "A slice containing only HUST mirror site collector URLs",
			args:              args{urls: []string{HUSTDownloadPageURL}},
			wantCollectorName: fancyindex.Name,
		},
		{
			name:              "A slice containing only NJU mirror site collector URLs",
			args:              args{urls: []string{NJUDownloadPageURL}},
			wantCollectorName: fancyindex.Name,
		},
		{
			name:              "A slice containing only USTC mirror site collector URLs",
			args:              args{urls: []string{USTCDownloadPageURL}},
			wantCollectorName: autoindex.Name,
		},
		{
			name:              "Collector not found",
			args:              args{urls: []string{"hello world"}},
			wantCollectorName: "",
			wantErr:           errs.ErrCollectorNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewCollector(tt.args.urls...)

			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				assert.Equal(t, tt.wantCollectorName, gotC.Name())
			}
		})
	}
}
