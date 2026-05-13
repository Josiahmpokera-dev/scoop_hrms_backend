package models

// Recruitment pipeline stages for JobApplication.Stage (string).
// UI should map tabs/filters to these values.
const (
	StageApplied            = "Applied"
	StageScreening          = "Screening"
	StageInterview          = "Interview"          // shortlisted / interview round
	StageInterviewCompleted = "InterviewCompleted" // all rounds done; awaiting decision
	StageOffer              = "Offer"
	StageHired              = "Hired"
	StageRejected           = "Rejected"
	StageTalentPool         = "TalentPool" // not advancing; retained for future roles
)
