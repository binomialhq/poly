package poly

import (
    "testing"
    "strings"
    "bufio"
    "errors"
    "unicode"
    "io"
    "fmt"
)

// protoErrorReader is a prototype Reader that errors when Read is called. This
// is used for testing, for the purposes of test coverage.
type protoErrorReader struct {
}

// A bogus Read method that conform with to Reader interface. This method simply
// returns an error. Useful for testing purposes only.

// Read is a bogus implementation that conforms to the Reader interface. This
// method simply returns an error. Useful for testing.
func (reader *protoErrorReader) Read(p []byte) (n int, err error) {
    return 0, errors.New("something went wrong")
}

// TestIsLetter tests a set of unicode characters from several classes to
// assess whether the IsLetter routine is likely to accept a proper range
// of unicode characters.
func TestIsLetter(t *testing.T) {

    // Must accept L-class unicode characters
    if !IsLetter('a') { t.Log("L-class character unexpectedly rejected"); t.Fail() }
    if !IsLetter('ɮ') { t.Log("L-class character unexpectedly rejected"); t.Fail() }
    if !IsLetter('ʶ') { t.Log("L-class character unexpectedly rejected"); t.Fail() }
    if !IsLetter('ꫳ') { t.Log("L-class character unexpectedly rejected"); t.Fail() }
    if !IsLetter('ن') { t.Log("L-class character unexpectedly rejected"); t.Fail() }
    if !IsLetter('ݼ') { t.Log("L-class character unexpectedly rejected"); t.Fail() }

    // Must reject Z-class unicode characters
    if IsLetter('\t') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('\n') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('\v') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('\f') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('\r') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsLetter(' ') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('\u0085') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('\u00A0') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }

    // Must reject N-class unicode characters
    if IsLetter('0') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('𐒦') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('꧳') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('ⅷ') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('〥') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('᱙') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('¼') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('𒑣') { t.Log("N-class character unexpectedly accepted"); t.Fail() }

    // Must accept P-class character class only
    if !IsLetter('_') { t.Log("P-class character unexpectedly rejected"); t.Fail() }
    if IsLetter('—') { t.Log("P-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('゠') { t.Log("P-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('[') { t.Log("P-class character unexpectedly accepted"); t.Fail() }
    if IsLetter(']') { t.Log("P-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('»') { t.Log("P-class character unexpectedly accepted"); t.Fail() }
    if IsLetter('«') { t.Log("P-class character unexpectedly accepted"); t.Fail() }
}

// TestIsDecimal tests a set of unicode characters from several classes to
// assess whether the IsDecimal routine is likely to accept a proper range
// of unicode characters.
func TestIsDecimal(t *testing.T) {

    // Must accept all arab decimal digits
    if !IsDecimal('0') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('1') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('2') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('3') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('4') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('5') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('6') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('7') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('8') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }
    if !IsDecimal('9') { t.Log("arab numeral character unexpectedly rejected"); t.Fail() }

    // Everything else must be rejected
    if IsDecimal('a') { t.Log("L-class character unexpectedly accepted"); t.Fail() }
    if IsDecimal('ɮ') { t.Log("L-class character unexpectedly accepted"); t.Fail() }
    if IsDecimal('\t') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsDecimal('\n') { t.Log("Z-class character unexpectedly accepted"); t.Fail() }
    if IsDecimal('𐒦') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsDecimal('¼') { t.Log("N-class character unexpectedly accepted"); t.Fail() }
    if IsDecimal('_') { t.Log("P-class character unexpectedly accepted"); t.Fail() }
    if IsDecimal('—') { t.Log("P-class character unexpectedly accepted"); t.Fail() }
}

// TestIsWhitespace tests a set of unicode characters from several classes to
// assess whether the IsWhitespace routine is likely to accept a proper range
// of unicode characters.
func TestIsWhitespace(t *testing.T) {

    // Must accept Latin-1 spaces
    if !IsWhitespace('\u0009') { t.Log("unexpectedly rejected Latin-1 space character"); }  // \t
    if !IsWhitespace('\u00A0') { t.Log("unexpectedly rejected Latin-1 space character"); }  // \n
    if !IsWhitespace('\u000B') { t.Log("unexpectedly rejected Latin-1 space character"); }  // \v
    if !IsWhitespace('\u000C') { t.Log("unexpectedly rejected Latin-1 space character"); }  // \f
    if !IsWhitespace('\u000D') { t.Log("unexpectedly rejected Latin-1 space character"); }  // \r
    if !IsWhitespace('\u0020') { t.Log("unexpectedly rejected Latin-1 space character"); }  // Space
    if !IsWhitespace('\u0085') { t.Log("unexpectedly rejected Latin-1 space character"); }  // NEL
    if !IsWhitespace('\u00A0') { t.Log("unexpectedly rejected Latin-1 space character"); }  // NBSP

    // Must accept category Z
    if !IsWhitespace('\u2028') { t.Log("unexpectedly rejected category Z space character"); }  // Line separator
    if !IsWhitespace('\u2029') { t.Log("unexpectedly rejected category Z space character"); }  // Paragraph separator
    if !IsWhitespace('\u1680') { t.Log("unexpectedly rejected category Z space character"); }  // Ogham space mark
    if !IsWhitespace('\u2000') { t.Log("unexpectedly rejected category Z space character"); }  // En quad
    if !IsWhitespace('\u2001') { t.Log("unexpectedly rejected category Z space character"); }  // Em quad
    if !IsWhitespace('\u3000') { t.Log("unexpectedly rejected category Z space character"); }  // Ideaographic space

    // Must reject L-class unicode characters
    if IsWhitespace('a') { t.Log("unexpectedly accepted category L character"); t.Fail() }
    if IsWhitespace('ɮ') { t.Log("unexpectedly accepted category L character"); t.Fail() }

    // Must reject N-class unicode characters
    if IsWhitespace('0') { t.Log("unexpectedly accepted category N character"); t.Fail() }
    if IsWhitespace('𐒦') { t.Log("unexpectedly accepted category N character"); t.Fail() }

    // Must reject P-class unicode characters
    if IsWhitespace('_') { t.Log("unexpectedly accepted category P character"); t.Fail() }
    if IsWhitespace('—') { t.Log("unexpectedly accepted category P character"); t.Fail() }
}

// TestReadRuneAndCount tests normal behaviour for readRuneAndCount when reading
// valid unicode without BOM.
func TestReadRuneAndCountNormal(t *testing.T) {

    str := strings.NewReader("pØly")
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // Read 'p'
    code, size, err := reader.readRuneAndCount()

    if code != 'p' { t.Log("wrong code point read first from input string"); t.Fail() }
    if size != 1 { t.Log("wrong code point size read first from input string"); t.Fail() }
    if err != nil { t.Log("unexpected error when reading valid unicode"); t.Fail() }

    if reader.Line != 1 { t.Log("wrong line count after reading one character"); t.Fail() }
    if reader.Column != 1 { t.Log("wrong column count after reading one character"); t.Fail() }
    if reader.Offset != 1 { t.Log("wrong offset count after reading one byte"); t.Fail() }

    // Read 'o' (shows advancement)
    code, size, err = reader.readRuneAndCount()

    if code != 'Ø' { t.Log("wrong code point read second from input string"); t.Fail() }
    if size != 2 { t.Log("wrong code point size read second from input string"); t.Fail() }
    if err != nil { t.Log("unexpected error when reading valid unicode"); t.Fail() }

    if reader.Line != 1 { t.Log("wrong line count after reading two characters"); t.Fail() }
    if reader.Column != 2 { t.Log("wrong column count after reading two characters"); t.Fail() }
    if reader.Offset != 3 { t.Log("wrong offset count after reading three bytes"); t.Fail() }
}

// TestReadRuneAndCountWithValidBOM asserts that readRuneAndCount properly reads
// and advances a valid BOM marker.
func TestReadRuneAndCountWithValidBOM(t *testing.T) {

    str := strings.NewReader("\uFEFFpoly")
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // Read 'BOM'
    code, size, err := reader.readRuneAndCount()

    if code != BOM { t.Log("wrong BOM code point read from input string"); t.Fail() }
    if size != BOMSize { t.Log("wrong BOM code point size read from input string"); t.Fail() }
    if err != nil { t.Log("unexpected error when reading BOM"); t.Fail() }

    if reader.Line != 1 { t.Log("wrong line count after reading BOM"); t.Fail() }
    if reader.Column != 0 { t.Log("wrong column count after reading BOM"); t.Fail() }
    if reader.Offset != BOMSize { t.Log("wrong offset count after reading BOM"); t.Fail() }

    // Read 'p'
    code, size, err = reader.readRuneAndCount()

    if code != 'p' { t.Log("wrong code point read from input string"); t.Fail() }
    if size != 1 { t.Log("wrong code point size read from input string"); t.Fail() }
    if err != nil { t.Log("unexpected error when reading from input string"); t.Fail() }

    if reader.Line != 1 { t.Log("wrong line count when reading after BOM"); t.Fail() }
    if reader.Column != 1 { t.Log("wrong column count when reading after BOM"); t.Fail() }
    if reader.Offset != BOMSize + 1 { t.Log("wrong offset count when reading after BOM"); t.Fail() }
}

// TestReadRuneAndCountWithInvalidBOM asserts that an invalid BOM marker is
// properly detected.
func TestReadRuneAndCountWithInvalidBOM(t *testing.T) {

    input := fmt.Sprintf("%cp%coly", BOM, BOM)

    str := strings.NewReader(input)
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // Read BOM and 'p'
    reader.readRuneAndCount()
    reader.readRuneAndCount()

    // Reading BOM out of position is an error
    _, _, err := reader.readRuneAndCount()

    if err == nil {
        t.Log("expected error caused by invalid BOM")
        t.Fail()
    }
}

// TestReadRuneAndCountWithInvalidUnicode asserts that invalid unicode is
// handled properly.
func TestReadRuneAndCountWithInvalidUnicode(t *testing.T) {

    str := strings.NewReader("\xc3poly")
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // Read invalid code point
    code, size, err := reader.readRuneAndCount()

    if code != unicode.ReplacementChar { t.Log("unexpected code point read from input string"); t.Fail() }
    if size != 1 { t.Log("wrong code point size read from input string"); t.Fail() }
    if err == nil { t.Log("expected an error when reading invalid unicode code point"); t.Fail() }

    if reader.Line != 1 { t.Log("wrong line count after reading invalid unicode code point"); t.Fail() }
    if reader.Column != 1 { t.Log("wrong column count after reading invalid unicode code point"); t.Fail() }
    if reader.Offset != 1 { t.Log("wrong offset count after reading invalid unicode code point"); t.Fail() }

    // Read the second (valid) character to make sure the caret advances
    code, size, err = reader.readRuneAndCount()

    if code != 'p' { t.Log("unexpected code point read from input string"); t.Fail() }
    if size != 1 { t.Log("wrong code point size read from input string"); t.Fail() }
    if err != nil { t.Log("unexpected error when reading valid unicode code point"); t.Fail() }

    if reader.Line != 1 { t.Log("wrong line count after reading valid unicode code point"); t.Fail() }
    if reader.Column != 2 { t.Log("wrong column count after reading valid unicode code point"); t.Fail() }
    if reader.Offset != 2 { t.Log("wrong offset count after reading valid unicode code point"); t.Fail() }
}

// TestReadRuneAndCountWithNewLine asserts that new line characters (\n) are
// read properly and increment line counts.
func TestReadRuneAndCountWithNewLine(t *testing.T) {

    str := strings.NewReader("p\noly")
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // Read
    reader.readRuneAndCount() // p
    reader.readRuneAndCount() // \n

    if reader.Line != 2 { t.Log("wrong line count after line change"); t.Fail() }
    if reader.Column != 0 { t.Log("wrong column count after line change"); t.Fail() }
}

// TestReadRuneAndCountWithError asserts that readRuneAndCount panics on I/O
// errors. In that event, the implementation should call a delegate reporting
// the error.
func TestReadRuneAndCountWithError(t *testing.T) {

    callsDelegate := false

    defer func() {

        // A fatal error is expected
        if r := recover() ; r == nil {
            t.Log("expected panic for I/O error")
            t.Fail()
        }

        // Make sure the delegate was notified
        if !callsDelegate {
            t.Log("expected a delegate report for fatal error")
            t.Fail()
        }
    }()

    // Initialize with a bogus reader
    str := &protoErrorReader{}
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // A fatal error is expected
    reader.FatalError = func(reader *UTF8Reader, err error) {
        callsDelegate = true
    }

    // Read
    reader.readRuneAndCount()

    // If we get here, the panic was never triggered
    t.Log("expected panic for I/O error")
    t.Fail()
}

// TestReadRuneAndCountEOF asserts that EOF is properly detected and an error
// is returned for it.
func TestReadRuneAndCountEOF(t *testing.T) {

    str := strings.NewReader("")    // Empty
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // Read
    _, _, err := reader.readRuneAndCount()

    // Should return an error
    if err != io.EOF {
        t.Log("expected EOF error when reading empty string")
        t.Fail()
    }
}

// TestReadBasic checks that the Read routine reads properly from valid input.
func TestReadBasic(t *testing.T) {

    str := strings.NewReader("poly")
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    readOne := func(expected rune) {

        code, _, _ := reader.Read()

        if code != expected {
            t.Logf("expected unicode code point %U and got %U instead", expected, code)
            t.Fail()
        }
    }

    readOne('p')
    readOne('o')
    readOne('l')
    readOne('y')
}

// TestReadSkipBOM asserts that the Read routine properly skips BOM sequences,
// and that only those not showing at the first position are reported.
func TestReadSkipBOM(t *testing.T) {

    input := fmt.Sprintf("%c%c%cpoly", BOM, BOM, BOM)

    // Input starts with three BOM markers
    str := strings.NewReader(input)
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // Count error reports
    errorCount := 0

    reader.Error = func(reader *UTF8Reader, err error) {
        errorCount++
    }

    // Should skip BOM 3 times and read p
    code, _, _ := reader.Read()

    if code != 'p' {
        t.Log("unexpected character read after BOM sequence")
        t.Fail()
    }

    // Three BOM markers, but the first one is valid
    if errorCount != 2 {
        t.Log("unexpected error report count")
        t.Fail()
    }
}

// TestReadSkipInvalid asserts that the Read routine properly skips invalid
// unicode sequences and that it properly reports them as errors. Ignoring
// invalid input, the reader should return EOF.
func TestReadSkipInvalidAndEmpty(t *testing.T) {

    // The input consists of two invalid octets
    str := strings.NewReader("\xf0\x8c")
    buf := bufio.NewReader(str)

    reader := NewUTF8Reader(buf)

    // Count error reports
    errorCount := 0

    reader.Error = func(reader *UTF8Reader, err error) {
        errorCount++
    }

    // Should skip all input and return EOF
    code, _, err := reader.Read()

    if err != io.EOF {
        t.Logf("unexpected input continuation %U read after invalid unicode", code)
        t.Fail()
    }

    // 2 octects are invalid
    if errorCount != 2 {
        t.Log("unexpected error report count")
        t.Fail()
    }
}
