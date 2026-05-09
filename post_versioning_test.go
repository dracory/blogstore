package blogstore

import (
	"encoding/json"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestPostMarshalToVersioning(t *testing.T) {
	post := NewPost()
	post.SetTitle("Test Post").
		SetContent("Test Content").
		SetSummary("Test Summary").
		SetStatus(POST_STATUS_PUBLISHED).
		SetFeatured(YES)

	// Test marshaling
	result, err := post.MarshalToVersioning()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == "" {
		t.Fatalf("Expected non-empty result")
	}

	// Parse the result to verify structure
	var versionedData map[string]string
	err = json.Unmarshal([]byte(result), &versionedData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that timestamp fields are excluded
	if _, ok := versionedData[COLUMN_CREATED_AT]; ok {
		t.Errorf("Expected %s to not be in versioned data", COLUMN_CREATED_AT)
	}
	if _, ok := versionedData[COLUMN_UPDATED_AT]; ok {
		t.Errorf("Expected %s to not be in versioned data", COLUMN_UPDATED_AT)
	}
	if _, ok := versionedData[COLUMN_SOFT_DELETED_AT]; ok {
		t.Errorf("Expected %s to not be in versioned data", COLUMN_SOFT_DELETED_AT)
	}

	// Verify that important fields are included
	if versionedData["title"] != "Test Post" {
		t.Errorf("Expected title to be 'Test Post', got %q", versionedData["title"])
	}
	if versionedData["content"] != "Test Content" {
		t.Errorf("Expected content to be 'Test Content', got %q", versionedData["content"])
	}
	if versionedData["summary"] != "Test Summary" {
		t.Errorf("Expected summary to be 'Test Summary', got %q", versionedData["summary"])
	}
	if versionedData["status"] != POST_STATUS_PUBLISHED {
		t.Errorf("Expected status to be %s, got %s", POST_STATUS_PUBLISHED, versionedData["status"])
	}
	if versionedData["featured"] != YES {
		t.Errorf("Expected featured to be %s, got %s", YES, versionedData["featured"])
	}
}

func TestPostUnmarshalFromVersioning(t *testing.T) {
	// Create original post
	originalPost := NewPost()
	originalPost.SetTitle("Original Title").
		SetContent("Original Content").
		SetSummary("Original Summary").
		SetStatus(POST_STATUS_DRAFT).
		SetFeatured(NO)

	// Marshal to versioning format
	versionedData, err := originalPost.MarshalToVersioning()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Create new post and unmarshal from versioning
	newPost := NewPost()
	err = newPost.UnmarshalFromVersioning(versionedData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify data was restored correctly
	title := newPost.GetTitle()
	if title != "Original Title" {
		t.Errorf("Expected title to be 'Original Title', got %q", title)
	}

	content := newPost.GetContent()
	if content != "Original Content" {
		t.Errorf("Expected content to be 'Original Content', got %q", content)
	}

	summary := newPost.GetSummary()
	if summary != "Original Summary" {
		t.Errorf("Expected summary to be 'Original Summary', got %q", summary)
	}

	status := newPost.GetStatus()
	if status != POST_STATUS_DRAFT {
		t.Errorf("Expected status to be %s, got %s", POST_STATUS_DRAFT, status)
	}

	featured := newPost.GetFeatured()
	if featured != NO {
		t.Errorf("Expected featured to be %s, got %s", NO, featured)
	}

	// Verify that updated_at was set to current time
	updatedAt := newPost.GetUpdatedAt()

	// Parse the updated_at timestamp
	parsedTime := carbon.Parse(updatedAt)

	// Should be very recent (within 1 second)
	now := carbon.Now(carbon.UTC)
	if now.DiffInSeconds(parsedTime) > 1 {
		t.Errorf("Expected updated_at to be within 1 second of now")
	}
}

func TestPostUnmarshalFromVersioningWithInvalidTimestamps(t *testing.T) {
	// Create versioned data with invalid timestamp format
	versionedData := map[string]string{
		"title":      "Test Title",
		"content":    "Test Content",
		"created_at": "2026-01-15 08:08:36 +0000 +0000", // Invalid format
		"updated_at": "2026-01-15 08:08:36 +0000 +0000", // Invalid format
		"status":     POST_STATUS_PUBLISHED,
		"featured":   YES,
	}

	data, err := json.Marshal(versionedData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Create post and unmarshal
	post := NewPost()
	err = post.UnmarshalFromVersioning(string(data))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify data was restored (except timestamps)
	title := post.GetTitle()
	if title != "Test Title" {
		t.Errorf("Expected title to be 'Test Title', got %q", title)
	}

	status := post.GetStatus()
	if status != POST_STATUS_PUBLISHED {
		t.Errorf("Expected status to be %s, got %s", POST_STATUS_PUBLISHED, status)
	}

	// Verify that updated_at was set to current time (not the invalid timestamp)
	updatedAt := post.GetUpdatedAt()

	parsedTime := carbon.Parse(updatedAt)

	now := carbon.Now(carbon.UTC)
	if now.DiffInSeconds(parsedTime) > 1 {
		t.Errorf("Expected updated_at to be within 1 second of now")
	}
}

func TestPostUnmarshalFromVersioningEmptyData(t *testing.T) {
	post := NewPost()

	// Test with empty string
	err := post.UnmarshalFromVersioning("")
	if err == nil {
		t.Errorf("Expected error for empty input")
	}

	// Test with invalid JSON
	err = post.UnmarshalFromVersioning("{invalid json}")
	if err == nil {
		t.Errorf("Expected error for invalid JSON")
	}
}

func TestPostVersioningRoundTrip(t *testing.T) {
	// Create original post with all fields
	originalPost := NewPost()
	originalPost.SetTitle("Round Trip Test").
		SetContent("This is test content for round trip testing").
		SetSummary("Test summary").
		SetStatus(POST_STATUS_PUBLISHED).
		SetFeatured(YES).
		SetMetaDescription("Meta description").
		SetMetaKeywords("test,keywords").
		SetMetaRobots("index,follow").
		SetCanonicalURL("https://example.com/test").
		SetImageUrl("https://example.com/image.jpg").
		SetMemo("Test memo")

	// Add some metas
	err := originalPost.AddMetas(map[string]string{
		"custom_field1": "value1",
		"custom_field2": "value2",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Marshal to versioning
	versionedData, err := originalPost.MarshalToVersioning()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Create new post and restore from versioning
	newPost := NewPost()
	err = newPost.UnmarshalFromVersioning(versionedData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify all fields were restored correctly
	title := newPost.GetTitle()
	if title != "Round Trip Test" {
		t.Errorf("Expected title to be 'Round Trip Test', got %q", title)
	}

	content := newPost.GetContent()
	if content != "This is test content for round trip testing" {
		t.Errorf("Expected content to be 'This is test content for round trip testing', got %q", content)
	}

	summary := newPost.GetSummary()
	if summary != "Test summary" {
		t.Errorf("Expected summary to be 'Test summary', got %q", summary)
	}

	status := newPost.GetStatus()
	if status != POST_STATUS_PUBLISHED {
		t.Errorf("Expected status to be %s, got %s", POST_STATUS_PUBLISHED, status)
	}

	featured := newPost.GetFeatured()
	if featured != YES {
		t.Errorf("Expected featured to be %s, got %s", YES, featured)
	}

	metaDesc := newPost.GetMetaDescription()
	if metaDesc != "Meta description" {
		t.Errorf("Expected meta description to be 'Meta description', got %q", metaDesc)
	}

	metaKeywords := newPost.GetMetaKeywords()
	if metaKeywords != "test,keywords" {
		t.Errorf("Expected meta keywords to be 'test,keywords', got %q", metaKeywords)
	}

	metaRobots := newPost.GetMetaRobots()
	if metaRobots != "index,follow" {
		t.Errorf("Expected meta robots to be 'index,follow', got %q", metaRobots)
	}

	canonicalURL := newPost.GetCanonicalURL()
	if canonicalURL != "https://example.com/test" {
		t.Errorf("Expected canonical URL to be 'https://example.com/test', got %q", canonicalURL)
	}

	imageUrl := newPost.GetImageUrl()
	if imageUrl != "https://example.com/image.jpg" {
		t.Errorf("Expected image URL to be 'https://example.com/image.jpg', got %q", imageUrl)
	}

	memo := newPost.GetMemo()
	if memo != "Test memo" {
		t.Errorf("Expected memo to be 'Test memo', got %q", memo)
	}

	// Verify metas
	metas, err := newPost.GetMetas()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if metas["custom_field1"] != "value1" {
		t.Errorf("Expected custom_field1 to be 'value1', got %q", metas["custom_field1"])
	}
	if metas["custom_field2"] != "value2" {
		t.Errorf("Expected custom_field2 to be 'value2', got %q", metas["custom_field2"])
	}
}

func TestPostVersioningExcludesTimestamps(t *testing.T) {
	post := NewPost()

	// Set specific timestamps
	testTime := carbon.Parse("2026-01-15 10:00:00")
	post.SetCreatedAt(testTime.ToDateTimeString(carbon.UTC))
	post.SetUpdatedAt(testTime.ToDateTimeString(carbon.UTC))

	// Add other data
	post.SetTitle("Timestamp Test")

	// Marshal to versioning
	versionedData, err := post.MarshalToVersioning()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Parse and verify timestamps are excluded
	var data map[string]string
	err = json.Unmarshal([]byte(versionedData), &data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, ok := data[COLUMN_CREATED_AT]; ok {
		t.Errorf("Expected %s to not be in data", COLUMN_CREATED_AT)
	}
	if _, ok := data[COLUMN_UPDATED_AT]; ok {
		t.Errorf("Expected %s to not be in data", COLUMN_UPDATED_AT)
	}
	if _, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		t.Errorf("Expected %s to not be in data", COLUMN_SOFT_DELETED_AT)
	}
	if _, ok := data["title"]; !ok {
		t.Errorf("Expected title to be in data")
	}
	if data["title"] != "Timestamp Test" {
		t.Errorf("Expected title to be 'Timestamp Test', got %q", data["title"])
	}
}

func BenchmarkPostMarshalToVersioning(b *testing.B) {
	post := NewPost()
	post.SetTitle("Benchmark Test").
		SetContent("This is benchmark content with enough length to be realistic").
		SetSummary("Benchmark summary").
		SetStatus(POST_STATUS_PUBLISHED).
		SetFeatured(YES)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := post.MarshalToVersioning()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPostUnmarshalFromVersioning(b *testing.B) {
	post := NewPost()
	post.SetTitle("Benchmark Test").
		SetContent("This is benchmark content with enough length to be realistic").
		SetSummary("Benchmark summary").
		SetStatus(POST_STATUS_PUBLISHED).
		SetFeatured(YES)

	versionedData, err := post.MarshalToVersioning()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newPost := NewPost()
		err := newPost.UnmarshalFromVersioning(versionedData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
