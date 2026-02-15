# AI Agent Handoff: Manapool Client v0.2.0 → v0.3.0 Upgrade

## Overview

This document provides instructions for AI agents to perform automated upgrades of the Manapool Go client from v0.2.0 to v0.3.0.

**Change Summary:** API inventory limits updated (max: 500→10000, default: 500→100)

## Prerequisites

Before executing the upgrade:

1. **Repository access**: Read/write access to the target repository
2. **Go version**: Confirm Go 1.23+ is available
3. **Test environment**: Ability to run `go test ./...`
4. **API credentials**: Optional, for integration testing (not required for upgrade)

## Upgrade Steps

### Step 1: Identify Affected Code

Search for code that uses inventory operations without explicit limits:

```bash
# Find InventoryOptions without explicit Limit field
grep -r "InventoryOptions{" . --include="*.go" | grep -v "Limit:"

# Find uses of GetSellerInventory
grep -r "GetSellerInventory" . --include="*.go"

# Find uses of IterateInventory
grep -r "IterateInventory" . --include="*.go"
```

**Expected patterns to review:**

```go
// Pattern 1: No limit specified (uses default)
opts := manapool.InventoryOptions{
    Offset: 0,
}

// Pattern 2: Limit set to 0 (uses default)
opts := manapool.InventoryOptions{
    Limit:  0,
    Offset: 0,
}

// Pattern 3: Empty struct (uses all defaults)
opts := manapool.InventoryOptions{}
```

### Step 2: Analyze Impact

For each identified usage, determine if the code is affected:

**Not affected (no changes needed):**
- ✅ Code explicitly sets `Limit` to any value (e.g., `Limit: 100`, `Limit: 500`)
- ✅ Code uses `IterateInventory()` without customization
- ✅ Tests that don't assert on specific page sizes

**Potentially affected (review required):**
- ⚠️ Code that omits `Limit` or sets it to 0
- ⚠️ Code that processes results assuming 500 items per page
- ⚠️ Tests that assert `len(items) == 500` for default pagination
- ⚠️ Performance-critical code that may benefit from larger pages

### Step 3: Update Dependency

```bash
# Update to v0.3.0
go get -u github.com/repricah/manapool@v0.3.0

# Clean up
go mod tidy
```

### Step 4: Code Modifications

#### Option A: Maintain v0.2.0 Behavior (Conservative)

For each affected usage, explicitly set `Limit: 500`:

```go
// Before (v0.2.0)
opts := manapool.InventoryOptions{
    Offset: offset,
}

// After (v0.3.0) - Explicit limit to maintain behavior
opts := manapool.InventoryOptions{
    Limit:  500,  // Explicitly set to maintain v0.2.0 default
    Offset: offset,
}
```

#### Option B: Adopt v0.3.0 Defaults (Recommended)

Leave the code as-is to use the new default, but verify the logic handles variable page sizes:

```go
// No changes needed - will use new default of 100
opts := manapool.InventoryOptions{
    Offset: offset,
}

// Ensure pagination logic doesn't assume 500 items per page
for {
    resp, err := client.GetSellerInventory(ctx, opts)
    if err != nil {
        return err
    }
    
    processItems(resp.Inventory)  // Handle any number of items
    
    // Correct: Check using pagination metadata
    if opts.Offset + resp.Pagination.Returned >= resp.Pagination.Total {
        break
    }
    opts.Offset += resp.Pagination.Returned
}
```

#### Option C: Optimize with Larger Pages

For batch processing, increase page sizes for better performance:

```go
// Optimize for large inventory processing
opts := manapool.InventoryOptions{
    Limit:  2000,  // Larger pages for efficiency
    Offset: offset,
}
```

### Step 5: Update Tests

Find and update tests that assume specific page sizes:

```bash
# Find tests that check length
grep -r "len.*Inventory.*==" . --include="*_test.go"

# Find tests that assert 500
grep -r "500" . --include="*_test.go" | grep -i "inventory\|limit"
```

Update test assertions:

```go
// Before (v0.2.0)
func TestInventoryPagination(t *testing.T) {
    opts := manapool.InventoryOptions{}  // Used default of 500
    resp, err := client.GetSellerInventory(ctx, opts)
    
    if resp.Pagination.Limit != 500 {
        t.Errorf("expected default limit of 500")
    }
}

// After (v0.3.0) - Option 1: Update expectation
func TestInventoryPagination(t *testing.T) {
    opts := manapool.InventoryOptions{}  // Uses new default of 100
    resp, err := client.GetSellerInventory(ctx, opts)
    
    if resp.Pagination.Limit != 100 {
        t.Errorf("expected default limit of 100")
    }
}

// After (v0.3.0) - Option 2: Explicit limit in test
func TestInventoryPagination(t *testing.T) {
    opts := manapool.InventoryOptions{Limit: 500}  // Explicit for test
    resp, err := client.GetSellerInventory(ctx, opts)
    
    if resp.Pagination.Limit != 500 {
        t.Errorf("expected limit of 500")
    }
}
```

### Step 6: Validation

Run comprehensive validation:

```bash
# 1. Compile check
go build ./...

# 2. Run all tests
go test ./... -v

# 3. Run with race detector
go test ./... -race

# 4. Run with coverage
go test ./... -cover

# 5. Check for common issues
go vet ./...

# 6. Format code
go fmt ./...
```

**Expected results:**
- ✅ All tests pass
- ✅ No race conditions detected
- ✅ Coverage maintained or improved
- ✅ No vet warnings

### Step 7: Integration Testing (Optional)

If API credentials are available, test against the real API:

```bash
export MANAPOOL_API_TOKEN="your-token"
export MANAPOOL_API_EMAIL="your-email"

# Run a real inventory query with different page sizes
go run -tags=integration ./cmd/test-inventory/
```

### Step 8: Commit Changes

```bash
# Stage changes
git add go.mod go.sum

# If code changes were made
git add <modified-files>

# Commit with clear message
git commit -m "chore: upgrade manapool client to v0.3.0

- Update dependency: github.com/repricah/manapool v0.2.0 → v0.3.0
- Updated inventory options to use new default limit (100)
- Verified all tests pass with new pagination defaults

See MIGRATION_GUIDE.md for details on the changes."
```

## Automated Upgrade Script

For AI agents that can execute shell scripts:

```bash
#!/bin/bash
set -e

echo "=== Manapool Client v0.2.0 → v0.3.0 Upgrade ==="

# Step 1: Backup
echo "Creating backup branch..."
git checkout -b backup-before-manapool-upgrade

# Step 2: Return to main branch
git checkout main

# Step 3: Update dependency
echo "Updating dependency..."
go get -u github.com/repricah/manapool@v0.3.0
go mod tidy

# Step 4: Run tests
echo "Running tests..."
if ! go test ./... -v; then
    echo "❌ Tests failed. Please review manually."
    exit 1
fi

# Step 5: Check for common issues
echo "Checking for issues..."
go vet ./...

# Step 6: Format code
echo "Formatting code..."
go fmt ./...

# Step 7: Commit
echo "Committing changes..."
git add go.mod go.sum
git commit -m "chore: upgrade manapool client to v0.3.0

- Update dependency: github.com/repricah/manapool v0.2.0 → v0.3.0
- Verified all tests pass with new defaults

See MIGRATION_GUIDE.md for details."

echo "✅ Upgrade complete!"
echo "Next steps:"
echo "  1. Review changes: git diff backup-before-manapool-upgrade"
echo "  2. Test manually if integration tests are available"
echo "  3. Push changes: git push origin main"
```

## Decision Tree for AI Agents

```
START: Upgrade manapool v0.2.0 → v0.3.0

1. Scan codebase for InventoryOptions usage
   ├─ All uses have explicit Limit set?
   │  ├─ YES → Proceed with simple upgrade (Step 3)
   │  └─ NO → Continue to step 2
   │
2. Analyze impact of default limit change
   ├─ Code logic assumes 500 items per page?
   │  ├─ YES → Apply Option A (explicit Limit: 500)
   │  └─ NO → Apply Option B (use new defaults)
   │
3. Update dependency (go get -u)
   │
4. Update tests with changed assertions
   │
5. Run validation suite
   ├─ Tests pass?
   │  ├─ YES → Proceed to commit
   │  └─ NO → Review failures, apply fixes, retry
   │
6. Commit changes with descriptive message
   │
END: Upgrade complete
```

## Common Pitfalls

### Pitfall 1: Hard-coded Expectations

**Problem:**
```go
// Code assumes exactly 500 items per page
items := make([]Item, 500)
for i, inventoryItem := range resp.Inventory {
    items[i] = convertItem(inventoryItem)
}
```

**Solution:**
```go
// Use actual returned count
items := make([]Item, len(resp.Inventory))
for i, inventoryItem := range resp.Inventory {
    items[i] = convertItem(inventoryItem)
}
```

### Pitfall 2: Test Data Assumptions

**Problem:**
```go
// Test data prepared for 500 items
mockInventory := generateMockInventory(500)
```

**Solution:**
```go
// Test data should match new default or be explicit
mockInventory := generateMockInventory(100)  // Or set explicit limit in test
```

### Pitfall 3: Progress Indicators

**Problem:**
```go
// Progress calculation assumes 500 per page
progress := (page * 500) / totalItems
```

**Solution:**
```go
// Use actual pagination metadata
progress := opts.Offset / resp.Pagination.Total
```

## Rollback Procedure

If issues are encountered:

```bash
# Option 1: Revert to v0.2.0
go get github.com/repricah/manapool@v0.2.0
go mod tidy

# Option 2: Restore from backup branch
git checkout backup-before-manapool-upgrade
git branch -D main
git checkout -b main
```

## Verification Checklist

- [ ] Dependency updated to v0.3.0 in go.mod
- [ ] All uses of InventoryOptions reviewed
- [ ] Tests updated for new default limit (if applicable)
- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Code formatted: `go fmt ./...`
- [ ] No vet warnings: `go vet ./...`
- [ ] Integration tests pass (if available)
- [ ] Changes committed with clear message
- [ ] MIGRATION_GUIDE.md reviewed (if manual steps needed)

## Success Criteria

The upgrade is successful when:

1. ✅ Application compiles without errors
2. ✅ All unit tests pass
3. ✅ Integration tests pass (if available)
4. ✅ No regression in functionality
5. ✅ Performance is maintained or improved
6. ✅ Code follows project style guidelines

## Documentation Updates

After upgrade, update these documents if present:

- [ ] README.md - Update version references
- [ ] CHANGELOG.md - Add entry for dependency update
- [ ] API documentation - Update examples if needed
- [ ] Integration docs - Update page size recommendations

## Support

For issues during automated upgrade:

1. **Review logs**: Check test output for specific failures
2. **Manual review**: Flag complex cases for human review
3. **Documentation**: Refer to MIGRATION_GUIDE.md
4. **Issue tracker**: https://github.com/repricah/manapool/issues

## Agent Capabilities Required

This upgrade requires an AI agent with:

- ✅ File reading/writing
- ✅ Command execution (bash, go commands)
- ✅ Code analysis (pattern matching, AST parsing)
- ✅ Git operations (commit, branch)
- ✅ Test execution and result parsing
- ⚠️ Optional: Integration test execution with API credentials

## Estimated Time

- **Simple upgrade** (no code changes needed): 2-5 minutes
- **Standard upgrade** (minor code adjustments): 10-20 minutes
- **Complex upgrade** (extensive code review needed): 30-60 minutes

---

**Version:** 1.0  
**Last Updated:** 2026-02-15  
**Compatibility:** Manapool client v0.2.0 → v0.3.0
