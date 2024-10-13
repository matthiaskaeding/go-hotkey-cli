package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSetDataFile(t *testing.T) {
	err := setDataFile()
	assert.NoError(t, err)
	assert.NotEmpty(t, dataFile)
	assert.Contains(t, dataFile, DataFileName)
}

func TestWriteAndSetData(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "sdb-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dataFile = filepath.Join(tempDir, DataFileName)
	data = map[string]string{"test": "value"}

	// Test writeData
	err = writeData()
	assert.NoError(t, err)

	// Verify file contents
	content, err := os.ReadFile(dataFile)
	assert.NoError(t, err)
	var readData map[string]string
	err = json.Unmarshal(content, &readData)
	assert.NoError(t, err)
	assert.Equal(t, data, readData)

	// Test setData
	data = make(map[string]string) // Clear data
	err = setData()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"test": "value"}, data)
}

func TestSetCommand(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "sdb-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dataFile = filepath.Join(tempDir, DataFileName)
	data = make(map[string]string)

	cmd := &cobra.Command{}
	args := []string{"testKey", "testValue"}

	err = setCmd.RunE(cmd, args)
	assert.NoError(t, err)
	assert.Equal(t, "testValue", data["testKey"])
}

func TestGetCommand(t *testing.T) {
	// Setup
	data = map[string]string{"testKey": "testValue"}

	cmd := &cobra.Command{}
	args := []string{"testKey"}

	err := getCmd.RunE(cmd, args)
	assert.NoError(t, err)

	args = []string{"nonExistentKey"}
	err = getCmd.RunE(cmd, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDelCommand(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "sdb-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dataFile = filepath.Join(tempDir, DataFileName)
	data = map[string]string{"testKey": "testValue"}

	cmd := &cobra.Command{}
	args := []string{"testKey"}

	err = delCmd.RunE(cmd, args)
	assert.NoError(t, err)
	assert.NotContains(t, data, "testKey")

	// Test deleting non-existent key
	err = delCmd.RunE(cmd, args)
	assert.NoError(t, err) // It should not return an error for non-existent key
}
