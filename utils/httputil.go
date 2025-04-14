package utils

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Range represents a byte range
type Range struct {
	Start int64
	End   int64
}

// ParseRangeHeader parses an HTTP Range header
func ParseRangeHeader(rangeHeader string, fileSize int64) ([]Range, error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return nil, errors.New("invalid range header format")
	}

	rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")
	ranges := strings.Split(rangeHeader, ",")
	parsedRanges := make([]Range, 0, len(ranges))

	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			return nil, errors.New("invalid range format")
		}

		var start, end int64 = 0, fileSize - 1

		// Parse start
		if parts[0] != "" {
			s, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid range start: %w", err)
			}
			start = s
		}

		// Parse end
		if parts[1] != "" {
			e, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid range end: %w", err)
			}
			end = e
		}

		// Validate range
		if start >= fileSize {
			return nil, errors.New("range start exceeds file size")
		}

		if end >= fileSize {
			end = fileSize - 1
		}

		if start > end {
			return nil, errors.New("range start greater than end")
		}

		parsedRanges = append(parsedRanges, Range{Start: start, End: end})
	}

	return parsedRanges, nil
}

// CopyN copies exactly n bytes from src to dst
func CopyN(dst io.Writer, src io.Reader, n int64) (int64, error) {
	// Use a reasonably sized buffer for streaming
	const bufSize = 64 * 1024 // 64KB buffer
	buf := make([]byte, bufSize)

	var totalWritten int64 = 0

	for totalWritten < n {
		// Calculate how much to read
		toRead := bufSize
		if n-totalWritten < int64(bufSize) {
			toRead = int(n - totalWritten)
		}

		// Read a chunk
		read, err := src.Read(buf[:toRead])
		if err != nil && err != io.EOF {
			return totalWritten, err
		}

		if read == 0 {
			break // End of file
		}

		// Write the chunk
		written, err := dst.Write(buf[:read])
		if err != nil {
			return totalWritten, err
		}

		totalWritten += int64(written)

		// If we couldn't write the full chunk, stop
		if written != read {
			return totalWritten, io.ErrShortWrite
		}

		// If we're at EOF, stop
		if err == io.EOF {
			break
		}
	}

	return totalWritten, nil
}
