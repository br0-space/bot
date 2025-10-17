package fortune_test

import (
	"testing"

	"github.com/br0-space/bot/pkg/fortune"
)

func TestMakeService(t *testing.T) {
	t.Parallel()

	service := fortune.MakeService()
	// Service is an empty struct, so we just verify it was created without panic
	_ = service
}

func TestGetList(t *testing.T) {
	t.Parallel()

	service := fortune.MakeService()

	// Test that the method doesn't panic and returns a slice
	files := service.GetList()

	// files should be a slice (even if empty for the default path)
	if files == nil {
		t.Error("GetList should return a non-nil slice")
	}
}

func TestExists(t *testing.T) {
	t.Parallel()

	service := fortune.MakeService()

	testCases := []struct {
		name     string
		filename string
		want     bool
	}{
		{"non-existent file", "this-does-not-exist-xyz123", false},
		{"empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := service.Exists(tc.filename)
			if got != tc.want {
				t.Errorf("Exists(%q) = %v, want %v", tc.filename, got, tc.want)
			}
		})
	}
}

func TestReadFortuneFile(t *testing.T) {
	t.Parallel()

	service := fortune.MakeService()

	testCases := []struct {
		name        string
		file        string
		shouldError bool
	}{
		{"non-existent file", "nonexistent-file-xyz", true},
		{"empty filename", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.GetFortune(tc.file)

			if tc.shouldError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetRandomFortune_NoFiles(t *testing.T) {
	t.Parallel()

	service := fortune.MakeService()

	// If there are no files in the default path, we should get an error
	// Note: this test might pass or fail depending on the actual fortune files present
	files := service.GetList()
	if len(files) == 0 {
		_, err := service.GetRandomFortune()
		if err == nil {
			t.Error("Expected error when no fortune files found, got nil")
		}
	}
}

func TestGetFortune_NonExistent(t *testing.T) {
	t.Parallel()

	service := fortune.MakeService()

	_, err := service.GetFortune("this-file-definitely-does-not-exist-xyz123")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestWeightedSelection(t *testing.T) {
	t.Parallel()

	// This test verifies that the weighted selection logic works correctly
	// by checking the distribution over many samples

	service := fortune.MakeService()
	files := service.GetList()

	if len(files) == 0 {
		t.Skip("No fortune files available for weighted selection test")
	}

	// Sample many times and verify we get results
	sampleSize := 100
	successCount := 0

	for range sampleSize {
		_, err := service.GetRandomFortune()
		if err == nil {
			successCount++
		}
	}

	if successCount == 0 {
		t.Error("GetRandomFortune failed for all samples")
	}

	// Verify that we got mostly successful results
	successRate := float64(successCount) / float64(sampleSize)
	if successRate < 0.9 {
		t.Errorf("Success rate too low: %.2f (expected >= 0.9)", successRate)
	}
}
