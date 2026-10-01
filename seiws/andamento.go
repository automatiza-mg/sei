package seiws

import (
	"context"
	"encoding/xml"
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
