package is_test

import (
	"strings"
	"testing"

	"github.com/prodadidb/go-validation"
	"github.com/prodadidb/go-validation/is"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	tests := []struct {
		tag            string
		rule           validation.Rule
		valid, invalid string
		err            string
	}{
		{"Email", is.Email, "test@example.com", "example.com", "must be a valid email address"},
		{"EmailFormat", is.EmailFormat, "test@example.com", "example.com", "must be a valid email address"},
		{"URL", is.URL, "http://example.com", "examplecom", "must be a valid URL"},
		{"RequestURL", is.RequestURL, "http://example.com", "examplecom", "must be a valid request URL"},
		{"RequestURI", is.RequestURI, "http://example.com", "examplecom", "must be a valid request URI"},
		{"Alpha", is.Alpha, "abcd", "ab12", "must contain English letters only"},
		{"Digit", is.Digit, "123", "12ab", "must contain digits only"},
		{"Alphanumeric", is.Alphanumeric, "abc123", "abc.123", "must contain English letters and digits only"},
		{"UTFLetter", is.UTFLetter, "ａｂｃ", "１２３", "must contain unicode letter characters only"},
		{"UTFDigit", is.UTFDigit, "１２３", "ａｂｃ", "must contain unicode decimal digits only"},
		{"UTFNumeric", is.UTFNumeric, "１２３", "ａｂｃ.１２３", "must contain unicode number characters only"},
		{"UTFLetterNumeric", is.UTFLetterNumeric, "ａｂｃ１２３", "ａｂｃ.１２３", "must contain unicode letters and numbers only"},
		{"LowerCase", is.LowerCase, "ａｂc", "Aｂｃ", "must be in lower case"},
		{"UpperCase", is.UpperCase, "ABC", "ABｃ", "must be in upper case"},
		{"IP", is.IP, "74.125.19.99", "74.125.19.999", "must be a valid IP address"},
		{"IPv4", is.IPv4, "74.125.19.99", "2001:4860:0:2001::68", "must be a valid IPv4 address"},
		{"IPv6", is.IPv6, "2001:4860:0:2001::68", "74.125.19.99", "must be a valid IPv6 address"},
		{"MAC", is.MAC, "0123.4567.89ab", "74.125.19.99", "must be a valid MAC address"},
		{"Subdomain", is.Subdomain, "example-subdomain", "example.com", "must be a valid subdomain"},
		{"Domain", is.Domain, "example-domain.com", "localhost", "must be a valid domain"},
		{"Domain", is.Domain, "example-domain.com", strings.Repeat("a", 256), "must be a valid domain"},
		{"DNSName", is.DNSName, "example.com", "abc%", "must be a valid DNS name"},
		{"Host", is.Host, "example.com", "abc%", "must be a valid IP address or DNS name"},
		{"Port", is.Port, "123", "99999", "must be a valid port number"},
		{"Latitude", is.Latitude, "23.123", "100", "must be a valid latitude"},
		{"Longitude", is.Longitude, "123.123", "abc", "must be a valid longitude"},
		{"SSN", is.SSN, "100-00-1000", "100-0001000", "must be a valid social security number"},
		{"Semver", is.Semver, "1.0.0", "1.0.0.0", "must be a valid semantic version"},
		{"ISBN", is.ISBN, "1-61729-085-8", "1-61729-085-81", "must be a valid ISBN"},
		{"ISBN10", is.ISBN10, "1-61729-085-8", "1-61729-085-81", "must be a valid ISBN-10"},
		{"ISBN13", is.ISBN13, "978-4-87311-368-5", "978-4-87311-368-a", "must be a valid ISBN-13"},
		{"UUID", is.UUID, "a987fbc9-4bed-3078-cf07-9141ba07c9f1", "a987fbc9-4bed-3078-cf07-9141ba07c9f3a", "must be a valid UUID"},
		{"UUIDv3", is.UUIDv3, "b987fbc9-4bed-3078-cf07-9141ba07c9f3", "b987fbc9-4bed-4078-cf07-9141ba07c9f3", "must be a valid UUID v3"},
		{"UUIDv4", is.UUIDv4, "57b73598-8764-4ad0-a76a-679bb6640eb1", "b987fbc9-4bed-3078-cf07-9141ba07c9f3", "must be a valid UUID v4"},
		{"UUIDv5", is.UUIDv5, "987fbc97-4bed-5078-af07-9141ba07c9f3", "b987fbc9-4bed-3078-cf07-9141ba07c9f3", "must be a valid UUID v5"},
		{"MongoID", is.MongoID, "507f1f77bcf86cd799439011", "507f1f77bcf86cd79943901", "must be a valid hex-encoded MongoDB ObjectId"},
		{"CreditCard", is.CreditCard, "375556917985515", "375556917985516", "must be a valid credit card number"},
		{"JSON", is.JSON, "[1, 2]", "[1, 2,]", "must be in valid JSON format"},
		{"ASCII", is.ASCII, "abc", "ａabc", "must contain ASCII characters only"},
		{"PrintableASCII", is.PrintableASCII, "abc", "ａabc", "must contain printable ASCII characters only"},
		{"E164", is.E164, "+19251232233", "+00124222333", "must be a valid E164 number"},
		{"CountryCode2", is.CountryCode2, "US", "XY", "must be a valid two-letter country code"},
		{"CountryCode3", is.CountryCode3, "USA", "XYZ", "must be a valid three-letter country code"},
		{"CurrencyCode", is.CurrencyCode, "USD", "USS", "must be valid ISO 4217 currency code"},
		{"DialString", is.DialString, "localhost.local:1", "localhost.loc:100000", "must be a valid dial string"},
		{"DataURI", is.DataURI, "data:image/png;base64,TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4=", "image/gif;base64,U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", "must be a Base64-encoded data URI"},
		{"Base64", is.Base64, "TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4=", "image", "must be encoded in Base64"},
		{"Multibyte", is.Multibyte, "ａｂｃ", "abc", "must contain multibyte characters"},
		{"FullWidth", is.FullWidth, "３ー０", "abc", "must contain full-width characters"},
		{"HalfWidth", is.HalfWidth, "abc123い", "００１１", "must contain half-width characters"},
		{"VariableWidth", is.VariableWidth, "３ー０123", "abc", "must contain both full-width and half-width characters"},
		{"Hexadecimal", is.Hexadecimal, "FEF", "FTF", "must be a valid hexadecimal number"},
		{"HexColor", is.HexColor, "F00", "FTF", "must be a valid hexadecimal color code"},
		{"RGBColor", is.RGBColor, "rgb(100, 200, 1)", "abc", "must be a valid RGB color code"},
		{"Int", is.Int, "100", "1.1", "must be an integer number"},
		{"Float", is.Float, "1.1", "a.1", "must be a floating point number"},
		{"VariableWidth", is.VariableWidth, "", "", ""},
	}

	for _, test := range tests {
		err := test.rule.Validate("")
		assert.Nil(t, err, test.tag)
		err = test.rule.Validate(test.valid)
		assert.Nil(t, err, test.tag)
		err = test.rule.Validate(&test.valid)
		assert.Nil(t, err, test.tag)
		err = test.rule.Validate(test.invalid)
		assertError(t, test.err, err, test.tag)
		err = test.rule.Validate(&test.invalid)
		assertError(t, test.err, err, test.tag)
	}
}

func assertError(t *testing.T, expected string, err error, tag string) {
	if expected == "" {
		assert.Nil(t, err, tag)
	} else if assert.NotNil(t, err, tag) {
		assert.Equal(t, expected, err.Error(), tag)
	}
}
