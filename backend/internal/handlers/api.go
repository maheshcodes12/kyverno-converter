package handlers

import (
	"kyverno-converter-backend/internal/converter"
	"kyverno-converter-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// ConvertRequest defines the structure for incoming conversion requests.
type ConvertRequest struct {
	YAML string `json:"yaml" binding:"required"`
}

// ConvertResponse defines the structure for the conversion response.
type ConvertResponse struct {
	ConvertedYAML string `json:"convertedYaml"`
}

// ErrorResponse defines a generic error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ConvertPolicyHandler is the main HTTP handler for policy conversion.
func ConvertPolicyHandler(c *gin.Context) {
	var request ConvertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	// Unmarshal the input YAML into our legacy policy struct
	var legacyPolicy models.LegacyClusterPolicy
	if err := yaml.Unmarshal([]byte(request.YAML), &legacyPolicy); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Kyverno policy YAML: " + err.Error()})
		return
	}

	// Perform the conversion
	validatingPolicy, err := converter.ToValidatingPolicy(legacyPolicy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to convert policy: " + err.Error()})
		return
	}

	// Marshal the new policy struct back to YAML
	outputYAML, err := yaml.Marshal(validatingPolicy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate output YAML: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ConvertResponse{ConvertedYAML: string(outputYAML)})
}
