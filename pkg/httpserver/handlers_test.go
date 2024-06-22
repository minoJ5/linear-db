package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	sr "linear-db/pkg/structure"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/joho/godotenv"
)

func makeRequest(bodyString, posturl string) (*http.Request, error) {
	body := []byte(bodyString)
	r, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	return r, nil
}

func sendRequest(body, endpoint string, c *http.Client, t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("HTTP server: Error loading .env file", err)
	}
	posturl := fmt.Sprintf("http://%s/%s", os.Getenv("URL_HTTP"), endpoint)

	r, err := makeRequest(body, posturl)
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
	client := &http.Client{}
	//debug.SetMemoryLimit(200 * 1 << 20)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		// dbname := fmt.Sprintf("db%d", i)
		// sendCreateDatabaseRequest(dbname, client, t)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bodyString := fmt.Sprintf(`{
				"name": "db%d"
			}`, i)
			sendRequest(bodyString, "createdb", client, t)
		}(i)
	}
	wg.Wait()
}

func TestCreateTables(t *testing.T) {
	client := &http.Client{}
	//debug.SetMemoryLimit(200 * 1 << 20)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		// dbname := fmt.Sprintf("db%d", i)
		// sendCreateDatabaseRequest(dbname, client, t)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func(j int) {
					defer wg.Done()
					bodyString := fmt.Sprintf(`{
						"name": "table%d",
						"database_name": "db%d",
						"columns": [
								{
								"index": 0,
								"name": "col1",
								"type": "int",
								"values": [[1,4], 2, 3]
								},
								{
								"index": 1,
								"Name": "col2",
								"Type": "string",
								"Values": ["a", "b", "v"]
								}
						]
					}`, i, j)
					sendRequest(bodyString, "createtable", client, t)
				}(j)
			}
		}(i)
	}
	wg.Wait()
}

func TestDeleteDatabase(t *testing.T) {
	client := &http.Client{}
	//debug.SetMemoryLimit(200 * 1 << 20)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		// dbname := fmt.Sprintf("db%d", i)
		// sendCreateDatabaseRequest(dbname, client, t)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bodyString := fmt.Sprintf(`{
				"name": "db%d"
			}`, i)
			sendRequest(bodyString, "deletedb", client, t)
		}(i)
	}
	wg.Wait()
}

func TestDeleteTables(t *testing.T) {
	client := &http.Client{}
	//debug.SetMemoryLimit(200 * 1 << 20)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		// dbname := fmt.Sprintf("db%d", i)
		// sendCreateDatabaseRequest(dbname, client, t)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func(j int) {
					defer wg.Done()
					bodyString := fmt.Sprintf(`{
						"name": "table%d",
						"database_name": "db%d"
					}`, i, j)
					sendRequest(bodyString, "deletetable", client, t)
				}(j)
			}
		}(i)
	}
	wg.Wait()
}

func TestAppended(t *testing.T) {
	client := &http.Client{}
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("HTTP server: Error loading .env file", err)
	}
	geturl := fmt.Sprintf("http://%s/listdbs", os.Getenv("URL_HTTP"))
	res, err := client.Get(geturl)
	if err != nil {
		t.Fatalf("could not send get request: %s", err)
	}
	defer res.Body.Close()
	databases := new(sr.Databases)
	b, _ := io.ReadAll(res.Body)
	json.Unmarshal(b, &databases.Databases)
	fmt.Println("appended:", len(databases.Databases))
}
