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

// ListarAtividadesParams reúne os parâmetros de [Client.ListarAtividades].
//
// Campos com valor zero (0) são omitidos da requisição, exceto
// [ListarAtividadesParams.Procedimento], que é obrigatório.
type ListarAtividadesParams struct {
	// Procedimento é o id do processo. Obrigatório.
	Procedimento int
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados
func (p ListarAtividadesParams) values() url.Values {
	q := make(url.Values)
	q.Set("procedimento", strconv.Itoa(p.Procedimento))
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}

	return q
}

// ListarAtividades retorna a lista de Atividades do Processo.
func (c *Client) ListarAtividades(ctx context.Context, params ListarAtividadesParams) ([]Atividade, int, error) {
	endpoint := c.endpoint + "/atividade/listar"
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

	var env Envelope[[]Atividade]
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

// Atividade é um andamento registrado pelo SEI para um processo.
type Atividade struct {
	ID        string             `json:"id"`
	Atributos AtributosAtividade `json:"atributos"`
}

// AtributosAtividade reúne os campos aninhados sob `atributos` em cada
type AtributosAtividade struct {
	IDProcesso string `json:"idProcesso"`
	Usuario    string `json:"usuario"`
	Data       string `json:"data"`
	Hora       string `json:"hora"`
	Unidade    string `json:"unidade"`
	Informacao string `json:"informacao"`
}
