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

// Marcador representa o marcador do processo retornada pelo WSSEI.
type Marcador struct {
	IDMarcador   string `json:"idMarcador"`
	IDProtocolo  string `json:"idProtocolo"`
	Texto        string `json:"texto"`
	IDCor        string `json:"idCor"`
	DescricaoCor string `json:"descricaoCor"`
	ArquivoCor   string `json:"arquivoCor"`
}

// MarcadorCor representa a cor de marcador retornada pelo WSSEI.
type MarcadorCor struct {
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
	Arquivo   string `json:"arquivo"`
}

// MarcadorHistorico representa o histórico de marcador do processo pelo WSSEI.
type MarcadorHistorico struct {
	MarcadorAtivo string `json:"marcadorAtivo"`
	Data          string `json:"data"`
	Texto         string `json:"texto"`
	NomeMarcador  string `json:"nomeMarcador"`
	NomeUsuario   string `json:"nomeUsuario"`
	SiglaUsuario  string `json:"siglaUsuario"`
}

// MarcadorProcessoParams reúne os parâmetros para vincular um marcador a um processo.
type MarcadorProcessoParams struct {
	Texto    string `json:"texto"`
	Marcador int    `json:"marcador"`
}

// ConsultarMarcador retorna o marcador associado ao processo identificado por protocolo.
func (c *Client) ConsultarMarcador(ctx context.Context, protocolo int) (*Marcador, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("protocolo inválido: %d", protocolo)
	}

	endpoint := fmt.Sprintf("%s/marcador/processo/%d/consultar", c.endpoint, protocolo)

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

	var env Envelope[Marcador]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}

// ListarCores retorna a lista de [MarcadorCor] disponíveis e o total de registros.
func (c *Client) ListarCores(ctx context.Context) ([]MarcadorCor, int, error) {
	endpoint := fmt.Sprintf("%s/marcador/cores/listar", c.endpoint)

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

	var env Envelope[[]MarcadorCor]
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

// ListarHistoricoMarcador retorna o histórico de marcadores do processo identificado por protocolo.
func (c *Client) ListarHistoricoMarcador(ctx context.Context, protocolo int) ([]MarcadorHistorico, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("protocolo inválido: %d", protocolo)
	}

	endpoint := fmt.Sprintf("%s/marcador/processo/%d/historico/listar", c.endpoint, protocolo)

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

	var env Envelope[[]MarcadorHistorico]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return env.Data, nil
}

// MarcarProcesso vincula um marcador ao processo identificado por protocolo.
func (c *Client) MarcarProcesso(ctx context.Context, protocolo int, params MarcadorProcessoParams) error {
	if protocolo <= 0 {
		return fmt.Errorf("protocolo inválido: %d", protocolo)
	}
	if strings.TrimSpace(params.Texto) == "" {
		return fmt.Errorf("texto required")
	}
	if params.Marcador <= 0 {
		return fmt.Errorf("marcador inválido: %d", params.Marcador)
	}
	endpoint := fmt.Sprintf("%s/marcador/processo/%d/marcar", c.endpoint, protocolo)

	jsonBody, err := json.Marshal(params)
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

// MarcadorDetalhe representa a pesquisa de marcador do processo pelo WSSEI.
type MarcadorDetalhe struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	Ativo        string `json:"ativo"`
	IDCor        string `json:"idCor"`
	DescricaoCor string `json:"descricaoCor"`
	ArquivoCor   string `json:"arquivoCor"`
}

// TipoAtivo representa o status de ativo usado nas requisições.
type TipoAtivo string

// Tipos de status aceitos para o atributo ativo.
const (
	TipoAtivoSim TipoAtivo = "S"
	TipoAtivoNao TipoAtivo = "N"
)

// PesquisarMarcadorParams reúne os parâmetros opcionais de [Client.PesquisarMarcador].
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type PesquisarMarcadorParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o id do marcador para detalhamento.
	ID int
	// Ativo representa o status do marcador.
	Ativo TipoAtivo
}

// Converte os parâmetros em [url.Values], omitindo os campos opcionais
// zerados.
func (p PesquisarMarcadorParams) values() url.Values {
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
	if p.Ativo != "" {
		q.Set("ativo", string(p.Ativo))
	}

	return q
}

// PesquisarMarcador retorna a lista de marcadores e o total de registros.
func (c *Client) PesquisarMarcador(ctx context.Context, params PesquisarMarcadorParams) ([]MarcadorDetalhe, int, error) {
	endpoint := c.endpoint + "/marcador/pesquisar"
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

	var env Envelope[[]MarcadorDetalhe]
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

// CadastrarMarcadorParams reúne os parâmetros necessários
// para cadastrar um marcador no WSSEI.
type CadastrarMarcadorParams struct {
	Nome  string `json:"nome"`
	IDCor int    `json:"idCor"`
}

// marcadorDetalheRequest executa a requisição de cadastro ou alteração
// de um marcador.
func (c *Client) marcadorDetalheRequest(ctx context.Context, endpoint string, params CadastrarMarcadorParams) (*MarcadorDetalhe, error) {
	if strings.TrimSpace(params.Nome) == "" {
		return nil, fmt.Errorf("nome invalido: %s", params.Nome)
	}
	if params.IDCor <= 0 {
		return nil, fmt.Errorf("idCor invalido: %d", params.IDCor)
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

	var env Envelope[MarcadorDetalhe]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}

// CadastrarMarcador cadastra um novo marcador no WSSEI.
func (c *Client) CadastrarMarcador(ctx context.Context, params CadastrarMarcadorParams) (*MarcadorDetalhe, error) {
	endpoint := c.endpoint + "/marcador/criar"
	return c.marcadorDetalheRequest(ctx, endpoint, params)
}

// AlterarMarcador altera um marcador no WSSEI.
func (c *Client) AlterarMarcador(ctx context.Context, marcador int, params CadastrarMarcadorParams) (*MarcadorDetalhe, error) {
	if marcador <= 0 {
		return nil, fmt.Errorf("marcador invalido: %d", marcador)
	}

	endpoint := fmt.Sprintf("%s/marcador/%d/alterar", c.endpoint, marcador)
	return c.marcadorDetalheRequest(ctx, endpoint, params)
}

// MarcadorRequestParams reúne os parâmetros necessários para
// alterar um ou mais marcadores no WSSEI.
type MarcadorRequestParams struct {
	Marcadores []string
}

// marcadorRequestParams representa o formato esperado pelo WSSEI
// para o corpo da requisição.
type marcadorRequestParams struct {
	Marcadores string `json:"marcadores"`
}

// marcadorRequest executa a requisição de alteração de status
// de um ou mais marcadores no WSSEI.
func (c *Client) marcadorRequest(ctx context.Context, endpoint string, params MarcadorRequestParams) error {
	if len(params.Marcadores) == 0 {
		return fmt.Errorf("lista de marcadores obrigatória")
	}

	for _, marcador := range params.Marcadores {
		if strings.TrimSpace(marcador) == "" {
			return fmt.Errorf("marcador vazio na lista de marcadores")
		}
	}

	jsonBody, err := json.Marshal(marcadorRequestParams{
		Marcadores: strings.Join(params.Marcadores, ","),
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

// ExcluirMarcadores exclui um ou mais marcadores no WSSEI.
func (c *Client) ExcluirMarcadores(ctx context.Context, params MarcadorRequestParams) error {
	endpoint := c.endpoint + "/marcador/excluir"
	return c.marcadorRequest(ctx, endpoint, params)
}

// DesativarMarcadores desativa um ou mais marcadores no WSSEI.
func (c *Client) DesativarMarcadores(ctx context.Context, params MarcadorRequestParams) error {
	endpoint := c.endpoint + "/marcador/desativar"
	return c.marcadorRequest(ctx, endpoint, params)
}

// ReativarMarcadores reativa um ou mais marcadores no WSSEI.
func (c *Client) ReativarMarcadores(ctx context.Context, params MarcadorRequestParams) error {
	endpoint := c.endpoint + "/marcador/reativar"
	return c.marcadorRequest(ctx, endpoint, params)
}
