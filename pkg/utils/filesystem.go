package utils

var targetDirectories = []string{
	"node_modules",       // common Node.js dependencies folder
	"node_modules_cache", // alternative Node.js cache folder
	"vendor",             // Go/Php vendor folders
	".venv",              // Python virtual environment
	"__pycache__",        // Python cache folder
	"venv",               // Python virtual environment
	"target",             // Rust target folder
}

var folderTypeMap = map[string]string{
	"node_modules":       "Node.js",
	"node_modules_cache": "Node.js",
	"vendor":             "Go/PHP",
	".venv":              "Python",
	"__pycache__":        "Python",
	"venv":               "Python",
	"target":             "Rust",
}

func DetectType(folderName string) string {

	if typValue, ok := folderTypeMap[folderName]; ok {
		return typValue
	}
	return "Unknown"

}

func IsTargetDirectory(name string) bool {

	for _, target := range targetDirectories {
		if name == target {
			return true
		}
	}

	return false

}
