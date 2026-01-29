package middleware

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

// TestUploadConfig 测试上传配置
func TestUploadConfig(t *testing.T) {
	config := &UploadConfig{
		MaxSize:      10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif"},
		AllowedExts:  []string{".jpg", ".jpeg", ".png", ".gif"},
	}

	if config.MaxSize != 10*1024*1024 {
		t.Errorf("Expected MaxSize to be 10MB, got %d", config.MaxSize)
	}

	if len(config.AllowedTypes) != 3 {
		t.Errorf("Expected 3 allowed types, got %d", len(config.AllowedTypes))
	}

	if len(config.AllowedExts) != 4 {
		t.Errorf("Expected 4 allowed extensions, got %d", len(config.AllowedExts))
	}
}

// TestValidateUpload 测试文件上传验证
func TestValidateUpload(t *testing.T) {
	config := &UploadConfig{
		MaxSize:      10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif"},
		AllowedExts:  []string{".jpg", ".jpeg", ".png", ".gif"},
	}

	tests := []struct {
		name       string
		filename   string
		content    []byte
		maxSize    int64
		wantErr    bool
		errMessage string
	}{
		{
			name:     "valid JPEG file",
			filename: "test.jpg",
			content:  []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}, // JPEG header
			maxSize:  10 * 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "valid PNG file",
			filename: "test.png",
			content:  []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG header
			maxSize:  10 * 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "file too large",
			filename: "large.jpg",
			content:  bytes.Repeat([]byte{0xFF}, 11*1024*1024), // 11MB
			maxSize:  10 * 1024 * 1024,
			wantErr:  true,
		},
		{
			name:     "invalid extension",
			filename: "test.exe",
			content:  []byte{0xFF, 0xD8, 0xFF, 0xE0},
			maxSize:  10 * 1024 * 1024,
			wantErr:  true,
		},
		{
			name:     "no extension",
			filename: "test",
			content:  []byte{0xFF, 0xD8, 0xFF, 0xE0},
			maxSize:  10 * 1024 * 1024,
			wantErr:  true,
		},
		{
			name:     "empty file",
			filename: "empty.jpg",
			content:  []byte{},
			maxSize:  10 * 1024 * 1024,
			wantErr:  true,
		},
		{
			name:     "extension mismatch - exe claiming to be jpg",
			filename: "malicious.exe.jpg",
			content:  []byte{0x4D, 0x5A}, // PE header (executable)
			maxSize:  10 * 1024 * 1024,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file header
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", tt.filename)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}

			_, err = part.Write(tt.content)
			if err != nil {
				t.Fatalf("Failed to write content: %v", err)
			}

			writer.Close()

			// Create request
			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Parse multipart form
			err = req.ParseMultipartForm(tt.maxSize + 1)
			if err != nil {
				t.Fatalf("Failed to parse multipart form: %v", err)
			}

			file, header, err := req.FormFile("file")
			if err != nil {
				t.Fatalf("Failed to get form file: %v", err)
			}
			defer file.Close()

			// Test validation
			testConfig := &UploadConfig{
				MaxSize:      tt.maxSize,
				AllowedTypes: config.AllowedTypes,
				AllowedExts:  config.AllowedExts,
			}

			err = ValidateUpload(header, testConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateFileExtension 测试文件扩展名验证
func TestValidateFileExtension(t *testing.T) {
	config := &UploadConfig{
		AllowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"},
	}

	tests := []struct {
		filename string
		want     bool
	}{
		{"test.jpg", true},
		{"test.jpeg", true},
		{"test.png", true},
		{"test.gif", true},
		{"test.pdf", true},
		{"test.JPG", true},  // Case insensitive
		{"test.Png", true},  // Case insensitive
		{"test.exe", false},
		{"test.txt", false},
		{"test", false},
		{"", false},
		{".hidden", false},
		{"test.tar.gz", false}, // No support for double extensions
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := ValidateFileExtension(tt.filename, config.AllowedExts)
			if got != tt.want {
				t.Errorf("ValidateFileExtension(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// TestValidateFileType 测试文件类型验证（MIME类型）
func TestValidateFileType(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		wantType string
		wantErr  bool
	}{
		{
			name:     "JPEG file",
			content:  []byte{0xFF, 0xD8, 0xFF, 0xE0},
			wantType: "image/jpeg",
			wantErr:  false,
		},
		{
			name:     "PNG file",
			content:  []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			wantType: "image/png",
			wantErr:  false,
		},
		{
			name:     "GIF file",
			content:  []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61},
			wantType: "image/gif",
			wantErr:  false,
		},
		{
			name:     "PDF file",
			content:  []byte{0x25, 0x50, 0x44, 0x46},
			wantType: "application/pdf",
			wantErr:  false,
		},
		{
			name:     "empty file",
			content:  []byte{},
			wantType: "",
			wantErr:  true,
		},
		{
			name:     "unknown file type",
			content:  []byte{0x00, 0x01, 0x02, 0x03},
			wantType: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileType, err := DetectFileType(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFileType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && fileType != tt.wantType {
				t.Errorf("DetectFileType() = %v, want %v", fileType, tt.wantType)
			}
		})
	}
}

// TestValidateFileSize 测试文件大小验证
func TestValidateFileSize(t *testing.T) {
	config := &UploadConfig{
		MaxSize: 10 * 1024 * 1024, // 10MB
	}

	tests := []struct {
		name    string
		size    int64
		wantErr bool
	}{
		{"1KB file", 1024, false},
		{"1MB file", 1024 * 1024, false},
		{"5MB file", 5 * 1024 * 1024, false},
		{"10MB file", 10 * 1024 * 1024, false},
		{"10MB + 1 byte", 10*1024*1024 + 1, true},
		{"11MB file", 11 * 1024 * 1024, true},
		{"100MB file", 100 * 1024 * 1024, true},
		{"empty file", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileSize(tt.size, config.MaxSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileSize(%d) error = %v, wantErr %v", tt.size, err, tt.wantErr)
			}
		})
	}
}

// TestGenerateSafeFilename 测试安全文件名生成
func TestGenerateSafeFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "normal filename",
			filename: "test.jpg",
			want:     "test.jpg",
		},
		{
			name:     "filename with path",
			filename: "../../../etc/passwd.jpg",
			want:     "passwd.jpg",
		},
		{
			name:     "filename with special chars",
			filename: "test@#$%file.jpg",
			want:     "testfile.jpg",
		},
		{
			name:     "filename with spaces",
			filename: "my file name.jpg",
			want:     "my_file_name.jpg",
		},
		{
			name:     "filename with unicode",
			filename: "文件名.jpg",
			want:     ".jpg",
		},
		{
			name:     "empty filename",
			filename: "",
			want:     "",
		},
		{
			name:     "filename with dots",
			filename: "test.file.name.jpg",
			want:     "test.file.name.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSafeFilename(tt.filename)
			if got != tt.want {
				t.Errorf("GenerateSafeFilename(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}
