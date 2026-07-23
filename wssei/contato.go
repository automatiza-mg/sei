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

// PesquisarContatoResult representa os dados retornados na pesquisa de
// contatos.
type PesquisarContatoResult struct {
	NomeFormatado string `json:"nomeformatado"`
	Nome          string `json:"nome"`
	Sigla         string `json:"sigla"`
	ID            string `json:"id"`
}

// PesquisarContatoParams reúne os parâmetros opcionais da pesquisa de
// contatos.
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type PesquisarContatoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o identificador do contato para detalhamento.
	ID int
	// IDGrupoContato é o identificador do grupo de contato.
	IDGrupoContato int
}

// values converte os parâmetros da pesquisa em query params,
// omitindo campos que possuem valor zero.
func (p PesquisarContatoParams) values() url.Values {
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
	if p.IDGrupoContato != 0 {
		q.Set("idGrupoContato", strconv.Itoa(p.IDGrupoContato))
	}
	return q
}

// PesquisarContato retorna os contatos encontrados e o total de registros.
func (c *Client) PesquisarContato(ctx context.Context, params PesquisarContatoParams) ([]PesquisarContatoResult, int, error) {
	endpoint := c.endpoint + "/contato/pesquisar"
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

	var env Envelope[[]PesquisarContatoResult]
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

// CadastrarContatoResult representa os dados retornados após o cadastro de um
// contato.
type CadastrarContatoResult struct {
	ID string `json:"id"`
}

// CadastrarContatoParams reúne os dados necessários para cadastrar um
// contato.
type CadastrarContatoParams struct {
	Nome string `json:"nome"`
}

// CadastrarContato cria um novo contato.
func (c *Client) CadastrarContato(ctx context.Context, params CadastrarContatoParams) (*CadastrarContatoResult, error) {
	if strings.TrimSpace(params.Nome) == "" {
		return nil, fmt.Errorf("nome required")
	}
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json body: %w", err)
	}

	endpoint := c.endpoint + "/contato/criar"

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

	var env Envelope[CadastrarContatoResult]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}
