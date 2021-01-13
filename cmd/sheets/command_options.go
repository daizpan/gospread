package sheets

import (
	"net/url"
	"regexp"

	"github.com/daizpan/gospread"
)

type CommandOptions struct {
	Credential string
	GoSpread   *gospread.GoSpread
}

func spreadIDFromArg(arg string) string {
	spreadID := spreadIDFromURL(arg)
	if spreadID == "" {
		spreadID = arg
	}
	return spreadID
}

var spreadURLRE = regexp.MustCompile(`^/[^/]+/d/([^/]+)`)

func spreadIDFromURL(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return ""
	}

	m := spreadURLRE.FindStringSubmatch(u.Path)
	if m == nil {
		return ""
	}
	return m[1]
}
