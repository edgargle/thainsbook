package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"thainsbook/internal/models"
	"thainsbook/internal/utils"

	"github.com/google/uuid"
)

func (a *Application) HandleGetUserEntries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(UserIdKey).(string)

	entries, err := a.Entries.GetEntriesByUser(userId)
	if err != nil {
		log.Printf("Error getting entries by user: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Error in retrieving entries.")
		return
	}

	if entries == nil {
		entries = []models.EntryResponse{}
	}

	response := map[string]interface{}{
		"entries": entries,
		"count":   len(entries),
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (a *Application) HandleGetUserEntry(w http.ResponseWriter, r *http.Request) {
	entryId := r.PathValue("id")
	ctx := r.Context()
	userId := ctx.Value(UserIdKey).(string)

	entry, err := a.Entries.GetEntry(entryId, userId)
	if err != nil {
		log.Println("Error Retrieving Entry:", err)
		utils.WriteError(w, http.StatusNotFound, "Entry not found.")
		return
	}

	response := map[string]interface{}{
		"entry": entry,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (a *Application) HandleCreateEntry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(UserIdKey).(string)
	decoder := json.NewDecoder(r.Body)
	var e models.EntryRequest
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)
	if err != nil {
		log.Println("JSON Decode Error:", err)
		utils.WriteError(w, http.StatusBadRequest, "Unable to process request.")
		return
	}

	if e.Content == "" {
		utils.WriteError(w, http.StatusBadRequest, "Cannot add empty content.")
		return
	}

	time, err := utils.ParseEntryDate(e.EntryDate)
	if err != nil {
		log.Println("Error Parsing EntryDate:", err)
		utils.WriteError(w, http.StatusInternalServerError, "Unable to process request.")
	}

	newEntry := models.EntryDto{
		Id:        uuid.NewString(),
		Title:     e.Title,
		Content:   e.Content,
		EntryDate: time,
		UserId:    userId,
	}

	err = a.Entries.AddEntry(&newEntry)
	if err != nil {
		log.Println("Error Adding EntryDto:", err)
		utils.WriteError(w, http.StatusInternalServerError, "Unable to process request.")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "New entry created successfully."})
}

func (a *Application) HandleUpdateEntry(w http.ResponseWriter, r *http.Request) {
	entryId := r.PathValue("id")

	if entryId == "" {
		utils.WriteError(w, http.StatusBadRequest, "Please provide a valid entry id.")
		return
	}

	ctx := r.Context()
	userId := ctx.Value(UserIdKey).(string)
	decoder := json.NewDecoder(r.Body)
	var e models.EntryRequest
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)
	if err != nil {
		log.Println("JSON Decode Error:", err)
		utils.WriteError(w, http.StatusBadRequest, "Unable to process request.")
		return
	}

	if e.Content == "" && e.EntryDate == "" && e.Title == "" {
		utils.WriteError(w, http.StatusBadRequest, "Please provide at least one field to update: Title, Content or Entry Date.")
		return
	}

	time := ""
	if e.EntryDate != "" {
		time, err = utils.ParseEntryDate(e.EntryDate)
		if err != nil {
			log.Println("Error Parsing EntryDate:", err)
			utils.WriteError(w, http.StatusInternalServerError, "Unable to process request.")
		}
	}

	UpdatedEntry := models.EntryDto{
		Id:        "",
		Title:     e.Title,
		Content:   e.Content,
		EntryDate: time,
		UserId:    userId,
	}

	err = a.Entries.UpdateEntry(entryId, &UpdatedEntry)
	if err != nil {
		log.Println("Error Updating EntryDto:", err)
		utils.WriteError(w, http.StatusInternalServerError, "Unable to process request.")
		return
	}

	succMsg := "Entry " + entryId + " updated successfully."
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": succMsg})
}

func (a *Application) HandleDeleteEntry(w http.ResponseWriter, r *http.Request) {
	entryId := r.PathValue("id")
	if entryId == "" {
		utils.WriteError(w, http.StatusBadRequest, "Please provide a valid entry id.")
		return
	}
	ctx := r.Context()
	userId := ctx.Value(UserIdKey).(string)

	err := a.Entries.DeleteEntry(entryId, userId)
	if err != nil {
		log.Println("Error Deleting EntryDto:", err)
		utils.WriteError(w, http.StatusInternalServerError, "Unable to process request.")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Deleted entry " + entryId + " successfully."})
}
