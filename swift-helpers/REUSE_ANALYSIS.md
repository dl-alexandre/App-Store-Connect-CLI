# Swift Helper Tools - What We Brought Over vs. What We Won't Use

## Overview

This document tracks what was copied from the original `swift-asc-cli/` project and what we're intentionally not using.

## Architecture Difference

**swift-asc-cli**: Full Swift CLI replacement (standalone, 70+ commands)
**swift-helpers**: Minimal helper tools that Go CLI shells out to (3 focused helpers)

## What We Brought Over

### From `swift-asc-cli/Sources/ASCAuth/`

| Original File | What We Used | Where It Went |
|---------------|--------------|---------------|
| `JWT.swift` | JWT generation logic, base64url encoding, P-256 signing | `swift-helpers/Sources/asc-jwt-sign/main.swift` (simplified) |
| `Keychain.swift` | Keychain storage patterns, OSStatus error mapping, credential structure | `swift-helpers/Sources/asc-keychain/main.swift` (simplified) |

### Key Reuse Patterns

1. **JWT Structure**: Same header/claims format, 10-minute token lifetime, ES256 algorithm
2. **Keychain Service Name**: `asc:credential:` prefix, `asc` service name
3. **Error Mapping**: Same OSStatus codes (-25293 access denied, -25308 interaction not allowed, etc.)
4. **Base64URL Encoding**: Same character replacement logic (`+` → `-`, `/` → `_`, strip `=`)

## What We Won't Use

### Full Libraries Not Copied

The following from `swift-asc-cli/` are NOT being used in the helper tools approach:

#### ASCAuth Library (not used)
- `ASCAuth.swift` - Full authentication command wrapper
- `CredentialResolver.swift` - Complex credential resolution with priority ordering
- `KeyLoader.swift` - Key loading utilities (simplified in our version)

#### ASCCommands Library (not used)
- `ASCCommands.swift` - Full command implementations (70+ commands)
- `ASCCredentialTokenProvider.swift` - Token provider for API calls
- `SharedCommandUtils.swift` - Command utilities, output formatters
- `Auth/` - Auth subcommands (add, list, remove, use)
- `Apps/` - App management subcommands
- `Builds/` - Build operations subcommands
- `Certificates/` - Certificate management
- `Notarization/` - Notarization commands
- `Profiles/` - Provisioning profiles
- `Reviews/` - App review management
- `Screenshots/` - Screenshot management (different from framing)
- `TestFlight/` - TestFlight operations
- `Users/` - User management
- `Versions/` - Version management

#### ASCConfig Library (not used)
- `ConfigFile.swift` - Config file reading/writing
- `Credential.swift` - Credential model with sources
- `DurationValue.swift` - Duration parsing for config

#### ASCCore Library (not used)
- `ASCCore.swift` - Core library types
- `HTTPClient.swift` - HTTP client implementation
- `OutputFormatters.swift` - JSON, table, markdown formatters

#### ASCMain (not used)
- `asc-swift/main.swift` - Full CLI entry point with all subcommands

## Why We're Not Using These

### Design Philosophy

The helper tools approach is different from a full Swift CLI replacement:

1. **Minimal Surface Area**: Only 3 helpers instead of 70+ commands
2. **Single Responsibility**: Each helper does one thing well (JWT, keychain, screenshots)
3. **Go Remains Primary**: Go CLI is the interface; Swift is just for performance-critical operations
4. **No State Management**: Go handles credentials, config, API clients; Swift handles compute

### What Go Still Does

- **Command parsing** (ffcli)
- **API client** (HTTP, retries, pagination)
- **Output formatting** (JSON, table, markdown)
- **Configuration** (profiles, defaults)
- **Authentication resolution** (env vars, config, keychain lookup)
- **Workflow orchestration**

### What Swift Helpers Do

- **JWT signing** (2-3x faster with CryptoKit hardware acceleration)
- **Keychain raw operations** (native Security.framework, no CGO)
- **Screenshot image processing** (Core Image/Metal, replaces Python)

## Performance Comparison

| Operation | swift-asc-cli (full) | swift-helpers | Go CLI |
|-----------|---------------------|---------------|---------|
| JWT Sign | ✅ CryptoKit | ✅ CryptoKit | golang-jwt (slower) |
| Keychain | ✅ Security.framework | ✅ Security.framework | 99designs/keyring (CGO) |
| HTTP Client | URLSession | ❌ Not needed | Go net/http |
| JSON Formatting | Custom | ❌ Not needed | Go encoding/json |
| Screenshot Frame | Core Image | ✅ Core Image | Python/Koubou |
| Command Parsing | ArgumentParser | ❌ Not needed | ffcli |

## Benefits of This Approach

1. **Incremental Adoption**: Can add Swift helpers one at a time
2. **Cross-Platform**: Go CLI works on Linux/Windows; helpers only on macOS
3. **Falls Back Gracefully**: Go implementations remain as fallbacks
4. **Testability**: Each helper is a standalone CLI tool
5. **No Lock-in**: Can swap Swift for Go or vice versa per-operation

## Future Considerations

If we need more operations later, we can add more helpers:
- `asc-image-compress` - Metal-accelerated image compression
- `asc-payload-sign` - Notary/signing operations
- `asc-simulator` - Simulator control (XCTest framework)

But these would be added as new executables, not as part of a full Swift CLI.

## File Inventory

### Copied & Adapted
```
swift-helpers/Sources/
├── asc-jwt-sign/main.swift      (adapted from ASCAuth/JWT.swift)
├── asc-keychain/main.swift       (adapted from ASCAuth/Keychain.swift)
└── asc-screenshot-frame/main.swift (new, Core Image based)
```

### Not Copied
```
swift-asc-cli/Sources/
├── ASCAuth/
│   ├── ASCAuth.swift             ❌ Not needed
│   ├── CredentialResolver.swift  ❌ Not needed
│   └── KeyLoader.swift           ❌ Not needed
├── ASCCommands/
│   ├── ASCCommands.swift         ❌ Not needed
│   ├── ASCCredentialTokenProvider.swift ❌ Not needed
│   ├── SharedCommandUtils.swift  ❌ Not needed
│   ├── Auth/                     ❌ Not needed
│   ├── Apps/                     ❌ Not needed
│   ├── Builds/                   ❌ Not needed
│   ├── Certificates/             ❌ Not needed
│   ├── Notarization/             ❌ Not needed
│   ├── Profiles/                 ❌ Not needed
│   ├── Reviews/                  ❌ Not needed
│   ├── Screenshots/              ❌ Not needed (different from framing)
│   ├── TestFlight/               ❌ Not needed
│   ├── Users/                    ❌ Not needed
│   └── Versions/                 ❌ Not needed
├── ASCConfig/
│   ├── ConfigFile.swift          ❌ Not needed
│   ├── Credential.swift          ❌ Not needed
│   └── DurationValue.swift       ❌ Not needed
├── ASCCore/
│   ├── ASCCore.swift             ❌ Not needed
│   ├── HTTPClient.swift          ❌ Not needed
│   └── OutputFormatters.swift    ❌ Not needed
└── asc-swift/
    └── main.swift                ❌ Not needed (full CLI entry point)
```

## Summary

- **Reused**: Core JWT signing logic, keychain error handling patterns
- **Adapted**: Simplified for single-purpose CLI tools
- **Not Used**: Full command implementations, HTTP clients, formatters, config management
- **Rationale**: Helpers are focused performance tools, not a full CLI replacement
