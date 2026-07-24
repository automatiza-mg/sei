package wssei

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Upload representa os parâmetros de upload configurados no
// sistema.
type Upload struct {
	Extensoes         []UploadExtensao `json:"extensoes"`
	TamanhoDocDefault string           `json:"tamanhoDocDefault"`
	ValidarExtensoes  bool             `json:"validarExtensoes"`
	Info              string           `json:"info"`
}

// UploadExtensao representa uma extensão de arquivo permitida para upload e
// seu tamanho máximo.
type UploadExtensao struct {
	Extensao string `json:"extensao"`
	Tamanho  string `json:"tamanho"`
}

// ConsultarParametrosUpload retorna os parâmetros de upload configurados no
// sistema.
func (c *Client) ConsultarParametrosUpload(ctx context.Context) (*Upload, error) {
	endpoint := c.endpoint + "/upload/parametros"

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

	var env Envelope[Upload]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return &env.Data, nil
}
