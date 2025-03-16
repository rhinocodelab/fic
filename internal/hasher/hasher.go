/* github.com/rhinocodelab/fic/internal/hasher/hasher.go */

package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// CalculateHash generates a SHA256 hash for a given file.
func CalculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file '%s': %v", filePath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info for '%s': %v", filePath, err)
	}

	// Handle empty file scenario
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("file '%s' is empty", filePath)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to read file '%s' during hashing: %v", filePath, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
