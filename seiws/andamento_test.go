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
			"<SinRetornarAtributos>S</SinRetornarAtributos>",
			"<Tarefas><item>1</item><item>2</item></Tarefas>",
		} {
			if !strings.Contains(string(body), value) {
				t.Errorf("request body does not contain %q: %s", value, body)
			}
		}

		w.Header().Set("Content-Type", "text/xml")
		_, _ = io.WriteString(w, `<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:ns1="Sei">
  <soap:Body>
    <ns1:listarAndamentosResponse>
      <parametros SOAP-ENC:arrayType="ns1:Andamento[1]" xsi:type="ns1:ArrayOfAndamento">
        <item xsi:type="ns1:Andamento">
          <IdAndamento xsi:type="xsd:string">243251</IdAndamento>
          <Descricao xsi:type="xsd:string">Processo recebido</Descricao>
          <DataHora xsi:type="xsd:string">01/10/2026 12:23:14</DataHora>
          <Atributos SOAP-ENC:arrayType="ns1:AtributoAndamento[2]" xsi:type="ns1:ArrayOfAtributoAndamento">
            <item xsi:type="ns1:AtributoAndamento"><Nome xsi:type="xsd:string">DESCRICAO</Nome><Valor xsi:type="xsd:string">Texto do andamento</Valor><IdOrigem xsi:nil="true"/></item>
            <item xsi:type="ns1:AtributoAndamento"><Nome xsi:type="xsd:string">UNIDADE</Nome><Valor xsi:type="xsd:string">ABC/DCD¥Automatiza</Valor><IdOrigem xsi:type="xsd:string">110001112</IdOrigem></item>
          </Atributos>
        </item>
      </parametros>
    </ns1:listarAndamentosResponse>
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
	gotAtributos := response.Parametros.Items[0].Atributos
	if len(gotAtributos) != 2 {
		t.Fatalf("got %d attributes, want 2", len(gotAtributos))
	}
	if got := gotAtributos[0]; got.Nome != "DESCRICAO" || got.Valor != "Texto do andamento" || got.IDOrigem != "" {
		t.Errorf("first attribute = %+v, want DESCRICAO / Texto do andamento / empty origin", got)
	}
	if got := gotAtributos[1]; got.Nome != "UNIDADE" || got.Valor != "ABC/DCD¥Automatiza" || got.IDOrigem != "110001112" {
		t.Errorf("second attribute = %+v, want UNIDADE / ABC/DCD¥Automatiza / 110001112", got)
	}
}
