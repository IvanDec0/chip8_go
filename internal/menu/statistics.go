package menu

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// PlaySession represents a single ROM play session
type PlaySession struct {
	ROMPath   string        `json:"rom_path"`
	ROMName   string        `json:"rom_name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
}

// ROMStats represents statistics for a specific ROM
type ROMStats struct {
	Path           string        `json:"path"`
	Name           string        `json:"name"`
	PlayCount      int           `json:"play_count"`
	TotalPlayTime  time.Duration `json:"total_play_time"`
	AverageSession time.Duration `json:"average_session"`
	LastPlayed     time.Time     `json:"last_played"`
	FirstPlayed    time.Time     `json:"first_played"`
}

// SystemStats represents overall system usage statistics
type SystemStats struct {
	TotalSessions      int           `json:"total_sessions"`
	TotalPlayTime      time.Duration `json:"total_play_time"`
	AverageSessionTime time.Duration `json:"average_session_time"`
	UniqueROMsPlayed   int           `json:"unique_roms_played"`
	MostPlayedROM      string        `json:"most_played_rom"`
	LongestSession     time.Duration `json:"longest_session"`
	FirstLaunch        time.Time     `json:"first_launch"`
	LastActivity       time.Time     `json:"last_activity"`
}

// StatisticsManager manages usage statistics and analytics
type StatisticsManager struct {
	romStats        map[string]*ROMStats
	sessions        []PlaySession
	dataPath        string
	maxSessions     int
	currentSession  *PlaySession
	autoSaveEnabled bool
}

// NewStatisticsManager creates a new statistics manager
func NewStatisticsManager(dataPath string) *StatisticsManager {
	manager := &StatisticsManager{
		romStats:        make(map[string]*ROMStats),
		sessions:        make([]PlaySession, 0),
		dataPath:        dataPath,
		maxSessions:     1000, // Keep last 1000 sessions
		autoSaveEnabled: true,
	}

	// Ensure data directory exists
	if err := os.MkdirAll(dataPath, 0755); err == nil {
		manager.loadStatistics()
	}

	return manager
}

// StartSession begins a new play session for a ROM
func (sm *StatisticsManager) StartSession(romPath, romName string) {
	sm.EndCurrentSession() // End any existing session

	sm.currentSession = &PlaySession{
		ROMPath:   romPath,
		ROMName:   romName,
		StartTime: time.Now(),
	}
}

// EndCurrentSession ends the current play session
func (sm *StatisticsManager) EndCurrentSession() {
	if sm.currentSession == nil {
		return
	}

	sm.currentSession.EndTime = time.Now()
	sm.currentSession.Duration = sm.currentSession.EndTime.Sub(sm.currentSession.StartTime)

	// Only record sessions longer than 5 seconds
	if sm.currentSession.Duration >= 5*time.Second {
		sm.recordSession(*sm.currentSession)
	}

	sm.currentSession = nil

	if sm.autoSaveEnabled {
		sm.saveStatistics()
	}
}

// recordSession records a completed session and updates statistics
func (sm *StatisticsManager) recordSession(session PlaySession) {
	// Add to sessions list
	sm.sessions = append(sm.sessions, session)

	// Limit session history
	if len(sm.sessions) > sm.maxSessions {
		sm.sessions = sm.sessions[len(sm.sessions)-sm.maxSessions:]
	}

	// Update ROM statistics
	romStats, exists := sm.romStats[session.ROMPath]
	if !exists {
		romStats = &ROMStats{
			Path:        session.ROMPath,
			Name:        session.ROMName,
			FirstPlayed: session.StartTime,
		}
		sm.romStats[session.ROMPath] = romStats
	}

	romStats.PlayCount++
	romStats.TotalPlayTime += session.Duration
	romStats.LastPlayed = session.StartTime

	if romStats.PlayCount > 0 {
		romStats.AverageSession = romStats.TotalPlayTime / time.Duration(romStats.PlayCount)
	}
}

// GetROMStats returns statistics for a specific ROM
func (sm *StatisticsManager) GetROMStats(romPath string) (*ROMStats, bool) {
	stats, exists := sm.romStats[romPath]
	return stats, exists
}

// GetTopROMs returns the most played ROMs
func (sm *StatisticsManager) GetTopROMs(limit int) []*ROMStats {
	roms := make([]*ROMStats, 0, len(sm.romStats))
	for _, stats := range sm.romStats {
		roms = append(roms, stats)
	}

	// Sort by play count (descending)
	sort.Slice(roms, func(i, j int) bool {
		if roms[i].PlayCount == roms[j].PlayCount {
			// If play count is equal, sort by total play time
			return roms[i].TotalPlayTime > roms[j].TotalPlayTime
		}
		return roms[i].PlayCount > roms[j].PlayCount
	})

	if limit > 0 && len(roms) > limit {
		roms = roms[:limit]
	}

	return roms
}

// GetRecentSessions returns recent play sessions
func (sm *StatisticsManager) GetRecentSessions(limit int) []PlaySession {
	sessions := make([]PlaySession, len(sm.sessions))
	copy(sessions, sm.sessions)

	// Sort by start time (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.After(sessions[j].StartTime)
	})

	if limit > 0 && len(sessions) > limit {
		sessions = sessions[:limit]
	}

	return sessions
}

// GetSystemStats returns overall system usage statistics
func (sm *StatisticsManager) GetSystemStats() SystemStats {
	stats := SystemStats{
		TotalSessions:    len(sm.sessions),
		UniqueROMsPlayed: len(sm.romStats),
	}

	if len(sm.sessions) == 0 {
		return stats
	}

	var totalDuration time.Duration
	var longestSession time.Duration
	var firstLaunch, lastActivity time.Time

	for i, session := range sm.sessions {
		totalDuration += session.Duration

		if session.Duration > longestSession {
			longestSession = session.Duration
		}

		if i == 0 || session.StartTime.Before(firstLaunch) {
			firstLaunch = session.StartTime
		}

		if i == 0 || session.StartTime.After(lastActivity) {
			lastActivity = session.StartTime
		}
	}

	stats.TotalPlayTime = totalDuration
	stats.LongestSession = longestSession
	stats.FirstLaunch = firstLaunch
	stats.LastActivity = lastActivity

	if stats.TotalSessions > 0 {
		stats.AverageSessionTime = totalDuration / time.Duration(stats.TotalSessions)
	}

	// Find most played ROM
	topROMs := sm.GetTopROMs(1)
	if len(topROMs) > 0 {
		stats.MostPlayedROM = topROMs[0].Name
	}

	return stats
}

// GetSessionsInTimeRange returns sessions within a specific time range
func (sm *StatisticsManager) GetSessionsInTimeRange(start, end time.Time) []PlaySession {
	var sessions []PlaySession

	for _, session := range sm.sessions {
		if session.StartTime.After(start) && session.StartTime.Before(end) {
			sessions = append(sessions, session)
		}
	}

	// Sort by start time (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.After(sessions[j].StartTime)
	})

	return sessions
}

// GetDailyStats returns statistics grouped by day
func (sm *StatisticsManager) GetDailyStats(days int) map[string]interface{} {
	dailyStats := make(map[string]interface{})
	now := time.Now()

	for i := 0; i < days; i++ {
		day := now.AddDate(0, 0, -i)
		dayKey := day.Format("2006-01-02")

		startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		sessions := sm.GetSessionsInTimeRange(startOfDay, endOfDay)

		var totalDuration time.Duration
		uniqueROMs := make(map[string]bool)

		for _, session := range sessions {
			totalDuration += session.Duration
			uniqueROMs[session.ROMPath] = true
		}

		dailyStats[dayKey] = map[string]interface{}{
			"date":            dayKey,
			"sessions":        len(sessions),
			"total_play_time": totalDuration,
			"unique_roms":     len(uniqueROMs),
			"average_session": func() time.Duration {
				if len(sessions) > 0 {
					return totalDuration / time.Duration(len(sessions))
				}
				return 0
			}(),
		}
	}

	return dailyStats
}

// ExportStatistics exports statistics to JSON
func (sm *StatisticsManager) ExportStatistics() (map[string]interface{}, error) {
	export := map[string]interface{}{
		"system_stats":    sm.GetSystemStats(),
		"rom_stats":       sm.romStats,
		"recent_sessions": sm.GetRecentSessions(100),
		"top_roms":        sm.GetTopROMs(20),
		"daily_stats":     sm.GetDailyStats(30),
		"export_time":     time.Now(),
	}

	return export, nil
}

// ClearStatistics clears all statistics (with confirmation)
func (sm *StatisticsManager) ClearStatistics() error {
	sm.romStats = make(map[string]*ROMStats)
	sm.sessions = make([]PlaySession, 0)
	sm.currentSession = nil

	// Remove statistics files
	statsFile := filepath.Join(sm.dataPath, "statistics.json")
	sessionsFile := filepath.Join(sm.dataPath, "sessions.json")

	os.Remove(statsFile)
	os.Remove(sessionsFile)

	return nil
}

// SetAutoSave enables or disables automatic saving
func (sm *StatisticsManager) SetAutoSave(enabled bool) {
	sm.autoSaveEnabled = enabled
}

// SaveStatistics manually saves statistics to disk
func (sm *StatisticsManager) SaveStatistics() error {
	return sm.saveStatistics()
}

// saveStatistics saves statistics to disk
func (sm *StatisticsManager) saveStatistics() error {
	// Save ROM statistics
	statsFile := filepath.Join(sm.dataPath, "statistics.json")
	statsData, err := json.MarshalIndent(sm.romStats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ROM stats: %v", err)
	}

	if err := os.WriteFile(statsFile, statsData, 0644); err != nil {
		return fmt.Errorf("failed to save ROM stats: %v", err)
	}

	// Save sessions
	sessionsFile := filepath.Join(sm.dataPath, "sessions.json")
	sessionsData, err := json.MarshalIndent(sm.sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sessions: %v", err)
	}

	if err := os.WriteFile(sessionsFile, sessionsData, 0644); err != nil {
		return fmt.Errorf("failed to save sessions: %v", err)
	}

	return nil
}

// loadStatistics loads statistics from disk
func (sm *StatisticsManager) loadStatistics() {
	// Load ROM statistics
	statsFile := filepath.Join(sm.dataPath, "statistics.json")
	if data, err := os.ReadFile(statsFile); err == nil {
		json.Unmarshal(data, &sm.romStats)
	}

	// Load sessions
	sessionsFile := filepath.Join(sm.dataPath, "sessions.json")
	if data, err := os.ReadFile(sessionsFile); err == nil {
		json.Unmarshal(data, &sm.sessions)
	}
}
