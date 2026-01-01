package cmd

import (
	"fmt"

	"github.com/ratneshrt/github-activity/activity"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "github-activity",
		Short: "Github User Ctivity is a CLI tool for fetching user activity",
		Long: `Github User Ctivity is a CLI tool for fetching user activity. It allows you to fetch user activity by providing the username.

		Example: 
		> github-activity ratneshrt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunDisplayActivityCmd(args)
		},
	}

	return cmd
}

func RunDisplayActivityCmd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please provide username")
	}

	username := args[0]
	activites, err := activity.FetchGithubActivity(username)

	if err != nil {
		return err
	}

	return activity.DisplayActivity(username, activites)
}
