# Migration Guide

## Upgrading to v0.3.0 from v0.2.0

### Overview

Version 0.3.0 updates the inventory API limits to align with the Manapool API v0.24.1 specification. The API now supports fetching up to 10,000 inventory items per request (previously limited to 500), and the default page size has been changed from 500 to 100 items.

### Breaking Changes

#### Default Inventory Limit Changed: 500 → 100

**Impact:** Code that relies on the implicit default limit will fetch 100 items per request instead of 500.

**Who is affected:**
- Applications using `InventoryOptions{}` or `InventoryOptions{Limit: 0}` without explicitly setting the limit
- Applications that assumed paginated requests would return 500 items by default

**Migration:**

If your application depends on receiving 500 items per page by default, explicitly set the limit:

```go
// Before (v0.2.0)
opts := manapool.InventoryOptions{
    Offset: 0,
    // Implicitly used limit of 500
}

// After (v0.3.0) - Explicit limit if you need 500 items
opts := manapool.InventoryOptions{
    Limit:  500,  // Explicitly set to 500 if required
    Offset: 0,
}
```

**Recommended approach:**

Most applications should embrace the new default of 100 or choose an optimal page size:

```go
// Recommended: Use the new default (100)
opts := manapool.InventoryOptions{
    Offset: 0,
    // Will use default of 100
}

// Or: Choose a page size appropriate for your use case
opts := manapool.InventoryOptions{
    Limit:  1000,  // Larger pages for batch processing
    Offset: 0,
}
```

### New Capabilities

#### Maximum Limit Increased: 500 → 10,000

**Benefit:** You can now fetch up to 10,000 inventory items in a single request, significantly reducing the number of API calls needed for large inventories.

**Example:**

```go
// Now valid (was rejected in v0.2.0)
opts := manapool.InventoryOptions{
    Limit:  5000,  // Fetch 5,000 items at once
    Offset: 0,
}
inventory, err := client.GetSellerInventory(ctx, opts)
```

**Performance optimization:**

For large inventory processing, use larger page sizes to reduce API overhead:

```go
// Efficient batch processing of large inventories
opts := manapool.InventoryOptions{
    Limit:  5000,  // Larger page size
    Offset: 0,
}

for {
    resp, err := client.GetSellerInventory(ctx, opts)
    if err != nil {
        return err
    }
    
    // Process items...
    for _, item := range resp.Inventory {
        processItem(item)
    }
    
    // Check if done
    if opts.Offset + resp.Pagination.Returned >= resp.Pagination.Total {
        break
    }
    
    opts.Offset += resp.Pagination.Returned
}
```

#### IterateInventory Performance Improvement

**What changed:** The `IterateInventory()` helper function now fetches 1,000 items per page instead of 500, reducing the number of API calls by half for large inventories.

**Impact:** Positive performance improvement with no code changes required.

```go
// No changes needed - automatically benefits from larger page size
err := manapool.IterateInventory(ctx, client, func(item *manapool.InventoryItem) error {
    fmt.Printf("%s: $%.2f\n", item.Product.Single.Name, item.PriceDollars())
    return nil
})
```

### Compatibility Matrix

| Scenario | v0.2.0 | v0.3.0 | Action Required |
|----------|--------|--------|-----------------|
| Using `Limit: 100` | ✅ Works | ✅ Works | None |
| Using `Limit: 500` | ✅ Works | ✅ Works | None |
| Using `Limit: 0` (default) | Returns 500 | Returns 100 | Review if dependent on 500-item pages |
| Omitting `Limit` field | Returns 500 | Returns 100 | Review if dependent on 500-item pages |
| Using `Limit: 501-10000` | ❌ Validation error | ✅ Works | Update to benefit from larger pages |
| Using `IterateInventory()` | 500 per page | 1000 per page | No changes needed |

### Testing Your Migration

#### 1. Unit Tests

Update tests that assume the default limit:

```go
// Before (v0.2.0)
func TestMyInventoryFunction(t *testing.T) {
    // Test assumed 500 items per page
    if len(result.Inventory) != 500 {
        t.Errorf("expected 500 items")
    }
}

// After (v0.3.0)
func TestMyInventoryFunction(t *testing.T) {
    // Explicitly set limit in test
    opts := manapool.InventoryOptions{Limit: 500}
    // ... or adjust expectations ...
    if len(result.Inventory) > 100 {
        t.Errorf("expected default page size of 100 or less")
    }
}
```

#### 2. Integration Tests

Test with your actual inventory to ensure pagination works correctly:

```go
func TestPaginationWithNewLimits(t *testing.T) {
    client := manapool.NewClient(token, email)
    
    // Test with new default
    opts := manapool.InventoryOptions{Offset: 0}
    resp, err := client.GetSellerInventory(ctx, opts)
    if err != nil {
        t.Fatal(err)
    }
    
    // Verify pagination metadata
    if resp.Pagination.Limit != 100 {
        t.Errorf("expected default limit of 100, got %d", resp.Pagination.Limit)
    }
}
```

#### 3. Performance Testing

Benchmark the impact of different page sizes on your workload:

```go
func BenchmarkInventoryProcessing(b *testing.B) {
    pageSizes := []int{100, 500, 1000, 5000}
    
    for _, size := range pageSizes {
        b.Run(fmt.Sprintf("PageSize_%d", size), func(b *testing.B) {
            opts := manapool.InventoryOptions{Limit: size}
            // Benchmark your processing logic
        })
    }
}
```

### Rollback Plan

If you encounter issues after upgrading to v0.3.0, you can:

#### Option 1: Pin to v0.2.0

```bash
go get github.com/repricah/manapool@v0.2.0
```

#### Option 2: Explicitly Set Limits

Update your code to explicitly specify the desired limit:

```go
// Maintain v0.2.0 behavior in v0.3.0
opts := manapool.InventoryOptions{
    Limit: 500,  // Explicit limit
    Offset: 0,
}
```

### Recommended Upgrade Path

1. **Review your code** for uses of `InventoryOptions` without explicit `Limit`
2. **Update tests** that assume 500-item default pages
3. **Run your test suite** to catch any pagination-dependent logic
4. **Update to v0.3.0** using `go get -u github.com/repricah/manapool@v0.3.0`
5. **Verify in staging** with real API calls
6. **Deploy to production** after validation
7. **Monitor** for any unexpected behavior related to pagination

### Optimization Opportunities

Take advantage of the increased limits:

```go
// For batch processing: Use larger pages to reduce API overhead
opts := manapool.InventoryOptions{
    Limit: 5000,  // Fewer API calls for large inventories
    Offset: 0,
}

// For UI pagination: Use smaller pages for better UX
opts := manapool.InventoryOptions{
    Limit: 50,  // Faster initial response
    Offset: page * 50,
}

// For streaming: Use the helper function (automatically optimized)
err := manapool.IterateInventory(ctx, client, func(item *manapool.InventoryItem) error {
    // Process each item efficiently
    return processItem(item)
})
```

### Getting Help

If you encounter issues during migration:

1. Check the [CHANGELOG.md](./CHANGELOG.md) for detailed changes
2. Review the [README.md](./README.md) for updated examples
3. Open an issue at https://github.com/repricah/manapool/issues
4. See [API documentation](https://pkg.go.dev/github.com/repricah/manapool) for reference

### FAQ

**Q: Do I need to change my code if I'm already setting `Limit` explicitly?**  
A: No, if you're explicitly setting `Limit` to any value ≤ 500, your code will work without changes.

**Q: Will this break my application?**  
A: Unlikely. The change is backward compatible unless your application logic assumes exactly 500 items per page when using the default.

**Q: Should I increase my page size to 10,000?**  
A: It depends on your use case. Larger pages reduce API calls but increase memory usage and response time. For most use cases, 500-2000 items per page is optimal.

**Q: Does this affect rate limiting?**  
A: No, rate limits are still based on the number of requests, not the number of items per request. Larger pages mean fewer requests.

**Q: How do I know what page size to use?**  
A: 
- **Small inventories (<1000 items)**: Default of 100 is fine
- **Medium inventories (1000-10000 items)**: Use 500-1000 per page
- **Large inventories (>10000 items)**: Use 2000-5000 per page
- **Streaming/processing**: Use `IterateInventory()` (automatically optimized)
