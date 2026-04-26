package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

type Format string

const (
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatTable Format = "table"
)

type PrintOptions struct {
	Format  Format
	Columns []string // nil = auto-detect and show all
}

func PrintOutput(data []byte, format Format) error {
	return PrintWithOpts(data, PrintOptions{Format: format, Columns: nil})
}

func PrintWithOpts(data []byte, opts PrintOptions) error {
	switch opts.Format {
	case FormatJSON:
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, data, "", "  "); err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
		fmt.Println(pretty.String())
		return nil

	case FormatYAML:
		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
		yamlData, err := yaml.Marshal(jsonData)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
		fmt.Print(string(yamlData))
		return nil

	case FormatTable:
		return printTable(data, opts.Columns)

	default:
		fmt.Println(string(data))
		return nil
	}
}

func printTable(data []byte, columns []string) error {
	var results []map[string]interface{}

	var pr struct {
		Results []map[string]interface{} `json:"results"`
	}
	if err := json.Unmarshal(data, &pr); err == nil && len(pr.Results) > 0 {
		results = pr.Results
	} else {
		var single map[string]interface{}
		if err := json.Unmarshal(data, &single); err == nil {
			if len(single) > 0 {
				results = []map[string]interface{}{single}
			}
		}
	}

	if len(results) == 0 {
		// Try as array
		var arr []map[string]interface{}
		if err := json.Unmarshal(data, &arr); err == nil && len(arr) > 0 {
			results = arr
		}
	}

	if len(results) == 0 {
		fmt.Println("No results.")
		return nil
	}

	// Determine columns to display
	availableCols := collectColumns(results)
	cols := resolveColumns(columns, availableCols)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	header := strings.Join(cols, "\t")
	fmt.Fprintln(w, header)

	sep := make([]string, len(cols))
	for i, c := range cols {
		sep[i] = strings.Repeat("-", len(c))
	}
	fmt.Fprintln(w, strings.Join(sep, "\t"))

	for _, r := range results {
		row := make([]string, len(cols))
		for i, col := range cols {
			row[i] = formatValue(r[col])
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	return w.Flush()
}

func collectColumns(results []map[string]interface{}) []string {
	seen := map[string]bool{}
	var cols []string
	for _, r := range results {
		for k := range r {
			if !seen[k] {
				seen[k] = true
				cols = append(cols, k)
			}
		}
	}

	priority := []string{"id", "name", "title", "severity", "status", "description", "url", "product", "engagement", "test", "scan_type", "created", "updated", "active", "verified", "mitigated", "found_by", "test_type", "environment", "branch_tag", "host", "port", "protocol", "username", "email", "is_active", "last_login", "prod_type", "lifecycle", "platform", "findings_count", "tags"}
	inserted := map[string]bool{}
	var ordered []string
	for _, p := range priority {
		if seen[p] {
			ordered = append(ordered, p)
			inserted[p] = true
		}
	}
	for _, c := range cols {
		if !inserted[c] {
			ordered = append(ordered, c)
		}
	}
	return ordered
}

func resolveColumns(requested, available []string) []string {
	if len(requested) == 0 {
		return available
	}

	availSet := make(map[string]bool, len(available))
	for _, c := range available {
		availSet[c] = true
	}

	var result []string
	for _, c := range requested {
		if availSet[c] {
			result = append(result, c)
		}
	}
	if len(result) == 0 {
		return available
	}
	return result
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		if len(val) > 80 {
			return val[:77] + "..."
		}
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%.2f", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case map[string]interface{}:
		if id, ok := val["id"]; ok {
			return fmt.Sprintf("#%v", id)
		}
		if name, ok := val["name"]; ok {
			return fmt.Sprintf("%v", name)
		}
		b, _ := json.Marshal(val)
		s := string(b)
		if len(s) > 40 {
			return s[:37] + "..."
		}
		return s
	case []interface{}:
		if len(val) == 0 {
			return ""
		}
		b, _ := json.Marshal(val)
		s := string(b)
		if len(s) > 40 {
			return s[:37] + "..."
		}
		return s
	default:
		return fmt.Sprintf("%v", v)
	}
}
