package utils

import "testing"

func TestDetectType(t *testing.T) {

	tests := []struct {
		name       string
		folderName string
		expected   string
	}{
		{
			name:       "Detect Node.js folder",
			folderName: "node_modules", expected: "Node.js",
		},
		{
			name:       "Detect Node.js cache folder",
			folderName: "node_modules_cache",
			expected:   "Node.js",
		},
		{
			name:       "Detect Go/PHP vendor folder",
			folderName: "vendor",
			expected:   "Go/PHP"},
		{
			name:       "Detect Python .venv folder",
			folderName: ".venv",
			expected:   "Python"},
		{
			name:       "Detect Python __pycache__ folder",
			folderName: "__pycache__",
			expected:   "Python",
		},
		{
			name:       "Detect Python venv folder",
			folderName: "venv",
			expected:   "Python",
		},
		{
			name:       "Detect Rust target folder",
			folderName: "target",
			expected:   "Rust",
		},
		{
			name:       "Detect unknown folder",
			folderName: "unknown_folder",
			expected:   "Unknown",
		},
		{
			name:       "Detect empty folder name",
			folderName: "",
			expected:   "Unknown",
		},
		{
			name:       "Detect uppercase NODE_MODULES folder",
			folderName: "NODE_MODULES",
			expected:   "Unknown",
		},
		{
			name:       "Detect random folder name",
			folderName: "random_name",
			expected:   "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectType(tt.folderName)
			if result != tt.expected {
				t.Errorf("DetectType(%q) = %q; want %q", tt.folderName, result, tt.expected)
			}
		})
	}
}

func TestIsTargetDirectory(t *testing.T) {

	tests := []struct {
		name       string
		folderName string
		expected   bool
	}{
		{
			name:       "Is Node.js folder",
			folderName: "node_modules", expected: true,
		},
		{
			name:       "Is Node.js cache folder",
			folderName: "node_modules_cache",
			expected:   true,
		},
		{
			name:       "Is Go/PHP vendor folder",
			folderName: "vendor",
			expected:   true,
		},
		{
			name:       "Is Python .venv folder",
			folderName: ".venv",
			expected:   true,
		},
		{
			name:       "Is Python __pycache__ folder",
			folderName: "__pycache__",
			expected:   true,
		},
		{
			name:       "Is Python venv folder",
			folderName: "venv",
			expected:   true,
		},
		{
			name:       "Is Rust target folder",
			folderName: "target",
			expected:   true,
		},
		{
			name:       "Is unknown folder",
			folderName: "unknown_folder",
			expected:   false,
		},
		{
			name:       "Is empty folder name",
			folderName: "",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTargetDirectory(tt.folderName)
			if result != tt.expected {
				t.Errorf("IsTargetDirectory(%q) = %v; want %v", tt.folderName, result, tt.expected)
			}
		})
	}

}

// Benchmark tests

func BenchmarkDetectType(b *testing.B) {
	folderNames := []string{
		"node_modules",
		"node_modules_cache",
		"vendor",
		".venv",
		"__pycache__",
		"venv",
		"target",
		"unknown_folder",
		"",
		"NODE_MODULES",
		"random_name",
	}

	for _, folderName := range folderNames {
		b.Run(folderName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				DetectType(folderName)
			}
		})

	}

}

func BenchmarkIsTargetDirectory(b *testing.B) {
	folderNames := []string{
		"node_modules",
		"node_modules_cache",
		"vendor",
		".venv",
		"__pycache__",
		"venv",
		"target",
		"unknown_folder",
		"",
	}

	for _, folderName := range folderNames {
		b.Run(folderName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IsTargetDirectory(folderName)
			}
		})

	}

}
