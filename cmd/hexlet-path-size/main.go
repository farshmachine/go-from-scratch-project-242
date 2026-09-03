package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory",
		ArgsUsage: "<path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursive size of directories",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 1 {
				return fmt.Errorf("missing required argument: <path>")
			}

			path := c.Args().First()
			human := c.Bool("human")
			includeHidden := c.Bool("all")
			recursive := c.Bool("recursive")

			size, err := GetPathSize(path, includeHidden, recursive)
			result := formatResult(size, path, human)

			if err != nil {
				return err
			}

			fmt.Print(result)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func GetPathSize(path string, includeHidden bool, recursive bool) (int64, error) {
	fileList, err := getFileList(path, recursive)

	if err != nil {
		return 0, err
	}

	filteredFileList := filterFileList(fileList, includeHidden)

	return getFileListSize(filteredFileList), nil
}

func getFileList(path string, recursive bool) ([]os.FileInfo, error) {
	fileList := []os.FileInfo{}

	fileInfo, err := os.Lstat(path)

	if err != nil {
		return []os.FileInfo{}, nil
	}

	if !fileInfo.IsDir() {
		return append(fileList, fileInfo), nil
	}

	dirEntryList, err := os.ReadDir(path)

	for _, entry := range dirEntryList {
		if !entry.IsDir() {
			fileInfo, err := entry.Info()

			if err != nil {
				return []os.FileInfo{}, err
			}

			fileList = append(fileList, fileInfo)
		}

		if recursive && entry.IsDir() {
			dirPath := filepath.Join(path, entry.Name())
			innerFileList, _ := getFileList(dirPath, recursive)
			fileList = append(fileList, innerFileList...)
		}
	}

	return fileList, nil

}

func filterFileList(fileList []os.FileInfo, includeHidden bool) []os.FileInfo {
	filtered := make([]os.FileInfo, 0, len(fileList))

	for _, file := range fileList {
		isHidden := isFileHidden(file.Name())

		if !isHidden || includeHidden {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

func getFileListSize(fileList []os.FileInfo) int64 {
	var size int64

	for _, file := range fileList {
		size += file.Size()
	}

	return size
}

func getFileSize(file os.FileInfo) int64 {
	if file.IsDir() {
		return 0
	}

	return file.Size()
}

func isFileHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

func formatResult(bytes int64, path string, human bool) string {
	if !human {
		return fmt.Sprintf("%dB \t %s\n", bytes, path)
	}

	sizeList := [7]string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	size := float64(bytes)
	i := 0

	for size >= 1024 {
		size /= 1024
		i++
	}

	if size == float64(int64(size)) {
		return fmt.Sprintf("%.0f%s \t %s\n", size, sizeList[i], path)
	}

	fmt.Println(path)

	return fmt.Sprintf("%.2f%s \t %s\n", size, sizeList[i], path)

}
