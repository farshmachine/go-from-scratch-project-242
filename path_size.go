package code

import (
	"errors"
	"fmt"
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
	if path == "" {
		return "", errors.New("Path is not provided")
	}

	fileList, err := getFileList(path, recursive)

	if err != nil {
		return "", err
	}

	filteredFileList := filterFileList(fileList, all)

	size := getFileListSize(filteredFileList)

	return formatResult(size, human), nil
}

func getFileList(dir string, recursive bool) ([]os.FileInfo, error) {
	info, err := os.Stat(dir)

	if err != nil {
		return nil, fmt.Errorf("Can't get info about %s: %w", dir, err)
	}

	if !info.IsDir() {
		return []os.FileInfo{info}, nil
	}

	if recursive {
		return walkRecursive(dir)
	}

	return walkShallow(dir)
}

func walkRecursive(root string) ([]os.FileInfo, error) {
	filePathList := []os.FileInfo{}

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() {
			info, err := entry.Info()

			if err != nil {
				return err
			}

			filePathList = append(filePathList, info)
		}

		return nil
	})

	return filePathList, err
}

func walkShallow(root string) ([]os.FileInfo, error) {
	entryList, err := os.ReadDir(root)

	if err != nil {
		return nil, err
	}

	var files []os.FileInfo

	for _, entry := range entryList {
		if !entry.IsDir() {
			info, err := entry.Info()

			if err != nil {
				return nil, err
			}

			files = append(files, info)
		}
	}

	return files, nil
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
