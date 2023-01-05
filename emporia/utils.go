package emporia

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// ReadStringFromFile reads a value from filename.
func (e *Emporia) ReadStringFromFile(filename string) (string, error) {
	if e.FileExists(filename) {
		valueString, fileErr := os.ReadFile(fmt.Sprintf("%s/%s", e.RootTempDir, filename))
		if fileErr != nil {
			return "", fileErr
		}
		return string(valueString), nil
	}
	return "", errors.New("file does not exist")
}

// FileExists returns true if a filepath exists.
func (e *Emporia) FileExists(name string) bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", e.RootTempDir, name)); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

// WriteValueToFile writes a value to filename.
func (e *Emporia) WriteValueToFile(filename string, value string) error {
	content := []byte(value)
	err := os.WriteFile(fmt.Sprintf("%s/%s", e.RootTempDir, filename), content, 0644)

	if err != nil {
		return err
	}

	return nil
}

// GetModifiedTime reads when the filename was last modified to see if the file contains
// today's max/min. On midnight, this will be yesterday and the file removed for the new
// days recordings.
func (e *Emporia) GetModifiedTime(filename string) (*time.Time, error) {
	file, err := os.Stat(fmt.Sprintf("%s/%s", e.RootTempDir, filename))
	if err != nil {
		return nil, err
	}

	modTime := file.ModTime()
	return &modTime, nil
}

// LocalTime sets a time in location per specified config.
func (e *Emporia) LocalTime() (*time.Time, error) {
	loc, err := time.LoadLocation(e.Timezone)
	if err != nil {
		return nil, err
	}

	localTime := time.Now().In(loc)
	return &localTime, nil
}
