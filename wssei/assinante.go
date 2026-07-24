package wssei

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Assinante representa um assinante retornado pelo WSSEI.
type Assinante struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

// ListarAssinantesParams reúne os parâmetros opcionais da listagem de
// assinantes.
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type ListarAssinantesParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o identificador do assinante para detalhamento.
	ID int
}

// values converte os parâmetros da pesquisa em query params,
// omitindo campos que possuem valor zero.
func (p ListarAssinantesParams) values() url.Values {
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

// ListarAssinantes retorna os assinantes cadastrados na unidade e o total de
// registros.
func (c *Client) ListarAssinantes(ctx context.Context, params ListarAssinantesParams) ([]Assinante, int, error) {
	endpoint := c.endpoint + "/assinante/listar"
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

	var env Envelope[[]Assinante]
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

// ListarOrgaosAssinanteParams reúne os parâmetros opcionais da listagem de
// órgãos para assinantes.
//
// Campos com valor zero (0) são omitidos da requisição.
type ListarOrgaosAssinanteParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
}

// values converte os parâmetros da pesquisa em query params,
// omitindo campos que possuem valor zero.
func (p ListarOrgaosAssinanteParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	return q
}

// ListarOrgaosAssinante retorna os órgãos disponíveis para assinantes e o
// total de registros.
func (c *Client) ListarOrgaosAssinante(ctx context.Context, params ListarOrgaosAssinanteParams) ([]Orgao, int, error) {
	endpoint := c.endpoint + "/assinante/orgao"
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

	var env Envelope[[]Orgao]
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
