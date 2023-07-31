package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Define your secret key for signing the JWT token
var cli *client.Contract

// User represents the user object
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	router := gin.Default()

	// Route for generating the JWT token

	// Authenticated route
	auth := router.Group("/api")
	auth.Use(authMiddleware()) // Apply authentication middleware

	auth.GET("/protected", protectedHandler)
	auth.GET("/getalldetails", getAllDetailsHandler)
	auth.POST("/submittxn", createTxnHandler)
	auth.GET("/readasset/:id", readAssetHandler)
	auth.GET("/transferassets", transferAssetsHandler)
	router.Run(":8080")
}

// Handler for the login route

// Middleware to check if the request is authenticated
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		clientConnection := newGrpcConnection()
		defer clientConnection.Close()

		id := newIdentity()
		sign := newSign()

		// Create a Gateway connection for a specific client identity
		gw, err := client.Connect(
			id,
			client.WithSign(sign),
			client.WithClientConnection(clientConnection),
			// Default timeouts for different gRPC calls
			client.WithEvaluateTimeout(5*time.Second),
			client.WithEndorseTimeout(15*time.Second),
			client.WithSubmitTimeout(5*time.Second),
			client.WithCommitStatusTimeout(1*time.Minute),
		)
		if err != nil {
			panic(err)
		}
		defer gw.Close()

		// Override default values for chaincode and channel name as they may differ in testing contexts.
		chaincodeName := "basic"
		if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
			chaincodeName = ccname
		}

		channelName := "mychannel"
		if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
			channelName = cname
		}

		network := gw.GetNetwork(channelName)
		contract := network.GetContract(chaincodeName)
		cli = contract
		c.Next()
	}
}

// Handler for the protected route
func protectedHandler(c *gin.Context) {
	initLedger(cli)
	c.JSON(http.StatusOK, gin.H{})
}
