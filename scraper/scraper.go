// Package scraper extrai dados diretamente das páginas HTML do SEI, para
// informações que não são expostas nem pelo módulo WSSEI (REST) nem pela API
// SOAP legada (SeiWS.php) — por exemplo, a lista de órgãos na tela de login,
// o download do PDF de um processo ou a listagem de documentos na página de
// acesso externo.
package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
)

const (
	// O número aproximado de órgãos no SEI-MG.
	orgaosLen = 80
)

var (
	// ErrProcessoVazio é o erro retornado quando o processo não possui
	// conteúdo para download.
	ErrProcessoVazio = errors.New("processo is empty")
	// ErrInvalidContentType é o erro retornado quando o download do processo
	// não retorna um PDF.
	ErrInvalidContentType = errors.New("invalid content-type")
)

// Scraper é responsável por extrair dados direto das páginas HTML do SEI.
type Scraper struct {
	baseURL string
	http    *http.Client
}

// NewScraper cria um [*Scraper] para a URL base do SEI informada, usando o
// [http.DefaultClient].
func NewScraper(baseURL string) *Scraper {
	return &Scraper{
		baseURL: baseURL,
		http:    http.DefaultClient,
	}
}

// Orgao representa um órgão exposto na página de login do SEI.
type Orgao struct {
	ID   int
	Nome string
}

// ListOrgaos lista os órgãos disponíveis no SEI na página de login.
func (s *Scraper) ListOrgaos(ctx context.Context) ([]Orgao, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	orgaos := make([]Orgao, 0, orgaosLen)
	doc.Find("#selOrgao > option").Each(func(i int, s *goquery.Selection) {
		value, ok := s.Attr("value")
		if !ok {
			return
		}

		// Pegamos apenas valores inteiros válidos.
		id, err := strconv.Atoi(value)
		if err != nil {
			return
		}

		orgaos = append(orgaos, Orgao{
			ID:   id,
			Nome: strings.TrimSpace(s.Text()),
		})
	})

	return orgaos, nil
}

// DownloadProcedimento baixa o PDF de um processo a partir do seu Link de
// Acesso Externo. É uma extensão da API do SEI: a página é raspada para
// coletar os identificadores dos documentos e então um POST gera o PDF.
//
// O chamador é responsável por fechar o [io.ReadCloser] retornado.
func (s *Scraper) DownloadProcedimento(ctx context.Context, linkAcessoExterno string) (io.ReadCloser, error) {
	u, err := url.Parse(linkAcessoExterno)
	if err != nil {
		return nil, err
	}

	// Adiciona www ao host em produção. Essa mudança foi introduzida após
	// a atualização do SEI versão 5 em MG.
	if u.Host == "sei.mg.gov.br" {
		u.Host = fmt.Sprintf("www.%s", u.Host)
	}
	linkAcessoExterno = u.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, linkAcessoExterno, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get info: status %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	formData := make(url.Values)
	doc.Find("#hdnInfraItensHash, #hdnInfraItens").Each(func(i int, sel *goquery.Selection) {
		name, ok := sel.Attr("name")
		if ok {
			val, _ := sel.Attr("value")
			formData.Set(name, val)
		}
	})

	var listaIDs []string
	doc.Find("#tblDocumentos tr input[type='checkbox']").Each(func(i int, sel *goquery.Selection) {
		if val, ok := sel.Attr("value"); ok {
			listaIDs = append(listaIDs, val)
		}
	})

	if len(listaIDs) == 0 {
		return nil, ErrProcessoVazio
	}

	formData.Set("hdnInfraItensSelecionados", strings.Join(listaIDs, ","))
	formData.Set("hdnFlagGerar", "1")

	reqPost, err := http.NewRequestWithContext(ctx, http.MethodPost, linkAcessoExterno, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resPost, err := s.http.Do(reqPost)
	if err != nil {
		return nil, err
	}

	contentType := resPost.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/pdf") {
		defer resPost.Body.Close()
		return nil, ErrInvalidContentType
	}

	return resPost.Body, nil
}

// LinhaDocumento representa uma linha da tabela de documentos na página de
// acesso externo de um processo.
type LinhaDocumento struct {
	Numero  string `json:"numero"`
	Link    string `json:"link"`
	Tipo    string `json:"tipo"`
	Data    string `json:"data"`
	Unidade string `json:"unidade"`
}

// ListarDocumentos retorna a lista de documentos raspando a página de acesso
// externo de um processo, a partir do seu Link de Acesso Externo. Os links dos
// documentos são resolvidos de forma absoluta em relação ao linkAcessoExterno
// informado.
func (s *Scraper) ListarDocumentos(ctx context.Context, linkAcessoExterno string) ([]LinhaDocumento, error) {
	base, err := url.Parse(linkAcessoExterno)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, linkAcessoExterno, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	dec := charmap.ISO8859_1.NewDecoder()
	rd := dec.Reader(res.Body)

	doc, err := goquery.NewDocumentFromReader(rd)
	if err != nil {
		return nil, err
	}

	documentos := make([]LinhaDocumento, 0)
	doc.Find("#tblDocumentos tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}

		var documento LinhaDocumento
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				return
			case 1:
				link := s.Children().First()
				numero := link.Text()
				href, ok := link.Attr("href")
				if ok {
					documento.Numero = numero
					if ref, err := url.Parse(href); err == nil {
						documento.Link = base.ResolveReference(ref).String()
					} else {
						documento.Link = href
					}
				}
			case 2:
				documento.Tipo = s.Text()
			case 3:
				documento.Data = s.Text()
			case 4:
				documento.Unidade = s.Text()
			default:
				return
			}
		})

		documentos = append(documentos, documento)
	})

	return documentos, nil
}
