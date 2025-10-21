package main

import (
	"flag"

	"github.com/shreyaskaundinya/garlic/cmd"
	"github.com/shreyaskaundinya/garlic/models"
	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func main() {
	log := utils.NewLogger()

	// parse arguments for the source and destination paths
	sourcePath := flag.String("src-folder", "", "The source path of the project")
	destinationPath := flag.String("dest-folder", "", "The destination path of the project")
	shouldServe := flag.Bool("serve", false, "Whether to serve the project")
	shouldSeedFiles := flag.Bool("seed-files", false, "Whether to seed the project")
	flag.Parse()

	if sourcePath == nil || destinationPath == nil {
		log.Fatal("Source and destination paths are required")
	}

	config := &models.Config{
		SrcPath:         *sourcePath,
		DestPath:        *destinationPath,
		ShouldServe:     *shouldServe,
		ShouldSeedFiles: *shouldSeedFiles,
	}

	log.Infow("Config: ", "config", config)

	cmd.StartGarlic(config)
}
