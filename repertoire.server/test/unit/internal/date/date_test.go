package date

import (
	"bytes"
	"testing"
	"time"

	"repertoire/server/internal/date"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScan_WhenValueIsTimeTime_ShouldSetCalendarDateInUTC(t *testing.T) {
	var d date.Date
	input := time.Date(2024, time.March, 14, 10, 30, 0, 0, time.UTC)

	err := d.Scan(input)

	require.NoError(t, err)
	assert.True(t, time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC).Equal(time.Time(d)))
}

func TestScan_WhenValueIsTimeTimeInNonUTCLocation_ShouldPreserveLocalCalendarDate(t *testing.T) {
	// 23:30 in a UTC-5 zone is 04:30 UTC the *next* day. Scan should keep
	// the calendar date as observed in the value's own location, not the
	// date you'd get by naively converting to UTC first.
	loc := time.FixedZone("UTC-5", -5*60*60)
	input := time.Date(2024, time.March, 14, 23, 30, 0, 0, loc)
	var d date.Date

	err := d.Scan(input)

	require.NoError(t, err)
	assert.True(t, time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC).Equal(time.Time(d)))
}

func TestScan_WhenValueIsByteSlice_ShouldParseAsDate(t *testing.T) {
	var d date.Date

	err := d.Scan([]byte("2024-03-14"))

	require.NoError(t, err)
	assert.True(t, time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC).Equal(time.Time(d)))
}

func TestScan_WhenValueIsString_ShouldParseAsDate(t *testing.T) {
	var d date.Date

	err := d.Scan("2024-03-14")

	require.NoError(t, err)
	assert.True(t, time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC).Equal(time.Time(d)))
}

func TestScan_WhenValueIsMalformedString_ShouldReturnError(t *testing.T) {
	var d date.Date

	err := d.Scan("not-a-date")

	assert.Error(t, err)
}

func TestScan_WhenValueIsNil_ShouldLeaveDateUnchanged(t *testing.T) {
	d := date.Date(time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC))

	err := d.Scan(nil)

	require.NoError(t, err)
	assert.True(t, time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC).Equal(time.Time(d)))
}

func TestScan_WhenValueIsUnsupportedType_ShouldReturnError(t *testing.T) {
	var d date.Date

	err := d.Scan(12345)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type for Date")
}

func TestValue_WhenDateIsZero_ShouldReturnNil(t *testing.T) {
	var d date.Date

	value, err := d.Value()

	require.NoError(t, err)
	assert.Nil(t, value)
}

func TestValue_WhenDateIsNotZero_ShouldReturnUTCMidnight(t *testing.T) {
	d := date.Date(time.Date(2024, time.March, 14, 15, 45, 0, 0, time.UTC))

	value, err := d.Value()

	require.NoError(t, err)
	assert.Equal(t, time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC), value)
}

func TestGormDataType_ShouldReturnDate(t *testing.T) {
	var d date.Date

	assert.Equal(t, "date", d.GormDataType())
}

func TestMarshalJSON_WhenDateIsZero_ShouldReturnNull(t *testing.T) {
	var d date.Date

	b, err := d.MarshalJSON()

	require.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestMarshalJSON_WhenDateIsNotZero_ShouldReturnISO8601String(t *testing.T) {
	d := date.Date(time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC))

	b, err := d.MarshalJSON()

	require.NoError(t, err)
	assert.Equal(t, `"2024-03-14"`, string(b))
}

func TestUnmarshalJSON_WhenDataIsNull_ShouldSetZeroDate(t *testing.T) {
	d := date.Date(time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC))

	err := d.UnmarshalJSON([]byte("null"))

	require.NoError(t, err)
	assert.True(t, time.Time(d).IsZero())
}

func TestUnmarshalJSON_WhenDataIsValidDateString_ShouldParseDate(t *testing.T) {
	var d date.Date

	err := d.UnmarshalJSON([]byte(`"2024-03-14"`))

	require.NoError(t, err)
	assert.True(t, time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC).Equal(time.Time(d)))
}

func TestUnmarshalJSON_WhenDataIsInvalidDateFormat_ShouldReturnError(t *testing.T) {
	var d date.Date

	err := d.UnmarshalJSON([]byte(`"not-a-date"`))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date format")
}

func TestUnmarshalJSON_WhenDataIsNotAJSONString_ShouldReturnError(t *testing.T) {
	var d date.Date

	err := d.UnmarshalJSON([]byte(`2024-03-14`)) // unquoted: not a valid JSON string

	assert.Error(t, err)
}

func TestString_ShouldReturnISO8601FormattedDate(t *testing.T) {
	d := date.Date(time.Date(2024, time.March, 14, 8, 0, 0, 0, time.UTC))

	assert.Equal(t, "2024-03-14", d.String())
}

func TestGobEncode_ShouldMatchUnderlyingTimeEncoding(t *testing.T) {
	underlying := time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC)
	d := date.Date(underlying)

	got, err := d.GobEncode()

	require.NoError(t, err)
	want, err := underlying.GobEncode()
	require.NoError(t, err)
	assert.True(t, bytes.Equal(want, got))
}

func TestGobDecode_ShouldDecodeBytesProducedByTimeGobEncode(t *testing.T) {
	underlying := time.Date(2024, time.March, 14, 0, 0, 0, 0, time.UTC)
	encoded, err := underlying.GobEncode()
	require.NoError(t, err)
	var d date.Date

	err = d.GobDecode(encoded)

	require.NoError(t, err)
	assert.True(t, underlying.Equal(time.Time(d)))
}
