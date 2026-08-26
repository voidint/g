// Copyright (c) 2019 voidint <voidint@126.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cli

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/version"
)

func TestSuperseded(t *testing.T) {
	tests := []struct {
		name    string
		install []string
		inuse   string
		want    []string
	}{
		{
			name:    "no version installed",
			install: nil,
			inuse:   "",
			want:    nil,
		},
		{
			name:    "single version",
			install: []string{"1.24.6"},
			inuse:   "1.24.6",
			want:    nil,
		},
		{
			name:    "one version per minor series",
			install: []string{"1.22.12", "1.23.8", "1.24.6"},
			inuse:   "1.24.6",
			want:    nil,
		},
		{
			name:    "keeps the newest of each minor series",
			install: []string{"1.19.10", "1.19.13", "1.20.4", "1.20.5", "1.24.5", "1.24.6"},
			inuse:   "1.24.6",
			want:    []string{"1.19.10", "1.20.4", "1.24.5"},
		},
		{
			name:    "keeps the version in use even when superseded",
			install: []string{"1.24.3", "1.24.5", "1.24.6"},
			inuse:   "1.24.3",
			want:    []string{"1.24.5"},
		},
		{
			name:    "prerelease is superseded by the final release",
			install: []string{"1.25rc1", "1.25.0"},
			inuse:   "1.25.0",
			want:    []string{"1.25rc1"},
		},
		{
			name:    "keeps the newest prerelease within an unreleased series",
			install: []string{"1.26beta1", "1.26rc1"},
			inuse:   "1.26rc1",
			want:    []string{"1.26beta1"},
		},
		{
			name:    "different major versions are independent series",
			install: []string{"1.24.1", "2.0.0", "2.0.1"},
			inuse:   "1.24.1",
			want:    []string{"2.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := make([]*version.Version, 0, len(tt.install))
			for _, vname := range tt.install {
				items = append(items, version.MustNew(vname))
			}
			sort.Sort(version.Collection(items))

			doomed := superseded(items, tt.inuse)
			got := make([]string, 0, len(doomed))
			for _, v := range doomed {
				got = append(got, v.Name())
			}
			if tt.want == nil {
				assert.Empty(t, got)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
