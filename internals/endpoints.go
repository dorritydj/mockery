package internals

import (
	"fmt"
	"strings"
)

type Endpoint struct {
	Id       string `json:"id"`
	Uri      string `json:"uri"`
	Method   string `json:"method"`
	Variants map[string]Variant
}

type Endpoints map[string]Endpoint

type Variant struct {
	Status  int  `json:"status"`
	Default bool `json:"default"`
	Body    any  `json:"body"`
}

func (e Endpoint) GetDefaultResponse() Variant {
	var found Variant
	for _, v := range e.Variants {
		if v.Default {
			found = v
		}
	}

	return found
}

func (e Endpoint) ParseUri() string {
	var parts []string
	res := strings.Fields(strings.ReplaceAll(e.Uri, "/", " "))

	for _, partition := range res {
		colon := string(partition[0])
		rest := string(partition[1:])

		if strings.Compare(colon, ":") == 0 {
			parts = append(parts, fmt.Sprintf("/{%s}", rest))
			continue
		}

		parts = append(parts, fmt.Sprintf("/%s", partition))
	}

	return strings.Join(parts, "")
}
