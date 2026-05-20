package bootstrapconnect

/*
Functions to interact with the bootstrap node of the network
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"bartering/utils"
)

func GetPeersFromBootstrapHTTP(IP string, port string) []string {

	/*
		Function to get peers from the bootstrap node via HTTP
		Arguments : IP of bootsrap as string, port as string
		Returns : bootstrap's response as string
	*/

	bootstrapUrl := IP + ":" + port

	bootstrapResponse, err := http.Get("http://" + bootstrapUrl)
	utils.ErrorHandler(err)

	defer func() { _ = bootstrapResponse.Body.Close() }()

	if bootstrapResponse.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status code:", bootstrapResponse.StatusCode)
		panic(-1)
	}

	bootstrapResponseBody, err := io.ReadAll(bootstrapResponse.Body)
	utils.ErrorHandler(err)

	var peers []string

	err = json.Unmarshal(bootstrapResponseBody, &peers)

	utils.ErrorHandler(err)

	return peers

}
