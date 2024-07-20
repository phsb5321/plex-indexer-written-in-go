package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"plex-indexer-written-in-go/internal/config"
	"plex-indexer-written-in-go/internal/models"
	"plex-indexer-written-in-go/pkg/fileutils"
)

var videoExtensions = map[string]bool{".mp4": true, ".mkv": true, ".avi": true, ".mov": true}

func cleanName(name string) string {
	name = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(name, "")
	return strings.Title(strings.ToLower(name))
}

func getAllVideoFiles(folder models.Folder) []models.File {
	var videoFiles []models.File
	for _, file := range folder.Files {
		if videoExtensions[strings.ToLower(filepath.Ext(file.Name))] {
			videoFiles = append(videoFiles, file)
		}
	}
	for _, subFolder := range folder.SubFolders {
		videoFiles = append(videoFiles, getAllVideoFiles(subFolder)...)
	}
	return videoFiles
}

func cleanAndReorderSeasons(series models.SeriesStructure) models.SeriesStructure {
	var cleanedSeasons []models.Season
	for _, season := range series.Seasons {
		if len(season.Episodes) > 0 {
			cleanedSeasons = append(cleanedSeasons, season)
		}
	}

	sort.Slice(cleanedSeasons, func(i, j int) bool {
		numI := extractSeasonNumber(cleanedSeasons[i].Name)
		numJ := extractSeasonNumber(cleanedSeasons[j].Name)
		return numI < numJ
	})

	for i := range cleanedSeasons {
		oldName := cleanedSeasons[i].Name
		newName := fmt.Sprintf("S%02d%s", i+1, oldName[3:])
		cleanedSeasons[i].Name = newName

		for j := range cleanedSeasons[i].Episodes {
			oldEpisodeName := cleanedSeasons[i].Episodes[j].Name
			newEpisodeName := fmt.Sprintf("S%02d%s", i+1, oldEpisodeName[3:])
			cleanedSeasons[i].Episodes[j].Name = newEpisodeName
		}
	}

	return models.SeriesStructure{Seasons: cleanedSeasons}
}

func extractSeasonNumber(seasonName string) int {
	re := regexp.MustCompile(`S(\d+)`)
	match := re.FindStringSubmatch(seasonName)
	if len(match) > 1 {
		num := 0
		fmt.Sscanf(match[1], "%d", &num)
		return num
	}
	return 0
}

func RunFull(_ *cobra.Command, _ []string) {
	rootFolder, err := fileutils.GetFolderStructure(config.CourseDir)
	if err != nil {
		fmt.Printf("Error getting folder structure: %v\n", err)
		os.Exit(1)
	}

	seriesStructure := generateSeriesStructure(rootFolder)
	seriesStructure = cleanAndReorderSeasons(seriesStructure)

	if config.GenerateJSON {
		fileutils.WriteJSONToFile(rootFolder, "original_structure.json")
		fileutils.WriteJSONToFile(seriesStructure, "series_structure.json")
		fmt.Println("Original structure has been written to original_structure.json")
		fmt.Println("Series structure has been written to series_structure.json")
	}

	createSymlinks(seriesStructure, config.OutputDir, rootFolder.Name)
}
