# Exports Directory

This directory stores exported test results in JSON and CSV formats.

## File Naming Convention

```
test_YYYYMMDD_HHMMSS.json
test_YYYYMMDD_HHMMSS.csv

Example:
test_20241221_143052.json
test_20241221_143052.csv
```

## Export Formats

### JSON Format

- **Complete data preservation**
- Includes all request details
- Easy to parse programmatically
- Can be re-imported for analysis

### CSV Format

- **Spreadsheet-friendly**
- Opens in Excel, Google Sheets
- Easy to analyze with pivot tables
- Summary statistics + individual requests

## What Gets Exported

Both formats include:

### Test Configuration

- Target URL
- Total requests
- Duration
- Traffic type (normal/attack)
- Test mode (baseline/burst)
- HTTP method
- Error mode
- User-agent settings

### Test Results

- Test ID and timestamps
- Success/error/blocked counts
- Response time statistics (avg, min, max, p50, p95, p99)
- Rate limit detection
- Actual vs expected rate

### Individual Requests

- Request ID and timestamp
- Status code and text
- Response time
- URL (including attack payloads if applicable)
- Headers
- Error messages
- Attack info (if attack mode)

## Using Exports

### View in Terminal

```bash
# View JSON (pretty-printed)
cat test_20241221_143052.json | json_pp

# View CSV
cat test_20241221_143052.csv
```

### Open in Excel/Sheets

1. Download the CSV file
2. Open in Excel or Google Sheets
3. Use pivot tables, charts for analysis

### Compare Tests

```bash
# Compare two JSON exports
diff test_20241221_143052.json test_20241221_144523.json
```

### Parse Programmatically

```python
import json

# Load JSON export
with open('test_20241221_143052.json', 'r') as f:
    data = json.load(f)

print(f"Success rate: {data['results']['success_count']}/{data['results']['total_requests']}")
print(f"p95 response time: {data['results']['p95_response']}ms")
```

## 🧹 Maintenance

### Clean Old Exports

```bash
# Delete exports older than 30 days
find . -name "test_*.json" -mtime +30 -delete
find . -name "test_*.csv" -mtime +30 -delete
```

### Archive Exports

```bash
# Create archive of all exports
tar -czf exports_backup_$(date +%Y%m%d).tar.gz test_*.json test_*.csv
```

## Best Practices

1. **Export after each test** - Don't lose important results
2. **Use descriptive naming** - Timestamp helps identify tests
3. **Regular cleanup** - Remove old exports to save space
4. **JSON for automation** - Use JSON for scripts and tools
5. **CSV for analysis** - Use CSV for manual analysis in spreadsheets

## Security Note

Exported files may contain:

- Target URLs
- API endpoints
- Attack payloads
- Response data

**Keep exports secure** and don't commit them to version control!

Add to `.gitignore`:

```
exports/*.json
exports/*.csv
!exports/README.md
```
