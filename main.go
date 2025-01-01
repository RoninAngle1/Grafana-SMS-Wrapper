package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"
)

// Config holds the configuration data (API URL, server settings, etc.)
type Config struct {
        Server struct {
                IP   string `json:"ip"`
                Port string `json:"port"`
                Path string `json:"path"`
        } `json:"server"`
        PhoneNumbers []string `json:"phone_numbers"`
        URL          string   `json:"url"`
        UserName     string   `json:"username"`
        Password     string   `json:"password"`
}

// GrafanaAlert defines the structure for the incoming Grafana webhook payload
type GrafanaAlert struct {
        Alerts []struct {
                Title       string `json:"title"`
                Annotations struct {
                        Description string `json:"description"`
                } `json:"annotations"`
        } `json:"alerts"`
}

func main() {
        // Read configuration file
        configFile := "./config.json"
        config, err := readConfig(configFile)
        if err != nil {
                log.Fatalf("Error reading config: %v", err)
        }

        // Set up the HTTP handler for Grafana webhooks
        http.HandleFunc(config.Server.Path, func(w http.ResponseWriter, r *http.Request) {
                if r.Method == http.MethodPost {
                        // Parse incoming JSON message (Grafana webhook)
                        var alertPayload GrafanaAlert
                        if err := json.NewDecoder(r.Body).Decode(&alertPayload); err != nil {
                                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                                return
                        }

                        // Process each alert
                        for _, alert := range alertPayload.Alerts {
                                // Create the message by combining title and description
                                message := fmt.Sprintf("Title: %s\nDescription: %s", alert.Title, alert.Annotations.Description)

                                // Send the message to the destination API, passing each phone number individually
                                for _, phoneNumber := range config.PhoneNumbers {
                                        if err := sendMessage(config, message, phoneNumber); err != nil {
                                                http.Error(w, fmt.Sprintf("Error sending message: %v", err), http.StatusInternalServerError)
                                                return
                                        }
                                }
                        }

                        // Respond to the sender
                        w.WriteHeader(http.StatusOK)
                        w.Write([]byte("Alert titles and descriptions sent successfully"))
                } else {
                        http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
                }
        })

        // Start the server with the IP and port from config
        serverAddress := fmt.Sprintf("%s:%s", config.Server.IP, config.Server.Port)
        log.Printf("Starting server at %s...", serverAddress)
        log.Fatal(http.ListenAndServe(serverAddress, nil))
}

// readConfig reads the configuration file
func readConfig(filePath string) (Config, error) {
        file, err := os.Open(filePath)
        if err != nil {
                return Config{}, err
        }
        defer file.Close()

        var config Config
        decoder := json.NewDecoder(file)
        if err := decoder.Decode(&config); err != nil {
                return Config{}, err
        }

        return config, nil
}

// sendMessage sends the title-description message along with a single phone number
func sendMessage(config Config, message, phoneNumber string) error {
        // Create the JSON body for the request containing "Message" and "PhoneNumber"
        body := map[string]string{
                "Message":   message,
                "PhoneNumber": phoneNumber,
        }

        // Marshal the body into JSON
        jsonBody, err := json.Marshal(body)
        if err != nil {
                return fmt.Errorf("failed to marshal message: %v", err)
        }

        // Log the JSON body to the console (before sending it)
        log.Printf("Sending request to destination API with body: %s", string(jsonBody))

        // Create the HTTP POST request
        req, err := http.NewRequest(http.MethodPost, config.URL, bytes.NewBuffer(jsonBody))
        if err != nil {
                return fmt.Errorf("failed to create HTTP request: %v", err)
        }

        // Set the content-type header
        req.Header.Set("Content-Type", "application/json")

        // Add the custom headers (UserName and Password)
        req.Header.Set("UserName", config.UserName)
        req.Header.Set("Password", config.Password)

        // Send the request using http.DefaultClient
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                return fmt.Errorf("failed to send HTTP request: %v", err)
        }
        defer resp.Body.Close()

        // Read the response (optional)
        respBody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return fmt.Errorf("failed to read response body: %v", err)
        }

        // Log the response from the API for debugging purposes
        log.Printf("Response from destination API: %s", string(respBody))

        // Check for a successful response (200 OK)
        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("received non-OK response: %d %s", resp.StatusCode, string(respBody))
        }

        return nil
}

root@nexus-dariche:/mnt#
