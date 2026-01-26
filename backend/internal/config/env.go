package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv loads environment variables from a .env file.
func LoadEnv() error {
	file, err := os.Open("../.env")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip the comments, empty strings
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// remove quotes, if present
		value = strings.Trim(value, `"'`)

		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
