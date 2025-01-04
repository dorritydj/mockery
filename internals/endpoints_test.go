package internals

import "testing"

func TestParseUri(t *testing.T) {
	endpoint := Endpoint{
		Uri: "/path",
	}
	want := "/path"

	if endpoint.ParseUri() != want {
		t.Errorf("uri: %s parsed %s", endpoint.Uri, endpoint.ParseUri())
	}

	endpoint.Uri = "/path/to/route"
	want = "/path/to/route"

	if endpoint.ParseUri() != want {
		t.Errorf("uri: %s parsed %s", endpoint.Uri, endpoint.ParseUri())
	}

	endpoint.Uri = "/path/with/:slug"
	want = "/path/with/{slug}"

	if endpoint.ParseUri() != want {
		t.Errorf("uri: %s parsed %s", endpoint.Uri, endpoint.ParseUri())
	}

	endpoint.Uri = "/:slug"
	want = "/{slug}"

	if endpoint.ParseUri() != want {
		t.Errorf("uri: %s parsed %s", endpoint.Uri, endpoint.ParseUri())
	}

}

func TestGetDefaultResponse(t *testing.T) {
	endpoint := Endpoint{
		Variants: map[string]Variant{
			"success": {},
			"error": {
				Default: true,
			},
		},
	}
	want := Variant{
		Default: true,
	}

	if def, err := endpoint.GetDefaultResponse(); err != nil && def != want {
		t.Errorf("endpoint %#v want %#v", endpoint, want)
	}

	endpoint = Endpoint{
		Variants: map[string]Variant{
			"success": {},
			"error":   {},
		},
	}

	if _, err := endpoint.GetDefaultResponse(); err != nil {
		t.Errorf("endpoint %#v want %#v", endpoint, want)
	}
}
