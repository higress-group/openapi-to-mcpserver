package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/higress-group/openapi-to-mcpserver/pkg/converter"
	"github.com/higress-group/openapi-to-mcpserver/pkg/models"
	"github.com/higress-group/openapi-to-mcpserver/pkg/parser"
	"github.com/higress-group/openapi-to-mcpserver/pkg/validator"
	"gopkg.in/yaml.v3"
)

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the OpenAPI specification file (JSON or YAML)")
	outputFile := flag.String("output", "", "Path to the output MCP configuration file (YAML)")
	serverName := flag.String("server-name", "openapi-server", "Name of the MCP server")
	toolNamePrefix := flag.String("tool-prefix", "", "Prefix for tool names")
	format := flag.String("format", "yaml", "Output format (yaml or json)")
	validate := flag.Bool("validate", false, "Validate the OpenAPI specification")
	mcpCompat := flag.String("mcp-compat", "", "MCP compatibility check mode: 'strict' (error on incompatible operations) or 'warn' (skip with warnings)")
	templateFile := flag.String("template", "", "Path to a template file to patch the output")

	// Parse command-line flags
	flag.Parse()

	// Validate required flags
	if *inputFile == "" {
		fmt.Println("Error: input file is required")
		flag.Usage()
		os.Exit(1)
	}

	if *outputFile == "" {
		fmt.Println("Error: output file is required")
		flag.Usage()
		os.Exit(1)
	}

	// Validate MCP compatibility mode flag
	compatMode, valid := validator.ParseMCPCompatMode(*mcpCompat)
	if !valid {
		fmt.Printf("Error: invalid mcp-compat mode %q (must be 'strict' or 'warn')\n", *mcpCompat)
		os.Exit(1)
	}

	// Create a new parser
	p := parser.NewParser()

	// Set validation option
	p.SetValidation(*validate)

	// Parse the OpenAPI specification
	err := p.ParseFile(*inputFile)
	if err != nil {
		fmt.Printf("Error parsing OpenAPI specification: %v\n", err)
		os.Exit(1)
	}

	// Run MCP compatibility validation if enabled
	var compatResult *validator.MCPCompatResult
	if compatMode != validator.MCPCompatOff {
		compatResult = validator.ValidateDocument(p.GetDocument())

		if len(compatResult.Issues) > 0 {
			if compatMode == validator.MCPCompatStrict {
				fmt.Println("Error: OpenAPI specification contains MCP-incompatible operations:")
				for _, issue := range compatResult.Issues {
					fmt.Printf("  %s\n", issue.String())
				}
				os.Exit(1)
			}
			// Warn mode: print warnings, will skip these operations during conversion
			for _, issue := range compatResult.Issues {
				fmt.Printf("Warning: %s\n", issue.String())
			}
			fmt.Println()
		}
	}

	// Create a new converter
	c := converter.NewConverter(p, models.ConvertOptions{
		ServerName:     *serverName,
		ToolNamePrefix: *toolNamePrefix,
		TemplatePath:   *templateFile,
	})

	// Set skip filter if in warn mode with incompatible operations
	if compatResult != nil && compatMode == validator.MCPCompatWarn {
		c.SkipFilter = func(path, method string) bool {
			return !compatResult.IsCompatible(path, method)
		}
	}

	// Convert the OpenAPI specification to an MCP configuration
	config, err := c.Convert()
	if err != nil {
		fmt.Printf("Error converting OpenAPI specification: %v\n", err)
		os.Exit(1)
	}

	// Create the output directory if it doesn't exist
	outputDir := filepath.Dir(*outputFile)
	if outputDir != "" && outputDir != "." {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Write the MCP configuration to the output file
	var data []byte
	if *format == "json" {
		data, err = json.MarshalIndent(config, "", "  ")
	} else {
		var buffer bytes.Buffer
		encoder := yaml.NewEncoder(&buffer)
		encoder.SetIndent(2)

		if err := encoder.Encode(config); err != nil {
			fmt.Printf("Error encoding YAML: %v\n", err)
			return
		}
		data = buffer.Bytes()
	}
	if err != nil {
		fmt.Printf("Error marshaling MCP configuration: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing MCP configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted OpenAPI specification to MCP configuration: %s\n", *outputFile)
}
