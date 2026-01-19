package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// Current cache version
	CacheVersion = "1.0.0"
	// Cache filename
	CacheFileName = "weapon_codes.json"
)

// WeaponCodeCache represents the cache structure with version control
type WeaponCodeCache struct {
	Version     string       `json:"version"`
	LastUpdated string       `json:"last_updated"`
	TotalCount  int          `json:"total_count"`
	DataSource  string       `json:"data_source"` // "local", "api", etc.
	WeaponCodes []WeaponCode `json:"weapon_codes"`
}

// CacheManager manages the weapon codes cache
type CacheManager struct {
	cachePath string
	mu        sync.RWMutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	return &CacheManager{
		cachePath: getCacheFilePath(),
	}
}

// getCacheFilePath returns the path to the cache file
// It checks multiple locations: local data dir, executable dir, etc.
func getCacheFilePath() string {
	possiblePaths := []string{
		filepath.Join("data", CacheFileName),
		CacheFileName,
	}

	// Check executable directory
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		possiblePaths = append(possiblePaths, filepath.Join(exeDir, CacheFileName))
		possiblePaths = append(possiblePaths, filepath.Join(exeDir, "data", CacheFileName))
		possiblePaths = append(possiblePaths, filepath.Join(filepath.Dir(exeDir), "data", CacheFileName))
	}

	// Return first existing path, or default to data dir
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Default to data directory in current working directory
	return filepath.Join("data", CacheFileName)
}

// Load loads weapon codes from cache
// Returns the codes and a boolean indicating if cache was used
func (cm *CacheManager) Load() ([]WeaponCode, bool, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Check if cache file exists
	if _, err := os.Stat(cm.cachePath); os.IsNotExist(err) {
		return nil, false, fmt.Errorf("cache file not found")
	}

	// Read cache file
	data, err := os.ReadFile(cm.cachePath)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read cache file: %w", err)
	}

	// Parse cache
	var cache WeaponCodeCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, false, fmt.Errorf("failed to parse cache file: %w", err)
	}

	// Validate cache version (for future compatibility)
	if cache.Version != CacheVersion {
		// Version mismatch, but we can still try to use it
		fmt.Printf("Warning: Cache version mismatch. Expected %s, got %s\n", CacheVersion, cache.Version)
	}

	fmt.Printf("Loaded %d weapon codes from cache (version: %s, updated: %s)\n",
		len(cache.WeaponCodes), cache.Version, cache.LastUpdated)

	return cache.WeaponCodes, true, nil
}

// Save saves weapon codes to cache
func (cm *CacheManager) Save(codes []WeaponCode, dataSource string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Ensure directory exists
	dir := filepath.Dir(cm.cachePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Create cache structure
	cache := WeaponCodeCache{
		Version:     CacheVersion,
		LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
		TotalCount:  len(codes),
		DataSource:  dataSource,
		WeaponCodes: codes,
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	// Write to file
	if err := os.WriteFile(cm.cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	fmt.Printf("Saved %d weapon codes to cache: %s\n", len(codes), cm.cachePath)

	return nil
}

// GetCachePath returns the current cache file path
func (cm *CacheManager) GetCachePath() string {
	return cm.cachePath
}

// IsCacheExpired checks if the cache is older than the specified duration
func (cm *CacheManager) IsCacheExpired(maxAge time.Duration) (bool, error) {
	info, err := os.Stat(cm.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // Cache doesn't exist, consider it expired
		}
		return false, err
	}

	age := time.Since(info.ModTime())
	return age > maxAge, nil
}

// ClearCache removes the cache file
func (cm *CacheManager) ClearCache() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if err := os.Remove(cm.cachePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache file: %w", err)
	}

	fmt.Printf("Cache file removed: %s\n", cm.cachePath)
	return nil
}
