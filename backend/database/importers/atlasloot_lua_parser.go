package importers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type LuaLootItem struct {
	ID       int
	SpellID  int
	DropRate string
	Name     string
	Icon     string
	Quality  int
}

type SpecialMenuItem struct {
	Key  string
	Name string
}

type AtlasLootLuaParser struct {
	AddonDir string
}

func NewAtlasLootLuaParser(addonDir string) *AtlasLootLuaParser {
	return &AtlasLootLuaParser{AddonDir: addonDir}
}

// ParseTableRegister parses TableRegister.lua for display names
func (p *AtlasLootLuaParser) ParseTableRegister() (map[string]string, error) {
	displayNames := make(map[string]string)
	path := filepath.Join(p.AddonDir, "Database", "TableRegister.lua")

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	keyRegex := regexp.MustCompile(`\["(\w+)"\]\s*=`)
	alRegex := regexp.MustCompile(`AL\["([^"]+)"\]`)
	quoteRegex := regexp.MustCompile(`\{\s*"([^"]+)"`)

	var currentKey string
	var buffer []string

	for scanner.Scan() {
		line := scanner.Text()

		keyMatch := keyRegex.FindStringSubmatch(line)
		if len(keyMatch) > 1 {
			currentKey = keyMatch[1]
			buffer = []string{line}
		} else if currentKey != "" {
			buffer = append(buffer, strings.TrimSpace(line))
			combined := strings.Join(buffer, " ")

			if strings.Contains(combined, "AtlasLoot") && strings.Contains(combined, "Items") {
				// Extract display name
				alMatches := alRegex.FindAllStringSubmatch(combined, -1)
				if len(alMatches) > 0 {
					var parts []string
					for _, m := range alMatches {
						if m[1] != "Rare" && m[1] != "Summon" && m[1] != "Quest" && m[1] != "Enchants" {
							parts = append(parts, m[1])
						}
					}
					if len(parts) > 0 {
						displayNames[currentKey] = strings.Join(parts, " - ")
					}
				} else {
					qMatch := quoteRegex.FindStringSubmatch(combined)
					if len(qMatch) > 1 {
						displayNames[currentKey] = qMatch[1]
					}
				}
				currentKey = ""
				buffer = nil
			}
		}
	}

	return displayNames, scanner.Err()
}

// ParseLootTables parses a Lua file in Database/ to extract loot tables
func (p *AtlasLootLuaParser) ParseLootTables(filename string) (map[string][]LuaLootItem, error) {
	tables := make(map[string][]LuaLootItem)
	path := filepath.Join(p.AddonDir, "Database", filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return tables, nil
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(contentBytes)

	// Regex patterns
	tableStartRegex := regexp.MustCompile(`^\s*(\w+)\s*=\s*\{`)

	// { 16686, "0.02%" } - Sets style
	shortItemRegex := regexp.MustCompile(`\{\s*(\d+)\s*,\s*"([^"]*%)"\s*\}`)

	// { "s2329", "inv_potion_56", "=q1=Name", ... } - Crafting style
	// { "e7421", ... } - Enchanting style
	// { 22384, "INV_Hammer_08", "=q4=Name", ... }
	detailedItemRegex := regexp.MustCompile(`\{\s*"?([se]?\d+)"?\s*,\s*"([^"]*)"\s*,\s*"([^"]*)"`)

	// { 12345, ... } - Fallback
	simpleIDRegex := regexp.MustCompile(`\{\s*(\d+)\s*,`)

	lines := strings.Split(content, "\n")
	var currentTable string
	var currentItems []LuaLootItem

	for _, line := range lines {
		// Table start
		if match := tableStartRegex.FindStringSubmatch(line); len(match) > 1 {
			// Save previous table if it wasn't closed properly or handled
			if currentTable != "" && len(currentItems) > 0 {
				tables[currentTable] = currentItems
			}
			currentTable = match[1]
			currentItems = []LuaLootItem{}
			// fmt.Printf("DEBUG: Found table start: %s\n", currentTable)
			continue
		}

		// Item entry
		if strings.Contains(line, "{") && currentTable != "" {
			item := LuaLootItem{}

			// Try detailed first (most restrictive)
			if match := detailedItemRegex.FindStringSubmatch(line); len(match) > 3 {
				rawID := match[1]
				item.Icon = match[2]
				rawName := match[3]

				// Handle ID
				if strings.HasPrefix(rawID, "s") {
					item.SpellID, _ = strconv.Atoi(rawID[1:])
				} else if strings.HasPrefix(rawID, "e") {
					item.SpellID, _ = strconv.Atoi(rawID[1:])
				} else {
					item.ID, _ = strconv.Atoi(rawID)
				}

				// Handle Name & Quality
				// =q1=Name
				if strings.HasPrefix(rawName, "=q") && len(rawName) > 2 {
					qChar := rawName[2] // '1'
					item.Quality, _ = strconv.Atoi(string(qChar))
					if len(rawName) > 4 {
						item.Name = rawName[4:]
					}
				} else {
					item.Name = rawName
				}
			} else if match := shortItemRegex.FindStringSubmatch(line); len(match) > 2 {
				item.ID, _ = strconv.Atoi(match[1])
				item.DropRate = match[2]
			} else if match := simpleIDRegex.FindStringSubmatch(line); len(match) > 1 {
				item.ID, _ = strconv.Atoi(match[1])
			}

			if item.ID > 0 || item.SpellID > 0 {
				currentItems = append(currentItems, item)
			}
		}

		// Table end
		if strings.Contains(line, "};") && currentTable != "" {
			if len(currentItems) > 0 {
				tables[currentTable] = currentItems
				// fmt.Printf("DEBUG: Saved table: %s with %d items\n", currentTable, len(currentItems))
			} else {
				// fmt.Printf("DEBUG: Table %s has 0 items\n", currentTable)
			}
			currentTable = ""
			currentItems = nil
		}
	}

	if currentTable != "" && len(currentItems) > 0 {
		tables[currentTable] = currentItems
	}

	return tables, nil
}

// ParseSpecialMenus extracts the mapping of Special Menus (e.g. T0SET) to their sub-tables
func (p *AtlasLootLuaParser) ParseSpecialMenus() (map[string][]SpecialMenuItem, error) {
	// 1. Get AtlasLoot_MenuList from Core/AtlasLoot.lua
	corePath := filepath.Join(p.AddonDir, "Core", "AtlasLoot.lua")
	coreContentBytes, err := os.ReadFile(corePath)
	if err != nil {
		return nil, err
	}
	coreContent := string(coreContentBytes)

	// Extract AtlasLoot_MenuList = { ... }
	menuListRegex := regexp.MustCompile(`(?s)AtlasLoot_MenuList\s*=\s*{(.*?)}`)
	match := menuListRegex.FindStringSubmatch(coreContent)
	if len(match) < 2 {
		return nil, fmt.Errorf("AtlasLoot_MenuList not found")
	}

	menuListContent := match[1]
	menuMap := make(map[string]string)
	// ["KEY"] = "Function"
	entryRegex := regexp.MustCompile(`\["([^"]+)"\]\s*=\s*"([^"]+)"`)

	for _, line := range strings.Split(menuListContent, "\n") {
		m := entryRegex.FindStringSubmatch(line)
		if len(m) > 2 {
			menuMap[m[1]] = m[2]
		}
	}

	specialMenus := make(map[string][]SpecialMenuItem)

	// 2. Scan Core/*.lua for function definitions
	err = filepath.Walk(filepath.Join(p.AddonDir, "Core"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".lua") {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(contentBytes)

		for menuKey, funcName := range menuMap {
			// Find function funcName() ... end
			funcStartRegex := regexp.MustCompile(`(?m)^function\s+` + regexp.QuoteMeta(funcName) + `\s*\(`)
			loc := funcStartRegex.FindStringIndex(content)
			if loc == nil {
				continue
			}

			startIdx := loc[0]
			// Find end of function (naive: next function start or EOF)
			nextFuncRegex := regexp.MustCompile(`(?m)^function\s+`)
			nextLoc := nextFuncRegex.FindStringIndex(content[loc[1]:])

			var body string
			if nextLoc == nil {
				body = content[startIdx:]
			} else {
				endIdx := loc[1] + nextLoc[0]
				body = content[startIdx:endIdx]
			}

			// Parse items in body
			lines := strings.Split(body, "\n")
			var items []SpecialMenuItem
			var currentName string

			setTextRegex := regexp.MustCompile(`:SetText\((.*)\)`)
			lootPageRegex := regexp.MustCompile(`\.lootpage\s*=\s*"([^"]+)"`)
			alRegex := regexp.MustCompile(`AL\["([^"]+)"\]`)
			strRegex := regexp.MustCompile(`"([^"]+)"`)

			for _, line := range lines {
				// Parse Name
				if match := setTextRegex.FindStringSubmatch(line); len(match) > 1 {
					rawText := match[1]
					alMatches := alRegex.FindAllStringSubmatch(rawText, -1)
					var parts []string
					if len(alMatches) > 0 {
						for _, m := range alMatches {
							parts = append(parts, m[1])
						}
					} else {
						// Fallback to literal strings
						strMatches := strRegex.FindAllStringSubmatch(rawText, -1)
						for _, m := range strMatches {
							// skip colors
							if !strings.HasPrefix(m[1], "|c") {
								parts = append(parts, m[1])
							}
						}
					}
					if len(parts) > 0 {
						currentName = strings.Join(parts, " ")
					}
				}

				// Parse LootPage
				if match := lootPageRegex.FindStringSubmatch(line); len(match) > 1 {
					key := match[1]
					name := currentName
					if name == "" {
						name = key
					}
					items = append(items, SpecialMenuItem{Key: key, Name: name})
					currentName = ""
				}
			}

			if len(items) > 0 {
				specialMenus[menuKey] = items
			}
		}
		return nil
	})

	return specialMenus, err
}
