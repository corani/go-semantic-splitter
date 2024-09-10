package main

import (
	"context"
	"fmt"
	"io"
	"os"
)

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func main1(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: %s <filename>", args[0])
	}

	text, err := readFile(args[1])
	if err != nil {
		return err
	}

	emb, err := NewCybertronEmbedder("sentence-transformers/all-MiniLM-L6-v2")
	if err != nil {
		return err
	}

	cnt, err := NewTikTokenCounter("cl100k_base")
	if err != nil {
		return err
	}

	splitter := NewSemanticSplitter(emb, cnt)

	docs, err := splitter.Split(ctx, text)
	if err != nil {
		return err
	}

	for i, doc := range docs {
		fmt.Printf("Split %d: %d docs, %d tokens, score: %f\n", i, len(doc.Documents), doc.Count, doc.Score)
		fmt.Println(doc.Documents)
	}

	return nil
}

func main() {
	if err := main1(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
