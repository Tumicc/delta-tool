package app

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx          context.Context
	codeLoader   *WeaponCodeLoader
	cacheManager *CacheManager
	enableExcel  bool // Set to true only in development mode
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		codeLoader:   NewWeaponCodeLoader(),
		cacheManager: NewCacheManager(),
		enableExcel:  false, // Default to false for production builds
	}
}

// NewAppWithExcel creates a new App with Excel support enabled
// Use this only in development mode for generating cache
func NewAppWithExcel() *App {
	return &App{
		codeLoader:   NewWeaponCodeLoader(),
		cacheManager: NewCacheManager(),
		enableExcel:  true,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Log cache status
	cachePath := a.cacheManager.GetCachePath()
	fmt.Printf("Cache location: %s\n", cachePath)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// LoadWeaponCodesFromExcel is a helper function for development mode
// It loads all weapon codes from Excel files
func LoadWeaponCodesFromExcel(a *App) ([]WeaponCode, error) {
	return a.LoadWeaponCodes()
}

// GetWeaponCodes returns all weapon modification codes
// In production, it loads from local cache
// In development (if enableExcel is true), it can load from Excel
func (a *App) GetWeaponCodes() []WeaponCode {
	// Try loading from cache first
	codes, found, err := a.cacheManager.Load()
	if err == nil && found {
		return codes
	}

	// If cache not found and Excel is disabled, show error
	if !a.enableExcel {
		fmt.Printf("Error: Weapon codes cache not found. Please ensure the cache file exists at: %s\n",
			a.cacheManager.GetCachePath())
		fmt.Println("To generate the cache file, run: go run . generate-cache")
		return []WeaponCode{}
	}

	// Development mode: load from Excel as fallback
	fmt.Println("Development mode: Loading weapon codes from Excel...")
	codes, err = a.LoadWeaponCodes()
	if err != nil {
		fmt.Printf("Error loading weapon codes: %v\n", err)
		return []WeaponCode{}
	}

	// Save to cache for next time
	if err := a.cacheManager.Save(codes, "local-excel"); err != nil {
		fmt.Printf("Warning: Failed to save cache: %v\n", err)
	}

	return codes
}

// GetWeaponCodesFromDaoZai returns weapon codes from 刀仔 data source
// Filters cached data by source "刀仔"
func (a *App) GetWeaponCodesFromDaoZai() []WeaponCode {
	// Try loading from cache first
	codes, found, err := a.cacheManager.Load()
	fmt.Printf("[DEBUG] GetWeaponCodesFromDaoZai: found=%v, err=%v, cachePath=%s\n", found, err, a.cacheManager.GetCachePath())
	if err == nil && found {
		fmt.Printf("[DEBUG] Loaded %d codes from cache\n", len(codes))
		result := filterBySource(codes, "刀仔")
		fmt.Printf("[DEBUG] Filtered to %d codes from 刀仔\n", len(result))
		return result
	}

	// If cache not found and Excel is enabled, load from Excel
	if a.enableExcel {
		codes, err := a.LoadWeaponCodesFromDaoZai()
		if err != nil {
			fmt.Printf("Error loading 刀仔 weapon codes: %v\n", err)
			return []WeaponCode{}
		}
		// Add source field to each code
		for i := range codes {
			codes[i].Source = "刀仔"
		}
		return codes
	}

	fmt.Printf("Error: Weapon codes cache not found. Please ensure the cache file exists at: %s\n",
		a.cacheManager.GetCachePath())
	return []WeaponCode{}
}

// GetWeaponCodesFromWeaponMaster returns weapon codes from 武器大师 data source
// Filters cached data by source "武器大师"
func (a *App) GetWeaponCodesFromWeaponMaster() []WeaponCode {
	// Try loading from cache first
	codes, found, err := a.cacheManager.Load()
	fmt.Printf("[DEBUG] GetWeaponCodesFromWeaponMaster: found=%v, err=%v, cachePath=%s\n", found, err, a.cacheManager.GetCachePath())
	if err == nil && found {
		fmt.Printf("[DEBUG] Loaded %d codes from cache\n", len(codes))
		result := filterBySource(codes, "武器大师")
		fmt.Printf("[DEBUG] Filtered to %d codes from 武器大师\n", len(result))
		return result
	}

	// If cache not found and Excel is enabled, load from Excel
	if a.enableExcel {
		codes, err := a.LoadWeaponCodesFromWeaponMaster()
		if err != nil {
			fmt.Printf("Error loading 武器大师 weapon codes: %v\n", err)
			return []WeaponCode{}
		}
		// Add source field to each code
		for i := range codes {
			codes[i].Source = "武器大师"
		}
		return codes
	}

	fmt.Printf("Error: Weapon codes cache not found. Please ensure the cache file exists at: %s\n",
		a.cacheManager.GetCachePath())
	return []WeaponCode{}
}

// filterBySource filters weapon codes by data source
func filterBySource(codes []WeaponCode, source string) []WeaponCode {
	var result []WeaponCode
	for _, code := range codes {
		if code.Source == source {
			result = append(result, code)
		}
	}
	return result
}

// GetCacheInfo returns information about the current cache
func (a *App) GetCacheInfo() map[string]interface{} {
	codes, found, err := a.cacheManager.Load()
	info := map[string]interface{}{
		"cache_path":   a.cacheManager.GetCachePath(),
		"cache_found":  found,
		"cache_loaded": err == nil,
		"version":      CacheVersion,
	}

	if err == nil && found {
		info["code_count"] = len(codes)
		// Count by source
		daoZaiCount := 0
		weaponMasterCount := 0
		for _, code := range codes {
			if code.Source == "刀仔" {
				daoZaiCount++
			} else if code.Source == "武器大师" {
				weaponMasterCount++
			}
		}
		info["dao_zai_count"] = daoZaiCount
		info["weapon_master_count"] = weaponMasterCount
	}

	if a.enableExcel {
		info["excel_enabled"] = true
		info["mode"] = "development"
	} else {
		info["excel_enabled"] = false
		info["mode"] = "production"
	}

	return info
}
