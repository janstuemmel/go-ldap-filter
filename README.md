# LDAP filter

A work-in-progress [RFC4515](https://datatracker.ietf.org/doc/html/rfc4515) ldap filter parser

Only a few filters are implemented at the moment (And, Or, Equality)

## Usage

```go
import "github.com/janstuemmel/go-ldap-filter"

filter, err := ldapfilter.NewParser("|(name=Jon)(name=Foo)").Parse()

if err != nil {
  panic(err)
}

ok := filter.Match(map[string][]string{
  "name": {"Jon"}
})

fmt.Println(ok)
```