package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/service/generator"
)

type GeneratorHandler struct {
	gen    generator.GeneratorService
	logger lib.Logger
}

func NewGeneratorHandler(gen generator.GeneratorService, logger lib.Logger) GeneratorHandler {
	return GeneratorHandler{gen: gen, logger: logger}
}

func (h GeneratorHandler) Schema(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	schema, err := h.gen.GetSchema(r.Context(), slug)
	if err != nil {
		respondGeneratorError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, schema)
}

func (h GeneratorHandler) Generate(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var body struct {
		Values map[string]any `json:"values"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	accountID, _ := middleware.GetAccountID(r.Context())
	var aid *uuid.UUID
	if accountID != uuid.Nil {
		aid = &accountID
	}

	result, err := h.gen.Generate(r.Context(), aid, generator.GenerateInput{
		PlateSlug: slug,
		Values:    body.Values,
	})
	if err != nil {
		h.logger.Debugf("generate error for %s: %v", slug, err)
		respondGeneratorError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+slug+".zip\"")
	w.Header().Set("X-Generation-ID", result.GenerationID.String())
	w.WriteHeader(http.StatusOK)
	w.Write(result.ZipBytes) //nolint:errcheck
}

func respondGeneratorError(w http.ResponseWriter, err error) {
	switch {
	case err == generator.ErrPlateNotFound:
		respondError(w, http.StatusNotFound, err.Error())
	case err == generator.ErrPlateNotApproved:
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case err == generator.ErrNoRepoURL:
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case err == generator.ErrMissingYAML:
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case err == generator.ErrFetchFailed:
		respondError(w, http.StatusBadGateway, err.Error())
	case err == generator.ErrInvalidInput:
		respondError(w, http.StatusBadRequest, err.Error())
	case err == generator.ErrTemplateFailed:
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "generation failed")
	}
}
