package version

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {

	vs := []*Version{
		MustNew("1.21.0"),
		MustNew("1.20.10"),
		MustNew("1.21.4"),
		MustNew("1.20rc1"),
		MustNew("1.21rc2"),
		MustNew("1.19.12"),
		MustNew("1.21rc4"),
		MustNew("1.20.1"),
	}

	sort.Sort(Collection(vs))

	assert.Equal(t, vs[0].name, "1.19.12")
	assert.Equal(t, vs[1].name, "1.20rc1")
	assert.Equal(t, vs[2].name, "1.20.1")
	assert.Equal(t, vs[3].name, "1.20.10")
	assert.Equal(t, vs[4].name, "1.21rc2")
	assert.Equal(t, vs[5].name, "1.21rc4")
	assert.Equal(t, vs[6].name, "1.21.0")
	assert.Equal(t, vs[7].name, "1.21.4")
}
