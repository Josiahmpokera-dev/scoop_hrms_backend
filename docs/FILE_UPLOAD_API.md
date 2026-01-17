# File Upload API Documentation

## Overview

The File Upload API provides local object storage for employee onboarding files, including photos and documents. Files are stored locally in a structured directory system and can be accessed via HTTP URLs.

## Base URL
All API endpoints use the base URL: `{{BASE_URL}}/api/v1/employees/onboarding`

## Storage Configuration

Files are stored in the local filesystem with the following structure:
```
storage/
├── photos/
│   └── EMP001/
│       └── photo_20240115_143045_abc123.jpg
└── documents/
    └── EMP001/
        ├── national_id_20240115_143100_def456.pdf
        └── degree_20240115_143115_ghi789.pdf
```

### Environment Variables

- `STORAGE_PATH` (default: `./storage`): Base path for file storage
- `STORAGE_URL` (default: `/storage`): Base URL for accessing files

### File Access

Files are served statically at: `http://localhost:8080/storage/{category}/{employee_id}/{filename}`

## Authentication
All requests require Bearer token authentication:
```
Authorization: Bearer <access_token>
```

---

## 1. Upload Photo API (Step 1)

### Endpoint
```
POST /api/v1/employees/onboarding/upload-photo
```

### Description
Uploads a photo for employee onboarding (Step 1 - Personal Information). The photo is stored and a URL is returned that can be used in the onboarding step data.

### Request Format
`multipart/form-data`

### Form Fields
- `employee_id` (string, required): Employee ID (e.g., "EMP001")
- `file` (file, required): Photo file (image)

### Supported Image Types
- JPEG/JPG
- PNG
- GIF
- WebP

### File Size Limit
- Maximum: 10MB

### Success Response (200)
```json
{
  "success": true,
  "message": "Photo uploaded successfully",
  "data": {
    "photo_url": "/storage/photos/EMP001/photo_20240115_143045_abc123.jpg",
    "file_name": "employee-photo.jpg",
    "file_size": 245678,
    "mime_type": "image/jpeg",
    "employee_id": "EMP001"
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/employees/onboarding/upload-photo" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "employee_id=EMP001" \
  -F "file=@/path/to/photo.jpg"
```

### Example Request (JavaScript)
```javascript
const formData = new FormData();
formData.append('employee_id', 'EMP001');
formData.append('file', fileInput.files[0]);

const response = await fetch('http://localhost:8080/api/v1/employees/onboarding/upload-photo', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
});

const result = await response.json();
console.log(result.data.photo_url); // Use this URL in Step 1 data
```

### Usage in Step 1
After uploading the photo, use the returned `photo_url` in your Step 1 data:

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "photo_url": "/storage/photos/EMP001/photo_20240115_143045_abc123.jpg",
  ...
}
```

---

## 2. Upload Document API (Step 6)

### Endpoint
```
POST /api/v1/employees/onboarding/upload-document
```

### Description
Uploads a document for employee onboarding (Step 6 - Documents). The document is stored and a URL is returned that can be used in the onboarding step data.

### Request Format
`multipart/form-data`

### Form Fields
- `employee_id` (string, required): Employee ID (e.g., "EMP001")
- `document_type` (string, required): Document type - `identity`, `work_permit`, `education`, `contract`, `tax_statutory`, `other`
- `file` (file, required): Document file
- `description` (string, optional): Document description

### Supported Document Types
- PDF
- DOC/DOCX
- XLS/XLSX

### File Size Limit
- Maximum: 10MB

### Success Response (200)
```json
{
  "success": true,
  "message": "Document uploaded successfully",
  "data": {
    "file_url": "/storage/documents/EMP001/national_id_20240115_143100_def456.pdf",
    "file_name": "national-id.pdf",
    "file_size": 512345,
    "mime_type": "application/pdf",
    "document_type": "identity",
    "description": "National ID card",
    "employee_id": "EMP001",
    "uploaded_by": 1
  }
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/employees/onboarding/upload-document" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "employee_id=EMP001" \
  -F "document_type=identity" \
  -F "description=National ID card" \
  -F "file=@/path/to/national-id.pdf"
```

### Example Request (JavaScript)
```javascript
const formData = new FormData();
formData.append('employee_id', 'EMP001');
formData.append('document_type', 'identity');
formData.append('description', 'National ID card');
formData.append('file', fileInput.files[0]);

const response = await fetch('http://localhost:8080/api/v1/employees/onboarding/upload-document', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
});

const result = await response.json();
console.log(result.data.file_url); // Use this URL in Step 6 data
```

### Usage in Step 6
After uploading documents, use the returned `file_url` in your Step 6 data:

```json
{
  "documents": [
    {
      "document_type": "identity",
      "file_name": "national-id.pdf",
      "file_url": "/storage/documents/EMP001/national_id_20240115_143100_def456.pdf",
      "file_size": 512345,
      "mime_type": "application/pdf",
      "description": "National ID card"
    }
  ]
}
```

---

## 3. Save Step with File Uploads API

### Endpoint
```
POST /api/v1/employees/onboarding/:employee_id/step/:step/upload
```

### Description
Saves step data with file uploads in a single request. Supports Step 1 (photo) and Step 6 (documents) with direct file uploads.

### Path Parameters
- `employee_id` (string, required): Employee ID (e.g., "EMP001")
- `step` (integer, required): Step number (1 for photo, 6 for documents)

### Request Format
`multipart/form-data`

### Form Fields

#### For Step 1 (Photo):
- `file` (file, optional): Photo file
- `data` (string, optional): JSON string of other step 1 data (excluding photo_url)

#### For Step 6 (Documents):
- `file` (file, optional): Single document file
- `files` (file[], optional): Multiple document files
- `document_type` (string, optional): Default document type for all files
- `files_document_type_0`, `files_document_type_1`, ... (string, optional): Document type for each file
- `description` (string, optional): Default description for all files
- `files_description_0`, `files_description_1`, ... (string, optional): Description for each file
- `data` (string, optional): JSON string of other step 6 data

### Success Response (200)
```json
{
  "success": true,
  "message": "Step saved successfully with file uploads",
  "data": {
    "draft_id": 1,
    "employee_id": "EMP001",
    "progress": 10.0,
    "completed_steps": [1],
    "finished_steps": [1],
    "unfinished_steps": [2, 3, 4, 5, 6, 7, 8, 9, 10],
    "is_completed": false,
    "steps": {
      "1": {
        "photo_url": "/storage/photos/EMP001/photo_20240115_143045_abc123.jpg",
        "first_name": "John",
        "last_name": "Doe",
        ...
      }
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T14:30:00Z"
  }
}
```

### Example Request - Step 1 with Photo (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/employees/onboarding/EMP001/step/1/upload" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/photo.jpg" \
  -F 'data={"first_name":"John","last_name":"Doe","gender":"male"}'
```

### Example Request - Step 6 with Documents (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/employees/onboarding/EMP001/step/6/upload" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "files=@/path/to/national-id.pdf" \
  -F "files=@/path/to/degree.pdf" \
  -F "files_document_type_0=identity" \
  -F "files_document_type_1=education" \
  -F "files_description_0=National ID card" \
  -F "files_description_1=University degree"
```

### Example Request - Step 6 with Documents (JavaScript)
```javascript
const formData = new FormData();
formData.append('files', file1);
formData.append('files', file2);
formData.append('files_document_type_0', 'identity');
formData.append('files_document_type_1', 'education');
formData.append('files_description_0', 'National ID card');
formData.append('files_description_1', 'University degree');

const response = await fetch('http://localhost:8080/api/v1/employees/onboarding/EMP001/step/6/upload', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
});

const result = await response.json();
console.log(result.data);
```

---

## 4. Delete File API

### Endpoint
```
POST /api/v1/employees/onboarding/delete-file
```

### Description
Deletes a file from storage using its file URL.

### Request Body
```json
{
  "file_url": "/storage/photos/EMP001/photo_20240115_143045_abc123.jpg"
}
```

### Required Fields
- `file_url` (string, required): Full file URL to delete

### Success Response (200)
```json
{
  "success": true,
  "message": "File deleted successfully"
}
```

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/employees/onboarding/delete-file" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "file_url": "/storage/photos/EMP001/photo_20240115_143045_abc123.jpg"
  }'
```

---

## Complete Workflow Examples

### Workflow 1: Upload Photo for Step 1

```javascript
// Step 1: Upload photo
const photoFormData = new FormData();
photoFormData.append('employee_id', 'EMP001');
photoFormData.append('file', photoFile);

const photoResponse = await fetch('/api/v1/employees/onboarding/upload-photo', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: photoFormData
});

const photoData = await photoResponse.json();
const photoURL = photoData.data.photo_url;

// Step 2: Save Step 1 with photo URL
const step1Data = {
  first_name: 'John',
  last_name: 'Doe',
  photo_url: photoURL,
  // ... other fields
};

await fetch('/api/v1/employees/onboarding/EMP001/step/1', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ step: 1, data: step1Data })
});
```

### Workflow 2: Upload Photo Directly in Step 1 (Single Request)

```javascript
// Upload photo and save step 1 in one request
const formData = new FormData();
formData.append('file', photoFile);
formData.append('data', JSON.stringify({
  first_name: 'John',
  last_name: 'Doe',
  gender: 'male',
  // ... other fields (photo_url will be added automatically)
}));

await fetch('/api/v1/employees/onboarding/EMP001/step/1/upload', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: formData
});
```

### Workflow 3: Upload Multiple Documents for Step 6

```javascript
// Upload multiple documents
const docFormData = new FormData();
docFormData.append('employee_id', 'EMP001');
docFormData.append('document_type', 'identity');
docFormData.append('file', nationalIdFile);

const doc1Response = await fetch('/api/v1/employees/onboarding/upload-document', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: docFormData
});

const doc1Data = await doc1Response.json();

// Upload second document
const doc2FormData = new FormData();
doc2FormData.append('employee_id', 'EMP001');
doc2FormData.append('document_type', 'education');
doc2FormData.append('file', degreeFile);

const doc2Response = await fetch('/api/v1/employees/onboarding/upload-document', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: doc2FormData
});

const doc2Data = await doc2Response.json();

// Save Step 6 with document URLs
const step6Data = {
  documents: [
    {
      document_type: 'identity',
      file_name: 'national-id.pdf',
      file_url: doc1Data.data.file_url,
      file_size: doc1Data.data.file_size,
      mime_type: doc1Data.data.mime_type,
      description: 'National ID card'
    },
    {
      document_type: 'education',
      file_name: 'degree.pdf',
      file_url: doc2Data.data.file_url,
      file_size: doc2Data.data.file_size,
      mime_type: doc2Data.data.mime_type,
      description: 'University degree'
    }
  ]
};

await fetch('/api/v1/employees/onboarding/EMP001/step/6', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ step: 6, data: step6Data })
});
```

### Workflow 4: Upload Multiple Documents Directly in Step 6 (Single Request)

```javascript
// Upload multiple documents and save step 6 in one request
const formData = new FormData();
formData.append('files', nationalIdFile);
formData.append('files', degreeFile);
formData.append('files_document_type_0', 'identity');
formData.append('files_document_type_1', 'education');
formData.append('files_description_0', 'National ID card');
formData.append('files_description_1', 'University degree');

await fetch('/api/v1/employees/onboarding/EMP001/step/6/upload', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: formData
});
```

---

## File Storage Structure

### Directory Organization
```
storage/
├── photos/
│   └── {employee_id}/
│       └── {filename}
└── documents/
    └── {employee_id}/
        └── {filename}
```

### Filename Format
Files are automatically renamed to prevent conflicts:
```
{original_name}_{timestamp}_{unique_id}.{extension}
```

Example: `national_id_20240115_143100_def456.pdf`

---

## File Access

### Static File Serving
Files are served statically at:
```
http://localhost:8080/storage/{category}/{employee_id}/{filename}
```

### Example URLs
- Photo: `http://localhost:8080/storage/photos/EMP001/photo_20240115_143045_abc123.jpg`
- Document: `http://localhost:8080/storage/documents/EMP001/national_id_20240115_143100_def456.pdf`

### CORS Configuration
Make sure your CORS configuration allows access to the `/storage` path if accessing from a different origin.

---

## Error Responses

### File Too Large
```json
{
  "success": false,
  "message": "file size exceeds maximum allowed size of 10485760 bytes"
}
```

### Invalid File Type
```json
{
  "success": false,
  "message": "file type application/zip is not allowed"
}
```

### Missing File
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": "file is required"
}
```

### File Not Found (Delete)
```json
{
  "success": false,
  "message": "failed to delete file: file does not exist"
}
```

---

## Supported File Types

### Images (for Photos)
- `image/jpeg` / `image/jpg`
- `image/png`
- `image/gif`
- `image/webp`

### Documents
- `application/pdf`
- `application/msword` (.doc)
- `application/vnd.openxmlformats-officedocument.wordprocessingml.document` (.docx)
- `application/vnd.ms-excel` (.xls)
- `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` (.xlsx)

---

## Best Practices

1. **Upload Files First**: Upload files and get URLs before saving step data
2. **Use Direct Upload Endpoint**: For Step 1 and Step 6, use the `/upload` endpoint to upload and save in one request
3. **Validate File Types**: Check file types on the frontend before uploading
4. **Handle Errors**: Always handle file upload errors gracefully
5. **Delete Unused Files**: Clean up files if step data is deleted or changed
6. **File Size Validation**: Validate file sizes on the frontend before uploading

---

## Quick Reference

| Operation | Endpoint | Method | Content-Type |
|-----------|----------|--------|--------------|
| Upload Photo | `/onboarding/upload-photo` | POST | multipart/form-data |
| Upload Document | `/onboarding/upload-document` | POST | multipart/form-data |
| Save Step with Files | `/onboarding/:employee_id/step/:step/upload` | POST | multipart/form-data |
| Delete File | `/onboarding/delete-file` | POST | application/json |

---

## Integration with Onboarding Steps

### Step 1: Personal Information
- **Photo Upload**: Use `/upload-photo` or `/step/1/upload` with file
- **Photo URL**: Automatically added to step data when using upload endpoint

### Step 6: Documents
- **Document Upload**: Use `/upload-document` or `/step/6/upload` with files
- **Document URLs**: Automatically added to step data when using upload endpoint

### Other Steps
- Steps 2-5, 7-10: Continue using JSON format (no file uploads needed)

---

## Security Considerations

1. **File Size Limits**: Maximum 10MB per file (configurable)
2. **File Type Validation**: Only allowed file types are accepted
3. **Authentication Required**: All endpoints require JWT authentication
4. **Tenant Isolation**: Files are organized by employee_id (tenant-scoped)
5. **Unique Filenames**: Prevents file overwrites and conflicts

---

## Configuration

### Environment Variables
Add to your `.env` file:
```env
STORAGE_PATH=./storage
STORAGE_URL=/storage
```

### Production Considerations
- Use a dedicated storage server or cloud storage (S3, Azure Blob, etc.)
- Implement file backup and recovery
- Set up CDN for file delivery
- Configure proper file permissions
- Implement file cleanup for deleted employees

---

**Last Updated:** 2024-01-15  
**Version:** 1.0.0
