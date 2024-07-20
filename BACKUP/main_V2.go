package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"facette.io/natsort"
	"github.com/spf13/cobra"
)

type File struct {
	Name     string `json:"name"`
	FullPath string `json:"fullPath"`
}

type Folder struct {
	Name       string   `json:"name"`
	FullPath   string   `json:"fullPath"`
	Files      []File   `json:"files,omitempty"`
	SubFolders []Folder `json:"subFolders,omitempty"`
}

type Episode struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
}

type Season struct {
	Name     string    `json:"name"`
	Episodes []Episode `json:"episodes"`
}

type SeriesStructure struct {
	Seasons []Season `json:"seasons"`
}

var (
	videoExtensions = map[string]bool{".mp4": true, ".mkv": true, ".avi": true, ".mov": true}
	rootCmd         = &cobra.Command{Use: "plex-indexer-written-in-go"}
	courseDir       string
	outputDir       string
	jsonFile        string
	generateJSON    bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&courseDir, "course", "c", "", "Course directory path")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "Output directory path")

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate JSON files from course directory",
		Run:   runGenerate,
	}
	generateCmd.Flags().StringVarP(&jsonFile, "json", "j", "", "JSON file to write series structure")
	generateCmd.MarkFlagRequired("course")

	symlinkCmd := &cobra.Command{
		Use:   "symlink",
		Short: "Create symlinks based on JSON file",
		Run:   runSymlink,
	}
	symlinkCmd.Flags().StringVarP(&jsonFile, "json", "j", "", "JSON file containing series structure")
	symlinkCmd.MarkFlagRequired("json")
	symlinkCmd.MarkFlagRequired("output")

	fullCmd := &cobra.Command{
		Use:   "full",
		Short: "Generate JSON files and create symlinks",
		Run:   runFull,
	}
	fullCmd.Flags().BoolVarP(&generateJSON, "generate-json", "g", false, "Generate JSON files")
	fullCmd.MarkFlagRequired("course")
	fullCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(generateCmd, symlinkCmd, fullCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runGenerate(_ *cobra.Command, _ []string) {
	rootFolder, err := getFolderStructure(courseDir)
	if err != nil {
		fmt.Printf("Error getting folder structure: %v\n", err)
		os.Exit(1)
	}

	seriesStructure := generateSeriesStructure(rootFolder)
	seriesStructure = cleanAndReorderSeasons(seriesStructure)

	writeJSONToFile(rootFolder, "original_structure.json")
	if jsonFile == "" {
		jsonFile = "series_structure.json"
	}
	writeJSONToFile(seriesStructure, jsonFile)

	fmt.Println("Original structure has been written to original_structure.json")
	fmt.Printf("Series structure has been written to %s\n", jsonFile)
}

func runSymlink(_ *cobra.Command, _ []string) {
	seriesStructure, err := readJSONFromFile(jsonFile)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		os.Exit(1)
	}

	createSymlinks(seriesStructure, outputDir, filepath.Base(outputDir))
}

func runFull(_ *cobra.Command, _ []string) {
	rootFolder, err := getFolderStructure(courseDir)
	if err != nil {
		fmt.Printf("Error getting folder structure: %v\n", err)
		os.Exit(1)
	}

	seriesStructure := generateSeriesStructure(rootFolder)
	seriesStructure = cleanAndReorderSeasons(seriesStructure)

	if generateJSON {
		writeJSONToFile(rootFolder, "original_structure.json")
		writeJSONToFile(seriesStructure, "series_structure.json")
		fmt.Println("Original structure has been written to original_structure.json")
		fmt.Println("Series structure has been written to series_structure.json")
	}

	createSymlinks(seriesStructure, outputDir, rootFolder.Name)
}

func getFolderStructure(dir string) (Folder, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return Folder{}, fmt.Errorf("getting absolute path: %w", err)
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return Folder{}, fmt.Errorf("reading directory: %w", err)
	}

	folder := Folder{Name: filepath.Base(absDir), FullPath: absDir}

	for _, entry := range entries {
		fullPath := filepath.Join(absDir, entry.Name())
		if entry.IsDir() {
			subFolder, err := getFolderStructure(fullPath)
			if err != nil {
				return Folder{}, err
			}
			folder.SubFolders = append(folder.SubFolders, subFolder)
		} else {
			folder.Files = append(folder.Files, File{Name: entry.Name(), FullPath: fullPath})
		}
	}

	sortFolderContents(&folder)
	return folder, nil
}

func sortFolderContents(folder *Folder) {
	sort.Slice(folder.Files, func(i, j int) bool {
		return natsort.Compare(folder.Files[i].Name, folder.Files[j].Name)
	})
	sort.Slice(folder.SubFolders, func(i, j int) bool {
		return natsort.Compare(folder.SubFolders[i].Name, folder.SubFolders[j].Name)
	})
}

func generateSeriesStructure(rootFolder Folder) SeriesStructure {
	var series SeriesStructure
	for i, seasonFolder := range rootFolder.SubFolders {
		season := Season{
			Name:     fmt.Sprintf("S%02d • %s", i+1, cleanName(seasonFolder.Name)),
			Episodes: generateEpisodes(seasonFolder, i+1, rootFolder.Name),
		}
		series.Seasons = append(series.Seasons, season)
	}
	return series
}

func generateEpisodes(folder Folder, seasonNumber int, rootName string) []Episode {
	var episodes []Episode
	for i, file := range getAllVideoFiles(folder) {
		parentFolderName := filepath.Base(filepath.Dir(file.FullPath))
		baseName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))

		episodeName := fmt.Sprintf("S%02dE%02d", seasonNumber, i+1)
		if parentFolderName != folder.Name && parentFolderName != rootName {
			episodeName += " • " + cleanName(parentFolderName)
		}
		episodeName += " • " + cleanName(baseName)

		episodes = append(episodes, Episode{Name: episodeName, Files: []File{file}})
	}
	return episodes
}

func getAllVideoFiles(folder Folder) []File {
	var videoFiles []File
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

func cleanName(name string) string {
	name = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(name, "")
	return strings.Title(strings.ToLower(name))
}

func writeJSONToFile(data interface{}, filename string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func readJSONFromFile(filename string) (SeriesStructure, error) {
	var seriesStructure SeriesStructure

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

func createSymlinks(series SeriesStructure, outputDir, rootName string) {
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

func cleanAndReorderSeasons(series SeriesStructure) SeriesStructure {
	var cleanedSeasons []Season
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

	return SeriesStructure{Seasons: cleanedSeasons}
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
