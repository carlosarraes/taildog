package tailing

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/carlosarraes/taildog/internal/client"
	"github.com/carlosarraes/taildog/internal/output"
)

type Tailer struct {
	client    *client.Client
	formatter *output.LogFormatter
	interval  time.Duration
}

func NewTailer(client *client.Client, interval time.Duration) *Tailer {
	return &Tailer{
		client:    client,
		formatter: output.NewFormatter(),
		interval:  interval,
	}
}

func (t *Tailer) Start(ctx context.Context, query string) error {
	var cursor string
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	response, err := t.client.FetchLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("error fetching logs: %w", err)
	}
	t.formatter.PrintLogs(response)
	cursor = extractCursor(response)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			response, err := t.client.FetchLogs(ctx, query, cursor)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching logs: %v\n", err)
				continue
			}
			t.formatter.PrintLogs(response)
			cursor = extractCursor(response)
		}
	}
}

func extractCursor(response *datadogV2.LogsListResponse) string {
	if response == nil || response.Meta == nil || response.Meta.Page == nil || response.Meta.Page.After == nil {
		return ""
	}
	return *response.Meta.Page.After
}
