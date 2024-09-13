package ssl

import (
	"reflect"
	"testing"
)

func TestExtractOpenSSHVersion(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Multiple version",
			input:    `Openssh_8.4.1p1, OpenSSH_8.4.1p1, OpenSSL 3.3.1 4 Jun 2024`,
			expected: []string{"8.4", "8.4"},
		},
		{
			name:     "version not match",
			input:    "This string has no OpenSSH version.",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			versions, err := ExtractOpenSSHVersion(tc.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(versions, tc.expected) {
				t.Errorf("Expected versions %v, got %v", tc.expected, versions)
			}
		})
	}
}

func TestExtractOpenSSHVersionWithMajorAndMinor(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *SSHVersion
	}{
		{
			name:  "Mixed case",
			input: `Openssh_8.4.1p1, OpenSSH_8.4.1p1, OpenSSL 3.3.1 4 Jun 2024`,
			expected: &SSHVersion{
				Major: 8, Minor: 4,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			versions, err := ExtractOpenSSHVersionWithMajorAndMinor(tc.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(versions, tc.expected) {
				t.Errorf("Expected versions %v, got %v", tc.expected, versions)
			}
		})
	}
}
