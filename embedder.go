package main

import (
	"context"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/embeddings/cybertron"
)

type Embedder interface {
	EmbedDocuments(ctx context.Context, docs []string) ([][]float32, error)
}

func NewCybertronEmbedder(model string) (Embedder, error) {
	embc, err := cybertron.NewCybertron(
		cybertron.WithModelsDir("models"),
		cybertron.WithModel(string(model)),
	)
	if err != nil {
		return nil, err
	}

	emb, err := embeddings.NewEmbedder(
		embc,
		embeddings.WithStripNewLines(false),
	)
	if err != nil {
		return nil, err
	}

	return &CybertronEmbedder{emb}, nil
}

type CybertronEmbedder struct {
	embeddings.Embedder
}
