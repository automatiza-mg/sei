package seiws

import (
	"context"
	"encoding/xml"
)

// Marcador representa um marcador disponível para uma unidade do SEI.
type Marcador struct {
	IDMarcador string `xml:"IdMarcador" json:"id_marcador"`
	Nome       string `xml:"Nome" json:"nome"`
	Icone      string `xml:"Icone" json:"icone"`
	SinAtivo   string `xml:"SinAtivo" json:"sin_ativo"`
}

// DefinicaoMarcador contém os dados para associar um marcador a um processo.
type DefinicaoMarcador struct {
	ProtocoloProcedimento string `xml:"ProtocoloProcedimento" json:"protocolo_procedimento"`
	IDMarcador            string `xml:"IdMarcador" json:"id_marcador"`
	Texto                 string `xml:"Texto" json:"texto"`
}

// DefinirMarcadorRequest é o payload da operação definirMarcador do SeiWS.php.
type DefinirMarcadorRequest struct {
	XMLName              xml.Name            `xml:"Sei definirMarcador"`
	SiglaSistema         string              `xml:"SiglaSistema"`
	IdentificacaoServico string              `xml:"IdentificacaoServico"`
	IDUnidade            string              `xml:"IdUnidade"`
	Definicoes           []DefinicaoMarcador `xml:"Definicoes>item"`
}

// DefinirMarcadorResponse é o envelope de resposta da operação definirMarcador.
type DefinirMarcadorResponse struct {
	XMLName    xml.Name `xml:"Sei definirMarcadorResponse"`
	Parametros bool     `xml:"parametros" json:"parametros"`
}

// DefinirMarcador associa marcadores aos processos informados pela API SOAP
// legada (SeiWS.php).
func (c *Client) DefinirMarcador(ctx context.Context, idUnidade string, definicoes []DefinicaoMarcador) (*DefinirMarcadorResponse, error) {
	return doReq[DefinirMarcadorRequest, DefinirMarcadorResponse](ctx, c, DefinirMarcadorRequest{
		SiglaSistema:         c.cfg.SiglaSistema,
		IdentificacaoServico: c.cfg.IdentificacaoServico,
		IDUnidade:            idUnidade,
		Definicoes:           definicoes,
	})
}

// ListarMarcadoresUnidadeRequest é o payload da operação listarMarcadoresUnidade
// do SeiWS.php.
type ListarMarcadoresUnidadeRequest struct {
	XMLName              xml.Name `xml:"Sei listarMarcadoresUnidade"`
	SiglaSistema         string   `xml:"SiglaSistema"`
	IdentificacaoServico string   `xml:"IdentificacaoServico"`
	IDUnidade            string   `xml:"IdUnidade"`
}

// ListarMarcadoresUnidadeResponse é o envelope de resposta da operação
// listarMarcadoresUnidade.
type ListarMarcadoresUnidadeResponse struct {
	XMLName    xml.Name             `xml:"Sei listarMarcadoresUnidadeResponse"`
	Parametros Parametros[Marcador] `xml:"parametros" json:"parametros"`
}

// ListarMarcadoresUnidade lista os marcadores disponíveis em uma unidade do
// SEI pela API SOAP legada (SeiWS.php).
func (c *Client) ListarMarcadoresUnidade(ctx context.Context, idUnidade string) (*ListarMarcadoresUnidadeResponse, error) {
	return doReq[ListarMarcadoresUnidadeRequest, ListarMarcadoresUnidadeResponse](ctx, c, ListarMarcadoresUnidadeRequest{
		SiglaSistema:         c.cfg.SiglaSistema,
		IdentificacaoServico: c.cfg.IdentificacaoServico,
		IDUnidade:            idUnidade,
	})
}
