package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient represents the interface for fetching weapon codes from a remote API
// This is designed for future integration with a separate data service
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// APIResponse represents the response structure from the remote API
type APIResponse struct {
	Success    bool         `json:"success"`
	Version    string       `json:"version"`
	LastUpdated string      `json:"last_updated"`
	Data       []WeaponCode `json:"data"`
	Message    string       `json:"message,omitempty"`
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchWeaponCodes fetches weapon codes from the remote API
func (api *APIClient) FetchWeaponCodes() ([]WeaponCode, error) {
	url := fmt.Sprintf("%s/api/weapon-codes", api.baseURL)

	resp, err := api.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	fmt.Printf("Fetched %d weapon codes from API (version: %s)\n",
		len(apiResp.Data), apiResp.Version)

	return apiResp.Data, nil
}

// FetchWeaponCodesWithMode fetches weapon codes filtered by mode
func (api *APIClient) FetchWeaponCodesWithMode(mode string) ([]WeaponCode, error) {
	url := fmt.Sprintf("%s/api/weapon-codes?mode=%s", api.baseURL, mode)

	resp, err := api.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	return apiResp.Data, nil
}

// DataSourceConfig represents the configuration for data sources
type DataSourceConfig struct {
	// UseLocalCache indicates whether to use local cache file
	UseLocalCache bool
	// LocalCachePath specifies the path to local cache (empty for default)
	LocalCachePath string
	// APIBaseURL specifies the base URL for remote API (empty to disable)
	APIBaseURL string
	// CacheMaxAge specifies how long the local cache is valid (0 = forever)
	CacheMaxAge time.Duration
}

// WeaponCodeLoader handles loading weapon codes from various sources
type WeaponCodeLoader struct {
	cacheManager *CacheManager
	apiClient    *APIClient
	config       DataSourceConfig
}

// NewWeaponCodeLoader creates a new weapon code loader with default config
func NewWeaponCodeLoader() *WeaponCodeLoader {
	return &WeaponCodeLoader{
		cacheManager: NewCacheManager(),
		config: DataSourceConfig{
			UseLocalCache: true,
			CacheMaxAge:   24 * time.Hour, // Default: cache valid for 24 hours
		},
	}
}

// NewWeaponCodeLoaderWithConfig creates a new weapon code loader with custom config
func NewWeaponCodeLoaderWithConfig(config DataSourceConfig) *WeaponCodeLoader {
	loader := &WeaponCodeLoader{
		cacheManager: NewCacheManager(),
		config:       config,
	}

	if config.APIBaseURL != "" {
		loader.apiClient = NewAPIClient(config.APIBaseURL)
	}

	return loader
}

// Load loads weapon codes using the configured data sources
// It tries sources in this order:
// 1. API (if configured and cache is expired)
// 2. Local cache (if enabled and valid)
// 3. Returns error if no source is available
func (loader *WeaponCodeLoader) Load() ([]WeaponCode, error) {
	// Try API first if configured and cache should be refreshed
	if loader.apiClient != nil {
		shouldRefresh := true
		if loader.config.UseLocalCache && loader.config.CacheMaxAge > 0 {
			expired, err := loader.cacheManager.IsCacheExpired(loader.config.CacheMaxAge)
			if err == nil && !expired {
				shouldRefresh = false
			}
		}

		if shouldRefresh {
			fmt.Println("Fetching weapon codes from API...")
			codes, err := loader.apiClient.FetchWeaponCodes()
			if err == nil {
				// Update cache with fresh data
				if loader.config.UseLocalCache {
					go loader.cacheManager.Save(codes, "api")
				}
				return codes, nil
			}
			fmt.Printf("API fetch failed: %v, falling back to cache\n", err)
		}
	}

	// Try local cache
	if loader.config.UseLocalCache {
		codes, found, err := loader.cacheManager.Load()
		if err == nil && found {
			return codes, nil
		}
		if err != nil {
			fmt.Printf("Cache load failed: %v\n", err)
		}
	}

	return nil, fmt.Errorf("no weapon codes available (API not configured or failed, cache not available)")
}
