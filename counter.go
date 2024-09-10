package main

import "github.com/pkoukk/tiktoken-go"

type TokenCounter interface {
	Count(string) int
}

func NewTikTokenCounter(model string) (TokenCounter, error) {
	tkm, err := tiktoken.GetEncoding(model)
	if err != nil {
		return nil, err
	}

	return TikTokenCounter{tkm}, nil
}

type TikTokenCounter struct {
	tkm *tiktoken.Tiktoken
}

func (t TikTokenCounter) Count(s string) int {
	return len(t.tkm.Encode(s, nil, nil))
}
