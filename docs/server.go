package docs

import (
	"KamaiZen/settings"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var moduleDocsCache = make(map[string]ModuleDocs)

type FunctionDoc struct {
	Name        string
	Parameters  string
	Description string
	Example     string
}

func (f FunctionDoc) String() string {
	return fmt.Sprintf("Function: %s\n\nParameters: %s\n\nDescription: %s\n\nExample:\n```\n%s\n```", f.Name, f.Parameters, f.Description, f.Example)
}

type FunctionDocs struct {
	Functions map[string]FunctionDoc
}

func (f *FunctionDocs) AddFunctionDoc(functionDoc FunctionDoc) {
	f.Functions[functionDoc.Name] = functionDoc
}

func (f *FunctionDocs) GetFunctionDocAsString(functionName string) string {
	return f.Functions[functionName].String()
}

type ModuleDocs struct {
	Functions map[string]FunctionDocs
}

func (m *ModuleDocs) AddFunctionDoc(moduleName string, functionDocs FunctionDocs) {
	m.Functions[moduleName] = functionDocs
}

func (m *ModuleDocs) GetFunctionDocAsString(moduleName string, functionName string) string {
	return m.Functions[moduleName].Functions[functionName].String()
}

func NewModuleDocs() ModuleDocs {
	return ModuleDocs{Functions: make(map[string]FunctionDocs)}
}

func listTableOfContents(readme string) {
	lines := strings.Split(readme, "\n")
	tocPattern := regexp.MustCompile(`^\s*\d+\.\s*(\w+)\s*$`)
	for _, line := range lines {
		if match := tocPattern.FindStringSubmatch(line); match != nil {
			fmt.Println(match[1])
		}
	}
}

func listFunctions(readme string) []string {
	lines := strings.Split(readme, "\n")
	funcPattern := regexp.MustCompile(`^\s*\d+\.\d+\.\s*(\w+)\((.*)\)\s*$`)
	listofFunctions := []string{}
	for _, line := range lines {
		if match := funcPattern.FindStringSubmatch(line); match != nil {
			// don't add the function if it is a duplicate
			if !strings.Contains(strings.Join(listofFunctions, ","), match[1]) {
				listofFunctions = append(listofFunctions, match[1])
			}
		}
	}
	return listofFunctions
}

func getFunction(readme string, functionName string) FunctionDoc {
	lines := strings.Split(readme, "\n")
	funcPattern := regexp.MustCompile(`^\s*\d+\.\d+\.\s*(\w+)\((.*)\)\s*$`)
	var functionDoc FunctionDoc
	for _, line := range lines {
		if match := funcPattern.FindStringSubmatch(line); match != nil {
			if match[1] == functionName {
				functionDoc.Name = match[1]
				functionDoc.Parameters = match[2]
			}
		}
	}

	return functionDoc
}

// example doc
// 7.2.  acc_db_request(comment, table)
//
//	Like acc_log_request, acc_db_request reports on a request. The report
//	is sent to database at “db_url”, in the table referred to in the second
//	action parameter.
//
//	Meaning of the parameters is as follows:
//	  * comment - Comment to be appended. The string can contain any number
//	    of pseudo-variables.
//	  * table - Database table to be used. It can contain config variables
//	    that are evaluated at runtime.
//
//	This function can be used from ANY_ROUTE.
//
//	Example 1.52. acc_db_request usage
//
// ...
// acc_db_request("Some comment", "SomeTable");
// acc_db_request("Some comment", "acc_$time(year)_$time(mon)");
// acc_db_request("$var(code) Error: $avp(reason)", "SomeTable");
// ...
func extractFunctionDoc(lines []string) []FunctionDoc {
	funcPattern := regexp.MustCompile(`^\s*\d+\.\d+\.\s*(\w+)\((.*)\)\s*$`)
	var functionDocs []FunctionDoc
	var functionDoc FunctionDoc
	var inExample bool
	var example string
	var exampleLineCount int
	for _, line := range lines {
		if match := funcPattern.FindStringSubmatch(line); match != nil {
			if functionDoc.Name != "" {
				functionDocs = append(functionDocs, functionDoc)
			}
			functionDoc = FunctionDoc{Name: match[1], Parameters: match[2]}
			inExample = false
			example = ""
		} else if functionDoc.Name != "" {
			if !inExample {
				// example has indentation of some spaces
				if strings.Contains(line, "Example") {
					inExample = true
				} else {
					functionDoc.Description += line + "\n"
				}
			} else {
				if strings.Contains(line, "...") {
					exampleLineCount++
					// end of example if there are 2 lines
					if exampleLineCount == 2 {
						inExample = false
					}
					functionDoc.Example = example
				} else {
					example += line + "\n"
				}
			}

		}
	}
	if functionDoc.Name != "" {
		functionDocs = append(functionDocs, functionDoc)
	}
	return functionDocs
}

func Initialise(s settings.LSPSettings) {
	path := s.KamailioSourcePath() + "/src/modules"
	listOfModules, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory")
		return
	}
	// Get All Modules
	for _, module := range listOfModules {
		// Parse README file for each module
		readme, err := os.ReadFile(path + "/" + module.Name() + "/README")
		if err != nil {
			fmt.Println("Error reading file")
			return
		}

		moduleDocs := NewModuleDocs()
		functionDocs := extractFunctionDoc(strings.Split(string(readme), "\n"))
		functionDocsMap := FunctionDocs{Functions: make(map[string]FunctionDoc)}
		for _, functionDoc := range functionDocs {
			functionDocsMap.AddFunctionDoc(functionDoc)
		}
		moduleDocs.AddFunctionDoc(module.Name(), functionDocsMap)
		if module.Name() == "xlog" {
			fmt.Println(moduleDocs.GetFunctionDocAsString(module.Name(), "xwarn"))
		}
		moduleDocsCache[module.Name()] = moduleDocs
	}
}

func GetFunctionDoc(moduleName string, functionName string) string {
	moduleDocs, exists := moduleDocsCache[moduleName]
	if !exists {
		return "Module not found"
	}
	return moduleDocs.GetFunctionDocAsString(moduleName, functionName)
}

func FindFunctionInAllModules(functionName string) string {
	for moduleName, moduleDocs := range moduleDocsCache {
		if _, exists := moduleDocs.Functions[moduleName].Functions[functionName]; exists {
			return "Module: " + moduleName + "\n" + moduleDocs.GetFunctionDocAsString(moduleName, functionName)
		}
	}
	return "Function not found"
}
