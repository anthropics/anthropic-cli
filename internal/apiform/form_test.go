package apiform

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
	"testing"
)

// Define test cases
var tests = map[string]struct {
	value    any
	format   FormFormat
	expected string
}{
	"nil": {
		value:    nil,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\n\r\n--xxx--\r\n",
	},
	"string": {
		value:    "hello",
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nhello\r\n--xxx--\r\n",
	},
	"int": {
		value:    42,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\n42\r\n--xxx--\r\n",
	},
	"float": {
		value:    3.14,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\n3.14\r\n--xxx--\r\n",
	},
	"bool": {
		value:    true,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\ntrue\r\n--xxx--\r\n",
	},
	"empty slice": {
		value:    []string{},
		expected: "\r\n--xxx--\r\n",
	},
	"nil slice": {
		value:    []string(nil),
		expected: "\r\n--xxx--\r\n",
	},
	"slice with dot indices": {
		value:    []string{"a", "b", "c"},
		format:   FormatIndicesDots,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo.0\"\r\n\r\na\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo.1\"\r\n\r\nb\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo.2\"\r\n\r\nc\r\n--xxx--\r\n",
	},
	"slice with bracket indices": {
		value:    []int{10, 20, 30},
		format:   FormatIndicesBrackets,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo[0]\"\r\n\r\n10\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo[1]\"\r\n\r\n20\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo[2]\"\r\n\r\n30\r\n--xxx--\r\n",
	},
	"slice with repeat": {
		value:    []int{10, 20, 30},
		format:   FormatRepeat,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\n10\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\n20\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\n30\r\n--xxx--\r\n",
	},
	"slice with commas": {
		value:    []int{10, 20, 30},
		format:   FormatComma,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\n10,20,30\r\n--xxx--\r\n",
	},
	"empty map": {
		value:    map[string]any{},
		expected: "\r\n--xxx--\r\n",
	},
	"nil map": {
		value:    map[string]any(nil),
		expected: "\r\n--xxx--\r\n",
	},
	"map": {
		value:    map[string]any{"key1": "value1", "key2": "value2"},
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo.key1\"\r\n\r\nvalue1\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo.key2\"\r\n\r\nvalue2\r\n--xxx--\r\n",
	},
	"nested_map": {
		value:    map[string]any{"outer": map[string]int{"inner1": 10, "inner2": 20}},
		format:   FormatIndicesDots,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo.outer.inner1\"\r\n\r\n10\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo.outer.inner2\"\r\n\r\n20\r\n--xxx--\r\n",
	},
	"mixed_map": {
		value:    map[string]any{"name": "John", "ages": []int{25, 30, 35}},
		format:   FormatIndicesDots,
		expected: "--xxx\r\nContent-Disposition: form-data; name=\"foo.ages.0\"\r\n\r\n25\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo.ages.1\"\r\n\r\n30\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo.ages.2\"\r\n\r\n35\r\n--xxx\r\nContent-Disposition: form-data; name=\"foo.name\"\r\n\r\nJohn\r\n--xxx--\r\n",
	},
}

func TestEncode(t *testing.T) {
	t.Parallel()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer(nil)
			writer := multipart.NewWriter(buf)
			writer.SetBoundary("xxx")

			form := map[string]any{"foo": test.value}
			err := MarshalWithSettings(form, writer, test.format)
			if err != nil {
				t.Errorf("serialization of %v failed with error %v", test.value, err)
			}
			err = writer.Close()
			if err != nil {
				t.Errorf("serialization of %v failed with error %v", test.value, err)
			}
			result := buf.String()
			if result != test.expected {
				t.Errorf("expected %+#v to serialize to:\n\t%q\nbut got:\n\t%q", test.value, test.expected, result)
			}
		})
	}
}

// namedReader wraps an io.Reader with a Name() method to simulate os.File.
type namedReader struct {
	io.Reader
	name string
}

func (r *namedReader) Name() string { return r.name }

func TestEncodeFileUpload(t *testing.T) {
	t.Parallel()

	t.Run("single file uses filename in Content-Disposition", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(buf)
		writer.SetBoundary("xxx")

		file := &namedReader{Reader: strings.NewReader("file content"), name: "test.txt"}
		form := map[string]any{"file": file}
		if err := MarshalWithSettings(form, writer, FormatBrackets); err != nil {
			t.Fatal(err)
		}
		writer.Close()
		result := buf.String()
		if !strings.Contains(result, `name="file"`) {
			t.Errorf("expected field name=\"file\", got:\n%s", result)
		}
		if !strings.Contains(result, `filename="test.txt"`) {
			t.Errorf("expected filename=\"test.txt\", got:\n%s", result)
		}
		if !strings.Contains(result, "file content") {
			t.Errorf("expected file content in body, got:\n%s", result)
		}
	})

	t.Run("array of files uses brackets notation", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(buf)
		writer.SetBoundary("xxx")

		files := []any{
			&namedReader{Reader: strings.NewReader("content A"), name: "skill-dir/SKILL.md"},
			&namedReader{Reader: strings.NewReader("content B"), name: "skill-dir/ref/doc.md"},
		}
		form := map[string]any{"files": files}
		if err := MarshalWithSettings(form, writer, FormatBrackets); err != nil {
			t.Fatal(err)
		}
		writer.Close()
		result := buf.String()
		// Verify files[] bracket notation for each file
		if strings.Count(result, `name="files[]"`) != 2 {
			t.Errorf("expected 2 fields with name=\"files[]\", got:\n%s", result)
		}
		if !strings.Contains(result, `filename="skill-dir/SKILL.md"`) {
			t.Errorf("expected filename with relative path preserved, got:\n%s", result)
		}
		if !strings.Contains(result, `filename="skill-dir/ref/doc.md"`) {
			t.Errorf("expected filename with relative path preserved, got:\n%s", result)
		}
	})

	t.Run("absolute path stripped to basename", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(buf)
		writer.SetBoundary("xxx")

		file := &namedReader{Reader: strings.NewReader("hello"), name: "/tmp/upload.txt"}
		form := map[string]any{"file": file}
		if err := MarshalWithSettings(form, writer, FormatBrackets); err != nil {
			t.Fatal(err)
		}
		writer.Close()
		result := buf.String()
		if !strings.Contains(result, `filename="upload.txt"`) {
			t.Errorf("expected absolute path stripped to basename, got:\n%s", result)
		}
		if strings.Contains(result, `filename="/tmp/upload.txt"`) {
			t.Errorf("absolute path should not appear in filename, got:\n%s", result)
		}
	})

	t.Run("relative path preserved in filename", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(buf)
		writer.SetBoundary("xxx")

		file := &namedReader{Reader: strings.NewReader("hello"), name: "my-skill/SKILL.md"}
		form := map[string]any{"file": file}
		if err := MarshalWithSettings(form, writer, FormatBrackets); err != nil {
			t.Fatal(err)
		}
		writer.Close()
		result := buf.String()
		if !strings.Contains(result, `filename="my-skill/SKILL.md"`) {
			t.Errorf("expected relative path preserved, got:\n%s", result)
		}
	})
}
