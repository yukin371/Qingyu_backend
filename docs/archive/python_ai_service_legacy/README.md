# Python AI Service Legacy Archive

## Overview

This directory contains archived documentation and references for the legacy `python_ai_service` that was previously part of `Qingyu_backend`.

## Migration Status

**Status**: ✅ **Completed** - Migrated to standalone `Qingyu-Ai-Service`

The Python AI service has been successfully migrated to a standalone microservice architecture:

- **New Location**: `Qingyu-Ai-Service/` (separate repository)
- **Communication**: gRPC protocol
- **Protocol Version**: v1.1.0

## Migration Timeline

- **2025-01-24**: Phase 3 Backend Refactoring completed
  - ✅ Implemented AI error types (`pkg/errors/ai_errors.go`)
  - ✅ Implemented circuit breaker (`pkg/circuitbreaker/circuit_breaker.go`)
  - ✅ Marked adapter package as deprecated
  - ✅ Simplified AIService with gRPC client
  - ✅ Removed `python_ai_service` directory from backend

## Key Changes

### Before (Legacy)

```
Qingyu_backend/
├── python_ai_service/          # ❌ Removed
│   ├── src/
│   ├── tests/
│   ├── docker/
│   └── proto/
└── service/ai/
    ├── adapter/                # ❌ Deprecated
    └── ai_service.go           # ❌ Replaced
```

### After (New Architecture)

```
Qingyu-Ai-Service/              # ✅ Standalone service
├── src/
│   ├── agents/
│   ├── grpc_service/
│   └── services/
├── docker/
└── proto/

Qingyu_backend/
├── pkg/
│   ├── errors/
│   │   └── ai_errors.go        # ✅ New
│   └── circuitbreaker/
│       └── circuit_breaker.go  # ✅ New
└── service/ai/
    ├── adapter/                # ⚠️ Deprecated (emergency fallback)
    ├── ai_service_v2.go        # ✅ New gRPC-based service
    ├── grpc_client.go          # ✅ gRPC client
    └── quota_service.go        # ✅ Quota management
```

## Migration Guide

### For Backend Developers

**Old code**:
```go
adapter := ai.NewOpenAIAdapter(...)
result, err := adapter.GenerateText(ctx, req)
```

**New code**:
```go
client := ai.NewGRPCClient(conn)
result, err := client.ExecuteAgent(ctx, req)
```

### Configuration

**Environment Variables**:
- `AI_SERVICE_ENDPOINT`: Qingyu-Ai-Service gRPC endpoint (default: `localhost:50051`)
- `AI_ENABLE_FALLBACK`: Enable emergency fallback to legacy adapters (default: `false`)

## Documentation References

- **Migration Design**: `docs/plans/2026-01-24-ai-service-complete-migration-design.md`
- **Implementation**: `docs/plans/2026-01-24-ai-service-migration-implementation.md`
- **gRPC Integration Guide**: `docs/ai/GRPC_INTEGRATION_GUIDE.md`

## Rollback Plan

If needed, the legacy adapter code is still available in `service/ai/adapter/` (marked as deprecated) for emergency fallback:

1. Set `AI_ENABLE_FALLBACK=true` in environment
2. Restart backend service
3. The system will fall back to direct API calls

**Note**: This is temporary and will be removed in v2.0.0

## Archive Contents

This directory contains:
- Historical documentation
- Migration references
- Rollback instructions

## Contact

For questions or issues related to the migration, refer to:
- Project Architecture: `docs/architecture/`
- Migration Planning: `docs/plans/`
- AI Service Documentation: `docs/ai/`

---

**Last Updated**: 2025-01-24  
**Migration Version**: v1.1.0  
**Status**: ✅ Complete
