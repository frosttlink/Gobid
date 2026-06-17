package api

import (
	"context"
	"net/http"

	"github.com/frosttlink/gobid/internal/jsonutils"
	"github.com/frosttlink/gobid/internal/services"
	"github.com/frosttlink/gobid/internal/usecase/product"
	"github.com/google/uuid"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductReq](r)
	if err != nil {
		_ = jsonutils.Encondejson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)
	if !ok {
		_ = jsonutils.Encondejson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected error, try again later",
		})
		return
	}

	productId, err := api.ProductService.CreateProduct(
		r.Context(),
		userID,
		data.ProductName,
		data.Description,
		data.Baseprice,
		data.AuctionEnd,
	)
	if err != nil {
		_ = jsonutils.Encondejson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product auction try again later",
		})
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)
	auctionRoom := services.NewAuctionRoom(ctx, productId, api.BidsService)

	go auctionRoom.Run()

	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[productId] = auctionRoom
	api.AuctionLobby.Unlock()

	_ = jsonutils.Encondejson(w, r, http.StatusCreated, map[string]any{
		"message":    "Auction has started with success",
		"product_id": productId,
	})
}
