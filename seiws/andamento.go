package seiws

import "encoding/xml"

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
