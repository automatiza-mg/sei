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

// DocumentoBlocoAssinatura representa o conteúdo de um bloco de assinatura,
// reunindo as permissões do usuário sobre o bloco e a lista de documentos.
type DocumentoBlocoAssinatura struct {
	Permissoes DocumentoBlocoAssinaturaPermissoes   `json:"permissoes"`
	Dados      Slice[DocumentoBlocoAssinaturaDados] `json:"dados"`
}

// DocumentoBlocoAssinaturaPermissoes reúne as permissões do usuário logado
// sobre um [DocumentoBlocoAssinatura].
type DocumentoBlocoAssinaturaPermissoes struct {
	Assinar bool `json:"assinar"`
	Retirar bool `json:"retirar"`
	Anotar  bool `json:"anotar"`
}

// DocumentoBlocoAssinaturaDados representa um documento dentro de um bloco
// de assinatura.
type DocumentoBlocoAssinaturaDados struct {
	Sequencia          string                                     `json:"sequencia"`
	ID                 string                                     `json:"id"`
	Aberto             string                                     `json:"aberto"`
	Data               string                                     `json:"data"`
	IDDocumento        string                                     `json:"idDocumento"`
	IDProcesso         string                                     `json:"idProcesso"`
	NomeTipoProcesso   string                                     `json:"nomeTipoProcesso"`
	ProtocoloFormatado string                                     `json:"protocoloFormatado"`
	NumeroDocumento    string                                     `json:"numeroDocumento"`
	TipoDocumento      string                                     `json:"tipoDocumento"`
	Assinaturas        Slice[DocumentoBlocoAssinaturaAssinaturas] `json:"assinaturas"`
	Anotacao           string                                     `json:"anotacao"`
}

// DocumentoBlocoAssinaturaAssinaturas representa uma assinatura registrada
// em um [DocumentoBlocoAssinaturaDados].
type DocumentoBlocoAssinaturaAssinaturas struct {
	Nome      string `json:"nome"`
	Cargo     string `json:"cargo"`
	IDUsuario string `json:"idUsuario"`
}

// ListarDocumentosBlocoAssinatura retorna os documentos, as permissões do
// usuário e o total de registros de um bloco de assinatura.
func (c *Client) ListarDocumentosBlocoAssinatura(ctx context.Context, bloco int) (*DocumentoBlocoAssinatura, int, error) {
	endpoint := fmt.Sprintf("%s/bloco/assinatura/%d/documentos/listar", c.endpoint, bloco)

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

	var env Envelope[DocumentoBlocoAssinatura]
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

	return &env.Data, total, nil
}

// AssinarBlocoAssinaturaParams reúne os parâmetros para
// para assinar todos os documentos de um bloco de assinatura.
type AssinarBlocoAssinaturaParams struct {
	Orgao   int    `json:"orgao"`
	Cargo   string `json:"cargo"`
	Login   string `json:"login"`
	Senha   string `json:"senha"`
	Usuario int    `json:"usuario"`
}

// AssinarBlocoAssinatura realiza a assinatura de todos os documentos do
// bloco de assinatura identificado por bloco
func (c *Client) AssinarBlocoAssinatura(ctx context.Context, bloco int, params AssinarBlocoAssinaturaParams) error {
	endpoint := fmt.Sprintf("%s/bloco/assinatura/%d/assinar", c.endpoint, bloco)

	// O WSSEI (PHP legado) interpreta o valor de "cargo" como Latin-1
	// mesmo recebendo JSON. Transcodamos apenas esse campo para bytes
	// Latin-1 antes de montar o body.
	jsonBody, err := json.Marshal(struct {
		Orgao   int             `json:"orgao"`
		Cargo   json.RawMessage `json:"cargo"`
		Login   string          `json:"login"`
		Senha   string          `json:"senha"`
		Usuario int             `json:"usuario"`
	}{
		Orgao:   params.Orgao,
		Cargo:   jsonStringLatin1(params.Cargo),
		Login:   params.Login,
		Senha:   params.Senha,
		Usuario: params.Usuario,
	})
	if err != nil {
		return fmt.Errorf("json body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[struct{}]
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return nil
}

// BlocoAssinatura representa um bloco de assinatura retornado pelo WSSEI.
type BlocoAssinatura struct {
	ID        string                   `json:"id"`
	Atributos BlocoAssinaturaAtributos `json:"atributos"`
}

// BlocoAssinaturaAtributos reúne os atributos detalhados de um [BlocoAssinatura].
type BlocoAssinaturaAtributos struct {
	IDBloco          string                        `json:"idBloco"`
	IDUnidade        string                        `json:"idUnidade"`
	SiglaUnidade     string                        `json:"siglaUnidade"`
	Estado           string                        `json:"estado"`
	Descricao        string                        `json:"descricao"`
	Unidades         Slice[BlocoAssinaturaUnidade] `json:"unidades"`
	NumeroDocumentos string                        `json:"numeroDocumentos"`
}

// BlocoAssinaturaUnidade representa uma unidade vinculada a um [BlocoAssinatura].
type BlocoAssinaturaUnidade struct {
	IDUnidade string `json:"idUnidade"`
	Unidade   string `json:"unidade"`
}

// EstadoSituacao representa o estado de um bloco de assinatura.
type EstadoSituacao string

// Estados aceitos pelo endpoint de pesquisa de blocos de assinatura.
const (
	EstadoSituacaoAberto          EstadoSituacao = "A"
	EstadoSituacaoDisponibilizado EstadoSituacao = "D"
	EstadoSituacaoRetornado       EstadoSituacao = "R"
	EstadoSituacaoConcluido       EstadoSituacao = "C"
)

// PesquisarBlocoAssinaturaParams reúne os parâmetros opcionais de
// [Client.PesquisarBlocoAssinatura].
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type PesquisarBlocoAssinaturaParams struct {
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
	// Estado é o estado das situações.
	Estado EstadoSituacao
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p PesquisarBlocoAssinaturaParams) values() url.Values {
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
	if p.Estado != "" {
		q.Set("estado", string(p.Estado))
	}
	return q
}

// PesquisarBlocoAssinatura retorna a lista de blocos de assinatura e o total
// de registros, aplicando os filtros e a paginação informados em params.
func (c *Client) PesquisarBlocoAssinatura(ctx context.Context, params PesquisarBlocoAssinaturaParams) ([]BlocoAssinatura, int, error) {
	endpoint := c.endpoint + "/bloco/assinatura/pesquisar"
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

	var env Envelope[[]BlocoAssinatura]
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

// BlocoInterno representa um bloco interno retornado pelo WSSEI.
type BlocoInterno struct {
	ID              string `json:"id"`
	IDUnidade       string `json:"idUnidade"`
	SiglaUnidade    string `json:"siglaUnidade"`
	Estado          string `json:"estado"`
	Descricao       string `json:"descricao"`
	NumeroProcessos string `json:"numeroProcessos"`
}

// BlocoInternoParams reúne os parâmetros opcionais de [Client.PesquisarBlocoInterno].
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type BlocoInternoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o id do bloco para detalhamento.
	ID int
	// Estado representa a situação do bloco.
	Estado EstadoSituacao
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p BlocoInternoParams) values() url.Values {
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
	if p.Estado != "" {
		q.Set("estado", string(p.Estado))
	}
	return q
}

// PesquisarBlocoInterno retorna a lista de blocos internos encontrados e o total
// de registros, aplicando os filtros e a paginação informados em params.
func (c *Client) PesquisarBlocoInterno(ctx context.Context, params BlocoInternoParams) ([]BlocoInterno, int, error) {
	endpoint := c.endpoint + "/bloco/interno/pesquisar"
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

	var env Envelope[[]BlocoInterno]
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

// CadastrarBlocoInternoResult representa os dados retornados após o cadastro
// de um bloco interno.
type CadastrarBlocoInternoResult struct {
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
}

// CadastrarBlocoInternoParams reúne os parâmetros necessários para
// cadastrar um bloco interno.
type CadastrarBlocoInternoParams struct {
	Descricao string `json:"descricao"`
}

// blocoInternoRequest executa a requisição de cadastro ou alteração
// de um bloco interno.
func (c *Client) blocoInternoRequest(ctx context.Context, params CadastrarBlocoInternoParams, endpoint string) (*CadastrarBlocoInternoResult, error) {
	if strings.TrimSpace(params.Descricao) == "" {
		return nil, fmt.Errorf("descricao required")
	}
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

	var env Envelope[CadastrarBlocoInternoResult]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}

// CadastrarBlocoInterno cadastra um novo bloco interno e retorna
// os dados do bloco criado.
func (c *Client) CadastrarBlocoInterno(ctx context.Context, params CadastrarBlocoInternoParams) (*CadastrarBlocoInternoResult, error) {
	endpoint := c.endpoint + "/bloco/interno/criar"
	return c.blocoInternoRequest(ctx, params, endpoint)
}

// AlterarBlocoInterno altera um bloco interno existente e retorna
// os dados atualizados do bloco.
func (c *Client) AlterarBlocoInterno(ctx context.Context, params CadastrarBlocoInternoParams, bloco int) (*CadastrarBlocoInternoResult, error) {
	if bloco <= 0 {
		return nil, fmt.Errorf("bloco invalido: %d", bloco)
	}
	endpoint := fmt.Sprintf("%s/bloco/interno/%d/alterar", c.endpoint, bloco)
	return c.blocoInternoRequest(ctx, params, endpoint)
}

// BlocoAssinaturaResult representa os dados retornados após o cadastro
// de um bloco de assinatura.
type BlocoAssinaturaResult struct {
	ID        string        `json:"id"`
	Descricao string        `json:"descricao"`
	Unidades  Slice[string] `json:"unidades"`
}

// BlocoAssinaturaParams reúne os parâmetros necessários para
// cadastrar um bloco de assinatura.
type BlocoAssinaturaParams struct {
	Descricao string
	Unidades  []string
}

// blocoAssinaturaParams representa o formato esperado pelo WSSEI
// para o corpo da requisição.
type blocoAssinaturaParams struct {
	Descricao string `json:"descricao"`
	Unidades  string `json:"unidades"`
}

// blocoAssinaturaRequest executa a requisição de cadastro ou alteração
// de um bloco de assinatura.
func (c *Client) blocoAssinaturaRequest(ctx context.Context, params BlocoAssinaturaParams, endpoint string) (*BlocoAssinaturaResult, error) {
	if strings.TrimSpace(params.Descricao) == "" {
		return nil, fmt.Errorf("descricao required")
	}
	if len(params.Unidades) > 0 {
		for _, b := range params.Unidades {
			if strings.TrimSpace(b) == "" {
				return nil, fmt.Errorf("unidade vazia na lista de unidades")
			}
		}
	}

	jsonBody, err := json.Marshal(blocoAssinaturaParams{
		Descricao: params.Descricao,
		Unidades:  strings.Join(params.Unidades, ","),
	})
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

	var env Envelope[BlocoAssinaturaResult]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}

// CadastrarBlocoAssinatura cadastra um novo bloco de assinatura e retorna
// os dados do bloco criado.
func (c *Client) CadastrarBlocoAssinatura(ctx context.Context, params BlocoAssinaturaParams) (*BlocoAssinaturaResult, error) {
	endpoint := c.endpoint + "/bloco/assinatura/criar"
	return c.blocoAssinaturaRequest(ctx, params, endpoint)
}

// AlterarBlocoAssinatura altera um bloco de assinatura existente e retorna
// os dados atualizados do bloco.
func (c *Client) AlterarBlocoAssinatura(ctx context.Context, params BlocoAssinaturaParams, bloco int) (*BlocoAssinaturaResult, error) {
	if bloco <= 0 {
		return nil, fmt.Errorf("bloco invalido: %d", bloco)
	}
	endpoint := fmt.Sprintf("%s/bloco/assinatura/%d/alterar", c.endpoint, bloco)
	return c.blocoAssinaturaRequest(ctx, params, endpoint)
}

// BlocoActionParams reúne os parâmetros necessários para realizar
// ações sobre blocos que recebem uma lista de identificadores.
type BlocoActionParams struct {
	Blocos []string
}

// blocoActionParams representa o formato esperado pelo WSSEI
// para o corpo da requisição.
type blocoActionParams struct {
	Blocos string `json:"blocos"`
}

// blocoActionRequest executa uma ação sobre blocos enviando a lista de
// identificadores informada no formato esperado pelo WSSEI.
func (c *Client) blocoActionRequest(ctx context.Context, params BlocoActionParams, endpoint string) error {
	if len(params.Blocos) == 0 {
		return fmt.Errorf("blocos required")
	}
	for _, b := range params.Blocos {
		if strings.TrimSpace(b) == "" {
			return fmt.Errorf("bloco vazio na lista de blocos")
		}
	}

	jsonBody, err := json.Marshal(blocoActionParams{
		Blocos: strings.Join(params.Blocos, ","),
	})
	if err != nil {
		return fmt.Errorf("json body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[struct{}]
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return nil
}

// ExcluirBlocoAssinatura exclui um ou mais blocos de assinatura.
func (c *Client) ExcluirBlocoAssinatura(ctx context.Context, params BlocoActionParams) error {
	endpoint := c.endpoint + "/bloco/assinatura/excluir"
	return c.blocoActionRequest(ctx, params, endpoint)
}

// ExcluirBlocoInterno exclui um ou mais blocos internos.
func (c *Client) ExcluirBlocoInterno(ctx context.Context, params BlocoActionParams) error {
	endpoint := c.endpoint + "/bloco/interno/excluir"
	return c.blocoActionRequest(ctx, params, endpoint)
}

// ConcluirBlocoAssinatura conclui um ou mais blocos de assinatura.
func (c *Client) ConcluirBlocoAssinatura(ctx context.Context, params BlocoActionParams) error {
	endpoint := c.endpoint + "/bloco/assinatura/concluir"
	return c.blocoActionRequest(ctx, params, endpoint)
}

// ConcluirBlocoInterno conclui um ou mais blocos de assinatura.
func (c *Client) ConcluirBlocoInterno(ctx context.Context, params BlocoActionParams) error {
	endpoint := c.endpoint + "/bloco/interno/concluir"
	return c.blocoActionRequest(ctx, params, endpoint)
}