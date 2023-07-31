package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetDetail struct {
	AppraisedValue int    `json:"appraisedValue"`
	Color          string `json:"color"`
	ID             string `json:"id"`
	Owner          string `json:"owner"`
	Size           int    `json:"size"`
}

type AssetTransfer struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
}

func getAllDetailsHandler(c *gin.Context) {
	getAllAssets(cli)
	c.JSON(http.StatusOK, gin.H{})
}

func createTxnHandler(c *gin.Context) {
	var asset AssetDetail
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	createAsset(cli, asset.ID, asset.Color, strconv.Itoa(asset.Size), asset.Owner, strconv.Itoa(asset.AppraisedValue))
}

func readAssetHandler(c *gin.Context) {
	idParam := c.Param("id")
	readAssetByID(cli, idParam)
	c.JSON(http.StatusOK, gin.H{})
}

func transferAssetsHandler(c *gin.Context) {
	var newTxn AssetTransfer

	if err := c.ShouldBindJSON(&newTxn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	transferAssetAsync(cli, newTxn.ID, newTxn.Owner)

	c.JSON(http.StatusOK, gin.H{"Result": "Asset Transferred Successfully"})
}
