package seiws

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListarAndamentos(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}

		for _, value := range []string{
			"<SiglaSistema>SISTEMA</SiglaSistema>",
			"<IdentificacaoServico>token</IdentificacaoServico>",
			"<IdUnidade>42</IdUnidade>",
			"<ProtocoloProcedimento>123</ProtocoloProcedimento>",
			"<SinRetornarAtributos>N</SinRetornarAtributos>",
			"<Tarefas><item>1</item><item>2</item></Tarefas>",
		} {
			if !strings.Contains(string(body), value) {
				t.Errorf("request body does not contain %q: %s", value, body)
			}
		}

		w.Header().Set("Content-Type", "text/xml")
		_, _ = io.WriteString(w, `<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <listarAndamentosResponse xmlns="Sei">
      <parametros><item><Descricao>Processo recebido</Descricao><DataHora>2026-10-01 12:00:00</DataHora></item></parametros>
    </listarAndamentosResponse>
  </soap:Body>
</soap:Envelope>`)
	}))
	defer server.Close()

	client := NewClient(Config{
		URL:                  server.URL,
		SiglaSistema:         "SISTEMA",
		IdentificacaoServico: "token",
	})

	response, err := client.ListarAndamentos(context.Background(), "42", "123", []string{"1", "2"})
	if err != nil {
		t.Fatalf("ListarAndamentos() error = %v", err)
	}
	if len(response.Parametros.Items) != 1 {
		t.Fatalf("got %d items, want 1", len(response.Parametros.Items))
	}
	if got := response.Parametros.Items[0].Descricao; got != "Processo recebido" {
		t.Errorf("andamento description = %q, want %q", got, "Processo recebido")
	}
}
