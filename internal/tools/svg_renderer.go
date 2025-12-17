// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/goccy/go-graphviz"
)

// dotToSVG converts a DOT graph string to SVG format using graphviz.
func dotToSVG(ctx context.Context, dot string) (string, error) {
	g, err := graphviz.New(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create graphviz instance: %w", err)
	}
	defer func() {
		if closeErr := g.Close(); closeErr != nil {
			// If we're returning an error, combine them; otherwise just log/ignore
			// Since we're in a defer, we can't easily return this error
			// so we combine it with any existing error via a sentinel approach
			err = errors.Join(err, closeErr)
		}
	}()

	graph, err := graphviz.ParseBytes([]byte(dot))
	if err != nil {
		return "", fmt.Errorf("failed to parse DOT graph: %w", err)
	}
	defer func() {
		if closeErr := graph.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, graphviz.SVG, &buf); err != nil {
		return "", fmt.Errorf("failed to render SVG: %w", err)
	}

	return buf.String(), nil
}
