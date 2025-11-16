package testutils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Debugging migrations path detection...")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}
	fmt.Printf("Current working directory: %s\n", wd)

	possiblePaths := []string{
		filepath.Join(wd, "migrations"),
		filepath.Join(wd, "..", "migrations"),
		filepath.Join(wd, "..", "..", "migrations"),
		filepath.Join(wd, "backend", "migrations"),
	}

	fmt.Println("Checking possible paths:")
	for i, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("✅ [%d] FOUND: %s\n", i+1, path)

			// List files in migrations directory
			files, _ := filepath.Glob(filepath.Join(path, "*.sql"))
			fmt.Printf("   Files: %v\n", files)
		} else {
			fmt.Printf("❌ [%d] NOT FOUND: %s\n", i+1, path)
		}
	}
}