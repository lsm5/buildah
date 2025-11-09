package digestutils

import "fmt"

func AlgorithmPrefix(id string) (string, error) {
	idLength := len(id)
	switch idLength {
	case 64:
		return "sha256:", nil
	case 128:
		return "sha512:", nil
	default:
		return "", fmt.Errorf("invalid hash length: expected 64 or 128, got %d", idLength)
	}
}
