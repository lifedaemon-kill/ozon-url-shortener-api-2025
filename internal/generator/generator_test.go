package generator

import (
	"main/internal/config"
	"strings"
	"testing"
)

func TestUniqueSeqGenerator_Generate(t *testing.T) {
	cfg := config.UrlGenConfig{
		BaseHost: "http://localhost",
		Alphabet: "abc123",
		Length:   10,
	}
	gen := New(cfg)

	for i := 0; i < 100; i++ { // проверяем 100 сгенерированных ссылок
		result := gen.Generate()

		if len(result) != cfg.Length {
			t.Errorf("Expected length %d, got %d", cfg.Length, len(result))
		}

		for _, r := range result {
			if !strings.ContainsRune(cfg.Alphabet, r) {
				t.Errorf("Generated invalid character: %c", r)
			}
		}
	}
}

func TestUniqueSeqGenerator_BaseHost(t *testing.T) {
	cfg := config.UrlGenConfig{
		BaseHost: "http://localhost",
		Alphabet: "abc123",
		Length:   10,
	}
	gen := New(cfg)

	if gen.BaseHost() != cfg.BaseHost {
		t.Errorf("Expected BaseHost %s, got %s", cfg.BaseHost, gen.BaseHost())
	}
}
