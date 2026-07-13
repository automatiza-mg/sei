// Package wssei fornece um client Go para o módulo WSSEI do SEI (Sistema
// Eletrônico de Informações), cobrindo autenticação, processos, documentos,
// blocos de assinatura, marcadores e demais operações expostas pela API.
//
// O [Client] autentica com usuário e senha uma única vez e reaproveita o
// token gerado em cada requisição, renovando-o automaticamente quando o
// servidor responde 401/403. Todas as respostas do WSSEI seguem o formato
// [Envelope]; os métodos do [Client] já extraem o conteúdo útil e traduzem
// falhas para erros do Go.
//
// Referência da API: https://pengovbr.github.io/mod-wssei/#/
package wssei

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// AuthCallback é invocado pelo [Client] após cada autenticação bem-sucedida
// no WSSEI, incluindo as renovações automáticas por token expirado. Recebe o
// [*AuthResponse] retornado pelo endpoint POST /autenticar, permitindo
// persistir o token ou observar o resultado do login.
type AuthCallback func(ctx context.Context, resp *AuthResponse) error

// O caminho da API do módulo WSSEI relativo à URL base do SEI.
const apiBasePath = "/sei/modulos/wssei/controlador_ws.php/api/v2"

// Monta a URL base da API do WSSEI a partir da URL base do SEI.
func apiBaseURL(baseURL string) string {
	return strings.TrimRight(baseURL, "/") + apiBasePath
}

// Envelope é o formato padrão de resposta do WSSEI, presente em todas as
// chamadas HTTP.
// O conteúdo útil fica em Data, cujo tipo varia por endpoint.
type Envelope[T any] struct {
	Sucesso  bool   `json:"sucesso"`
	Mensagem string `json:"mensagem"`
	Total    string `json:"total"`
	Data     T      `json:"data"`
}

// Se total vazio, return 0  e sem erro
func (e *Envelope[T]) getTotal() (int, error) {
	if e.Total == "" {
		return 0, nil
	}
	return strconv.Atoi(e.Total)
}

// Config reúne os dados necessários para autenticar e acessar o WSSEI.
type Config struct {
	// BaseURL é a URL base do SEI (ex: https://www.sei.mg.gov.br).
	BaseURL string
	// Usuario é o login do usuário usado na autenticação.
	Usuario string
	// Senha é a senha do usuário usada na autenticação.
	Senha string
	// Orgao é o id do órgão da autenticação.
	Orgao int
	// OnAuthenticated, se não nulo, é chamado pelo [Client] após cada
	// autenticação bem-sucedida no WSSEI. Ver [AuthCallback].
	OnAuthenticated AuthCallback
}

// Client é o client HTTP do WSSEI. Ele encapsula a URL base da API, o
// [http.Client] usado nas requisições e a autenticação por token, feita de
// forma transparente pelo [tokenTransport].
type Client struct {
	endpoint string
	http     *http.Client
}

// NewClient cria um [Client] que autentica no WSSEI com usuário e senha,
// gerando e reaproveitando o token automaticamente em cada requisição.
//
// O token é cacheado em memória e renovado quando o servidor responde
// 401/403. Caso [Config.OnAuthenticated] esteja definido, o callback é
// invocado após cada autenticação bem-sucedida (inclusive nas renovações).
func NewClient(cfg Config) *Client {
	return &Client{
		endpoint: apiBaseURL(cfg.BaseURL),
		http: &http.Client{
			Transport: &tokenTransport{
				RoundTripper:    http.DefaultTransport,
				auth:            NewAuth(cfg.BaseURL),
				usuario:         cfg.Usuario,
				senha:           cfg.Senha,
				orgao:           cfg.Orgao,
				onAuthenticated: cfg.OnAuthenticated,
			},
		},
	}
}
