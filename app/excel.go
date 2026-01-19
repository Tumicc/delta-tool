package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// WeaponCode represents a single weapon modification code entry
type WeaponCode struct {
	ID         string  `json:"id"`
	Mode       string  `json:"mode"`        // 烽火地带 or 全面战场
	Name       string  `json:"name"`        // 枪械名称
	Tier       string  `json:"tier"`        // 版本排行/类型
	Price      *int    `json:"price"`       // 改装价格（万），null 表示无数据
	Build      string  `json:"build"`       // 改装描述
	Code       string  `json:"code"`        // 改枪码
	Range      *int    `json:"range"`       // 有效射程（米），null 表示无数据
	UpdateTime *string `json:"update_time"` // 更新时间，null 表示无数据
	Source     string  `json:"source"`      // 数据来源: "刀仔" or "武器大师"
}

// LoadWeaponCodes reads the Excel file and returns all weapon codes
// Supports two file formats:
// 1. 刀仔三角洲枪械改装.xlsx - Single sheet with data in multiple columns
// 2. 武器大师地板的改枪码合集.xlsx - Multiple sheets with single column format
// This method tries both sources and returns combined results
func (a *App) LoadWeaponCodes() ([]WeaponCode, error) {
	var allCodes []WeaponCode

	// Try to load from 刀仔 file
	daoZaiCodes, err := a.LoadWeaponCodesFromDaoZai()
	if err == nil && len(daoZaiCodes) > 0 {
		allCodes = append(allCodes, daoZaiCodes...)
	}

	// Try to load from 武器大师 file
	weaponMasterCodes, err := a.LoadWeaponCodesFromWeaponMaster()
	if err == nil && len(weaponMasterCodes) > 0 {
		allCodes = append(allCodes, weaponMasterCodes...)
	}

	if len(allCodes) == 0 {
		return nil, fmt.Errorf("no weapon codes found from any source")
	}

	return allCodes, nil
}

// LoadWeaponCodesFromDaoZai loads weapon codes from 刀仔 data source
func (a *App) LoadWeaponCodesFromDaoZai() ([]WeaponCode, error) {
	possiblePaths := []string{
		"data/刀仔三角洲枪械改装.xlsx",
		"",
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		possiblePaths[1] = filepath.Join(exeDir, "data", "刀仔三角洲枪械改装.xlsx")
		possiblePaths = append(possiblePaths, filepath.Join(filepath.Dir(exeDir), "data", "刀仔三角洲枪械改装.xlsx"))
	}

	var f *excelize.File
	var lastErr error
	for _, path := range possiblePaths {
		if path == "" {
			continue
		}
		f, err = excelize.OpenFile(path)
		if err == nil {
			break
		}
		lastErr = err
	}

	if f == nil {
		return nil, fmt.Errorf("failed to open 刀仔 file, last error: %w", lastErr)
	}
	defer f.Close()

	return parse刀仔File(f)
}

// LoadWeaponCodesFromWeaponMaster loads weapon codes from 武器大师 data source
func (a *App) LoadWeaponCodesFromWeaponMaster() ([]WeaponCode, error) {
	possiblePaths := []string{
		"data/武器大师地板的改枪码合集.xlsx",
		"",
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		possiblePaths[1] = filepath.Join(exeDir, "data", "武器大师地板的改枪码合集.xlsx")
		possiblePaths = append(possiblePaths, filepath.Join(filepath.Dir(exeDir), "data", "武器大师地板的改枪码合集.xlsx"))
	}

	var f *excelize.File
	var lastErr error
	for _, path := range possiblePaths {
		if path == "" {
			continue
		}
		f, err = excelize.OpenFile(path)
		if err == nil {
			break
		}
		lastErr = err
	}

	if f == nil {
		return nil, fmt.Errorf("failed to open 武器大师 file, last error: %w", lastErr)
	}
	defer f.Close()

	return parseWeaponMasterFile(f)
}

// parse刀仔File parses the 刀仔 format file (single sheet, multiple columns)
func parse刀仔File(f *excelize.File) ([]WeaponCode, error) {
	var codes []WeaponCode
	id := 0
	lastFireName := "" // Track last weapon name for 烽火地带
	lastFullName := "" // Track last weapon name for 全面战场

	// Read rows starting from row 11 (header is at row 10, 0-indexed as row 10)
	// Excel row 11 is where data starts
	for rowNum := 11; rowNum <= 500; rowNum++ {
		// Read row data directly using GetCellValue for reliability
		rowData := readRowDirectly(f, "工作表1", rowNum)

		// Check if row has any data
		hasData := false
		for _, cell := range rowData {
			if cell != "" {
				hasData = true
				break
			}
		}

		if !hasData {
			continue
		}

		// Check for ad rows and skip them
		if isAdRow(rowData) {
			continue
		}

		// Skip header row
		if rowNum == 11 {
			continue
		}

		// Parse 烽火地带 region (columns A-G, indices 0-6)
		if fireCodes := parseFireRegion(rowData, rowNum, &id, &lastFireName); len(fireCodes) > 0 {
			codes = append(codes, fireCodes...)
		}

		// Parse 全面战场 region (columns I-K, indices 8-10)
		if fullCodes := parseFullRegion(rowData, rowNum, &id, &lastFullName); len(fullCodes) > 0 {
			codes = append(codes, fullCodes...)
		}
	}

	return codes, nil
}

// parseWeaponMasterFile parses the 武器大师地板 format file (multiple sheets, single column)
func parseWeaponMasterFile(f *excelize.File) ([]WeaponCode, error) {
	var codes []WeaponCode
	id := 0

	// Parse both sheets
	sheets := map[string]string{
		"烽火地带": "烽火地带",
		"全面战场": "全面战场",
	}

	for sheetName, mode := range sheets {
		sheetCodes, err := parseWeaponMasterSheet(f, sheetName, mode, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sheet %s: %w", sheetName, err)
		}
		codes = append(codes, sheetCodes...)
	}

	return codes, nil
}

// parseWeaponMasterSheet parses a single sheet from 武器大师地板 format
func parseWeaponMasterSheet(f *excelize.File, sheetName, mode string, id *int) ([]WeaponCode, error) {
	var codes []WeaponCode

	// Find data start row (skip ads and headers)
	dataStartRow := 1
	for rowNum := 1; rowNum <= 10; rowNum++ {
		val, _ := f.GetCellValue(sheetName, fmt.Sprintf("A%d", rowNum))
		if strings.Contains(val, "步枪") {
			dataStartRow = rowNum + 1
			break
		}
	}

	// Read each row and parse cells directly
	for rowNum := dataStartRow; rowNum <= 500; rowNum++ {
		// Read columns A, B, C (枪名, 价格+描述, 改枪码)
		// Then skip to E, F, G (下一组数据)
		// Pattern: A,B,C | E,F,G | I,J,K
		colGroups := [][]int{{1, 2, 3}, {5, 6, 7}, {9, 10, 11}}

		for _, cols := range colGroups {
			nameCol, priceCol, codeCol := cols[0], cols[1], cols[2]

			colName1, _ := excelize.CoordinatesToCellName(nameCol, rowNum)
			colName2, _ := excelize.CoordinatesToCellName(priceCol, rowNum)
			colName3, _ := excelize.CoordinatesToCellName(codeCol, rowNum)

			nameVal, _ := f.GetCellValue(sheetName, colName1)
			name := strings.TrimSpace(nameVal)

			priceVal, _ := f.GetCellValue(sheetName, colName2)
			priceAndBuild := strings.TrimSpace(priceVal)

			codeVal, _ := f.GetCellValue(sheetName, colName3)
			code := strings.TrimSpace(codeVal)

			// Skip if no code or if it's an ad
			if code == "" || isWeaponMasterAdCell(code) {
				continue
			}

			// Check if this looks like valid data (code starts with 6 and is 21 chars)
			if !strings.HasPrefix(code, "6") || len(code) != 21 {
				continue
			}

			// Skip if no name
			if name == "" {
				continue
			}

			// Parse the data based on mode
			var sheetCodes []WeaponCode
			if mode == "烽火地带" {
				sheetCodes = parseWeaponMasterFireData(name, priceAndBuild, code, id)
			} else {
				sheetCodes = parseWeaponMasterFullData(name, priceAndBuild, code, id)
			}

			if len(sheetCodes) > 0 {
				codes = append(codes, sheetCodes...)
			}
		}
	}

	return codes, nil
}

// parseWeaponMasterFireData parses 烽火地带 data from 武器大师 file
func parseWeaponMasterFireData(name, priceAndBuild, code string, id *int) []WeaponCode {
	// Parse price
	var price *int
	build := "标准改装"

	if priceAndBuild != "" {
		if p := parsePrice(priceAndBuild); p != nil {
			price = p
		}
		// Remove price part to get build description
		build = strings.TrimSpace(priceAndBuild)
		re := regexp.MustCompile(`\d+W`)
		build = re.ReplaceAllString(build, "")
		build = strings.TrimSpace(build)
		if build == "" {
			build = "标准改装"
		}
	}

	// Infer tier from name
	tier := inferTierFromName(name)
	if tier == "-" {
		tier = "T0"
	}

	codes := []WeaponCode{
		{
			ID:         fmt.Sprintf("%d", *id),
			Mode:       "烽火地带",
			Name:       name,
			Tier:       tier,
			Price:      price,
			Build:      build,
			Code:       code,
			Range:      nil,
			UpdateTime: nil,
			Source:     "武器大师",
		},
	}
	*id++
	return codes
}

// parseWeaponMasterFullData parses 全面战场 data from 武器大师 file
func parseWeaponMasterFullData(name, priceAndBuild, code string, id *int) []WeaponCode {
	// Parse price
	var price *int
	build := "标准配置"

	if priceAndBuild != "" {
		// Try to extract price number
		re := regexp.MustCompile(`(\d+)`)
		matches := re.FindStringSubmatch(priceAndBuild)
		if len(matches) >= 2 {
			if val, err := strconv.Atoi(matches[1]); err == nil {
				price = &val
			}
		}

		// Extract build description (non-number part)
		build = re.ReplaceAllString(priceAndBuild, "")
		build = strings.TrimSpace(build)
		if build == "" {
			build = "标准配置"
		}
	}

	// Infer tier from name
	tier := inferTierFromName(name)

	codes := []WeaponCode{
		{
			ID:         fmt.Sprintf("%d", *id),
			Mode:       "全面战场",
			Name:       name,
			Tier:       tier,
			Price:      price,
			Build:      build,
			Code:       code,
			Range:      nil,
			UpdateTime: nil,
			Source:     "武器大师",
		},
	}
	*id++
	return codes
}

// parseWeaponMasterFireCell parses a single cell from 烽火地带 sheet
// Format: "MK47 22W青春版 6IDP1280B97T7MULLRJ3C" or "QCQ171 26W青春版 6IG8E6O07OULUBJA9PRPI"
func parseWeaponMasterFireCell(rowData []string, id *int) []WeaponCode {
	var codes []WeaponCode

	for _, cell := range rowData {
		cell = strings.TrimSpace(cell)

		// Skip ads and invalid content
		if isWeaponMasterAdCell(cell) {
			continue
		}

		// Parse the cell format: "枪名 价格+描述 改枪码"
		// Split by spaces and parse
		parts := strings.Fields(cell)
		if len(parts) < 2 {
			continue
		}

		// Find the code (usually starts with 6 and is long)
		var code string
		var nameParts []string
		var priceAndBuild string

		for i, part := range parts {
			if strings.HasPrefix(part, "6") && len(part) > 20 {
				code = part
				// Everything before code is name + price+build
				if i > 0 {
					priceAndBuild = parts[i-1]
					nameParts = parts[:i-1]
				}
				break
			}
		}

		if code == "" {
			continue
		}

		// Extract name from nameParts
		name := strings.Join(nameParts, " ")
		if name == "" {
			continue
		}

		// Parse price and build from priceAndBuild
		var price *int
		var build string
		if priceAndBuild != "" {
			if p := parsePrice(priceAndBuild); p != nil {
				price = p
			}
			// Remove price part to get build description
			build = strings.TrimSpace(priceAndBuild)
			re := regexp.MustCompile(`\d+W`)
			build = re.ReplaceAllString(build, "")
			build = strings.TrimSpace(build)
			if build == "" {
				build = "标准改装"
			}
		} else {
			build = "标准改装"
		}

		// Infer tier from name
		tier := inferTierFromName(name)
		if tier == "-" {
			tier = "T0" // 武器大师 file defaults to T0
		}

		codes = append(codes, WeaponCode{
			ID:         fmt.Sprintf("%d", *id),
			Mode:       "烽火地带",
			Name:       name,
			Tier:       tier,
			Price:      price,
			Build:      build,
			Code:       code,
			Range:      nil,
			UpdateTime: nil,
			Source:     "武器大师",
		})
		*id++
	}

	return codes
}

// parseWeaponMasterFullCell parses a single cell from 全面战场 sheet
// Format: "M4A1 6IENQK0097PFORHQ0UK53 30" or "M4A1 6IENQT0097PFORHQ0UK53 60腰射"
func parseWeaponMasterFullCell(rowData []string, id *int) []WeaponCode {
	var codes []WeaponCode

	for _, cell := range rowData {
		cell = strings.TrimSpace(cell)

		// Skip ads and invalid content
		if isWeaponMasterAdCell(cell) {
			continue
		}

		// Parse the cell format: "枪名 改枪码 价格" or "枪名 改枪码 价格+描述"
		parts := strings.Fields(cell)
		if len(parts) < 2 {
			continue
		}

		// First part is name
		name := parts[0]
		if name == "" {
			continue
		}

		// Second part is code (starts with 6)
		var code string
		var priceAndBuild string

		if len(parts) >= 2 && strings.HasPrefix(parts[1], "6") && len(parts[1]) > 20 {
			code = parts[1]
		}

		if code == "" {
			continue
		}

		// Third part (if exists) is price + build
		if len(parts) >= 3 {
			priceAndBuild = parts[2]
		}

		// Parse price
		var price *int
		build := "标准配置"

		if priceAndBuild != "" {
			// Try to extract price number
			re := regexp.MustCompile(`(\d+)`)
			matches := re.FindStringSubmatch(priceAndBuild)
			if len(matches) >= 2 {
				if val, err := strconv.Atoi(matches[1]); err == nil {
					price = &val
				}
			}

			// Extract build description (non-number part)
			build = re.ReplaceAllString(priceAndBuild, "")
			build = strings.TrimSpace(build)
			if build == "" {
				build = "标准配置"
			}
		}

		// Infer tier from name
		tier := inferTierFromName(name)

		codes = append(codes, WeaponCode{
			ID:         fmt.Sprintf("%d", *id),
			Mode:       "全面战场",
			Name:       name,
			Tier:       tier,
			Price:      price,
			Build:      build,
			Code:       code,
			Range:      nil,
			UpdateTime: nil,
			Source:     "武器大师",
		})
		*id++
	}

	return codes
}

// isWeaponMasterAdCell checks if a cell contains ad content from 武器大师 file
func isWeaponMasterAdCell(cell string) bool {
	if len(cell) > 100 {
		return true // Very long content is likely an ad
	}

	adKeywords := []string{
		"抖音", "刀仔", "武器大师", "地板", "改枪码大全",
		"每次使用点链接", "在线文档", "失效",
		"保存好链接", "永久更新", "S7最新版",
		"被抄袭", "被超越", "屏息",
		"射手步枪以及狙击步枪", "霰弹枪以及其它",
		"高手版", "陈泽杯",
	}

	cellLower := strings.ToLower(cell)
	for _, keyword := range adKeywords {
		if strings.Contains(cellLower, keyword) {
			return true
		}
	}

	return false
}

// readRowDirectly reads a row by accessing each cell directly
func readRowDirectly(f *excelize.File, sheet string, rowNum int) []string {
	// Read columns A through L (0-11)
	var row []string
	for colIdx := 0; colIdx <= 11; colIdx++ {
		colName, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
		val, err := f.GetCellValue(sheet, colName)
		if err != nil {
			val = ""
		}
		row = append(row, strings.TrimSpace(val))
	}
	return row
}

// isAdRow checks if a row contains advertisement content
func isAdRow(row []string) bool {
	adKeywords := []string{"抖音搜", "画质调整", "刀仔", "关注", "群", "频道"}
	for _, cell := range row {
		cellLower := strings.ToLower(cell)
		for _, keyword := range adKeywords {
			if strings.Contains(cellLower, keyword) {
				return true
			}
		}
	}
	return false
}

// parseFireRegion parses 烽火地带 data from columns A-G
func parseFireRegion(row []string, rowIndex int, id *int, lastName *string) []WeaponCode {
	// Need at least column E (code)
	if len(row) < 5 {
		return nil
	}

	name := strings.TrimSpace(getCellValue(row, 0))
	code := strings.TrimSpace(getCellValue(row, 4))

	// Skip if no code
	if code == "" {
		return nil
	}

	// If name is empty, use the last name (continuation from previous row)
	if name == "" {
		name = *lastName
	}

	// Skip if still no weapon name
	if name == "" {
		return nil
	}

	// Check for header row or invalid data
	if strings.Contains(name, "枪械名称") || strings.Contains(code, "枪械代码") {
		return nil
	}

	// Check if this looks like valid data (name shouldn't be too long and shouldn't contain ad keywords)
	if len(name) > 50 || isAdRow([]string{name}) {
		return nil
	}

	// Update last name if this row has a new name
	if strings.TrimSpace(getCellValue(row, 0)) != "" {
		*lastName = name
	}

	// Parse tier (version)
	tier := strings.TrimSpace(getCellValue(row, 1))
	if tier == "" {
		tier = "-"
	}

	// Parse price (column 2) - handle "85w" format
	var price *int
	priceStr := strings.TrimSpace(getCellValue(row, 2))
	if priceStr != "" {
		if p := parsePrice(priceStr); p != nil {
			price = p
		}
	}

	// Parse build description (column 3)
	build := strings.TrimSpace(getCellValue(row, 3))
	if build == "" {
		build = "标准改装"
	}

	// Parse range (column 5) - handle "52米" format
	var rangeValue *int
	rangeStr := strings.TrimSpace(getCellValue(row, 5))
	if rangeStr != "" {
		if r := parseRange(rangeStr); r != nil {
			rangeValue = r
		}
	}

	// Parse update time (column 6)
	var updateTime *string
	timeStr := strings.TrimSpace(getCellValue(row, 6))
	if timeStr != "" {
		updateTime = &timeStr
	}

	result := []WeaponCode{
		{
			ID:         fmt.Sprintf("%d", *id),
			Mode:       "烽火地带",
			Name:       name,
			Tier:       tier,
			Price:      price,
			Build:      build,
			Code:       code,
			Range:      rangeValue,
			UpdateTime: updateTime,
			Source:     "刀仔",
		},
	}
	*id++
	return result
}

// parseFullRegion parses 全面战场 data from columns I-K
func parseFullRegion(row []string, rowIndex int, id *int, lastName *string) []WeaponCode {
	// Need at least column K (code)
	if len(row) < 11 {
		return nil
	}

	name := strings.TrimSpace(getCellValue(row, 8))
	code := strings.TrimSpace(getCellValue(row, 10))

	// Skip if no code
	if code == "" {
		return nil
	}

	// If name is empty, use the last name (continuation from previous row)
	if name == "" {
		name = *lastName
	}

	// Skip if still no weapon name
	if name == "" {
		return nil
	}

	// Check for header row or invalid data
	if strings.Contains(name, "枪械名称") || strings.Contains(code, "改枪码") {
		return nil
	}

	// Check if this looks like valid data
	if len(name) > 50 || isAdRow([]string{name}) {
		return nil
	}

	// Update last name if this row has a new name
	if strings.TrimSpace(getCellValue(row, 8)) != "" {
		*lastName = name
	}

	// Parse build/style (column 9)
	build := strings.TrimSpace(getCellValue(row, 9))
	if build == "" {
		build = "标准配置"
	}

	// Infer tier from weapon name for 全面战场
	tier := inferTierFromName(name)

	result := []WeaponCode{
		{
			ID:         fmt.Sprintf("%d", *id),
			Mode:       "全面战场",
			Name:       name,
			Tier:       tier,
			Price:      nil, // 全面战场 has no price data
			Build:      build,
			Code:       code,
			Range:      nil, // 全面战场 has no range data
			UpdateTime: nil, // 全面战场 has no update time
			Source:     "刀仔",
		},
	}
	*id++
	return result
}

// parsePrice converts price string like "85w" to integer 85
func parsePrice(s string) *int {
	s = strings.ToLower(strings.TrimSpace(s))
	re := regexp.MustCompile(`(\d+)w`)
	matches := re.FindStringSubmatch(s)
	if len(matches) >= 2 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			return &val
		}
	}
	// Try parsing as plain number
	if val, err := strconv.Atoi(s); err == nil {
		return &val
	}
	return nil
}

// parseRange converts range string like "52米" to integer 52
func parseRange(s string) *int {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) >= 2 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			return &val
		}
	}
	return nil
}

// inferTierFromName infers weapon tier/type from weapon name
func inferTierFromName(name string) string {
	nameLower := strings.ToLower(name)

	// Check for weapon types
	typeMap := map[string]string{
		"狙击":    "狙击",
		"连狙":    "连狙",
		"弓":     "弓弩",
		"手枪":    "手枪",
		"冲锋枪":   "冲锋枪",
		"步枪":    "步枪",
		"机枪":    "机枪",
		"霰弹枪":   "霰弹枪",
		"发射器":   "发射器",
	}

	for keyword, tier := range typeMap {
		if strings.Contains(nameLower, keyword) {
			return tier
		}
	}

	return "-"
}

// getCellValue safely gets a cell value by index
func getCellValue(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return row[index]
}
