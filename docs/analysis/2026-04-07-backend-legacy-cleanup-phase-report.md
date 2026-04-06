# Backend Legacy Cleanup Phase Report

Date: 2026-04-07
Branch: `chore/backend-legacy-cleanup`
Scope: backend legacy cleanup phase 1

## Phase Boundary

This phase closes a low-risk cleanup slice and intentionally avoids broad architectural rewiring.

Included:

- retirement of legacy compatibility aliases under `repository/interfaces/shared/shared_repository.go`
- auth repository migration from shared compatibility aliases to `repository/interfaces/auth.RoleRepository`
- writer `project/document` request DTO compatibility retirement
- cleanup of dead search compatibility stubs
- synchronization of core architecture and migration documents

Excluded:

- relocating `TokenBlacklistRepository` out of `repository/interfaces/shared`
- `service/shared/stats` refactor
- removal of the AI legacy `Service` compatibility layer

## Delivered Changes

### 1. Shared compatibility alias retirement

- deleted `repository/interfaces/shared/shared_repository.go`
- retired storage aliases:
  - `StorageRepository`
  - `FileFilter`
- retired recommendation alias:
  - `RecommendationRepository`
- retired messaging aliases:
  - `MessageRepository`
  - `MessageFilter`
  - `NotificationFilter`
- retired finance aliases:
  - `WalletRepository`
  - `TransactionFilter`
  - `WithdrawFilter`

### 2. Auth repository migration

- `CreateAuthRepository()` now returns `auth.RoleRepository`
- `service/auth/*` now imports `repository/interfaces/auth` directly
- `service/user/*` now imports `repository/interfaces/auth` directly
- legacy auth compatibility files were removed where no longer needed

### 3. Writer DTO cleanup

- retired request DTO aliases from `service/writer/project/project_dto.go`
- retired request DTO aliases from `service/writer/document/document_dto.go`
- API, service, implementation, and tests now use `models/dto` request types directly
- document tree response types remain in place because they are still active module-specific DTOs

### 4. Dead file cleanup

- deleted `repository/search/elasticsearch_repository.go`
- deleted `repository/search/search_cache_repository.go`
- deleted `repository/mongodb/auth/auth_repository_compat.go`
- deleted `service/shared/auth/compat.go`
- deleted `service/shared/messaging_compat.go`

### 5. Documentation sync

- updated central repository and service architecture documents
- updated auth and user module READMEs to reflect `RoleRepository`
- updated cleanup assessment with the current deferred/high-risk boundary

## Verification

The following checks passed in this worktree:

```bash
go test -work ./service/auth ./service/user ./service/container ./router/user ./repository/mongodb ./repository/interfaces/... -run '^$' -count=1
go test -work ./api/v1/writer ./service/writer/document ./service/writer/impl ./router/writer -run '^$' -count=1
go test -work ./service/writer/project -run 'TestProjectService_(CreateProject|UpdateProject|ListMyProjects)' -count=1
go test -work ./repository/redis -run '^TestTokenBlacklistRepository' -count=1
```

Notes:

- `-work` was required because this Windows environment intermittently fails while cleaning temporary Go build directories, even when package compilation and tests themselves succeed
- an unrelated pre-existing document test instability remains outside this phase boundary: `service/writer/document` full test suite still has a known assertion mismatch in `TestVerifyDocumentEdit_Success`

## Deferred Items

### `TokenBlacklistRepository`

- GitNexus `impact` on `TokenBlacklistRepository` returned `HIGH`
- the interface still lives under the legacy `repository/interfaces/shared` package, so a package-path move would fan out across many `shared` imports
- this phase leaves `repository/interfaces/shared/token_blacklist_repository.go` in place and documents it as the only remaining production interface in that package

## Recommended Next Directions

1. isolate `TokenBlacklistRepository` migration into a dedicated PR, with explicit package-path updates and a smaller caller set
2. audit and retire remaining legacy mocks under `repository/interfaces/shared/mocks`
3. continue writer response-side DTO cleanup only after separating active module-specific types from true compatibility aliases
4. re-evaluate the AI legacy service wrapper and `service/shared/stats` only after a fresh impact scan confirms a bounded blast radius
