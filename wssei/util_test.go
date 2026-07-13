package wssei

import (
	"encoding/json"
	"testing"
)

func TestJSONStringLatin1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string // bytes esperados como string (podem incluir bytes >127)
	}{
		{
			name: "ASCII puro",
			in:   "Servidor",
			want: `"Servidor"`,
		},
		{
			name: "Acentos Latin-1",
			// ú = U+00FA, byte Latin-1 0xFA
			in:   "Público",
			want: "\"P\xFAblico\"",
		},
		{
			name: "Composto real do SEI",
			in:   "Servidor (a) Público (a)",
			want: "\"Servidor (a) P\xFAblico (a)\"",
		},
		{
			name: "Vários acentos",
			// á=E1 é=E9 í=ED ó=F3 ú=FA ç=E7 ã=E3 õ=F5
			in:   "áéíóúçãõ",
			want: "\"\xE1\xE9\xED\xF3\xFA\xE7\xE3\xF5\"",
		},
		{
			name: "Escapa aspas e barra",
			in:   `a"b\c`,
			want: `"a\"b\\c"`,
		},
		{
			name: "Escapa controles",
			in:   "a\nb\tc",
			want: `"a\nb\tc"`,
		},
		{
			name: "Fora do Latin-1 mantém UTF-8",
			// U+20AC (€) = bytes UTF-8 E2 82 AC
			in:   "€",
			want: "\"\xE2\x82\xAC\"",
		},
		{
			name: "String vazia",
			in:   "",
			want: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jsonStringLatin1(tt.in)
			if string(got) != tt.want {
				t.Errorf("jsonStringLatin1(%q):\n  got  = % x (%q)\n  want = % x (%q)",
					tt.in, []byte(got), string(got), []byte(tt.want), tt.want)
			}
		})
	}
}

func TestJSONStringLatin1_UsableAsRawMessage(t *testing.T) {
	// Garante que o valor produzido é usável dentro de um json.Marshal
	// via json.RawMessage sem que o marshal padrão re-valide/escape os
	// bytes Latin-1 fora do ASCII.
	payload := struct {
		Cargo json.RawMessage `json:"cargo"`
	}{
		Cargo: jsonStringLatin1("Público"),
	}
	out, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// Deve conter o byte 0xFA cru (Latin-1) e não os bytes UTF-8 0xC3 0xBA.
	want := "{\"cargo\":\"P\xFAblico\"}"
	if string(out) != want {
		t.Errorf("marshal:\n  got  = % x (%q)\n  want = % x (%q)",
			out, string(out), []byte(want), want)
	}
}
