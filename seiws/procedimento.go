package seiws

import (
	"context"
	"encoding/xml"
)

// UnidadeProcedimentoAberto representa uma unidade em que um processo está
// aberto, junto do usuário responsável (atribuição).
type UnidadeProcedimentoAberto struct {
	Unidade Unidade `xml:"Unidade" json:"unidade"`
	Usuario Usuario `xml:"Usuario" json:"usuario"`
}

// ConsultarProcedimentoRequest é o payload da operação consultarProcedimento
// do SeiWS.php.
type ConsultarProcedimentoRequest struct {
	XMLName                               xml.Name `xml:"Sei consultarProcedimento"`
	SiglaSistema                          string
	IdentificacaoServico                  string
	IDUnidade                             string `xml:"IdUnidade,omitempty"`
	ProtocoloProcedimento                 string
	SinRetornarAssuntos                   string
	SinRetornarInteressados               string
	SinRetornarObservacoes                string
	SinRetornarAndamentoGeracao           string
	SinRetornarAndamentoConclusao         string
	SinRetornarUltimoAndamento            string
	SinRetornarUnidadesProcedimentoAberto string
	SinRetornarProcedimentosRelacionados  string
}

// RetornoConsultaProcedimento contém os metadados de um processo retornados
// pela API SOAP do SEI.
type RetornoConsultaProcedimento struct {
	IDProcedimento             string                                `xml:"IdProcedimento" json:"id_procedimento"`
	ProcedimentoFormatado      string                                `xml:"ProcedimentoFormatado" json:"procedimento_formatado"`
	Especificacao              string                                `xml:"Especificacao" json:"especificacao"`
	DataAutuacao               string                                `xml:"DataAutuacao" json:"data_autuacao"`
	NivelAcessoLocal           int                                   `xml:"NivelAcessoLocal" json:"nivel_acesso_local"`
	NivelAcessoGlobal          int                                   `xml:"NivelAcessoGlobal" json:"nivel_acesso_global"`
	LinkAcesso                 string                                `xml:"LinkAcesso" json:"link_acesso"`
	AndamentoGeracao           Andamento                             `xml:"AndamentoGeracao" json:"andamento_geracao"`
	UnidadesProcedimentoAberto Parametros[UnidadeProcedimentoAberto] `xml:"UnidadesProcedimentoAberto" json:"unidades_procedimento_aberto"`
}

// ConsultarProcedimentoResponse é o envelope de resposta da operação
// consultarProcedimento.
type ConsultarProcedimentoResponse struct {
	XMLName    xml.Name                    `xml:"Sei consultarProcedimentoResponse"`
	Parametros RetornoConsultaProcedimento `xml:"parametros" json:"parametros"`
}

// ConsultarProcedimento consulta os metadados de um processo no SEI pela API
// SOAP legada (SeiWS.php).
func (c *Client) ConsultarProcedimento(ctx context.Context, protocolo string) (*ConsultarProcedimentoResponse, error) {
	return doReq[ConsultarProcedimentoRequest, ConsultarProcedimentoResponse](ctx, c, ConsultarProcedimentoRequest{
		SiglaSistema:                          c.cfg.SiglaSistema,
		IdentificacaoServico:                  c.cfg.IdentificacaoServico,
		ProtocoloProcedimento:                 protocolo,
		SinRetornarAndamentoGeracao:           "S",
		SinRetornarUnidadesProcedimentoAberto: "S",
	})
}

// UnidadeDestino agrupa os ids das unidades de destino de um envio de
// processo.
type UnidadeDestino struct {
	IDUnidade []string `xml:"IdUnidade"`
}

// EnviarProcessoRequest é o payload da operação enviarProcesso do SeiWS.php.
type EnviarProcessoRequest struct {
	XMLName                       xml.Name `xml:"Sei enviarProcesso"`
	SiglaSistema                  string
	IdentificacaoServico          string
	IDUnidade                     string `xml:"IdUnidade"`
	ProtocoloProcedimento         string
	UnidadesDestino               UnidadeDestino
	SinManterAbertoUnidade        string
	SinRemoverAnotacao            string
	SinEnviarEmailNotificacao     string
	SinDiasUteisRetornoProgramado string
	SinReabrir                    string
}

// EnviarProcessoResponse é o envelope de resposta da operação enviarProcesso.
type EnviarProcessoResponse struct {
	XMLName    xml.Name `xml:"Sei enviarProcessoResponse"`
	Parametros string   `xml:"parametros" json:"parametros"`
}

// EnviarProcesso movimenta um processo da unidade de origem para uma ou mais
// unidades de destino no SEI pela API SOAP legada (SeiWS.php).
func (c *Client) EnviarProcesso(ctx context.Context, protocolo string, unidadeOrigem string, unidadesDestino []string) (*EnviarProcessoResponse, error) {
	return doReq[EnviarProcessoRequest, EnviarProcessoResponse](ctx, c, EnviarProcessoRequest{
		SiglaSistema:          c.cfg.SiglaSistema,
		IdentificacaoServico:  c.cfg.IdentificacaoServico,
		IDUnidade:             unidadeOrigem,
		ProtocoloProcedimento: protocolo,
		UnidadesDestino: UnidadeDestino{
			IDUnidade: unidadesDestino,
		},
		SinManterAbertoUnidade:        "N",
		SinRemoverAnotacao:            "N",
		SinEnviarEmailNotificacao:     "N",
		SinDiasUteisRetornoProgramado: "N",
		SinReabrir:                    "N",
	})
}
