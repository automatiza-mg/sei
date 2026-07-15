package seiws

import (
	"context"
	"encoding/xml"
)

// ListarUnidadesRequest é o payload da operação listarUnidades do SeiWS.php.
type ListarUnidadesRequest struct {
	XMLName              xml.Name `xml:"Sei listarUnidades"`
	SiglaSistema         string   `xml:"SiglaSistema"`
	IdentificacaoServico string   `xml:"IdentificacaoServico"`
}

// ListarUnidadesResponse é o envelope de resposta da operação listarUnidades.
type ListarUnidadesResponse struct {
	XMLName    xml.Name            `xml:"Sei listarUnidadesResponse"`
	Parametros Parametros[Unidade] `xml:"parametros" json:"parametros"`
}

// ListarUnidades lista as unidades disponíveis no SEI pela API SOAP legada
// (SeiWS.php).
func (c *Client) ListarUnidades(ctx context.Context) (*ListarUnidadesResponse, error) {
	return doReq[ListarUnidadesRequest, ListarUnidadesResponse](ctx, c, ListarUnidadesRequest{
		SiglaSistema:         c.cfg.SiglaSistema,
		IdentificacaoServico: c.cfg.IdentificacaoServico,
	})
}
