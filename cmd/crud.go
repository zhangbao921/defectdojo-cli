package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/defectdojo-cli/pkg"
	"github.com/spf13/cobra"
)

func registerResourceCommands() {
	for _, r := range allResources {
		registerResource(r)
	}
}

func registerResource(r ResourceDef) {
	resourceCmd := &cobra.Command{
		Use:     r.Name,
		Short:   fmt.Sprintf("Manage %s", displayName(r)),
		Long:    fmt.Sprintf("List, get, create, update, delete %s.", displayName(r)),
		Aliases: []string{strings.ReplaceAll(r.Name, "-", "_")},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: fmt.Sprintf("List %s", displayName(r)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listResources(r, cmd)
		},
	}
	listCmd.Flags().StringArrayP("filter", "f", nil, "Filter by field (e.g. --filter name=MyApp)")
	listCmd.Flags().Int("limit", 0, "Max results (0 = all)")
	listCmd.Flags().Int("page", 1, "Page number (when not listing all)")
	listCmd.Flags().String("columns", "", "Comma-separated columns to show (default: auto)")
	resourceCmd.AddCommand(listCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: fmt.Sprintf("Get a %s by ID", r.Singular),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getResource(r, args[0], cmd)
		},
	}
	resourceCmd.AddCommand(getCmd)

	createCmd := &cobra.Command{
		Use:   "create",
		Short: fmt.Sprintf("Create a %s", r.Singular),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createResource(r, cmd)
		},
	}
	createCmd.Flags().String("data", "", "JSON data (or use --data-file)")
	createCmd.Flags().String("data-file", "", "Read JSON data from file")
	createCmd.MarkFlagsOneRequired("data", "data-file")
	resourceCmd.AddCommand(createCmd)

	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: fmt.Sprintf("Update a %s (full update, PUT)", r.Singular),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateResource(r, args[0], cmd)
		},
	}
	updateCmd.Flags().String("data", "", "JSON data (or use --data-file)")
	updateCmd.Flags().String("data-file", "", "Read JSON data from file")
	updateCmd.MarkFlagsOneRequired("data", "data-file")
	resourceCmd.AddCommand(updateCmd)

	patchCmd := &cobra.Command{
		Use:   "patch <id>",
		Short: fmt.Sprintf("Partially update a %s (PATCH)", r.Singular),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return patchResource(r, args[0], cmd)
		},
	}
	patchCmd.Flags().String("data", "", "JSON data (or use --data-file)")
	patchCmd.Flags().String("data-file", "", "Read JSON data from file")
	patchCmd.MarkFlagsOneRequired("data", "data-file")
	resourceCmd.AddCommand(patchCmd)

	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: fmt.Sprintf("Delete a %s", r.Singular),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteResource(r, args[0])
		},
	}
	resourceCmd.AddCommand(deleteCmd)

	rootCmd.AddCommand(resourceCmd)
}

func getFilters(cmd *cobra.Command) url.Values {
	filters, _ := cmd.Flags().GetStringArray("filter")
	params := url.Values{}
	for _, f := range filters {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) == 2 {
			params.Set(parts[0], parts[1])
		}
	}
	return params
}

func readData(cmd *cobra.Command) ([]byte, error) {
	dataStr, _ := cmd.Flags().GetString("data")
	dataFile, _ := cmd.Flags().GetString("data-file")
	if dataFile != "" {
		return readFile(dataFile)
	}
	return []byte(dataStr), nil
}

func listResources(r ResourceDef, cmd *cobra.Command) error {
	params := getFilters(cmd)
	limit, _ := cmd.Flags().GetInt("limit")
	page, _ := cmd.Flags().GetInt("page")
	columnsStr, _ := cmd.Flags().GetString("columns")

	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if page > 1 {
		params.Set("page", fmt.Sprintf("%d", page))
	}

	var data []byte
	var err error

	if r.ListAll && limit == 0 {
		items, err := client.ListAll(r.Path, params)
		if err != nil {
			return err
		}
		wrapper := map[string]interface{}{"results": items, "count": len(items)}
		data, err = json.Marshal(wrapper)
		if err != nil {
			return err
		}
	} else {
		data, err = client.Get(r.Path, params)
		if err != nil {
			return err
		}
	}

	outFmt, _ := getOutputFormat(cmd)

	// Resolve columns
	var cols []string
	if columnsStr != "" {
		cols = strings.Split(columnsStr, ",")
		for i := range cols {
			cols[i] = strings.TrimSpace(cols[i])
		}
	} else if len(r.Columns) > 0 && outFmt == pkg.FormatTable {
		cols = r.Columns
	}

	return pkg.PrintWithOpts(data, pkg.PrintOptions{
		Format:  outFmt,
		Columns: cols,
	})
}

func getResource(r ResourceDef, id string, cmd *cobra.Command) error {
	data, err := client.GetByID(r.Path, id)
	if err != nil {
		return err
	}
	outFmt, _ := getOutputFormat(cmd)
	return pkg.PrintOutput(data, outFmt)
}

func createResource(r ResourceDef, cmd *cobra.Command) error {
	rawData, err := readData(cmd)
	if err != nil {
		return fmt.Errorf("reading data: %w", err)
	}
	var data interface{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	result, err := client.Post(r.Path, data)
	if err != nil {
		return err
	}
	fmt.Println("Created:")
	return pkg.PrintOutput(result, pkg.FormatJSON)
}

func updateResource(r ResourceDef, id string, cmd *cobra.Command) error {
	rawData, err := readData(cmd)
	if err != nil {
		return fmt.Errorf("reading data: %w", err)
	}
	var data interface{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	result, err := client.Put(r.Path, id, data)
	if err != nil {
		return err
	}
	fmt.Println("Updated:")
	return pkg.PrintOutput(result, pkg.FormatJSON)
}

func patchResource(r ResourceDef, id string, cmd *cobra.Command) error {
	rawData, err := readData(cmd)
	if err != nil {
		return fmt.Errorf("reading data: %w", err)
	}
	var data interface{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	result, err := client.Patch(r.Path, id, data)
	if err != nil {
		return err
	}
	fmt.Println("Patched:")
	return pkg.PrintOutput(result, pkg.FormatJSON)
}

func deleteResource(r ResourceDef, id string) error {
	if err := client.Delete(r.Path, id); err != nil {
		return err
	}
	fmt.Printf("Deleted %s #%s\n", r.Singular, id)
	return nil
}

func displayName(r ResourceDef) string {
	switch r.Singular {
	case "endpoint status", "metadata":
		return r.Singular
	}
	if strings.HasSuffix(r.Singular, "analysis") {
		return strings.TrimSuffix(r.Singular, "analysis") + "analyses"
	}
	return r.Singular + "s"
}

func getOutputFormat(cmd *cobra.Command) (pkg.Format, error) {
	val, _ := cmd.Flags().GetString("output")
	if val == "" {
		val = rootCmd.PersistentFlags().Lookup("output").Value.String()
	}
	switch strings.ToLower(val) {
	case "json":
		return pkg.FormatJSON, nil
	case "yaml":
		return pkg.FormatYAML, nil
	case "table":
		return pkg.FormatTable, nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", val)
	}
}

func readFile(path string) ([]byte, error) {
	return readFileImpl(path)
}
