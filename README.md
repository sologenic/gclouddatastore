## Cloud Datastore Go Client (Fork)

[![Go Reference](https://pkg.go.dev/badge/github.com/sologenic/gclouddatastore.svg)](https://pkg.go.dev/github.com/sologenic/gclouddatastore)

This is a fork of the [Google Cloud Datastore Go client library](https://github.com/googleapis/google-cloud-go/tree/main/datastore) with enhanced functionality for partial object parsing, improved query filtering, and optimized count operations.

The changes here (relative to Google’s `datastore` package) come from several years of running **Cloud Datastore in Datastore mode**. It is an interesting database to operate, but the official Go client has limitations and bugs; this fork addresses those. Because it is maintained by someone who uses Datastore daily, the fixes remove practical, day-to-day annoyances that Go developers can hit with the upstream package.

As a bonus, the fork also drives down **data storage, backup, and snapshot** costs—by an estimated **~95%** in the author’s workloads (on the order of ~1 billion records total). Hypothetically, writes can be faster too: less I/O and CPU.

The package surface is **100% backwards compatible** with Google’s client and works as a **drop-in replacement**.

**Original Package**: `cloud.google.com/go/datastore`  
**This Fork**: `github.com/sologenic/gclouddatastore`

### Key Enhancements

This fork includes the following improvements over the original Google package:

1. **Partial Object Parsing**: Load entities into structs that only contain a subset of the datastore fields. Extra fields from datastore are gracefully skipped instead of causing errors.

2. **GetAllWithUnparsedFields**: New query function that returns both the loaded entities and information about fields that couldn't be parsed (with their types).

3. **Automatic Type Casting in Filters**: `FilterField` now automatically casts custom types (enums, custom int/float types) to base types, eliminating the need for manual casting before querying.

4. **Optimized Count Function**: The `Count()` function now uses GQL aggregation queries for efficient server-side counting instead of loading all records into memory.

### License

This package is licensed under the Apache License 2.0, same as the original Google package. See [LICENSE](LICENSE) for details.

### Original Documentation

- [About Cloud Datastore](https://cloud.google.com/datastore/)
- [Activating the API for your project](https://cloud.google.com/datastore/docs/activate)
- [API documentation](https://cloud.google.com/datastore/docs)
- [Original Go client documentation](https://pkg.go.dev/cloud.google.com/go/datastore)
- [Complete sample program](https://github.com/GoogleCloudPlatform/golang-samples/tree/main/datastore/tasks)

### Installation

```bash
go get github.com/sologenic/gclouddatastore
```

### Example Usage

#### Basic Usage

First create a `datastore.Client` to use throughout your application:

```go
import "github.com/sologenic/gclouddatastore"

client, err := datastore.NewClient(ctx, "my-project-id")
if err != nil {
	log.Fatal(err)
}
```

#### Partial Object Parsing

Load entities into structs that only contain a subset of fields:

```go
// Datastore entity has: Name, Age, Email, Address, Phone
type UserSummary struct {
	Name  string
	Email string
	// Age, Address, Phone are skipped - no error!
}

var users []UserSummary
keys, err := client.GetAll(ctx, datastore.NewQuery("User"), &users)
// Works even if datastore has extra fields not in UserSummary
```

#### Get Unparsed Fields Information

Get information about fields that couldn't be parsed:

```go
type PartialUser struct {
	Name string
}

var users []PartialUser
keys, unparsed, err := client.GetAllWithUnparsedFields(ctx, 
	datastore.NewQuery("User"), &users)

// unparsed map contains: {"Age": "int64", "Email": "string", ...}
for fieldName, fieldType := range unparsed {
	fmt.Printf("Unparsed field: %s (type: %s)\n", fieldName, fieldType)
}
```

#### Automatic Type Casting in Filters

Use custom types and enums directly in filters without manual casting:

```go
type Status int
const (
	StatusActive Status = 1
	StatusInactive Status = 2
)

// No need to cast to int64 - auto-casting handles it!
query := datastore.NewQuery("Post").
	FilterField("Status", "=", StatusActive)

// Works with arrays too
query := datastore.NewQuery("Post").
	FilterField("Status", "in", []Status{StatusActive, StatusInactive})
```

#### Optimized Count Queries

Count operations are now efficient and don't load records into memory:

```go
count, err := client.Count(ctx, datastore.NewQuery("User").
	FilterField("Active", "=", true))
// Uses server-side aggregation - fast even for large datasets
```

### Migration from Original Package

To migrate from `cloud.google.com/go/datastore`:

1. Update import: `cloud.google.com/go/datastore` → `github.com/sologenic/gclouddatastore`
2. Done

Other fork features: partial parsing, auto-casting in filters, optimized `Count`, and `GetAllWithUnparsedFields` behave as documented above.

### Changes from Original Package

- **Partial Object Support**: Missing fields no longer cause `ErrFieldMismatch` errors
- **GetAllWithUnparsedFields**: New function for getting unparsed field information
- **FilterField Auto-Casting**: Custom types automatically converted to base types
- **Count Optimization**: Uses GQL aggregation instead of loading records
