package rules

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	versions = []version{
		{0, 0, 0},
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 1},
		{1, 0, 0},
		{1, 0, 1},
		{1, 1, 0},
		{1, 1, 1},
	}

	ranges = make(map[int]map[int]Range)
)

func init() {
	for i := range versions {
		for x := range versions {
			if x > i {
				if _, ok := ranges[i]; !ok {
					ranges[i] = make(map[int]Range)
				}

				upp := versions[x]
				lwr := versions[i]
				rnge, err := newVersionRange(
					&rangeLimit{upp, true},
					&rangeLimit{lwr, true})

				if err != nil {
					panic(err)
				}

				ranges[i][x] = rnge
			}
		}
	}
}

func TestRules_IsDeprecated_Whitelist(t *testing.T) {
	mod := module{Path: "abc", Version: versions[0]}

	rls := rules{modules: map[string]*rule{mod.Path: {Whitelist: []version{versions[0]}}}}
	assert.False(t, rls.IsDeprecated(mod))

	rls = rules{modules: map[string]*rule{mod.Path: {Whitelist: []version{versions[1]}}}}
	assert.True(t, rls.IsDeprecated(mod))
}

func TestRules_IsDeprecated_Blacklist(t *testing.T) {
	mod := module{Path: "abc", Version: versions[0]}

	rls := rules{modules: map[string]*rule{mod.Path: {Blacklist: []version{versions[0]}}}}
	assert.True(t, rls.IsDeprecated(mod))

	rls = rules{modules: map[string]*rule{mod.Path: {Blacklist: []version{versions[1]}}}}
	assert.False(t, rls.IsDeprecated(mod))
}

func TestRules_IsDeprecated_WhiteAndBlacklist(t *testing.T) {
	mod := module{Path: "abc", Version: versions[0]}

	rules := rules{modules: map[string]*rule{mod.Path: {
		Whitelist: []version{versions[0]},
		Blacklist: []version{versions[0]}}}}

	assert.True(t, rules.IsDeprecated(mod))
}

func Test_IsIncluded_VersionList(t *testing.T) {
	for i, v := range versions {
		assert.True(t, versions[i].isIncluded([]Range{}, versions))
		for y := range versions {
			assert.Equal(t, versions[i].isIncluded([]Range{}, []version{versions[y]}), i == y)
			assert.Equal(t, v.isIncluded([]Range{}, remove(versions, y)), i != y)
		}
	}
}

func Test_IsIncluded_Ranges(t *testing.T) {
	for v := range versions {
		for x := range ranges {
			for y := range ranges[x] {
				assert.Equal(t, versions[v].isIncluded([]Range{ranges[x][y]}, []version{}), v >= x && v <= y)
			}
		}
	}
}

func remove(s []version, i int) []version {
	n := make([]version, len(s))
	copy(n, s)
	n[len(s)-1], n[i] = s[i], s[len(s)-1]
	return n[:len(s)-1]
}

func Test_Version_IsLessThan(t *testing.T) {
	for i, version := range versions {
		for y := range versions {
			act, _ := version.IsLessThan(versions[y])
			assert.Equal(t, i < y, act)

			if y > 0 {
				act, _ := versions[y].IsLessThan(version)
				assert.Equal(t, y < i, act)
			}
		}
		break
	}
}

func Test_Version_IsMoreThan(t *testing.T) {
	for i, version := range versions {
		for y := range versions {
			act, _ := version.IsMoreThan(versions[y])
			assert.Equal(t, i > y, act)

			if y > 0 {
				act, _ := versions[y].IsMoreThan(version)
				assert.Equal(t, y > i, act)
			}
		}
		break
	}
}

func Test_Range_In(t *testing.T) {
	for i := range versions {
		for y := range versions {
			if y > i {
				rnge, err := newVersionRange(
					&rangeLimit{versions[y], true},
					&rangeLimit{versions[i], true})

				assert.Nil(t, err)

				for x := range versions {
					act := rnge.In(versions[x])
					ex := x >= i && x <= y
					act = rnge.In(versions[x])
					assert.Equal(t, ex, act)
				}
			}
		}
		break
	}
}

func Test_Parse_Rules(t *testing.T) {
	file, err := os.Open("./test_data/rules.modcop")
	assert.Nil(t, err)

	r, err := ParseRules(file)
	assert.Nil(t, err)

	path := "github.com/myles-mcdonnell/blondie"
	assert.True(t, r.IsDeprecated(module{Path: path, Version: version{Major: 2, Minor: 5, Build: 0}}))
	assert.True(t, r.IsDeprecated(module{Path: path, Version: version{Major: 0, Minor: 7, Build: 0}}))
	assert.True(t, r.IsDeprecated(module{Path: path, Version: version{Major: 1, Minor: 5, Build: 7}}))
	assert.True(t, r.IsDeprecated(module{Path: path, Version: version{Major: 1, Minor: 8, Build: 2}}))
	assert.False(t, r.IsDeprecated(module{Path: path, Version: version{Major: 1, Minor: 8, Build: 3}}))
	assert.False(t, r.IsDeprecated(module{Path: path, Version: version{Major: 2, Minor: 5, Build: 1}}))
	assert.False(t, r.IsDeprecated(module{Path: path, Version: version{Major: 0, Minor: 8, Build: 0}}))
	assert.False(t, r.IsDeprecated(module{Path: path, Version: version{Major: 0, Minor: 9, Build: 3}}))
}

func Test_Parse_Mod(t *testing.T) {
	file, err := os.Open("./go.mod")
	assert.Nil(t, err)

	mods, err := ParseModFile(file)

	assert.Equal(t, 1, len(mods))
}
