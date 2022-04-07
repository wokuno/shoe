package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type SensorData struct {
	TrialID int `db:"trialid"`
	Time    int `db:"time"`
	A1      int `db:"a1"`
	A2      int `db:"a2"`
	A3      int `db:"a3"`
	A4      int `db:"a4"`
	A5      int `db:"a5"`
	LR      int `db:"LR"`
}

type Trial struct {
	ID         int          `db:"id"`
	StartTime  int          `db:"start_time"`
	SensorData []SensorData `json:"SensorData"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "shoe"
)

var (
	connectionString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	db               *sqlx.DB
)

func main() {
	var err error
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sqlx.Connect("pgx", connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// defer db.Close(context.Background())

	fmt.Println("Successfully connected!")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.File("/", "index.html")

	e.GET("/newTrial", newTrial)

	e.POST("/newData", newData)

	e.GET("/latestSteps", latestSteps)

	e.GET("/totalSteps", totalSteps)

	e.GET("/totalStepsO", totalStepsO)

	e.GET("/latestStepsO", latestStepsO)

	e.GET("/pronation", pronation)

	e.GET("/supination", supination)

	e.GET("/normal", normal)

	e.GET("/Tpronation", Tpronation)

	e.GET("/Tsupination", Tsupination)

	e.GET("/Tnormal", Tnormal)

	e.POST("/selectSteps", selectSteps)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func newTrial(c echo.Context) error {
	_, err := db.NamedExec(`INSERT INTO trials (start_time) VALUES (:time)`,
		map[string]interface{}{
			"time": time.Now().Unix(),
		})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	id := 0
	err = db.Get(&id, "SELECT id FROM trials WHERE id = (SELECT MAX(id) FROM trials)")

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.String(http.StatusOK, strconv.Itoa(id))
}

func newData(c echo.Context) error {
	u := []SensorData{}
	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&u)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	_, err = db.NamedExec(`INSERT INTO data (trialid, time, a1, a2, a3, a4, a5) VALUES (:trialid, :time, :a1, :a2, :a3, :a4, :a5)`, u)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.String(http.StatusOK, "1")
}

func getRows(trialid string) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}
	count := 0
	if trialid != "0" {
		err := db.Get(&count, "SELECT COUNT(*) FROM data", trialid)

		if err != nil {
			return -1
		}
	} else {
		err := db.Get(&count, "SELECT COUNT(*) FROM data WHERE trialid=$1", trialid)

		if err != nil {
			return -1
		}
	}

	return count
}

func getThreshold(trialid string) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}
	threshold := 0
	max := 0
	for i := 100; i < 500; i++ {
		steps := 0

		if trialid != "0" {
			err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data WHERE trialid=$2) AS data WHERE steps", i, trialid)

			if err != nil {
				return -1
			}
		} else {
			err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data) AS data WHERE steps", i)

			if err != nil {
				return -1
			}
		}

		if steps > max {
			max = steps
			threshold = i
		}
	}
	return threshold
}

func getPronation(trialid string, pThreshold int) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}

	threshold := getThreshold(trialid)

	pronation := 0

	if trialid != "0" {
		err := db.Get(&pronation, "SELECT COUNT(pronation) FROM (SELECT ((((a2 + a3) - (a1 + a5)) > $1) AND ((a1 + a2 + a3 + a5) > $2)) AS pronation FROM data WHERE trialid=$3) AS data WHERE pronation", strconv.Itoa(pThreshold), strconv.Itoa(threshold), trialid)
		if err != nil {
			return -1
		}
	} else {
		err := db.Get(&pronation, "SELECT COUNT(pronation) FROM (SELECT ((((a2 + a3) - (a1 + a5)) > $1) AND ((a1 + a2 + a3 + a5) > $2)) AS pronation FROM data) AS data WHERE pronation", strconv.Itoa(pThreshold), strconv.Itoa(threshold))
		if err != nil {
			return -1
		}
	}

	return pronation
}

func getSupination(trialid string, sThreshold int) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}

	threshold := getThreshold(trialid)

	supination := 0

	if trialid != "0" {
		err := db.Get(&supination, "SELECT COUNT(supination) FROM (SELECT ((((a2 + a3) - (a1 + a5)) < $1) AND ((a1 + a2 + a3 + a5) > $2)) AS supination FROM data WHERE trialid=$3) AS data WHERE supination", strconv.Itoa(sThreshold), strconv.Itoa(threshold), trialid)
		if err != nil {
			return -1
		}
	} else {
		err := db.Get(&supination, "SELECT COUNT(supination) FROM (SELECT ((((a2 + a3) - (a1 + a5)) < $1) AND ((a1 + a2 + a3 + a5) > $2)) AS supination FROM data) AS data WHERE supination", strconv.Itoa(sThreshold), strconv.Itoa(threshold))
		if err != nil {
			return -1
		}
	}

	return supination
}

func getNormal(trialid string, pThreshold int, sThreshold int) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}

	threshold := getThreshold(trialid)

	normal := 0

	if trialid != "0" {
		err := db.Get(&normal, "SELECT COUNT(normal) FROM (SELECT ((NOT((((a2 + a3) - (a1 + a5)) < $1) AND ((a1 + a2 + a3 + a5) > $3)) AND NOT((((a2 + a3) - (a1 + a5)) > $2) AND ((a1 + a2 + a3 + a5) > $3))) AND ((a1 + a2 + a3 + a5) > $3)) AS normal FROM data WHERE trialid=$4) AS data WHERE normal", strconv.Itoa(sThreshold), strconv.Itoa(pThreshold), strconv.Itoa(threshold), trialid)
		if err != nil {
			return -1
		}
	} else {
		err := db.Get(&normal, "SELECT COUNT(normal) FROM (SELECT ((NOT((((a2 + a3) - (a1 + a5)) < $1) AND ((a1 + a2 + a3 + a5) > $3)) AND NOT((((a2 + a3) - (a1 + a5)) > $2) AND ((a1 + a2 + a3 + a5) > $3))) AND ((a1 + a2 + a3 + a5) > $3)) AS normal FROM data) AS data WHERE normal", strconv.Itoa(sThreshold), strconv.Itoa(pThreshold), strconv.Itoa(threshold))
		if err != nil {
			return -1
		}
	}

	return normal
}

func getSteps(trialid string, threshold int) string {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return "-1"
		}
	}
	if trialid == "0" && threshold != -1 {
		steps := 0
		err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY $1) < $1)) AS steps FROM data) AS data WHERE steps", threshold)

		if err != nil {
			return "-1"
		}
		return strconv.Itoa(steps)
	}
	if threshold == -1 {
		max := 0
		for i := 100; i < 500; i++ {
			steps := 0

			if trialid != "0" {
				err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data WHERE trialid=$2) AS data WHERE steps", i, trialid)

				if err != nil {
					return "-1"
				}
			} else {
				err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data) AS data WHERE steps", i)

				if err != nil {
					return "-1"
				}
			}

			if steps > max {
				max = steps
			}
		}
		return strconv.Itoa(max)
	}
	steps := 0
	err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data WHERE trialid=$2) AS data WHERE steps", threshold, trialid)

	if err != nil {
		return "-1"
	}

	return strconv.Itoa(steps)
}

func selectSteps(c echo.Context) error {
	data := Trial{}
	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	trialID := strconv.Itoa(data.ID)
	steps := getSteps(trialID, -1)

	if steps == "-1" {
		return echo.NewHTTPError(http.StatusBadRequest, steps)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.String(http.StatusOK, steps)
}

func latestSteps(c echo.Context) error {
	steps := getSteps("-1", 420)
	if steps == "-1" {
		return echo.NewHTTPError(http.StatusBadRequest, steps)
	}

	return c.String(http.StatusOK, steps)
}

func latestStepsO(c echo.Context) error {
	steps := getSteps("-1", -1)
	if steps == "-1" {
		return echo.NewHTTPError(http.StatusBadRequest, steps)
	}
	return c.String(http.StatusOK, steps)
}

// SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > 420) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < 420)) AS steps FROM data) AS data WHERE steps
func totalSteps(c echo.Context) error {
	steps := getSteps("0", 420)
	if steps == "-1" {
		return echo.NewHTTPError(http.StatusBadRequest, steps)
	}

	return c.String(http.StatusOK, steps)
}
func totalStepsO(c echo.Context) error {
	steps := getSteps("0", -1)
	if steps == "-1" {
		return echo.NewHTTPError(http.StatusBadRequest, steps)
	}

	return c.String(http.StatusOK, steps)
}

func pronation(c echo.Context) error {
	pronation := getPronation("-1", 4200)
	if pronation == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, pronation)
	}

	return c.String(http.StatusOK, strconv.Itoa(pronation))
}

func supination(c echo.Context) error {
	supination := getSupination("-1", -200)
	if supination == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, supination)
	}

	return c.String(http.StatusOK, strconv.Itoa(supination))
}

func normal(c echo.Context) error {
	normal := getNormal("-1", 4200, -200)
	if normal == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, normal)
	}

	return c.String(http.StatusOK, strconv.Itoa(normal))
}

func Tpronation(c echo.Context) error {
	pronation := getPronation("0", 4200)
	if pronation == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, pronation)
	}

	return c.String(http.StatusOK, strconv.Itoa(pronation))
}

func Tsupination(c echo.Context) error {
	supination := getSupination("0", -200)
	if supination == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, supination)
	}

	return c.String(http.StatusOK, strconv.Itoa(supination))
}

func Tnormal(c echo.Context) error {
	normal := getNormal("0", 4200, -200)
	if normal == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, normal)
	}

	return c.String(http.StatusOK, strconv.Itoa(normal))
}

// INSERT INTO data (TrialID, Time, A1, A2, A3, A4, A5) Values ($1, $2, $3, $4, $5, $6, $7)
// SELECT * FROM data WHERE TrialID=$1 ORDER BY Time
// SELECT SUM(A1) FROM data WHERE TrialID=$1 ORDER BY TIME
//INSERT INTO trials (start_time) VALUES (select extract(epoch from now()))
// curl -d '{"trialid":5, "time":1, "a1":5, "a2":5, "a3":0, "a4":0, "a5":0 }' -H "Content-Type: application/json" -X POST http://localhost:1323/newData
// SELECT * FROM data WHERE trialid=(SELECT max(trialid) FROM data);
// SELECT (a1+a2+a3+a5) FROM data WHERE trialid=(SELECT max(trialid) FROM data);
// SELECT ((a5+a1)-(a2+a3)) FROM data WHERE trialid=(SELECT max(trialid) FROM data);
// SELECT ((a1+a2+a3+a5),((a5+a1)-(a2+a3))) FROM data WHERE trialid=(SELECT max(trialid) FROM data);
// SELECT (((a1 + a2 + a3 + a5) > 420) AND ((LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid)) < 420)) FROM data WHERE trialid=(SELECT max(trialid) FROM data);

// COPY (SELECT * FROM data WHERE trialid=(SELECT max(trialid) FROM data)) TO '/tmp/data.csv' WITH DELIMITER ',' CSV HEADER;
// rsync -vae "ssh" ~/Dropbox/Dev/School/EE4951W wokuno@104.248.49.139:~/
// rsync -vae "ssh" wokuno@104.248.49.139:/tmp/data.csv ~/Desktop
// curl -d '[{"trialid":5, "time":5, "a1":5, "a2":5, "a3":0, "a4":0, "a5":0 },{"trialid":5, "time":2, "a1":5, "a2":5, "a3":0, "a4":0, "a5":0 },{"trialid":5, "time":3, "a1":5, "a2":5, "a3":0, "a4":0, "a5":0 }]' -H "Content-Type: application/json" -X POST http://localhost:1323/newData
// SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > 420) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < 420)) AS steps FROM data WHERE trialid=(SELECT max(trialid) FROM data)) AS data WHERE steps;
// curl -d '{"id":19}' -H "Content-Type: application/json" -X POST http://104.248.49.139/selectSteps
// SELECT COUNT(pronation) FROM (SELECT ((((a2 + a3) - (a1 + a5)) > 4500) AND ((a1 + a2 + a3 + a5) > 420)) AS pronation FROM data WHERE trialid=26) AS data WHERE pronation;
// SELECT COUNT(normal) FROM (SELECT ((NOT((((a2 + a3) - (a1 + a5)) < -200) AND ((a1 + a2 + a3 + a5) > 420)) AND NOT((((a2 + a3) - (a1 + a5)) > 4200) AND ((a1 + a2 + a3 + a5) > 420))) AND ((a1 + a2 + a3 + a5) > 420)) AS normal FROM data WHERE trialid=26) AS data WHERE normal;
