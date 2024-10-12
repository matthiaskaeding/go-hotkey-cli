package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	data     = make(map[string]string)
	dataFile string
)

func initDataFile() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	dataFile = filepath.Join(homeDir, ".sdb_data.json")
}

func loadData() {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return
	}
	json.Unmarshal(file, &data)
}

func saveData() {
	file, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(dataFile, file, 0644)
}

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"set"},
	Short:   "Set value",
	Long:    "Set key value pair in db",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]

		loadData()
		data[key] = value
		saveData()
		loadData()
		if data[key] != value {
			panic("Setting didn't work")
		} else {
			fmt.Printf("Set key value pair %s, %s", key, value)
		}
	},
}

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"get"},
	Short:   "get value",
	Long:    "Get value from db",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		loadData()
		value, ok := data[key]
		if ok {
			fmt.Printf("Value = %s", value)
		} else {
			fmt.Printf("Key '%s' not found.", key)
		}
	},
}

func init() {
	initDataFile()
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
}
