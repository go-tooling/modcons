package rules

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type (
	badRangeError struct{}

	version struct {
		Major, Minor, Build int
	}

	module struct {
		Path    string
		Version version
	}

	rangeLimit struct {
		version   version
		Inclusive bool
	}

	versionRange struct {
		upper *rangeLimit
		lower *rangeLimit
	}

	rule struct {
		WhitelistRanges []Range
		Whitelist       []version
		BlacklistRanges []Range
		Blacklist       []version
	}

	rules struct {
		modules map[string]*rule
	}

	Range interface {
		In(version version) bool
	}

	Rules interface {
		IsDeprecated(module module) bool
	}
)

func (*badRangeError) Error() string {
	return "Bad range"
}

func IsBadRange(err error) bool {
	_, ok := err.(*badRangeError)

	return ok
}

func newVersionRange(upper, lower *rangeLimit) (Range, error) {
	if upper == nil && lower == nil {
		return nil, &badRangeError{}
	}

	if !(upper == nil || lower == nil) {
		lessThan, equal := upper.version.IsLessThan(lower.version)
		if lessThan || equal {
			return nil, &badRangeError{}
		}
	}

	return &versionRange{upper: upper, lower: lower}, nil
}

func (rules *rules) IsDeprecated(module module) bool {
	if rule, ok := rules.modules[module.Path]; ok {
		if len(rule.Whitelist) > 0 || len(rule.WhitelistRanges) > 0 {
			if !module.Version.isIncluded(rule.WhitelistRanges, rule.Whitelist) {
				return true
			}
		}

		return module.Version.isIncluded(rule.BlacklistRanges, rule.Blacklist)
	}

	return false
}

func (version version) isIncluded(ranges []Range, versions []version) bool {
	isIncluded := false
	for _, whiteRange := range ranges {
		if whiteRange.In(version) {
			isIncluded = true
			break
		}
	}

	if !isIncluded {
		for _, whiteVersion := range versions {
			if whiteVersion == version {
				isIncluded = true
				break
			}
		}
	}

	return isIncluded
}

func (rnge versionRange) In(version version) bool {
	if rnge.lower != nil {
		lessThan, equal := rnge.lower.version.IsLessThan(version)

		if !lessThan && !(rnge.lower.Inclusive && equal) {
			return false
		}
	}

	if rnge.upper != nil {
		moreThan, equal := rnge.upper.version.IsMoreThan(version)

		if !moreThan {
			return rnge.upper.Inclusive && equal
		}
	}

	return true
}

func (versionX version) IsMoreThan(versionY version) (bool, bool) {
	return versionY.IsLessThan(versionX)
}

func (x version) IsLessThan(y version) (bool, bool) {
	for _, c := range []struct{ x, y int }{
		{x.Major, y.Major},
		{x.Minor, y.Minor},
		{x.Build, y.Build}} {
		if c.x > c.y {
			return false, false
		} else if c.x < c.y {
			return true, false
		}
	}

	return false, true
}

func ParseRules(rdr io.Reader) (Rules, error) {
	scanner := bufio.NewScanner(rdr)
	rules := &rules{modules: make(map[string]*rule)}
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) < 2 {
			continue
		}

		path := parts[0]
		var rl *rule
		rl, ok := rules.modules[path]
		if !ok {
			rl = &rule{}
			rules.modules[path] = rl
		}

		ruleExp := parts[1]
		black := false
		if strings.Index(ruleExp, "!") == 0 {
			black = true
			ruleExp = ruleExp[1:]
		}

		isRange, inclusive := false, false
		var operatorIndex int
		if oi := strings.Index(ruleExp, ">="); oi > -1 {
			isRange = true
			inclusive = true
			operatorIndex = oi
		} else if oi := strings.Index(ruleExp, ">"); oi > -1 {
			isRange = true
			operatorIndex = oi
		}

		if isRange {
			fromVersionToken := ruleExp[:operatorIndex]
			fromVersion, err := parseSemver(fromVersionToken)
			if err != nil {
				return nil, err
			}
			ruleExp = ruleExp[operatorIndex+2:]
			toVersion, err := parseSemver(ruleExp)
			if err != nil {
				return nil, err
			}

			rl.addVersionRange(black, inclusive, fromVersion, toVersion)
			continue
		}

		if strings.Index(ruleExp, "=") == 0 {
			ruleExp = ruleExp[1:]
			for _, token := range strings.Split(ruleExp, ",") {
				version, err := parseSemver(token)
				if err != nil {
					return nil, err
				}
				rl.addVersion(black, version)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func (r *rule) addVersion(black bool, v version) {
	if black {
		if r.Blacklist == nil {
			r.Blacklist = make([]version, 0)
		}
		r.Blacklist = append(r.Blacklist, v)
	} else {
		if r.Whitelist == nil {
			r.Whitelist = make([]version, 0)
		}
		r.Whitelist = append(r.Whitelist, v)
	}
}

func (r *rule) addVersionRange(black, inclusive bool, lower, upper version) {
	rnge := versionRange{
		lower: &rangeLimit{version: lower, Inclusive: inclusive},
		upper: &rangeLimit{version: upper}}

	if black {
		if r.BlacklistRanges == nil {
			r.BlacklistRanges = make([]Range, 0)
		}
		r.BlacklistRanges = append(r.BlacklistRanges, rnge)
	} else {
		if r.WhitelistRanges == nil {
			r.WhitelistRanges = make([]Range, 0)
		}
		r.WhitelistRanges = append(r.WhitelistRanges, rnge)
	}
}

func parseSemver(token string) (version, error) {
	version := version{}
	if strings.Index(token, "v") != 0 {
		return version, fmt.Errorf("Unable to parse semver token: %v ", token)
	}

	token = token[1:]
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return version, fmt.Errorf("Unable to parse semver token: %v ", token)
	}

	var err error
	version.Major, err = strconv.Atoi(parts[0])
	if err != nil {
		return version, fmt.Errorf("Unable to parse semver token: %v ", token)
	}

	version.Minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return version, fmt.Errorf("Unable to parse semver token: %v ", token)
	}

	version.Build, err = strconv.Atoi(parts[2])
	if err != nil {
		return version, fmt.Errorf("Unable to parse semver token: %v ", token)
	}

	return version, nil
}

func ParseModFile(r io.Reader) ([]module, error) {
	scanner := bufio.NewScanner(r)
	mods := []module{}
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")

		if len(parts) < 3 || parts[0] != "require" {
			continue
		}

		version, err := parseSemver(parts[2])
		if err != nil {
			continue
		}

		mods = append(mods, module{Version: version, Path: parts[1]})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mods, nil
}

func (v version) String() string {
	return fmt.Sprintf("v%v.%v.%v", v.Major, v.Minor, v.Build)
}
