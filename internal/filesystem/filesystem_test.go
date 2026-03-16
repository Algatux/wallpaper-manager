package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListSupportedImages_Empty(t *testing.T) {
	tmpDir := t.TempDir()

	result := ListSupportedImages(tmpDir)

	if result == nil {
		t.Error("ListSupportedImages should return empty slice, not nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 images, got %d", len(result))
	}
}

func TestListSupportedImages_OnlyDirs(t *testing.T) {
	tmpDir := t.TempDir()

	os.Mkdir(filepath.Join(tmpDir, "subdir1"), 0755)
	os.Mkdir(filepath.Join(tmpDir, "subdir2"), 0755)

	result := ListSupportedImages(tmpDir)

	if len(result) != 0 {
		t.Errorf("Expected 0 images when only directories exist, got %d", len(result))
	}
}

func TestListSupportedImages_SupportedFormats(t *testing.T) {
	tmpDir := t.TempDir()

	supportedFormats := []string{".png", ".jpg", ".jpeg", ".bmp", ".pnm", ".tga", ".tiff", ".webp", ".gif"}

	for i, ext := range supportedFormats {
		filename := filepath.Join(tmpDir, "image"+string(rune('0'+i))+ext)
		if err := os.WriteFile(filename, []byte("fake image"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	result := ListSupportedImages(tmpDir)

	if len(result) != len(supportedFormats) {
		t.Errorf("Expected %d images, got %d", len(supportedFormats), len(result))
	}
}

func TestListSupportedImages_UnsupportedFormats(t *testing.T) {
	tmpDir := t.TempDir()

	unsupportedFormats := []string{".txt", ".pdf", ".doc", ".mp4", ".mov"}

	for _, ext := range unsupportedFormats {
		filename := filepath.Join(tmpDir, "file"+ext)
		if err := os.WriteFile(filename, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	result := ListSupportedImages(tmpDir)

	if len(result) != 0 {
		t.Errorf("Expected 0 images with unsupported formats, got %d", len(result))
	}
}

func TestListSupportedImages_Mixed(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "pic1.png"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "pic2.jpg"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "pic3.gif"), []byte(""), 0644)

	os.WriteFile(filepath.Join(tmpDir, "doc.pdf"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "script.sh"), []byte(""), 0644)

	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "nested.png"), []byte(""), 0644)

	result := ListSupportedImages(tmpDir)

	if len(result) != 3 {
		t.Errorf("Expected 3 images, got %d", len(result))
	}

	expectedFiles := map[string]bool{
		"pic1.png": false,
		"pic2.jpg": false,
		"pic3.gif": false,
	}

	for _, imagePath := range result {
		filename := filepath.Base(imagePath)
		if _, exists := expectedFiles[filename]; !exists {
			t.Errorf("Unexpected file in result: %s", filename)
		} else {
			expectedFiles[filename] = true
		}
	}

	for file, found := range expectedFiles {
		if !found {
			t.Errorf("Expected file not found in result: %s", file)
		}
	}
}

func TestListSupportedImages_TrailingSlash(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "image.png"), []byte(""), 0644)

	dirWithSlash := tmpDir + string(os.PathSeparator)
	result := ListSupportedImages(dirWithSlash)

	if len(result) != 1 {
		t.Errorf("Expected 1 image with trailing slash, got %d", len(result))
	}
}

func TestListSupportedImages_NoTrailingSlash(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "image.png"), []byte(""), 0644)

	result := ListSupportedImages(tmpDir)

	if len(result) != 1 {
		t.Errorf("Expected 1 image without trailing slash, got %d", len(result))
	}
}

func TestListSupportedImages_PathCorrectness(t *testing.T) {
	tmpDir := t.TempDir()

	filename := "test_image.png"
	os.WriteFile(filepath.Join(tmpDir, filename), []byte(""), 0644)

	result := ListSupportedImages(tmpDir)

	if len(result) != 1 {
		t.Errorf("Expected 1 image, got %d", len(result))
	}

	returnedPath := result[0]
	expectedPath := filepath.Join(tmpDir, filename)

	if returnedPath != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, returnedPath)
	}
}

func TestListSupportedImages_MultipleImages(t *testing.T) {
	tmpDir := t.TempDir()

	for i := 0; i < 10; i++ {
		filename := filepath.Join(tmpDir, "image"+string(rune('0'+i))+".jpg")
		os.WriteFile(filename, []byte(""), 0644)
	}

	result := ListSupportedImages(tmpDir)

	if len(result) != 10 {
		t.Errorf("Expected 10 images, got %d", len(result))
	}
}

func TestListSupportedImages_SpecialCharactersInFilename(t *testing.T) {
	tmpDir := t.TempDir()

	filenames := []string{
		"image (1).png",
		"photo_2024-03-16.jpg",
		"wallpaper-dark.gif",
	}

	for _, filename := range filenames {
		os.WriteFile(filepath.Join(tmpDir, filename), []byte(""), 0644)
	}

	result := ListSupportedImages(tmpDir)

	if len(result) != 3 {
		t.Errorf("Expected 3 images with special characters, got %d", len(result))
	}
}

func TestListSupportedImages_NonexistentDirectory(t *testing.T) {
	result := ListSupportedImages("/nonexistent/path/that/does/not/exist")

	if result == nil {
		t.Error("ListSupportedImages should return empty slice, not nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 images for non-existent directory, got %d", len(result))
	}
}
