package repl

// func TestReplStart(t *testing.T) {
// 	type testCase struct {
// 		name           string
// 		input          []string
// 		expectedOutput []string
// 		expectedErrs   []string
// 	}

// 	testCases := []testCase{
// 		{
// 			name: "test REPL can evaluate simple expressions",
// 			input: []string{
// 				"var x = 10;",
// 				"x;",
// 			},
// 			expectedOutput: []string{
// 				"10",
// 				"10",
// 			},
// 			expectedErrs: []string{},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			inReader, inWriter := io.Pipe()
// 			outReader, outWriter := io.Pipe()

// 			r := New(inReader, outWriter, nil)
// 			go r.Start()

// 			drainDonkeyASCII(t, outReader)

// 			for i, cmd := range tc.input {
// 				// Send the command
// 				_, err := fmt.Fprintf(inWriter, "%s\n", cmd)
// 				assert.NoError(t, err)

// 				// read output until next prompt
// 				output := readUntilPrompt(t, outReader, i+1)
// 				output = strings.TrimSpace(output)

// 				assert.Equal(t, tc.expectedOutput[i], output)
// 			}

// 			inWriter.Close()
// 		})
// 	}
// }

// func readUntilPrompt(t *testing.T, r io.Reader, lineNumber int) string {
// 	t.Helper()
// 	var buf strings.Builder
// 	b := make([]byte, 1)

// 	nextPrompt := fmt.Sprintf(inPromptTemplate, lineNumber+1)

// 	for {
// 		_, err := r.Read(b)
// 		if err != nil {
// 			return buf.String()
// 		}

// 		buf.WriteByte(b[0])
// 		s := buf.String()
// 		if strings.HasSuffix(s, nextPrompt) {
// 			// Strip the trailing next prompt before returning
// 			return s[:len(s)-len(nextPrompt)]
// 		}
// 	}
// }

// func drainDonkeyASCII(t *testing.T, r io.Reader) {
// 	t.Helper()
// 	var buf strings.Builder
// 	b := make([]byte, 1)

// 	// read until we encounter the welcome message, which indicates the end of the ASCII art
// 	for {
// 		_, err := r.Read(b)
// 		if err != nil {
// 			return
// 		}

// 		buf.WriteByte(b[0])
// 		s := buf.String()
// 		if strings.Contains(s, "Welcome to the Donkey programming language!") {
// 			return
// 		}
// 	}
// }
