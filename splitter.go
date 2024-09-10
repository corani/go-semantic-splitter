package main

import (
	"context"
	"strings"

	"github.com/gosbd/gosbd"
)

type DocumentSplit struct {
	Documents []string
	Score     float32
	Count     int
}

func NewSemanticSplitter(emb Embedder, cnt TokenCounter) *SemanticSplitter {
	return &SemanticSplitter{
		emb:          emb,
		cnt:          cnt,
		seg:          gosbd.NewSegmenter("en"),
		lang:         "en",
		window:       5,
		threshold:    0.2,
		minChunkSize: 100,
		maxChunkSize: 500,
	}
}

type SemanticSplitter struct {
	emb          Embedder
	cnt          TokenCounter
	seg          gosbd.Segmenter
	lang         string
	window       int
	threshold    float32
	minChunkSize int
	maxChunkSize int
}

func (ss *SemanticSplitter) Split(ctx context.Context, text string) ([]DocumentSplit, error) {
	text = ss.clean(text)
	splits := ss.splitToSentences(text)

	encodedSplits, err := ss.emb.EmbedDocuments(ctx, splits)
	if err != nil {
		return nil, err
	}

	similarities := ss.calculateSimilarityScores(encodedSplits)

	return ss.splitDocuments(splits, similarities)
}

func (ss *SemanticSplitter) clean(text string) string {
	// replace multiple whitespace (including newlines) with a single space
	text = strings.Join(strings.Fields(text), " ")

	return text
}

func (ss *SemanticSplitter) splitToSentences(text string) []string {
	parts := ss.seg.Segment(text)

	for i, v := range parts {
		parts[i] = strings.TrimSpace(v)
	}

	return parts
}

func (ss *SemanticSplitter) calculateSimilarityScores(docs [][]float32) []float32 {
	var result []float32

	result = append(result, 0)

	for i := 1; i < len(docs); i++ {
		start := max(0, i-ss.window)

		cumulativeContext := mean(docs[start:i])
		score := dot(cumulativeContext, docs[i]) / (norm(cumulativeContext)*norm(docs[i]) + 1e-10)

		result = append(result, score)
	}

	return result
}

func (ss *SemanticSplitter) splitDocuments(docs []string, scores []float32) ([]DocumentSplit, error) {
	var result []DocumentSplit

	current := DocumentSplit{}

	for i, doc := range docs {
		docTokenCount := ss.cnt.Count(doc)

		// Check if the current split is a split point based on similarity
		if i+1 < len(scores) && scores[i+1] < ss.threshold {
			// If so, and the current chunk plus this split exceeds the min
			// token limit, emit the current chunk and start a new one.
			if current.Count+docTokenCount >= ss.minChunkSize {
				current.Documents = append(current.Documents, doc)
				current.Count += docTokenCount
				current.Score = scores[i]

				result = append(result, current)
				current = DocumentSplit{}

				continue
			}
		}

		// Check if adding the current document exceeds the max token limit
		if current.Count+docTokenCount >= ss.maxChunkSize && current.Count >= ss.minChunkSize {
			// If so, and the current chunk already exceeds the min token limit,
			// emit the current chunk and start a new one.
			if current.Count >= ss.minChunkSize {
				// triggered by token limit, not by similarity
				current.Score = 0

				result = append(result, current)
				current = DocumentSplit{}
			}
		}

		current.Documents = append(current.Documents, doc)
		current.Count += docTokenCount
	}

	// handle last split
	if len(current.Documents) > 0 {
		result = append(result, current)
	}

	return result, nil
}
