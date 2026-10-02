package soap

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestEncodedPlacesEncodingStyleOnOperation(t *testing.T) {
	type operation struct {
		XMLName xml.Name `xml:"Sei listarAndamentosMarcadores"`
		Foo     string
	}
	got, err := xml.Marshal(Envelope[Encoded[operation]]{
		Body: Body[Encoded[operation]]{Content: Encoded[operation]{Value: operation{Foo: "bar"}}},
	})
	if err != nil {
		t.Fatalf("xml.Marshal() error = %v", err)
	}

	xmlText := string(got)
	if !strings.Contains(xmlText, `<listarAndamentosMarcadores xmlns="Sei" xmlns:envelope="http://schemas.xmlsoap.org/soap/envelope/" envelope:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">`) {
		t.Errorf("encodingStyle is not on operation element: %s", xmlText)
	}
	if strings.Contains(xmlText, `<Body xmlns:envelope="http://schemas.xmlsoap.org/soap/envelope/" envelope:encodingStyle=`) {
		t.Errorf("encodingStyle must not be on Body: %s", xmlText)
	}
}
