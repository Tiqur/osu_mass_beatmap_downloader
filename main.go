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
  _ "database/sql"
  "os"
)

type GameMode string;
const (
  standard GameMode = "S"
  mania GameMode    = "SM"
  taiko GameMode    = "ST"
  catch GameMode    = "SC"
)

func main() {
  //get_all_ranked_beatmap_ids_of_gamemode(standard);
  init_database();
}

func init_database() {
  if _, err := os.Stat("./data/data.db"); err != nil {
    fmt.Println("Creating db");
    os.Mkdir("./data", 0755);
    os.Create("./data/data.db");
  } else {
    fmt.Println("db exists");
  }
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
