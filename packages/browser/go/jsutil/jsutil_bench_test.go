package jsutil

import (
	"fmt"
	"math/cmplx"
	"math/rand"
	"testing"
	"time"

	"github.com/pborman/uuid"
)

type testValue struct {
	Uint    uint      `json:"uint"`
	Uint8   uint8     `json:"uint8"`
	Uint16  uint16    `json:"uint16"`
	Uint32  uint32    `json:"uint32"`
	Uint64  uint64    `json:"uint64"`
	Int     int       `json:"int"`
	Int8    int8      `json:"int8"`
	Int16   int16     `json:"int16"`
	Int32   int32     `json:"int32"`
	Int64   int64     `json:"int64"`
	Float32 float32   `json:"float32"`
	Float64 float64   `json:"float64"`
	Byte    byte      `json:"byte"`
	Rune    rune      `json:"rune"`
	String  string    `json:"string"`
	Bool    bool      `json:"bool"`
	Time    time.Time `json:"time"`
}

func newTestValue() *testValue {
	return &testValue{
		Uint:    uint(randomInt()),
		Uint8:   uint8(randomInt()),
		Uint16:  uint16(randomInt()),
		Uint32:  uint32(randomInt()),
		Uint64:  uint64(randomInt()),
		Int:     randomInt(),
		Int8:    int8(randomInt()),
		Int16:   int16(randomInt()),
		Int32:   int32(randomInt()),
		Int64:   int64(randomInt()),
		Float32: float32(randomInt()),
		Float64: float64(randomInt()),
		Byte:    []byte(randomString())[0],
		Rune:    []rune(randomString())[0],
		String:  randomString(),
		Bool:    randomBool(),
		Time:    time.Now(),
	}
}

func BenchmarkInefficientlyConvertToJS(b *testing.B) {
	counts := []int{
		1,
		10,
		100,
		1000,
	}
	for _, count := range counts {
		name := fmt.Sprintf("BenchmarkInefficientlyConvertToJS_%d", count)
		b.Run(name, runBenchmarkInefficientlyConvertToJS(count))
	}
}

func runBenchmarkInefficientlyConvertToJS(count int) func(*testing.B) {
	return func(b *testing.B) {
		values := []*testValue{}
		for i := 0; i < count; i++ {
			values = append(values, newTestValue())
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = InefficientlyConvertToJS(values)
		}
	}
}

// randomInt returns a pseudo-random int between the minimum and maximum
// possible values.
func randomInt() int {
	return rand.Int()
}

// randomString returns a random string of length 16
func randomString() string {
	return uuid.New()
}

// randomBool returns a random bool
func randomBool() bool {
	return rand.Int()%2 == 0
}

// randomFloat returns a random float64
func randomFloat() float64 {
	return rand.Float64()
}

// randomComplex returns a random complex128
func randomComplex() complex128 {
	return cmplx.Rect(randomFloat(), randomFloat())
}
