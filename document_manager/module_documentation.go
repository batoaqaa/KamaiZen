package document_manager

import (
	"KamaiZen/logger"
	"errors"
)

/**
 * @brief Holds the documentation for all modules.
 * It maps module names to their corresponding ModuleDocs structs.
 */
type moduleDocumentationMap struct {
	ModuleDocs map[string]ModuleDocs
}

/**
 * @brief Retrieves the documentation for a specific module from the module documentation map.
 * It searches for the module by its name and returns the corresponding ModuleDocs struct and a boolean
 * indicating whether the module was found.
 *
 * @param moduleName The name of the module to retrieve documentation for.
 * @return A ModuleDocs struct containing the documentation for the specified module.
 * @return A boolean indicating whether the module was found (true) or not (false).
 */
func (m *moduleDocumentationMap) GetModuleDocs(moduleName string) (ModuleDocs, bool) {
	for key, value := range m.ModuleDocs {
		if key == moduleName {
			return value, true
		}
	}
	return ModuleDocs{}, false
}

/**
 * @brief Adds a module's documentation to the module documentation map.
 * If the module documentation already exists and overwrite is set to false, it returns an error.
 * If overwrite is set to true, it overwrites the existing documentation.
 *
 * @param moduleName The name of the module to add documentation for.
 * @param moduleDocs The ModuleDocs struct containing the documentation to be added.
 * @param overwrite A boolean indicating whether to overwrite existing documentation if it already exists.
 * @return An error if the module documentation already exists and overwrite is false, otherwise nil.
 */
func (m *moduleDocumentationMap) AddModuleDocs(moduleName string, moduleDocs ModuleDocs, overwrite bool) error {
	if _, exists := m.ModuleDocs[moduleName]; exists && !overwrite {
		return errors.New("Module already exists")
	}
	m.ModuleDocs[moduleName] = moduleDocs
	return nil
}

type ModuleDocs struct {
	Functions map[string]FunctionDocumentationMap
}

/**
 * @brief Adds a function's documentation to the module's function documentation map.
 * If the function documentation already exists and overwrite is set to false, it returns an error.
 * If overwrite is set to true, it logs a message and overwrites the existing documentation.
 *
 * @param moduleName The name of the module to add the function documentation to.
 * @param functionDocs The FunctionDocumentationMap struct containing the function documentation to be added.
 * @param overwrite A boolean indicating whether to overwrite existing documentation if it already exists.
 * @return An error if the function documentation already exists and overwrite is false, otherwise nil.
 */
func (m *ModuleDocs) AddFunctionDoc(moduleName string, functionDocs FunctionDocumentationMap, overwrite bool) error {
	if _, exists := m.Functions[moduleName]; exists {
		if !overwrite {
			return errors.New("Function Documentation already exists")
		} else {
			logger.Info("Overwrting function documentation for module: ", moduleName)
		}
	}
	m.Functions[moduleName] = functionDocs
	return nil
}

/**
 * @brief Retrieves the documentation for a specific function within a module as a formatted string.
 * It looks up the function in the module's function documentation map and returns its string representation.
 *
 * @param moduleName The name of the module containing the function.
 * @param functionName The name of the function to retrieve documentation for.
 * @return A string containing the formatted documentation for the specified function.
 */
func (m *ModuleDocs) GetFunctionDocAsString(moduleName string, functionName string) string {
	return m.Functions[moduleName].Functions[functionName].String()
}

/**
 * @brief Initializes and returns a new ModuleDocs struct.
 * It creates an empty map for storing function documentation.
 *
 * @return A ModuleDocs struct with an initialized Functions map.
 */
func newModuleDocs() ModuleDocs {
	return ModuleDocs{Functions: make(map[string]FunctionDocumentationMap)}
}
