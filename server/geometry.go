package server

import (
    "github.com/paulmach/go.geojson"
    "github.com/BurntSushi/toml"
    _ "github.com/lib/pq"
    "database/sql"
    //"fmt"
)

type GeoDB struct {
    db *sql.DB
    state *sql.Stmt
    stateCity *sql.Stmt
    statePrefecture *sql.Stmt
}

type Config struct {
    Database DbConfig
}

type DbConfig struct {
    User string `toml:"user"`
    Name string `toml:"name"`
    Pass string `toml:"pass"`
}

func ConnectDB() (*GeoDB, error) {
    var config Config
    _, err := toml.DecodeFile("./config/dbConfig.toml", &config)
    if err != nil {
        panic(err)
    }

    db, err := sql.Open("postgres", "user=" + config.Database.User + " dbname=" + config.Database.Name + " password=" + config.Database.Pass + " sslmode=disable")
    if err != nil {
        panic(err)
    }

    // state1, err := db.Prepare("SELECT ST_AsGeoJSON(geom) FROM region WHERE ken_name=$1 AND gst_name=$2 AND css_name=$3 AND moji=$4 LIMIT 1;")
    // state2, err := db.Prepare("SELECT ST_AsGeoJSON(geom) FROM region WHERE ken_name=$1 AND gst_name=$2 AND moji=$3 LIMIT 1;")

    // state, err := db.Prepare("SELECT row_to_json(featurecollection) FROM (SELECT 'FeatureCollection' AS type, array_to_json(array_agg(feature)) AS features FROM ( SELECT 'Feature' AS type, ST_AsGeoJSON(geom)::json AS geometry, row_to_json((SELECT p FROM (SELECT ken_name AS prefecture, gst_name AS city, css_name AS region1, moji AS region2) AS p)) AS properties FROM region WHERE ST_Intersects(ST_setSRID(ST_MakeBox2D(ST_MakePoint($1, $2), ST_MakePoint($3, $4)), 4612), geom)) AS feature) AS featurecollection;")
    state, err := db.Prepare("SELECT ST_AsGeoJSON(geom), ken_name AS prefecture, gst_name AS city, css_name AS region1, moji AS region2 FROM region WHERE ST_Intersects(ST_setSRID(ST_MakeBox2D(ST_MakePoint($1, $2), ST_MakePoint($3, $4)), 4612), geom);")
    if err != nil {
      panic(err)
    }

    stateCity, err := db.Prepare("SELECT ST_AsGeoJSON(geom), ken_name AS prefecture, gst_name AS city FROM city WHERE ST_Intersects(ST_setSRID(ST_MakeBox2D(ST_MakePoint($1, $2), ST_MakePoint($3, $4)), 4612), geom);")
    if err != nil {
      panic(err)
    }

    statePrefecture, err := db.Prepare("SELECT ST_AsGeoJSON(geom), ken_name AS prefecture FROM prefecture WHERE ST_Intersects(ST_setSRID(ST_MakeBox2D(ST_MakePoint($1, $2), ST_MakePoint($3, $4)), 4612), geom);")
    if err != nil {
      panic(err)
    }

    return &GeoDB{
        db: db,
        state: state,
        stateCity: stateCity,
        statePrefecture: statePrefecture,
    }, nil
}

type FeatureCollection struct {
  Type string `json:"type"`
  Features []*Feature `json:"features"`
}

type Feature struct {
  Type string `json:"type"`
  Geometry *geojson.Geometry `json:"geometry"`
  Properties map[string]string `json:"properties"`
}

func (gc *GeoDB) GetBoundedArea(southWestLon float64, southWestLat float64, northEastLon float64, northEastLat float64) (*FeatureCollection) {
    var features []*Feature
    row, err := gc.state.Query(southWestLon, southWestLat, northEastLon, northEastLat)
    if err != nil {
        panic(err)
    }
    defer row.Close()

    for {
        ok := row.Next()
        if !ok {
            break
        }

        var geometry *geojson.Geometry
        var prefecture *string
        var city *string
        var region1 *string
        var region2 *string
        var sregion1 string
        var sregion2 string

        err := row.Scan(&geometry, &prefecture, &city, &region1, &region2)
        if err != nil {
            panic(err)
        }

        if region1 == nil {
          sregion1 = ""
        } else {
          sregion1 = *region1
        }

        if region2 == nil {
          sregion2 = ""
        } else {
          sregion2 = *region2
        }

        features = append(features, &Feature{
          Type: "Feature",
          Geometry: geometry,
          Properties: map[string]string{
            "prefecture": *prefecture,
            "city": *city,
            "region1": sregion1,
            "region2": sregion2,
          },
        })
    }

    return &FeatureCollection{
      Type: "FeatureCollection",
      Features: features,
    }
}

func (gc *GeoDB) GetBoundedAreaCity(southWestLon float64, southWestLat float64, northEastLon float64, northEastLat float64) (*FeatureCollection) {
    var features []*Feature
    row, err := gc.stateCity.Query(southWestLon, southWestLat, northEastLon, northEastLat)
    if err != nil {
        panic(err)
    }
    defer row.Close()

    for {
        ok := row.Next()
        if !ok {
            break
        }

        var geometry *geojson.Geometry
        var prefecture *string
        var city *string

        err := row.Scan(&geometry, &prefecture, &city)
        if err != nil {
            panic(err)
        }

        features = append(features, &Feature{
          Type: "Feature",
          Geometry: geometry,
          Properties: map[string]string{
            "prefecture": *prefecture,
            "city": *city,
          },
        })
    }

    return &FeatureCollection{
      Type: "FeatureCollection",
      Features: features,
    }
}

func (gc *GeoDB) GetBoundedAreaPrefecture(southWestLon float64, southWestLat float64, northEastLon float64, northEastLat float64) (*FeatureCollection) {
    var features []*Feature
    row, err := gc.statePrefecture.Query(southWestLon, southWestLat, northEastLon, northEastLat)
    if err != nil {
        panic(err)
    }
    defer row.Close()

    for {
        ok := row.Next()
        if !ok {
            break
        }

        var geometry *geojson.Geometry
        var prefecture *string

        err := row.Scan(&geometry, &prefecture)
        if err != nil {
            panic(err)
        }

        features = append(features, &Feature{
          Type: "Feature",
          Geometry: geometry,
          Properties: map[string]string{
            "prefecture": *prefecture,
          },
        })
    }

    return &FeatureCollection{
      Type: "FeatureCollection",
      Features: features,
    }
}

// func (gc *GeoDB) SearchDB(prefecture string, city string, region1 string, region2 string) (*geojson.Geometry, error) {
//     var geometry *geojson.Geometry
//     var row *sql.Rows
//     if region1 == "" {
//       row, _ = gc.state2.Query(prefecture, city, region2)

//     } else {
//       row, _ = gc.state1.Query(prefecture, city, region1, region2)

//     }

//    defer row.Close()

//     ok := row.Next()
//     if !ok {
//       panic(ok)
//     }

//     err := row.Scan(&geometry)
//     if err != nil {
//         panic(err)
//     }

//     return geometry, nil
// }
