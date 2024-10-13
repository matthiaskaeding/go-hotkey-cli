package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.design/x/clipboard"

	"github.com/spf13/cobra"
)

const (
	DataFileName = ".sdb_data.json"
)

var (
	dataFile string
	data     map[string]string
)

func setDataFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dataFile = filepath.Join(homeDir, DataFileName)
	return nil
}
func setData() error {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &data)
}
func initializeApp() error {
	errDataFile := setDataFile()
	if errDataFile != nil {
		return errDataFile
	}
	data = make(map[string]string)
	if err := setData(); err != nil {
		return fmt.Errorf("error loading data: %w", err)
	}
	return nil
}

func writeData() error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}
	write_err := os.WriteFile(dataFile, file, 0644)
	if write_err != nil {
		return fmt.Errorf("error writing data file: %w", write_err)
	}
	return nil
}

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"set"},
	Short:   "Set value",
	Long:    "Set key value pair in db",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		k, v := args[0], args[1]
		data[k] = v
		return writeData()
	},
}

var delCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"del"},
	Short:   "Delete value",
	Long:    "Delete key value pair in db",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		k := args[0]
		delete(data, k)
		_, ok := data[k]
		if !ok {
			return writeData()
		} else {
			return fmt.Errorf("deletion didn't work")
		}
	},
}

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"get"},
	Short:   "get value",
	Long:    "Get value from db",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value, ok := data[key]
		if !ok {
			return fmt.Errorf("key '%s' not found", key)
		}
		fmt.Printf("Value = %s", value)
		return nil
	},
}

var cpCmd = &cobra.Command{
	Use:     "cp",
	Aliases: []string{"cp"},
	Short:   "Copy value to clipboard",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		k := args[0]
		v, ok := data[k]
		if ok {
			clipboard.Write(clipboard.FmtText, []byte(v))
			fmt.Printf("Value '%s' copied to clipboard", v)
			return nil
		} else {
			return fmt.Errorf("key '%s' not found", k)
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		if err := initializeApp(); err != nil {
			fmt.Println("Error initializing app:", err)
			os.Exit(1)
		}
	})
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(delCmd)
	rootCmd.AddCommand(cpCmd)
}
