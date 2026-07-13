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
	Marcador         Object[ProcessoMarcador]         `json:"marcador"`
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
