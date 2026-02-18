package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/models"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/modules/settings/repositories"
)

type MenuVisibilityService struct {
	repo *repositories.MenuVisibilityRepository
}

func NewMenuVisibilityService() *MenuVisibilityService {
	return &MenuVisibilityService{repo: repositories.NewMenuVisibilityRepository()}
}

// GetFullTree returns the complete menu tree with visibility applied.
func (s *MenuVisibilityService) GetFullTree() (map[string]interface{}, error) {
	hiddenSet, err := s.hiddenSet()
	if err != nil {
		return nil, err
	}

	tree := models.DefaultMenuTree()
	total, visible, hidden := applyVisibility(tree, hiddenSet)

	return map[string]interface{}{
		"items": tree,
		"summary": map[string]int{
			"total":   total,
			"visible": visible,
			"hidden":  hidden,
		},
	}, nil
}

// GetHiddenKeys returns the list of hidden menu keys (lightweight endpoint).
func (s *MenuVisibilityService) GetHiddenKeys() ([]string, error) {
	return s.repo.GetHiddenKeys()
}

// BulkUpdate applies a list of visibility changes with cascade and validation.
func (s *MenuVisibilityService) BulkUpdate(items []models.VisibilityItem, callerID uint) ([]models.VisibilityItem, error) {
	allKeys := models.AllMenuKeys()

	var invalidKeys []string
	var protectedKeys []string
	for _, item := range items {
		if !allKeys[item.MenuKey] {
			invalidKeys = append(invalidKeys, item.MenuKey)
		}
		if !item.Visible && models.ProtectedKeys[item.MenuKey] {
			protectedKeys = append(protectedKeys, item.MenuKey)
		}
	}
	if len(invalidKeys) > 0 {
		return nil, fmt.Errorf("invalid menu keys: %s", strings.Join(invalidKeys, ", "))
	}
	if len(protectedKeys) > 0 {
		return nil, fmt.Errorf("protected menu items cannot be hidden: %s", strings.Join(protectedKeys, ", "))
	}

	// Build a desired-state map from the request, then cascade.
	desired := make(map[string]bool)
	for _, item := range items {
		desired[item.MenuKey] = item.Visible
	}
	s.cascade(desired)

	// Convert to DB models.
	now := time.Now()
	var dbItems []models.MenuVisibilitySetting
	var result []models.VisibilityItem
	for key, vis := range desired {
		dbItems = append(dbItems, models.MenuVisibilitySetting{
			MenuKey:   key,
			Visible:   vis,
			UpdatedBy: &callerID,
			UpdatedAt: now,
		})
		result = append(result, models.VisibilityItem{MenuKey: key, Visible: vis})
	}

	if err := s.repo.UpsertMany(dbItems); err != nil {
		return nil, err
	}
	return result, nil
}

// ToggleSingle toggles a single menu item with cascade.
func (s *MenuVisibilityService) ToggleSingle(menuKey string, visible bool, callerID uint) ([]models.VisibilityItem, error) {
	allKeys := models.AllMenuKeys()
	if !allKeys[menuKey] {
		return nil, fmt.Errorf("unknown menu key: %s", menuKey)
	}
	if !visible && models.ProtectedKeys[menuKey] {
		return nil, fmt.Errorf("protected menu item cannot be hidden: %s", menuKey)
	}

	desired := map[string]bool{menuKey: visible}
	s.cascade(desired)

	now := time.Now()
	var dbItems []models.MenuVisibilitySetting
	var result []models.VisibilityItem
	for key, vis := range desired {
		dbItems = append(dbItems, models.MenuVisibilitySetting{
			MenuKey:   key,
			Visible:   vis,
			UpdatedBy: &callerID,
			UpdatedAt: now,
		})
		result = append(result, models.VisibilityItem{MenuKey: key, Visible: vis})
	}

	if err := s.repo.UpsertMany(dbItems); err != nil {
		return nil, err
	}
	return result, nil
}

// ResetAll resets every item to visible.
func (s *MenuVisibilityService) ResetAll() (int64, error) {
	allKeys := models.AllMenuKeys()
	resetCount := int64(len(allKeys))
	if _, err := s.repo.ResetAll(); err != nil {
		return 0, err
	}
	return resetCount, nil
}

// ---------- helpers ----------

func (s *MenuVisibilityService) hiddenSet() (map[string]bool, error) {
	keys, err := s.repo.GetHiddenKeys()
	if err != nil {
		return nil, err
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m, nil
}

// cascade enforces parent/child visibility rules.
func (s *MenuVisibilityService) cascade(desired map[string]bool) {
	childToParent, parentToChildren := models.ParentMap()

	// Rule 1 & 2: Hiding a parent hides all descendants.
	for key, vis := range desired {
		if !vis {
			s.hideDescendants(key, parentToChildren, desired)
		}
	}

	// Rule 3: Showing a child auto-shows ancestors.
	for key, vis := range desired {
		if vis {
			parent := childToParent[key]
			for parent != "" {
				desired[parent] = true
				parent = childToParent[parent]
			}
		}
	}

	// Never hide protected keys.
	for key := range models.ProtectedKeys {
		if v, ok := desired[key]; ok && !v {
			desired[key] = true
		}
	}
}

func (s *MenuVisibilityService) hideDescendants(key string, parentToChildren map[string][]string, desired map[string]bool) {
	for _, child := range parentToChildren[key] {
		if !models.ProtectedKeys[child] {
			desired[child] = false
		}
		s.hideDescendants(child, parentToChildren, desired)
	}
}

// applyVisibility walks the tree, sets Visible flags, and counts items.
func applyVisibility(nodes []*models.MenuNode, hiddenSet map[string]bool) (total, visible, hidden int) {
	for _, node := range nodes {
		total++
		if hiddenSet[node.Key] {
			node.Visible = false
			hidden++
		} else {
			node.Visible = true
			visible++
		}
		ct, cv, ch := applyVisibility(node.Children, hiddenSet)
		total += ct
		visible += cv
		hidden += ch
	}
	return
}
