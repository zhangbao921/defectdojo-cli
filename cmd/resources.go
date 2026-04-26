package cmd

type ResourceDef struct {
	Name     string   // CLI command name (e.g. "product-types")
	Path     string   // API path (e.g. "product_types")
	Singular string   // Description for help text
	ListAll  bool     // Use paginated list-all by default
	Columns  []string // Default table columns (empty = auto)
}

var allResources = []ResourceDef{
	{
		Name: "product-types", Path: "product_types", Singular: "product type", ListAll: true,
		Columns: []string{"id", "name", "description", "critical_product", "key_product"},
	},
	{
		Name: "products", Path: "products", Singular: "product", ListAll: true,
		Columns: []string{"id", "name", "description", "prod_type", "lifecycle", "platform", "findings_count", "tags"},
	},
	{
		Name: "engagements", Path: "engagements", Singular: "engagement", ListAll: true,
		Columns: []string{"id", "name", "status", "description", "product", "target_start", "target_end", "lead"},
	},
	{
		Name: "tests", Path: "tests", Singular: "test", ListAll: true,
		Columns: []string{"id", "title", "test_type", "engagement", "environment", "branch_tag", "created"},
	},
	{
		Name: "test-types", Path: "test_types", Singular: "test type", ListAll: true,
		Columns: []string{"id", "name", "static_tool", "dynamic_tool", "active"},
	},
	{
		Name: "findings", Path: "findings", Singular: "finding", ListAll: true,
		Columns: []string{"id", "title", "severity", "status", "test", "found_by", "date", "mitigated", "cwe", "cvssv3_score"},
	},
	{
		Name: "finding-templates", Path: "finding_templates", Singular: "finding template", ListAll: true,
		Columns: []string{"id", "title", "severity", "description", "cwe"},
	},
	{
		Name: "endpoints", Path: "endpoints", Singular: "endpoint", ListAll: true,
		Columns: []string{"id", "host", "port", "protocol", "product", "findings_count"},
	},
	{
		Name: "endpoint-status", Path: "endpoint_status", Singular: "endpoint status", ListAll: true,
		Columns: []string{"id", "endpoint", "finding", "mitigated", "false_positive", "out_of_scope", "risk_accepted"},
	},
	{
		Name: "users", Path: "users", Singular: "user", ListAll: true,
		Columns: []string{"id", "username", "email", "first_name", "last_name", "is_active", "is_superuser", "last_login"},
	},
	{
		Name: "development-environments", Path: "development_environments", Singular: "development environment", ListAll: true,
		Columns: []string{"id", "name"},
	},
	{Name: "notes", Path: "notes", Singular: "note", ListAll: false},
	{Name: "metadata", Path: "metadata", Singular: "metadata", ListAll: false},
	{
		Name: "risk-acceptance", Path: "risk_acceptance", Singular: "risk acceptance", ListAll: true,
		Columns: []string{"id", "name", "recommendation", "decision", "accepted_by", "expiration_date"},
	},
	{
		Name: "jira-instances", Path: "jira_instances", Singular: "JIRA instance", ListAll: true,
		Columns: []string{"id", "url", "username"},
	},
	{
		Name: "jira-projects", Path: "jira_projects", Singular: "JIRA project", ListAll: true,
		Columns: []string{"id", "product", "engagement", "project_key", "enabled", "push_all_issues"},
	},
	{
		Name: "jira-finding-mappings", Path: "jira_finding_mappings", Singular: "JIRA finding mapping", ListAll: true,
		Columns: []string{"id", "jira_id", "jira_key", "finding", "engagement"},
	},
	{
		Name: "sonarqube-issues", Path: "sonarqube_issues", Singular: "SonarQube issue", ListAll: true,
		Columns: []string{"id", "key", "status", "type"},
	},
	{
		Name: "sonarqube-transitions", Path: "sonarqube_transitions", Singular: "SonarQube transition", ListAll: true,
	},
	{
		Name: "tool-configurations", Path: "tool_configurations", Singular: "tool configuration", ListAll: true,
		Columns: []string{"id", "name", "tool_type", "url", "authentication_type"},
	},
	{
		Name: "tool-product-settings", Path: "tool_product_settings", Singular: "tool product setting", ListAll: true,
		Columns: []string{"id", "name", "product", "tool_configuration", "tool_project_id"},
	},
	{
		Name: "tool-types", Path: "tool_types", Singular: "tool type", ListAll: true,
		Columns: []string{"id", "name", "description"},
	},
	{
		Name: "regulations", Path: "regulations", Singular: "regulation", ListAll: true,
		Columns: []string{"id", "name", "description"},
	},
	{Name: "roles", Path: "roles", Singular: "role", ListAll: true, Columns: []string{"id", "name"}},
	{
		Name: "dojo-groups", Path: "dojo_groups", Singular: "dojo group", ListAll: true,
		Columns: []string{"id", "name", "social_provider"},
	},
	{
		Name: "dojo-group-members", Path: "dojo_group_members", Singular: "dojo group member", ListAll: true,
		Columns: []string{"id", "group_id", "user_id", "role"},
	},
	{
		Name: "global-roles", Path: "global_roles", Singular: "global role", ListAll: true,
		Columns: []string{"id", "user", "group", "role"},
	},
	{
		Name: "product-members", Path: "product_members", Singular: "product member", ListAll: true,
		Columns: []string{"id", "product_id", "user_id", "role"},
	},
	{
		Name: "product-groups", Path: "product_groups", Singular: "product group", ListAll: true,
		Columns: []string{"id", "product_id", "group_id", "role"},
	},
	{
		Name: "product-type-members", Path: "product_type_members", Singular: "product type member", ListAll: true,
		Columns: []string{"id", "product_type_id", "user_id", "role"},
	},
	{
		Name: "product-type-groups", Path: "product_type_groups", Singular: "product type group", ListAll: true,
		Columns: []string{"id", "product_type_id", "group_id", "role"},
	},
	{
		Name: "product-api-scan-configs", Path: "product_api_scan_configurations", Singular: "API scan config", ListAll: true,
		Columns: []string{"id", "product", "tool_configuration", "service_key_1"},
	},
	{
		Name: "app-analysis", Path: "technologies", Singular: "application analysis", ListAll: true,
		Columns: []string{"id", "product", "name", "version"},
	},
	{Name: "credentials", Path: "credentials", Singular: "credential", ListAll: true},
	{Name: "credential-mappings", Path: "credential_mappings", Singular: "credential mapping", ListAll: true},
	{Name: "stub-findings", Path: "stub_findings", Singular: "stub finding", ListAll: true},
	{Name: "system-settings", Path: "system_settings", Singular: "system setting", ListAll: false},
	{Name: "announcements", Path: "announcements", Singular: "announcement", ListAll: true},
	{Name: "notifications", Path: "notifications", Singular: "notification", ListAll: true},
	{Name: "notification-webhooks", Path: "notification_webhooks", Singular: "notification webhook", ListAll: true},
	{Name: "languages", Path: "languages", Singular: "language", ListAll: true},
	{Name: "language-types", Path: "language_types", Singular: "language type", ListAll: true},
	{
		Name: "sla-configurations", Path: "sla_configurations", Singular: "SLA configuration", ListAll: true,
		Columns: []string{"id", "name", "critical", "high", "medium", "low"},
	},
	{Name: "network-locations", Path: "network_locations", Singular: "network location", ListAll: true},
	{Name: "engagement-presets", Path: "engagement_presets", Singular: "engagement preset", ListAll: true},
	{Name: "test-imports", Path: "test_imports", Singular: "test import", ListAll: true},
	{Name: "configuration-permissions", Path: "configuration_permissions", Singular: "configuration permission", ListAll: true},
	{Name: "request-response-pairs", Path: "request_response_pairs", Singular: "request/response pair", ListAll: true},
	{Name: "user-contact-infos", Path: "user_contact_infos", Singular: "user contact info", ListAll: true},
	{Name: "questionnaire-answers", Path: "questionnaire_answers", Singular: "questionnaire answer", ListAll: true},
	{Name: "questionnaire-answered-questionnaires", Path: "questionnaire_answered_questionnaires", Singular: "answered questionnaire", ListAll: true},
	{Name: "questionnaire-engagement-questionnaires", Path: "questionnaire_engagement_questionnaires", Singular: "engagement questionnaire", ListAll: true},
	{Name: "questionnaire-general-questionnaires", Path: "questionnaire_general_questionnaires", Singular: "general questionnaire", ListAll: true},
	{Name: "questionnaire-questions", Path: "questionnaire_questions", Singular: "questionnaire question", ListAll: true},
}
