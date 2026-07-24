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

// TipoDocumentoExterno representa os dados retornados na
// pesquisa de tipos de documento externo.
type TipoDocumentoExterno struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

// TipoDocumentoExternoParams reúne os parâmetros opcionais da
// pesquisa de tipos de documento externo.
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type TipoDocumentoExternoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o identificador do tipo de documento para detalhamento.
	ID int
}

// values converte os parâmetros da pesquisa em query params,
// omitindo campos que possuem valor zero.
func (p TipoDocumentoExternoParams) values() url.Values {
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

// PesquisarTipoDocumentoExterno retorna os tipos de documento externo
// cadastrados e o total de registros.
func (c *Client) PesquisarTipoDocumentoExterno(ctx context.Context, params TipoDocumentoExternoParams) ([]TipoDocumentoExterno, int, error) {
	endpoint := c.endpoint + "/serie/externo/pesquisar"
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

	var env Envelope[[]TipoDocumentoExterno]
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
