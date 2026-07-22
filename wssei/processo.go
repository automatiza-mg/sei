package wssei

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Processo representa um processo retornado pelo WSSEI.
type Processo struct {
	ID                 string            `json:"id"`
	Status             string            `json:"status"`
	SeiNumMaxDocsPasta string            `json:"seiNumMaxDocsPasta"`
	Atributos          ProcessoAtributos `json:"atributos"`
}

// ProcessoAtributos reúne os atributos detalhados de um [Processo].
type ProcessoAtributos struct {
	IDProcedimento   string                           `json:"idProcedimento"`
	IDProtocolo      string                           `json:"idProtocolo"`
	Numero           string                           `json:"numero"`
	TipoProcesso     string                           `json:"tipoProcesso"`
	Descricao        string                           `json:"descricao"`
	UsuarioAtribuido Object[ProcessoUsuarioAtribuido] `json:"usuarioAtribuido"`
	Unidade          Object[ProcessoUnidade]          `json:"unidade"`
	Ciencias         Slice[ProcessoCiencia]           `json:"ciencias"`
	Marcadores       Slice[ProcessoMarcador]          `json:"marcador"`
	DadosAbertura    Object[ProcessoDadosAbertura]    `json:"dadosAbertura"`
	Anotacoes        Slice[ProcessoAnotacao]          `json:"anotacoes"`
	Status           ProcessoStatusFlags              `json:"status"`
}

// ProcessoUsuarioAtribuido representa o usuário atribuído a um [Processo].
type ProcessoUsuarioAtribuido struct {
	IDAtividade   string `json:"idAtividade"`
	IDUsuario     string `json:"idUsuario"`
	Sigla         string `json:"sigla"`
	Nome          string `json:"nome"`
	Nomeformatado string `json:"nomeformatado"`
}

// ProcessoUnidade representa a unidade atual de um [Processo].
type ProcessoUnidade struct {
	IDUnidade string `json:"idUnidade"`
	Sigla     string `json:"sigla"`
}

// ProcessoCiencia representa um registro de ciência em um [Processo].
type ProcessoCiencia struct {
	IDProtocolo  string `json:"idProtocolo"`
	IDAtividade  string `json:"idAtividade"`
	Data         string `json:"data"`
	IDUnidade    string `json:"idUnidade"`
	Unidade      string `json:"unidade"`
	SiglaUnidade string `json:"siglaUnidade"`
	IDUsuario    string `json:"idUsuario"`
	SiglaUsuario string `json:"siglaUsuario"`
	NomeUsuario  string `json:"nomeUsuario"`
	Descricao    string `json:"descricao"`
}

// ProcessoMarcador representa o marcador aplicado a um [Processo].
type ProcessoMarcador struct {
	IDMarcador   string `json:"idMarcador"`
	Nome         string `json:"nome"`
	Texto        string `json:"texto"`
	IDCor        string `json:"idCor"`
	DescricaoCor string `json:"descricaoCor"`
	ArquivoCor   string `json:"arquivoCor"`
}

// ProcessoDadosAbertura reúne os dados de abertura de um [Processo].
type ProcessoDadosAbertura struct {
	Info     string                                `json:"info"`
	Unidades Slice[ProcessoDadosAberturaUnidade]   `json:"unidades"`
	Lista    Slice[ProcessoDadosAberturaListaItem] `json:"lista"`
}

// ProcessoDadosAberturaUnidade representa uma unidade nos dados de abertura de um [Processo].
type ProcessoDadosAberturaUnidade struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

// ProcessoDadosAberturaListaItem representa um item da lista nos dados de abertura de um [Processo].
type ProcessoDadosAberturaListaItem struct {
	Sigla string `json:"sigla"`
}

// ProcessoAnotacao representa uma anotação de um [Processo].
type ProcessoAnotacao struct {
	IDAnotacao    string `json:"idAnotacao"`
	IDProtocolo   string `json:"idProtocolo"`
	Descricao     string `json:"descricao"`
	IDUnidade     string `json:"idUnidade"`
	IDUsuario     string `json:"idUsuario"`
	DthAnotacao   string `json:"dthAnotacao"`
	SinPrioridade string `json:"sinPrioridade"`
	StaAnotacao   string `json:"staAnotacao"`
}

// ProcessoStatusFlags reúne os indicadores de situação de um [Processo].
type ProcessoStatusFlags struct {
	DocumentoSigiloso                 string `json:"documentoSigiloso"`
	DocumentoRestrito                 string `json:"documentoRestrito"`
	DocumentoNovo                     string `json:"documentoNovo"`
	DocumentoPublicado                string `json:"documentoPublicado"`
	Anotacao                          string `json:"anotacao"`
	AnotacaoPrioridade                string `json:"anotacaoPrioridade"`
	Ciencia                           string `json:"ciencia"`
	RetornoProgramado                 string `json:"retornoProgramado"`
	RetornoData                       any    `json:"retornoData"`
	RetornoAtrasado                   string `json:"retornoAtrasado"`
	ProcessoAcessadoUsuario           string `json:"processoAcessadoUsuario"`
	ProcessoAcessadoUnidade           string `json:"processoAcessadoUnidade"`
	ProcessoRemocaoSobrestamento      string `json:"processoRemocaoSobrestamento"`
	ProcessoBloqueado                 string `json:"processoBloqueado"`
	ProcessoDocumentoIncluidoAssinado string `json:"processoDocumentoIncluidoAssinado"`
	ProcessoPublicado                 string `json:"processoPublicado"`
	NivelAcessoGlobal                 string `json:"nivelAcessoGlobal"`
	PodeGerenciarCredenciais          string `json:"podeGerenciarCredenciais"`
	ProcessoAberto                    string `json:"processoAberto"`
	ProcessoEmTramitacao              string `json:"processoEmTramitacao"`
	ProcessoSobrestado                string `json:"processoSobrestado"`
	ProcessoAnexado                   string `json:"processoAnexado"`
	PodeReabrirProcesso               string `json:"podeReabrirProcesso"`
	PodeRegistrarAnotacao             string `json:"podeRegistrarAnotacao"`
	PodeRemoverSobrestamento          bool   `json:"podeRemoverSobrestamento"`
	Tipo                              string `json:"tipo"`
	ProcessoGeradoRecebido            string `json:"processoGeradoRecebido"`
}

// TipoBusca representa o tipo de busca usado em [Client.ListarProcessos].
type TipoBusca string

// Tipos de busca aceitos pelo endpoint de listagem de processos.
const (
	TipoBuscaTotal         TipoBusca = "T"
	TipoBuscaParcial       TipoBusca = "P"
	TipoBuscaResumido      TipoBusca = "R"
	TipoBuscaExterno       TipoBusca = "E"
	TipoBuscaAuditoria     TipoBusca = "A"
	TipoBuscaUnidade       TipoBusca = "U"
	TipoBuscaPersonalizado TipoBusca = "Z"
)

// ListarProcessosParams reúne os parâmetros opcionais de [Client.ListarProcessos].
//
// Campos com valor zero (0, "" ou false) são omitidos da requisição.
type ListarProcessosParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o id do processo para detalhamento.
	ID int
	// Usuario é o id do usuário de atribuição.
	Usuario int
	// Tipo é o tipo de busca.
	Tipo TipoBusca
	// ApenasMeus, quando verdadeiro, retorna apenas os processos do usuário.
	ApenasMeus bool
	// Unidade é o id da unidade.
	Unidade int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListarProcessosParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.Filter != "" {
		q.Set("filter", p.Filter)
	}
	if p.ID != 0 {
		q.Set("id", strconv.Itoa(p.ID))
	}
	if p.Usuario != 0 {
		q.Set("usuario", strconv.Itoa(p.Usuario))
	}
	if p.Tipo != "" {
		q.Set("tipo", string(p.Tipo))
	}
	if p.ApenasMeus {
		q.Set("apenasMeus", "S")
	}

	return q
}

// ListarProcessos retorna a lista de processos e o total de registros,
// aplicando os filtros e a paginação informados em params.
func (c *Client) ListarProcessos(ctx context.Context, params ListarProcessosParams) ([]Processo, int, error) {
	endpoint := c.endpoint + "/processo/listar"
	if q := params.values().Encode(); q != "" {
		endpoint += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[[]Processo]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, 0, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, 0, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	total, err := env.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("parse total %q: %w", env.Total, err)
	}

	return env.Data, total, nil
}

// ConsultarProcesso retorna os dados detalhados de um único processo
// com base no número do protocolo informado.
func (c *Client) ConsultarProcesso(ctx context.Context, protocolo int) (*Processo, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("protocolo inválido: %d", protocolo)
	}

	endpoint := fmt.Sprintf("%s/processo/%d", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[Processo]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}

// CadastrarProcessoResult representa o processo criado pelo WSSEI.
type CadastrarProcessoResult struct {
	IDProcedimento     string `json:"IdProcedimento"`
	ProtocoloFormatado string `json:"ProtocoloFormatado"`
}

// CadastrarProcessoParams reúne os parâmetros para cadastrar um processo.
type CadastrarProcessoParams struct {
	Assuntos      []CadastrarProcessoAssuntos     `json:"assuntos"`
	Interessados  []CadastrarProcessoInteressados `json:"interessados"`
	Especificacao string                          `json:"especificacao"`
	Observacoes   string                          `json:"observacoes"`
	TipoProcesso  int                             `json:"tipoProcesso"`
	NivelAcesso   NivelAcesso                     `json:"nivelAcesso"`
	HipoteseLegal int                             `json:"hipoteseLegal"`
}

// CadastrarProcessoAssuntos representa um assunto associado ao processo.
type CadastrarProcessoAssuntos struct {
	ID string `json:"id"`
}

// CadastrarProcessoInteressados representa um interessado associado ao processo.
type CadastrarProcessoInteressados struct {
	ID string `json:"id"`
}

// NivelAcesso representa o nível de acesso de um processo.
type NivelAcesso int

// Níveis de acesso aceitos pelo endpoint de cadastro de processo.
const (
	NivelAcessoPublico  NivelAcesso = 0
	NivelAcessoRestrito NivelAcesso = 1
	NivelAcessoSigiloso NivelAcesso = 2
)

// CadastrarProcesso cadastra um novo processo e retorna seus dados de identificação.
func (c *Client) CadastrarProcesso(ctx context.Context, params CadastrarProcessoParams) (*CadastrarProcessoResult, error) {
	if params.TipoProcesso <= 0 {
		return nil, fmt.Errorf("tipo processo inválido: %d", params.TipoProcesso)
	}
	if params.NivelAcesso < NivelAcessoPublico || params.NivelAcesso > NivelAcessoSigiloso {
		return nil, fmt.Errorf("nivel acesso inválido: %d", params.NivelAcesso)
	}
	if params.NivelAcesso != NivelAcessoPublico && params.HipoteseLegal <= 0 {
		return nil, fmt.Errorf("hipotese legal inválido: %d", params.HipoteseLegal)
	}
	endpoint := c.endpoint + "/processo/criar"

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[CadastrarProcessoResult]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}

// PesquisarAssuntoResult representa um assunto retornado pela pesquisa.
type PesquisarAssuntoResult struct {
	CodigoEstruturadoFormatado string `json:"codigoestruturadoformatado"`
	Descricao                  string `json:"descricao"`
	CodigoEstruturado          string `json:"codigoestruturado"`
	ID                         string `json:"id"`
}

// PesquisarAssuntoParams reúne os parâmetros opcionais de [Client.PesquisarAssunto].
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type PesquisarAssuntoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o id do assunto para detalhamento.
	ID int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p PesquisarAssuntoParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.Filter != "" {
		q.Set("filter", p.Filter)
	}
	if p.ID != 0 {
		q.Set("id", strconv.Itoa(p.ID))
	}
	return q
}

// PesquisarAssunto retorna a lista de assuntos encontrados e o total de registros,
// aplicando os filtros e a paginação informados em params.
func (c *Client) PesquisarAssunto(ctx context.Context, params PesquisarAssuntoParams) ([]PesquisarAssuntoResult, int, error) {
	endpoint := c.endpoint + "/processo/assunto/pesquisar"
	if q := params.values().Encode(); q != "" {
		endpoint += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[[]PesquisarAssuntoResult]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, 0, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, 0, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	total, err := env.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("parse total %q: %w", env.Total, err)
	}

	return env.Data, total, nil
}

// PesquisarProcessoResult representa um processo retornado pela pesquisa geral.
type PesquisarProcessoResult struct {
	IDProcedimento                 string                             `json:"idProcedimento"`
	IDTipoProcedimento             string                             `json:"idTipoProcedimento"`
	NomeTipoProcedimento           string                             `json:"nomeTipoProcedimento"`
	SiglaUnidadeGeradora           string                             `json:"siglaUnidadeGeradora"`
	IDUnidadeGeradora              string                             `json:"idUnidadeGeradora"`
	ProtocoloFormatadoProcedimento string                             `json:"protocoloFormatadoProcedimento"`
	IDUsuarioGerador               string                             `json:"idUsuarioGerador"`
	NomeUsuarioGerador             string                             `json:"nomeUsuarioGerador"`
	SiglaUsuarioGerador            string                             `json:"siglaUsuarioGerador"`
	DataGeracao                    string                             `json:"dataGeracao"`
	Documento                      Object[PesquisarProcessoDocumento] `json:"documento"`
}

// PesquisarProcessoDocumento representa o documento associado ao resultado
// da pesquisa geral de processos.
type PesquisarProcessoDocumento struct {
	IDDocumento                 string `json:"idDocumento"`
	IDSerieDocumento            string `json:"idSerieDocumento"`
	NomeSerieDocumento          string `json:"nomeSerieDocumento"`
	ProtocoloFormatadoDocumento string `json:"protocoloFormatadoDocumento"`
	NumeroDocumento             string `json:"numeroDocumento"`
	StaDocumento                string `json:"staDocumento"`
	DtaGeracao                  string `json:"dtaGeracao"`
	DadosAnexo                  any    `json:"dadosAnexo"`
}

// StaTipoData representa o tipo de busca por data em [PesquisarProcessoParams].
type StaTipoData int

// Tipos de busca por data aceitos pelo endpoint de pesquisa geral de processos.
const (
	// StaTipoDataPeriodo habilita o período informado em DataInicio e DataFim.
	StaTipoDataPeriodo StaTipoData = 0
	// StaTipoData30Dias filtra os últimos 30 dias.
	StaTipoData30Dias StaTipoData = 30
	// StaTipoData60Dias filtra os últimos 60 dias.
	StaTipoData60Dias StaTipoData = 60
)

// PesquisarProcessoParams reúne os parâmetros opcionais de [Client.PesquisarProcesso].
//
// Campos com valor zero (0 ou "") são omitidos da requisição, com exceção
// de StaTipoData, que é incluído sempre que DataInicio e DataFim são informados.
type PesquisarProcessoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Grupo é o id do grupo de acompanhamento.
	Grupo int
	// PalavrasChave são as palavras-chave usadas na pesquisa.
	PalavrasChave string
	// Descricao é o texto de descrição a ser pesquisado.
	Descricao string
	// StaTipoData define o tipo de busca por data. É obrigatório quando
	// DataInicio e DataFim são informados.
	StaTipoData StaTipoData
	// DataInicio é a data de início do período (formato dd/mm/aaaa).
	DataInicio string
	// DataFim é a data de término do período (formato dd/mm/aaaa).
	DataFim string
	// IDUnidadeGeradora é o id da unidade geradora.
	IDUnidadeGeradora int
	// IDAssunto é o id do assunto.
	IDAssunto int
	// BuscaRapida é o texto para comportamento igual à busca rápida
	// (não deve ser combinado com os outros filtros).
	BuscaRapida string
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p PesquisarProcessoParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.Grupo != 0 {
		q.Set("grupo", strconv.Itoa(p.Grupo))
	}
	if p.PalavrasChave != "" {
		q.Set("palavrasChave", p.PalavrasChave)
	}
	if p.Descricao != "" {
		q.Set("descricao", p.Descricao)
	}
	if p.DataInicio != "" {
		q.Set("dataInicio", p.DataInicio)
	}
	if p.DataFim != "" {
		q.Set("dataFim", p.DataFim)
	}
	if p.DataInicio != "" || p.DataFim != "" {
		q.Set("staTipoData", strconv.Itoa(int(p.StaTipoData)))
	}
	if p.IDUnidadeGeradora != 0 {
		q.Set("idUnidadeGeradora", strconv.Itoa(p.IDUnidadeGeradora))
	}
	if p.IDAssunto != 0 {
		q.Set("idAssunto", strconv.Itoa(p.IDAssunto))
	}
	if p.BuscaRapida != "" {
		q.Set("buscaRapida", p.BuscaRapida)
	}
	return q
}

// PesquisarProcesso retorna a lista de processos encontrados e o total de
// registros, aplicando os filtros e a paginação informados em params.
func (c *Client) PesquisarProcesso(ctx context.Context, params PesquisarProcessoParams) ([]PesquisarProcessoResult, int, error) {
	endpoint := c.endpoint + "/processo/pesquisar"
	if q := params.values().Encode(); q != "" {
		endpoint += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[[]PesquisarProcessoResult]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, 0, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, 0, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	total, err := env.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("parse total %q: %w", env.Total, err)
	}

	return env.Data, total, nil
}

// TipoProcessoResult representa um tipo de processo retornado pela listagem.
//
// Este endpoint está marcado como deprecated pela API e pode ser removido
// em versões futuras.
type TipoProcessoResult struct {
	ID              string `json:"id"`
	Nome            string `json:"nome"`
	PermiteSigiloso bool   `json:"permiteSigiloso"`
}

// TipoProcessoParams reúne os parâmetros opcionais de
// [Client.ListarTipoProcesso].
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type TipoProcessoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o id do tipo de processo para detalhamento.
	ID int
	// Favoritos retorna apenas os tipos de processo favoritos.
	Favoritos Favoritos
	// Aplicabilidade filtra os tipos de processo pela aplicabilidade.
	Aplicabilidade []Aplicabilidade
}

// Favoritos representa o filtro para retornar apenas os tipos de processo
// favoritos.
type Favoritos string

// Tipos de busca aceitos pelo endpoint de listagem de tipos de processos.
const (
	Favorito Favoritos = "S"
)

// Aplicabilidade representa a aplicabilidade de um tipo de processo.
type Aplicabilidade string

// Tipos de busca aceitos pelo endpoint de listagem de tipos de processos.
const (
	AplicabilidadeInternoExterno Aplicabilidade = "T"
	AplicabilidadeInterno        Aplicabilidade = "I"
	AplicabilidadeExterno        Aplicabilidade = "E"
	AplicabilidadeFormulario     Aplicabilidade = "F"
)

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p TipoProcessoParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.Filter != "" {
		q.Set("filter", p.Filter)
	}
	if p.ID != 0 {
		q.Set("id", strconv.Itoa(p.ID))
	}
	if p.Favoritos != "" {
		q.Set("favoritos", string(p.Favoritos))
	}

	if len(p.Aplicabilidade) > 0 {
		var valores []string
		for _, aplicabilidade := range p.Aplicabilidade {
			valores = append(valores, string(aplicabilidade))
		}
		q.Set("aplicabilidade", strings.Join(valores, ","))
	}

	return q
}

// ListarTipoProcesso retorna a lista de tipos de processo e o total de
// registros, aplicando os filtros e a paginação informados em params.
//
// Este endpoint está marcado como deprecated pela API e pode ser removido
// em versões futuras.
func (c *Client) ListarTipoProcesso(ctx context.Context, params TipoProcessoParams) ([]TipoProcessoResult, int, error) {
	endpoint := c.endpoint + "/processo/tipo/listar"
	if q := params.values().Encode(); q != "" {
		endpoint += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[[]TipoProcessoResult]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, 0, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, 0, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	total, err := env.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("parse total %q: %w", env.Total, err)
	}

	return env.Data, total, nil
}

// ProcessoInteressados representa um interessado retornado pela listagem
// de interessados de um [Processo].
type ProcessoInteressados struct {
	ID            string `json:"id"`
	Nome          string `json:"nome"`
	Sigla         string `json:"sigla"`
	NomeFormatado string `json:"nomeformatado"`
}

// ListarInteressadosProcessoParams reúne os parâmetros opcionais de
// [Client.ListarInteressadosProcesso].
//
// Campos com valor zero (0) são omitidos da requisição.
type ListarInteressadosProcessoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListarInteressadosProcessoParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	return q
}

// ListarInteressadosProcesso retorna a lista de interessados de um processo
// e o total de registros, aplicando a paginação informada em params.
func (c *Client) ListarInteressadosProcesso(ctx context.Context, protocolo int, params ListarInteressadosProcessoParams) ([]ProcessoInteressados, int, error) {
	if protocolo <= 0 {
		return nil, 0, fmt.Errorf("protocolo inválido: %d", protocolo)
	}

	endpoint := fmt.Sprintf("%s/processo/%d/interessados/listar", c.endpoint, protocolo)
	if q := params.values().Encode(); q != "" {
		endpoint += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[[]ProcessoInteressados]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, 0, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, 0, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	total, err := env.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("parse total %q: %w", env.Total, err)
	}

	return env.Data, total, nil
}

// ListaSugestaoAssunto tipo utilizado funcao ListarSugestaoAssuntosProcesso
type ListaSugestaoAssunto struct {
	CodigoEstruturadoFormatado string `json:"codigoestruturadoformatado"`
	Descricao                  string `json:"descricao"`
	CodigoEstruturado          string `json:"codigoestruturado"`
	ID                         string `json:"id"`
}

// ListaAssuntosProcessoParams parametros da query da funcao ListarSugestaoAssuntosProcesso
type ListaAssuntosProcessoParams struct {
	// TipoProcedimento obrigatorio
	TipoProcedimento int
	Limit            int
	Start            int
	Filter           string
	ID               int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListaAssuntosProcessoParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Start))
	}
	if p.Filter != "" {
		q.Set("tipoprocedimento", p.Filter)
	}
	if p.ID != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.ID))
	}
	return q
}

// ListarSugestaoAssuntosProcesso Retorna a lista de Sugestão de Assuntos pelo Tipo de Processo.
func (c Client) ListarSugestaoAssuntosProcesso(ctx context.Context, params ListaAssuntosProcessoParams) ([]ListaSugestaoAssunto, int, error) {
	if params.TipoProcedimento == 0 {
		return nil, 0, fmt.Errorf("invalid tipo-procedimento : %d", params.TipoProcedimento)
	}

	url := fmt.Sprintf("%s/processo/assunto/sugestao/%d/listar", c.endpoint, params.TipoProcedimento)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]ListaSugestaoAssunto]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("consulta failed %d: %s", params.TipoProcedimento, result.Mensagem)
	}

	total, err := strconv.Atoi(result.Total)
	if err != nil {
		return nil, 0, fmt.Errorf("error: %w", err)
	}
	return result.Data, total, nil
}

// CienciasProcesso tipo utilizado funcao ListarCiencasProcesso.
type CienciasProcesso struct {
	IDProtocolo  string `json:"idProtocolo"`
	IDAtividade  string `json:"idAtividade"`
	Data         string `json:"data"`
	IDUnidade    string `json:"idUnidade"`
	Unidade      string `json:"unidade"`
	SiglaUnidade string `json:"siglaUnidade"`
	IDUsuario    string `json:"idUsuario"`
	SiglaUsuario string `json:"siglaUsuario"`
	NomeUsuario  string `json:"nomeUsuario"`
	Descricao    string `json:"descricao"`
}

// ListaCienciaProcessoParams parametros da query da funcao ListarCiencasProcesso.
type ListaCienciaProcessoParams struct {
	// Protocolo obrigatorio
	Protocolo int
	Limit     int
	Start     int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListaCienciaProcessoParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Start))
	}
	return q
}

// ListarCiencasProcesso funcao Retorna a lista de Ciências do Processo.
func (c Client) ListarCiencasProcesso(ctx context.Context, params ListaCienciaProcessoParams) ([]CienciasProcesso, int, error) {
	if params.Protocolo == 0 {
		return nil, 0, fmt.Errorf("invalid protocolo : %d", params.Protocolo)
	}

	url := fmt.Sprintf("%s/processo/%d/ciencia/listar", c.endpoint, params.Protocolo)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]CienciasProcesso]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("consulta failed %d: %s", params.Protocolo, result.Mensagem)
	}

	total, err := strconv.Atoi(result.Total)
	if err != nil {
		return nil, 0, fmt.Errorf("error: %w", err)
	}
	return result.Data, total, nil
}

// ListaMeusAcompanhamentosParams parametros da query da funcao ListarCiencasProcesso.
type ListaMeusAcompanhamentosParams struct {
	// Protocolo obrigatorio
	Limit   int
	Start   int
	Grupo   int
	Usuario int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListaMeusAcompanhamentosParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Start))
	}
	if p.Grupo != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Grupo))
	}
	if p.Usuario != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Usuario))
	}
	return q
}

// ListarMeusAcompanhamentos funcao Retorna a lista de Ciências do Processo.
func (c Client) ListarMeusAcompanhamentos(ctx context.Context, params ListaMeusAcompanhamentosParams) ([]Processo, int, error) {
	url := fmt.Sprintf("%s/processo/listar/meus/acompanhamentos", c.endpoint)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]Processo]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("consulta failed: %s", result.Mensagem)
	}

	total, err := strconv.Atoi(result.Total)
	if err != nil {
		return nil, 0, fmt.Errorf("error: %w", err)
	}
	return result.Data, total, nil
}

// ListaAcompanhamentosParams parametros da query da funcao ListarAcompanhamentos.
type ListaAcompanhamentosParams struct {
	// Grupo obrigatorio
	Grupo int
	Limit int
	Start int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListaAcompanhamentosParams) values() url.Values {
	q := make(url.Values)
	if p.Grupo != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Grupo))
	}
	if p.Limit != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Start))
	}
	return q
}

// ListarAcompanhamentos Retorna a lista de Processos Acompanhados na Unidade.
func (c Client) ListarAcompanhamentos(ctx context.Context, params ListaAcompanhamentosParams) ([]Processo, int, error) {
	if params.Grupo == 0 {
		return nil, 0, fmt.Errorf("invalid Grupo: %d", params.Grupo)
	}

	url := fmt.Sprintf("%s/processo/listar/acompanhamentos", c.endpoint)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]Processo]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("consulta failed: %s", result.Mensagem)
	}

	total, err := strconv.Atoi(result.Total)
	if err != nil {
		return nil, 0, fmt.Errorf("error: %w", err)
	}
	return result.Data, total, nil
}

// ListaCredenciamentoParams parametros da query da funcao ListarAcompanhamentos.
type ListaCredenciamentoParams struct {
	// Procedimento obrigatorio
	Procedimento int
	Limit        int
	Start        int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListaCredenciamentoParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("tipoprocedimento", strconv.Itoa(p.Start))
	}
	return q
}

// ListaCredenciamento tipo utilizado na funcao ListarCredenciamentoProcesso.
type ListaCredenciamento struct {
	Atividade         string `json:"atividade"`
	SiglaUsuario      string `json:"siglaUsuario"`
	SiglaUnidade      string `json:"siglaUnidade"`
	NomeUsuario       string `json:"nomeUsuario"`
	DescricaoUnidade  string `json:"descricaoUnidade"`
	DataAbertura      string `json:"dataAbertura"`
	CredencialCassada bool   `json:"credencialCassada"`
	DataCassacao      string `json:"dataCassacao"`
}

// ListarCredenciamentoProcesso Retorna a lista de Usuário com Credênciais de Acesso ao Processo.
func (c Client) ListarCredenciamentoProcesso(ctx context.Context, params ListaCredenciamentoParams) ([]ListaCredenciamento, int, error) {
	if params.Procedimento == 0 {
		return nil, 0, fmt.Errorf("invalid Procedimento: %d", params.Procedimento)
	}

	url := fmt.Sprintf("%s/processo/%d/credenciamento/listar", c.endpoint, params.Procedimento)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]ListaCredenciamento]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("consulta failed: %s", result.Mensagem)
	}

	total, err := strconv.Atoi(result.Total)
	if err != nil {
		return nil, 0, fmt.Errorf("error: %w", err)
	}
	return result.Data, total, nil
}

// ConsultarAtribuicaoProcesso Retorna os dados da Atribuição do Processo.
func (c Client) ConsultarAtribuicaoProcesso(ctx context.Context, protocolo int) (*ProcessoUsuarioAtribuido, error) {
	if protocolo == 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}

	url := fmt.Sprintf("%s/processo/%d/consultar/atribuicao", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[ProcessoUsuarioAtribuido]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, fmt.Errorf("consulta failed: %s", result.Mensagem)
	}

	return &result.Data, nil
}

// AcessoProcesso tipo utilizado na funcao VerificarAcessoProcesso.
type AcessoProcesso struct {
	AcessoLiberado     bool `json:"acessoLiberado"`
	ChamarAutenticacao bool `json:"chamarAutenticacao"`
}

// VerificarAcessoProcesso Verifica se o Usuário pode acessar o Processo.
func (c Client) VerificarAcessoProcesso(ctx context.Context, protocolo int) (*AcessoProcesso, error) {
	if protocolo == 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}

	url := fmt.Sprintf("%s/processo/verifica/acesso/%d", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[AcessoProcesso]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, fmt.Errorf("consulta failed: %s", result.Mensagem)
	}

	return &result.Data, nil
}

// ConsultaAcompanhamento tipo utilizado na funcao ConsultarAcompanhamentosProcesso.
type ConsultaAcompanhamento struct {
	IDAcompanhamento      string `json:"idAcompanhamento"`
	IDProtocolo           string `json:"idProtocolo"`
	IDUnidade             string `json:"idUnidade"`
	Observacao            string `json:"observacao"`
	Visualizacao          string `json:"visualizacao"`
	IDGrupoAcompanhamento string `json:"idGrupoAcompanhamento"`
	NomeGrupo             string `json:"nomeGrupo"`
}

// ConsultarAcompanhamentosProcesso Retorna os dados de Acompanhamento do Processo.
func (c Client) ConsultarAcompanhamentosProcesso(ctx context.Context, protocolo int) (*ConsultaAcompanhamento, error) {
	if protocolo == 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}

	url := fmt.Sprintf("%s/processo/acompanhamento/consultar?%d", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[ConsultaAcompanhamento]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, fmt.Errorf("consulta failed: %s", result.Mensagem)
	}

	return &result.Data, nil
}

// SobrestoParams tipo utilizado na funcao SobrestarProcesso.
type SobrestoParams struct {
	ProtocoloDestino int
	// Motivo obrigatorio
	Motivo string
}

// PostProcesso tipo utilizado na funcao SobrestarProcesso.
type PostProcesso struct {
	Mensagem string `json:"mensagem"`
	Total    string `json:"total"`
}

// SobrestarProcesso Realiza o Sobrestamento do Processo.
func (c *Client) SobrestarProcesso(ctx context.Context, protocolo int, params SobrestoParams) (*PostProcesso, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}
	if strings.TrimSpace(params.Motivo) == "" {
		return nil, fmt.Errorf("Motivo required")
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/%d/sobrestar/processo", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// AcessoParams tipo utilizado na funcao ConcederAcesso.
type AcessoParams struct {
	// Todos sao obrigatorios
	Unidade int
	Usuario int
}

// ConcederAcesso Concede Acesso a um Usuário no Processo.
func (c *Client) ConcederAcesso(ctx context.Context, procedimento int, params AcessoParams) (*PostProcesso, error) {
	if procedimento <= 0 {
		return nil, fmt.Errorf("invalid Procedimento: %d", procedimento)
	}
	if params.Unidade <= 0 {
		return nil, fmt.Errorf("invalid Unidade: %d", params.Unidade)
	}
	if params.Usuario <= 0 {
		return nil, fmt.Errorf("invalid Usuario: %d", params.Usuario)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/%d/credenciamento/conceder", c.endpoint, procedimento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// RenunciarAcesso O Usuário Renuncia das credênciais de acesso ao Processo.
func (c *Client) RenunciarAcesso(ctx context.Context, procedimento int) (*PostProcesso, error) {
	if procedimento <= 0 {
		return nil, fmt.Errorf("invalid Procedimento: %d", procedimento)
	}

	url := fmt.Sprintf("%s/processo/%d/credenciamento/renunciar", c.endpoint, procedimento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// CassarAcessoParams tipo utilizado na funcao CassarAcesso.
type CassarAcessoParams struct {
	// Todos sao obrigatorios
	Atividade int
}

// CassarAcesso Inativa o Acesso de um Usuário a um Processo.
func (c *Client) CassarAcesso(ctx context.Context, procedimento int, params CassarAcessoParams) (*PostProcesso, error) {
	if procedimento <= 0 {
		return nil, fmt.Errorf("invalid Procedimento: %d", procedimento)
	}
	if params.Atividade <= 0 {
		return nil, fmt.Errorf("invalid Unidade: %d", params.Atividade)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/%d/credenciamento/cassar", c.endpoint, procedimento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// EniviarProcessoParams tipo utilizado na funcao EnviarProcesso.
type EniviarProcessoParams struct {
	// NumeroProcesso obrigatorio
	NumeroProcesso string
	// UnidadesDestino obrigatorio
	UnidadesDestino               string
	SinManterAbertoUnidade        string
	SinRemoverAnotacao            string
	SinEnviarEmailNotificacao     string
	DataRetornoProgramado         string
	DiasRetornoProgramado         string
	SinDiasUteisRetornoProgramado string
	SinReabrir                    string
}

// EnviarProcesso Envia o Processo para outras Unidades.
func (c *Client) EnviarProcesso(ctx context.Context, params EniviarProcessoParams) (*PostProcesso, error) {
	if strings.TrimSpace(params.NumeroProcesso) == "" {
		return nil, fmt.Errorf("NumeroProcesso required")
	}
	if strings.TrimSpace(params.UnidadesDestino) == "" {
		return nil, fmt.Errorf("UnidadesDestino required")
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/enviar", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// ConcluirProcessoParams tipo utilizado na funcao ConcluirProcesso.
type ConcluirProcessoParams struct {
	// NumeroProcesso obrigatorio
	NumeroProcesso string
}

// ConcluirProcesso Conclui o Processo.
func (c *Client) ConcluirProcesso(ctx context.Context, params ConcluirProcessoParams) (*PostProcesso, error) {
	if strings.TrimSpace(params.NumeroProcesso) == "" {
		return nil, fmt.Errorf("NumeroProcesso required")
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/concluir", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// ReceberProcessoParams tipo utilizado na funcao ReceberProcesso.
type ReceberProcessoParams struct {
	// Procedimento obrigatorio
	Procedimento int
}

// ReceberProcesso Recebe um Processo de outra Unidade.
func (c *Client) ReceberProcesso(ctx context.Context, params ReceberProcessoParams) (*PostProcesso, error) {
	if params.Procedimento <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", params.Procedimento)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/receber", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// AcessoProcessoParams tipo utilizado na funcao AcessoProcesso.
type AcessoProcessoParams struct {
	// Todos obrigatorios
	Protocolo int
	Senha     string
}

// AcessoProcesso Realiza a verificação de acesso do Usuário ao Processo.
func (c *Client) AcessoProcesso(ctx context.Context, params AcessoProcessoParams) (*PostProcesso, error) {
	if params.Protocolo <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", params.Protocolo)
	}
	if strings.TrimSpace(params.Senha) == "" {
		return nil, fmt.Errorf("invalid senha: %s", params.Senha)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/identificar/acesso", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// AcompanharProcessoParams tipo utilizado na funcao AcompanharProcesso.
type AcompanharProcessoParams struct {
	// Protocolo obrigatorio
	Protocolo int
	// Grupo obrigatorio
	Grupo      int
	observacao string
}

// AcompanharProcesso Realiza o Acompanhamento do Processo.
func (c *Client) AcompanharProcesso(ctx context.Context, params AcompanharProcessoParams) (*PostProcesso, error) {
	if params.Protocolo <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", params.Protocolo)
	}
	if params.Grupo <= 0 {
		return nil, fmt.Errorf("invalid Grupo: %d", params.Grupo)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/acompanhar", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// AlterarAcompanhamentoProcessoParams tipo utilizado na funcao AlterarAcompanhamentoProcesso.
type AlterarAcompanhamentoProcessoParams struct {
	// Protocolo obrigatorio
	Protocolo int
	// Grupo obrigatorio
	Grupo      int
	observacao string
}

// AlterarAcompanhamentoProcesso Realiza a edição do Acompanhamento de um Processo.
func (c *Client) AlterarAcompanhamentoProcesso(ctx context.Context, params AcompanharProcessoParams) (*PostProcesso, error) {
	if params.Protocolo <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", params.Protocolo)
	}
	if params.Grupo <= 0 {
		return nil, fmt.Errorf("invalid Grupo: %d", params.Grupo)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/acompanhamento/alterar", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// ReabrirProcessoParams tipo utilizado na funcao ReabrirProcesso.
type ReabrirProcessoParams struct {
	// NumeroProcesso obrigatorio
	NumeroProcesso string
}

// ReabrirProcesso Reabre o Processo na Unidade.
func (c *Client) ReabrirProcesso(ctx context.Context, procedimento int, params ReabrirProcessoParams) (*PostProcesso, error) {
	if procedimento <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", procedimento)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/reabrir/%d", c.endpoint, procedimento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// AtribuirProcessoParams tipo utilizado na funcao AtribuirProcesso.
type AtribuirProcessoParams struct {
	// Todos obrigatorios
	NumeroProcesso string
	Usuario        int
}

// AtribuirProcesso Atribui um Processo a um Usuário.
func (c *Client) AtribuirProcesso(ctx context.Context, params AtribuirProcessoParams) (*PostProcesso, error) {
	if strings.TrimSpace(params.NumeroProcesso) == "" {
		return nil, fmt.Errorf("invalid NumeroProcesso: %s", params.NumeroProcesso)
	}
	if params.Usuario <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", params.Usuario)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/atribuir", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// RemoverAtribuicaoProcesso Remove a Atribuição de um Processo.
func (c *Client) RemoverAtribuicaoProcesso(ctx context.Context, protocolo int) (*PostProcesso, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}
	url := fmt.Sprintf("%s/processo/%d/remover/atribuicao", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// CancelarSobrestamentoProcesso Cancela o Sobrestamento do Processo.
func (c *Client) CancelarSobrestamentoProcesso(ctx context.Context, protocolo int) (*PostProcesso, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}
	url := fmt.Sprintf("%s/processo/%d/cancelar/sobrestamento", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// ExcluirAcompanhamento Exclui o Acompanhamento de um Processo.
func (c *Client) ExcluirAcompanhamento(ctx context.Context, acompanhamento int) (*PostProcesso, error) {
	if acompanhamento <= 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", acompanhamento)
	}
	url := fmt.Sprintf("%s/processo/acompanhamento/%d/excluir", c.endpoint, acompanhamento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// DarCienciaProcesso Da Ciencia no Processo.
func (c *Client) DarCienciaProcesso(ctx context.Context, procedimento int) (*PostProcesso, error) {
	if procedimento <= 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", procedimento)
	}
	url := fmt.Sprintf("%s/processo/%d/ciencia", c.endpoint, procedimento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// AlterarProcessoParams tipo utilizado na funcao AlterarProcesso.
type AlterarProcessoParams struct {
	Assuntos      string
	Interessados  string
	Especificacao string
	Observacao    string
	// IDTipoProcesso obrigatorio
	IDTipoProcesso int
	// NivelACesso obrigatorio
	NivelACesso     int
	IDHipoteseLegal int
	// GrauSigilo obrigatorio
	GrauSigilo string
}

// AlterarProcesso Realiza a edição de um Processo.
func (c *Client) AlterarProcesso(ctx context.Context, protocolo int, params AlterarProcessoParams) (*PostProcesso, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}
	if params.IDTipoProcesso <= 0 {
		return nil, fmt.Errorf("invalid IDTipoProcesso: %d", params.IDTipoProcesso)
	}
	if params.NivelACesso <= 0 {
		return nil, fmt.Errorf("invalid NivelACesso: %d", params.NivelACesso)
	}
	if strings.TrimSpace(params.GrauSigilo) == "" {
		return nil, fmt.Errorf("invalid GrauSigilo: %s", params.GrauSigilo)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/processo/%d/alterar", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[PostProcesso]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}
