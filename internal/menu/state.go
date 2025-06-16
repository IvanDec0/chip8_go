package menu

// StateManager handles transitions between emulator states
type StateManager struct {
	currentState  EmulatorState
	previousState EmulatorState
	menuData      *MenuData
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{
		currentState:  StateMenu, // Start in menu by default
		previousState: StateMenu,
		menuData: &MenuData{
			SelectedROM:   "",
			ShouldLoadROM: false,
			ShouldExit:    false,
		},
	}
}

// GetCurrentState returns the current emulator state
func (sm *StateManager) GetCurrentState() EmulatorState {
	return sm.currentState
}

// GetPreviousState returns the previous emulator state
func (sm *StateManager) GetPreviousState() EmulatorState {
	return sm.previousState
}

// TransitionTo transitions to a new state
func (sm *StateManager) TransitionTo(newState EmulatorState) {
	sm.previousState = sm.currentState
	sm.currentState = newState
}

// GetMenuData returns the current menu data
func (sm *StateManager) GetMenuData() *MenuData {
	return sm.menuData
}

// SetSelectedROM sets the selected ROM path
func (sm *StateManager) SetSelectedROM(romPath string) {
	sm.menuData.SelectedROM = romPath
}

// SetShouldLoadROM sets whether a ROM should be loaded
func (sm *StateManager) SetShouldLoadROM(shouldLoad bool) {
	sm.menuData.ShouldLoadROM = shouldLoad
}

// SetShouldExit sets whether the application should exit
func (sm *StateManager) SetShouldExit(shouldExit bool) {
	sm.menuData.ShouldExit = shouldExit
}

// ShouldLoadROM returns whether a ROM should be loaded
func (sm *StateManager) ShouldLoadROM() bool {
	return sm.menuData.ShouldLoadROM
}

// ShouldExit returns whether the application should exit
func (sm *StateManager) ShouldExit() bool {
	return sm.menuData.ShouldExit
}

// GetSelectedROM returns the selected ROM path
func (sm *StateManager) GetSelectedROM() string {
	return sm.menuData.SelectedROM
}

// Reset resets the menu data
func (sm *StateManager) Reset() {
	sm.menuData.SelectedROM = ""
	sm.menuData.ShouldLoadROM = false
	sm.menuData.ShouldExit = false
}
