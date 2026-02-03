package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
)

// RequesterInfo holds requester details for response (camelCase)
type RequesterInfo struct {
	ID         string
	Name       string
	Email      string
	Phone      string
	Department string
}

// AssignedToInfo holds assignee details for response (camelCase)
type AssignedToInfo struct {
	ID    interface{}
	Name  string
	Email string
	Team  string
}

// formatTicketResponseHelper formats a ticket for response using camelCase as per API doc
func formatTicketResponseHelper(ticket *models.Ticket, includeDetails bool, requester *RequesterInfo, assignedTo *AssignedToInfo) map[string]interface{} {
	ticketMap := map[string]interface{}{
		"id":           ticket.ID,
		"ticketNumber": ticket.TicketNumber,
		"title":        ticket.Title,
		"description":  ticket.Description,
		"category":     ticket.Category,
		"priority":     string(ticket.Priority),
		"status":       string(ticket.Status),
		"channel":      string(ticket.Channel),
		"createdAt":    ticket.CreatedAt,
		"updatedAt":    ticket.UpdatedAt,
	}

	if ticket.SubCategory != nil {
		ticketMap["subCategory"] = *ticket.SubCategory
	}

	if requester != nil {
		ticketMap["requester"] = map[string]interface{}{
			"id":         requester.ID,
			"name":       requester.Name,
			"email":      requester.Email,
			"phone":      requester.Phone,
			"department": requester.Department,
		}
	} else if ticket.RequesterUserID != nil {
		req := map[string]interface{}{"id": fmt.Sprintf("USER-%d", *ticket.RequesterUserID)}
		if ticket.RequesterName != nil {
			req["name"] = *ticket.RequesterName
		}
		if ticket.RequesterEmail != nil {
			req["email"] = *ticket.RequesterEmail
		}
		req["phone"], req["department"] = "", ""
		ticketMap["requester"] = req
	} else if ticket.RequesterEmployeeID != nil {
		ticketMap["requester"] = map[string]interface{}{
			"id": *ticket.RequesterEmployeeID,
		}
	}

	if assignedTo != nil {
		ticketMap["assignedTo"] = map[string]interface{}{
			"id":    assignedTo.ID,
			"name":  assignedTo.Name,
			"email": assignedTo.Email,
			"team":  assignedTo.Team,
		}
	} else if ticket.AssignedToName != nil {
		at := map[string]interface{}{"name": *ticket.AssignedToName}
		if ticket.AssignedToEmail != nil {
			at["email"] = *ticket.AssignedToEmail
		}
		if ticket.AssignedToTeam != nil {
			at["team"] = *ticket.AssignedToTeam
		}
		if ticket.AssignedToID != nil {
			at["id"] = *ticket.AssignedToID
		}
		ticketMap["assignedTo"] = at
	}

	if ticket.DueDate != nil {
		ticketMap["dueDate"] = *ticket.DueDate
	}
	if ticket.SLAStatus != nil {
		ticketMap["slaStatus"] = *ticket.SLAStatus
	}
	if ticket.FirstResponseTime != nil {
		ticketMap["firstResponseTime"] = formatDurationMinutes(ticket.CreatedAt, ticket.FirstResponseTime)
	}
	if ticket.ResolutionTime != nil {
		ticketMap["resolutionTime"] = formatDurationMinutes(ticket.CreatedAt, ticket.ResolutionTime)
	}

	if ticket.TagsJSON != nil && *ticket.TagsJSON != "" {
		var tags []string
		_ = json.Unmarshal([]byte(*ticket.TagsJSON), &tags)
		ticketMap["tags"] = tags
	}

	if includeDetails {
		if ticket.ResolvedAt != nil {
			ticketMap["resolvedAt"] = *ticket.ResolvedAt
		}
		if ticket.ClosedAt != nil {
			ticketMap["closedAt"] = *ticket.ClosedAt
		}
		if ticket.CSATRating != nil {
			ticketMap["csatRating"] = *ticket.CSATRating
		}
		if ticket.CSATComment != nil {
			ticketMap["csatComment"] = *ticket.CSATComment
		}
	}

	return ticketMap
}

// formatFirstResponseDuration returns a string like "15 minutes" from ticket createdAt to firstResponseTime
func formatDurationMinutes(createdAt time.Time, firstResponseTime *time.Time) string {
	if firstResponseTime == nil {
		return ""
	}
	d := firstResponseTime.Sub(createdAt)
	m := int(d.Minutes())
	if m < 1 {
		m = 1
	}
	if m < 60 {
		if m == 1 {
			return "1 minute"
		}
		return formatInt(m) + " minutes"
	}
	h := m / 60
	m = m % 60
	if m == 0 {
		if h == 1 {
			return "1 hour"
		}
		return formatInt(h) + " hours"
	}
	if h == 1 {
		return "1 hour " + formatInt(m) + " minutes"
	}
	return formatInt(h) + " hours " + formatInt(m) + " minutes"
}

func formatInt(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

// FormatTicketResponse formats a ticket for response (exported for use in other handlers).
// requester and assignedTo may be nil; then minimal info is built from ticket.
func FormatTicketResponse(ticket *models.Ticket, includeDetails bool, requester *RequesterInfo, assignedTo *AssignedToInfo) map[string]interface{} {
	return formatTicketResponseHelper(ticket, includeDetails, requester, assignedTo)
}

// formatAttachmentsForList returns attachment list with camelCase keys and human-readable size
func formatAttachmentsForList(attachments []models.Attachment) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(attachments))
	for i := range attachments {
		a := &attachments[i]
		out = append(out, map[string]interface{}{
			"id":         a.ID,
			"name":       a.Name,
			"url":        a.URL,
			"size":       formatAttachmentSize(a.Size),
			"type":       a.Type,
			"uploadedAt": a.UploadedAt,
		})
	}
	return out
}

func formatAttachmentSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	kb := bytes / 1024
	if kb < 1024 {
		return fmt.Sprintf("%d KB", kb)
	}
	mb := kb / 1024
	return fmt.Sprintf("%d MB", mb)
}

// BuildAssignedToFromTicket builds AssignedToInfo from ticket fields (for when no extra lookup is done)
func BuildAssignedToFromTicket(ticket *models.Ticket) *AssignedToInfo {
	if ticket.AssignedToName == nil {
		return nil
	}
	at := &AssignedToInfo{Name: *ticket.AssignedToName}
	if ticket.AssignedToID != nil {
		at.ID = *ticket.AssignedToID
	}
	if ticket.AssignedToEmail != nil {
		at.Email = *ticket.AssignedToEmail
	}
	if ticket.AssignedToTeam != nil {
		at.Team = *ticket.AssignedToTeam
	}
	return at
}
