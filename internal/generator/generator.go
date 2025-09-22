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

type UniqueSeqGenerator struct {
	baseHost string
	alphabet string
	length   int
}

func New(cfg config.UrlGenConfig) *UniqueSeqGenerator {
	return &UniqueSeqGenerator{
		baseHost: cfg.BaseHost,
		alphabet: cfg.Alphabet,
		length:   cfg.Length,
	}
}

func (g *UniqueSeqGenerator) Generate() string {
	var builder strings.Builder
	builder.Grow(g.length)

	symbolsLen := len(g.alphabet)

	for range g.length {
		builder.WriteRune(rune(g.alphabet[rand.Intn(symbolsLen)]))
	}

	return builder.String()
}

func (g *UniqueSeqGenerator) BaseHost() string {
	return g.baseHost
}
