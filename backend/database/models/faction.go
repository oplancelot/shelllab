package models

// Faction represents a WoW reputation faction
type Faction struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Side        int    `json:"side"` // 1=Alliance, 2=Horde, 3=Both
	CategoryId  int    `json:"categoryId"`
}

// FactionEntry represents a faction for JSON import
type FactionEntry struct {
	FactionID   int    `json:"factionID"`
	Name        string `json:"name_loc0"`
	Description string `json:"description1_loc0"`
	Side        int    `json:"side"`
	Team        int    `json:"team"`
}

// FactionDetail represents detailed faction info for detail view
type FactionDetail struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Side        int              `json:"side"` // 1=Alliance, 2=Horde, 3=Both
	SideName    string           `json:"sideName"`
	CategoryId  int              `json:"categoryId"`
	Creatures   []*Creature      `json:"creatures,omitempty"`
	Quests      []*QuestRelation `json:"quests,omitempty"`
}
