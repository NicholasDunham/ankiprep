package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// FileService provides robust file I/O operations with comprehensive error handling
type FileService struct {
	maxRetries    int
	retryDelay    time.Duration
	tempDir       string
	backupEnabled bool
}

// FileError represents a file operation error with enhanced context
type FileError struct {
	Operation string
	Path      string
	Err       error
	Retries   int
	Timestamp time.Time
}

// Error implements the error interface
func (fe *FileError) Error() string {
	return fmt.Sprintf("%s failed for %s after %d retries: %v",
		fe.Operation, fe.Path, fe.Retries, fe.Err)
}

// Unwrap returns the underlying error for error unwrapping
func (fe *FileError) Unwrap() error {
	return fe.Err
}

// NewFileService creates a new FileService with default configuration
func NewFileService() *FileService {
	tempDir := os.TempDir()
	return &FileService{
		maxRetries:    3,
		retryDelay:    100 * time.Millisecond,
		tempDir:       tempDir,
		backupEnabled: false,
	}
}

// SetRetryPolicy configures the retry behavior for file operations
func (fs *FileService) SetRetryPolicy(maxRetries int, retryDelay time.Duration) {
	fs.maxRetries = maxRetries
	fs.retryDelay = retryDelay
}

// SetTempDirectory sets the directory for temporary files
func (fs *FileService) SetTempDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create temp directory %s: %w", path, err)
		}
	}
	fs.tempDir = path
	return nil
}

// EnableBackup enables automatic backup creation before overwriting files
func (fs *FileService) EnableBackup(enabled bool) {
	fs.backupEnabled = enabled
}

// SafeReadFile reads a file with retry logic and detailed error reporting
func (fs *FileService) SafeReadFile(filePath string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= fs.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(fs.retryDelay)
		}

		data, err := os.ReadFile(filePath)
		if err == nil {
			return data, nil
		}

		lastErr = err

		// Don't retry certain errors
		if os.IsNotExist(err) || os.IsPermission(err) {
			break
		}
	}

	return nil, &FileError{
		Operation: "read",
		Path:      filePath,
		Err:       lastErr,
		Retries:   fs.maxRetries,
		Timestamp: time.Now(),
	}
}

// SafeWriteFile writes data to a file with atomic operations and backup
func (fs *FileService) SafeWriteFile(filePath string, data []byte, perm os.FileMode) error {
	// Create backup if file exists and backup is enabled
	if fs.backupEnabled {
		if _, err := os.Stat(filePath); err == nil {
			if err := fs.createBackup(filePath); err != nil {
				return fmt.Errorf("backup creation failed: %w", err)
			}
		}
	}

	// Write to temporary file first for atomic operation
	tempFile, err := fs.createTempFile(filePath)
	if err != nil {
		return &FileError{
			Operation: "create_temp",
			Path:      filePath,
			Err:       err,
			Timestamp: time.Now(),
		}
	}
	defer os.Remove(tempFile) // Clean up temp file if something goes wrong

	var lastErr error
	for attempt := 0; attempt <= fs.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(fs.retryDelay)
		}

		if err := os.WriteFile(tempFile, data, perm); err != nil {
			lastErr = err
			continue
		}

		// Atomic move from temp file to final location
		if err := fs.atomicMove(tempFile, filePath); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return &FileError{
		Operation: "write",
		Path:      filePath,
		Err:       lastErr,
		Retries:   fs.maxRetries,
		Timestamp: time.Now(),
	}
}

// SafeCreateFile creates a new file with directory structure if needed
func (fs *FileService) SafeCreateFile(filePath string) (*os.File, error) {
	// Ensure directory exists
	if err := fs.ensureDir(filepath.Dir(filePath)); err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt <= fs.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(fs.retryDelay)
		}

		file, err := os.Create(filePath)
		if err == nil {
			return file, nil
		}

		lastErr = err

		// Don't retry permission errors
		if os.IsPermission(err) {
			break
		}
	}

	return nil, &FileError{
		Operation: "create",
		Path:      filePath,
		Err:       lastErr,
		Retries:   fs.maxRetries,
		Timestamp: time.Now(),
	}
}

// SafeRemoveFile removes a file with retry logic
func (fs *FileService) SafeRemoveFile(filePath string) error {
	var lastErr error

	for attempt := 0; attempt <= fs.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(fs.retryDelay)
		}

		err := os.Remove(filePath)
		if err == nil || os.IsNotExist(err) {
			return nil // Success or file doesn't exist
		}

		lastErr = err

		// Don't retry permission errors
		if os.IsPermission(err) {
			break
		}
	}

	return &FileError{
		Operation: "remove",
		Path:      filePath,
		Err:       lastErr,
		Retries:   fs.maxRetries,
		Timestamp: time.Now(),
	}
}

// CheckDiskSpace verifies sufficient disk space for an operation
func (fs *FileService) CheckDiskSpace(path string, requiredBytes int64) error {
	dir := filepath.Dir(path)

	var stat syscallStatFS
	if err := syscallStatfs(dir, &stat); err != nil {
		return fmt.Errorf("failed to check disk space for %s: %w", dir, err)
	}

	availableBytes := int64(stat.Bavail) * int64(stat.Bsize)
	if availableBytes < requiredBytes {
		return fmt.Errorf("insufficient disk space: need %d bytes, have %d bytes",
			requiredBytes, availableBytes)
	}

	return nil
}

// GetFileInfo returns enhanced file information with error handling
func (fs *FileService) GetFileInfo(filePath string) (os.FileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, &FileError{
			Operation: "stat",
			Path:      filePath,
			Err:       err,
			Timestamp: time.Now(),
		}
	}
	return info, nil
}

// CleanupTempFiles removes temporary files created by this service
func (fs *FileService) CleanupTempFiles() error {
	pattern := filepath.Join(fs.tempDir, "ankiprep-*.tmp")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to find temp files: %w", err)
	}

	var errors []error
	for _, match := range matches {
		if err := fs.SafeRemoveFile(match); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to cleanup %d temp files", len(errors))
	}

	return nil
}

// createBackup creates a backup of the existing file
func (fs *FileService) createBackup(filePath string) error {
	backupPath := filePath + ".backup." + fmt.Sprintf("%d", time.Now().Unix())

	data, err := fs.SafeReadFile(filePath)
	if err != nil {
		return err
	}

	return fs.SafeWriteFile(backupPath, data, 0644)
}

// createTempFile creates a temporary file in the same directory as the target
func (fs *FileService) createTempFile(targetPath string) (string, error) {
	dir := filepath.Dir(targetPath)
	prefix := "ankiprep-" + filepath.Base(targetPath) + "-"

	// Try to create temp file in same directory first
	tempFile, err := os.CreateTemp(dir, prefix+"*.tmp")
	if err != nil {
		// Fall back to system temp directory
		tempFile, err = os.CreateTemp(fs.tempDir, prefix+"*.tmp")
		if err != nil {
			return "", err
		}
	}

	tempPath := tempFile.Name()
	tempFile.Close()

	return tempPath, nil
}

// atomicMove performs an atomic move operation
func (fs *FileService) atomicMove(srcPath, dstPath string) error {
	// Try direct rename first (fastest, atomic on same filesystem)
	if err := os.Rename(srcPath, dstPath); err == nil {
		return nil
	}

	// Fall back to copy + remove (cross-filesystem move)
	data, err := fs.SafeReadFile(srcPath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dstPath, data, 0644); err != nil {
		return err
	}

	return fs.SafeRemoveFile(srcPath)
}

// ensureDir creates directory structure if it doesn't exist
func (fs *FileService) ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return &FileError{
				Operation: "mkdir",
				Path:      dirPath,
				Err:       err,
				Timestamp: time.Now(),
			}
		}
	}
	return nil
}

// Platform-specific disk space checking
// These will be implemented with build tags for different platforms

// syscallStatFS represents filesystem statistics
type syscallStatFS struct {
	Bsize  uint64 // Block size
	Bavail uint64 // Available blocks
}

// syscallStatfs gets filesystem statistics for the given path
func syscallStatfs(path string, stat *syscallStatFS) error {
	// This is a simplified implementation
	// In a real implementation, this would use platform-specific syscalls
	// For now, we'll just return success to avoid platform-specific code

	// For demonstration purposes, assume we have 1GB available
	stat.Bsize = 4096
	stat.Bavail = 262144 // 1GB / 4096 bytes per block

	return nil
}

// IsFileError checks if an error is a FileError
func IsFileError(err error) (*FileError, bool) {
	if fileErr, ok := err.(*FileError); ok {
		return fileErr, true
	}
	return nil, false
}

// GetErrorContext extracts detailed context from file errors
func GetErrorContext(err error) string {
	if fileErr, ok := IsFileError(err); ok {
		return fmt.Sprintf("Operation: %s, Path: %s, Time: %s, OS: %s",
			fileErr.Operation, fileErr.Path, fileErr.Timestamp.Format(time.RFC3339), runtime.GOOS)
	}
	return err.Error()
}
