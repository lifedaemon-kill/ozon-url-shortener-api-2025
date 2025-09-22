package generator

import (
	"main/internal/config"
	"math/rand"
	"strings"
)

type Generator interface {
	Generate() string
	BaseHost() string
}

type SeqGenerator struct {
	baseHost string
	alphabet string
	length   int
}

func New(cfg config.UrlGenConfig) *SeqGenerator {
	return &SeqGenerator{
		baseHost: cfg.BaseHost,
		alphabet: cfg.Alphabet,
		length:   cfg.Length,
	}
}

func (g *SeqGenerator) Generate() string {
	var builder strings.Builder
	builder.Grow(g.length)

	symbolsLen := len(g.alphabet)

	for range g.length {
		builder.WriteRune(rune(g.alphabet[rand.Intn(symbolsLen)]))
	}

	return builder.String()
}

func (g *SeqGenerator) BaseHost() string {
	return g.baseHost
}
