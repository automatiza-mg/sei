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

// ListarGrupoAcompanhamentoResult representa os dados retornados na listagem de grupos de acompanhamento.
type ListarGrupoAcompanhamentoResult struct {
	IDGrupoAcompanhamento string `json:"idGrupoAcompanhamento"`
	Nome                  string `json:"nome"`
}

// ListarGrupoAcompanhamentoParams reúne os parâmetros opcionais da listagem de grupos de acompanhamento.
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type ListarGrupoAcompanhamentoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o id do grupo de acompanhamento para detalhamento.
	ID int
}

// values converte os parâmetros da pesquisa em query params,
// omitindo campos que possuem valor zero.
func (p ListarGrupoAcompanhamentoParams) values() url.Values {
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

// ListarGrupoAcompanhamento retorna os grupos de acompanhamento cadastrados.
func (c *Client) ListarGrupoAcompanhamento(ctx context.Context, params ListarGrupoAcompanhamentoParams) ([]ListarGrupoAcompanhamentoResult, int, error) {
	endpoint := c.endpoint + "/grupoacompanhamento/listar"
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

	var env Envelope[[]ListarGrupoAcompanhamentoResult]
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

// CadastrarGrupoAcompanhamentoResult representa os dados retornados após o cadastro de um grupo de acompanhamento.
type CadastrarGrupoAcompanhamentoResult struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

// CadastrarGrupoAcompanhamentoParams reúne os dados necessários para cadastrar um grupo de acompanhamento.
type CadastrarGrupoAcompanhamentoParams struct {
	Nome string `json:"nome"`
}

// CadastrarGrupoAcompanhamento cria um novo grupo de acompanhamento.
func (c *Client) CadastrarGrupoAcompanhamento(ctx context.Context, params CadastrarGrupoAcompanhamentoParams) (*CadastrarGrupoAcompanhamentoResult, error) {
	if strings.TrimSpace(params.Nome) == "" {
		return nil, fmt.Errorf("nome required")
	}
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json body: %w", err)
	}

	endpoint := c.endpoint + "/grupoacompanhamento/cadastrar"

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

	var env Envelope[CadastrarGrupoAcompanhamentoResult]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}

// ExcluirGrupoAcompanhamentoParams reúne os parâmetros necessários para
// excluir um ou mais grupos de acompanhamento.
type ExcluirGrupoAcompanhamentoParams struct {
	Grupos []string
}

// excluirGrupoAcompanhamentoParams representa o formato esperado pelo WSSEI
// para o corpo da requisição.
type excluirGrupoAcompanhamentoParams struct {
	Grupos string `json:"grupos"`
}

// ExcluirGrupoAcompanhamento exclui um ou mais grupos de acompanhamento.
func (c *Client) ExcluirGrupoAcompanhamento(ctx context.Context, params ExcluirGrupoAcompanhamentoParams) error {
	if len(params.Grupos) == 0 {
		return fmt.Errorf("grupos required")
	}
	for _, grupo := range params.Grupos {
		if strings.TrimSpace(grupo) == "" {
			return fmt.Errorf("grupo vazio na lista de grupos")
		}
	}

	jsonBody, err := json.Marshal(excluirGrupoAcompanhamentoParams{
		Grupos: strings.Join(params.Grupos, ","),
	})
	if err != nil {
		return fmt.Errorf("json body: %w", err)
	}

	endpoint := c.endpoint + "/grupoacompanhamento/excluir"

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

// AlterarGrupoAcompanhamentoParams reúne os parâmetros necessários para
// alterar um ou mais grupos de acompanhamento.
type AlterarGrupoAcompanhamentoParams struct {
	Nome []string
}

// alterarGrupoAcompanhamentoParams representa o payload enviado ao WSSEI
// para alteração de grupo de acompanhamento.
type alterarGrupoAcompanhamentoParams struct {
	Nome string `json:"nome"`
}

// AlterarGrupoAcompanhamento altera um grupo de acompanhamento.
func (c *Client) AlterarGrupoAcompanhamento(ctx context.Context, grupoacompanhamento int, params AlterarGrupoAcompanhamentoParams) error {
	if grupoacompanhamento <= 0 {
		return fmt.Errorf("grupo acompanhamento invalido: %d", grupoacompanhamento)
	}
	if len(params.Nome) == 0 {
		return fmt.Errorf("nome required")
	}
	for _, nome := range params.Nome {
		if strings.TrimSpace(nome) == "" {
			return fmt.Errorf("nome vazio na lista de nomes")
		}
	}

	jsonBody, err := json.Marshal(alterarGrupoAcompanhamentoParams{
		Nome: strings.Join(params.Nome, ","),
	})
	if err != nil {
		return fmt.Errorf("json body: %w", err)
	}

	endpoint := fmt.Sprintf("%s/grupoacompanhamento/%d/alterar", c.endpoint, grupoacompanhamento)

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
