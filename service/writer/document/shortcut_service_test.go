package document

import (
	"context"
	"testing"

	"Qingyu_backend/models/writer"
)

func TestShortcutServiceValidateShortcutsAcceptsDefaultPatterns(t *testing.T) {
	svc := NewShortcutService()

	shortcuts := map[string]writer.Shortcut{
		"save": {
			Action: "save",
			Key:    "Ctrl+S",
		},
		"heading1": {
			Action: "heading1",
			Key:    "Ctrl+Alt+1",
		},
		"indent": {
			Action: "indent",
			Key:    "Tab",
		},
		"fullscreen": {
			Action: "fullscreen",
			Key:    "F11",
		},
		"sidebar": {
			Action: "sidebar",
			Key:    "Ctrl+\\",
		},
	}

	if err := svc.UpdateUserShortcuts(context.Background(), "user-1", shortcuts); err != nil {
		t.Fatalf("expected shortcuts to pass validation, got error: %v", err)
	}
}

func TestShortcutServiceValidateShortcutsRejectsInvalidFormat(t *testing.T) {
	svc := NewShortcutService()

	shortcuts := map[string]writer.Shortcut{
		"save": {
			Action: "save",
			Key:    "Ctrl+",
		},
	}

	if err := svc.UpdateUserShortcuts(context.Background(), "user-1", shortcuts); err == nil {
		t.Fatal("expected invalid shortcut format to fail validation")
	}
}

func TestShortcutServiceValidateShortcutsRejectsDuplicateModifiers(t *testing.T) {
	svc := NewShortcutService()

	shortcuts := map[string]writer.Shortcut{
		"save": {
			Action: "save",
			Key:    "Ctrl+Ctrl+S",
		},
	}

	if err := svc.UpdateUserShortcuts(context.Background(), "user-1", shortcuts); err == nil {
		t.Fatal("expected duplicate modifiers to fail validation")
	}
}

func TestShortcutServiceValidateShortcutsRejectsUnknownModifier(t *testing.T) {
	svc := NewShortcutService()

	shortcuts := map[string]writer.Shortcut{
		"save": {
			Action: "save",
			Key:    "Hyper+S",
		},
	}

	if err := svc.UpdateUserShortcuts(context.Background(), "user-1", shortcuts); err == nil {
		t.Fatal("expected unknown modifier to fail validation")
	}
}
