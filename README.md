# DefectDojo CLI

A command-line client for [DefectDojo](https://github.com/DefectDojo/django-DefectDojo) API v2. Manage products, engagements, tests, findings, and more — all from the terminal.

## Quick Start

```bash
# Set credentials
export DD_HOST=https://demo.defectdojo.org
export DD_API_KEY=your-api-key

# List resources
dd products list
dd findings list --limit 20
dd engagements list

# Get single resource
dd findings get 42 -o json
```

## Quick Reference

```
dd --help
```

```
A command-line client for DefectDojo API v2.
Manage products, engagements, tests, findings, and more.

Usage:
  dd [command]

Available Commands:
  announcements                           Manage announcements
  app-analysis                            Manage application analyses
  celery                                  Check Celery task status
  completion                              Generate the autocompletion script for the specified shell
  configuration-permissions               Manage configuration permissions
  credential-mappings                     Manage credential mappings
  credentials                             Manage credentials
  development-environments                Manage development environments
  dojo-group-members                      Manage dojo group members
  dojo-groups                             Manage dojo groups
  endpoint-meta-import                    Import endpoint metadata
  endpoint-status                         Manage endpoint status
  endpoints                               Manage endpoints
  engagement-checklist                    Get or update engagement checklist
  engagement-close                        Close an engagement
  engagement-notes                        List or add notes to an engagement
  engagement-presets                      Manage engagement presets
  engagement-reopen                       Reopen an engagement
  engagements                             Manage engagements
  finding-close                           Close a finding
  finding-duplicate                       Get duplicate cluster for a finding
  finding-metadata                        Get finding metadata
  finding-notes                           List or add notes to a finding
  finding-tags                            Manage finding tags
  finding-templates                       Manage finding templates
  finding-verify                          Verify a finding
  findings                                Manage findings
  global-roles                            Manage global roles
  help                                    Help about any command
  import-languages                        Import languages for a product
  import-scan                             Import a scan report into an engagement
  jira-finding-mappings                   Manage JIRA finding mappings
  jira-instances                          Manage JIRA instances
  jira-projects                           Manage JIRA projects
  language-types                          Manage language types
  languages                               Manage languages
  metadata                                Manage metadata
  network-locations                       Manage network locations
  notes                                   Manage notes
  notification-webhooks                   Manage notification webhooks
  notifications                           Manage notifications
  product-api-scan-configs                Manage API scan configs
  product-groups                          Manage product groups
  product-members                         Manage product members
  product-type-groups                     Manage product type groups
  product-type-members                    Manage product type members
  product-types                           Manage product types
  products                                Manage products
  questionnaire-answered-questionnaires   Manage answered questionnaires
  questionnaire-answers                   Manage questionnaire answers
  questionnaire-engagement-questionnaires Manage engagement questionnaires
  questionnaire-general-questionnaires    Manage general questionnaires
  questionnaire-questions                 Manage questionnaire questions
  regulations                             Manage regulations
  reimport-scan                           Re-import a scan report to update an existing test
  request-response-pairs                  Manage request/response pairs
  risk-acceptance                         Manage risk acceptances
  roles                                   Manage roles
  scan                                    Quick scan: auto-create product + engagement + import scan
  sla-configurations                      Manage SLA configurations
  sonarqube-issues                        Manage SonarQube issues
  sonarqube-transitions                   Manage SonarQube transitions
  stub-findings                           Manage stub findings
  system-settings                         Manage system settings
  test-imports                            Manage test imports
  test-types                              Manage test types
  tests                                   Manage tests
  tool-configurations                     Manage tool configurations
  tool-product-settings                   Manage tool product settings
  tool-types                              Manage tool types
  user-contact-infos                      Manage user contact infos
  user-profile                            Get current user profile
  user-reset-token                        Reset user API token
  users                                   Manage users

Flags:
      --api-key string   API key (env: DD_API_KEY)
  -h, --help             help for dd
      --host string      DefectDojo host URL (env: DD_HOST)
  -o, --output string    Output format: table, json, yaml (default "table")
```

## Installation

```bash
git clone <repo-url> defectdojo-cli
cd defectdojo-cli
go build -o dd .
```

### Requirements

- Go 1.21+

## Configuration

Configure via (in order of precedence):

1. **CLI flags**: `--host`, `--api-key`
2. **Environment variables**: `DD_HOST`, `DD_API_KEY`
3. **Config file**: `~/.dd/config.yaml`
4. **Default host**: `https://demo.defectdojo.org`

```bash
# Environment variables (recommended)
export DD_HOST=https://your-dojo.example.com
export DD_API_KEY=your-api-key-here

# Config file (~/.dd/config.yaml)
host: https://your-dojo.example.com
api-key: your-api-key-here
```

## Usage

### Resource CRUD

Every resource supports `list`, `get`, `create`, `update` (PUT), `patch` (PATCH), and `delete`:

```bash
# List
dd products list
dd products list --limit 20
dd products list --filter name=MyApp
dd products list -f name__contains=App -f lifecycle=production
dd products list --columns id,name,lifecycle,platform

# Get by ID
dd products get 1
dd products get 1 -o json     # JSON output
dd products get 1 -o yaml     # YAML output
dd products get 1 -o table    # Table output

# Create
dd products create --data '{"name":"MyApp","description":"My application","prod_type":1}'
dd products create --data-file product.json

# Update (full)
dd products update 1 --data '{"name":"MyApp","description":"updated","prod_type":1}'

# Partial update
dd products patch 1 --data '{"description":"Updated description"}'

# Delete
dd products delete 1
```

### Output Formats

| Flag | Format | Description |
|------|--------|-------------|
| `-o table` | Table | Default. Compact columns, human-readable |
| `-o json` | JSON | Full object data, indented |
| `-o yaml` | YAML | YAML format |

Use `--columns` to customize table output:

```bash
dd findings list --columns id,title,severity,test,active,found_by --limit 20
```

### Filtering

Use `--filter` (or `-f`) with Django-style field lookups:

```bash
dd findings list --filter severity=High
dd findings list -f severity=High -f active=true
dd findings list -f severity__in=Critical,High
dd findings list -f created__gt=2024-01-01
dd products list -f name__contains=App
```

### Pagination

By default, list commands fetch all pages. Use `--limit` to restrict:

```bash
dd findings list --limit 20          # First 20 results
dd findings list --limit 20 --page 2 # Page 2, 20 per page
```

### Available Resources

All 54 resources support `list`, `get`, `create`, `update` (PUT), `patch` (PATCH), `delete` by default. Resources marked † are deprecated (EOL 2026-06-01).

| Resource | CLI Name | Default Table Columns |
|----------|----------|----------------------|
| Product Types | `product-types` | id, name, description, critical_product, key_product |
| Products | `products` | id, name, description, prod_type, lifecycle, platform, findings_count, tags |
| Engagements | `engagements` | id, name, status, description, product, target_start, target_end, lead |
| Tests | `tests` | id, title, test_type, engagement, environment, branch_tag, created |
| Test Types | `test-types` | id, name, static_tool, dynamic_tool, active |
| Findings | `findings` | id, title, severity, status, test, found_by, date, mitigated, cwe, cvssv3_score |
| Finding Templates | `finding-templates` | id, title, severity, description, cwe |
| Endpoints | `endpoints` | id, host, port, protocol, product, findings_count |
| Endpoint Status | `endpoint-status` | id, endpoint, finding, mitigated, false_positive, out_of_scope, risk_accepted |
| Users | `users` | id, username, email, first_name, last_name, is_active, is_superuser, last_login |
| Development Environments | `development-environments` | id, name |
| Notes | `notes` | _(auto)_ |
| Metadata | `metadata` | _(auto)_ |
| Risk Acceptances | `risk-acceptance` | id, name, recommendation, decision, accepted_by, expiration_date |
| JIRA Instances | `jira-instances` | id, url, username |
| JIRA Projects | `jira-projects` | id, product, engagement, project_key, enabled, push_all_issues |
| JIRA Finding Mappings | `jira-finding-mappings` | id, jira_id, jira_key, finding, engagement |
| SonarQube Issues | `sonarqube-issues` | id, key, status, type |
| SonarQube Transitions | `sonarqube-transitions` | _(auto)_ |
| Tool Configurations | `tool-configurations` | id, name, tool_type, url, authentication_type |
| Tool Product Settings | `tool-product-settings` | id, name, product, tool_configuration, tool_project_id |
| Tool Types | `tool-types` | id, name, description |
| Regulations | `regulations` | id, name, description |
| Roles | `roles` | id, name |
| Dojo Groups | `dojo-groups` | id, name, social_provider |
| Dojo Group Members | `dojo-group-members` | id, group_id, user_id, role |
| Global Roles | `global-roles` | id, user, group, role |
| Product Members | `product-members` | id, product_id, user_id, role |
| Product Groups | `product-groups` | id, product_id, group_id, role |
| Product Type Members | `product-type-members` | id, product_type_id, user_id, role |
| Product Type Groups | `product-type-groups` | id, product_type_id, group_id, role |
| API Scan Configurations | `product-api-scan-configs` | id, product, tool_configuration, service_key_1 |
| App Analysis (Technologies) | `app-analysis` | id, product, name, version |
| Credentials † | `credentials` | _(auto)_ |
| Credential Mappings † | `credential-mappings` | _(auto)_ |
| Stub Findings † | `stub-findings` | _(auto)_ |
| System Settings | `system-settings` | _(auto)_ |
| Announcements | `announcements` | _(auto)_ |
| Notifications | `notifications` | _(auto)_ |
| Notification Webhooks | `notification-webhooks` | _(auto)_ |
| Languages | `languages` | _(auto)_ |
| Language Types | `language-types` | _(auto)_ |
| SLA Configurations | `sla-configurations` | id, name, critical, high, medium, low |
| Network Locations | `network-locations` | _(auto)_ |
| Engagement Presets | `engagement-presets` | _(auto)_ |
| Test Imports | `test-imports` | _(auto)_ |
| Configuration Permissions | `configuration-permissions` | _(auto)_ |
| Request / Response Pairs | `request-response-pairs` | _(auto)_ |
| User Contact Infos | `user-contact-infos` | _(auto)_ |
| Questionnaire Answers | `questionnaire-answers` | _(auto)_ |
| Answered Questionnaires | `questionnaire-answered-questionnaires` | _(auto)_ |
| Engagement Questionnaires | `questionnaire-engagement-questionnaires` | _(auto)_ |
| General Questionnaires | `questionnaire-general-questionnaires` | _(auto)_ |
| Questionnaire Questions | `questionnaire-questions` | _(auto)_ |

### Scan Import

```bash
# Import scan into existing engagement
dd import-scan \
  --engagement-id 1 \
  --scan-type "ZAP Scan" \
  --file results.json

# With options
dd import-scan \
  --engagement-id 1 \
  --scan-type "Nmap Scan" \
  --file scan.xml \
  --active \
  --verified \
  --close-old-findings \
  --minimum-severity 3 \
  --tags "automated,weekly" \
  --version "v1.2.3"

# Re-import to update existing test
dd reimport-scan \
  --test-id 1 \
  --scan-type "ZAP Scan" \
  --file results.json

# Quick scan — auto-creates product + engagement
dd scan \
  --product "MyApp" \
  --engagement "Sprint 42" \
  --scan-type "ZAP Scan" \
  --file results.json
```

### Finding Actions

```bash
dd finding-close 42                  # Close a finding
dd finding-verify 42                 # Verify a finding
dd finding-notes 42                  # List notes
dd finding-notes 42 --body "Fixed"   # Add note
dd finding-tags 42 --add "critical,auth"        # Add tags
dd finding-tags 42 --remove "old"               # Remove tags
dd finding-duplicate 42              # View duplicate cluster
dd finding-metadata 42               # View metadata
```

### Engagement Actions

```bash
dd engagement-close 1                # Close engagement
dd engagement-reopen 1               # Reopen engagement
dd engagement-notes 1                # List notes
dd engagement-notes 1 --body "Done"  # Add note
dd engagement-checklist 1            # View checklist
dd engagement-checklist 1 --data '{"status":"checked"}'  # Update
```

### User Actions

```bash
dd user-profile                      # Current user profile
dd user-reset-token 1                # Reset API token
```

### Other Commands

```bash
dd celery <task-id>                  # Check Celery task status
dd import-languages <product-id> <file>  # Import languages
dd endpoint-meta-import --product-id 1 --file meta.csv  # Import endpoint metadata
```

### Shell Completion

```bash
dd completion bash > /etc/bash_completion.d/dd
dd completion zsh > /usr/local/share/zsh/site-functions/_dd
dd completion fish > ~/.config/fish/completions/dd.fish
```

## Examples

### Typical Workflow

```bash
# 1. View available scan types
dd test-types list --columns id,name --limit 10

# 2. List products and pick one
dd products list

# 3. View engagements for that product
dd engagements list --filter product=1

# 4. Import a scan
dd import-scan --engagement-id 1 --scan-type "ZAP Scan" --file zap.json

# 5. Check findings
dd findings list --filter severity=High --columns id,title,severity,test

# 6. Investigate a finding
dd findings get 42 -o json

# 7. Close false positives
dd finding-close 42
```

### Scripting

```bash
# Export findings to JSON file
dd findings list --filter severity=Critical -o json > critical-findings.json

# Count findings per severity
dd findings list -o json | jq '.results | group_by(.severity) | map({severity: .[0].severity, count: length})'

# Bulk close inactive findings
dd findings list -f active=false --limit 100 -o json | \
  jq -r '.results[].id' | xargs -I{} dd finding-close {}
```

## API Reference

This CLI covers all DefectDojo API v2 endpoints defined in `dojo/api_v2/views.py` (37 ViewSets) and registered in `dojo/urls.py`. Resources are auto-generated from the router registrations, and special actions (close, verify, notes, tags, etc.) are exposed as dedicated commands.
