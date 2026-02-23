# AWS S3 Storage & File Upload API

> All file uploads in the HRMS are now backed by AWS S3. The system stores the S3 URL in the database and returns it in API responses.

---

## Architecture

```
                  Frontend
                     │
                     ▼
        ┌──────────────────────┐
        │   Upload Handler     │  POST /uploads/file
        │   (or any handler)   │  POST /uploads/image
        └──────────┬───────────┘  POST /uploads/document
                   │
                   ▼
        ┌──────────────────────┐
        │   StorageService     │  internal/utils/storage/storage.go
        │                      │
        │  ┌─── S3 configured? │
        │  │   YES  │   NO     │
        │  ▼        ▼          │
        │ S3Client  Local Disk │
        └──────────────────────┘
                   │
          ┌────────┴────────┐
          ▼                 ▼
   ┌────────────┐   ┌────────────┐
   │  AWS S3    │   │  ./storage │
   │  Bucket    │   │  (local)   │
   └────────────┘   └────────────┘
```

When `AWS_S3_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_REGION` are all set, **every upload goes to S3**. Otherwise it falls back to local disk.

---

## Configuration

Add these to your `.env` file:

```env
# AWS S3 Configuration
AWS_REGION=eu-west-1
AWS_S3_BUCKET=scoop-hrms
AWS_S3_BASE_URL=https://scoop-hrms.s3.eu-west-1.amazonaws.com
AWS_S3_DOCS_PREFIX=gt/
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

| Variable               | Required | Description                                        |
|------------------------|----------|----------------------------------------------------|
| `AWS_REGION`           | Yes      | AWS region (e.g., `eu-west-1`)                     |
| `AWS_S3_BUCKET`        | Yes      | S3 bucket name                                     |
| `AWS_S3_BASE_URL`      | No       | Custom base URL (auto-generated if omitted)        |
| `AWS_S3_DOCS_PREFIX`   | No       | Path prefix inside bucket (e.g., `gt/`)            |
| `AWS_ACCESS_KEY_ID`    | Yes      | IAM access key                                     |
| `AWS_SECRET_ACCESS_KEY` | Yes     | IAM secret key                                     |

> **Security:** Never commit secrets to git. The `.env` file is already in `.gitignore`. Use `.env.example` as a template.

---

## S3 Key Structure

Files are stored under the prefix with a predictable path:

```
{AWS_S3_DOCS_PREFIX}{folder}/{identifier}/{timestamp}{extension}
```

Examples:
```
gt/documents/EMP001/1738856400000000000.pdf
gt/photos/EMP001/1738856400000000000.jpg
gt/helpdesk/42/1738856400000000000.png
gt/general/misc/1738856400000000000.docx
```

---

## API Endpoints — Quick Reference

| # | Method   | Endpoint                   | Description               | Access         |
|---|----------|----------------------------|---------------------------|----------------|
| 1 | `POST`   | `/api/v1/uploads/file`     | Upload any allowed file   | Authenticated  |
| 2 | `POST`   | `/api/v1/uploads/image`    | Upload image only         | Authenticated  |
| 3 | `POST`   | `/api/v1/uploads/document` | Upload document only      | Authenticated  |
| 4 | `DELETE` | `/api/v1/uploads`          | Delete file by URL        | Authenticated  |
| 5 | `GET`    | `/api/v1/uploads/info`     | Storage configuration     | Authenticated  |

All endpoints require `Authorization: Bearer <token>`.

---

## Detailed API Reference

---

### 1. Upload File (Any Type)

Upload any allowed file type (images + documents).

```
POST /api/v1/uploads/file
Content-Type: multipart/form-data
```

**Form Fields:**

| Field      | Type   | Required | Default   | Description                     |
|------------|--------|----------|-----------|---------------------------------|
| file       | file   | Yes      | —         | The file to upload              |
| folder     | string | No       | `general` | S3 folder/category              |
| identifier | string | No       | `misc`    | Sub-folder (e.g., employee ID)  |

**Response (201 Created):**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "url": "https://scoop-hrms.s3.eu-west-1.amazonaws.com/gt/documents/EMP001/1738856400000000000.pdf",
    "file_name": "contract.pdf",
    "file_size": 245760,
    "mime_type": "application/pdf",
    "folder": "documents",
    "storage": "s3"
  }
}
```

**Limits:** Max 25 MB.

**Allowed Types:**
- Images: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`, `.svg`
- Documents: `.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx`, `.txt`, `.csv`, `.zip`, `.rar`

---

### 2. Upload Image

Image-only upload with stricter validation.

```
POST /api/v1/uploads/image
Content-Type: multipart/form-data
```

**Form Fields:** Same as Upload File.

**Limits:** Max 10 MB. Only `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`, `.svg`.

**Response:** Same structure as Upload File.

---

### 3. Upload Document

Document-only upload.

```
POST /api/v1/uploads/document
Content-Type: multipart/form-data
```

**Form Fields:** Same as Upload File.

**Limits:** Max 25 MB. Only `.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx`, `.txt`, `.csv`, `.zip`, `.rar`.

**Response:** Same structure as Upload File.

---

### 4. Delete File

Delete a file from storage (S3 or local) by its URL.

```
DELETE /api/v1/uploads
Content-Type: application/json
```

**Request Body:**
```json
{
  "url": "https://scoop-hrms.s3.eu-west-1.amazonaws.com/gt/documents/EMP001/1738856400000000000.pdf"
}
```

**Response:**
```json
{
  "success": true,
  "message": "File deleted successfully",
  "data": null
}
```

---

### 5. Storage Info

Get current storage configuration.

```
GET /api/v1/uploads/info
```

**Response:**
```json
{
  "success": true,
  "message": "Storage configuration",
  "data": {
    "backend": "s3",
    "max_file_size_mb": 25,
    "allowed_images": ["jpg", "jpeg", "png", "gif", "webp", "svg"],
    "allowed_documents": ["pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "csv", "zip", "rar"]
  }
}
```

---

## Integration with Existing Features

All existing upload functionality now automatically uses S3:

| Feature                 | Endpoint                                     | What Changed                              |
|-------------------------|----------------------------------------------|-------------------------------------------|
| Employee Documents      | `POST /employees/onboarding/upload-document`  | URLs are now S3 URLs stored in DB         |
| Employee Photos         | `POST /employees/onboarding/upload-photo`     | URLs are now S3 URLs stored in DB         |
| Helpdesk Attachments    | `POST /helpdesk/tickets/:id/attachments`      | URLs are now S3 URLs stored in DB         |
| File Deletion           | `POST /employees/onboarding/delete-file`      | Deletes from S3 instead of local disk     |

**No API contract changes.** The only visible difference is that `file_url` / `url` fields now contain S3 URLs instead of local paths.

---

## Frontend Integration

### Before (Local Storage)
```
"file_url": "http://localhost:8080/storage/documents/EMP001/1738856400000000000.pdf"
```

### After (S3)
```
"file_url": "https://scoop-hrms.s3.eu-west-1.amazonaws.com/gt/documents/EMP001/1738856400000000000.pdf"
```

The frontend should:
1. Use the `url` / `file_url` from API responses directly — it's always a full URL
2. For images, use the URL in `<img src="...">` tags
3. For documents, use the URL for download links

### Uploading a File (Frontend Example)

```javascript
const uploadFile = async (file, folder, identifier) => {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('folder', folder);      // e.g., "documents", "photos"
  formData.append('identifier', identifier); // e.g., employee ID

  const res = await fetch('/api/v1/uploads/file', {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` },
    body: formData,
  });

  const data = await res.json();
  // data.data.url → the S3 URL to store/display
  return data.data.url;
};
```

---

## File Structure

```
internal/
├── config/
│   └── config.go              ← Added AWSConfig struct
├── utils/
│   └── storage/
│       ├── storage.go         ← Updated: S3-aware StorageService + ResolveURL
│       └── s3.go              ← NEW: S3 upload/delete client
└── modules/
    └── uploads/
        └── handlers/
            └── upload_handler.go  ← NEW: Generic upload endpoints
```

---

## Error Responses

| Code | Scenario                                     |
|------|----------------------------------------------|
| 400  | No file provided                             |
| 400  | File too large (> 25 MB for docs, > 10 MB for images) |
| 400  | File type not allowed                        |
| 400  | Missing URL in delete request                |
| 500  | S3 upload failed (network, permissions, etc.) |
| 500  | S3 delete failed                             |

---

## Fallback Behavior

If the AWS environment variables are not set, the system **automatically falls back to local disk storage** (`./storage` directory). This is useful for development without AWS credentials. The `GET /uploads/info` endpoint reports the active backend (`"s3"` or `"local"`).
