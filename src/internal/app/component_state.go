// src/internal/app/component_state.go
// Component-specific state management

package app

import (
	"nextui-themes/internal/logging"
)

// ComponentType represents the type of theme component
type ComponentType int

const (
	ComponentTypeFullTheme ComponentType = iota + 1
	ComponentTypeAccent
	ComponentTypeLED
	ComponentTypeWallpaper
	ComponentTypeIcon
	ComponentTypeFont
)

// Component state variables
type componentState struct {
	// Import state
	ImportComponentType       ComponentType
	SelectedImportItem        string
	ImportAllComponents       bool
	SelectedImportComponents  map[ComponentType]bool

	// Export state
	ExportComponentType       ComponentType
	ExportName                string
	ExportAllComponents       bool
	SelectedExportComponents  map[ComponentType]bool

	// Convert/deconstruct state
	SelectedConvertTheme      string
	ConvertAllComponents      bool
	SelectedConvertComponents map[ComponentType]bool
}

// Global component state
var compState = componentState{
	SelectedImportComponents:  make(map[ComponentType]bool),
	SelectedExportComponents:  make(map[ComponentType]bool),
	SelectedConvertComponents: make(map[ComponentType]bool),
}

// Import component state functions

// SetImportComponentType sets the type of component to import
func SetImportComponentType(componentType ComponentType) {
	logging.LogDebug("Setting import component type: %d", componentType)
	compState.ImportComponentType = componentType
}

// GetImportComponentType returns the current import component type
func GetImportComponentType() ComponentType {
	return compState.ImportComponentType
}

// SetSelectedImportItem sets the selected item to import
func SetSelectedImportItem(item string) {
	logging.LogDebug("Setting selected import item: %s", item)
	compState.SelectedImportItem = item
}

// GetSelectedImportItem returns the currently selected import item
func GetSelectedImportItem() string {
	return compState.SelectedImportItem
}

// SetImportAllComponents sets whether to import all components from a theme
func SetImportAllComponents(importAll bool) {
	logging.LogDebug("Setting import all components: %v", importAll)
	compState.ImportAllComponents = importAll
}

// GetImportAllComponents returns whether all components should be imported
func GetImportAllComponents() bool {
	return compState.ImportAllComponents
}

// AddImportComponent adds a component type to the import selection
func AddImportComponent(componentType ComponentType) {
	compState.SelectedImportComponents[componentType] = true
}

// RemoveImportComponent removes a component type from the import selection
func RemoveImportComponent(componentType ComponentType) {
	delete(compState.SelectedImportComponents, componentType)
}

// ToggleImportComponent toggles a component type's selection state
func ToggleImportComponent(componentType ComponentType) {
	if compState.SelectedImportComponents[componentType] {
		delete(compState.SelectedImportComponents, componentType)
		logging.LogDebug("Removed component %d from import selection", componentType)
	} else {
		compState.SelectedImportComponents[componentType] = true
		logging.LogDebug("Added component %d to import selection", componentType)
	}
}

// GetSelectedImportComponents returns the currently selected import components
func GetSelectedImportComponents() map[ComponentType]bool {
	return compState.SelectedImportComponents
}

// ClearSelectedImportComponents clears all selected import components
func ClearSelectedImportComponents() {
	compState.SelectedImportComponents = make(map[ComponentType]bool)
}

// Export component state functions

// SetExportComponentType sets the type of component to export
func SetExportComponentType(componentType ComponentType) {
	logging.LogDebug("Setting export component type: %d", componentType)
	compState.ExportComponentType = componentType
}

// GetExportComponentType returns the current export component type
func GetExportComponentType() ComponentType {
	return compState.ExportComponentType
}

// SetExportName sets the name for the exported component/theme
func SetExportName(name string) {
	logging.LogDebug("Setting export name: %s", name)
	compState.ExportName = name
}

// GetExportName returns the current export name
func GetExportName() string {
	return compState.ExportName
}

// SetExportAllComponents sets whether to export all components for a theme
func SetExportAllComponents(exportAll bool) {
	logging.LogDebug("Setting export all components: %v", exportAll)
	compState.ExportAllComponents = exportAll
}

// GetExportAllComponents returns whether all components should be exported
func GetExportAllComponents() bool {
	return compState.ExportAllComponents
}

// AddExportComponent adds a component type to the export selection
func AddExportComponent(componentType ComponentType) {
	compState.SelectedExportComponents[componentType] = true
}

// RemoveExportComponent removes a component type from the export selection
func RemoveExportComponent(componentType ComponentType) {
	delete(compState.SelectedExportComponents, componentType)
}

// ToggleExportComponent toggles a component type's selection state
func ToggleExportComponent(componentType ComponentType) {
	if compState.SelectedExportComponents[componentType] {
		delete(compState.SelectedExportComponents, componentType)
		logging.LogDebug("Removed component %d from export selection", componentType)
	} else {
		compState.SelectedExportComponents[componentType] = true
		logging.LogDebug("Added component %d to export selection", componentType)
	}
}

// GetSelectedExportComponents returns the currently selected export components
func GetSelectedExportComponents() map[ComponentType]bool {
	return compState.SelectedExportComponents
}

// ClearSelectedExportComponents clears all selected export components
func ClearSelectedExportComponents() {
	compState.SelectedExportComponents = make(map[ComponentType]bool)
}

// Convert/deconstruct theme state functions

// SetSelectedConvertTheme sets the theme selected for conversion/deconstruction
func SetSelectedConvertTheme(theme string) {
	logging.LogDebug("Setting selected convert theme: %s", theme)
	compState.SelectedConvertTheme = theme
}

// GetSelectedConvertTheme returns the currently selected theme for conversion
func GetSelectedConvertTheme() string {
	return compState.SelectedConvertTheme
}

// SetConvertAllComponents sets whether to convert all components from a theme
func SetConvertAllComponents(convertAll bool) {
	logging.LogDebug("Setting convert all components: %v", convertAll)
	compState.ConvertAllComponents = convertAll
}

// GetConvertAllComponents returns whether all components should be converted
func GetConvertAllComponents() bool {
	return compState.ConvertAllComponents
}

// AddConvertComponent adds a component type to the conversion selection
func AddConvertComponent(componentType ComponentType) {
	compState.SelectedConvertComponents[componentType] = true
}

// RemoveConvertComponent removes a component type from the conversion selection
func RemoveConvertComponent(componentType ComponentType) {
	delete(compState.SelectedConvertComponents, componentType)
}

// ToggleConvertComponent toggles a component type's selection state
func ToggleConvertComponent(componentType ComponentType) {
	if compState.SelectedConvertComponents[componentType] {
		delete(compState.SelectedConvertComponents, componentType)
		logging.LogDebug("Removed component %d from convert selection", componentType)
	} else {
		compState.SelectedConvertComponents[componentType] = true
		logging.LogDebug("Added component %d to convert selection", componentType)
	}
}

// GetSelectedConvertComponents returns the currently selected conversion components
func GetSelectedConvertComponents() []ComponentType {
	var components []ComponentType
	for comp := range compState.SelectedConvertComponents {
		components = append(components, comp)
	}
	return components
}

// ClearSelectedConvertComponents clears all selected conversion components
func ClearSelectedConvertComponents() {
	compState.SelectedConvertComponents = make(map[ComponentType]bool)
}

// ResetComponentState resets all component state to defaults
func ResetComponentState() {
	compState.ImportComponentType = ComponentTypeFullTheme
	compState.SelectedImportItem = ""
	compState.ImportAllComponents = true
	compState.SelectedImportComponents = make(map[ComponentType]bool)

	compState.ExportComponentType = ComponentTypeFullTheme
	compState.ExportName = ""
	compState.ExportAllComponents = true
	compState.SelectedExportComponents = make(map[ComponentType]bool)

	compState.SelectedConvertTheme = ""
	compState.ConvertAllComponents = true
	compState.SelectedConvertComponents = make(map[ComponentType]bool)
}