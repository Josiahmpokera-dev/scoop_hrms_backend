package handlers

import (
	"encoding/json"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/helpdesk/models"
)

// formatTicketResponseHelper formats a ticket for response
func formatTicketResponseHelper(ticket *models.Ticket, includeDetails bool) map[string]interface{} {
	ticketMap := map[string]interface{}{
		"id":            ticket.ID,
		"ticket_number": ticket.TicketNumber,
		"title":         ticket.Title,
		"description":   ticket.Description,
		"category":      ticket.Category,
		"priority":      string(ticket.Priority),
		"status":        string(ticket.Status),
		"channel":       string(ticket.Channel),
		"created_at":    ticket.CreatedAt,
		"updated_at":    ticket.UpdatedAt,
	}

	if ticket.SubCategory != nil {
		ticketMap["sub_category"] = *ticket.SubCategory
	}
	if ticket.RequesterEmployeeID != nil {
		ticketMap["requester"] = map[string]interface{}{
			"id": *ticket.RequesterEmployeeID,
		}
	}
	if ticket.AssignedToName != nil {
		assignedToMap := map[string]interface{}{
			"name": *ticket.AssignedToName,
		}
		if ticket.AssignedToEmail != nil {
			assignedToMap["email"] = *ticket.AssignedToEmail
		}
		if ticket.AssignedToTeam != nil {
			assignedToMap["team"] = *ticket.AssignedToTeam
		}
		ticketMap["assigned_to"] = assignedToMap
	}
	if ticket.DueDate != nil {
		ticketMap["due_date"] = *ticket.DueDate
	}
	if ticket.SLAStatus != nil {
		ticketMap["sla_status"] = *ticket.SLAStatus
	}
	if ticket.FirstResponseTime != nil {
		ticketMap["first_response_time"] = ticket.FirstResponseTime.Format("15:04")
	}
	if ticket.ResolutionTime != nil {
		ticketMap["resolution_time"] = ticket.ResolutionTime.Format("15:04")
	}

	// Parse tags
	if ticket.TagsJSON != nil && *ticket.TagsJSON != "" {
		var tags []string
		_ = json.Unmarshal([]byte(*ticket.TagsJSON), &tags)
		ticketMap["tags"] = tags
	}

	if includeDetails {
		if ticket.ResolvedAt != nil {
			ticketMap["resolved_at"] = *ticket.ResolvedAt
		}
		if ticket.ClosedAt != nil {
			ticketMap["closed_at"] = *ticket.ClosedAt
		}
		if ticket.CSATRating != nil {
			ticketMap["csat_rating"] = *ticket.CSATRating
		}
		if ticket.CSATComment != nil {
			ticketMap["csat_comment"] = *ticket.CSATComment
		}
	}

	return ticketMap
}
