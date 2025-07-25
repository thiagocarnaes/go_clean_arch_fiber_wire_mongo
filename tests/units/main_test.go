package units

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	buildOutputMsg   = "Build output: %s"
	mainGoFile       = "main.go"
	mainFunctionText = "func main()"
)

func TestMainBuild(t *testing.T) {
	// Test that main function exists and the project can be built

	// Get project root directory
	projectRoot := filepath.Join("..", "..")

	// Create a temporary directory for building
	tempDir := t.TempDir()

	// Build the application
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(tempDir, "app"), ".")
	buildCmd.Dir = projectRoot

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf(buildOutputMsg, output)
		require.NoError(t, err, "Failed to build the application")
	}

	// Test that the binary was created
	binaryPath := filepath.Join(tempDir, "app")
	_, err = os.Stat(binaryPath)
	assert.NoError(t, err, "Binary should be created successfully")
}

func TestMainPackageStructure(t *testing.T) {
	// Test that main.go has the expected structure

	// Read main.go file from the project root
	projectRoot := filepath.Join("..", "..")
	mainFile := filepath.Join(projectRoot, mainGoFile)

	mainContent, err := os.ReadFile(mainFile)
	require.NoError(t, err)

	mainString := string(mainContent)

	// Check that it contains expected elements
	assert.Contains(t, mainString, "package main", "Should have main package declaration")
	assert.Contains(t, mainString, mainFunctionText, "Should have main function")
	assert.Contains(t, mainString, "cmd.Execute()", "Should call cmd.Execute()")
	assert.Contains(t, mainString, "user-management/cmd", "Should import cmd package")
}

func TestMainSimplicity(t *testing.T) {
	// Test that the main package properly uses dependency injection

	// Read main.go file from the project root
	projectRoot := filepath.Join("..", "..")
	mainFile := filepath.Join(projectRoot, mainGoFile)

	mainContent, err := os.ReadFile(mainFile)
	require.NoError(t, err)

	mainString := string(mainContent)

	// Verify that main.go is simple and delegates to cmd package
	assert.NotContains(t, mainString, "database", "Main should not directly handle database")
	assert.NotContains(t, mainString, "config", "Main should not directly handle config")
	assert.NotContains(t, mainString, "repository", "Main should not directly handle repositories")
	assert.NotContains(t, mainString, "mongodb", "Main should not directly handle mongodb")
	assert.NotContains(t, mainString, "fiber", "Main should not directly handle fiber")

	// Should only contain the basic structure
	lines := len(strings.Split(strings.TrimSpace(mainString), "\n"))
	assert.Less(t, lines, 15, "Main.go should be simple and short")

	// Should be mostly empty except for package, import, and main function
	assert.Contains(t, mainString, "package main")
	assert.Contains(t, mainString, "import")
	assert.Contains(t, mainString, mainFunctionText)
}

func TestMainProjectStructure(t *testing.T) {
	// Test that the project has the expected structure for main.go to work

	projectRoot := filepath.Join("..", "..")

	// Check that essential files exist
	requiredFiles := []string{
		mainGoFile,
		"go.mod",
		"cmd/root.go",
	}

	for _, file := range requiredFiles {
		filePath := filepath.Join(projectRoot, file)
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "Required file should exist: %s", file)
	}
}

func TestMainGoMod(t *testing.T) {
	// Test that go.mod is consistent and the project compiles

	projectRoot := filepath.Join("..", "..")

	// Test that the project can be compiled
	buildCmd := exec.Command("go", "build", ".")
	buildCmd.Dir = projectRoot

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf(buildOutputMsg, output)
		require.NoError(t, err, "Project should compile successfully")
	}

	t.Log("Project builds successfully")
}

func TestMainHasMinimalDependencies(t *testing.T) {
	// Test that main.go doesn't import unnecessary packages

	projectRoot := filepath.Join("..", "..")
	mainFile := filepath.Join(projectRoot, mainGoFile)

	mainContent, err := os.ReadFile(mainFile)
	require.NoError(t, err)

	mainString := string(mainContent)

	// Count imports - should be minimal
	importLines := strings.Count(mainString, "import")
	assert.LessOrEqual(t, importLines, 1, "Main should have minimal imports")

	// Should not import complex packages directly
	prohibitedImports := []string{
		"github.com/gofiber/fiber",
		"go.mongodb.org/mongo-driver",
		"github.com/google/wire",
		"database/sql",
	}

	for _, prohibitedImport := range prohibitedImports {
		assert.NotContains(t, mainString, prohibitedImport,
			"Main should not directly import %s", prohibitedImport)
	}
}

// Benchmark test for building the application
func BenchmarkMainBuild(b *testing.B) {
	projectRoot := filepath.Join("..", "..")

	// Create a temporary directory for building
	tempDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Build the application
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tempDir, "app"), ".")
		buildCmd.Dir = projectRoot

		output, err := buildCmd.CombinedOutput()
		if err != nil {
			b.Logf(buildOutputMsg, output)
			b.Fatal(err)
		}

		// Clean up the binary for next iteration
		os.Remove(filepath.Join(tempDir, "app"))
	}
}

func TestMainFileSize(t *testing.T) {
	// Test that main.go is appropriately small

	projectRoot := filepath.Join("..", "..")
	mainFile := filepath.Join(projectRoot, mainGoFile)

	fileInfo, err := os.Stat(mainFile)
	require.NoError(t, err)

	// Main.go should be small (less than 1KB for a simple main function)
	assert.Less(t, fileInfo.Size(), int64(1024), "Main.go should be small and simple")
}

func TestMainFunctionExists(t *testing.T) {
	// Test that main function signature is correct

	projectRoot := filepath.Join("..", "..")
	mainFile := filepath.Join(projectRoot, mainGoFile)

	mainContent, err := os.ReadFile(mainFile)
	require.NoError(t, err)

	mainString := string(mainContent)

	// Check for correct main function signature
	assert.Contains(t, mainString, mainFunctionText, "Should have main function with no parameters")
	assert.NotContains(t, mainString, "func main(args", "Main function should not have parameters")
	assert.NotContains(t, mainString, "func main(os.Args", "Main function should not take os.Args")
}
