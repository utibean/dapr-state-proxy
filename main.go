package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

const ENV_KEY_DAPR_HTTP_PORT = "DAPR_HTTP_PORT"
const DEFAULT_DAPR_HTTP_PORT = "3500"
const ENV_KEY_STATE_STORE_NAME = "STATE_STORE_NAME"
const DEFAULT_STATE_STORE_NAME = "statestore"
const ENV_KEY_LOCAL_PORT = "APP_PORT"
const DEFAULT_LOCAL_PORT = "8999"

var daprPort = getEnvOrDefault(ENV_KEY_DAPR_HTTP_PORT, DEFAULT_DAPR_HTTP_PORT)
var stateStoreName = getEnvOrDefault(ENV_KEY_STATE_STORE_NAME, DEFAULT_STATE_STORE_NAME)
var localPort = getEnvOrDefault(ENV_KEY_LOCAL_PORT, DEFAULT_LOCAL_PORT)
var daprUrl = "http://localhost:" + daprPort + "/v1.0/state/" + stateStoreName

type stateItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	router := gin.Default()
	router.GET("/get/:key", getByKey)
	router.POST("/set", setState)
	router.Run("0.0.0.0:" + localPort)
}

func getByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	log.Println("Getting state for key " + key)
	response, err := http.Get(daprPort + "/" + key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
		return
	}
	valueBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.IndentedJSON(http.StatusOK, stateItem{
		Key:   key,
		Value: string(valueBytes),
	})
}

func setState(ctx *gin.Context) {

}

func getEnvOrDefault(key, defValue string) string {
	val, exist := os.LookupEnv(key)
	if !exist {
		return defValue
	}
	if val == "" {
		return defValue
	}
	return val
}
