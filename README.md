## modcons 
[![Go Report Card](https://goreportcard.com/badge/github.com/go-tooling/modcons)](https://goreportcard.com/report/github.com/go-tooling)

Modcons is a CLI tool that inspects go.mod files for deprecated versions according to a set of deprecation rules.

For example:

```$xslt
github.com/myles-mcdonnell/blondie v1.0.0>=v3.0.0         #whitelist range equal or more than v1 less than v3
github.com/myles-mcdonnell/blondie =v0.8.0,v0.9.3         #whitelist v0.8.0 and v0.9.3
github.com/myles-mcdonnell/blondie !v1.5.7>=v1.8.3        #blacklist range equal or more than v1.5.7 less than 1.8.3
github.com/myles-mcdonnell/blondie !=v2.5.0               #blacklist

```

### Install
The latest binaries for all supported operating systems are [here](https://github.com/go-tooling/modcons/releases)

If you have Go tool installed you may also run:
```
go get -u github.com/go-tooling/modcons/...
```


### Usage

Example CLI usage:

```
modcons --rulepath=http://my.domain.com/myrules.modcons --modpath=./go.mod --parseOnly=false
```

- Note that both path args may be local or http(s) urls.  
- `parseOnly` will only parse the rule file and will not inspect the go.mod file
- Both `modpath` and `parseOnly` flag are optional, the default values are shown above
- If any deprecated references are identified modcons exits with code 1.  

