# WSSEI

Client Go para o [SEI](https://www.gov.br/gestao/pt-br/assuntos/sei) (Sistema
Eletrônico de Informações), cobrindo as três interfaces disponíveis:

- [`wssei`](./wssei) — módulo [WSSEI](https://pengovbr.github.io/mod-wssei/#/)
  (REST), autenticado por usuário e senha. É a interface principal, usada
  para processos, documentos, blocos de assinatura, marcadores, etc.
- [`seiws`](./seiws) — API SOAP legada (`SeiWS.php`), autenticada por sigla
  de sistema e token de serviço. Usada apenas para operações que ainda não
  têm equivalente no WSSEI.
- [`scraper`](./scraper) — extração de dados diretamente das páginas HTML
  do SEI, para informações que não são expostas por nenhuma das APIs (ex:
  lista de órgãos na tela de login).
- [`soap`](./soap) — tipos auxiliares para serializar e deserializar
  envelopes SOAP, compartilhados pelo `seiws`.

Projetado para ser compartilhado entre múltiplos projetos.

## Instalação

```bash
go get github.com/automatiza-mg/WSSEI
```

Requer Go 1.25 ou superior.

## wssei — módulo WSSEI (REST)

Autentica com usuário e senha uma única vez e reaproveita o token gerado em
cada requisição, renovando-o automaticamente quando o servidor responde
401/403.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strconv"

    "github.com/automatiza-mg/WSSEI/wssei"
)

func main() {
    orgao, _ := strconv.Atoi(os.Getenv("SEI_ORGAO"))

    client := wssei.NewClient(wssei.Config{
        BaseURL: os.Getenv("SEI_BASE_URL"),
        Usuario: os.Getenv("SEI_USUARIO"),
        Senha:   os.Getenv("SEI_SENHA"),
        Orgao:   orgao,
    })

    ctx := context.Background()

    processo, err := client.ConsultarProcesso(ctx, 1234567)
    if err != nil {
        panic(err)
    }

    fmt.Println(processo.ProtocoloFormatado)
}
```

### Consultar e listar processos

```go
processos, total, err := client.ListarProcessos(ctx, wssei.ListarProcessosParams{
    Limit: 10,
    Start: 0,
})
```

### Documentos

```go
doc, err := client.ConsultarDocumentoInterno(ctx, protocolo)

// Baixa o conteúdo de um documento externo (anexo).
body, contentType, err := client.BaixarAnexo(ctx, protocolo)
defer body.Close()
```

### Blocos de assinatura

```go
blocos, total, err := client.PesquisarBlocoAssinatura(ctx, wssei.PesquisarBlocoAssinaturaParams{
    Estado: wssei.EstadoSituacaoDisponibilizado,
    Limit:  20,
})

err = client.AssinarBlocoAssinatura(ctx, bloco, wssei.AssinarBlocoAssinaturaParams{
    Orgao:   orgao,
    Cargo:   "Servidor (a) Público (a)",
    Login:   usuario,
    Senha:   senha,
    Usuario: idUsuario,
})
```

### Marcadores

```go
marcador, err := client.ConsultarMarcador(ctx, protocolo)

err = client.MarcarProcesso(ctx, protocolo, wssei.MarcadorProcessoParams{
    Texto:    "Aguardando análise",
    Marcador: idMarcador,
})
```

### Configuração

O client é configurado com uma struct `wssei.Config`:

| Campo             | Descrição                                                                     |
| ----------------- | ----------------------------------------------------------------------------- |
| `BaseURL`         | URL base do SEI (ex: `https://www.sei.mg.gov.br`).                            |
| `Usuario`         | Login usado na autenticação.                                                  |
| `Senha`           | Senha usada na autenticação.                                                  |
| `Orgao`           | Id do órgão da autenticação.                                                  |
| `OnAuthenticated` | Callback opcional, invocado após cada autenticação bem-sucedida (ver abaixo). |

### AuthCallback

`OnAuthenticated` recebe os dados retornados pelo WSSEI (`AuthResponse`) e
permite persistir ou observar o token, o `IdUsuario`, as unidades e os
perfis do usuário autenticado. É chamado no login inicial e a cada renovação
automática do token.

```go
cfg.OnAuthenticated = func(ctx context.Context, resp *wssei.AuthResponse) error {
    return cache.SaveToken(ctx, resp.Token)
}
```

### Autenticação isolada

Se você precisa apenas autenticar (sem executar chamadas subsequentes), use
`wssei.Auth` diretamente:

```go
auth := wssei.NewAuth(baseURL)
resp, err := auth.Autenticar(ctx, usuario, senha, orgao)
```

### Envelope

Todas as respostas do WSSEI seguem o formato:

```json
{
    "sucesso": true,
    "mensagem": "",
    "total": "42",
    "data": { ... }
}
```

Os métodos do `Client` já extraem `data` e `total` para o chamador,
transformando `sucesso: false` em erro. O tipo genérico `wssei.Envelope[T]`
é exportado para uso em cenários onde a resposta bruta é necessária.

### Compatibilidade Latin-1

Alguns endpoints do WSSEI (PHP legado) interpretam campos específicos como
Latin-1 mesmo recebendo JSON UTF-8 — o exemplo mais conhecido é o `cargo`
em assinaturas. O client transcodifica automaticamente esses campos,
garantindo que acentos sejam aceitos pelo servidor.

## seiws — API SOAP legada

Autentica com sigla do sistema e token de serviço (credenciais fixas por
aplicação, não por usuário). Usado apenas para operações ainda não expostas
pelo WSSEI.

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/automatiza-mg/WSSEI/seiws"
)

func main() {
    client := seiws.NewClient(seiws.Config{
        URL:                  os.Getenv("SEI_WS_URL"),
        SiglaSistema:         os.Getenv("SEI_SIGLA_SISTEMA"),
        IdentificacaoServico: os.Getenv("SEI_IDENTIFICACAO_SERVICO"),
    })

    ctx := context.Background()

    resp, err := client.ConsultarDocumento(ctx, "0000001")
    if err != nil {
        panic(err)
    }

    fmt.Println(resp.Parametros.DocumentoFormatado)
}
```

Erros de Fault SOAP são retornados como `*soap.Error`, que carrega o status
HTTP e o envelope de Fault decodificado.

## scraper — extração via HTML

Para dados que não estão em nenhuma API. Hoje suporta apenas a listagem de
órgãos na tela de login.

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/automatiza-mg/WSSEI/scraper"
)

func main() {
    s := scraper.NewScraper(os.Getenv("SEI_BASE_URL"))

    orgaos, err := s.ListOrgaos(context.Background())
    if err != nil {
        panic(err)
    }

    for _, o := range orgaos {
        fmt.Printf("%d\t%s\n", o.ID, o.Nome)
    }
}
```

## soap — envelopes SOAP

Tipos auxiliares (`Envelope`, `Body`, `Fault`, `Error`) usados internamente
por `seiws`. Não é necessário usar diretamente, exceto para inspecionar
detalhes de um Fault:

```go
var se *soap.Error
if errors.As(err, &se) {
    fmt.Println(se.Status, se.Fault.Body.Content.Message)
}
```

## Referência

- [Documentação oficial do módulo WSSEI](https://pengovbr.github.io/mod-wssei/#/)

## Licença

[MIT](./LICENSE)
