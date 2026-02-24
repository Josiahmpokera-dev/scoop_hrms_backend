# Recruitment Module - Frontend Integration Guide

This document is designed to assist frontend developers in integrating with the Recruitment Module API. It provides TypeScript interfaces, enum definitions, and usage flows.

## Base URLs
- **Development**: `http://localhost:8080/api/v1/recruitment`
- **Production**: `https://api.scoop-hrms.com/v1/recruitment`

## Authentication
Include the `Authorization` header in all requests (except `/public/apply`).
```javascript
headers: {
  "Authorization": "Bearer <your_jwt_token>"
}
```

---

## 1. Type Definitions (TypeScript)

### Enums & Constants

```typescript
export enum RequisitionStatus {
  Draft = "Draft",
  PendingApproval = "Pending Approval",
  Approved = "Approved",
  Rejected = "Rejected",
  Closed = "Closed"
}

export enum JobOpeningStatus {
  Draft = "Draft",
  Published = "Published",
  Closed = "Closed"
}

export enum ApplicationStage {
  Applied = "Applied",
  Screening = "Screening",
  Interview = "Interview",
  Offer = "Offer",
  Hired = "Hired",
  Rejected = "Rejected"
}

export enum InterviewType {
  PhoneScreen = "Phone Screen",
  Technical = "Technical",
  HR = "HR",
  Managerial = "Managerial"
}

export enum OfferStatus {
  Draft = "Draft",
  PendingApproval = "Pending Approval",
  Approved = "Approved",
  Sent = "Sent",
  Accepted = "Accepted",
  Rejected = "Rejected"
}
```

### Interfaces

```typescript
// --- Requisition ---
export interface JobRequisition {
  id: string;
  requisitionNo: string;
  jobTitle: string;
  department: string;
  location: string;
  grade: string;
  headcount: number;
  employmentType: string;
  hiringManager: string;
  requestedBy: string;
  targetStartDate: string; // ISO Date
  estimatedBudget: number;
  currency: string;
  justification: string;
  status: RequisitionStatus;
  approvalFlow: ApprovalStep[];
  createdAt: string;
}

export interface ApprovalStep {
  requisitionId: string;
  approver: string;
  role: string;
  status: string;
  comments: string;
  timestamp: string;
}

// --- Job Opening ---
export interface JobOpening {
  id: string;
  requisitionId: string;
  jobTitle: string;
  jobDescription: string;
  responsibilities: string[];
  mustHaveSkills: string[];
  niceToHaveSkills: string[];
  experienceMin: number;
  experienceMax: number;
  education: string[];
  salaryMin: number;
  salaryMax: number;
  salaryCurrency: string;
  expiryDate: string;
  status: JobOpeningStatus;
  platforms: string[];
}

// --- Candidate & Application ---
export interface Candidate {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  resumeUrl: string;
  linkedInUrl?: string;
  portfolioUrl?: string;
  coverLetter?: string;
}

export interface JobApplication {
  id: string;
  candidateId: string;
  jobOpeningId: string;
  stage: ApplicationStage;
  appliedDate: string;
  notes?: string;
  interview?: Interview[];
  offer?: Offer;
}

// --- Interview ---
export interface Interview {
  id: string;
  candidateId: string;
  jobOpeningId: string;
  applicationId: string;
  interviewType: InterviewType;
  round: number;
  scheduledDate: string;
  scheduledTime: string;
  durationMin: number;
  mode: string;
  interviewers: string[];
  status: string;
  feedback?: InterviewFeedback[];
}

export interface InterviewFeedback {
  interviewerEmail: string;
  rating: number; // 1-5
  strengths: string[];
  concerns: string[];
  recommendation: string;
  comments: string;
}

// --- Offer ---
export interface Offer {
  id: string;
  candidateId: string;
  jobOpeningId: string;
  applicationId: string;
  joiningDate: string;
  probationPeriod: number; // Months
  annualCTC: number;
  currency: string;
  components: SalaryComponent[];
  benefits: string[];
  expiryDate: string;
  status: OfferStatus;
}

export interface SalaryComponent {
  component: string;
  amount: number;
}
```

---

## 2. API Integration Flows

### Flow 1: Creating and Publishing a Job
1. **Create Requisition**: `POST /requisitions`
   - Initial status: `Pending Approval`.
2. **Approve Requisition**: `POST /requisitions/{id}/approval`
   - Required status for next step: `Approved`.
3. **Create Job Opening**: `POST /openings`
   - Links to `requisitionId`.
   - Initial status: `Draft`.
4. **Publish Job**: `POST /openings/{id}/publish`
   - Updates status to `Published`.
   - Job is now visible for applications.

### Flow 2: Candidate Application & Hiring
1. **Submit Application**: `POST /public/apply` (Public Endpoint)
   - Creates `Candidate` and `JobApplication`.
   - Initial stage: `Applied`.
2. **Review & Update Stage**: `PATCH /candidates/{id}/stage`
   - Move to `Screening` or `Interview`.
   - *Note*: Use `applicationId` or `candidateId` as needed by implementation.
3. **Schedule Interview**: `POST /interviews`
   - Requires valid `candidateId` and `jobId`.
4. **Submit Feedback**: `POST /interviews/{id}/feedback`
   - After interview is conducted.
5. **Create Offer**: `POST /offers`
   - When candidate passes all rounds.
6. **Approve & Send Offer**: `POST /offers/{id}/approval` -> `POST /offers/{id}/send`.

---

## 3. Endpoints Summary

| Feature | Method | Endpoint | Description |
|---|---|---|---|
| **Requisitions** | GET | `/requisitions` | List with filters (status, department) |
| | POST | `/requisitions` | Create new requisition |
| | GET | `/requisitions/:id` | Get details |
| | PUT | `/requisitions/:id` | Update details |
| | POST | `/requisitions/:id/approval` | Approve/Reject |
| **Job Openings** | GET | `/openings` | List openings |
| | POST | `/openings` | Create from approved requisition |
| | POST | `/openings/:id/publish` | Publish to platforms |
| **Candidates** | GET | `/candidates` | List candidates (filter by job/stage) |
| | POST | `/public/apply` | **Public**: Submit application |
| | PATCH | `/candidates/:id/stage` | Update application stage |
| **Interviews** | POST | `/interviews` | Schedule interview |
| | POST | `/interviews/:id/feedback` | Submit feedback |
| **Offers** | POST | `/offers` | Create offer |
| | POST | `/offers/:id/approval` | Approve offer |
| | POST | `/offers/:id/send` | Email offer to candidate |
| **Talent Pool** | POST | `/talent-pool` | Add candidate manually |
| | GET | `/talent-pool/search` | Search by skills/location |

---

## 4. Key Notes for Frontend
- **Dates**: All dates should be sent in ISO 8601 format (`YYYY-MM-DD` or `YYYY-MM-DDTHH:mm:ssZ`) unless specified otherwise.
- **Rich Text**: Fields like `jobDescription` and `responsibilities` may contain HTML or Markdown.
- **Resume Upload**: The `/public/apply` endpoint expects a `resume` string (Base64) or URL. If implementing file upload, ensure the file is uploaded to storage first (e.g., AWS S3) and the URL is sent to this API.
- **Error Handling**: Check `error.response.data.message` for user-friendly error messages.
