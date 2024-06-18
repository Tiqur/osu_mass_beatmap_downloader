package main

import (
  "log"
  "fmt"
  "net/http"
  "regexp"
  "io"
  "time"
  "errors"
  "strings"
)

type GameMode string;
const (
  standard GameMode = "S"
  mania GameMode    = "SM"
  taiko GameMode    = "ST"
  catch GameMode    = "SC"
)

func main() {
  get_all_ranked_beatmap_ids_of_gamemode(standard);
}

func get_all_ranked_beatmap_ids_of_gamemode(gm GameMode) {
  const uri = "https://osu.ppy.sh/beatmaps/packs/";
  for i := 1; true; i++ {
    var pack_url = uri + string(gm) + fmt.Sprint(i);

    err := get_map_ids_from_pack_url(pack_url);
    if err != nil {
      log.Fatalln(err);
      break;
    }

    time.Sleep(4 * time.Second);
  }
}

func get_map_ids_from_pack_url(pack_url string) error {
  resp, err := http.Get(pack_url);

  if err != nil {
    return errors.New("Error sending get request.")
  }

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return errors.New("Error reading response body.")
  }

  body_string := string(body);

  if strings.Contains(body_string, "Page Missing") {
    return errors.New("Page not found.")
  }

  pattern := regexp.MustCompile(`https:\/\/osu\.ppy\.sh\/beatmapsets\/([^\/"]+)`);
  matches := pattern.FindAllStringSubmatch(body_string, -1);

  for _, match := range matches {
    fmt.Println("Captured segment:", match[0]);
  }
  
  return nil;
}
