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

package main

import (
    "github.com/DataDrake/unsh/parser"
    "github.com/DataDrake/unsh/prompt"
    "io"
    "os"
)

func main() {

    // Loop until error or close
    for {
        // Print Prompt
        prompt.Print()
        // Get Command
        _, cmd, err := parser.Next()
        // Ctrl-D
        if err == io.EOF {
            os.Exit(0)
        }
        if err != nil {
           panic(err.Error())
        }
        // Run Command
        if cmd != nil {
            // Setup
            cmd.Stderr = os.Stderr
            cmd.Stdout = os.Stdout
            cmd.Stdin = os.Stdin
            // Run
            err := cmd.Start()
            if err != nil {
                println(err.Error())
            } else {
                // Wait until finished
                err := cmd.Wait()
                if err != nil {
                    println(err.Error())
                }
            }
        }
    }
}
