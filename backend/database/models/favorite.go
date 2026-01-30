// Package models contains all database entity definitions
package models

// FavoriteItem represents a user's favorite item
type FavoriteItem struct {
	ID        int    `json:"id"`
	ItemEntry int    `json:"itemEntry"`
	Category  string `json:"category"`
	AddedAt   string `json:"addedAt"`
	// Joined data from item_template
	ItemName    string `json:"itemName,omitempty"`
	ItemQuality int    `json:"itemQuality,omitempty"`
	IconPath    string `json:"iconPath,omitempty"`
	ItemLevel   int    `json:"itemLevel,omitempty"`
	// Status: 0=None, 1=Obtained, 2=Abandoned
	Status int `json:"status"`
}

// FavoriteCategory represents a grouping for favorites
type FavoriteCategory struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// FavoriteResult represents the result of adding/removing a favorite
type FavoriteResult struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
