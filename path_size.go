package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	fileList, err := getFileList(path, recursive)

	if err != nil {
		return "", err
	}

	filteredFileList := filterFileList(fileList, all)

	size := getFileListSize(filteredFileList)

	return formatResult(size, human), nil
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

func formatResult(bytes int64, human bool) string {
	sizeList := [7]string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	size := float64(bytes)
	i := 0

	for size >= 1024 {
		size /= 1024
		i++
	}

	if !human || i == 0 {
		return fmt.Sprintf("%dB\n", bytes)
	}

	return fmt.Sprintf("%.1f%s\n", size, sizeList[i])
}
