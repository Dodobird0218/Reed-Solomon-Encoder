# Reed-Solomon Encoder

## 1. Project Overview

This project implements a Reed-Solomon error correction code encoder in Go. Reed-Solomon codes are a type of forward error correction (FEC) codes widely used in digital storage and communication systems to detect and correct errors during data transmission or storage.

The encoder takes a message and produces redundant parity data that can be used to recover the original message even if parts of it are corrupted (up to a certain limit). This implementation uses a systematic encoding approach, meaning the original data is preserved as the first part of the encoded output.

Key features:
- Implementation of Galois Field GF(2^8) arithmetic
- Support for configurable data shards and parity shards
- JSON input/output support
- Systematic encoding (original data is preserved)

## 2. Mathematical Background

### Galois Field Theory

Reed-Solomon encoding operates on Galois Fields (finite fields). This implementation uses GF(2^8), which:
- Contains 256 elements (0 to 255)
- Uses polynomial arithmetic modulo an irreducible polynomial
- All operations (addition, subtraction, multiplication, division) are closed within the field

### Reed-Solomon Coding Theory

Reed-Solomon codes are based on polynomial interpolation over finite fields:
- Data is treated as coefficients of a polynomial
- Encoding involves evaluating this polynomial at specific points
- The encoded message consists of the original data plus the polynomial evaluations
- Error correction capability depends on the number of redundant (parity) symbols
- With `n` parity symbols, up to `n/2` errors can be corrected

### Lagrange Interpolation

The encoder uses Lagrange interpolation to calculate parity shards:
- Creates a polynomial that passes through given data points
- Evaluates this polynomial at different points to generate parity data
- This approach ensures the systematic property of the encoding (original data preserved)

## 3. Implementation Details

### GF(2^8) Implementation

- Exponent and logarithm tables are pre-computed for efficient field operations
- Field operations implemented:
  - Addition (XOR operation)
  - Subtraction (same as addition in GF(2^8))
  - Multiplication (using logarithm tables)
  - Division (using logarithm tables)
  - Power and inverse functions

### Reed-Solomon Encoder

The encoder implements the following encoding method:
1. **Lagrange Interpolation Method**:
   - Constructs a polynomial passing through data points
   - Evaluates the polynomial at specific points to generate parity shards

### File I/O Handling

- JSON format used for input/output
- Support for reading messages from JSON files
- Support for writing encoded data to JSON files with original message included

## 4. Usage Instructions

### Prerequisites

- Go programming language (1.16 or later recommended)
- No external dependencies required

### Building the Project

```bash
# Clone the repository
git clone [repository-url]
cd rs-encoder

# Build the project
go build -o rs-encoder cmd/main.go
```

### Running the Encoder

```bash
# Basic usage
./rs-encoder <input-file> <output-file>

# Example
./rs-encoder data/message.json encoded.json
```

### Input Format

The input JSON file should have the following format:
```json
{
  "message": [
    "0x00", "0x01", "0x02", "0x03", "0x04", "0x05"
  ]
}
```

### Output Format

The output JSON file will have the following format:
```json
{
  "message": [
    "0x00", "0x01", "0x02", "0x03", "0x04", "0x05"
  ],
  "encoded": [
    "0x00", "0x01", "0x02", "0x03", "0x04", "0x05",
    "0xXX", "0xXX", "0xXX", "0xXX", "0xXX", "0xXX",
    "0xXX", "0xXX", "0xXX", "0xXX", "0xXX", "0xXX"
  ]
}
```

Where the first 6 values in `encoded` are the original message, and the last 12 values are the parity data.

### Configuration

Constants in `main.go` can be modified to change:
- The primitive polynomial used in the finite field
- The number of data shards
- The number of parity shards