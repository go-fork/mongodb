# Migration Guide - v0.1.2

## Overview
This guide helps you migrate from v0.1.1 to v0.1.2. Most changes are internal improvements and GitHub integration that should not affect existing code.

## Prerequisites
- Go 1.23 or later
- Previous version (v0.1.1) installed

## Quick Migration Checklist
- [ ] Update import statements if using mocks (package name changed to mongodb_mocks)
- [ ] Update any hardcoded repository URLs to use go-fork organization
- [ ] Regenerate mocks if customized (new .mockery.yaml configuration)
- [ ] Update CI/CD references if pointing to old repository structure
- [ ] Run tests to ensure compatibility

## Changes That May Affect You

### Mock Package Naming
**Before (v0.1.1):**
```go
import "go.fork.vn/mongodb/mocks"
```

**After (v0.1.2):**
```go
import "go.fork.vn/mongodb/mocks" // Package is now mongodb_mocks
```

### Repository URL Changes
If you have any references to the repository in documentation or configuration:

**Before:**
- Old organization URLs

**After:**
- `github.com/go-fork/mongodb`
- All references now use go-fork organization format

### Mock Generation
If you regenerate mocks, use the updated `.mockery.yaml` configuration which provides better mock generation.

## Breaking Changes

### API Changes
#### Changed Functions
```go
// Old way (previous version)
oldFunction(param1, param2)

// New way (v0.1.2)
newFunction(param1, param2, newParam)
```

#### Removed Functions
- `removedFunction()` - Use `newAlternativeFunction()` instead

#### Changed Types
```go
// Old type definition
type OldConfig struct {
    Field1 string
    Field2 int
}

// New type definition
type NewConfig struct {
    Field1 string
    Field2 int64 // Changed from int
    Field3 bool  // New field
}
```

### Configuration Changes
If you're using configuration files:

```yaml
# Old configuration format
old_setting: value
deprecated_option: true

# New configuration format
new_setting: value
# deprecated_option removed
new_option: false
```

## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get go.fork.vn/mongodb@v0.1.2
go mod tidy
```

### Step 2: Update Import Statements
```go
// If import paths changed
import (
    "go.fork.vn/mongodb" // Updated import
)
```

### Step 3: Update Code
Replace deprecated function calls:

```go
// Before
result := mongodb.OldFunction(param)

// After
result := mongodb.NewFunction(param, defaultValue)
```

### Step 4: Update Configuration
Update your configuration files according to the new schema.

### Step 5: Run Tests
```bash
go test ./...
```

## Common Issues and Solutions

### Issue 1: Function Not Found
**Problem**: `undefined: mongodb.OldFunction`  
**Solution**: Replace with `mongodb.NewFunction`

### Issue 2: Type Mismatch
**Problem**: `cannot use int as int64`  
**Solution**: Cast the value or update variable type

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/mongodb@v0.1.2)
- Search [existing issues](https://github.com/go-fork/mongodb/issues)
- Create a [new issue](https://github.com/go-fork/mongodb/issues/new) if needed

## Rollback Instructions
If you need to rollback:

```bash
go get go.fork.vn/mongodb@previous-version
go mod tidy
```

Replace `previous-version` with your previous version tag.

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
