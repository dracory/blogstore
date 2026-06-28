package blogstore

import (
	"strconv"
	"testing"
)

// TestNewMedia tests that NewMedia() initializes all fields with correct defaults.
func TestNewMedia(t *testing.T) {
	m := NewMedia()

	if m.GetID() == "" {
		t.Error("NewMedia() ID should not be empty")
	}
	if m.GetEntityID() != "" {
		t.Errorf("NewMedia() EntityID = %q, want empty", m.GetEntityID())
	}
	if m.GetTitle() != "" {
		t.Errorf("NewMedia() Title = %q, want empty", m.GetTitle())
	}
	if m.GetDescription() != "" {
		t.Errorf("NewMedia() Description = %q, want empty", m.GetDescription())
	}
	if m.GetMemo() != "" {
		t.Errorf("NewMedia() Memo = %q, want empty", m.GetMemo())
	}
	if m.GetURL() != "" {
		t.Errorf("NewMedia() URL = %q, want empty", m.GetURL())
	}
	if m.GetType() != "" {
		t.Errorf("NewMedia() Type = %q, want empty", m.GetType())
	}
	if m.GetSize() != "0" {
		t.Errorf("NewMedia() Size = %q, want %q", m.GetSize(), "0")
	}
	if m.GetExtension() != "" {
		t.Errorf("NewMedia() Extension = %q, want empty", m.GetExtension())
	}
	if m.GetSequence() != 0 {
		t.Errorf("NewMedia() Sequence = %d, want 0", m.GetSequence())
	}
	if m.GetStatus() != MEDIA_STATUS_DRAFT {
		t.Errorf("NewMedia() Status = %q, want %q", m.GetStatus(), MEDIA_STATUS_DRAFT)
	}
	if m.GetCreatedAt() == "" {
		t.Error("NewMedia() CreatedAt should not be empty")
	}
	if m.GetUpdatedAt() == "" {
		t.Error("NewMedia() UpdatedAt should not be empty")
	}
	if m.GetSoftDeletedAt() != MAX_DATETIME {
		t.Errorf("NewMedia() SoftDeletedAt = %q, want %q", m.GetSoftDeletedAt(), MAX_DATETIME)
	}
	if m.IsSoftDeleted() {
		t.Error("NewMedia() IsSoftDeleted() = true, want false")
	}

	metas, err := m.GetMetas()
	if err != nil {
		t.Fatalf("NewMedia() GetMetas() error = %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("NewMedia() Metas len = %d, want 0", len(metas))
	}
}

// TestMediaSettersAndGetters tests all setter/getter pairs.
func TestMediaSettersAndGetters(t *testing.T) {
	m := NewMedia()

	m.SetID("media-123").
		SetEntityID("entity-456").
		SetTitle("Test Image").
		SetDescription("A test image").
		SetMemo("internal note").
		SetURL("https://example.com/image.jpg").
		SetType("image/jpeg").
		SetSize("1024").
		SetExtension("jpg").
		SetSequence(5).
		SetStatus(MEDIA_STATUS_ACTIVE)

	if got := m.GetID(); got != "media-123" {
		t.Errorf("GetID() = %q, want %q", got, "media-123")
	}
	if got := m.GetEntityID(); got != "entity-456" {
		t.Errorf("GetEntityID() = %q, want %q", got, "entity-456")
	}
	if got := m.GetTitle(); got != "Test Image" {
		t.Errorf("GetTitle() = %q, want %q", got, "Test Image")
	}
	if got := m.GetDescription(); got != "A test image" {
		t.Errorf("GetDescription() = %q, want %q", got, "A test image")
	}
	if got := m.GetMemo(); got != "internal note" {
		t.Errorf("GetMemo() = %q, want %q", got, "internal note")
	}
	if got := m.GetURL(); got != "https://example.com/image.jpg" {
		t.Errorf("GetURL() = %q, want %q", got, "https://example.com/image.jpg")
	}
	if got := m.GetType(); got != "image/jpeg" {
		t.Errorf("GetType() = %q, want %q", got, "image/jpeg")
	}
	if got := m.GetSize(); got != "1024" {
		t.Errorf("GetSize() = %q, want %q", got, "1024")
	}
	if got := m.GetExtension(); got != "jpg" {
		t.Errorf("GetExtension() = %q, want %q", got, "jpg")
	}
	if got := m.GetSequence(); got != 5 {
		t.Errorf("GetSequence() = %d, want %d", got, 5)
	}
	if got := m.GetStatus(); got != MEDIA_STATUS_ACTIVE {
		t.Errorf("GetStatus() = %q, want %q", got, MEDIA_STATUS_ACTIVE)
	}
}

// TestMediaTimestamps tests timestamp setters and getters.
func TestMediaTimestamps(t *testing.T) {
	m := NewMedia()

	m.SetCreatedAt("2025-01-10 08:00:00")
	m.SetUpdatedAt("2025-03-15 12:00:00")
	m.SetSoftDeletedAt(MAX_DATETIME)

	if got := m.GetCreatedAt(); got != "2025-01-10 08:00:00" {
		t.Errorf("GetCreatedAt() = %q, want %q", got, "2025-01-10 08:00:00")
	}
	if got := m.GetUpdatedAt(); got != "2025-03-15 12:00:00" {
		t.Errorf("GetUpdatedAt() = %q, want %q", got, "2025-03-15 12:00:00")
	}
	if got := m.GetSoftDeletedAt(); got != MAX_DATETIME {
		t.Errorf("GetSoftDeletedAt() = %q, want %q", got, MAX_DATETIME)
	}

	// Test Carbon accessors
	if m.GetCreatedAtCarbon() == nil {
		t.Error("GetCreatedAtCarbon() should not be nil")
	}
	if m.GetUpdatedAtCarbon() == nil {
		t.Error("GetUpdatedAtCarbon() should not be nil")
	}
	if m.GetSoftDeletedAtCarbon() == nil {
		t.Error("GetSoftDeletedAtCarbon() should not be nil")
	}
}

// TestMediaGetData tests that GetData() returns all fields as a map.
func TestMediaGetData(t *testing.T) {
	m := NewMedia()
	m.SetID("media-data-test").
		SetEntityID("entity-data-test").
		SetTitle("test.png").
		SetDescription("desc").
		SetMemo("memo").
		SetURL("https://example.com/test.png").
		SetType("image/png").
		SetSize("999").
		SetExtension("png").
		SetSequence(7).
		SetStatus(MEDIA_STATUS_ACTIVE)

	data := m.GetData()

	if data[COLUMN_ID] != "media-data-test" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_ID, data[COLUMN_ID], "media-data-test")
	}
	if data[COLUMN_ENTITY_ID] != "entity-data-test" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_ENTITY_ID, data[COLUMN_ENTITY_ID], "entity-data-test")
	}
	if data[COLUMN_TITLE] != "test.png" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_TITLE, data[COLUMN_TITLE], "test.png")
	}
	if data[COLUMN_MEDIA_URL] != "https://example.com/test.png" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_MEDIA_URL, data[COLUMN_MEDIA_URL], "https://example.com/test.png")
	}
	if data[COLUMN_MEDIA_TYPE] != "image/png" {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_MEDIA_TYPE, data[COLUMN_MEDIA_TYPE], "image/png")
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
	if data[COLUMN_STATUS] != MEDIA_STATUS_ACTIVE {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_STATUS, data[COLUMN_STATUS], MEDIA_STATUS_ACTIVE)
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

// TestMediaIsSoftDeleted tests the IsSoftDeleted() method.
func TestMediaIsSoftDeleted(t *testing.T) {
	m := NewMedia()

	if m.IsSoftDeleted() {
		t.Error("IsSoftDeleted() with MAX_DATETIME = true, want false")
	}

	m.SetSoftDeletedAt("2020-01-01 00:00:00")
	if !m.IsSoftDeleted() {
		t.Error("IsSoftDeleted() with past date = false, want true")
	}

	m.SetSoftDeletedAt(MAX_DATETIME)
	if m.IsSoftDeleted() {
		t.Error("IsSoftDeleted() after reset to MAX_DATETIME = true, want false")
	}
}

// TestMediaSequence tests sequence edge cases.
func TestMediaSequence(t *testing.T) {
	m := NewMedia()

	if got := m.GetSequence(); got != 0 {
		t.Errorf("default GetSequence() = %d, want 0", got)
	}

	m.SetSequence(-1)
	if got := m.GetSequence(); got != -1 {
		t.Errorf("GetSequence() = %d, want %d", got, -1)
	}

	m.SetSequence(999999)
	if got := m.GetSequence(); got != 999999 {
		t.Errorf("GetSequence() = %d, want %d", got, 999999)
	}
}

// TestMediaGetDataSequenceString tests that GetData() returns sequence as a string.
func TestMediaGetDataSequenceString(t *testing.T) {
	m := NewMedia()
	m.SetSequence(42)

	data := m.GetData()
	if data[COLUMN_SEQUENCE] != strconv.Itoa(42) {
		t.Errorf("GetData()[%s] = %q, want %q", COLUMN_SEQUENCE, data[COLUMN_SEQUENCE], strconv.Itoa(42))
	}
}

// TestMediaStatusPredicates tests status predicate methods.
func TestMediaStatusPredicates(t *testing.T) {
	m := NewMedia()

	m.SetStatus(MEDIA_STATUS_DRAFT)
	if !m.IsDraft() {
		t.Error("IsDraft() should be true for draft status")
	}
	if m.IsActive() {
		t.Error("IsActive() should be false for draft status")
	}

	m.SetStatus(MEDIA_STATUS_ACTIVE)
	if !m.IsActive() {
		t.Error("IsActive() should be true for active status")
	}
	if m.IsDraft() {
		t.Error("IsDraft() should be false for active status")
	}

	m.SetStatus(MEDIA_STATUS_INACTIVE)
	if !m.IsInactive() {
		t.Error("IsInactive() should be true for inactive status")
	}
}

// TestMediaTypePredicates tests type predicate methods.
func TestMediaTypePredicates(t *testing.T) {
	m := NewMedia()

	m.SetType("image/jpeg")
	if !m.IsImage() {
		t.Error("IsImage() should be true for image/jpeg")
	}
	if m.IsVideo() {
		t.Error("IsVideo() should be false for image/jpeg")
	}

	m.SetType("video/mp4")
	if !m.IsVideo() {
		t.Error("IsVideo() should be true for video/mp4")
	}
	if m.IsImage() {
		t.Error("IsImage() should be false for video/mp4")
	}

	m.SetType("application/pdf")
	if m.IsImage() {
		t.Error("IsImage() should be false for application/pdf")
	}
	if m.IsVideo() {
		t.Error("IsVideo() should be false for application/pdf")
	}
}

// TestMediaMetas tests metadata operations.
func TestMediaMetas(t *testing.T) {
	m := NewMedia()

	// SetMetas and GetMetas
	metaMap := map[string]string{"k1": "v1"}
	if err := m.SetMetas(metaMap); err != nil {
		t.Fatalf("SetMetas() error = %v", err)
	}

	metas, err := m.GetMetas()
	if err != nil {
		t.Fatalf("GetMetas() error = %v", err)
	}
	if metas["k1"] != "v1" {
		t.Errorf("metas[k1] = %q, want %q", metas["k1"], "v1")
	}

	// SetMeta
	if err := m.SetMeta("editor", "alice"); err != nil {
		t.Fatalf("SetMeta() error = %v", err)
	}
	if m.GetMeta("editor") != "alice" {
		t.Errorf("GetMeta('editor') = %q, want %q", m.GetMeta("editor"), "alice")
	}

	// MetasUpsert
	if err := m.MetasUpsert(map[string]string{"k2": "v2"}); err != nil {
		t.Fatalf("MetasUpsert() error = %v", err)
	}
	metas, _ = m.GetMetas()
	if metas["k2"] != "v2" {
		t.Errorf("metas[k2] = %q, want %q", metas["k2"], "v2")
	}
	if metas["k1"] != "v1" {
		t.Errorf("metas[k1] should still be %q after upsert", "v1")
	}

	// MetaRemove
	if err := m.MetaRemove("k1"); err != nil {
		t.Fatalf("MetaRemove() error = %v", err)
	}
	metas, _ = m.GetMetas()
	if _, exists := metas["k1"]; exists {
		t.Error("metas[k1] should not exist after removal")
	}

	// MetasRemove
	if err := m.MetasRemove([]string{"k2", "editor"}); err != nil {
		t.Fatalf("MetasRemove() error = %v", err)
	}
	metas, _ = m.GetMetas()
	if len(metas) != 0 {
		t.Errorf("metas should be empty after MetasRemove, len = %d", len(metas))
	}
}
