package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/RameezAhmad007/go-rest-crud/internal/model"
	"github.com/RameezAhmad007/go-rest-crud/internal/repository"
	"github.com/RameezAhmad007/go-rest-crud/internal/response"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterCardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/card", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createCard(w, r)
		case http.MethodGet:
			getAllCards(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/card/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/card/")
		switch r.Method {
		case http.MethodGet:
			getCard(w, r, id)
		case http.MethodPut:
			updateCard(w, r, id)
		case http.MethodDelete:
			deleteCard(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func createCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create card endpoint hit")
	ctx := r.Context()
	var card model.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		writeResponse(w, response.CardResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    err.Error(),
		})
		return
	}

	newCard, err := repository.CreateCard(ctx, card)
	if err != nil {
		if err == mongo.ErrClientDisconnected {
			writeResponse(w, response.CardResponse{
				Status:  http.StatusConflict,
				Message: "error",
				Data:    fmt.Sprintf("Card with name '%s' already exists", card.Name),
			})
		} else {
			writeResponse(w, response.CardResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    err.Error(),
			})
		}
		return
	}

	writeResponse(w, response.CardResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    newCard,
	})
}

func getCard(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("get card endpoint hit")
	ctx := r.Context()
	card, err := repository.GetCard(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeResponse(w, response.CardResponse{
				Status:  http.StatusNotFound,
				Message: "error",
				Data:    "Card not found",
			})
		} else {
			writeResponse(w, response.CardResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    err.Error(),
			})
		}
		return
	}

	writeResponse(w, response.CardResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    card,
	})
}

func getAllCards(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all cards endpoint hit")
	ctx := r.Context()
	cards, err := repository.GetAllCards(ctx)
	if err != nil {
		writeResponse(w, response.CardResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    err.Error(),
		})
		return
	}

	writeResponse(w, response.CardResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    cards,
	})
}

func updateCard(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("update card endpoint hit")
	ctx := r.Context()
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		writeResponse(w, response.CardResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    err.Error(),
		})
		return
	}

	updatedCard, err := repository.UpdateCard(ctx, id, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeResponse(w, response.CardResponse{
				Status:  http.StatusNotFound,
				Message: "error",
				Data:    "Card not found",
			})
		} else {
			writeResponse(w, response.CardResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    err.Error(),
			})
		}
		return
	}

	writeResponse(w, response.CardResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    updatedCard,
	})
}

func deleteCard(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("Delete card endpoint hit")
	ctx := r.Context()
	deletedCount, err := repository.DeleteCard(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeResponse(w, response.CardResponse{
				Status:  http.StatusNotFound,
				Message: "error",
				Data:    "Card not found",
			})
		} else {
			writeResponse(w, response.CardResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    err.Error(),
			})
		}
		return
	}

	if deletedCount == 0 {
		writeResponse(w, response.CardResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    "Card not found",
		})
		return
	}

	writeResponse(w, response.CardResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    "Card deleted successfully",
	})
}

func writeResponse(w http.ResponseWriter, resp response.CardResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}
