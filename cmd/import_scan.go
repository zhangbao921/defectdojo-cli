package cmd

import (
	"fmt"

	"github.com/defectdojo-cli/pkg"
	"github.com/spf13/cobra"
)

func registerScanCommands() {
	registerImportScan()
	registerReimportScan()
	registerQuickScan()
}

func registerImportScan() {
	cmd := &cobra.Command{
		Use:   "import-scan",
		Short: "Import a scan report into an engagement",
		Long: `Import a scan report into an engagement.

Examples:
  dd import-scan --engagement-id 1 --scan-type "ZAP Scan" --file results.json
  dd import-scan --engagement-id 1 --scan-type "Nmap Scan" --file scan.xml
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImportScan(cmd)
		},
	}
	cmd.Flags().Int("engagement-id", 0, "Engagement ID (required)")
	cmd.Flags().Int("product-id", 0, "Product ID")
	cmd.Flags().Int("test-id", 0, "Test ID to append findings to")
	cmd.Flags().String("scan-type", "", "Scan type (required)")
	cmd.Flags().String("file", "", "Scan report file path (required)")
	cmd.Flags().String("file-type", "", "Report file type extension")
	cmd.Flags().Bool("active", true, "Mark findings as active")
	cmd.Flags().Bool("verified", false, "Mark findings as verified")
	cmd.Flags().Int("minimum-severity", 0, "Minimum severity (0=all, 1=Info, 2=Low, 3=Medium, 4=High, 5=Critical)")
	cmd.Flags().String("version", "", "Product version")
	cmd.Flags().Bool("close-old-findings", false, "Close old findings")
	cmd.Flags().String("source-code-management-uri", "", "SCM URI")
	cmd.Flags().Int("push-to-jira", 0, "Push to JIRA (0=no, 1=yes)")
	cmd.Flags().String("tags", "", "Comma-separated tags")
	cmd.MarkFlagRequired("engagement-id")
	cmd.MarkFlagRequired("scan-type")
	cmd.MarkFlagRequired("file")
	rootCmd.AddCommand(cmd)
}

func registerReimportScan() {
	cmd := &cobra.Command{
		Use:   "reimport-scan",
		Short: "Re-import a scan report to update an existing test",
		Long: `Re-import a scan report to update an existing test.

Examples:
  dd reimport-scan --test-id 1 --scan-type "ZAP Scan" --file results.json
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReimportScan(cmd)
		},
	}
	cmd.Flags().Int("test-id", 0, "Test ID (required)")
	cmd.Flags().String("scan-type", "", "Scan type (required)")
	cmd.Flags().String("file", "", "Scan report file path (required)")
	cmd.Flags().String("file-type", "", "Report file type extension")
	cmd.Flags().Bool("active", true, "Mark findings as active")
	cmd.Flags().Bool("verified", false, "Mark findings as verified")
	cmd.Flags().Int("minimum-severity", 0, "Minimum severity to import")
	cmd.Flags().String("version", "", "Product version")
	cmd.Flags().Bool("close-old-findings", true, "Close old findings that are no longer present")
	cmd.Flags().String("tags", "", "Comma-separated tags")
	cmd.MarkFlagRequired("test-id")
	cmd.MarkFlagRequired("scan-type")
	cmd.MarkFlagRequired("file")
	rootCmd.AddCommand(cmd)
}

func registerQuickScan() {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Quick scan: auto-create product + engagement + import scan",
		Long: `Quickly scan by auto-creating product and engagement if needed.

Examples:
  dd scan --product "MyApp" --engagement "Sprint 42" --scan-type "ZAP Scan" --file results.json
  dd scan --product-id 1 --engagement-id 2 --scan-type "ZAP Scan" --file results.json
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuickScan(cmd)
		},
	}
	cmd.Flags().String("product", "", "Product name (created if doesn't exist)")
	cmd.Flags().String("engagement", "", "Engagement name (created if doesn't exist)")
	cmd.Flags().String("scan-type", "", "Scan type (required)")
	cmd.Flags().String("file", "", "Scan report file path (required)")
	cmd.Flags().Int("product-id", 0, "Existing product ID")
	cmd.Flags().Int("engagement-id", 0, "Existing engagement ID")
	cmd.MarkFlagsOneRequired("product", "product-id")
	cmd.MarkFlagsOneRequired("engagement", "engagement-id")
	cmd.MarkFlagRequired("scan-type")
	cmd.MarkFlagRequired("file")
	rootCmd.AddCommand(cmd)
}

func runImportScan(cmd *cobra.Command) error {
	fields, filePath := collectScanFields(cmd, "import")
	data, err := client.PostMultipart("import-scan/", fields, "file", filePath)
	if err != nil {
		return fmt.Errorf("import scan failed: %w", err)
	}
	fmt.Println("Scan imported successfully:")
	return pkg.PrintOutput(data, pkg.FormatJSON)
}

func runReimportScan(cmd *cobra.Command) error {
	fields, filePath := collectScanFields(cmd, "reimport")
	data, err := client.PostMultipart("reimport-scan/", fields, "file", filePath)
	if err != nil {
		return fmt.Errorf("reimport scan failed: %w", err)
	}
	fmt.Println("Scan reimported successfully:")
	return pkg.PrintOutput(data, pkg.FormatJSON)
}

func collectScanFields(cmd *cobra.Command, mode string) (map[string]string, string) {
	f := map[string]string{}

	if mode == "import" {
		f["engagement_id"] = flagStr(cmd, "engagement-id")
		if cmd.Flags().Changed("product-id") {
			f["product_id"] = flagStr(cmd, "product-id")
		}
		if cmd.Flags().Changed("test-id") {
			f["test_id"] = flagStr(cmd, "test-id")
		}
	} else {
		f["test_id"] = flagStr(cmd, "test-id")
	}

	f["scan_type"] = flagStr(cmd, "scan-type")
	filePath := flagStr(cmd, "file")

	if cmd.Flags().Changed("file-type") {
		f["file_type"] = flagStr(cmd, "file-type")
	}
	f["active"] = flagBoolStr(cmd, "active")
	f["verified"] = flagBoolStr(cmd, "verified")
	if cmd.Flags().Changed("minimum-severity") {
		f["minimum_severity"] = flagStr(cmd, "minimum-severity")
	}
	if cmd.Flags().Changed("version") {
		f["version"] = flagStr(cmd, "version")
	}
	f["close_old_findings"] = flagBoolStr(cmd, "close-old-findings")
	if cmd.Flags().Changed("source-code-management-uri") {
		f["source_code_management_uri"] = flagStr(cmd, "source-code-management-uri")
	}
	if cmd.Flags().Changed("push-to-jira") {
		f["push_to_jira"] = flagStr(cmd, "push-to-jira")
	}
	if cmd.Flags().Changed("tags") {
		f["tags"] = flagStr(cmd, "tags")
	}

	return f, filePath
}

func runQuickScan(cmd *cobra.Command) error {
	scanType, _ := cmd.Flags().GetString("scan-type")
	filePath, _ := cmd.Flags().GetString("file")

	fields := map[string]string{
		"scan_type": scanType,
		"file":      filePath,
		"active":    "true",
	}

	productID, _ := cmd.Flags().GetInt("product-id")
	if productID == 0 {
		productName, _ := cmd.Flags().GetString("product")
		prodData, err := client.Post("products/", map[string]string{
			"name": productName, "description": productName,
		})
		if err != nil {
			return fmt.Errorf("create product: %w", err)
		}
		var prod map[string]interface{}
		if err := jsonUnmarshalImpl(prodData, &prod); err != nil {
			return err
		}
		productID = int(prod["id"].(float64))
		fmt.Printf("Created product #%d\n", productID)
	}
	fields["product_id"] = fmt.Sprintf("%d", productID)

	engagementID, _ := cmd.Flags().GetInt("engagement-id")
	if engagementID == 0 {
		engagementName, _ := cmd.Flags().GetString("engagement")
		engData, err := client.Post("engagements/", map[string]interface{}{
			"name":        engagementName,
			"product":     productID,
			"target_start": currentDate(),
			"target_end":   currentDate(),
			"status":      "In Progress",
		})
		if err != nil {
			return fmt.Errorf("create engagement: %w", err)
		}
		var eng map[string]interface{}
		if err := jsonUnmarshalImpl(engData, &eng); err != nil {
			return err
		}
		engagementID = int(eng["id"].(float64))
		fmt.Printf("Created engagement #%d\n", engagementID)
	}
	fields["engagement_id"] = fmt.Sprintf("%d", engagementID)

	data, err := client.PostMultipart("import-scan/", fields, "file", filePath)
	if err != nil {
		return fmt.Errorf("import scan failed: %w", err)
	}
	fmt.Println("Scan imported successfully:")
	return pkg.PrintOutput(data, pkg.FormatJSON)
}

func flagStr(cmd *cobra.Command, name string) string {
	if !cmd.Flags().Changed(name) {
		return ""
	}
	f := cmd.Flags().Lookup(name)
	if f == nil {
		return ""
	}
	switch f.Value.Type() {
	case "int":
		v, _ := cmd.Flags().GetInt(name)
		return fmt.Sprintf("%d", v)
	case "string":
		v, _ := cmd.Flags().GetString(name)
		return v
	default:
		return f.Value.String()
	}
}

func flagBoolStr(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetBool(name)
	if v {
		return "true"
	}
	return "false"
}
