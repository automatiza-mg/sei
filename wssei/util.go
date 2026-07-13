package wssei

import (
	"bytes"
	"encoding/json"
	"unicode/utf8"
)

// jsonStringLatin1 codifica s como um literal de string JSON (com aspas)
// convertendo cada rune UTF-8 para 1 byte Latin-1 (ISO-8859-1) quando o code
// point couber em Latin-1 (U+0000..U+00FF). Runes fora dessa faixa são
// preservados como seus bytes UTF-8 originais. Bytes inválidos de UTF-8
// também são preservados verbatim.
//
// Usado para campos específicos (ex.: cargo de assinatura) em endpoints do
// WSSEI/SEI (PHP legado) que interpretam o payload como Latin-1 mesmo quando
// o Content-Type declara JSON.
//
// O retorno é um json.RawMessage pronto para uso em structs marshalados com
// encoding/json, evitando que o marshal padrão re-valide os bytes como UTF-8
// (o que corromperia os bytes Latin-1 fora do ASCII).
func jsonStringLatin1(s string) json.RawMessage {
	buf := make([]byte, 0, len(s)+2)
	buf = append(buf, '"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			// Byte inválido de UTF-8: preserva verbatim.
			buf = append(buf, s[i])
			i++
			continue
		}
		// Escapes obrigatórios em JSON.
		switch r {
		case '"':
			buf = append(buf, '\\', '"')
		case '\\':
			buf = append(buf, '\\', '\\')
		case '\b':
			buf = append(buf, '\\', 'b')
		case '\f':
			buf = append(buf, '\\', 'f')
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		default:
			if r < 0x20 {
				// Controles: escapa como \u00XX.
				const hex = "0123456789abcdef"
				buf = append(buf, '\\', 'u', '0', '0', hex[byte(r)>>4], hex[byte(r)&0xF])
			} else if r <= 0xFF {
				buf = append(buf, byte(r))
			} else {
				// Fora do Latin-1: mantém os bytes UTF-8 originais.
				buf = append(buf, s[i:i+size]...)
			}
		}
		i += size
	}
	buf = append(buf, '"')
	return json.RawMessage(buf)
}

// Object embrulha um valor que o WSSEI pode enviar como objeto ou, quando
// ausente, como string vazia "", array vazio [] ou null.
//
// Segue o padrão dos tipos sql.Null: Valid indica se há um valor presente e,
// em caso negativo, Value permanece zerado.
type Object[T any] struct {
	// Value é o valor decodificado. Só é significativo quando Valid é true.
	Value T
	// Valid indica se o WSSEI enviou um valor de fato, em vez de uma das formas
	// vazias.
	Valid bool
}

// UnmarshalJSON decodifica o objeto em Value e marca Valid como true. Se o
// WSSEI enviar uma das formas vazias ("", [] ou null), Value é zerado e Valid
// fica false.
//
// Algumas respostas do WSSEI também enviam o objeto embrulhado em um array
// (ex: [{...}]). Nesses casos o primeiro elemento é usado como Value; se o
// array vier vazio, Valid permanece false.
func (o *Object[T]) UnmarshalJSON(data []byte) error {
	var zero T
	trimmed := bytes.TrimSpace(data)
	switch string(trimmed) {
	case "", `""`, "[]", "null":
		o.Value = zero
		o.Valid = false
		return nil
	}
	if len(trimmed) > 0 && trimmed[0] == '[' {
		var items []T
		if err := json.Unmarshal(trimmed, &items); err != nil {
			return err
		}
		if len(items) == 0 {
			o.Value = zero
			o.Valid = false
			return nil
		}
		o.Value = items[0]
		o.Valid = true
		return nil
	}
	if err := json.Unmarshal(trimmed, &o.Value); err != nil {
		return err
	}
	o.Valid = true
	return nil
}

// Slice é uma lista que o WSSEI pode enviar como array ou, quando vazia, como
// string vazia "" ou null. Nesses casos vazios decodifica para um slice nil.
type Slice[T any] []T

// UnmarshalJSON decodifica o array nos elementos ou, se o WSSEI enviar uma das
// formas vazias ("", {} ou null), mantém o slice nil.
func (s *Slice[T]) UnmarshalJSON(data []byte) error {
	switch string(bytes.TrimSpace(data)) {
	case "", `""`, "{}", "null":
		*s = nil
		return nil
	default:
		return json.Unmarshal(data, (*[]T)(s))
	}
}
