package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type GameMode string;
const (
  standard GameMode = "S"
  mania GameMode    = "SM"
  taiko GameMode    = "ST"
  catch GameMode    = "SC"
)

func main() {
  db_path :=  "./data/data.db";
  create_database(db_path);

  db, err := sql.Open("sqlite3", db_path);
  if err != nil {
    log.Fatalln(err);
  }
  defer db.Close();

  init_database(db);

  get_all_ranked_beatmap_ids_of_gamemode(standard, db);
  print_db(db);
}

func insert_beatmap_id(db *sql.DB, id int) {
  insertMapIDSQL := `INSERT OR REPLACE INTO beatmaps(id) VALUES (?)`
  statement, err := db.Prepare(insertMapIDSQL);
  if err != nil {
    log.Fatalln(err);
  }
  _, err = statement.Exec(id);
  if err != nil {
    log.Fatalln(err);
  }

}

func print_db(db *sql.DB) {
	row, err := db.Query("SELECT * FROM beatmaps")

	if err != nil {
		log.Fatal(err)
	}

	defer row.Close()
	for row.Next() {
		var id int
		row.Scan(&id)
		log.Println("Beatmap ID: ", id);
	}
}

func create_database(db_path string) {
  if _, err := os.Stat(db_path); err != nil {
    fmt.Println("Creating db...");
    os.Mkdir("./data", 0755);
    file, err := os.Create(db_path);
    if err != nil {
      log.Fatal(err);
    }
    file.Close();
    fmt.Println(db_path + " created.");
  } else {
    fmt.Println("db exists");
  }
}

func init_database(db *sql.DB) {
  const create string = `
  CREATE TABLE IF NOT EXISTS beatmaps (
  id INTEGER NOT NULL PRIMARY KEY
  );`;

  if _, err := db.Exec(create); err != nil {
    log.Fatalln(err);
  }

}

func get_all_ranked_beatmap_ids_of_gamemode(gm GameMode, db *sql.DB) {
  const uri = "https://osu.ppy.sh/beatmaps/packs/";
  for i := 1; true; i++ {
    var pack_url = uri + string(gm) + fmt.Sprint(i);

    ids, err := get_map_ids_from_pack_url(pack_url);
    if err != nil {
      log.Fatalln(err);
      break;
    }

    for _, id := range ids {
      insert_beatmap_id(db, id);
      fmt.Printf("Inserted beatmap ID: %d\n", id);
    }

    time.Sleep(4 * time.Second);
  }
}

func get_map_ids_from_pack_url(pack_url string) ([]int, error) {
  resp, err := http.Get(pack_url);

  if err != nil {
    return nil, errors.New("Error sending get request.")
  }

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, errors.New("Error reading response body.")
  }

  body_string := string(body);

  if strings.Contains(body_string, "Page Missing") {
    return nil, errors.New("Page not found.")
  }

  pattern := regexp.MustCompile(`https:\/\/osu\.ppy\.sh\/beatmapsets\/([^\/"]+)`);
  matches := pattern.FindAllStringSubmatch(body_string, -1);

  var ids []int;
  for _, match := range matches {
    id, err := strconv.Atoi(match[1])
    if err != nil {
      return nil, errors.New("Error converting beatmap ID to int.");
    }

    ids = append(ids, id);
  }
  
  return ids, nil;
}
