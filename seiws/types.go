package seiws

import "encoding/xml"

// Serie identifica o tipo de documento (série) no SEI.
type Serie struct {
	IDSerie string `xml:"IdSerie" json:"id_serie"`
	Nome    string `xml:"Nome" json:"nome"`
}

// Unidade é uma unidade organizacional do SEI.
type Unidade struct {
	IDUnidade       string `xml:"IdUnidade" json:"id_unidade"`
	Sigla           string `xml:"Sigla,omitempty" json:"sigla"`
	Descricao       string `xml:"Descricao,omitempty" json:"descricao"`
	SinProtocolo    string `xml:"SinProtocolo,omitempty" json:"-"`
	SinArquivamento string `xml:"SinArquivamento,omitempty" json:"-"`
	SinOuvidoria    string `xml:"SinOuvidoria,omitempty" json:"-"`
}

// Usuario é um usuário do SEI.
type Usuario struct {
	IDUsuario string `xml:"IdUsuario" json:"id_usuario"`
	Sigla     string `xml:"Sigla" json:"sigla"`
	Nome      string `xml:"Nome" json:"nome"`
}

// Andamento é um andamento (movimentação) registrado em um processo do SEI.
type Andamento struct {
	IDAndamento    string  `xml:"IdAndamento" json:"id_andamento,omitempty"`
	IDTarefa       string  `xml:"IdTarefa" json:"id_tarefa,omitempty"`
	IDTarefaModulo string  `xml:"IdTarefaModulo" json:"id_tarefa_modulo,omitempty"`
	Descricao      string  `xml:"Descricao" json:"descricao"`
	DataHora       string  `xml:"DataHora" json:"data_hora"`
	Unidade        Unidade `xml:"Unidade" json:"unidade"`
}

// UnidadeElaboradora é a unidade do SEI responsável pela elaboração do
// documento.
type UnidadeElaboradora struct {
	IDUnidade string `xml:"IdUnidade" json:"id_unidade"`
	Sigla     string `xml:"Sigla" json:"sigla"`
	Descricao string `xml:"Descricao" json:"descricao"`
}

// Assinaturas é a coleção de [Assinatura] que constam em um documento.
type Assinaturas struct {
	XMLName xml.Name     `xml:"Assinaturas" json:"-"`
	Itens   []Assinatura `xml:"item" json:"itens"`
}

// Assinatura é uma assinatura registrada em um documento do SEI.
type Assinatura struct {
	Nome        string `xml:"Nome" json:"nome"`
	CargoFuncao string `xml:"CargoFuncao" json:"cargo_funcao"`
	DataHora    string `xml:"DataHora" json:"data_hora"`
	IDUsuario   string `xml:"IdUsuario" json:"id_usuario"`
	IDOrigem    string `xml:"IdOrigem" json:"id_origem"`
	IDOrgao     string `xml:"IdOrgao" json:"id_orgao"`
	Sigla       string `xml:"Sigla" json:"sigla"`
}
