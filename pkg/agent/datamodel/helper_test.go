// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package datamodel

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestValidateDNSPrefix(t *testing.T) {
	cases := []struct {
		name        string
		dnsPrefix   string
		expectedErr error
	}{
		{
			"valid DNS prefix",
			"validDnsPrefix",
			nil,
		},
		{
			"empty string",
			"",
			errors.New("DNSPrefix '' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 0)"), //nolint:lll
		},
		{
			"one char",
			"a",
			errors.New("DNSPrefix 'a' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 1)"), //nolint:lll
		},
		{
			"numbers",
			"1234",
			errors.New("DNSPrefix '1234' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 4)"), //nolint:lll
		},
		{
			"too many chars",
			"verylongdnsprefixthatismorethan45characterslong",
			errors.New("DNSPrefix 'verylongdnsprefixthatismorethan45characterslong' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 47)"), //nolint:lll
		},
		{
			"invalid special character",
			"dnswith_special?char",
			errors.New("DNSPrefix 'dnswith_special?char' is invalid. The DNSPrefix must contain between 3 and 45 characters and can contain only letters, numbers, and hyphens.  It must start with a letter and must end with a letter or a number. (length was 20)"), //nolint:lll
		},
		{
			"valid with numbers",
			"myDNS-1234",
			nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateDNSPrefix(c.dnsPrefix)
			if err != nil && c.expectedErr != nil {
				if err.Error() != c.expectedErr.Error() {
					t.Fatalf("expected validateDNSPrefix to return error %s, but instead got %s", c.expectedErr.Error(), err.Error())
				}
			} else {
				if c.expectedErr != nil {
					t.Fatalf("expected validateDNSPrefix to return error %s, but instead got no error", c.expectedErr.Error())
				} else if err != nil {
					t.Fatalf("expected validateDNSPrefix to return no error, but instead got %s", err.Error())
				}
			}
		})
	}
}

func getCSeriesVMCasesForTesting() []struct {
	name     string
	VMSKU    string
	Expected bool
} {
	cases := []struct {
		name     string
		VMSKU    string
		Expected bool
	}{
		{
			"Standard_DC2s",
			"Standard_DC2s",
			true,
		},
		{
			"Standard_DC4s",
			"Standard_DC4s",
			true,
		},
		{
			"Standard_D2_v2",
			"Standard_D2_v2",
			false,
		},
		{
			"gobledygook",
			"gobledygook",
			false,
		},
		{
			"empty string",
			"",
			false,
		},
	}
	return cases
}

/*
GetDCSeriesVMCasesForTesting returns a struct w/ VM SKUs and whether or not we expect them
to be SGX-enabled.
*/
func GetDCSeriesVMCasesForTesting() []struct {
	VMSKU    string
	Expected bool
} {
	cases := []struct {
		VMSKU    string
		Expected bool
	}{
		{
			"Standard_DC2s",
			true,
		},
		{
			"Standard_DC4s",
			true,
		},
		{
			"Standard_NC12",
			false,
		},
		{
			"gobledygook",
			false,
		},
		{
			"",
			false,
		},
	}

	return cases
}

func TestIsSGXEnabledSKU(t *testing.T) {
	cases := getCSeriesVMCasesForTesting()

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ret := IsSgxEnabledSKU(c.VMSKU)
			if ret != c.Expected {
				t.Fatalf("expected IsSgxEnabledSKU(%s) to return %t, but instead got %t", c.VMSKU, c.Expected, ret)
			}
		})
	}
}

func TestGetOrderedEscapedKeyValsString(t *testing.T) {
	alphabetizedString := `"foo=bar", "yes=please"`
	cases := []struct {
		name     string
		input    map[string]string
		expected string
	}{
		{
			name:     "nil input",
			input:    map[string]string{},
			expected: "",
		},
		{
			name: "valid input",
			input: map[string]string{
				"foo": "bar",
				"yes": "please",
			},
			expected: alphabetizedString,
		},
		{
			name: "valid input re-ordered",
			input: map[string]string{
				"yes": "please",
				"foo": "bar",
			},
			expected: alphabetizedString,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ret := GetOrderedEscapedKeyValsString(c.input)
			if ret != c.expected {
				t.Fatalf("expected GetOrderedEscapedKeyValsString(%s) to return %s, but instead got %s", c.input, c.expected, ret)
			}
		})
	}
}

func TestGetStorageAccountType(t *testing.T) {
	validPremiumVMSize := "Standard_DS2_v2"
	validStandardVMSize := "Standard_D2_v2"
	expectedPremiumTier := "Premium_LRS"
	expectedStandardTier := "Standard_LRS"
	invalidVMSize := "D2v2"

	// test premium VMSize returns premium managed disk tier
	premiumTier, err := GetStorageAccountType(validPremiumVMSize)
	if err != nil {
		t.Fatalf("Invalid sizeName: %s", err)
	}

	if premiumTier != expectedPremiumTier {
		t.Fatalf("premium VM did no match premium managed storage tier")
	}

	// test standard VMSize returns standard managed disk tier
	standardTier, err := GetStorageAccountType(validStandardVMSize)
	if err != nil {
		t.Fatalf("Invalid sizeName: %s", err)
	}

	if standardTier != expectedStandardTier {
		t.Fatalf("standard VM did no match standard managed storage tier")
	}

	// test invalid VMSize
	result, err := GetStorageAccountType(invalidVMSize)
	if err == nil {
		t.Errorf("GetStorageAccountType() = (%s, nil), want error", result)
	}
}

func TestSliceIntIsNonEmpty(t *testing.T) {
	cases := []struct {
		name     string
		input    []int
		expected bool
	}{
		{
			name: "valid slice",
			input: []int{
				1, 2, 3,
			},
			expected: true,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: false,
		},
		{
			name:     "nil slice",
			input:    nil,
			expected: false,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ret := SliceIntIsNonEmpty(c.input)
			if ret != c.expected {
				t.Fatalf("expected SliceIntIsNonEmpty(%v) to return %t, but instead got %t", c.input, c.expected, ret)
			}
		})
	}
}

func TestWrapAsVerbatim(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{
			name:     "just a string",
			s:        "foo",
			expected: "',foo,'",
		},
		{
			name:     "empty string",
			s:        "",
			expected: "',,'",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ret := WrapAsVerbatim(test.s)
			if test.expected != ret {
				t.Errorf("expected %s, instead got : %s", test.expected, ret)
			}
		})
	}
}

func TestIndentString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		count int
		want  string
	}{
		{
			name:  "should leave empty string alone",
			input: "",
			count: 4,
			want:  "",
		},
		{
			name:  "should indent single line string 4 spaces",
			input: "foo",
			count: 4,
			want:  "    foo\n",
		},
		{
			name:  "should indent multi-line string 4 spaces",
			input: "foo\nbar",
			count: 4,
			want:  "    foo\n    bar\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := IndentString(test.input, test.count)
			diff := cmp.Diff(test.want, got)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestGetContainerImageNameFromURL(t *testing.T) {
	tests := []struct {
		name        string
		downloadURL string
		expected    string
		expectedErr string
	}{
		{
			name:        "empty URL",
			downloadURL: "",
			expected:    "",
			expectedErr: "container image component URL is not in the expected format: ",
		},
		{
			name:        "valid URL",
			downloadURL: "mcr.microsoft.com/oss/kubernetes/autoscaler/addon-resizer:*",
			expected:    "addon-resizer",
			expectedErr: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := getContainerImageNameFromURL(test.downloadURL)
			if err != nil {
				if test.expectedErr != err.Error() {
					t.Fatalf("expected error %s, instead got %s", test.expectedErr, err.Error())
				}
			} else {
				if test.expected != got {
					t.Fatalf("expected %s, instead got %s", test.expected, got)
				}
			}
		})
	}
}

func TestGetComponentNameFromURL(t *testing.T) {
	tests := []struct {
		name        string
		downloadURL string
		expected    string
		expectedErr string
	}{
		{
			name:        "empty URL",
			downloadURL: "",
			expected:    "",
			expectedErr: "download file image URL is not in the expected format: ",
		},
		{
			name:        "valid URL",
			downloadURL: "https://acs-mirror.azureedge.net/cni-plugins/v*/binaries",
			expected:    "cni-plugins",
			expectedErr: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := getComponentNameFromURL(test.downloadURL)
			if err != nil {
				if test.expectedErr != err.Error() {
					t.Fatalf("expected error %s, instead got %s", test.expectedErr, err.Error())
				}
			} else {
				if test.expected != got {
					t.Fatalf("expected %s, instead got %s", test.expected, got)
				}
			}
		})
	}
}
