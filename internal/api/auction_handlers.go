package api

import (
	"errors"
	"net/http"

	"github.com/frosttlink/gobid/internal/jsonutils"
	"github.com/frosttlink/gobid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *Api) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductId := chi.URLParam(r, "product_id")

	producId, err := uuid.Parse(rawProductId)
	if err != nil {
		_ = jsonutils.Encondejson(w, r, http.StatusBadRequest, map[string]any{
			"message": "invalid product id - must be a valid uuid",
		})
		return
	}

	_, err = api.ProductService.GetProductById(r.Context(), producId)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			_ = jsonutils.Encondejson(w, r, http.StatusNotFound, map[string]any{
				"message": "no product with given id",
			})
			return
		}
		_ = jsonutils.Encondejson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "unexpected error, try again later",
		})
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)
	if !ok {
		_ = jsonutils.Encondejson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "unexpected error, try again later",
		})
	}

	conn, err := api.WsUpgradedr.Upgrade(w, r, nil)
	if err != nil {
		_ = jsonutils.Encondejson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "could not upgrade connection to a websocker protocol",
		})
	}

}
