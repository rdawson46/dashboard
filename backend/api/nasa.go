package api

import (
	"encoding/json"
	"fmt"
	"net/http"
    "os"
    "io"
)

var api_key = "DEMO_KEY"

type ApodResponse struct {
    URL string `json:"url"`
    MediaType string `json:"media_type"`
    Title string `json:"title"`
}

func Get_image_of_the_day() error {
    apiUrl := fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", api_key)

    resp, err := http.Get(apiUrl)

    fmt.Println(resp.StatusCode)
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    var apod ApodResponse 
    if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
        // TODO: move this print to the main logger when connected up
        fmt.Println("JSON Decode error:", err)
        return err
    }

    if apod.MediaType != "image" {
        // TODO: move this print to the main logger when connected up
        // and create an error for this
        fmt.Println("Media type is not image", apod.MediaType)
        fmt.Println("title", apod.Title)
        fmt.Println("URL", apod.URL)
        return nil
    }

    fmt.Println("Downloading image:", apod.URL)
    imageResp, err := http.Get(apod.URL)

    if err != nil {
        // TODO: move this print to the main logger when connected up
        fmt.Println("Failed to download image:", err)
        return err
    }

    defer imageResp.Body.Close()

    outFile, err := os.Create("apod_image.jpg")
    if err != nil {
        fmt.Println("Failed making file", err)
        return err
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, imageResp.Body)
    if err != nil {
        fmt.Println("Error when copying", err)
        return err
    }

    return nil
}
