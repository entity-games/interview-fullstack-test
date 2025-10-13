package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"entity/interview/cmd/server/data"
	"entity/interview/cmd/server/utils"
)

type PurchaseRequest struct {
	GameID    string `json:"game_id"`
	ItemID    string `json:"item_id"`
	OrderType string `json:"order_type"`
}

func MakeVirtualPurchase(c *gin.Context, a utils.App) {
	ctx := c.Request.Context()

	sessionId := data.SessionReadIdFromCookie(c)
	sessionData := &data.SessionData{}
	err := data.SessionGet(c.Request.Context(), a.Redis, sessionData, sessionId)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Unable to retrieve session: %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session"})
	}

	request := PurchaseRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var cost int
	var gameId int

	if request.OrderType == data.OrderTypeAccess {
		gameId, _ = strconv.Atoi(request.ItemID)
		game := &data.Game{GameID: gameId}
		err = data.GamesGet(c, a.Postgres, game)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Unable to retrieve game: %s", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game"})
			return
		}

		if game.Cost == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Game cannot be purchased with coins"})
			return
		}

		cost = game.Cost
	} else if request.OrderType == data.OrderTypePlatformGoods {
		storeItem := &data.StoreItem{ItemID: request.ItemID}
		err = data.StoreItemsGet(c, a.Postgres, storeItem)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Unable to retrieve store item: %s", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item"})
			return
		}

		if storeItem.Cost == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item cannot be purchased with coins"})
			return
		}

		cost = storeItem.Cost
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect order type"})
		return
	}

	user := &data.User{UserID: sessionData.UserID}
	err = data.UsersGet(c, a.Postgres, user)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Unable to retrieve user: %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	err = data.UsersUpdateCoins(ctx, a.Postgres, user, -cost)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Unable to update coins: %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update the user"})
		return
	}

	order := &data.Order{
		ItemID: request.ItemID,
		UserID: sessionData.UserID,
		GameID: gameId,
		Type:   request.OrderType,
	}
	err = data.ApplyOrder(ctx, a.Postgres, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not apply goods"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}
