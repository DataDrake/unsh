//
// Copyright 2017-2018 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package parser

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	idle = uint8(iota)
	envVars
	cmdName
	argVals
)

// Parser reads user input for a new command
type Parser struct {
	in    *bufio.Scanner
	state uint8
}

var p *Parser

func init() {
	p = NewParser(os.Stdin)
}

// NewParser returns a fully initialized Parser
func NewParser(in io.Reader) *Parser {
	return &Parser{bufio.NewScanner(in), idle}
}

// Next gets the next command to run
func Next() (env []string, cmd *exec.Cmd, err error) {
	return p.Parse()
}

// Parse scans a line of input and then tries to parse changes to the environment,
// the cmd, and any arguments
func (p *Parser) Parse() (env []string, cmd *exec.Cmd, err error) {

	// Check if there's nothing to read
	if !p.in.Scan() {
		// Get current error
		err = p.in.Err()
		// If there was no error, then stdin was closed
		if p.in.Err() == nil {
			err = io.EOF
		}
		return
	}
	p.state = envVars
	var name string
	args := make([]string, 0)
	// Parse cmdline
	for _, field := range strings.Fields(p.in.Text()) {
		switch p.state {
		case envVars:
			// Environment Variables
			if strings.ContainsRune(field, '=') {
				env = append(env, field)
				continue
			}
			// Done with environment variables
			p.state = cmdName
			fallthrough
		case cmdName:
			// Command Name/Path
			name = field
			p.state = argVals
		case argVals:
			// Arguments
			args = append(args, field)
		default:
			// do nothing
		}
	}
	// Return to IDLE
	p.state = idle

	// Check for possible command
	if len(name) > 0 {
		cmd = exec.Command(name, args...)
	}
	return
}
