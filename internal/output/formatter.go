package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"golang.org/x/term"
)

type LogFormatter struct {
	terminalWidth int
}

func NewFormatter() *LogFormatter {
	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		width = w
	}
	return &LogFormatter{
		terminalWidth: width,
	}
}

func (f *LogFormatter) PrintLogs(response *datadogV2.LogsListResponse) {
	if response == nil || response.Data == nil || len(response.Data) == 0 {
		return
	}

	for _, log := range response.Data {
		f.printSingleLog(log)
	}
}

func (f *LogFormatter) printSingleLog(log datadogV2.Log) {
	if log.Attributes == nil || log.Attributes.Message == nil {
		return
	}

	var parts []string

	if log.Attributes.Timestamp != nil {
		parts = append(parts, log.Attributes.Timestamp.Format(time.RFC3339))
	} else {
		parts = append(parts, time.Now().Format(time.RFC3339))
	}

	if log.Attributes.Service != nil && *log.Attributes.Service != "" {
		parts = append(parts, *log.Attributes.Service)
	}

	if log.Attributes.Status != nil && *log.Attributes.Status != "" {
		parts = append(parts, *log.Attributes.Status)
	}

	message := f.processMessage(*log.Attributes.Message)
	parts = append(parts, message)

	fmt.Println(f.formatLogLine(parts))
}

func (f *LogFormatter) processMessage(message string) string {
	if strings.HasPrefix(message, "[") {
		if idx := strings.Index(message, "]"); idx > 0 && idx < len(message)-1 {
			message = strings.TrimSpace(message[idx+1:])
		}
	}

	if strings.HasPrefix(message, "{") || strings.HasPrefix(message, "[") {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(message), &jsonData); err == nil {
			return f.extractFromJSON(jsonData)
		}
	}

	message = strings.TrimSpace(message)
	return message
}

func (f *LogFormatter) extractFromJSON(data map[string]interface{}) string {
	var parts []string

	fields := []string{"module", "method", "info", "error", "msg", "message", "action", "status"}
	
	for _, field := range fields {
		if val, ok := data[field]; ok && val != nil {
			strVal := fmt.Sprintf("%v", val)
			if strVal != "" && strVal != "<nil>" {
				parts = append(parts, fmt.Sprintf("%s=%s", field, strVal))
			}
		}
	}

	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}

	return fmt.Sprintf("%v", data)
}

func (f *LogFormatter) formatLogLine(parts []string) string {
	line := strings.Join(parts, " ")
	
	maxWidth := max(f.terminalWidth - 5, 40)
	
	if len(line) > maxWidth {
		line = line[:maxWidth] + "..."
	}
	
	return line
}
