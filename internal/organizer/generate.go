package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"plex-indexer-written-in-go/internal/config"
	"plex-indexer-written-in-go/internal/models"
	"plex-indexer-written-in-go/pkg/fileutils"
)

func RunGenerate(_ *cobra.Command, _ []string) {
	rootFolder, err := fileutils.GetFolderStructure(config.CourseDir)
	if err != nil {
		fmt.Printf("Error getting folder structure: %v\n", err)
		os.Exit(1)
	}

	seriesStructure := generateSeriesStructure(rootFolder)
	seriesStructure = cleanAndReorderSeasons(seriesStructure)

	fileutils.WriteJSONToFile(rootFolder, "original_structure.json")
	if config.JSONFile == "" {
		config.JSONFile = "series_structure.json"
	}
	fileutils.WriteJSONToFile(seriesStructure, config.JSONFile)

	fmt.Println("Original structure has been written to original_structure.json")
	fmt.Printf("Series structure has been written to %s\n", config.JSONFile)
}

func generateSeriesStructure(rootFolder models.Folder) models.SeriesStructure {
	var series models.SeriesStructure
	for i, seasonFolder := range rootFolder.SubFolders {
		season := models.Season{
			Name:     fmt.Sprintf("S%02d • %s", i+1, cleanName(seasonFolder.Name)),
			Episodes: generateEpisodes(seasonFolder, i+1, rootFolder.Name),
		}
		series.Seasons = append(series.Seasons, season)
	}
	return series
}

func generateEpisodes(folder models.Folder, seasonNumber int, rootName string) []models.Episode {
	var episodes []models.Episode
	for i, file := range getAllVideoFiles(folder) {
		parentFolderName := filepath.Base(filepath.Dir(file.FullPath))
		baseName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))

		episodeName := fmt.Sprintf("S%02dE%02d", seasonNumber, i+1)
		if parentFolderName != folder.Name && parentFolderName != rootName {
			episodeName += " • " + cleanName(parentFolderName)
		}
		episodeName += " • " + cleanName(baseName)

		episodes = append(episodes, models.Episode{Name: episodeName, Files: []models.File{file}})
	}
	return episodes
}
