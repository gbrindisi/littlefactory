package tui

import "testing"

func TestStatusIcon(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{
			name:     "closed status",
			status:   "closed",
			expected: "[x]",
		},
		{
			name:     "done status",
			status:   "done",
			expected: "[x]",
		},
		{
			name:     "blocked status",
			status:   "blocked",
			expected: "[!]",
		},
		{
			name:     "in_progress status",
			status:   "in_progress",
			expected: "[>]",
		},
		{
			name:     "active status",
			status:   "active",
			expected: "[>]",
		},
		{
			name:     "pending status",
			status:   "pending",
			expected: "[ ]",
		},
		{
			name:     "empty status",
			status:   "",
			expected: "[ ]",
		},
		{
			name:     "unknown status",
			status:   "unknown",
			expected: "[ ]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := statusIcon(tt.status)
			if result != tt.expected {
				t.Errorf("statusIcon(%q) = %q, want %q", tt.status, result, tt.expected)
			}
		})
	}
}
