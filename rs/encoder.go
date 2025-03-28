package rs

import (
	"rs-encoder/gf"
)

// RSEncoder Reed-Solomon encoder
type RSEncoder struct {
	field        *gf.GF
	dataShards   int
	parityShards int
	totalShards  int
	alphaPoints  []byte
}

// NewRSEncoder creates a new Reed-Solomon encoder
func NewRSEncoder(field *gf.GF, dataShards, parityShards int) *RSEncoder {
	encoder := &RSEncoder{
		field:        field,
		dataShards:   dataShards,
		parityShards: parityShards,
		totalShards:  dataShards + parityShards,
	}
	encoder.generateAlphaPoints()
	return encoder
}

// Encode encodes the message using Reed-Solomon encoding
func (enc *RSEncoder) Encode(message []byte) []byte {
	if len(message) != enc.dataShards {
		panic("Message length must equal the number of data shards")
	}

	// Create encoded result array, first dataShards items are same as original message (systematic encoding)
	encoded := make([]byte, enc.totalShards)
	copy(encoded, message)

	// Use Lagrange interpolation to calculate the remaining redundant shards
	enc.lagrangeInterpolation(message, encoded)

	return encoded
}

// generateAlphaPoints generates alpha evaluation points
func (enc *RSEncoder) generateAlphaPoints() {
	enc.alphaPoints = make([]byte, enc.totalShards)

	// Use a^0, a^1, a^2, ..., a^(totalShards-1) as evaluation points
	// For GF(2^8), we typically choose a=2 as the primitive element
	alpha := byte(2)

	// First point is a^0 = 1
	enc.alphaPoints[0] = 1

	// Calculate remaining evaluation points
	for i := 1; i < enc.totalShards; i++ {
		enc.alphaPoints[i] = enc.field.Mul(enc.alphaPoints[i-1], alpha)
	}
}

// lagrangeInterpolation calculates redundant shards using Lagrange interpolation
func (enc *RSEncoder) lagrangeInterpolation(message []byte, encoded []byte) {
	// For systematic Reed-Solomon encoding, ensure the first dataShards items are identical to the original message

	// For each redundant shard position
	for i := enc.dataShards; i < enc.totalShards; i++ {
		// Calculate polynomial value at this point
		result := byte(0)

		// Construct Lagrange interpolation polynomial
		for j := 0; j < enc.dataShards; j++ {
			term := message[j]

			// Calculate Lagrange basis function
			for k := 0; k < enc.dataShards; k++ {
				if j != k {
					// Calculate (x - x_k)
					numerator := enc.field.Sub(enc.alphaPoints[i], enc.alphaPoints[k])
					// Calculate (x_j - x_k)
					denominator := enc.field.Sub(enc.alphaPoints[j], enc.alphaPoints[k])
					// Division
					factor := enc.field.Div(numerator, denominator)
					// Multiply by the current term
					term = enc.field.Mul(term, factor)
				}
			}

			// Add this term to the result
			result = enc.field.Add(result, term)
		}

		encoded[i] = result
	}
}

// Another more efficient encoding implementation method (optional)
func (enc *RSEncoder) EncodeEfficient(message []byte) []byte {
	if len(message) != enc.dataShards {
		panic("Message length must equal the number of data shards")
	}

	// Create encoded result array, first dataShards items are same as original message
	encoded := make([]byte, enc.totalShards)
	copy(encoded, message)

	// Build encoding matrix
	// For each redundant shard position
	for i := enc.dataShards; i < enc.totalShards; i++ {
		result := byte(0)

		// Calculate polynomial value at evaluation point
		x := enc.alphaPoints[i]

		// Use Horner's method to calculate polynomial value
		// p(x) = message[0] + message[1]*x + message[2]*x^2 + ... + message[dataShards-1]*x^(dataShards-1)
		for j := enc.dataShards - 1; j >= 0; j-- {
			result = enc.field.Add(enc.field.Mul(result, x), message[j])
		}

		encoded[i] = result
	}

	return encoded
}
