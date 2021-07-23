package acccore

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// UniqueIDGenerator define the unique string generator. The unique string generated MUST be widely system unique, even
// accross nodes.
type UniqueIDGenerator interface {
	// NewUniqueID will produce a unique ID string.
	NewUniqueID() string
}

// UUIDUniqueIDGenerator the unique ID generator using UUID
type UUIDUniqueIDGenerator struct{}

// NewUniqueID will produce a unique ID string.
func (gen *UUIDUniqueIDGenerator) NewUniqueID() string {
	return uuid.New().String()
}

var (
	nanoSince = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
)

// NanoSecondUniqueIDGenerator the unique ID generator using NANO second number.
// This is relying on generated nanosecond generated. Its a number of nanosecond since January 1st 2021 at 0:0:0:0.
// We hope the generated ID is unique over time
type NanoSecondUniqueIDGenerator struct{}

// NewUniqueID will produce a unique ID string.
func (gen *NanoSecondUniqueIDGenerator) NewUniqueID() string {
	return fmt.Sprintf("%d", time.Now().Sub(nanoSince).Nanoseconds())
}

const (
	LowerAlphabet = "abcdefghijklmnopqrstuvwxyz"
	UpperAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers       = "1234567890"
	Symbols       = "!@#$%^&*(){}[]|;:<>,./?~"
)

// RandomGenUniqueIDGenerator the unique ID generator using NANO second number.
type RandomGenUniqueIDGenerator struct {
	Length        int
	LowerAlpha    bool
	UpperAlpha    bool
	Numeric       bool
	Symbols       bool
	CharSetBuffer []byte
}

func (gen *RandomGenUniqueIDGenerator) NewUniqueID() string {
	var buff bytes.Buffer
	if gen.CharSetBuffer == nil {
		if gen.LowerAlpha == false && gen.UpperAlpha == false && gen.Numeric == false {
			buff.WriteString(UpperAlphabet)
			buff.WriteString(Numbers)
		} else {
			if gen.LowerAlpha {
				buff.WriteString(LowerAlphabet)
			}
			if gen.UpperAlpha {
				buff.WriteString(UpperAlphabet)
			}
			if gen.Numeric {
				buff.WriteString(Numbers)
			}
			if gen.Symbols {
				buff.WriteString(Symbols)
			}
		}
		gen.CharSetBuffer = buff.Bytes()
	}
	if gen.Length == 0 {
		gen.Length = 16
	}
	l := len(gen.CharSetBuffer)
	buff.Reset()
	for buff.Len() < gen.Length {
		r := rand.Intn(l)
		buff.Write(gen.CharSetBuffer[r : r+1])
	}
	return buff.String()
}
