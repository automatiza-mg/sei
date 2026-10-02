package seiws

import (
	"context"
	"encoding/xml"
	"strconv"
)

// ListarAndamentosRequest é o payload da operação listarAndamentos do SeiWS.php.
type ListarAndamentosRequest struct {
	XMLName               xml.Name `xml:"Sei listarAndamentos"`
	SiglaSistema          string
	IdentificacaoServico  string
	IDUnidade             string `xml:"IdUnidade"`
	ProtocoloProcedimento string
	SinRetornarAtributos  string
	Andamentos            []string `xml:"Andamentos>item,omitempty"`
	Tarefas               []string `xml:"Tarefas>item,omitempty"`
	TarefasModulos        []string `xml:"TarefasModulos>item,omitempty"`
}

// ListarAndamentosResponse é o envelope de resposta da operação
// listarAndamentos.
type ListarAndamentosResponse struct {
	XMLName    xml.Name              `xml:"Sei listarAndamentosResponse"`
	Parametros Parametros[Andamento] `xml:"parametros" json:"parametros"`
}

// ListarAndamentos lista os andamentos de um processo no SEI pela API SOAP
// legada (SeiWS.php), filtrando pelas tarefas informadas.
func (c *Client) ListarAndamentos(ctx context.Context, idUnidade, protocolo string, tarefas []string) (*ListarAndamentosResponse, error) {
	return doReq[ListarAndamentosRequest, ListarAndamentosResponse](ctx, c, ListarAndamentosRequest{
		SiglaSistema:          c.cfg.SiglaSistema,
		IdentificacaoServico:  c.cfg.IdentificacaoServico,
		IDUnidade:             idUnidade,
		ProtocoloProcedimento: protocolo,
		SinRetornarAtributos:  "S",
		Tarefas:               tarefas,
	})
}

// ListarAndamentosMarcadoresRequest é o payload da operação
// listarAndamentosMarcadores do SeiWS.php.
type ListarAndamentosMarcadoresRequest struct {
	XMLName               xml.Name `xml:"Sei listarAndamentosMarcadores"`
	SiglaSistema          string
	IdentificacaoServico  string
	IDUnidade             string `xml:"IdUnidade"`
	ProtocoloProcedimento string
	Marcadores            MarcadoresFiltro `xml:"Marcadores,omitempty"`
}

// MarcadorFiltro representa um marcador usado para filtrar andamentos. Um
// IdMarcador nil corresponde a um marcador removido.
type MarcadorFiltro struct {
	IdMarcador *string `xml:"IdMarcador"`
}

// MarcadoresFiltro representa o tipo ArrayOfIdMarcadores definido no WSDL.
type MarcadoresFiltro []MarcadorFiltro

// MarshalXML serializa Marcadores como um array SOAP de xsd:string, de acordo
// com o tipo ArrayOfIdMarcadores declarado no WSDL.
func (m MarcadoresFiltro) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr,
		xml.Attr{Name: xml.Name{Local: "xmlns:SOAP-ENC"}, Value: "http://schemas.xmlsoap.org/soap/encoding/"},
		xml.Attr{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
		xml.Attr{Name: xml.Name{Local: "xmlns:xsd"}, Value: "http://www.w3.org/2001/XMLSchema"},
		xml.Attr{Name: xml.Name{Local: "xmlns:tns"}, Value: "Sei"},
		xml.Attr{Name: xml.Name{Local: "xsi:type"}, Value: "tns:ArrayOfIdMarcadores"},
		xml.Attr{Name: xml.Name{Local: "SOAP-ENC:arrayType"}, Value: "xsd:string[" + strconv.Itoa(len(m)) + "]"},
	)
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for _, marcador := range m {
		item := xml.StartElement{Name: xml.Name{Local: "item"}}
		item.Attr = append(item.Attr, xml.Attr{Name: xml.Name{Local: "xsi:type"}, Value: "xsd:string"})
		if marcador.IdMarcador == nil {
			item.Attr = append(item.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
			if err := e.EncodeToken(item); err != nil {
				return err
			}
			if err := e.EncodeToken(item.End()); err != nil {
				return err
			}
			continue
		}
		if err := e.EncodeElement(*marcador.IdMarcador, item); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

// ListarAndamentosMarcadoresResponse é o envelope de resposta da operação
// listarAndamentosMarcadores.
type ListarAndamentosMarcadoresResponse struct {
	XMLName    xml.Name                      `xml:"Sei listarAndamentosMarcadoresResponse"`
	Parametros Parametros[AndamentoMarcador] `xml:"parametros" json:"parametros"`
}

// ListarAndamentosMarcadores lista o histórico de marcadores de um processo
// no SEI. Quando marcadores é informado, os resultados são filtrados por seus
// identificadores; um elemento vazio também solicita andamentos de remoção.
func (c *Client) ListarAndamentosMarcadores(ctx context.Context, idUnidade, protocolo string, marcadores []MarcadorFiltro) (*ListarAndamentosMarcadoresResponse, error) {
	// O WSDL permite itens nulos no array, mas o SEI em homologação interpreta
	// xsi:nil como array vazio para IdMarcador e responde com Fault. Quando o
	// filtro inclui eventos de remoção, busca sem Marcadores e filtra o retorno
	// localmente para contornar esse comportamento.
	filtrarLocalmente := false
	idsMarcadores := make(map[string]struct{}, len(marcadores))
	for _, marcador := range marcadores {
		if marcador.IdMarcador == nil {
			filtrarLocalmente = true
			continue
		}
		idsMarcadores[*marcador.IdMarcador] = struct{}{}
	}

	requestMarcadores := MarcadoresFiltro(marcadores)
	if filtrarLocalmente {
		requestMarcadores = nil
	}

	response, err := doReq[ListarAndamentosMarcadoresRequest, ListarAndamentosMarcadoresResponse](ctx, c, ListarAndamentosMarcadoresRequest{
		SiglaSistema:          c.cfg.SiglaSistema,
		IdentificacaoServico:  c.cfg.IdentificacaoServico,
		IDUnidade:             idUnidade,
		ProtocoloProcedimento: protocolo,
		Marcadores:            requestMarcadores,
	})
	if err != nil {
		return nil, err
	}
	if !filtrarLocalmente {
		return response, nil
	}

	filtrados := make([]AndamentoMarcador, 0, len(response.Parametros.Items))
	for _, andamento := range response.Parametros.Items {
		if andamento.Marcador == nil {
			filtrados = append(filtrados, andamento)
			continue
		}
		if _, ok := idsMarcadores[andamento.Marcador.IDMarcador]; ok {
			filtrados = append(filtrados, andamento)
		}
	}
	response.Parametros.Items = filtrados
	return response, nil
}
