package main

import (
	"log"

	"github.com/cloud-barista/poc-cicd-ladybug/src/grpc-api/cbadm/cmd"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// main - Entrypoint
func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Println("cbadm terminated with error: ", err.Error())
	}
}

// ===== [ Public Functions ] =====
