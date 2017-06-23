package newStreamGenerate
import (
    "github.com/BurntSushi/toml"
    _ "github.com/lib/pq"
    "database/sql"
    "math"
    //"fmt"
)

const (
    EQUATORIAL_RADIUS    = 6378137.0            // 赤道半径 GRS80
    POLAR_RADIUS         = 6356752.314          // 極半径 GRS80
    ECCENTRICITY         = 0.081819191042815790 // 第一離心率 GRS80
    ECCENTRICITY2        = 0.006694380022900788 // 第二離心率 GRS80
)

type GeoCoder struct {
    db *sql.DB
    state *sql.Stmt
    cache []map[int]map[string]interface{}
    distance float64
    trueNum int
    num int
}

type Config struct {
    Database DbConfig
    Cache CacheConfig
}

type DbConfig struct {
    Use bool `toml:"use"`
    User string `toml:"user"`
    Name string `toml:"name"`
    Pass string `toml:"pass"`
}

type CacheConfig struct {
    Num int `toml:"num"`
    Distance float64 `toml:"distance"`
}

func ConnectDB(topology string) (*GeoCoder, error) {
    var config Config
    _, err := toml.DecodeFile("./config/" + topology + "/config.toml", &config)
    if err != nil {
        panic(err)
    }

    db, err := sql.Open("postgres", "user=" + config.Database.User + " dbname=" + config.Database.Name + " password=" + config.Database.Pass + " sslmode=disable")
    if err != nil {
        panic(err)
    }

    state, err := db.Prepare("SELECT ken_name,gst_name,css_name,moji FROM region WHERE ST_Within(ST_setSRID(ST_MakePoint($1,$2),4612),geom) LIMIT 1;")
    if err != nil {
        panic(err)
    }

    cache := make([]map[int]map[string]interface{}, config.Cache.Num)
    for i := 0; i < len(cache); i++ {
        cache[i] = map[int]map[string]interface{}{}
    }

    return &GeoCoder{
        db: db,
        state: state,
        cache: cache,
        distance: config.Cache.Distance,
        trueNum: 0,
        num: 0,
    }, nil
}

func (gc *GeoCoder) GeoCoding(id int, lon float64, lat float64) (string, string, string, string, error) {
    var region2, region1, city, prefecture string
    if m, ok := gc.cache[id % len(gc.cache)][id]; ok {
        if gc.IsMinDistance(lon, lat, m["lon"].(float64), m["lat"].(float64)) {
            region2 = m["region2"].(string)
            region1 = m["region1"].(string)
            city = m["city"].(string)
            prefecture = m["prefecture"].(string)

            //_, _, _, tprefecture, _ := gc.SearchDB(lon, lat)
            //tregion2, tregion1, tcity, tprefecture, _ := gc.SearchDB(lon, lat)
            //if prefecture == tprefecture {
            //if city == tcity && prefecture == tprefecture {
            //if region1 == tregion1 && city == tcity && prefecture == tprefecture {
            //if region2 == tregion2 && region1 == tregion1 && city == tcity && prefecture == tprefecture {
                //gc.trueNum = gc.trueNum + 1
                //gc.num = gc.num + 1
            //} else {
                //gc.num = gc.num + 1
            //}
            //fmt.Println(float64(gc.trueNum) / float64(gc.num))
        } else {
            region2, region1, city, prefecture, _ = gc.SearchDB(lon, lat)
            m["region2"] = region2
            m["region1"] = region1
            m["city"] = city
            m["prefecture"] = prefecture
            m["lon"] = lon
            m["lat"] = lat
            //gc.trueNum = gc.trueNum + 1
            //gc.num = gc.num + 1
            //fmt.Println(float64(gc.trueNum) / float64(gc.num))
        }
    } else {
        region2, region1, city, prefecture, _ = gc.SearchDB(lon, lat)
        cacheValue := map[string]interface{}{}

        cacheValue["region2"] = region2
        cacheValue["region1"] = region1
        cacheValue["city"] = city
        cacheValue["prefecture"] = prefecture
        cacheValue["lon"] = lon
        cacheValue["lat"] = lat

        gc.cache[id % len(gc.cache)][id] = cacheValue
        //gc.trueNum = gc.trueNum + 1
        //gc.num = gc.num + 1
        //fmt.Println(float64(gc.trueNum) / float64(gc.num))
    }

    return region2, region1, city, prefecture, nil
}

func (gc *GeoCoder) IsMinDistance(newLon float64, newLat float64, oldLon float64, oldLat float64) bool {
    dx := degree2radian(newLon - oldLon)
    dy := degree2radian(newLat - oldLat)
    my := degree2radian((oldLat + newLat) / 2)

    W := math.Sqrt(1 - (math.Pow(ECCENTRICITY, 2) * math.Pow(math.Sin(my), 2))) // 卯酉線曲率半径の分母
    m_numer := EQUATORIAL_RADIUS * (1 - math.Pow(ECCENTRICITY, 2))         // 子午線曲率半径の分子

    M := m_numer / math.Pow(W, 3) // 子午線曲率半径
    N := EQUATORIAL_RADIUS / W    // 卯酉線曲率半径

    d := math.Sqrt(math.Pow(dy * M, 2) + math.Pow(dx * N * math.Cos(my), 2))

    return d < gc.distance
}

func degree2radian(x float64) float64 {
    return x * math.Pi / 180
}

func (gc *GeoCoder) GeoCodingR(id int, lon float64, lat float64) (string, string, string, string, error) {
    var region2, region1, city, prefecture string
    if m, ok := gc.cache[id % len(gc.cache)][id]; ok {
        if gc.IsContainedRect(lon, lat, m) {
            region2 = m["region2"].(string)
            region1 = m["region1"].(string)
            city = m["city"].(string)
            prefecture = m["prefecture"].(string)

            //tregion2, tregion1, tcity, tprefecture, _ := gc.SearchDB(lon, lat)
            //if region2 == tregion2 && region1 == tregion1 && city == tcity && prefecture == tprefecture {
                //gc.trueNum = gc.trueNum + 1
                //gc.num = gc.num + 1
            //} else {
                //gc.num = gc.num + 1
            //}
            //fmt.Println(float64(gc.trueNum) / float64(gc.num))
        } else {
            region2, region1, city, prefecture, _ = gc.SearchDB(lon, lat)
            lon1, lon2, lat1, lat2 := gc.CreateRectAngle(lon, lat)
            m["region2"] = region2
            m["region1"] = region1
            m["city"] = city
            m["prefecture"] = prefecture
            m["lon1"] = lon1
            m["lon2"] = lon2
            m["lat1"] = lat1
            m["lat2"] = lat2

            //gc.trueNum = gc.trueNum + 1
            //gc.num = gc.num + 1
            //fmt.Println(float64(gc.trueNum) / float64(gc.num))
        }
    } else {
        region2, region1, city, prefecture, _ = gc.SearchDB(lon, lat)
        lon1, lon2, lat1, lat2 := gc.CreateRectAngle(lon, lat)
        cacheValue := map[string]interface{}{}

        cacheValue["region2"] = region2
        cacheValue["region1"] = region1
        cacheValue["city"] = city
        cacheValue["prefecture"] = prefecture
        cacheValue["lon1"] = lon1
        cacheValue["lon2"] = lon2
        cacheValue["lat1"] = lat1
        cacheValue["lat2"] = lat2

        gc.cache[id % len(gc.cache)][id] = cacheValue
        //gc.trueNum = gc.trueNum + 1
        //gc.num = gc.num + 1
        //fmt.Println(float64(gc.trueNum) / float64(gc.num))
    }

    return region2, region1, city, prefecture, nil
}

func (gc *GeoCoder) IsContainedRect(lon float64, lat float64, m map[string]interface{}) bool {
    return lon < m["lon1"].(float64) && lon > m["lon2"].(float64) && lat < m["lat1"].(float64) && lat > m["lat2"].(float64)
}

func (gc *GeoCoder) CreateRectAngle(lon float64, lat float64) (float64, float64, float64, float64) {
    wt := math.Sqrt(1.0 - ECCENTRICITY2 * math.Pow(math.Sin(lat * math.Pi / 180), 2))
    mt := EQUATORIAL_RADIUS * (1.0 - ECCENTRICITY2) / math.Pow(wt, 3)
    dit1 := math.Sqrt2 * gc.distance * math.Cos(45 * math.Pi / 180) / mt
    dit2 := math.Sqrt2 * gc.distance * math.Cos(225 * math.Pi / 180) / mt

    i1 := lat * math.Pi / 180 + dit1 / 2
    i2 := lat * math.Pi / 180 + dit2 / 2

    w1 := math.Sqrt(1.0 - ECCENTRICITY2 * math.Pow(math.Sin(i1), 2))
    w2 := math.Sqrt(1.0 - ECCENTRICITY2 * math.Pow(math.Sin(i2), 2))

    lon1, lon2 := gc.getLon(lon, lat, w1, w2, i1, i2)
    lat1, lat2 := gc.getLat(lon, lat, w1, w2)

    return lon1, lon2, lat1, lat2
}

func (gc *GeoCoder) getLon(lon float64, lat float64, w1 float64, w2 float64, i1 float64, i2 float64) (float64, float64) {
    n1 := EQUATORIAL_RADIUS / w1
    n2 := EQUATORIAL_RADIUS / w2

    dk1 := math.Sqrt2 * gc.distance * math.Sin(45 * math.Pi / 180) / (n1 * math.Cos(i1))
    dk2 := math.Sqrt2 * gc.distance * math.Sin(225 * math.Pi / 180) / (n2 * math.Cos(i2))

    return (lon + dk1 * 180 / math.Pi), (lon + dk2 * 180 / math.Pi)
}

func (gc *GeoCoder) getLat(lon float64, lat float64, w1 float64, w2 float64) (float64, float64) {
    m1 := EQUATORIAL_RADIUS * (1.0 - ECCENTRICITY2)/ math.Pow(w1, 3)
    m2 := EQUATORIAL_RADIUS * (1.0 - ECCENTRICITY2)/ math.Pow(w2, 3)

    di1 := math.Sqrt2 * gc.distance * math.Cos(45 * math.Pi / 180) / m1
    di2 := math.Sqrt2 * gc.distance * math.Cos(225 * math.Pi / 180) / m2

    return (lat + di1 * 180 / math.Pi), (lat + di2 * 180 / math.Pi)
}

func (gc *GeoCoder) SearchDB(lon float64, lat float64) (string, string, string, string, error) {
    var region2 interface{}
    var region1 interface{}
    var city interface{}
    var prefecture interface{}
    empty := ""

    row, err := gc.state.Query(lon, lat)

    defer row.Close()
    if err != nil {
        return empty, empty, empty, empty, nil
    }

    ok := row.Next()
    if !ok {
        return empty, empty, empty, empty, nil
    }

    err = row.Scan(&prefecture, &city, &region1, &region2)
    if err != nil {
        panic(err)
    }

    region2S, _ := region2.(string)
    region1S, _ := region1.(string)
    cityS, _ := city.(string)
    prefectureS, _ := prefecture.(string)

    return region2S, region1S, cityS, prefectureS, nil
}

func (gc *GeoCoder) Close() {
    gc.db.Close()
}
