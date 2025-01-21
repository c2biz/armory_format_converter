package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// BOFArgsBuffer handles packing arguments for the BOF
type BOFArgsBuffer struct {
	Buffer *bytes.Buffer
}

// AddInt adds a 4-byte integer to the buffer
func (b *BOFArgsBuffer) AddInt(val uint32) error {
	return binary.Write(b.Buffer, binary.LittleEndian, val)
}

// AddString adds a zero-terminated string to the buffer
func (b *BOFArgsBuffer) AddString(val string) error {
	_, err := b.Buffer.WriteString(val + "\x00")
	return err
}

// AddData adds binary data to the buffer
func (b *BOFArgsBuffer) AddData(data []byte) error {
	_, err := b.Buffer.Write(data)
	return err
}

// Command defines a BOF command schema
type Command struct {
	Name      string     `json:"name"`
	Help      string     `json:"help"`
	Arguments []Argument `json:"arguments"`
}

// Argument defines an argument schema for a command
type Argument struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // int, string, file, etc.
	Optional bool   `json:"optional"`
}

// ParseCommand dynamically parses a command and its arguments
func ParseCommand(input string, schema Command) ([]byte, error) {
	// Parse the input into flags
	args := strings.Fields(input)
	if len(args) == 0 || args[0] != schema.Name {
		return nil, fmt.Errorf("invalid command: expected %s", schema.Name)
	}

	// Set up a flag set for dynamic parsing
	flags := flag.NewFlagSet(schema.Name, flag.ContinueOnError)
	parsedArgs := make(map[string]interface{})
	for _, arg := range schema.Arguments {
		switch arg.Type {
		case "int":
			parsedArgs[arg.Name] = flags.Int(arg.Name, 0, "integer argument")
		case "string":
			parsedArgs[arg.Name] = flags.String(arg.Name, "", "string argument")
		case "file":
			parsedArgs[arg.Name] = flags.String(arg.Name, "", "file path argument")
		}
	}

	// Parse flags
	err := flags.Parse(args[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %v", err)
	}

	// Validate required arguments
	for _, arg := range schema.Arguments {
		if !arg.Optional {
			value := parsedArgs[arg.Name]
			switch v := value.(type) {
			case *int:
				if *v == 0 {
					return nil, fmt.Errorf("missing required argument: %s", arg.Name)
				}
			case *string:
				if *v == "" {
					return nil, fmt.Errorf("missing required argument: %s", arg.Name)
				}
			}
		}
	}

	// Pack the BOF arguments
	argsBuffer := &BOFArgsBuffer{Buffer: new(bytes.Buffer)}
	for _, arg := range schema.Arguments {
		value := parsedArgs[arg.Name]
		switch v := value.(type) {
		case *int:
			if err := argsBuffer.AddInt(uint32(*v)); err != nil {
				return nil, err
			}
		case *string:
			if arg.Type == "file" {
				data, err := os.ReadFile(*v)
				if err != nil {
					return nil, fmt.Errorf("failed to read file `%s`: %v", *v, err)
				}
				if err := argsBuffer.AddData(data); err != nil {
					return nil, err
				}
			} else {
				if err := argsBuffer.AddString(*v); err != nil {
					return nil, err
				}
			}
		}
	}

	return argsBuffer.Buffer.Bytes(), nil
}

// ExecuteBOF simulates executing the BOF
func ExecuteBOF(buffer []byte) {
	fmt.Printf("Executing BOF with buffer: %x\n", buffer)
}

func main() {
	// Simulated command schema (normally loaded from JSON)
	commandSchema := `{
		"name": "nanodump_ppl_dump",
		"help": "Bypass PPL and dump LSASS.",
		"arguments": [
			{"name": "write", "type": "string", "optional": false},
			{"name": "valid", "type": "int", "optional": true},
			{"name": "duplicate", "type": "int", "optional": true},
			{"name": "file", "type": "file", "optional": true}
		]
	}`

	// Load the schema
	var schema Command
	err := json.Unmarshal([]byte(commandSchema), &schema)
	if err != nil {
		fmt.Println("Failed to load command schema:", err)
		os.Exit(1)
	}

	// Simulated input (e.g., from CLI)
	input := "nanodump_ppl_dump --write C:\\Windows\\Temp\\dump.dmp --valid 1 --duplicate 1 --file test.bin"

	// Parse and pack the command
	packedArgs, err := ParseCommand(input, schema)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Execute the BOF
	ExecuteBOF(packedArgs)
}
