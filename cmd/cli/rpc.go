package main

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Take the server out of maintenance mode",
	Run: func(cmd *cobra.Command, args []string) {
		rpcClient(false)
	},
}

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Put the server into maintenance mode",
	Run: func(cmd *cobra.Command, args []string) {
		rpcClient(true)
	},
}

func rpcClient(inMaintenanceMode bool) {
	rpcPort := os.Getenv("RPC_PORT")
	c, err := rpc.Dial("tcp", "127.0.0.1:"+rpcPort)
	if err != nil {
		exitGracefully(err)
	}

	fmt.Println("Connected...")
	var result string
	err = c.Call("RPCServer.MaintenanceMode", inMaintenanceMode, &result)
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow(result)
}
