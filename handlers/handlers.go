package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josephmakin/monerochat/api"
	"github.com/josephmakin/monerochat/models"
	monerohub "github.com/josephmakin/monerohub/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CallbackURL string

type DonationsHandler struct {
	collection 		*mongo.Collection
	ctx 			context.Context
}

func NewDonationsHandler(ctx context.Context, collection *mongo.Collection) *DonationsHandler {
	return &DonationsHandler{
		collection: collection,
		ctx: 		ctx,
	}
}

func (handler *DonationsHandler) ListDonationsHandler (c *gin.Context) {
	cursor, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cursor.Close(handler.ctx)

	var donations []models.Donation
	for cursor.Next(handler.ctx) {
		var donation models.Donation
		if err := cursor.Decode(&donation); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		donations = append(donations, donation)
	}

	c.JSON(http.StatusOK, donations)
}

func (handler *DonationsHandler) GetOneDonationHandler (c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	var donation models.Donation
	err = handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectID,
	}).Decode(&donation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, donation)
}

func (handler *DonationsHandler) CreateOneDonationHandler (c *gin.Context) {
	var donation models.Donation

	if err := c.ShouldBindJSON(&donation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	payment, err := api.MakePaymentRequest(CallbackURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	donation.ID = payment.ID
	donation.Address = payment.Address
	donation.Timestamp = time.Now()
	donation.Transactions = []monerohub.Transaction{}

	_, err = handler.collection.InsertOne(handler.ctx, donation)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, donation)
}

func (handler *DonationsHandler) CallbackTransactionHandler (c *gin.Context) {
	var transaction monerohub.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	filter := bson.M{
		"address": transaction.Address,
		"transactions.txid": bson.M{"$ne" : transaction.TxID},
	}

	update := bson.M{
		"$addToSet": bson.M{
			"transactions": transaction,
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := handler.collection.UpdateOne(handler.ctx, filter, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var donation models.Donation
	err = handler.collection.FindOne(handler.ctx, bson.M{
		"address": transaction.Address,
	}).Decode(&donation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var total uint64
	for _, transaction := range donation.Transactions {
		total += transaction.Amount
	}

	if total >= donation.Amount {
		_, err := handler.collection.UpdateOne(
			handler.ctx,
			bson.M{"address": donation.Address},
			bson.M{"$set": bson.M{"paid": true}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, transaction)
}
