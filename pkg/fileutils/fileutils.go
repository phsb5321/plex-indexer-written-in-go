package fileutils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"plex-indexer-written-in-go/internal/models"

	"facette.io/natsort"
)

func GetFolderStructure(dir string) (models.Folder, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return models.Folder{}, fmt.Errorf("getting absolute path: %w", err)
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return models.Folder{}, fmt.Errorf("reading directory: %w", err)
	}

	folder := models.Folder{Name: filepath.Base(absDir), FullPath: absDir}

	for _, entry := range entries {
		fullPath := filepath.Join(absDir, entry.Name())
		if entry.IsDir() {
			subFolder, err := GetFolderStructure(fullPath)
			if err != nil {
				return models.Folder{}, err
			}
			folder.SubFolders = append(folder.SubFolders, subFolder)
		} else {
			folder.Files = append(folder.Files, models.File{Name: entry.Name(), FullPath: fullPath})
		}
	}

	sortFolderContents(&folder)
	return folder, nil
}

func sortFolderContents(folder *models.Folder) {
	sort.Slice(folder.Files, func(i, j int) bool {
		return natsort.Compare(folder.Files[i].Name, folder.Files[j].Name)
	})
	sort.Slice(folder.SubFolders, func(i, j int) bool {
		return natsort.Compare(folder.SubFolders[i].Name, folder.SubFolders[j].Name)
	})
}

func WriteJSONToFile(data interface{}, filename string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func ReadJSONFromFile(filename string) (models.SeriesStructure, error) {
	var seriesStructure models.SeriesStructure

	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return seriesStructure, fmt.Errorf("reading JSON file: %w", err)
	}

	err = json.Unmarshal(jsonData, &seriesStructure)
	if err != nil {
		return seriesStructure, fmt.Errorf("unmarshalling JSON: %w", err)
	}

	return seriesStructure, nil
}
