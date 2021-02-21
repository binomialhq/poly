package poly

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

// Position represents the caret's position in some input text.
type Position struct {
	Line   int // Line number (1-indexed)
	Column int // Column number (1-indexed)
	Offset int // Byte offset from the beginning of the input (0-indexed)
}

// String returns a string representation of the Position object.
func (pos *Position) String() string {
	return fmt.Sprintf("%d:%d", pos.Line, pos.Column)
}

// BOM is the unicode Byte Order Mark.
const BOM = '\uFEFF'

// BOMSize indicates BOM's size, which is 3 bytes long.
const BOMSize = 3

// UTF8Reader reads UTF-8 code points in sequence from a buffered reader. The
// implementation ignores BOM markers and invalid unicode code points, reporting
// encoding errors to a settable Error callback. The implementation does not
// attempt to recover from I/O errors and instead panics, after reporting those
// errors to the FatalError callback delegate, if one is set. When reading, the
// reader advances a count of Line, Column, and byte Offset, which help with
// mapping the relative position of the caret within the input text.
//
// The main point of entry is Read(), which returns a rune with a unicode code
// point and the number of bytes that it consumed from the input. The exception
// to this is EOF, which is returned a null code point (zero) with null size,
// and an error set to io.EOF.
type UTF8Reader struct {

	// reader is the input reader.
	reader *bufio.Reader

	// Error is a delegate routine called when an error occurs. When not set,
	// errors are reported to os.Stderr.
	Error func(reader *UTF8Reader, err error)

	// FatalError is a delegate routine called when a fatal error occurs. When
	// not set or if the delegate routine returns, the code panics.
	FatalError func(reader *UTF8Reader, err error)

	// Position indicates the caret's position with respect to the source input.
	Position
}

// NewUTF8Reader instantiates and initializes a new UTF8Reader with the given
// io.Reader as input source.
func NewUTF8Reader(reader *bufio.Reader) *UTF8Reader {

	return &UTF8Reader{
		reader:   reader,
		Position: Position{Line: 1, Column: 0, Offset: 0},
	}
}

// IsLetter checks if a given rune is a letter (underscore or category L).
func IsLetter(code rune) bool {
	return code == '_' || unicode.IsLetter(code)
}

// IsDecimal checks whether a given rune is a decimal digita (0 through 9).
func IsDecimal(code rune) bool {
	return '0' <= code && code <= '9'
}

// IsWhitespace checks whether a given rune is white space.
func IsWhitespace(code rune) bool {
	return unicode.IsSpace(code)
}

// IsReplacementCharacter checks whether a given rune is an unicode replacement
// character.
func IsReplacementCharacter(code rune) bool {
	return code == unicode.ReplacementChar
}

// IsBOM checks whether a given rune is a byte order marker.
func IsBOM(code rune) bool {
	return code == BOM
}

// IsEOF checks if a given error is end of file.
func IsEOF(err error) bool {
	return err == io.EOF
}

// readRuneAndCount reads a rune from the input source and updates the caret's
// position according to what is read. Every time an unicode code point is read,
// the offset is increment by the size of the code point. The routine skips all
// BOM markers found, reporting as an error those that appear out of place (e.g.
// not as the first character in the sequence). Unicode replacement characters
// are always reported, since they flag invalid unicode. Line changes occur when
// a new line (\n) character is found. In any other situation, the column
// counter is incremented by 1, even for code points that represent more or less
// than one characer space. I/O errors result in a fatal error.
func (reader *UTF8Reader) readRuneAndCount() (rune, int, error) {

	// Read
	code, size, err := reader.reader.ReadRune()

	if err != nil {

		// EOF reached
		if IsEOF(err) {
			return 0, 0, err
		}

		// ReadRune only returns read errors, not related with encoding; we
		// can't recover from those
		reader.reportFatalError(err)
	}

	// Increment offset, even for invalid input
	reader.Offset += size

	switch {

	case code == '\n':
		reader.Line++
		reader.Column = 0

	case IsBOM(code):

		// BOM is only an error when out of place (e.g. not the first character
		// in the input sequence).
		if reader.Offset != BOMSize {
			err = fmt.Errorf("invalid BOM at position %d", reader.Offset-size)
		}

	case IsReplacementCharacter(code):
		err = fmt.Errorf("invalid unicode code point %U at position %d", code, reader.Offset-size)

		// It still moves the column forward
		reader.Column++

	default:
		// All others count as one character (even tabs and other special cases)
		reader.Column++
	}

	return code, size, err
}

// Read returns the first non-BOM and non-replacement character from the input
// sequence, reporting errors until a valid code point is found. The second
// return value is the number of bytes in the rune. The error argument is only
// expected to produce EOF errors. The caret's position is updated according
// to the type of input read (see readRuneAndCount).
func (reader *UTF8Reader) Read() (rune, int, error) {

	code, size, err := reader.readRuneAndCount()

	// BOM and unicode replacement characters are skipped until neither is at
	// the top of the caret position. Although this situation is very unlikely,
	// it is also imperative that we don't return either as valid unicode.
	for IsBOM(code) || IsReplacementCharacter(code) {

		// This is never EOF, which is not supposed to be reported
		if err != nil {
			reader.reportError(err)
		}

		// Read next
		code, size, err = reader.readRuneAndCount()
	}

	// This is valid even for EOF
	return code, size, err
}

// reportError calls the Error delegate if one has been set. If not, the error
// is logged to os.Stderr.
func (reader *UTF8Reader) reportError(err error) {

	// Report the error if a delegate is registered
	if reader.Error != nil {
		reader.Error(reader, err)
	} else {

		// Report to os.Stderr
		fmt.Fprintf(os.Stderr, "%s: %v\n", reader.Position.String(), err)
	}
}

// reportFatalError calls the FatalError delegate if one has been set. If not,
// or if the callback returns, the code panics.
func (reader *UTF8Reader) reportFatalError(err error) {

	// Report the error if a delegate is registered
	if reader.FatalError != nil {
		reader.FatalError(reader, err)
	}

	// If FatalError returns, we panic
	panic(err)
}
