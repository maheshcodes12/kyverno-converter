package converter

import (
	"fmt"
	"kyverno-converter-backend/internal/models"
	"strings"
)

// GenerateCEL is the dispatcher that selects the correct generation logic
// based on the contents of the `validate` block.
func GenerateCEL(validate *models.ValidateBlock) (string, error) {
	if validate.Pattern != nil {
		// The rule uses a direct pattern.
		return generateFromPattern(validate.Pattern, "object")
	}
	if len(validate.ForEach) > 0 {
		// The rule uses a foreach loop.
		return generateFromForEach(validate.ForEach)
	}
	return "", fmt.Errorf("unsupported validation type: only 'pattern' and 'foreach' are implemented")
}

// generateFromForEach handles `foreach` validation rules.
func generateFromForEach(forEachs []models.ForEach) (string, error) {
	var allConditions []string
	for _, fe := range forEachs {
		// Convert JMESPath list to a CEL path (e.g., `request.object.spec.containers` -> `object.spec.containers`)
		listPath := strings.Replace(fe.List, "request.object", "object", 1)

		// The variable `element` will be used in the CEL `all()` function to refer to each item in the list.
		elementVar := "element"

		// Recursively generate the CEL expression for the pattern, applied to each `element`.
		subExpression, err := generateFromPattern(fe.Pattern, elementVar)
		if err != nil {
			return "", err
		}

		// Wrap the sub-expression in a `has()` check for safety, then in an `all()` quantifier.
		// This means: "for all elements in the list, the condition must be true".
		// The `!has(...)` part makes the rule pass if the list itself doesn't exist.
		condition := fmt.Sprintf("!has(%s) || %s.all(%s, %s)", listPath, listPath, elementVar, subExpression)
		allConditions = append(allConditions, condition)
	}

	return strings.Join(allConditions, " && "), nil
}

// generateFromPattern is the core recursive function that walks the pattern YAML
// and builds a CEL expression.
// `patternData` is the current piece of the pattern (e.g., a map, a string).
// `path` is the CEL path to the current data (e.g., `object.spec.metadata`).
func generateFromPattern(patternData interface{}, path string) (string, error) {
	switch p := patternData.(type) {
	case map[string]interface{}:
		return generateFromMap(p, path)
	case []interface{}:
		// This is a rare case in Kyverno patterns and is complex to generalize.
		// A proper implementation would need to know if it's an AND or OR condition.
		// For now, we'll return an error indicating it's not supported.
		return "", fmt.Errorf("array patterns are not supported; please use a `foreach` rule instead")
	case string:
		return generateValueCheck(p, path)
	case bool, int, int64, float64:
		return fmt.Sprintf("%s == %v", path, p), nil
	case nil:
		return fmt.Sprintf("!has(%s)", path), nil
	default:
		return "", fmt.Errorf("unsupported type in pattern: %T", p)
	}
}

// generateFromMap handles object patterns.
func generateFromMap(patternMap map[string]interface{}, path string) (string, error) {
	var conditions []string
	for key, value := range patternMap {
		var currentPath string
		var condition string
		var err error

		// Handle Kyverno's special anchors and operators in keys
		if strings.HasPrefix(key, "(") && strings.HasSuffix(key, ")") {
			// This is a conditional anchor: `(key): value`.
			// The logic is: IF the field `key` exists, THEN its value must match `value`.
			// This translates to CEL as: `!has(path.to.key) || (path.to.key matches value)`
			actualKey := strings.TrimSuffix(strings.TrimPrefix(key, "("), ")")
			currentPath = fmt.Sprintf("%s.%s", path, actualKey)
			subExpr, err := generateFromPattern(value, currentPath)
			if err != nil {
				return "", err
			}
			condition = fmt.Sprintf("!has(%s) || (%s)", currentPath, subExpr)
		} else if strings.HasPrefix(key, "+(") && strings.HasSuffix(key, ")") {
			// This is an "add if not present" anchor, used in mutation, not validation.
			// We'll treat it like a regular key for validation purposes.
			actualKey := strings.TrimSuffix(strings.TrimPrefix(key, "+("), ")")
			currentPath = fmt.Sprintf("%s.%s", path, actualKey)
			condition, err = generateFromPattern(value, currentPath)
		} else {
			// This is a regular key.
			currentPath = fmt.Sprintf("%s.%s", path, key)
			// The sub-pattern must be true, and the field itself must exist.
			subExpr, err := generateFromPattern(value, currentPath)
			if err != nil {
				return "", err
			}
			condition = fmt.Sprintf("has(%s) && %s", path, subExpr)
		}

		if err != nil {
			return "", err
		}
		conditions = append(conditions, condition)
	}

	if len(conditions) == 0 {
		return "true", nil
	}
	return fmt.Sprintf("(%s)", strings.Join(conditions, " && ")), nil
}

// generateValueCheck handles string values, which may contain wildcards.
func generateValueCheck(value string, path string) (string, error) {
	// Handle special image parsing
	if strings.Contains(path, "image") {
		if strings.HasPrefix(value, "!*:") {
			// Disallow a specific tag, e.g., `!*:latest`
			tag := strings.TrimPrefix(value, "!*:")
			return fmt.Sprintf("image(%s).tag != '%s'", path, tag), nil
		}
		if strings.HasSuffix(value, "*") && !strings.Contains(value, "|") {
			// Restrict to a registry, e.g., `my-registry.io/*`
			registry := strings.TrimSuffix(value, "/*")
			return fmt.Sprintf("image(%s).registry == '%s'", path, registry), nil
		}
	}

	// Handle standard wildcards
	switch value {
	case "?*": // Field must exist and be non-empty.
		return fmt.Sprintf("has(%s) && %s != ''", path, path), nil
	case "*": // Field must exist.
		return fmt.Sprintf("has(%s)", path), nil
	case "!*": // Field must not exist.
		return fmt.Sprintf("!has(%s)", path), nil
	default: // Strict equality check.
		return fmt.Sprintf("%s == '%s'", path, value), nil
	}
}
