package code

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// GetPathSize calculates the total size of files at the given path and
// returns it as a formatted string. If recursive is true, it descends into
// subdirectories. If all is true, hidden files are included in the
// calculation. If human is true, the result is formatted in a
// human-readable form (e.g. "1.2 MB") instead of raw bytes.
func GetPathSize(path string, recursive, human, all bool) (string, error) {
	fileList, err := getFileList(path, recursive)

	if err != nil {
		return "", err
	}

	filteredFileList := filterFileList(fileList, all)

	size := getFileListSize(filteredFileList)

	return formatResult(size, human), nil
}

func getFileList(dir string, recursive bool) ([]os.FileInfo, error) {
	fileList := []os.FileInfo{}

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !isInnerPath(path) {
			return nil
		}

		if recursive && info.IsDir() {
			innerFileList, _ := getFileList(info.Name(), recursive)
			fileList = append(fileList, innerFileList...)
			return nil
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		fileList = append(fileList, info)
		return nil
	})

	if err != nil {
		return []os.FileInfo{}, nil
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

func isInnerPath(path string) bool {
	return strings.Contains(path, "/")
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
		return fmt.Sprintf("%dB", bytes)
	}

	return fmt.Sprintf("%.1f%s", size, sizeList[i])
}
