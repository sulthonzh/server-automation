package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

// quotaCmd represents the quota command
var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Manage disk quotas for users",
	Long:  `A command to initialize and set disk quotas for system users.`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		limitGBStr, _ := cmd.Flags().GetString("limitGB")
		limitGB, err := strconv.Atoi(limitGBStr)
		if err != nil {
			log.Fatalf("Invalid limitGB value: %v", err)
		}

		mountPoint, _ := cmd.Flags().GetString("mount")

		// Initialize quotas if needed (could be optional based on your logic)
		if err := InitializeQuotas(mountPoint); err != nil {
			log.Fatalf("Error initializing quotas: %v", err)
		}

		if err := SetUserQuota(username, mountPoint, limitGB); err != nil {
			log.Fatalf("Error setting user quota: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(quotaCmd)
	quotaCmd.Flags().StringP("username", "u", "", "Specify the username for the quota")
	quotaCmd.Flags().StringP("limitGB", "l", "10", "Specify the disk quota limit in GB")
	quotaCmd.Flags().StringP("mount", "m", "/home", "Specify the mount point for the quota")
	quotaCmd.MarkFlagRequired("username")
}

// InitializeQuotas prepares the file system for quotas and enables them.
func InitializeQuotas(mountPoint string) error {
	_, err := exec.Command("quotacheck", "-cum", mountPoint).Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("quotaon", "-v", mountPoint).Output()
	return err
}

// SetUserQuota sets the disk quota for a specific user.
func SetUserQuota(username string, mountPoint string, limitGB int) error {
	limitKB := limitGB * 1024 * 1024 // Convert GB to KB
	cmd := fmt.Sprintf("setquota -u %s %d %d 0 0 %s", username, limitKB, limitKB, mountPoint)
	_, err := exec.Command("bash", "-c", cmd).Output()
	return err
}
