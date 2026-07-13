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

// HipoteseLegalResult representa uma hipótese legal retornada pela pesquisa pelo WSSEI
type HipoteseLegalResult struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	BaseLegal string `json:"baselegal"`
}

// HipoteseLegalParams reúne os parâmetros de [Client.PesquisarHipoteseLegal].
//
// Campos com valor zero (0, "" ou false) são omitidos da requisição.
type HipoteseLegalParams struct {
	// Limit é o limite de registros da paginação.
	Limit       int
	// Start é a página de início da paginação.
	Start       int
	// Filter é a palavra-chave da pesquisa.
	Filter      string
	// ID é o id da hipótese legal para detalhamento.
	ID          int
	// NivelAcesso é o nível de acesso da hipótese legal (obrigatório).
	NivelAcesso NivelAcesso
}

// Converte os parâmetros em [url.Values], omitindo os campos opcionais
// zerados. NivelAcesso é sempre incluído por ser obrigatório.
func (p HipoteseLegalParams) values() url.Values {
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
	q.Set("nivelAcesso", strconv.Itoa(int(p.NivelAcesso)))

	return q
}

// PesquisarHipoteseLegal retorna a lista de hipóteses legais e o total de
// registros.
func (c *Client) PesquisarHipoteseLegal(ctx context.Context, params HipoteseLegalParams) ([]HipoteseLegalResult, int, error) {
	endpoint := c.endpoint + "/hipoteseLegal/pesquisar"
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

	var env Envelope[[]HipoteseLegalResult]
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
