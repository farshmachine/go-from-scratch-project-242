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

	return formatResult(size, path, human), nil
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
