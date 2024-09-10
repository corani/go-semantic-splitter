# go-semantic-splitter

This is a Go implementation of the rolling window semantic splitter algorithm. The algorithm splits
a document into sentences then calculates an embedding vector for each of them. It calculates the
similarity between each sentence and the mean of the preceding window (of size 5), then uses this
information to do the final chunking. This takes into account the semantic similarity and the
min/max chunk size.

## Usage

```bash
./build.sh
./bin/splitter <input-file>
```

## Resources

- [aurelio-labs/semantic-router](https://github.com/aurelio-labs/semantic-router/blob/main/semantic_router/splitters/rolling_window.py)
