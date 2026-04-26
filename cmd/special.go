package cmd

import (
	"fmt"

	"github.com/defectdojo-cli/pkg"
	"github.com/spf13/cobra"
)

func registerSpecialCommands() {
	registerFindingActions()
	registerEngagementActions()
	registerUserActions()
	registerExtraCommands()
}

func registerFindingActions() {
	// `dd findings` is already registered as a resource command in crud.go.
	// These are subcommands that need to be added to the `findings` resource command.
	// But since resource commands are generated dynamically, we need to find it.
	// Instead, we add them with `finding` prefix directly.
	// Actually, let's add them as separate commands prefixed with "finding-".

	// dd finding-close <id>
	closeCmd := &cobra.Command{
		Use:   "finding-close <id>",
		Short: "Close a finding",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Post(fmt.Sprintf("findings/%s/close/", args[0]), map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("close finding: %w", err)
			}
			fmt.Println("Finding closed:")
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(closeCmd)

	// dd finding-verify <id>
	verifyCmd := &cobra.Command{
		Use:   "finding-verify <id>",
		Short: "Verify a finding",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Post(fmt.Sprintf("findings/%s/verify/", args[0]), map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("verify finding: %w", err)
			}
			fmt.Println("Finding verified:")
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(verifyCmd)

	// dd finding-notes <id> [--body "..."]
	notesCmd := &cobra.Command{
		Use:   "finding-notes <id>",
		Short: "List or add notes to a finding",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, _ := cmd.Flags().GetString("body")
			if body != "" {
				// Add note via POST
				data, err := client.Post(fmt.Sprintf("findings/%s/notes/", args[0]),
					map[string]string{"entry": body})
				if err != nil {
					return fmt.Errorf("add note: %w", err)
				}
				fmt.Println("Note added:")
				return pkg.PrintOutput(data, pkg.FormatJSON)
			}
			// List notes via GET
			data, err := client.Get(fmt.Sprintf("findings/%s/notes/", args[0]), nil)
			if err != nil {
				return fmt.Errorf("list notes: %w", err)
			}
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	notesCmd.Flags().String("body", "", "Note body text (omit to list notes)")
	rootCmd.AddCommand(notesCmd)

	// dd finding-tags <id> [--add tag1,tag2] [--remove tag1]
	tagsCmd := &cobra.Command{
		Use:   "finding-tags <id>",
		Short: "Manage finding tags",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addTags, _ := cmd.Flags().GetString("add")
			removeTags, _ := cmd.Flags().GetString("remove")

			if addTags != "" {
				data, err := client.Post(fmt.Sprintf("findings/%s/tags/", args[0]),
					map[string]string{"tags": addTags})
				if err != nil {
					return fmt.Errorf("add tags: %w", err)
				}
				fmt.Println("Tags added:")
				return pkg.PrintOutput(data, pkg.FormatJSON)
			}
			if removeTags != "" {
				data, err := client.Post(fmt.Sprintf("findings/%s/remove_tags/", args[0]),
					map[string]string{"tags": removeTags})
				if err != nil {
					return fmt.Errorf("remove tags: %w", err)
				}
				fmt.Println("Tags removed:")
				return pkg.PrintOutput(data, pkg.FormatJSON)
			}

			// List tags
			data, err := client.Get(fmt.Sprintf("findings/%s/", args[0]), nil)
			if err != nil {
				return err
			}
			var f map[string]interface{}
			if err := jsonUnmarshalImpl(data, &f); err != nil {
				return err
			}
			if tags, ok := f["tags"]; ok {
				fmt.Println("Tags:", tags)
			}
			return nil
		},
	}
	tagsCmd.Flags().String("add", "", "Tags to add (comma-separated)")
	tagsCmd.Flags().String("remove", "", "Tags to remove (comma-separated)")
	rootCmd.AddCommand(tagsCmd)

	// dd finding-duplicate <id>
	dupCmd := &cobra.Command{
		Use:   "finding-duplicate <id>",
		Short: "Get duplicate cluster for a finding",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Get(fmt.Sprintf("findings/%s/duplicate/", args[0]), nil)
			if err != nil {
				return fmt.Errorf("get duplicates: %w", err)
			}
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(dupCmd)

	// dd finding-metadata <id>
	metaCmd := &cobra.Command{
		Use:   "finding-metadata <id>",
		Short: "Get finding metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Get(fmt.Sprintf("findings/%s/metadata/", args[0]), nil)
			if err != nil {
				return fmt.Errorf("get metadata: %w", err)
			}
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(metaCmd)
}

func registerEngagementActions() {
	// dd engagement-close <id>
	closeCmd := &cobra.Command{
		Use:   "engagement-close <id>",
		Short: "Close an engagement",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Post(fmt.Sprintf("engagements/%s/close/", args[0]), map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("close engagement: %w", err)
			}
			fmt.Println("Engagement closed:")
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(closeCmd)

	// dd engagement-reopen <id>
	reopenCmd := &cobra.Command{
		Use:   "engagement-reopen <id>",
		Short: "Reopen an engagement",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Post(fmt.Sprintf("engagements/%s/reopen/", args[0]), map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("reopen engagement: %w", err)
			}
			fmt.Println("Engagement reopened:")
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(reopenCmd)

	// dd engagement-notes <id> [--body "..."]
	notesCmd := &cobra.Command{
		Use:   "engagement-notes <id>",
		Short: "List or add notes to an engagement",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, _ := cmd.Flags().GetString("body")
			if body != "" {
				data, err := client.Post(fmt.Sprintf("engagements/%s/notes/", args[0]),
					map[string]string{"entry": body})
				if err != nil {
					return fmt.Errorf("add note: %w", err)
				}
				fmt.Println("Note added:")
				return pkg.PrintOutput(data, pkg.FormatJSON)
			}
			data, err := client.Get(fmt.Sprintf("engagements/%s/notes/", args[0]), nil)
			if err != nil {
				return fmt.Errorf("list notes: %w", err)
			}
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	notesCmd.Flags().String("body", "", "Note body text (omit to list notes)")
	rootCmd.AddCommand(notesCmd)
}

func registerUserActions() {
	// dd user-reset-token <id>
	resetCmd := &cobra.Command{
		Use:   "user-reset-token <id>",
		Short: "Reset user API token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Post(fmt.Sprintf("users/%s/reset_api_token/", args[0]), map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("reset token: %w", err)
			}
			fmt.Println("Token reset:")
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(resetCmd)

	// dd user-profile
	profileCmd := &cobra.Command{
		Use:   "user-profile",
		Short: "Get current user profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Get("user_profile/", nil)
			if err != nil {
				return fmt.Errorf("get profile: %w", err)
			}
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(profileCmd)
}

func registerExtraCommands() {
	// dd engagement-checklist <id>
	checklistCmd := &cobra.Command{
		Use:   "engagement-checklist <id>",
		Short: "Get or update engagement checklist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			method := "GET"
			var body interface{}
			payload, _ := cmd.Flags().GetString("data")
			if payload != "" {
				method = "POST"
				var data map[string]interface{}
				if err := jsonUnmarshalImpl([]byte(payload), &data); err != nil {
					return fmt.Errorf("invalid JSON: %w", err)
				}
				body = data
			}
			if method == "POST" {
				data, err := client.Post(fmt.Sprintf("engagements/%s/complete_checklist/", args[0]), body)
				if err != nil {
					return fmt.Errorf("update checklist: %w", err)
				}
				return pkg.PrintOutput(data, pkg.FormatJSON)
			}
			data, err := client.Get(fmt.Sprintf("engagements/%s/complete_checklist/", args[0]), nil)
			if err != nil {
				return fmt.Errorf("get checklist: %w", err)
			}
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	checklistCmd.Flags().String("data", "", "Checklist JSON data (omit to retrieve)")
	rootCmd.AddCommand(checklistCmd)

	// dd import-languages <product-id> <file>
	importLangCmd := &cobra.Command{
		Use:   "import-languages <product-id> <file>",
		Short: "Import languages for a product",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.PostMultipart("import-languages/",
				map[string]string{"product": args[0]}, "file", args[1])
			if err != nil {
				return fmt.Errorf("import languages: %w", err)
			}
			fmt.Println("Languages imported:")
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(importLangCmd)

	// dd endpoint-meta-import
	endpointMetaCmd := &cobra.Command{
		Use:   "endpoint-meta-import",
		Short: "Import endpoint metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			productID, _ := cmd.Flags().GetInt("product-id")
			if productID == 0 {
				return fmt.Errorf("--product-id is required")
			}
			filePath, _ := cmd.Flags().GetString("file")
			if filePath == "" {
				return fmt.Errorf("--file is required")
			}
			createEndpoints, _ := cmd.Flags().GetBool("create-endpoints")
			fields := map[string]string{
				"product": fmt.Sprintf("%d", productID),
			}
			if createEndpoints {
				fields["create_endpoints"] = "true"
			}
			data, err := client.PostMultipart("endpoint_meta_import/", fields, "file", filePath)
			if err != nil {
				return fmt.Errorf("import endpoint meta: %w", err)
			}
			fmt.Println("Endpoint metadata imported:")
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	endpointMetaCmd.Flags().Int("product-id", 0, "Product ID (required)")
	endpointMetaCmd.Flags().String("file", "", "CSV file (required)")
	endpointMetaCmd.Flags().Bool("create-endpoints", false, "Create endpoints that don't exist")
	rootCmd.AddCommand(endpointMetaCmd)

	// dd celery <task-id>
	celeryCmd := &cobra.Command{
		Use:   "celery <task-id>",
		Short: "Check Celery task status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Get(fmt.Sprintf("celery/%s/", args[0]), nil)
			if err != nil {
				return fmt.Errorf("check celery task: %w", err)
			}
			return pkg.PrintOutput(data, pkg.FormatJSON)
		},
	}
	rootCmd.AddCommand(celeryCmd)
}
