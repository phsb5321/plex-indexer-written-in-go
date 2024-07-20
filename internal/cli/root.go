package cli

import (
	"github.com/spf13/cobra"

	"plex-indexer-written-in-go/internal/config"
	"plex-indexer-written-in-go/internal/organizer"
)

var rootCmd = &cobra.Command{Use: "course-organizer"}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.CourseDir, "course", "c", "", "Course directory path")
	rootCmd.PersistentFlags().StringVarP(&config.OutputDir, "output", "o", "", "Output directory path")

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate JSON files from course directory",
		Run:   organizer.RunGenerate,
	}
	generateCmd.Flags().StringVarP(&config.JSONFile, "json", "j", "", "JSON file to write series structure")
	generateCmd.MarkFlagRequired("course")

	symlinkCmd := &cobra.Command{
		Use:   "symlink",
		Short: "Create symlinks based on JSON file",
		Run:   organizer.RunSymlink,
	}
	symlinkCmd.Flags().StringVarP(&config.JSONFile, "json", "j", "", "JSON file containing series structure")
	symlinkCmd.MarkFlagRequired("json")
	symlinkCmd.MarkFlagRequired("output")

	fullCmd := &cobra.Command{
		Use:   "full",
		Short: "Generate JSON files and create symlinks",
		Run:   organizer.RunFull,
	}
	fullCmd.Flags().BoolVarP(&config.GenerateJSON, "generate-json", "g", false, "Generate JSON files")
	fullCmd.MarkFlagRequired("course")
	fullCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(generateCmd, symlinkCmd, fullCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
