package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// phpfpmCmd represents the phpfpm command
var phpfpmCmd = &cobra.Command{
	Use:   "phpfpm",
	Short: "Limit memory usage for PHP-FPM pools",
	Long:  `A command to modify PHP-FPM pool configurations to limit memory usage.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		version, _ := cmd.Flags().GetString("version")
		poolDir := fmt.Sprintf("/etc/php/%s/fpm/pool.d/", version)
		if dir != "" {
			poolDir = dir
		}

		err := ModifyPHPFPMPoolConfig(poolDir)
		if err != nil {
			log.Fatalf("Error modifying PHP-FPM pool configurations: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(phpfpmCmd)
	phpfpmCmd.Flags().StringP("dir", "d", "", "Specify the PHP-FPM pool configuration directory")
	phpfpmCmd.Flags().StringP("version", "v", "7.4", "Specify the PHP version")
}

// ModifyPHPFPMPoolConfig updates the PHP-FPM pool configuration to limit the number of child processes.
func ModifyPHPFPMPoolConfig(poolDir string) error {
	files, err := ioutil.ReadDir(poolDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".conf") {
			filePath := filepath.Join(poolDir, f.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			// Example modification: Adjusting max children (this is a simplification)
			newContent := strings.Replace(string(content), "pm.max_children = 5", "pm.max_children = 10", -1)

			err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
