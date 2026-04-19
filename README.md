## Cloud Datastore Go Client (Fork)

This is a fork of the [Google Cloud Datastore Go client library](https://github.com/googleapis/google-cloud-go/tree/main/datastore) with enhanced functionality for partial object parsing, improved query filtering, and optimized count operations.

**Original Package**: `cloud.google.com/go/datastore`  
**This Fork**: `github.com/norbertvannobelen/gclouddatastore`

### Key Enhancements

This fork includes the following improvements over the original Google package:

1. **Partial Object Parsing**: Load entities into structs that only contain a subset of the datastore fields. Extra fields from datastore are gracefully skipped instead of causing errors.

2. **GetAllWithUnparsedFields**: New query function that returns both the loaded entities and information about fields that couldn't be parsed (with their types).

3. **Automatic Type Casting in Filters**: `FilterField` now automatically casts custom types (enums, custom int/float types) to base types, eliminating the need for manual casting before querying.

4. **Optimized Count Function**: The `Count()` function now uses GQL aggregation queries for efficient server-side counting instead of loading all records into memory.

5. **Default no index**: Struct fields and `datastore.Property` values are excluded from Datastore indexes unless you opt in. Use the `index` struct tag, `datastore.RegisterIndexedFields`, or `datastore.Indexed(name, value)` for manual `PropertyLoadSaver` code. This reduces index storage and avoids tagging every large blob with `noindex`.

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
go get github.com/norbertvannobelen/gclouddatastore
```

### Example Usage

#### Basic Usage

First create a `datastore.Client` to use throughout your application:

```go
import "github.com/norbertvannobelen/gclouddatastore"

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

1. Update import: `cloud.google.com/go/datastore` → `github.com/norbertvannobelen/gclouddatastore`
2. Add `datastore:",index"` (or `name,index` with a name prefix) on **every struct field** you query with `Filter`, `FilterField`, `Order`, or projections. Fields omitted from indexes are not queryable the same way as before.
3. **Property / PropertyLoadSaver**: Replace `NoIndex bool` with `Index bool`. Old `NoIndex: false` (default) meant “indexed” → set `Index: true`. Old `NoIndex: true` → omit `Index` or `Index: false`.
4. Deploy or adjust composite indexes in your project if query fields change.

Other fork features: partial parsing, auto-casting in filters, optimized `Count`, and `GetAllWithUnparsedFields` behave as documented above.

### Changes from Original Package

- **Partial Object Support**: Missing fields no longer cause `ErrFieldMismatch` errors
- **GetAllWithUnparsedFields**: New function for getting unparsed field information
- **FilterField Auto-Casting**: Custom types automatically converted to base types
- **Count Optimization**: Uses GQL aggregation instead of loading records
- **Indexing defaults**: Excluded from indexes by default; explicit `index` tag or `RegisterIndexedFields` to opt in (`Property` uses `Index` instead of `NoIndex`)
