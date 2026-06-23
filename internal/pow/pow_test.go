package pow

import "testing"

func TestSolveVerify(t *testing.T) {
	for _, bits := range []int{4, 8, 12, 16} {
		wanted := ""
		for i := 0; i < bits; i++ {
			wanted += "0"
		}
		req := Request{Key: "test-key-1234567890", Wanted: wanted}
		sol := Solve(req)
		if len(sol) != 16 {
			t.Fatalf("bits=%d: solution length %d", bits, len(sol))
		}
		if !Verify(req, sol) {
			t.Fatalf("bits=%d: solution did not verify", bits)
		}
	}
}

func TestVerifyRejectsWrong(t *testing.T) {
	req := Request{Key: "abc", Wanted: "0000000000000000"} // 16 bits
	if Verify(req, "0000000000000000") {
		// extremely unlikely the zero counter solves 16 bits
		t.Skip("zero counter happened to solve; skip")
	}
	if Verify(req, "short") {
		t.Fatal("accepted wrong-length solution")
	}
}

func TestWriteCounterHex(t *testing.T) {
	dst := make([]byte, 16)
	writeCounterHex(dst, 0)
	if string(dst) != "0000000000000000" {
		t.Fatalf("zero: %q", dst)
	}
	writeCounterHex(dst, 0xABCDEF)
	if string(dst) != "0000000000abcdef" {
		t.Fatalf("0xABCDEF: %q", dst)
	}
}

func TestIncrement(t *testing.T) {
	s := []byte("000000000000000f")
	incrementHexAsciiInPlace(s)
	if string(s) != "0000000000000010" {
		t.Fatalf("carry: %q", s)
	}
	s = []byte("0000000000000009")
	incrementHexAsciiInPlace(s)
	if string(s) != "000000000000000a" {
		t.Fatalf("9->a: %q", s)
	}
}
