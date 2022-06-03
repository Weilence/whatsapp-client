/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"whatsapp-client/pkg/whatsapp"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run:  newRun(nil),
}

func (ctx *Context) logout() {
	client, err := whatsapp.GetClient(ctx.args[1])
	if err != nil {
		ctx.Write(err)
		return
	}
	err = client.Logout()
	if err != nil {
		ctx.Write(err)
		return
	}
	ctx.Write("登出成功")
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}