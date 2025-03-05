package cmd

import (
	"fmt"
	"net/rpc"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Take the server out of maintenance mode",
	Run: func(cmd *cobra.Command, args []string) {
		// Start spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Taking server out of maintenance mode..."
		s.Color("green")
		s.Start()

		err := rpcClient(false)

		s.Stop()
		if err != nil {
			exitGracefully(err)
		}

		color.Green("✓ Server is now live!")
	},
}

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Put the server into maintenance mode",
	Run: func(cmd *cobra.Command, args []string) {
		// Start spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Putting server into maintenance mode..."
		s.Color("yellow")
		s.Start()

		err := rpcClient(true)

		s.Stop()
		if err != nil {
			exitGracefully(err)
		}

		color.Yellow("⚠ Server is now in maintenance mode")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}

func rpcClient(inMaintenanceMode bool) error {
	setup(nil, nil)

	rpcPort := os.Getenv("RPC_PORT")
	if rpcPort == "" {
		return fmt.Errorf("RPC_PORT is not set in your .env file")
	}

	c, err := rpc.Dial("tcp", "127.0.0.1:"+rpcPort)
	if err != nil {
		return err
	}

	var result string
	err = c.Call("RPCServer.MaintenanceMode", inMaintenanceMode, &result)
	if err != nil {
		return err
	}

	return nil
}
