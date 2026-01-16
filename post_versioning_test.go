package blogstore

import (
	"encoding/json"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Parse the result to verify structure
	var versionedData map[string]string
	err = json.Unmarshal([]byte(result), &versionedData)
	require.NoError(t, err)

	// Verify that timestamp fields are excluded
	assert.NotContains(t, versionedData, COLUMN_CREATED_AT)
	assert.NotContains(t, versionedData, COLUMN_UPDATED_AT)
	assert.NotContains(t, versionedData, COLUMN_SOFT_DELETED_AT)

	// Verify that important fields are included
	assert.Equal(t, "Test Post", versionedData["title"])
	assert.Equal(t, "Test Content", versionedData["content"])
	assert.Equal(t, "Test Summary", versionedData["summary"])
	assert.Equal(t, POST_STATUS_PUBLISHED, versionedData["status"])
	assert.Equal(t, YES, versionedData["featured"])
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
	require.NoError(t, err)

	// Create new post and unmarshal from versioning
	newPost := NewPost()
	err = newPost.UnmarshalFromVersioning(versionedData)
	require.NoError(t, err)

	// Verify data was restored correctly
	title := newPost.Title()
	assert.Equal(t, "Original Title", title)

	content := newPost.Content()
	assert.Equal(t, "Original Content", content)

	summary := newPost.Summary()
	assert.Equal(t, "Original Summary", summary)

	status := newPost.Status()
	assert.Equal(t, POST_STATUS_DRAFT, status)

	featured := newPost.Featured()
	assert.Equal(t, NO, featured)

	// Verify that updated_at was set to current time
	updatedAt := newPost.UpdatedAt()

	// Parse the updated_at timestamp
	parsedTime := carbon.Parse(updatedAt)

	// Should be very recent (within 1 second)
	now := carbon.Now(carbon.UTC)
	assert.True(t, now.DiffInSeconds(parsedTime) <= 1)
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
	require.NoError(t, err)

	// Create post and unmarshal
	post := NewPost()
	err = post.UnmarshalFromVersioning(string(data))
	require.NoError(t, err)

	// Verify data was restored (except timestamps)
	title := post.Title()
	assert.Equal(t, "Test Title", title)

	status := post.Status()
	assert.Equal(t, POST_STATUS_PUBLISHED, status)

	// Verify that updated_at was set to current time (not the invalid timestamp)
	updatedAt := post.UpdatedAt()

	parsedTime := carbon.Parse(updatedAt)

	now := carbon.Now(carbon.UTC)
	assert.True(t, now.DiffInSeconds(parsedTime) <= 1)
}

func TestPostUnmarshalFromVersioningEmptyData(t *testing.T) {
	post := NewPost()

	// Test with empty string
	err := post.UnmarshalFromVersioning("")
	assert.Error(t, err)

	// Test with invalid JSON
	err = post.UnmarshalFromVersioning("{invalid json}")
	assert.Error(t, err)
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
	require.NoError(t, err)

	// Marshal to versioning
	versionedData, err := originalPost.MarshalToVersioning()
	require.NoError(t, err)

	// Create new post and restore from versioning
	newPost := NewPost()
	err = newPost.UnmarshalFromVersioning(versionedData)
	require.NoError(t, err)

	// Verify all fields were restored correctly
	title := newPost.Title()
	assert.Equal(t, "Round Trip Test", title)

	content := newPost.Content()
	assert.Equal(t, "This is test content for round trip testing", content)

	summary := newPost.Summary()
	assert.Equal(t, "Test summary", summary)

	status := newPost.Status()
	assert.Equal(t, POST_STATUS_PUBLISHED, status)

	featured := newPost.Featured()
	assert.Equal(t, YES, featured)

	metaDesc := newPost.MetaDescription()
	assert.Equal(t, "Meta description", metaDesc)

	metaKeywords := newPost.MetaKeywords()
	assert.Equal(t, "test,keywords", metaKeywords)

	metaRobots := newPost.MetaRobots()
	assert.Equal(t, "index,follow", metaRobots)

	canonicalURL := newPost.CanonicalURL()
	assert.Equal(t, "https://example.com/test", canonicalURL)

	imageUrl := newPost.ImageUrl()
	assert.Equal(t, "https://example.com/image.jpg", imageUrl)

	memo := newPost.Memo()
	assert.Equal(t, "Test memo", memo)

	// Verify metas
	metas, err := newPost.Metas()
	require.NoError(t, err)
	assert.Equal(t, "value1", metas["custom_field1"])
	assert.Equal(t, "value2", metas["custom_field2"])
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
	require.NoError(t, err)

	// Parse and verify timestamps are excluded
	var data map[string]string
	err = json.Unmarshal([]byte(versionedData), &data)
	require.NoError(t, err)

	assert.NotContains(t, data, COLUMN_CREATED_AT)
	assert.NotContains(t, data, COLUMN_UPDATED_AT)
	assert.NotContains(t, data, COLUMN_SOFT_DELETED_AT)
	assert.Contains(t, data, "title")
	assert.Equal(t, "Timestamp Test", data["title"])
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
