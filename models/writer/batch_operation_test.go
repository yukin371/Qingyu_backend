package writer

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBatchOperation_TouchForCreate(t *testing.T) {
	op := &BatchOperation{
		ProjectID: primitive.NewObjectID(),
		Type:      BatchOpTypeDelete,
		TargetIDs: []string{"doc-1", "doc-2"},
		Atomic:    true,
		CreatedBy: primitive.NewObjectID(),
	}

	op.TouchForCreate()

	if op.ID.IsZero() {
		t.Error("ID should be generated")
	}
	if op.Status != BatchOpStatusPending {
		t.Errorf("Status should be pending, got %s", op.Status)
	}
	if !op.Cancelable {
		t.Error("Should be cancelable by default")
	}
	if op.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if op.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

func TestBatchOperation_TouchForCreate_WithExistingStatus(t *testing.T) {
	op := &BatchOperation{
		ProjectID: primitive.NewObjectID(),
		Type:      BatchOpTypeDelete,
		TargetIDs: []string{"doc-1"},
		Status:    BatchOpStatusRunning,
		CreatedBy: primitive.NewObjectID(),
	}

	op.TouchForCreate()

	// Status should not be overwritten if already set
	if op.Status != BatchOpStatusRunning {
		t.Errorf("Status should remain running, got %s", op.Status)
	}
}

func TestBatchOperation_IsRunning(t *testing.T) {
	tests := []struct {
		name   string
		status BatchOperationStatus
		want   bool
	}{
		{"pending", BatchOpStatusPending, false},
		{"running", BatchOpStatusRunning, true},
		{"completed", BatchOpStatusCompleted, false},
		{"failed", BatchOpStatusFailed, false},
		{"cancelled", BatchOpStatusCancelled, false},
		{"partially_failed", BatchOpStatusPartiallyFailed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &BatchOperation{Status: tt.status}
			if got := op.IsRunning(); got != tt.want {
				t.Errorf("IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatchOperation_IsTerminal(t *testing.T) {
	tests := []struct {
		name   string
		status BatchOperationStatus
		want   bool
	}{
		{"pending", BatchOpStatusPending, false},
		{"running", BatchOpStatusRunning, false},
		{"completed", BatchOpStatusCompleted, true},
		{"failed", BatchOpStatusFailed, true},
		{"cancelled", BatchOpStatusCancelled, true},
		{"partially_failed", BatchOpStatusPartiallyFailed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &BatchOperation{Status: tt.status}
			if got := op.IsTerminal(); got != tt.want {
				t.Errorf("IsTerminal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatchOperation_CanCancel(t *testing.T) {
	tests := []struct {
		name       string
		status     BatchOperationStatus
		cancelable bool
		want       bool
	}{
		{"running cancelable", BatchOpStatusRunning, true, true},
		{"running not cancelable", BatchOpStatusRunning, false, false},
		{"completed cancelable", BatchOpStatusCompleted, true, false},
		{"pending cancelable", BatchOpStatusPending, true, false},
		{"failed cancelable", BatchOpStatusFailed, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &BatchOperation{
				Status:     tt.status,
				Cancelable: tt.cancelable,
			}
			if got := op.CanCancel(); got != tt.want {
				t.Errorf("CanCancel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatchOperationItem_TouchForCreate(t *testing.T) {
	item := &BatchOperationItem{
		BatchID:        primitive.NewObjectID(),
		TargetID:       "doc-1",
		TargetStableRef: "stable-ref-123",
	}

	item.TouchForCreate()

	if item.ID.IsZero() {
		t.Error("ID should be generated")
	}
	if item.Status != BatchItemStatusPending {
		t.Errorf("Status should be pending, got %s", item.Status)
	}
	if item.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if item.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

func TestBatchOperationItem_TouchForCreate_WithExistingStatus(t *testing.T) {
	item := &BatchOperationItem{
		BatchID:        primitive.NewObjectID(),
		TargetID:       "doc-1",
		TargetStableRef: "stable-ref-123",
		Status:         BatchItemStatusProcessing,
	}

	item.TouchForCreate()

	// Status should not be overwritten if already set
	if item.Status != BatchItemStatusProcessing {
		t.Errorf("Status should remain processing, got %s", item.Status)
	}
}

func TestOperationLog_TouchForCreate(t *testing.T) {
	log := &OperationLog{
		ProjectID:   primitive.NewObjectID(),
		UserID:      primitive.NewObjectID(),
		CommandType: CommandDelete,
		TargetIDs:   []string{"doc-1"},
	}

	log.TouchForCreate()

	if log.ID.IsZero() {
		t.Error("ID should be generated")
	}
	if log.ChainID == "" {
		t.Error("ChainID should be set")
	}
	if log.ChainID != log.ID.Hex() {
		t.Error("ChainID should equal ID when not set")
	}
	if log.Status != OpLogStatusExecuted {
		t.Errorf("Status should be executed, got %s", log.Status)
	}
	if log.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if log.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

func TestOperationLog_TouchForCreate_WithChainID(t *testing.T) {
	log := &OperationLog{
		ProjectID:   primitive.NewObjectID(),
		UserID:      primitive.NewObjectID(),
		CommandType: CommandDelete,
		TargetIDs:   []string{"doc-1"},
		ChainID:     "existing-chain-id",
	}

	log.TouchForCreate()

	// ChainID should not be overwritten if already set
	if log.ChainID != "existing-chain-id" {
		t.Errorf("ChainID should remain existing-chain-id, got %s", log.ChainID)
	}
}

func TestOperationLog_TouchForCreate_WithExistingStatus(t *testing.T) {
	log := &OperationLog{
		ProjectID:   primitive.NewObjectID(),
		UserID:      primitive.NewObjectID(),
		CommandType: CommandDelete,
		TargetIDs:   []string{"doc-1"},
		Status:      OpLogStatusUndone,
	}

	log.TouchForCreate()

	// Status should not be overwritten if already set
	if log.Status != OpLogStatusUndone {
		t.Errorf("Status should remain undone, got %s", log.Status)
	}
}

func TestOperationLog_IsUndoable(t *testing.T) {
	tests := []struct {
		name       string
		status     OperationLogStatus
		committed  bool
		want       bool
	}{
		{"executed and committed", OpLogStatusExecuted, true, true},
		{"executed but not committed", OpLogStatusExecuted, false, false},
		{"undone", OpLogStatusUndone, true, false},
		{"redone", OpLogStatusRedone, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &OperationLog{
				Status:      tt.status,
				IsCommitted: tt.committed,
			}
			if got := log.IsUndoable(); got != tt.want {
				t.Errorf("IsUndoable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationLog_IsRedoable(t *testing.T) {
	tests := []struct {
		name           string
		status         OperationLogStatus
		inverseCommand map[string]interface{}
		want           bool
	}{
		{"undone with inverse", OpLogStatusUndone, map[string]interface{}{"test": "value"}, true},
		{"undone without inverse", OpLogStatusUndone, nil, false},
		{"executed with inverse", OpLogStatusExecuted, map[string]interface{}{"test": "value"}, false},
		{"redone with inverse", OpLogStatusRedone, map[string]interface{}{"test": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &OperationLog{
				Status:         tt.status,
				InverseCommand: tt.inverseCommand,
			}
			if got := log.IsRedoable(); got != tt.want {
				t.Errorf("IsRedoable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationLog_IsRedoable_EmptyInverseCommand(t *testing.T) {
	log := &OperationLog{
		Status:         OpLogStatusUndone,
		InverseCommand: map[string]interface{}{}, // Empty map
	}

	// Empty map is not nil, so it should be considered as having an inverse command
	// This test verifies the behavior with empty inverse command
	if got := log.IsRedoable(); !got {
		t.Error("IsRedoable() should return true even when InverseCommand is empty (not nil)")
	}
}

func TestBatchOperationType_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value BatchOperationType
	}{
		{"delete", BatchOpTypeDelete},
		{"move", BatchOpTypeMove},
		{"copy", BatchOpTypeCopy},
		{"export", BatchOpTypeExport},
		{"apply_template", BatchOpTypeApplyTemplate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("BatchOperationType %s should not be empty", tt.name)
			}
		})
	}
}

func TestDocumentCommandType_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value DocumentCommandType
	}{
		{"create", CommandCreate},
		{"update", CommandUpdate},
		{"move", CommandMove},
		{"copy", CommandCopy},
		{"delete", CommandDelete},
		{"restore", CommandRestore},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("DocumentCommandType %s should not be empty", tt.name)
			}
		})
	}
}

func TestBatchOperationStatus_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value BatchOperationStatus
	}{
		{"pending", BatchOpStatusPending},
		{"running", BatchOpStatusRunning},
		{"completed", BatchOpStatusCompleted},
		{"failed", BatchOpStatusFailed},
		{"cancelled", BatchOpStatusCancelled},
		{"partially_failed", BatchOpStatusPartiallyFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("BatchOperationStatus %s should not be empty", tt.name)
			}
		})
	}
}

func TestBatchItemStatus_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value BatchItemStatus
	}{
		{"pending", BatchItemStatusPending},
		{"processing", BatchItemStatusProcessing},
		{"succeeded", BatchItemStatusSucceeded},
		{"failed", BatchItemStatusFailed},
		{"skipped", BatchItemStatusSkipped},
		{"cancelled", BatchItemStatusCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("BatchItemStatus %s should not be empty", tt.name)
			}
		})
	}
}

func TestOperationLogStatus_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value OperationLogStatus
	}{
		{"executed", OpLogStatusExecuted},
		{"undone", OpLogStatusUndone},
		{"redone", OpLogStatusRedone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("OperationLogStatus %s should not be empty", tt.name)
			}
		})
	}
}

func TestExecutionMode_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value ExecutionMode
	}{
		{"standard_atomic", ExecutionModeStandardAtomic},
		{"saga_atomic", ExecutionModeSagaAtomic},
		{"non_atomic", ExecutionModeNonAtomic},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("ExecutionMode %s should not be empty", tt.name)
			}
		})
	}
}

func TestConflictPolicy_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value ConflictPolicy
	}{
		{"abort", ConflictPolicyAbort},
		{"overwrite", ConflictPolicyOverwrite},
		{"skip", ConflictPolicySkip},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("ConflictPolicy %s should not be empty", tt.name)
			}
		})
	}
}

func TestPreflightSummary_Structure(t *testing.T) {
	summary := PreflightSummary{
		TotalCount:   100,
		ValidCount:   80,
		InvalidCount: 15,
		SkippedCount: 5,
	}

	if summary.TotalCount != 100 {
		t.Errorf("TotalCount = %d, want 100", summary.TotalCount)
	}
	if summary.ValidCount != 80 {
		t.Errorf("ValidCount = %d, want 80", summary.ValidCount)
	}
	if summary.InvalidCount != 15 {
		t.Errorf("InvalidCount = %d, want 15", summary.InvalidCount)
	}
	if summary.SkippedCount != 5 {
		t.Errorf("SkippedCount = %d, want 5", summary.SkippedCount)
	}
}

func TestBatchOperation_TimeFields(t *testing.T) {
	now := time.Now()
	op := &BatchOperation{
		Status:    BatchOpStatusRunning,
		StartedAt: &now,
	}

	if !op.IsRunning() {
		t.Error("Operation should be running")
	}

	later := now.Add(time.Hour)
	op.FinishedAt = &later
	op.Status = BatchOpStatusCompleted

	if op.IsRunning() {
		t.Error("Operation should not be running after completion")
	}
}

func TestBatchOperationItem_TimeFields(t *testing.T) {
	now := time.Now()
	item := &BatchOperationItem{
		Status:     BatchItemStatusProcessing,
		StartedAt:  &now,
	}

	if item.Status != BatchItemStatusProcessing {
		t.Errorf("Status = %s, want processing", item.Status)
	}

	later := now.Add(time.Minute)
	item.FinishedAt = &later
	item.Status = BatchItemStatusSucceeded

	if item.Status != BatchItemStatusSucceeded {
		t.Errorf("Status = %s, want succeeded", item.Status)
	}
}

func TestOperationLog_TimeFields(t *testing.T) {
	now := time.Now()
	log := &OperationLog{
		Status:      OpLogStatusExecuted,
		IsCommitted: true,
	}

	if !log.IsUndoable() {
		t.Error("Should be undoable when status is executed and committed")
	}

	// Change to undone status and add inverse command
	log.Status = OpLogStatusUndone
	log.UndoneAt = &now
	log.InverseCommand = map[string]interface{}{"test": "value"}
	if !log.IsRedoable() {
		t.Error("Should be redoable when status is undone and has inverse command")
	}
}

func TestBatchOperation_WithPreflightSummary(t *testing.T) {
	summary := &PreflightSummary{
		TotalCount:   50,
		ValidCount:   45,
		InvalidCount: 3,
		SkippedCount: 2,
	}

	op := &BatchOperation{
		ProjectID:        primitive.NewObjectID(),
		Type:             BatchOpTypeDelete,
		TargetIDs:        []string{"doc-1", "doc-2"},
		PreflightSummary: summary,
	}

	if op.PreflightSummary == nil {
		t.Error("PreflightSummary should not be nil")
	}

	if op.PreflightSummary.TotalCount != 50 {
		t.Errorf("PreflightSummary.TotalCount = %d, want 50", op.PreflightSummary.TotalCount)
	}
}

func TestBatchOperationItem_WithVersionControl(t *testing.T) {
	expectedVersion := 1
	actualVersion := 2

	item := &BatchOperationItem{
		BatchID:         primitive.NewObjectID(),
		TargetID:        "doc-1",
		TargetStableRef: "stable-ref-123",
		ExpectedVersion: &expectedVersion,
		ActualVersion:   &actualVersion,
	}

	if item.ExpectedVersion == nil {
		t.Error("ExpectedVersion should not be nil")
	}

	if *item.ExpectedVersion != 1 {
		t.Errorf("ExpectedVersion = %d, want 1", *item.ExpectedVersion)
	}

	if *item.ActualVersion != 2 {
		t.Errorf("ActualVersion = %d, want 2", *item.ActualVersion)
	}
}

func TestOperationLog_WithBatchOpID(t *testing.T) {
	batchOpID := primitive.NewObjectID()

	log := &OperationLog{
		ProjectID:   primitive.NewObjectID(),
		UserID:      primitive.NewObjectID(),
		CommandType: CommandDelete,
		TargetIDs:   []string{"doc-1"},
		BatchOpID:   &batchOpID,
	}

	if log.BatchOpID == nil {
		t.Error("BatchOpID should not be nil")
	}

	if *log.BatchOpID != batchOpID {
		t.Errorf("BatchOpID = %s, want %s", log.BatchOpID.Hex(), batchOpID.Hex())
	}
}

func TestOperationLog_WithCommandPayload(t *testing.T) {
	payload := map[string]interface{}{
		"reason": "user_request",
		"cascade": true,
	}

	log := &OperationLog{
		ProjectID:      primitive.NewObjectID(),
		UserID:         primitive.NewObjectID(),
		CommandType:    CommandDelete,
		TargetIDs:      []string{"doc-1"},
		CommandPayload: payload,
	}

	if log.CommandPayload == nil {
		t.Error("CommandPayload should not be nil")
	}

	if log.CommandPayload["reason"] != "user_request" {
		t.Errorf("CommandPayload[reason] = %v, want user_request", log.CommandPayload["reason"])
	}

	if log.CommandPayload["cascade"] != true {
		t.Errorf("CommandPayload[cascade] = %v, want true", log.CommandPayload["cascade"])
	}
}

func TestBatchOperationItem_WithInverseCommand(t *testing.T) {
	inverseCommand := map[string]interface{}{
		"action": "restore",
		"data":   "original_data",
	}

	item := &BatchOperationItem{
		BatchID:         primitive.NewObjectID(),
		TargetID:        "doc-1",
		TargetStableRef: "stable-ref-123",
		InverseCommand:  inverseCommand,
	}

	if item.InverseCommand == nil {
		t.Error("InverseCommand should not be nil")
	}

	if item.InverseCommand["action"] != "restore" {
		t.Errorf("InverseCommand[action] = %v, want restore", item.InverseCommand["action"])
	}
}

func TestBatchOperation_WithClientRequestID(t *testing.T) {
	clientRequestID := "client-request-123"

	op := &BatchOperation{
		ProjectID:       primitive.NewObjectID(),
		Type:            BatchOpTypeDelete,
		TargetIDs:       []string{"doc-1"},
		ClientRequestID: clientRequestID,
	}

	if op.ClientRequestID != clientRequestID {
		t.Errorf("ClientRequestID = %s, want %s", op.ClientRequestID, clientRequestID)
	}
}

func TestBatchOperation_WithExpectedVersions(t *testing.T) {
	expectedVersions := map[string]int{
		"doc-1": 1,
		"doc-2": 2,
		"doc-3": 1,
	}

	op := &BatchOperation{
		ProjectID:        primitive.NewObjectID(),
		Type:             BatchOpTypeDelete,
		TargetIDs:        []string{"doc-1", "doc-2", "doc-3"},
		ExpectedVersions: expectedVersions,
	}

	if op.ExpectedVersions == nil {
		t.Error("ExpectedVersions should not be nil")
	}

	if op.ExpectedVersions["doc-1"] != 1 {
		t.Errorf("ExpectedVersions[doc-1] = %d, want 1", op.ExpectedVersions["doc-1"])
	}

	if op.ExpectedVersions["doc-2"] != 2 {
		t.Errorf("ExpectedVersions[doc-2] = %d, want 2", op.ExpectedVersions["doc-2"])
	}
}

func TestBatchOperation_WithPayload(t *testing.T) {
	payload := map[string]interface{}{
		"recursive":     true,
		"skip_conflicts": false,
	}

	op := &BatchOperation{
		ProjectID: primitive.NewObjectID(),
		Type:      BatchOpTypeDelete,
		TargetIDs: []string{"doc-1"},
		Payload:   payload,
		Atomic:    true,
	}

	if !op.Atomic {
		t.Error("Operation should be atomic")
	}

	if op.Payload == nil {
		t.Error("Payload should not be nil")
	}

	if op.Payload["recursive"] != true {
		t.Errorf("Payload[recursive] = %v, want true", op.Payload["recursive"])
	}
}

func TestBatchOperation_OriginalTargetIDs(t *testing.T) {
	originalIDs := []string{"doc-1", "doc-2", "doc-3"}
	modifiedIDs := []string{"doc-1", "doc-2"}

	op := &BatchOperation{
		ProjectID:         primitive.NewObjectID(),
		Type:              BatchOpTypeDelete,
		TargetIDs:         modifiedIDs,
		OriginalTargetIDs: originalIDs,
	}

	if len(op.OriginalTargetIDs) != 3 {
		t.Errorf("OriginalTargetIDs length = %d, want 3", len(op.OriginalTargetIDs))
	}

	if len(op.TargetIDs) != 2 {
		t.Errorf("TargetIDs length = %d, want 2", len(op.TargetIDs))
	}
}
