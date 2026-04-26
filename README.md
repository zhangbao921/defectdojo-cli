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

| Resource | CLI Name | Default Columns |
|----------|----------|----------------|
| Product Types | `product-types` | id, name, description, critical_product, key_product |
| Products | `products` | id, name, description, prod_type, lifecycle, platform, findings_count |
| Engagements | `engagements` | id, name, status, description, product, target_start, target_end |
| Tests | `tests` | id, title, test_type, engagement, environment, branch_tag |
| Test Types | `test-types` | id, name, static_tool, dynamic_tool |
| Findings | `findings` | id, title, severity, status, test, found_by, date, cwe |
| Finding Templates | `finding-templates` | id, title, severity, description, cwe |
| Endpoints | `endpoints` | id, host, port, protocol, product |
| Endpoint Status | `endpoint-status` | id, endpoint, finding, mitigated, false_positive |
| Users | `users` | id, username, email, is_active, is_superuser |
| Development Environments | `development-environments` | id, name |
| JIRA Instances | `jira-instances` | id, url, username |
| JIRA Projects | `jira-projects` | id, product, engagement, project_key, enabled |
| JIRA Finding Mappings | `jira-finding-mappings` | id, jira_id, jira_key, finding |
| SonarQube Issues | `sonarqube-issues` | id, key, status, type |
| Tool Configurations | `tool-configurations` | id, name, tool_type, url |
| Tool Types | `tool-types` | id, name, description |
| Regulations | `regulations` | id, name, description |
| Roles | `roles` | id, name |
| Dojo Groups | `dojo-groups` | id, name, social_provider |
| Product Members | `product-members` | id, product_id, user_id, role |
| Product Groups | `product-groups` | id, product_id, group_id, role |
| … and 30+ more | | |

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
