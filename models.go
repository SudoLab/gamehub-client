package gamehub

import (
	"encoding/json"
	"time"
)

// User represents a GameHub user
type User struct {
	ID               int64     `json:"id"`
	TelegramID       int64     `json:"telegram_id"`
	Username         string    `json:"username"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	DisplayName      string    `json:"display_name"`
	ProfilePicURL    string    `json:"profile_pic_url,omitempty"`
	GlobalCoins      int       `json:"global_coins"`
	TotalGamesPlayed int       `json:"total_games_played"`
	AverageScore     float64   `json:"average_score"`
	ReferralCode     string    `json:"referral_code"`
	ReferredBy       *int64    `json:"referred_by"`
	WalletAddress    string    `json:"wallet_address,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Game represents a registered game
type Game struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
	IconURL     string    `json:"icon_url"`
	BannerURL   string    `json:"banner_url"`
	MinVersion  string    `json:"min_version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserGameStats represents user's stats for a specific game
type UserGameStats struct {
	GameID      string    `json:"game_id"`
	GameName    string    `json:"game_name"`
	TotalScore  int64     `json:"total_score"`
	GamesPlayed int       `json:"games_played"`
	BestScore   int64     `json:"best_score"`
	LastPlayed  time.Time `json:"last_played"`
}

// Transaction represents a coin transaction
type Transaction struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	GameID          string    `json:"game_id,omitempty"`
	Amount          int       `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	Description     string    `json:"description"`
	ReferenceID     string    `json:"reference_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// RankingEntry represents a leaderboard entry
type RankingEntry struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Score       int64  `json:"score"`
	GamesPlayed int    `json:"games_played"`
	Rank        int    `json:"rank"`
}

// UserRanking represents a user's ranking information
type UserRanking struct {
	UserID       int64   `json:"user_id"`
	GlobalRank   int     `json:"global_rank"`
	TotalScore   int64   `json:"total_score"`
	GamesPlayed  int     `json:"games_played"`
	AverageScore float64 `json:"average_score"`
}

// Referral represents a referral relationship
type Referral struct {
	ID           int64      `json:"id"`
	ReferrerID   int64      `json:"referrer_id"`
	ReferredID   int64      `json:"referred_id"`
	ReferredUser string     `json:"referred_user"`
	RewardAmount int        `json:"reward_amount"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
}

// ReferralStats represents referral statistics
type ReferralStats struct {
	TotalReferrals     int `json:"total_referrals"`
	CompletedReferrals int `json:"completed_referrals"`
	TotalRewards       int `json:"total_rewards"`
	PendingRewards     int `json:"pending_rewards"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
	Message string          `json:"message,omitempty"`
}
