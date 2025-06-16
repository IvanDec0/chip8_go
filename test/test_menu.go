package main

import (
	"chip8/internal/browser"
	"chip8/internal/menu"
	"fmt"
)

// Simple test to verify Phase 1 implementation
func main() {
	fmt.Println("=== CHIP-8 ROM Browser - Phase 1 Test ===")

	// Test 1: Browser functionality
	fmt.Println("\n1. Testing Browser Implementation...")
	br := browser.NewFileBrowser("./roms")

	roms, err := br.ScanDirectory("./roms")
	if err != nil {
		fmt.Printf("   ❌ Error scanning directory: %v\n", err)
	} else {
		fmt.Printf("   ✅ Found %d items in roms directory\n", len(roms))
		for i, rom := range roms {
			if i < 5 { // Show first 5 items
				icon := "🎮"
				if rom.IsDirectory {
					icon = "📁"
				}
				fmt.Printf("      %s %s\n", icon, rom.Name)
			}
		}
		if len(roms) > 5 {
			fmt.Printf("      ... and %d more items\n", len(roms)-5)
		}
	}

	// Test 2: State Management
	fmt.Println("\n2. Testing State Management...")
	sm := menu.NewStateManager()
	fmt.Printf("   ✅ Initial state: %s\n", sm.GetCurrentState())

	sm.TransitionTo(menu.StateEmulator)
	fmt.Printf("   ✅ Transitioned to: %s\n", sm.GetCurrentState())

	sm.SetSelectedROM("test.ch8")
	fmt.Printf("   ✅ Selected ROM: %s\n", sm.GetSelectedROM())

	// Test 3: Menu Management
	fmt.Println("\n3. Testing Menu Management...")
	manager := menu.NewManager()
	manager.SetBrowser(br)

	items := manager.GetMenuItems()
	fmt.Printf("   ✅ Menu has %d items:\n", len(items))
	for _, item := range items {
		fmt.Printf("      - %s (%s)\n", item.Text, item.Action)
	}

	// Test 4: Browser Integration
	fmt.Println("\n4. Testing Browser Integration...")
	err = manager.LoadROMBrowser()
	if err != nil {
		fmt.Printf("   ❌ Error loading ROM browser: %v\n", err)
	} else {
		browserItems := manager.GetMenuItems()
		fmt.Printf("   ✅ ROM browser loaded with %d items\n", len(browserItems))
	}

	// Test 5: Configuration Types
	fmt.Println("\n5. Testing Configuration Types...")
	theme := menu.DefaultMenuTheme()
	fmt.Printf("   ✅ Default theme - Text: RGB(%d,%d,%d)\n",
		theme.TextColor.R, theme.TextColor.G, theme.TextColor.B)

	fmt.Println("\n🎉 Phase 1 Implementation Test Complete!")
	fmt.Println("\n📋 Ready for Phase 2: Text Rendering & Menu UI")
}
