package client

import (
	"context"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

type Config struct {
	APIKey     string
	AppKey     string
	Site       string
	Timeout    time.Duration
	MaxRetries int
}

type Client struct {
	api    *datadogV2.LogsApi
	config *Config
}

func NewClient(cfg *Config) *Client {
	if cfg.Site == "" {
		cfg.Site = "datadoghq.com"
	}

	configuration := datadog.NewConfiguration()
	// TODO: Properly set timeout when we understand the configuration better

	if configuration.DefaultHeader == nil {
		configuration.DefaultHeader = make(map[string]string)
	}
	configuration.DefaultHeader["DD-API-KEY"] = cfg.APIKey
	configuration.DefaultHeader["DD-APPLICATION-KEY"] = cfg.AppKey

	configuration.Servers = datadog.ServerConfigurations{
		{
			URL: "https://api." + cfg.Site,
		},
	}

	logsApi := datadogV2.NewLogsApi(datadog.NewAPIClient(configuration))

	return &Client{
		api:    logsApi,
		config: cfg,
	}
}

func (c *Client) FetchLogs(ctx context.Context, query string, cursor ...string) (*datadogV2.LogsListResponse, error) {
	params := datadogV2.NewListLogsGetOptionalParameters()

	if query != "" {
		params = params.WithFilterQuery(query)
	}

	if len(cursor) > 0 && cursor[0] != "" {
		params = params.WithPageCursor(cursor[0])
	} else {
		now := time.Now()
		from := now.Add(-15 * time.Minute)
		params = params.WithFilterFrom(from)
		params = params.WithFilterTo(now)
	}

	params = params.WithPageLimit(100)
	
	params = params.WithSort(datadogV2.LOGSSORT_TIMESTAMP_ASCENDING)

	response, _, err := c.api.ListLogsGet(ctx, *params)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
