# Payroll Module

Comprehensive payroll management system with built-in data encryption for sensitive employee information.

## Features

- **Salary Structures**: Manage employee salary components and structures
- **Payroll Processing**: Automated payroll calculation and processing
- **Payslip Generation**: Generate and distribute digital payslips
- **Tax Calculations**: Tanzania PAYE tax calculations
- **Statutory Deductions**: NSSF, NHIF, and other statutory compliance
- **🔒 Data Encryption**: Automatic encryption of sensitive payroll data

## Data Encryption Security

### Overview
The payroll module includes built-in encryption for sensitive employee data to ensure compliance with data protection regulations and protect employee privacy.

### Encrypted Fields
The following sensitive fields are automatically encrypted:
- **Bank Account Information**: Employee bank account numbers
- **Tax Identification Numbers**: TIN numbers
- **Social Security Numbers**: NSSF membership numbers
- **Health Insurance Numbers**: NHIF membership numbers

### Encryption Implementation

#### Key Features
- **AES-256 Encryption**: Military-grade encryption algorithm
- **Employee-Specific Keys**: Each employee's data encrypted with unique key derived from `employee_id`
- **Automatic Encryption/Decryption**: No manual intervention required
- **Zero Frontend Impact**: Encryption happens transparently in backend

#### Technical Details
- **Encryption Service**: [encryption_service.go](file:///Users/josiampokera/GT%20Development/node/hrms-backend/internal/utils/encryption/encryption_service.go)
- **Payroll Encryption**: [payroll_encryption_service.go](file:///Users/josiampokera/GT%20Development/node/hrms-backend/internal/modules/payroll/services/payroll_encryption_service.go)
- **Repository Integration**: [payslip_repository.go](file:///Users/josiampokera/GT%20Development/node/hrms-backend/internal/modules/payroll/repositories/payslip_repository.go)

### Configuration

#### Environment Variables
```bash
# Encryption Key (32 characters recommended for AES-256)
ENCRYPTION_KEY=your-secure-encryption-key-32-chars
```

#### Fallback Behavior
- Development: Uses default key if `ENCRYPTION_KEY` not set
- Production: **REQUIRES** `ENCRYPTION_KEY` environment variable
- Graceful degradation: Operations continue if encryption fails (with warnings)

### Usage

#### Automatic Encryption
All sensitive data is automatically encrypted:
1. **Before Save**: Data encrypted before database persistence
2. **After Retrieval**: Data decrypted after database loading
3. **Transparent**: No code changes required in business logic

#### Manual Encryption (Advanced)
```go
// Get encryption service
encryptionService, err := services.NewPayrollEncryptionService()
if err != nil {
    // Handle error
}

// Encrypt specific data
err = encryptionService.EncryptPayslip(payslip)
if err != nil {
    // Handle error
}

// Decrypt specific data  
err = encryptionService.DecryptPayslip(payslip)
if err != nil {
    // Handle error
}
```

### Security Considerations

1. **Key Management**:
   - Store `ENCRYPTION_KEY` securely in production
   - Rotate keys periodically for enhanced security
   - Use different keys for different environments

2. **Data Recovery**:
   - Backup encryption keys securely
   - Without proper keys, encrypted data cannot be recovered

3. **Compliance**:
   - Meets data protection requirements for sensitive personal information
   - Encrypts financial and identification data at rest

### API Endpoints

All payroll API endpoints automatically handle encryption:
- `GET /api/payroll/payslips` - Returns decrypted data
- `POST /api/payroll/payslips` - Encrypts data before storage
- `PUT /api/payroll/payslips/{id}` - Handles encryption on update
- `GET /api/payroll/payslips/{id}` - Returns decrypted data

### Testing

Run encryption tests to verify functionality:
```bash
# Test encryption service
go test ./internal/utils/encryption/...

# Test payroll encryption  
go test ./internal/modules/payroll/services/...

# Test repository integration
go test ./internal/modules/payroll/repositories/...
```

## Status

✅ **Fully Implemented** - All features including data encryption are complete and operational.

---

**Security Note**: This implementation ensures that sensitive employee payroll data is protected with industry-standard encryption, meeting compliance requirements for data protection regulations.
