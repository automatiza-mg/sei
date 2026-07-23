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

// ListarGrupoModeloDocumentoResult representa os dados retornados na listagem
// de grupos de modelos de documento.
type ListarGrupoModeloDocumentoResult struct {
	IDGrupoProtocoloModelo string `json:"idGrupoProtocoloModelo"`
	Nome                   string `json:"nome"`
}

// ListarGrupoModeloDocumentoParams reúne os parâmetros opcionais da listagem
// de grupos de modelos de documento.
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type ListarGrupoModeloDocumentoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o identificador do grupo de modelo de documento para detalhamento.
	ID int
}

// values converte os parâmetros da pesquisa em query params,
// omitindo campos que possuem valor zero.
func (p ListarGrupoModeloDocumentoParams) values() url.Values {
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

// ListarGrupoModeloDocumento retorna os grupos de modelos de documento
// cadastrados.
func (c *Client) ListarGrupoModeloDocumento(ctx context.Context, params ListarGrupoModeloDocumentoParams) ([]ListarGrupoModeloDocumentoResult, int, error) {
	endpoint := c.endpoint + "/protocolomodelo/grupo/listar"
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

	var env Envelope[[]ListarGrupoModeloDocumentoResult]
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

// ListarModeloDocumentoResult representa os dados retornados na listagem de
// modelos de documento.
type ListarModeloDocumentoResult struct {
	IDProtocoloModelo        string `json:"idProtocoloModelo"`
	IDGrupoProtocoloModelo   string `json:"idGrupoProtocoloModelo"`
	NomeGrupoProtocoloModelo string `json:"nomeGrupoProtocoloModelo"`
	IDUsuario                string `json:"idUsuario"`
	NomeUsuario              string `json:"nomeUsuario"`
	SiglaUsuario             string `json:"siglaUsuario"`
	ProtocoloFormatado       string `json:"protocoloFormatado"`
	NomeSerie                string `json:"nomeSerie"`
	DataGeracao              string `json:"dataGeracao"`
}

// ListarModeloDocumentoParams reúne os parâmetros opcionais da listagem de
// modelos de documento.
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type ListarModeloDocumentoParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// Filter é a palavra-chave da pesquisa.
	Filter string
	// ID é o identificador do modelo de documento para detalhamento.
	ID int
	// GrupoProtocoloModelo é o identificador do grupo de modelo de documento.
	GrupoProtocoloModelo int
	// TipoFiltro é a sigla do tipo de filtro: T para todos ou M para meus.
	TipoFiltro TipoFiltro
}

type TipoFiltro string

const (
	Todos TipoFiltro = "T"
	Meus  TipoFiltro = "M"
)

// values converte os parâmetros da pesquisa em query params,
// omitindo campos que possuem valor zero.
func (p ListarModeloDocumentoParams) values() url.Values {
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
	if p.GrupoProtocoloModelo != 0 {
		q.Set("grupoProtocoloModelo", strconv.Itoa(p.GrupoProtocoloModelo))
	}
	if p.TipoFiltro != "" {
		q.Set("tipoFiltro", string(p.TipoFiltro))
	}
	return q
}

// ListarModeloDocumento retorna os modelos de documento cadastrados.
func (c *Client) ListarModeloDocumento(ctx context.Context, params ListarModeloDocumentoParams) ([]ListarModeloDocumentoResult, int, error) {
	endpoint := c.endpoint + "/protocolomodelo/listar"
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

	var env Envelope[[]ListarModeloDocumentoResult]
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
