package converter

import (
	"errors"
	"fmt"
	"kyverno-converter-backend/internal/models"
	"strings"
)

// ToValidatingPolicy is the main entry point for the conversion.
// It orchestrates the transformation from a legacy policy to a modern ValidatingPolicy.
func ToValidatingPolicy(legacy models.LegacyClusterPolicy) (*models.ValidatingPolicy, error) {
	if len(legacy.Spec.Rules) == 0 {
		return nil, errors.New("no rules found in the policy")
	}

	// Create the base structure for the new policy
	newPolicy := &models.ValidatingPolicy{
		APIVersion: "policies.kyverno.io/v1alpha1",
		Kind:       "ValidatingPolicy",
		Metadata:   legacy.Metadata,
		Spec: models.ValidatingSpec{
			ValidationActions:  []string{mapFailureAction(legacy.Spec.ValidationFailureAction)},
			Background:         legacy.Spec.Background,
			MatchConstraints:   convertMatchBlock(legacy.Spec.Rules[0].Match),
			ExcludeConstraints: convertMatchBlock(legacy.Spec.Rules[0].Exclude),
			Validations:        []models.Validation{},
		},
	}

	// Each rule in the legacy policy becomes a `validation` block in the new policy.
	for _, rule := range legacy.Spec.Rules {
		// Generate the core CEL expression from the `validate` block.
		celExpression, err := GenerateCEL(&rule.Validate)
		if err != nil {
			return nil, fmt.Errorf("error converting rule '%s': %w", rule.Name, err)
		}

		validation := models.Validation{
			Message:    rule.Validate.Message,
			Expression: celExpression,
		}
		newPolicy.Spec.Validations = append(newPolicy.Spec.Validations, validation)
	}

	return newPolicy, nil
}

// mapFailureAction converts the legacy action name to the new one.
func mapFailureAction(action string) string {
	if strings.ToLower(action) == "enforce" {
		return "Deny"
	}
	return "Audit" // Default to Audit
}

// convertMatchBlock transforms the legacy match/exclude blocks to the new format.
func convertMatchBlock(block models.MatchBlock) models.MatchConstraints {
	// This is a simplified conversion. A full implementation would handle `any` vs `all` logic.
	// For this tool, we'll assume a single `any` block is the common case.
	constraints := models.MatchConstraints{}
	if len(block.Any) > 0 {
		filter := block.Any[0]
		resourceRule := models.ResourceRule{
			APIGroups:   []string{""},   // Default, should be inferred if possible
			APIVersions: []string{"v1"}, // Default, should be inferred if possible
			Operations:  []string{"CREATE", "UPDATE"},
			Resources:   filter.Resources.Kinds,
		}
		constraints.ResourceRules = append(constraints.ResourceRules, resourceRule)
	}
	return constraints
}
