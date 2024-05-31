package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	sr "linear-db/pkg/structure"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func makeRequestCreateDatabase(dn string, posturl string) (*http.Request, error) {
	bodyString := fmt.Sprintf(`{
		"name": "%s"
	}`, dn)

	body := []byte(bodyString)
	r, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	return r, nil
}
func sendCreateDatabaseRequest(dn string, c *http.Client, t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("HTTP server: Error loading .env file", err)
	}
	posturl := fmt.Sprintf("http://%s/createdb", os.Getenv("URL_HTTP"))

	r, err := makeRequestCreateDatabase(dn, posturl)
	if err != nil {
		t.Fatalf("sending post failed %s\n", err)
	}
	res, err := c.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	//bodyString, _ := io.ReadAll(res.Body)
	//fmt.Println(string(bodyString))
}
func TestCreateDatabase(t *testing.T) {
	// tr := &http.Transport{
	// 	Proxy: http.ProxyFromEnvironment,
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   1000 * time.Second,
	// 		KeepAlive: 1000 * time.Second,
	// 	}).DialContext,
	// 	MaxIdleConnsPerHost:   100_000,
	// 	MaxIdleConns:          0,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// }

	client := &http.Client{Transport: &http.Transport{
		DisableKeepAlives: true,
	}}
	//debug.SetMemoryLimit(200 * 1 << 20)
	//var wg sync.WaitGroup

	for i := 0; i < 100_000; i++ {
		dbname := fmt.Sprintf("db%d", i)
		sendCreateDatabaseRequest(dbname, client, t)
		// 	wg.Add(1)
		// 	go func(i int) {
		// 		defer wg.Done()
		// 		dbname := fmt.Sprintf("db%d", i)
		// 		sendCreateDatabaseRequest(dbname, client, t)
		// 	}(i)
	}
	// wg.Wait()

}

func TestAppended(t *testing.T) {
	client := &http.Client{}
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("HTTP server: Error loading .env file", err)
	}
	geturl := fmt.Sprintf("http://%s/listdbs", os.Getenv("URL_HTTP"))
	res, _ := client.Get(geturl)
	defer res.Body.Close()
	databases := new(sr.Databases)
	b, _ := io.ReadAll(res.Body)
	json.Unmarshal(b, &databases.Databases)
	fmt.Println("appended:", len(databases.Databases))
}
