package wssei

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ConsultarDocumentoInterno retorna os metadados do Documento Interno.
func (c *Client) ConsultarDocumentoInterno(ctx context.Context, protocolo int) (*DocumentoInterno, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}

	url := fmt.Sprintf(
		"%s/documento/interno/consultar/%d",
		c.endpoint,
		protocolo,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[DocumentoInterno]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, fmt.Errorf("consultar failed %d: %s", protocolo, result.Mensagem)
	}

	return &result.Data, nil

}

// DocumentoInterno tipo utilizado na funcao "ConsultarDocumentoInterno".
type DocumentoInterno struct {
	NomeDocumento            string `json:"nomeDocumento"`
	Protocolo                string `json:"protocolo"`
	IDDocumento              string `json:"idDocumento"`
	IDSerie                  string `json:"idSerie"`
	NomeSerie                string `json:"nomeSerie"`
	Numero                   string `json:"numero"`
	IDTipoConferencia        string `json:"idTipoConferencia"`
	DescricaoTipoConferencia string `json:"descricaoTipoConferencia"`
	NivelAcesso              string `json:"nivelAcesso"`
	IDHipoteseLegal          string `json:"idHipoteseLegal"`
	NomeHipoteseLegal        string `json:"nomeHipoteseLegal"`
	BaseLegal                string `json:"baseLegal"`
	GrauSigilo               string `json:"grauSigilo"`
	Descricao                string `json:"descricao"`
	DataElaboracao           string `json:"dataElaboracao"`
	Observacao               string `json:"observacao"`

	Assuntos     []Assunto     `json:"assuntos"`
	Interessados []Interessado `json:"interessados"`
	//Destinatarios documentado como string, mas é identico ao Interessados
	Destinatarios       []Interessado `json:"destinatarios"`
	ObservacoesUnidades Slice[string] `json:"observacoesUnidades"`
}

// VisualizarDocumento retorna o HTML do Documento para visualização.
func (c *Client) VisualizarDocumento(ctx context.Context, documento int) (string, error) {
	if documento <= 0 {
		return "", fmt.Errorf("invalid documento: %d", documento)
	}

	url := fmt.Sprintf(
		"%s/documento/%d/interno/visualizar",
		c.endpoint,
		documento,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[string]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}

	if result.Sucesso != true {
		return "", fmt.Errorf("visualizar failed %d: %s", documento, result.Mensagem)
	}

	return result.Data, nil
}

// BaixarAnexo baixa um documento externo. Retorna o corpo da requisicao e o Content-Type.
func (c *Client) BaixarAnexo(ctx context.Context, protocolo int) (io.ReadCloser, string, error) {
	if protocolo <= 0 {
		return nil, "", fmt.Errorf("invalid protocolo: %d", protocolo)
	}

	url := fmt.Sprintf(
		"%s/documento/baixar/anexo/%d",
		c.endpoint,
		protocolo,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("response error: %w", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	return resp.Body, contentType, nil
}

// PesquisarTipoTemplateDocumento retorna a lista de Templates do Documento.
func (c *Client) PesquisarTipoTemplateDocumento(ctx context.Context, id int, procedimento int) (*TemplateDocumento, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %d", id)
	}
	if procedimento <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", procedimento)
	}

	q := url.Values{}
	q.Set("id", strconv.Itoa(id))
	q.Set("procedimento", strconv.Itoa(procedimento))

	url := fmt.Sprintf("%s/documento/tipo/template?%s", c.endpoint, q.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[TemplateDocumento]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if result.Sucesso != true {
		return nil, fmt.Errorf("pesquisar failed %d: %s", procedimento, result.Mensagem)
	}

	return &result.Data, nil
}

// TemplateDocumento tipo utilizado na funcao "PesquisarTipoTemplateDocumento".
type TemplateDocumento struct {
	Assuntos                     Assuntos              `json:"assuntos"`
	Interessados                 string                `json:"interessados"`
	NivelAcessoPermitido         NivelAcessoPermitido1 `json:"nivelAcessoPermitido"`
	PermiteInteressados          bool                  `json:"permiteInteressados"`
	PermiteDestinatarios         bool                  `json:"permiteDestinatarios"`
	ObrigatoriedadeHipoteseLegal string                `json:"obrigatoriedadeHipoteseLegal"`
	ObrigatoriedadeGrauSigilo    string                `json:"obrigatoriedadeGrauSigilo"`
}

// ConsultarDocumentoExterno consulta o Documento Externo.
func (c *Client) ConsultarDocumentoExterno(ctx context.Context, protocolo int) (*DocumentoExterno, error) {
	if protocolo <= 0 {
		return nil, fmt.Errorf("invalid protocolo: %d", protocolo)
	}

	url := fmt.Sprintf(
		"%s/documento/externo/consultar/%d",
		c.endpoint,
		protocolo,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[DocumentoExterno]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if result.Sucesso != true {
		return nil, fmt.Errorf("consultar failed %d: %s", protocolo, result.Mensagem)
	}

	return &result.Data, nil
}

// DocumentoExterno tipo utilizado na funcao "ConsultarDocumentoExterno".
type DocumentoExterno struct {
	NomeDocumento            string  `json:"nomeDocumento"`
	Protocolo                string  `json:"protocolo"`
	IDocumento               string  `json:"idDocumento"`
	IdSerie                  string  `json:"idSerie"`
	NomeSerie                string  `json:"nomeSerie"`
	Numero                   string  `json:"numero"`
	IdTipoConferencia        string  `json:"idTipoConferencia"`
	DescricaoTipoConferencia string  `json:"descricaoTipoConferencia"`
	NivelAcesso              string  `json:"nivelAcesso"`
	IdHipoteseLegal          string  `json:"idHipoteseLegal"`
	NomeHipoteseLegal        string  `json:"nomeHipoteseLegal"`
	GrauSigilo               string  `json:"grauSigilo"`
	Descricao                string  `json:"descricao"`
	DataElaboracao           string  `json:"dataElaboracao"`
	Observacao               string  `json:"observacao"`
	Assuntos                 string  `json:"assuntos"`
	Remetente                string  `json:"remetente"`
	Interessados             string  `json:"interessados"`
	Destinatarios            string  `json:"destinatarios"`
	ObservacoesUnidades      string  `json:"observacoesUnidades"`
	Anexo                    []Anexo `json:"anexo"`
}

// ListarDocumentosParams reúne os parâmetros opcionais de [Client.ListarProcessos].
//
// Campos com valor zero (0, "" ou false) são omitidos da requisição.
type ListarDocumentosParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	//Procedimento é o ID do processo. OBRIGATORIO
	Procedimento int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListarDocumentosParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	return q
}

// ListarDocumentosProcesso retorna a lista de Documentos do Processo.
func (c *Client) ListarDocumentosProcesso(ctx context.Context, params ListarDocumentosParams) ([]Documento, int, error) {
	if params.Procedimento == 0 {
		return nil, 0, fmt.Errorf("procedimento required")
	}

	endpoint := fmt.Sprintf("%s/documento/listar/%d", c.endpoint, params.Procedimento)
	if q := params.values().Encode(); q != "" {
		endpoint += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]Documento]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("decode error: %w", err)
	}

	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("listar failed %d: %s", params.Procedimento, result.Mensagem)
	}

	total, err := result.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("invalid total")
	}

	return result.Data, total, nil

}

// Documento tipo utilizado na funcao "ListarDocumentosProcesso".
type Documento struct {
	ID        string             `json:"id"`
	Atributos AtributosDocumento `json:"atributos"`
}

// AssinarDocParams tipo utilizado na funcao AssinarDocumentos.
type AssinarDocParams struct {
	// Todos os parametros sao obrigatorios
	Documento int    `json:"documento"`
	Orgao     int    `json:"orgao"`
	Cargo     string `json:"cargo"`
	Login     string `json:"login"`
	Senha     string `json:"senha"`
	Usuario   int    `json:"usuario"`
}

// AssinarDocumento realiza a assinatura de um documento.
func (c *Client) AssinarDocumento(ctx context.Context, params AssinarDocParams) error {
	if params.Documento <= 0 {
		return fmt.Errorf("invalid documento: %d", params.Documento)
	}
	if params.Orgao <= 0 {
		return fmt.Errorf("invalid orgao: %d", params.Orgao)
	}
	if strings.TrimSpace(params.Cargo) == "" {
		return fmt.Errorf("cargo required")
	}
	if strings.TrimSpace(params.Login) == "" {
		return fmt.Errorf("login required")
	}
	if strings.TrimSpace(params.Senha) == "" {
		return fmt.Errorf("senha required")
	}
	if params.Usuario <= 0 {
		return fmt.Errorf("invalid usuario: %d", params.Usuario)
	}

	// O WSSEI (PHP legado) interpreta o valor de "cargo" como Latin-1
	// mesmo recebendo JSON. Transcodamos apenas esse campo para bytes
	// Latin-1 antes de montar o body.
	bodyBytes, err := json.Marshal(struct {
		Documento int             `json:"documento"`
		Orgao     int             `json:"orgao"`
		Cargo     json.RawMessage `json:"cargo"`
		Login     string          `json:"login"`
		Senha     string          `json:"senha"`
		Usuario   int             `json:"usuario"`
	}{
		Documento: params.Documento,
		Orgao:     params.Orgao,
		Cargo:     jsonStringLatin1(params.Cargo),
		Login:     params.Login,
		Senha:     params.Senha,
		Usuario:   params.Usuario,
	})
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/documento/assinar", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("response error: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[struct{}]
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	if !env.Sucesso {
		return fmt.Errorf("assinar failed: %s", env.Mensagem)
	}

	return nil
}

// DocExternoParams tipo utilizado na funcao CadastrarDocumentoExterno
type DocExternoParams struct {
	// procedimento obrigatorio
	Procedimento               int
	IdUnidadeGeradoraProtocolo int
	Numero                     string
	// IdSerie obrigatorio
	IdSerie int
	// DataElaboracao obrigatorio
	DataElaboracao    string
	IdTipoConferencia int
	Assuntos          string
	Interessados      string
	Remetente         int
	Observacao        string
	// NivelAcesso obrigatorio
	NivelAcesso     int
	IdHipoteseLegal int
	// Anexo obrigatorio
	Anexo io.Reader
	// GrauSigilo obrigatorio
	GrauSigilo string
}

// DocumentoCadastrado representa um novo documento interno ou externo criado pela API
type DocumentoCadastrado struct {
	IDDocumento                 string `json:"idDocumento"`
	ProtocoloDocumentoFormatado string `json:"protocoloDocumentoFormatado"`
}

// CadastrarDocumentoExterno cadastra um novo documento externo.
func (c *Client) CadastrarDocumentoExterno(ctx context.Context, params DocExternoParams) (*DocumentoCadastrado, error) {
	if params.Procedimento <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", params.Procedimento)
	}
	if params.IdSerie <= 0 {
		return nil, fmt.Errorf("invalid idSerie: %d", params.IdSerie)
	}
	if strings.TrimSpace(params.DataElaboracao) == "" {
		return nil, fmt.Errorf("dataElaboracao required")
	}
	if params.NivelAcesso <= 0 {
		return nil, fmt.Errorf("invalid nivelAcesso: %d", params.NivelAcesso)
	}
	if params.Anexo == nil {
		return nil, fmt.Errorf("anexo required")
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	writer.WriteField("procedimento", strconv.Itoa(params.Procedimento))
	writer.WriteField("idSerie", strconv.Itoa(params.IdSerie))
	writer.WriteField("dataElaboracao", params.DataElaboracao)
	writer.WriteField("nivelAcesso", strconv.Itoa(params.NivelAcesso))
	writer.WriteField("grauSigilo", params.GrauSigilo)

	if params.IdUnidadeGeradoraProtocolo != 0 {
		writer.WriteField("idUnidadeGeradoraProtocolo", strconv.Itoa(params.IdUnidadeGeradoraProtocolo))
	}
	if params.Numero != "" {
		writer.WriteField("numero", params.Numero)
	}
	if params.IdTipoConferencia != 0 {
		writer.WriteField("idTipoConferencia", strconv.Itoa(params.IdTipoConferencia))
	}
	if params.Assuntos != "" {
		writer.WriteField("assuntos", params.Assuntos)
	}
	if params.Interessados != "" {
		writer.WriteField("interessados", params.Interessados)
	}
	if params.Remetente != 0 {
		writer.WriteField("remetente", strconv.Itoa(params.Remetente))
	}
	if params.Observacao != "" {
		writer.WriteField("observacao", params.Observacao)
	}
	if params.IdHipoteseLegal != 0 {
		writer.WriteField("idHipoteseLegal", strconv.Itoa(params.IdHipoteseLegal))
	}

	part, err := writer.CreateFormField("anexo")
	if err != nil {
		return nil, fmt.Errorf("create form field error: %w", err)
	}
	_, err = io.Copy(part, params.Anexo)
	if err != nil {
		return nil, fmt.Errorf("copy anexo error: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("writer close error: %w", err)
	}

	url := fmt.Sprintf("%s/documento/%d/externo/criar", c.endpoint, params.Procedimento)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[DocumentoCadastrado]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// DocInterno tipo utilizado na funcao CadastrarDocumentoInterno.
type DocInternoParams struct {
	// Procedimento obrigatorio
	Procedimento               int `json:"procedimento"`
	IdUnidadeGeradoraProtocolo int `json:"idUnidadeGeradoraProtocolo"`
	// IdSerie obrigatorio
	IdSerie      int    `json:"idSerie"`
	Assuntos     string `json:"assuntos"`
	Interessados string `json:"interessados"`
	// Observacao obrigatorio
	Observacao string `json:"observacao"`
	// NivelAcesso obrigatorio
	NivelAcesso              int    `json:"nivelAcesso"`
	IdHipoteseLegal          int    `json:"idHipoteseLegal"`
	IdTextoPadraoInterno     int    `json:"idTextoPadraoInterno"`
	ProtocoloDocumentoModelo string `json:"protocoloDocumentoModelo"`
	Descricao                string `json:"descricao"`
	Destinatarios            string `json:"destinatarios"`
}

// CadastrarDocumentoInterno cadastra um novo documento interno.
func (c *Client) CadastrarDocumentoInterno(ctx context.Context, params DocInternoParams) (*DocumentoCadastrado, error) {
	if params.Procedimento <= 0 {
		return nil, fmt.Errorf("invalid procedimento: %d", params.Procedimento)
	}
	if params.IdSerie <= 0 {
		return nil, fmt.Errorf("invalid idSerie: %d", params.IdSerie)
	}
	if strings.TrimSpace(params.Observacao) == "" {
		return nil, fmt.Errorf("observacao required")
	}
	if params.NivelAcesso <= 0 {
		return nil, fmt.Errorf("invalid nivelAcesso: %d", params.NivelAcesso)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/documento/%d/interno/criar", c.endpoint, params.Procedimento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Envelope[DocumentoCadastrado]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	return &result.Data, nil
}

// usuario tipo utilizado na funcao AssinarBlocoDocumentos.
type BlocoDocumentosParams struct {
	// Todos os parametros sao obrigatorios
	ArrDocumento string `json:"arrDocumento"`
	Orgao        int    `json:"orgao"`
	Cargo        string `json:"cargo"`
	Login        string `json:"login"`
	Senha        string `json:"senha"`
	Usuario      int    `json:"usuario"`
}

// AssinarBlocoDocumentos realiza a assinatura de um ou mais documentos.
func (c *Client) AssinarBlocoDocumentos(ctx context.Context, params BlocoDocumentosParams) error {
	if strings.TrimSpace(params.ArrDocumento) == "" {
		return fmt.Errorf("arrDocumento required")
	}
	if params.Orgao <= 0 {
		return fmt.Errorf("invalid orgao: %d", params.Orgao)
	}
	if strings.TrimSpace(params.Cargo) == "" {
		return fmt.Errorf("cargo required")
	}
	if strings.TrimSpace(params.Login) == "" {
		return fmt.Errorf("login required")
	}
	if strings.TrimSpace(params.Senha) == "" {
		return fmt.Errorf("senha required")
	}
	if params.Usuario <= 0 {
		return fmt.Errorf("invalid usuario: %d", params.Usuario)
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/documento/assinar/bloco", c.endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("response error: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[struct{}]
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	if !env.Sucesso {
		return fmt.Errorf("assinar bloco failed: %s", env.Mensagem)
	}

	return nil
}

// Doc tipo utilizado na funcao DocumentoCiencia.
type Doc struct {
	Documento int `json:"documento"`
}

// DocumentoCiencia da ciencia no documento.
func (c *Client) DocumentoCiencia(ctx context.Context, documento Doc) error {
	if documento.Documento <= 0 {
		return fmt.Errorf("invalid documento: %d", documento.Documento)
	}

	url := fmt.Sprintf("%s/documento/ciencia", c.endpoint)

	bodyBytes, err := json.Marshal(documento)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("response error: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[struct{}]
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	if !env.Sucesso {
		return fmt.Errorf("ciencia failed: %s", env.Mensagem)
	}

	return nil

}

// ListaAlteraDocExterno tipo utilizado na funcao AlterarDocumentoExterno.
type ListaAlteraDocExterno struct {
	// Documento obrigatorio
	Documento int    `jsoon:"documento"`
	Numero    string `json:"numero"`
	// IdSerie obrigatorio
	IDSerie int `json:"idSerie"`
	// DataElaboracao obrigatorio
	DataElaboracao    string `json:"dataElaboracao"`
	IDTipoConferencia int    `json:"idTipoConferencia"`
	Assuntos          string `json:"assuntos"`
	Interessados      string `json:"interessados"`
	Remetente         int    `json:"remetente"`
	Observacao        string `json:"observacao"`
	// NivelAcesso obrigatorio
	NivelAcesso     int `json:"nivelAcesso"`
	IDHipoteseLegal int `json:"idHipoteseLegal"`
	Anexo           io.Reader
	// GrauSigilo obrigatorio
	GrauSigilo string `json:"grauSigilo"`
}

// AlterarDocumentoExterno realiza a edição dos metadados de um documento externo.
func (c *Client) AlterarDocumentoExterno(ctx context.Context, params ListaAlteraDocExterno) error {
	if params.Documento <= 0 {
		return fmt.Errorf("invalid documento: %d", params.Documento)
	}
	if params.IDSerie <= 0 {
		return fmt.Errorf("invalid idSerie: %d", params.IDSerie)
	}
	if strings.TrimSpace(params.DataElaboracao) == "" {
		return fmt.Errorf("dataElaboracao required")
	}
	if params.NivelAcesso <= 0 {
		return fmt.Errorf("invalid nivelAcesso: %d", params.NivelAcesso)
	}
	if params.Anexo == nil {
		return fmt.Errorf("anexo required")
	}
	if strings.TrimSpace(params.GrauSigilo) == "" {
		return fmt.Errorf("grauSigilo required")
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	writer.WriteField("documento", strconv.Itoa(params.Documento))
	writer.WriteField("idSerie", strconv.Itoa(params.IDSerie))
	writer.WriteField("dataElaboracao", params.DataElaboracao)
	writer.WriteField("nivelAcesso", strconv.Itoa(params.NivelAcesso))
	writer.WriteField("grauSigilo", params.GrauSigilo)

	if params.Numero != "" {
		writer.WriteField("numero", params.Numero)
	}

	if params.IDTipoConferencia != 0 {
		writer.WriteField("idTipoConferencia", strconv.Itoa(params.IDTipoConferencia))
	}

	if params.Assuntos != "" {
		writer.WriteField("assuntos", params.Assuntos)
	}

	if params.Interessados != "" {
		writer.WriteField("interessados", params.Interessados)
	}

	if params.Remetente != 0 {
		writer.WriteField("remetente", strconv.Itoa(params.Remetente))
	}

	if params.Observacao != "" {
		writer.WriteField("observacao", params.Observacao)
	}

	if params.IDHipoteseLegal != 0 {
		writer.WriteField("idHipoteseLegal", strconv.Itoa(params.IDHipoteseLegal))
	}

	part, err := writer.CreateFormField("anexo")
	if err != nil {
		return fmt.Errorf("create form field error: %w", err)
	}

	_, err = io.Copy(part, params.Anexo)
	if err != nil {
		return fmt.Errorf("copy anexo error: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("writer close error: %w", err)
	}

	url := fmt.Sprintf("%s/documento/externo/%d/alterar", c.endpoint, params.Documento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// ListaAlteraDocInterno tipo utilizado na funcao AlterarDocumentoInterno.
type ListaAlteraDocInterno struct {
	// Documento obrigatorio
	Documento    int    `json:"documento"`
	Assuntos     string `json:"assuntos"`
	Interessados string `json:"interessados"`
	// Observacao obrigatorio
	Observacao string `json:"observacao"`
	// NivelAcesso obrigatorio
	NivelAcesso     int    `json:"nivelAcesso"`
	IDHipoteseLegal int    `json:"idHipoteseLegal"`
	Descricao       string `json:"descricao"`
	Destinatarios   string `json:"destinatarios"`
}

// AlterarDocumentoInterno realiza a edição dos metadados de um documento interno.
func (c *Client) AlterarDocumentoInterno(ctx context.Context, params ListaAlteraDocInterno) error {
	if params.Documento <= 0 {
		return fmt.Errorf("invalid documento: %d", params.Documento)
	}
	if strings.TrimSpace(params.Observacao) == "" {
		return fmt.Errorf("observacao required")
	}
	if params.NivelAcesso <= 0 {
		return fmt.Errorf("invalid nivelAcesso: %d", params.NivelAcesso)
	}
	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/documento/interno/%d/alterar", c.endpoint, params.Documento)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("response error: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[struct{}]
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	if !env.Sucesso {
		return fmt.Errorf("alterar failed: %s", env.Mensagem)
	}

	return nil
}

// BlocosAssinatura tipo utilizado no ListarBlocosAssinaturaDocumento.
type BlocosAssinatura struct {
	IdDocumento        string   `json:"idDocumento"`
	ProtocoloFormatado string   `json:"protocoloFormatado"`
	DataGeracao        string   `json:"dataGeracao"`
	IdSerie            string   `json:"idSerie"`
	NomeSerie          string   `json:"nomeSerie"`
	Blocos             []string `json:"blocos"`
}

// ListarBlocosAssinaturaDocumento Retorna a lista de Blocos de Assinatura do Documento.
func (c *Client) ListarBlocosAssinaturaDocumento(ctx context.Context, documento int) (*BlocosAssinatura, error) {
	if documento <= 0 {
		return nil, fmt.Errorf("invalid documento: %d", documento)
	}

	url := fmt.Sprintf(
		"%s/documento/%d/bloco/assinatura/listar", c.endpoint, documento)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[BlocosAssinatura]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, fmt.Errorf("consultar failed %d: %s", documento, result.Mensagem)
	}

	return &result.Data, nil

}

// ListarSugestoesParams tipo utilizado na funcao ListarSugestoesAssuntoTipoDocumento.
type ListarSugestoesParams struct {
	// Serie e obrigatorio
	Serie  int
	Limit  int
	Start  string
	ID     string
	Filter string
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListarSugestoesParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != "" {
		q.Set("start", p.Start)
	}
	if p.ID != "" {
		q.Set("id", p.ID)
	}
	if p.Filter != "" {
		q.Set("filter", p.Filter)
	}
	return q
}

// SugestoesAssunto representa o tipo utilizado na funcao ListarSugestoesAssuntoTipoDocumento.
type SugestoesAssunto struct {
	CodigoEstruturadoFormatado string `json:"codigoestruturadoformatado"`
	Descricao                  string `json:"descricao"`
	CodigoEstruturado          string `json:"codigoestruturado"`
	ID                         string `json:"id"`
}

// ListarSugestoesAssuntoTipoDocumento Retorna a Lista de Sugestões de Assuntos por Tipo de Documento.
func (c *Client) ListarSugestoesAssuntoTipoDocumento(ctx context.Context, params ListarSugestoesParams) ([]SugestoesAssunto, error) {
	if params.Serie <= 0 {
		return nil, fmt.Errorf("invalid serie: %d", params.Serie)
	}

	url := fmt.Sprintf("%s/documento/assunto/sugestao/%d/listar", c.endpoint, params.Serie)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]SugestoesAssunto]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	if result.Sucesso != true {
		return nil, fmt.Errorf("consultar failed %d: %s", params.Serie, result.Mensagem)
	}

	return result.Data, nil

}

// TiposConferenciasParams seleção de todos os params possíveis
type TiposConferenciasParams struct {
	// Query Params
	Limit  int
	Start  int
	Filter string
	ID     int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p TiposConferenciasParams) value() url.Values {
	q := make(url.Values)

	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.ID != 0 {
		q.Set("id", strconv.Itoa(p.ID))
	}
	if p.Filter != "" {
		q.Set("filter", p.Filter)
	}

	return q
}

// Tipo para representar cada item retornado dentro do array "data"
type ItemConferencia struct {
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
}

// Tipo utilizado como retorno da funcao PesquisarTipoConferencia
type TiposConferencia struct {
	Sucesso bool              `json:"sucesso"`
	Data    []ItemConferencia `json:"data"`
	Total   string            `json:"total"`
}

// PesquisarTipoConferencia Retorna a pesquisa de Tipo de Conferência
func (c *Client) PesquisarTipoConferencia(ctx context.Context, params TiposConferenciasParams) ([]ItemConferencia, int, error) {

	url := fmt.Sprintf("%s/documento/tipoconferencia/pesquisar", c.endpoint)
	if q := params.value().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]ItemConferencia]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("decode error: %w", err)
	}

	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("pesquisar failed: %s", result.Mensagem)
	}

	total, err := strconv.Atoi(result.Total)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid total: %w", err)
	}

	return result.Data, total, nil

}

// ListarAssinaturasDocumento retorna a lista de assinaturas de um documento específico.
func (c *Client) ListarAssinaturasDocumento(ctx context.Context, documento int) ([]Assinatura, error) {
	if documento <= 0 {
		return nil, fmt.Errorf("invalid documento: %d", documento)
	}

	url := fmt.Sprintf("%s/documento/listar/assinaturas/%d", c.endpoint, documento)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]Assinatura]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if !result.Sucesso {
		return nil, fmt.Errorf("listar assinaturas failed %d: %s", documento, result.Mensagem)
	}

	return result.Data, nil
}

// Assinatura representa os dados de uma assinatura retornada pela API.
type Assinatura struct {
	Nome    string `json:"nome"`
	Cargo   string `json:"cargo"`
	Unidade string `json:"unidade"`
}

// ListarCienciasDocumento retorna a lista de ciências de um documento específico através do seu protocolo.
func (c *Client) ListarCienciasDocumento(ctx context.Context, protocolo string) ([]Ciencia, error) {
	if protocolo == "" {
		return nil, fmt.Errorf("protocolo required")
	}

	url := fmt.Sprintf("%s/documento/listar/ciencia/%s", c.endpoint, protocolo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]Ciencia]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if !result.Sucesso {
		return nil, fmt.Errorf("listar ciencias failed %s: %s", protocolo, result.Mensagem)
	}

	return result.Data, nil
}

// Ciencia representa os dados de uma ciência retornada pela API.
type Ciencia struct {
	Data      string `json:"data"`
	Unidade   string `json:"unidade"`
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
}
