package main

import (
	"fmt"
	"os"

	"rs-encoder/gf"
	"rs-encoder/rs"
	"rs-encoder/util"
)

const (
	// Primitive polynomial x^8 + x^4 + x^3 + x^2 + 1
	PrimitivePolynomial = 0x1D // Low 8 bits of 0x11D, as we only need the low 8 bits in GF
	// Number of data shards
	DataShards = 6
	// Number of parity shards
	ParityShards = 12
	// Total shards
	TotalShards = DataShards + ParityShards
)

func main() {
	// Check command line arguments
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input file> <output file>\n", os.Args[0])
		fmt.Println("For example: ./rs-encoder message.json encoded.json")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Read input data from JSON file
	message, err := util.ReadMessageFromJSON(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read message: %v\n", err)
		os.Exit(1)
	}

	// Check message length
	if len(message) != DataShards {
		fmt.Fprintf(os.Stderr, "Input message must contain %d shards, but got %d\n", DataShards, len(message))
		os.Exit(1)
	}

	// Initialize GF(2^8) finite field
	field := gf.NewGF(PrimitivePolynomial)
	fmt.Println("Initialized GF(2^8) finite field, using primitive polynomial:", fmt.Sprintf("0x%x", PrimitivePolynomial))

	// Create RS encoder
	encoder := rs.NewRSEncoder(field, DataShards, ParityShards)
	fmt.Printf("Created Reed-Solomon encoder, data shards: %d, parity shards: %d\n", DataShards, ParityShards)

	// Display input message
	fmt.Println("Input message:")
	for i, b := range message {
		fmt.Printf("[%d] = 0x%02x\n", i, b)
	}

	// Perform RS encoding
	fmt.Println("Performing Reed-Solomon encoding...")
	encoded := encoder.Encode(message)

	// Display encoding result
	fmt.Println("Encoding result:")
	for i, b := range encoded {
		if i < DataShards {
			fmt.Printf("[%d] = 0x%02x (original data)\n", i, b)
		} else {
			fmt.Printf("[%d] = 0x%02x (parity data)\n", i, b)
		}
	}

	// Validate systematic encoding
	isSystematic := true
	for i := 0; i < DataShards; i++ {
		if encoded[i] != message[i] {
			isSystematic = false
			break
		}
	}
	if isSystematic {
		fmt.Println("Systematic encoding validation succeeded: first", DataShards, "shards are identical to original data")
	} else {
		fmt.Println("Systematic encoding validation failed: first", DataShards, "shards are different from original data")
	}

	// Write encoding result to JSON file
	if err := util.WriteEncodedToJSONWithOriginal(outputFile, encoded, message); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write encoding result: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Encoding result has been written to:", outputFile)
}
