# 2026-04-07 Backend Legacy Cleanup Assessment

## Worktree

- Repository: `E:\Github\Qingyu\Qingyu_backend`
- Worktree: `E:\Github\Qingyu\Qingyu_backend\.worktrees\backend-legacy-cleanup`
- Branch: `chore/backend-legacy-cleanup`
- Base commit: `1ba33f92`

## Assessment Method

This assessment was built from:

1. `npx gitnexus analyze` after the root index was reported stale
2. `gitnexus query/context/impact` on backend cleanup candidates
3. `rg` cross-checks for import paths and direct source references

Notes:

- GitNexus is indexed at the root repository `E:\Github\Qingyu`, so all graph queries were narrowed to `Qingyu_backend/` symbols and files.
- For compatibility wrappers, graph results were cross-checked with direct import/reference search before classifying anything as dead code.

## Backend Cleanup Hotspots

Keyword scan over backend code only, excluding `Qingyu_backend/docs/**`:

| Bucket | Hits |
| --- | ---: |
| `service/shared` | 53 |
| `service/ai` | 49 |
| `test/integration` | 32 |
| `api/v1` | 30 |
| `service/writer` | 28 |
| `repository/mongodb` | 21 |
| `scripts/check-dependencies` | 18 |
| `service/admin` | 17 |
| `repository/search` | 16 |
| `models/shared` | 10 |
| `service/user` | 10 |

Top files:

| File | Hits |
| --- | ---: |
| `service/shared/stats/stats_service_test.go` | 31 |
| `service/ai/ai_service.go` | 17 |
| `service/shared/stats/stats_service.go` | 17 |
| `scripts/check-dependencies/main.go` | 14 |
| `service/ai/context_service.go` | 12 |
| `repository/search/elasticsearch_repository.go` | 9 |
| `repository/search/search_cache_repository.go` | 7 |

## GitNexus Blast Radius Summary

### Low-risk first-batch candidates

#### `repository/search/search_cache_repository.go:NewSearchCacheRepository`

- GitNexus impact: `LOW`
- Direct callers: `0`
- Processes affected: `0`
- Recommendation: candidate for archive/removal in the first cleanup batch

#### `repository/search/elasticsearch_repository.go:NewElasticsearchRepository`

- GitNexus impact: `LOW`
- Direct callers: `0`
- Processes affected: `0`
- Recommendation: candidate for archive/removal in the first cleanup batch

#### `service/ai/adapter/manager.go:NewAdapterManager`

- GitNexus impact: `LOW`
- Direct callers: `0`
- Important caveat: the `service/ai/adapter` package is still imported by runtime code such as:
  - `service/ai/ai_service.go`
  - `service/ai/ai_gateway.go`
  - `service/ai/chat_service.go`
  - `service/ai/proofread_service.go`
  - `service/ai/sensitive_words_service.go`
  - `service/ai/summarize_service.go`
- Recommendation: do not delete the package in batch 1; only target obviously unused constructors or move toward package-internal consolidation

### High-risk architecture cleanup

#### `service/shared/stats/stats_service.go:NewPlatformStatsService`

- GitNexus impact: `HIGH`
- Direct callers: `24`
- Affected process: `RegisterRoutes`
- Affected modules: `Stats`, `Social`, `Integration`
- Known runtime reference: `router/enter.go` wires this service into the main route setup
- Recommendation: do not move/remove this service in the first cleanup batch

## Important Cross-check Findings

### AI legacy service is only partially dead

`service/ai/ai_service.go:NewService` has no graph callers, but the legacy compatibility surface is not removable as a whole:

- `service_container.go:911` still calls `NewServiceWithDependencies`
- the returned `Service` type still feeds `chatService := aiService.NewChatService(c.aiService, chatRepo)`

Conclusion:

- `NewService` alone looks removable
- the legacy `Service` compatibility layer is still part of runtime wiring and must be handled in a later migration batch

### `service/shared/auth` compatibility cleanup status (updated)

This worktree has completed the compatibility-path cleanup for auth:

- test and edge validation imports were migrated from `Qingyu_backend/service/shared/auth` to current paths (`service/auth` and `service/user` where appropriate)
- `service/shared/auth/compat.go` was deleted
- script hardcoded paths were migrated (notably in `scripts/testing/mvp_integration_test.sh` and `scripts/deployment/deployment_check.sh`)

Current state:

- no non-doc code imports remain for `Qingyu_backend/service/shared/auth`
- residual mentions are in docs and intentional dependency-check rules/tests only

### `service/shared/messaging_compat.go` retirement status (updated)

Pre-edit blast radius:

- GitNexus `impact` on `MessagingService` and `NewMessagingService` returned `LOW`
- local reference scan found no non-doc code usage of `shared.NewMessagingService`, `shared.NewNotificationService`, or other messaging aliases exported from `service/shared`

Current state:

- `service/shared/messaging_compat.go` was deleted
- `service/shared/README.md` now points callers directly to `service/channels`

Conclusion:

- the messaging compatibility shim was dead and safe to remove in this worktree

### `repository/mongodb/auth/auth_repository_compat.go` retirement status (updated)

Pre-edit blast radius:

- GitNexus `impact` on `NewAuthRepository` returned `LOW`
- direct callers: `3`
- all direct callers were test-only references in `tests/integration/permission_integration_test.go`

Current state:

- `tests/integration/permission_integration_test.go` now uses `auth.NewRoleRepository`
- `repository/mongodb/auth/auth_repository_compat.go` was deleted
- package README snippets were updated to point at `NewRoleRepository`

Conclusion:

- the Mongo auth compatibility constructor was dead outside tests and safe to remove in this worktree

### `repository/interfaces/shared` storage aliases retirement status (updated)

Pre-edit blast radius:

- GitNexus `impact` on `StorageRepository` returned `LOW`
- GitNexus `impact` on `FileFilter` returned `LOW`
- direct code references were limited to `service/shared/storage/*` and `repository/mongodb/storage/storage_repository_test.go`

Current state:

- `service/shared/storage/*` now imports `repository/interfaces/storage` directly
- `repository/mongodb/storage/storage_repository_test.go` now uses `storage.FileFilter` and `storage.StorageRepository`
- `StorageRepository` and `FileFilter` aliases were deleted from `repository/interfaces/shared/shared_repository.go`

Conclusion:

- the storage compatibility aliases are retired
- the only remaining production interface under `repository/interfaces/shared` is `TokenBlacklistRepository`, which is handled separately below

### `repository/interfaces/shared` recommendation alias retirement status (updated)

Pre-edit blast radius:

- GitNexus `impact` on `RecommendationRepository` returned `LOW`
- direct code references were limited to `service/recommendation/recommendation_service.go`, `repository/interfaces/RepoFactory_interface.go`, and `repository/mongodb/factory.go`

Current state:

- `service/recommendation/recommendation_service.go` now imports `repository/interfaces/recommendation` directly
- `repository/interfaces/RepoFactory_interface.go` now returns `recommendation.RecommendationRepository` from `CreateRecommendationRepository`
- `repository/mongodb/factory.go` now returns `recommendation.RecommendationRepository` from `CreateRecommendationRepository`
- `RecommendationRepository` alias was deleted from `repository/interfaces/shared/shared_repository.go`

Conclusion:

- the recommendation compatibility alias is retired

### `repository/interfaces/shared` messaging aliases retirement status (updated)

Pre-edit blast radius:

- GitNexus `impact` on `MessageRepository` returned `LOW`
- GitNexus `impact` on `NotificationFilter` returned `LOW`
- direct code references were limited to `service/channels/notification_service_complete.go`

Current state:

- `service/channels/notification_service_complete.go` now imports `repository/interfaces/messaging` directly
- `MessageRepository`, `MessageFilter`, and `NotificationFilter` aliases were deleted from `repository/interfaces/shared/shared_repository.go`

Conclusion:

- the messaging compatibility aliases are retired

### `repository/interfaces/shared` finance aliases retirement status (updated)

Pre-edit blast radius:

- GitNexus `impact` on `WalletRepository` returned `LOW`
- GitNexus `impact` on `TransactionFilter` returned `LOW`
- GitNexus `impact` on `WithdrawFilter` returned `LOW`
- direct code references were limited to `service/finance/*` and `repository/interfaces/shared/mocks/mock_wallet_repository.go`

Current state:

- `service/finance/*` now imports `repository/interfaces/finance` directly for wallet repository and filter types
- `repository/interfaces/shared/mocks/mock_wallet_repository.go` now uses `finance.TransactionFilter` and `finance.WithdrawFilter`
- `WalletRepository`, `TransactionFilter`, and `WithdrawFilter` aliases were deleted from `repository/interfaces/shared/shared_repository.go`

Conclusion:

- the finance compatibility aliases are retired

### `repository/interfaces/shared/shared_repository.go` retirement status (updated)

Pre-edit blast radius:

- `CreateAuthRepository` GitNexus `impact` returned `LOW`
- runtime callers for the final auth alias were concentrated in `service/container`, `service/auth/*`, and `service/user/*`
- the `shared` package also still contained `token_blacklist_repository.go`, so only the alias file itself was removable, not the whole package

Current state:

- `service/auth/*` now imports `repository/interfaces/auth` directly for `RoleRepository`
- `service/user/user_service.go` and `service/user/verification_service.go` now import `repository/interfaces/auth` directly for `RoleRepository`
- `repository/interfaces/RepoFactory_interface.go` and `repository/mongodb/factory.go` now return `auth.RoleRepository` from `CreateAuthRepository`
- `repository/interfaces/shared/shared_repository.go` was deleted
- `repository/interfaces/shared/` remains only because `token_blacklist_repository.go` still defines `TokenBlacklistRepository`

Conclusion:

- the shared alias file is retired
- the `repository/interfaces/shared` package is no longer a compatibility alias hub; it only remains as the current home of `TokenBlacklistRepository`

### `TokenBlacklistRepository` retirement status (updated in phase 2)

Fresh scan before phase-2 edit:

- earlier GitNexus output had classified the symbol as `HIGH` risk while it still sat under the legacy `repository/interfaces/shared` package
- current local reference scan over non-doc code showed no active runtime callers of:
  - `repository/interfaces/shared.TokenBlacklistRepository`
  - `repository/redis.NewTokenBlacklistRepository`
  - `repository/redis.NewTokenBlacklistRepositoryWithConfig`
- active auth runtime wiring already uses:
  - `service/auth.RedisAdapter`
  - `service/auth.InMemoryTokenBlacklist`

Current state:

- `repository/interfaces/shared/token_blacklist_repository.go` was deleted
- `repository/redis/token_blacklist_repository_redis.go` and its self-referential tests were deleted
- unused legacy mocks under `repository/interfaces/shared/mocks` were deleted
- new targeted tests were added around the live auth blacklist path in `service/auth/jwt_service_test.go`
- `service/user/user_service.go` now delegates both `LogoutUser` and `ValidateToken` to the live auth token lifecycle path via a narrow injected adapter instead of returning placeholder results

Conclusion:

- phase 2 does not migrate this interface; it retires the dead repository abstraction entirely
- `repository/interfaces/shared` is no longer needed as a live production package after this cleanup
- user-facing token lifecycle entry points now align with the active auth runtime path

### `service/writer/project` request DTO compatibility retirement status (updated)

Pre-edit blast radius:

- GitNexus on backend writer request symbols was noisy because of cross-repo name collisions, so the effective blast radius was verified with `rg`
- direct backend references were limited to `api/v1/writer/project_api.go`, `service/writer/project/project_service.go`, `service/writer/impl/project_management_impl.go`, and `service/writer/project/project_service_simple_test.go`
- risk level was treated as `MEDIUM` because the slice crossed API binding, service signatures, implementation adapters, and tests

Current state:

- `api/v1/writer/project_api.go` now passes `dto.CreateProjectRequest`, `dto.ListProjectsRequest`, and `dto.UpdateProjectRequest` directly
- `service/writer/project/project_service.go` now accepts `dto` request types directly
- `service/writer/impl/project_management_impl.go` now builds `dto` request types instead of `service/writer/project` aliases
- `service/writer/project/project_service_simple_test.go` now uses `dto` request types directly
- `CreateProjectRequest`, `ListProjectsRequest`, and `UpdateProjectRequest` aliases were deleted from `service/writer/project/project_dto.go`

Conclusion:

- the project request DTO compatibility aliases are retired
- `service/writer/project/project_dto.go` now only carries project-specific response compatibility types

### `service/writer/document` request DTO compatibility retirement status (updated)

Pre-edit blast radius:

- GitNexus on document request symbol names was partially polluted by frontend symbols, so backend scope was verified with `rg`
- `CreateDocumentRequest` / `UpdateDocumentRequest` backend references were limited to `api/v1/writer/document_api.go` and `service/writer/impl/document_management_impl.go`
- `ListDocumentsRequest` added one more direct backend reference in `service/writer/document/document_service.go`
- `MoveDocumentRequest` and `ReorderDocumentsRequest` were limited to `api/v1/writer/document_api.go` and `service/writer/impl/document_management_impl.go`
- `AutoSaveRequest`, `UpdateContentRequest`, and `ReplaceDocumentContentsRequest` were limited to `api/v1/writer/editor_api.go`, `api/v1/writer/document_api.go`, and `service/writer/impl/document_management_impl.go`
- `DocumentTreeNode` / `DocumentTreeResponse` remain live in `service/writer/document/document_service.go` and `service/writer/impl/document_management_impl.go`, so tree-response types were explicitly excluded from this slice

Current state:

- `api/v1/writer/document_api.go` now uses `dto.CreateDocumentRequest`, `dto.UpdateDocumentRequest`, `dto.ListDocumentsRequest`, `dto.MoveDocumentRequest`, `dto.ReorderDocumentsRequest`, and `dto.UpdateContentRequest` directly
- `api/v1/writer/editor_api.go` now uses `dto.AutoSaveRequest`, `dto.UpdateContentRequest`, and `dto.ReplaceDocumentContentsRequest` directly
- `service/writer/impl/document_management_impl.go` now builds `dto.CreateDocumentRequest`, `dto.UpdateDocumentRequest`, `dto.ListDocumentsRequest`, `dto.MoveDocumentRequest`, `dto.ReorderDocumentsRequest`, `dto.AutoSaveRequest`, and `dto.UpdateContentRequest`
- `service/writer/document/document_service.go` now accepts `dto.ListDocumentsRequest` directly
- request aliases were deleted from `service/writer/document/document_dto.go` for:
  - `CreateDocumentRequest`
  - `UpdateDocumentRequest`
  - `ListDocumentsRequest`
  - `MoveDocumentRequest`
  - `ReorderDocumentsRequest`
  - `AutoSaveRequest`
  - `UpdateContentRequest`
  - `ReplaceDocumentContentsRequest`

Conclusion:

- the document request DTO compatibility aliases for create, update, list, move, reorder, autosave, content update, and content replace are retired
- `service/writer/document/document_dto.go` now carries tree-response and response-side compatibility types only

## Recommended Cleanup Phases

### Phase A: safe cleanup in this worktree

Scope:

- remove or archive dead search repository stubs:
  - `repository/search/elasticsearch_repository.go`
  - `repository/search/search_cache_repository.go`
- clean low-value TODO placeholders around those dead paths
- move obviously obsolete example/demo files only after reference verification

Expected risk:

- low
- minimal runtime impact

### Phase B: compatibility path migration

Scope:

- keep `service/shared/auth` retired (code migration completed in this worktree)
- clean residual docs that still describe `service/shared/auth` as a live path
- migrate callers away from other deprecated compat-only paths
- reduce `_migration` surfaces that exist only for backward compatibility

Expected risk:

- medium
- test breakage likely if done incompletely

### Phase C: architecture cleanup

Scope:

- `service/shared/stats`
- writer DTO compatibility layers:
  - `service/interfaces/writer/dto.go`
  - `service/writer/document/document_dto.go`
- AI legacy `Service` compatibility layer
- `TokenBlacklistRepository` relocation / final `repository/interfaces/shared` retirement

Expected risk:

- high
- route wiring and integration test impact is likely

## Suggested First PR Slice

Keep the first PR in this worktree intentionally narrow:

1. delete/archive the two unused `repository/search` stub files
2. update any now-unused imports or documentation references inside backend code
3. run search-related tests or compile checks
4. produce a second impact pass before touching compatibility layers

This gives the cleanup branch a safe opening commit and avoids mixing dead-code pruning with deeper architecture migration.
