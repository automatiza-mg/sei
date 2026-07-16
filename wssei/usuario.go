package wssei

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ListarUsuariosParams seleciona a query do ListasUsuarios
type ListarUsuariosParams struct {
	Limit   int
	Start   int
	Unidade int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListarUsuariosParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.Unidade != 0 {
		q.Set("unidade", strconv.Itoa(p.Unidade))
	}
	return q
}

// ListarUsuarios retorna a lista de Usuários.
func (c *Client) ListarUsuarios(ctx context.Context, params ListarUsuariosParams) ([]Usuarios, int, error) {
	url := fmt.Sprintf("%s/usuario/listar", c.endpoint)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("erro request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("erro response: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]Usuarios]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("erro json decoder: %w", err)
	}

	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("erro listar ususarios : %s", result.Mensagem)
	}

	total, err := result.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("total invalido")
	}

	return result.Data, total, nil

}

// Usuarios tipo utilizado na funcao "ListarUsuarios".
type Usuarios struct {
	IDUsuario string `json:"id_usuario"`
	Sigla     string `json:"sigla"`
	Nome      string `json:"nome"`
	IDContato string `json:"id_contato"`
	Total     string `json:"total"`
}

// ListarUsuariosParams seleciona a query do ListasUsuarios.
type PesquisarUsuariosParams struct {
	// PalavraChave é obrigatório
	palavraChave string
	orgao        int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p PesquisarUsuariosParams) values() url.Values {
	q := make(url.Values)
	if p.palavraChave != "" {
		q.Set("limit", p.palavraChave)
	}
	if p.orgao != 0 {
		q.Set("start", strconv.Itoa(p.orgao))
	}
	return q
}

// PesquiarUsuarios retorna a pesquisa de Usuários.
func (c *Client) PesquiarUsuarios(ctx context.Context, params PesquisarUsuariosParams) (*UsuariosPesquisa, error) {
	if params.palavraChave == "" {
		return nil, fmt.Errorf("palavraChave required")
	}
	url := fmt.Sprintf("%s/usuario/pesquisar", c.endpoint)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro response: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[UsuariosPesquisa]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("erro json decoder: %w", err)
	}

	if result.Sucesso != true {
		return nil, fmt.Errorf("erro pesquisar ususarios : %s", result.Mensagem)
	}

	return &result.Data, nil

}

// UsuariosPesquisa tipo utilizado na funcao "PesquiarUsuarios".
type UsuariosPesquisa struct {
	IDContato string `json:"id_contato"`
	IDUsuario string `json:"id_usuario"`
	Sigla     string `json:"sigla"`
	Nome      string `json:"nome"`
}

// RetornarUnidadesUsuarios retorna as unidades de um Usuário.
func (c *Client) RetornarUnidadesUsuarios(ctx context.Context, usuario int) ([]UnidadesUsuarios, error) {
	if usuario == 0 {
		return nil, fmt.Errorf("usuario required")
	}
	url := fmt.Sprintf("%s/usuario/listar", c.endpoint)
	query := strconv.Itoa(usuario)
	url += "?" + query

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro response: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]UnidadesUsuarios]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("erro json decoder: %w", err)
	}

	if result.Sucesso != true {
		return nil, fmt.Errorf("erro listar ususarios : %s", result.Mensagem)
	}

	return result.Data, nil

}

// UnidadesUsuarios tipo utilizado na funcao RetornarUnidadesUsuarios.
type UnidadesUsuarios struct {
	ID    string `json:"id"`
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
}
