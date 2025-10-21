package cmd

import (
	"fmt"

	"github.com/shreyaskaundinya/garlic/models"
	"github.com/shreyaskaundinya/garlic/pkg/server"
)

func StartGarlic(config *models.Config) {
	s, err := server.NewServer(config)

	if err != nil || s == nil {
		// panic(err)
		fmt.Println("Error: ", err)
		return
	}

	s.Start(config)
}
