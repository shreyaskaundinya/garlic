package cmd

import (
	"fmt"

	"github.com/shreyaskaundinya/garlic/pkg/server"
)

func StartGarlic(SrcPath string, DestPath string) {
	s, err := server.NewServer(SrcPath, DestPath)

	if err != nil || s == nil {
		// panic(err)
		fmt.Println("Error: ", err)
		return
	}

	s.Start()
}
