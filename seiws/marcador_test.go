package seiws

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefinirMarcador(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}

		for _, value := range []string{
			"<SiglaSistema>SISTEMA</SiglaSistema>",
			"<IdentificacaoServico>token</IdentificacaoServico>",
			"<IdUnidade>42</IdUnidade>",
			"<Definicoes><item><ProtocoloProcedimento>123</ProtocoloProcedimento><IdMarcador>7</IdMarcador><Texto>Em análise</Texto></item><item><ProtocoloProcedimento>456</ProtocoloProcedimento><IdMarcador>8</IdMarcador><Texto>Urgente</Texto></item></Definicoes>",
		} {
			if !strings.Contains(string(body), value) {
				t.Errorf("request body does not contain %q: %s", value, body)
			}
		}

		w.Header().Set("Content-Type", "text/xml")
		_, _ = io.WriteString(w, `<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="Sei">
  <soap:Body>
    <ns1:definirMarcadorResponse>
      <parametros>true</parametros>
    </ns1:definirMarcadorResponse>
  </soap:Body>
</soap:Envelope>`)
	}))
	defer server.Close()

	client := NewClient(Config{
		URL:                  server.URL,
		SiglaSistema:         "SISTEMA",
		IdentificacaoServico: "token",
	})

	response, err := client.DefinirMarcador(context.Background(), "42", []DefinicaoMarcador{
		{ProtocoloProcedimento: "123", IDMarcador: "7", Texto: "Em análise"},
		{ProtocoloProcedimento: "456", IDMarcador: "8", Texto: "Urgente"},
	})
	if err != nil {
		t.Fatalf("DefinirMarcador() error = %v", err)
	}
	if !response.Parametros {
		t.Errorf("DefinirMarcador() parametros = false, want true")
	}
}

func TestListarMarcadoresUnidade(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}

		for _, value := range []string{
			"<SiglaSistema>SISTEMA</SiglaSistema>",
			"<IdentificacaoServico>token</IdentificacaoServico>",
			"<IdUnidade>42</IdUnidade>",
		} {
			if !strings.Contains(string(body), value) {
				t.Errorf("request body does not contain %q: %s", value, body)
			}
		}

		w.Header().Set("Content-Type", "text/xml")
		_, _ = io.WriteString(w, `<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:ns1="Sei">
  <soap:Body>
    <ns1:listarMarcadoresUnidadeResponse>
      <parametros SOAP-ENC:arrayType="ns1:Marcador[1]" xsi:type="ns1:ArrayOfMarcador">
        <item xsi:type="ns1:Marcador">
          <IdMarcador xsi:type="xsd:string">7</IdMarcador>
          <Nome xsi:type="xsd:string">Em análise</Nome>
          <Icone xsi:type="xsd:string">aWNvbmU=</Icone>
          <SinAtivo xsi:type="xsd:string">S</SinAtivo>
        </item>
      </parametros>
    </ns1:listarMarcadoresUnidadeResponse>
  </soap:Body>
</soap:Envelope>`)
	}))
	defer server.Close()

	client := NewClient(Config{
		URL:                  server.URL,
		SiglaSistema:         "SISTEMA",
		IdentificacaoServico: "token",
	})

	response, err := client.ListarMarcadoresUnidade(context.Background(), "42")
	if err != nil {
		t.Fatalf("ListarMarcadoresUnidade() error = %v", err)
	}
	if len(response.Parametros.Items) != 1 {
		t.Fatalf("got %d markers, want 1", len(response.Parametros.Items))
	}
	got := response.Parametros.Items[0]
	if got.IDMarcador != "7" || got.Nome != "Em análise" || got.Icone != "aWNvbmU=" || got.SinAtivo != "S" {
		t.Errorf("marker = %+v, want ID 7, Em análise, aWNvbmU=, active", got)
	}
}
