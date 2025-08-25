// encore:service
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"encore.dev/beta/auth"
	"encore.dev/storage/sqldb"
)

// Agent represents a Valorant agent
type Agent struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Description string   `json:"description"`
	Abilities   []string `json:"abilities"`
	ImageURL    string   `json:"imageUrl"`
}

// Weapon represents a Valorant weapon
type Weapon struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"` // Primary, Sidearm
	Cost     int    `json:"cost"`
	Damage   int    `json:"damage"`
	Accuracy int    `json:"accuracy"`
	ImageURL string `json:"imageUrl"`
}

// Loadout represents a user's selected loadout
type Loadout struct {
	ID      int       `json:"id"`
	UserID  string    `json:"userId"`
	Agent   string    `json:"agent"`
	Primary string    `json:"primary"`
	Sidearm string    `json:"sidearm"`
	Created time.Time `json:"created"`
}

// Database handle (Postgres)
var db = sqldb.Named("valorant")

// Initialize database tables (migration)
var _ = sqldb.Migration("001_create_tables", `
	CREATE TABLE IF NOT EXISTS loadouts (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL,
		agent TEXT NOT NULL,
		primary_weapon TEXT,
		sidearm TEXT,
		created TIMESTAMP DEFAULT NOW()
	);
`)

// --- Static data (demo) ---
var agents = []Agent{
	{
		ID: "jett", Name: "Jett", Role: "Duelist",
		Description: "Jett's agile and evasive fighting style lets her take risks no one else can.",
		Abilities:   []string{"Updraft", "Tailwind", "Cloudburst", "Blade Storm"},
		ImageURL:    "https://images.contentstack.io/v3/assets/bltb6530b271fddd0b1/blt5ebf40a2dfaffb4e/5f21297f5f0cb0629a5bfcb9/V_AGENTS_587x900_Jett.png",
	},
	{
		ID: "sova", Name: "Sova", Role: "Initiator",
		Description: "Sova tracks, finds, and eliminates enemies with ruthless efficiency.",
		Abilities:   []string{"Owl Drone", "Shock Bolt", "Recon Bolt", "Hunter's Fury"},
		ImageURL:    "https://images.contentstack.io/v3/assets/bltb6530b271fddd0b1/blt181ad63adc9976a4/5f2129b2e0999b628bc8eb4e/V_AGENTS_587x900_Sova.png",
	},
	{
		ID: "sage", Name: "Sage", Role: "Sentinel",
		Description: "Sage creates safety for herself and her team wherever she goes.",
		Abilities:   []string{"Barrier Orb", "Slow Orb", "Healing Orb", "Resurrection"},
		ImageURL:    "https://images.contentstack.io/v3/assets/bltb6530b271fddd0b1/blt2a1c7b18aa5b1a6b/5f21297f078a8b626859f4a8/V_AGENTS_587x900_Sage.png",
	},
	{
		ID: "omen", Name: "Omen", Role: "Controller",
		Description: "Omen hunts in the shadows. He renders enemies blind, teleports across the field.",
		Abilities:   []string{"Shrouded Step", "Paranoia", "Dark Cover", "From the Shadows"},
		ImageURL:    "https://images.contentstack.io/v3/assets/bltb6530b271fddd0b1/blt94dd043bce7fc9f2/5f21297f2ef66062fb6aa96c/V_AGENTS_587x900_Omen.png",
	},
}

var weapons = []Weapon{
	{ID: "classic", Name: "Classic", Type: "Sidearm", Cost: 0, Damage: 78, Accuracy: 85, ImageURL: ""},
	{ID: "sheriff", Name: "Sheriff", Type: "Sidearm", Cost: 800, Damage: 159, Accuracy: 79, ImageURL: ""},
	{ID: "spectre", Name: "Spectre", Type: "Primary", Cost: 1600, Damage: 78, Accuracy: 74, ImageURL: ""},
	{ID: "vandal", Name: "Vandal", Type: "Primary", Cost: 2900, Damage: 160, Accuracy: 73, ImageURL: ""},
	{ID: "phantom", Name: "Phantom", Type: "Primary", Cost: 2900, Damage: 156, Accuracy: 79, ImageURL: ""},
	{ID: "operator", Name: "Operator", Type: "Primary", Cost: 4700, Damage: 255, Accuracy: 76, ImageURL: ""},
}

// --- Auth handler (dev) ---
// encore:authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, *auth.Data, error) {
	// Simple dev auth: any Bearer token starting with "dev-" is accepted
	if strings.HasPrefix(token, "dev-") && len(token) > 4 {
		uid := auth.UID(token) // e.g. dev-alice
		return uid, &auth.Data{User: map[string]any{"role": "dev"}}, nil
	}
	return "", nil, auth.ErrUnauthenticated
}

// --- Public APIs ---

// encore:api public method=GET path=/agents
func GetAgents(ctx context.Context, params *GetAgentsParams) (*GetAgentsResponse, error) {
	var filtered []Agent
	for _, a := range agents {
		if params.Role != "" && a.Role != params.Role {
			continue
		}
		if params.Search != "" {
			q := strings.ToLower(params.Search)
			if !strings.Contains(strings.ToLower(a.Name), q) &&
				!strings.Contains(strings.ToLower(a.Description), q) {
				continue
			}
		}
		filtered = append(filtered, a)
	}
	return &GetAgentsResponse{Agents: filtered, Total: len(filtered)}, nil
}

type GetAgentsParams struct {
	Role   string `query:"role"`
	Search string `query:"search"`
}

type GetAgentsResponse struct {
	Agents []Agent `json:"agents"`
	Total  int     `json:"total"`
}

// encore:api public method=GET path=/weapons
func GetWeapons(ctx context.Context, params *GetWeaponsParams) (*GetWeaponsResponse, error) {
	var filtered []Weapon
	for _, w := range weapons {
		if params.Type != "" && w.Type != params.Type {
			continue
		}
		if params.MaxCost > 0 && w.Cost > params.MaxCost {
			continue
		}
		if params.Search != "" {
			q := strings.ToLower(params.Search)
			if !strings.Contains(strings.ToLower(w.Name), q) {
				continue
			}
		}
		filtered = append(filtered, w)
	}
	return &GetWeaponsResponse{Weapons: filtered, Total: len(filtered)}, nil
}

type GetWeaponsParams struct {
	Type    string `query:"type"`
	MaxCost int    `query:"maxCost"`
	Search  string `query:"search"`
}

type GetWeaponsResponse struct {
	Weapons []Weapon `json:"weapons"`
	Total   int      `json:"total"`
}

// encore:api public method=GET path=/health
func HealthCheck(ctx context.Context) (*HealthResponse, error) {
	return &HealthResponse{
		Status:    "healthy",
		Message:   "Valorant API is running",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}, nil
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// encore:api public method=GET path=/stats
func GetStats(ctx context.Context) (*StatsResponse, error) {
	var total int
	if err := db.QueryRow(ctx, "SELECT COUNT(*) FROM loadouts").Scan(&total); err != nil {
		total = 0
	}
	return &StatsResponse{
		TotalAgents:   len(agents),
		TotalWeapons:  len(weapons),
		TotalLoadouts: total,
		PopularAgent:  "Jett", // demo
	}, nil
}

type StatsResponse struct {
	TotalAgents   int    `json:"totalAgents"`
	TotalWeapons  int    `json:"totalWeapons"`
	TotalLoadouts int    `json:"totalLoadouts"`
	PopularAgent  string `json:"popularAgent"`
}

// --- Auth APIs ---

// encore:api auth method=POST path=/loadouts
func CreateLoadout(ctx context.Context, req *CreateLoadoutRequest) (*CreateLoadoutResponse, error) {
	userID, _ := auth.UserID()
	var id int
	err := db.QueryRow(ctx, `
		INSERT INTO loadouts (user_id, agent, primary_weapon, sidearm)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, string(userID), req.Agent, req.Primary, req.Sidearm).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create loadout: %v", err)
	}
	return &CreateLoadoutResponse{ID: id, Message: "Loadout saved successfully"}, nil
}

type CreateLoadoutRequest struct {
	Agent   string `json:"agent"`
	Primary string `json:"primary"`
	Sidearm string `json:"sidearm"`
}

type CreateLoadoutResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// encore:api auth method=GET path=/loadouts
func GetUserLoadouts(ctx context.Context) (*GetLoadoutsResponse, error) {
	userID, _ := auth.UserID()
	rows, err := db.Query(ctx, `
		SELECT id, agent, primary_weapon, sidearm, created
		FROM loadouts
		WHERE user_id = $1
		ORDER BY created DESC
	`, string(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get loadouts: %v", err)
	}
	defer rows.Close()

	var out []Loadout
	for rows.Next() {
		var l Loadout
		var primary, sidearm *string
		if err := rows.Scan(&l.ID, &l.Agent, &primary, &sidearm, &l.Created); err != nil {
			return nil, err
		}
		l.UserID = string(userID)
		if primary != nil {
			l.Primary = *primary
		}
		if sidearm != nil {
			l.Sidearm = *sidearm
		}
		out = append(out, l)
	}
	return &GetLoadoutsResponse{Loadouts: out, Total: len(out)}, nil
}

type GetLoadoutsResponse struct {
	Loadouts []Loadout `json:"loadouts"`
	Total    int       `json:"total"`
}
