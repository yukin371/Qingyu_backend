# Plans Directory Migration Notice

`Qingyu_backend/docs/plans` has been migrated to the parent repository:

- [docs/plans/submodules/backend/](../../../docs/plans/submodules/backend/README.md)

## Effective Rule

Do not add new backend plan/design documents in this submodule directory.

Create and maintain all new backend planning/design docs in:

- [docs/plans/submodules/backend/](../../../docs/plans/submodules/backend/README.md)

Use the parent taxonomy instead of recreating a local plan tree:

- `architecture/`: current architecture and model/versioning design
- `api-governance/`: API standardization, integration, and error handling governance
- `publication/`: publication/review workflow design
- `shared-and-layering/`: shared module and layering refactors
- `testing-and-quality/`: test strategy and quality plans
- `legacy-phases/`: historical rollout plans and completion reports

## Why

- centralizes cross-repo planning
- avoids duplicated plan trees between parent repo and submodule
- keeps backend implementation repo focused on code and local operational docs
