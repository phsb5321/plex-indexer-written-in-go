package organizer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"plex-indexer-written-in-go/internal/config"
	"plex-indexer-written-in-go/internal/models"
	"plex-indexer-written-in-go/pkg/fileutils"
)

func RunSymlink(_ *cobra.Command, _ []string) {
	seriesStructure, err := fileutils.ReadJSONFromFile(config.JSONFile)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		os.Exit(1)
	}

	createSymlinks(seriesStructure, config.OutputDir, filepath.Base(config.OutputDir))
}

func createSymlinks(series models.SeriesStructure, outputDir, rootName string) {
	seriesDir := filepath.Join(outputDir, rootName)
	if err := os.MkdirAll(seriesDir, 0755); err != nil {
		fmt.Printf("Error creating series directory: %v\n", err)
		return
	}

	for _, season := range series.Seasons {
		seasonDir := filepath.Join(seriesDir, season.Name)
		if err := os.MkdirAll(seasonDir, 0755); err != nil {
			fmt.Printf("Error creating season directory: %v\n", err)
			continue
		}

		for _, episode := range season.Episodes {
			for _, file := range episode.Files {
				symlinkPath := filepath.Join(seasonDir, episode.Name+filepath.Ext(file.Name))
				if err := os.Symlink(file.FullPath, symlinkPath); err != nil {
					fmt.Printf("Error creating symlink: %v\n", err)
				}
			}
		}
	}

	fmt.Println("Symlinks have been created in the output directory")
}
