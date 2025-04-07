package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RameezAhmad007/go-rest-crud/internal/api"
	"github.com/RameezAhmad007/go-rest-crud/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardHandler struct{}

func RegisterCardRoutes(mux *http.ServeMux) {
	r := chi.NewRouter()
	handler := &CardHandler{}
	api.HandlerFromMux(handler, r)
	mux.Handle("/", r)
}

func (h *CardHandler) PostCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create card endpoint hit")
	var card api.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		writeResponse(w, api.CardResponse{
			Status:  intPtr(http.StatusBadRequest),
			Message: stringPtr("error"),
			Data:    toMap(err.Error()),
		})
		return
	}

	newCard, err := repository.CreateCard(r.Context(), card)
	if err != nil {
		if err == mongo.ErrClientDisconnected {
			writeResponse(w, api.CardResponse{
				Status:  intPtr(http.StatusConflict),
				Message: stringPtr("error"),
				Data:    toMap(fmt.Sprintf("Card with name '%s' already exists", card.Name)),
			})
		} else {
			writeResponse(w, api.CardResponse{
				Status:  intPtr(http.StatusInternalServerError),
				Message: stringPtr("error"),
				Data:    toMap(err.Error()),
			})
		}
		return
	}

	writeResponse(w, api.CardResponse{
		Status:  intPtr(http.StatusCreated),
		Message: stringPtr("success"),
		Data:    toMap(newCard),
	})
}

func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all cards endpoint hit")
	cards, err := repository.GetAllCards(r.Context())
	if err != nil {
		writeResponse(w, api.CardResponse{
			Status:  intPtr(http.StatusInternalServerError),
			Message: stringPtr("error"),
			Data:    toMap(err.Error()),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}

func (h *CardHandler) GetCardId(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("get card endpoint hit")
	card, err := repository.GetCard(r.Context(), id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeResponse(w, api.CardResponse{
				Status:  intPtr(http.StatusNotFound),
				Message: stringPtr("error"),
				Data:    toMap("Card not found"),
			})
		} else {
			writeResponse(w, api.CardResponse{
				Status:  intPtr(http.StatusBadRequest),
				Message: stringPtr("error"),
				Data:    toMap(err.Error()),
			})
		}
		return
	}

	writeResponse(w, api.CardResponse{
		Status:  intPtr(http.StatusOK),
		Message: stringPtr("success"),
		Data:    toMap(card),
	})
}

func (h *CardHandler) PutCardId(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("update card endpoint hit")
	var updateData api.CardUpdate
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		writeResponse(w, api.CardResponse{
			Status:  intPtr(http.StatusBadRequest),
			Message: stringPtr("error"),
			Data:    toMap(err.Error()),
		})
		return
	}

	updateMap := make(map[string]interface{})
	if updateData.Name != nil {
		updateMap["name"] = *updateData.Name
	}
	if updateData.Number != nil {
		updateMap["number"] = *updateData.Number
	}

	updatedCard, err := repository.UpdateCard(r.Context(), id, updateMap)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeResponse(w, api.CardResponse{
				Status:  intPtr(http.StatusNotFound),
				Message: stringPtr("error"),
				Data:    toMap("Card not found"),
			})
		} else {
			writeResponse(w, api.CardResponse{
				Status:  intPtr(http.StatusInternalServerError),
				Message: stringPtr("error"),
				Data:    toMap(err.Error()),
			})
		}
		return
	}

	writeResponse(w, api.CardResponse{
		Status:  intPtr(http.StatusOK),
		Message: stringPtr("success"),
		Data:    toMap(updatedCard),
	})
}

func (h *CardHandler) DeleteCardId(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("Delete card endpoint hit")
	deletedCount, err := repository.DeleteCard(r.Context(), id)
	if err != nil || deletedCount == 0 {
		writeResponse(w, api.CardResponse{
			Status:  intPtr(http.StatusNotFound),
			Message: stringPtr("error"),
			Data:    toMap("Card not found"),
		})
		return
	}

	writeResponse(w, api.CardResponse{
		Status:  intPtr(http.StatusOK),
		Message: stringPtr("success"),
		Data:    toMap("Card deleted successfully"),
	})
}

func writeResponse(w http.ResponseWriter, resp api.CardResponse) {
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("writeResponse: Status is nil")
		return
	}
	if resp.Data == nil {
		fmt.Printf("writeResponse: Data is nil for status %d\n", *resp.Status)
	}
	w.WriteHeader(*resp.Status)
	json.NewEncoder(w).Encode(resp)
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func toMap(v interface{}) *map[string]interface{} {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case string:
		m := map[string]interface{}{"message": val}
		return &m
	default:
		data, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("toMap: failed to marshal %v: %v\n", v, err)
			return nil
		}
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			fmt.Printf("toMap: failed to unmarshal %v: %v\n", data, err)
			return nil
		}
		return &m
	}
}
