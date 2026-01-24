package adapter

// Deprecated Package Notice
//
// This package (adapter) is deprecated and will be removed in v2.0.0.
//
// Migration Path:
// 1. Use service/ai/grpc_client.go to call Qingyu-Ai-Service
// 2. Replace adapter calls with gRPC client calls
// 3. Enable emergency fallback with AI_ENABLE_FALLBACK=true if needed
//
// Migration Example:
//
// Old code:
//   adapter := ai.NewOpenAIAdapter(...)
//   result, err := adapter.GenerateText(ctx, req)
//
// New code:
//   client := ai.NewGRPCClient(...)
//   result, err := client.ExecuteAgent(ctx, req)
//
// For questions, refer to: docs/plans/2026-01-24-ai-service-complete-migration-design.md
//
// Deprecated: Use Qingyu-Ai-Service gRPC API instead.
// This adapter manager is kept only for emergency fallback.
// Will be removed in v2.0.0
