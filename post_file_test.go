package blogstore

import (
	"strconv"
	"testing"
)

// TestNewPostFileDefaults tests that NewPostFile() returns a PostFile with:
// - a non-empty ID,
// - empty string fields (PostID, Name, URL, Type, Extension),
// - Size set to "0",
// - Sequence set to 0,
// - non-empty CreatedAt and UpdatedAt,
// - SoftDeletedAt set to MAX_DATETIME (not soft-deleted).
func TestNewPostFileDefaults(t *testing.T) {
	f := NewPostFile()

	if f == nil {
		t.Fatalf("NewPostFile() returned nil")
	}

	if f.GetID() == "" {
		t.Errorf("NewPostFile() must set a non-empty ID")
	}

	if got := f.GetPostID(); got != "" {
		t.Errorf("NewPostFile() PostID = %q, want empty", got)
	}

	if got := f.GetName(); got != "" {
		t.Errorf("NewPostFile() Name = %q, want empty", got)
	}

	if got := f.GetURL(); got != "" {
		t.Errorf("NewPostFile() URL = %q, want empty", got)
	}

	if got := f.GetType(); got != "" {
		t.Errorf("NewPostFile() Type = %q, want empty", got)
	}

	if got := f.GetSize(); got != "0" {
		t.Errorf("NewPostFile() Size = %q, want %q", got, "0")
	}

	if got := f.GetExtension(); got != "" {
		t.Errorf("NewPostFile() Extension = %q, want empty", got)
	}

	if got := f.GetSequence(); got != 0 {
		t.Errorf("NewPostFile() Sequence = %d, want 0", got)
	}

	if got := f.GetCreatedAt(); got == "" {
		t.Errorf("NewPostFile() CreatedAt must not be empty")
	}

	if got := f.GetUpdatedAt(); got == "" {
		t.Errorf("NewPostFile() UpdatedAt must not be empty")
	}

	if got := f.GetSoftDeletedAt(); got != MAX_DATETIME {
		t.Errorf("NewPostFile() SoftDeletedAt = %q, want %q", got, MAX_DATETIME)
	}

	if f.IsSoftDeleted() {
		t.Errorf("NewPostFile() IsSoftDeleted() = true, want false")
	}
}

// TestPostFileGettersAndSetters tests all getter/setter pairs and method chaining.
func TestPostFileGettersAndSetters(t *testing.T) {
	f := NewPostFile()

	// Test method chaining returns the same instance
	result := f.SetID("test-id").
		SetPostID("post-123").
		SetName("photo.jpg").
		SetURL("https://example.com/photo.jpg").
		SetType("image/jpeg").
		SetSize("512000").
		SetExtension("jpg").
		SetSequence(5)

	if result != f {
		t.Errorf("method chaining should return the same instance")
	}

	// Verify all values
	if got := f.GetID(); got != "test-id" {
		t.Errorf("GetID() = %q, want %q", got, "test-id")
	}
	if got := f.GetPostID(); got != "post-123" {
		t.Errorf("GetPostID() = %q, want %q", got, "post-123")
	}
	if got := f.GetName(); got != "photo.jpg" {
		t.Errorf("GetName() = %q, want %q", got, "photo.jpg")
	}
	if got := f.GetURL(); got != "https://example.com/photo.jpg" {
		t.Errorf("GetURL() = %q, want %q", got, "https://example.com/photo.jpg")
	}
	if got := f.GetType(); got != "image/jpeg" {
		t.Errorf("GetType() = %q, want %q", got, "image/jpeg")
	}
	if got := f.GetSize(); got != "512000" {
		t.Errorf("GetSize() = %q, want %q", got, "512000")
	}
	if got := f.GetExtension(); got != "jpg" {
		t.Errorf("GetExtension() = %q, want %q", got, "jpg")
	}
	if got := f.GetSequence(); got != 5 {
		t.Errorf("GetSequence() = %d, want %d", got, 5)
	}
}

// TestPostFileTimestamps tests setting and getting timestamps, including carbon accessors.
func TestPostFileTimestamps(t *testing.T) {
	f := NewPostFile()

	// Set timestamps
	f.SetCreatedAt("2025-01-15 10:30:00")
	f.SetUpdatedAt("2025-06-20 14:45:00")
	f.SetSoftDeletedAt("2025-07-01 09:00:00")

	if got := f.GetCreatedAt(); got != "2025-01-15 10:30:00" {
		t.Errorf("GetCreatedAt() = %q, want %q", got, "2025-01-15 10:30:00")
	}
	if got := f.GetUpdatedAt(); got != "2025-06-20 14:45:00" {
		t.Errorf("GetUpdatedAt() = %q, want %q", got, "2025-06-20 14:45:00")
	}
	if got := f.GetSoftDeletedAt(); got != "2025-07-01 09:00:00" {
		t.Errorf("GetSoftDeletedAt() = %q, want %q", got, "2025-07-01 09:00:00")
	}

	// Carbon accessors should not be nil
	if f.GetCreatedAtCarbon() == nil {
		t.Errorf("GetCreatedAtCarbon() should not be nil")
	}
	if f.GetUpdatedAtCarbon() == nil {
		t.Errorf("GetUpdatedAtCarbon() should not be nil")
	}
	if f.GetSoftDeletedAtCarbon() == nil {
		t.Errorf("GetSoftDeletedAtCarbon() should not be nil")
	}

	// After setting soft delete timestamp, IsSoftDeleted should be true
	if !f.IsSoftDeleted() {
		t.Errorf("IsSoftDeleted() = false after setting past SoftDeletedAt, want true")
	}
}

// TestPostFileSetEmptyTimestamps tests that setting empty timestamps is a no-op.
func TestPostFileSetEmptyTimestamps(t *testing.T) {
	f := NewPostFile()

	originalCreatedAt := f.GetCreatedAt()
	originalUpdatedAt := f.GetUpdatedAt()
	originalSoftDeletedAt := f.GetSoftDeletedAt()

	// Setting empty should not change the value
	f.SetCreatedAt("")
	f.SetUpdatedAt("")
	f.SetSoftDeletedAt("")

	if got := f.GetCreatedAt(); got != originalCreatedAt {
		t.Errorf("SetCreatedAt(\"\") should be no-op, got %q, want %q", got, originalCreatedAt)
	}
	if got := f.GetUpdatedAt(); got != originalUpdatedAt {
		t.Errorf("SetUpdatedAt(\"\") should be no-op, got %q, want %q", got, originalUpdatedAt)
	}
	if got := f.GetSoftDeletedAt(); got != originalSoftDeletedAt {
		t.Errorf("SetSoftDeletedAt(\"\") should be no-op, got %q, want %q", got, originalSoftDeletedAt)
	}
}

// TestNewPostFileFromExistingData tests hydrating a PostFile from a data map.
func TestNewPostFileFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "file-abc123",
		COLUMN_POST_ID:         "post-xyz789",
		COLUMN_NAME:            "document.pdf",
		COLUMN_URL:             "https://example.com/doc.pdf",
		COLUMN_FILE_TYPE:       "application/pdf",
		COLUMN_FILE_SIZE:       "2048576",
		COLUMN_FILE_EXTENSION:  "pdf",
		COLUMN_SEQUENCE:        "3",
		COLUMN_CREATED_AT:      "2025-01-10 08:00:00",
		COLUMN_UPDATED_AT:      "2025-03-15 12:00:00",
		COLUMN_SOFT_DELETED_AT: MAX_DATETIME,
	}

	f := NewPostFileFromExistingData(data)

	if got := f.GetID(); got != "file-abc123" {
		t.Errorf("GetID() = %q, want %q", got, "file-abc123")
	}
	if got := f.GetPostID(); got != "post-xyz789" {
		t.Errorf("GetPostID() = %q, want %q", got, "post-xyz789")
	}
	if got := f.GetName(); got != "document.pdf" {
		t.Errorf("GetName() = %q, want %q", got, "document.pdf")
	}
	if got := f.GetURL(); got != "https://example.com/doc.pdf" {
		t.Errorf("GetURL() = %q, want %q", got, "https://example.com/doc.pdf")
	}
	if got := f.GetType(); got != "application/pdf" {
		t.Errorf("GetType() = %q, want %q", got, "application/pdf")
	}
	if got := f.GetSize(); got != "2048576" {
		t.Errorf("GetSize() = %q, want %q", got, "2048576")
	}
	if got := f.GetExtension(); got != "pdf" {
		t.Errorf("GetExtension() = %q, want %q", got, "pdf")
	}
	if got := f.GetSequence(); got != 3 {
		t.Errorf("GetSequence() = %d, want %d", got, 3)
	}
	if got := f.GetCreatedAt(); got != "2025-01-10 08:00:00" {
		t.Errorf("GetCreatedAt() = %q, want %q", got, "2025-01-10 08:00:00")
	}
	if got := f.GetUpdatedAt(); got != "2025-03-15 12:00:00" {
		t.Errorf("GetUpdatedAt() = %q, want %q", got, "2025-03-15 12:00:00")
	}
	if got := f.GetSoftDeletedAt(); got != MAX_DATETIME {
		t.Errorf("GetSoftDeletedAt() = %q, want %q", got, MAX_DATETIME)
	}
	if f.IsSoftDeleted() {
		t.Errorf("IsSoftDeleted() = true, want false for MAX_DATETIME")
	}
}

// TestNewPostFileFromExistingDataPartial tests hydrating with a partial data map.
func TestNewPostFileFromExistingDataPartial(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:      "file-123",
		COLUMN_POST_ID: "post-456",
		COLUMN_NAME:    "partial.jpg",
	}

	f := NewPostFileFromExistingData(data)

	if got := f.GetID(); got != "file-123" {
		t.Errorf("GetID() = %q, want %q", got, "file-123")
	}
	if got := f.GetPostID(); got != "post-456" {
		t.Errorf("GetPostID() = %q, want %q", got, "post-456")
	}
	if got := f.GetName(); got != "partial.jpg" {
		t.Errorf("GetName() = %q, want %q", got, "partial.jpg")
	}
	// Unset fields should be zero values
	if got := f.GetURL(); got != "" {
		t.Errorf("GetURL() = %q, want empty", got)
	}
	if got := f.GetSequence(); got != 0 {
		t.Errorf("GetSequence() = %d, want 0", got)
	}
}

// TestNewPostFileFromExistingDataInvalidSequence tests that invalid sequence value is handled gracefully.
func TestNewPostFileFromExistingDataInvalidSequence(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:       "file-123",
		COLUMN_SEQUENCE: "not-a-number",
	}

	f := NewPostFileFromExistingData(data)

	if got := f.GetSequence(); got != 0 {
		t.Errorf("GetSequence() with invalid value = %d, want 0", got)
	}
}

// TestPostFileGetData tests that GetData() returns all fields as a map.
func TestPostFileGetData(t *testing.T) {
	f := NewPostFile()
	f.SetID("file-data-test").
		SetPostID("post-data-test").
		SetName("test.png").
		SetURL("https://example.com/test.png").
		SetType("image/png").
		SetSize("999").
		SetExtension("png").
		SetSequence(7)

	data := f.GetData()

	if data[COLUMN_ID] != "file-data-test" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_ID, data[COLUMN_ID], "file-data-test")
	}
	if data[COLUMN_POST_ID] != "post-data-test" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_POST_ID, data[COLUMN_POST_ID], "post-data-test")
	}
	if data[COLUMN_NAME] != "test.png" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_NAME, data[COLUMN_NAME], "test.png")
	}
	if data[COLUMN_URL] != "https://example.com/test.png" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_URL, data[COLUMN_URL], "https://example.com/test.png")
	}
	if data[COLUMN_FILE_TYPE] != "image/png" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_FILE_TYPE, data[COLUMN_FILE_TYPE], "image/png")
	}
	if data[COLUMN_FILE_SIZE] != "999" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_FILE_SIZE, data[COLUMN_FILE_SIZE], "999")
	}
	if data[COLUMN_FILE_EXTENSION] != "png" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_FILE_EXTENSION, data[COLUMN_FILE_EXTENSION], "png")
	}
	if data[COLUMN_SEQUENCE] != "7" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_SEQUENCE, data[COLUMN_SEQUENCE], "7")
	}
	if data[COLUMN_CREATED_AT] == "" {
		t.Errorf("GetData()[%s] should not be empty", COLUMN_CREATED_AT)
	}
	if data[COLUMN_UPDATED_AT] == "" {
		t.Errorf("GetData()[%s] should not be empty", COLUMN_UPDATED_AT)
	}
	if data[COLUMN_SOFT_DELETED_AT] == "" {
		t.Errorf("GetData()[%s] should not be empty", COLUMN_SOFT_DELETED_AT)
	}
}

// TestPostFileGetDataRoundTrip tests that GetData() output can be used to create a new identical PostFile.
func TestPostFileGetDataRoundTrip(t *testing.T) {
	f1 := NewPostFile()
	f1.SetID("roundtrip-id").
		SetPostID("roundtrip-post").
		SetName("roundtrip.jpg").
		SetURL("https://example.com/roundtrip.jpg").
		SetType("image/jpeg").
		SetSize("12345").
		SetExtension("jpg").
		SetSequence(2)

	data := f1.GetData()
	f2 := NewPostFileFromExistingData(data)

	if f2.GetID() != f1.GetID() {
		t.Errorf("roundtrip GetID() = %q, want %q", f2.GetID(), f1.GetID())
	}
	if f2.GetPostID() != f1.GetPostID() {
		t.Errorf("roundtrip GetPostID() = %q, want %q", f2.GetPostID(), f1.GetPostID())
	}
	if f2.GetName() != f1.GetName() {
		t.Errorf("roundtrip GetName() = %q, want %q", f2.GetName(), f1.GetName())
	}
	if f2.GetURL() != f1.GetURL() {
		t.Errorf("roundtrip GetURL() = %q, want %q", f2.GetURL(), f1.GetURL())
	}
	if f2.GetType() != f1.GetType() {
		t.Errorf("roundtrip GetType() = %q, want %q", f2.GetType(), f1.GetType())
	}
	if f2.GetSize() != f1.GetSize() {
		t.Errorf("roundtrip GetSize() = %q, want %q", f2.GetSize(), f1.GetSize())
	}
	if f2.GetExtension() != f1.GetExtension() {
		t.Errorf("roundtrip GetExtension() = %q, want %q", f2.GetExtension(), f1.GetExtension())
	}
	if f2.GetSequence() != f1.GetSequence() {
		t.Errorf("roundtrip GetSequence() = %d, want %d", f2.GetSequence(), f1.GetSequence())
	}
}

// TestPostFileIsSoftDeleted tests the IsSoftDeleted() method with various soft delete timestamps.
func TestPostFileIsSoftDeleted(t *testing.T) {
	f := NewPostFile()

	// Default: MAX_DATETIME means not soft deleted
	if f.IsSoftDeleted() {
		t.Errorf("IsSoftDeleted() with MAX_DATETIME = true, want false")
	}

	// Set to a past date - should be soft deleted
	f.SetSoftDeletedAt("2020-01-01 00:00:00")
	if !f.IsSoftDeleted() {
		t.Errorf("IsSoftDeleted() with past date = false, want true")
	}

	// Set back to MAX_DATETIME - should not be soft deleted
	f.SetSoftDeletedAt(MAX_DATETIME)
	if f.IsSoftDeleted() {
		t.Errorf("IsSoftDeleted() after reset to MAX_DATETIME = true, want false")
	}
}

// TestPostFileSequence tests sequence edge cases.
func TestPostFileSequence(t *testing.T) {
	f := NewPostFile()

	// Default 0
	if got := f.GetSequence(); got != 0 {
		t.Errorf("default GetSequence() = %d, want 0", got)
	}

	// Negative
	f.SetSequence(-1)
	if got := f.GetSequence(); got != -1 {
		t.Errorf("GetSequence() = %d, want %d", got, -1)
	}

	// Large value
	f.SetSequence(999999)
	if got := f.GetSequence(); got != 999999 {
		t.Errorf("GetSequence() = %d, want %d", got, 999999)
	}
}

// TestPostFileGetDataSequenceString tests that GetData() returns sequence as a string.
func TestPostFileGetDataSequenceString(t *testing.T) {
	f := NewPostFile()
	f.SetSequence(42)

	data := f.GetData()
	if data[COLUMN_SEQUENCE] != strconv.Itoa(42) {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_SEQUENCE, data[COLUMN_SEQUENCE], strconv.Itoa(42))
	}
}
