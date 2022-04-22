package main

import (
	"fmt"
	"os"

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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	e.Static("/", "assets")

	e.GET("/newTrial", newTrial)

	e.POST("/newData", newData)

	e.POST("/steps", steps)

	e.GET("/pronation", pronation)

	e.GET("/supination", supination)

	e.GET("/normal", normal)

	e.GET("/Tpronation", Tpronation)

	e.GET("/Tsupination", Tsupination)

	e.GET("/Tnormal", Tnormal)

	e.GET("/trials", trials)

	e.GET("/getAllData", getAllData)

	e.GET("/getOutliers", getOutliersRoute)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
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
// curl -d '{"id":28}' -H "Content-Type: application/json" -X POST https://trana.fit/selectSteps
// SELECT COUNT(pronation) FROM (SELECT ((((a2 + a3) - (a1 + a5)) > 4500) AND ((a1 + a2 + a3 + a5) > 420)) AS pronation FROM data WHERE trialid=26) AS data WHERE pronation;
// SELECT COUNT(normal) FROM (SELECT ((NOT((((a2 + a3) - (a1 + a5)) < -200) AND ((a1 + a2 + a3 + a5) > 420)) AND NOT((((a2 + a3) - (a1 + a5)) > 4200) AND ((a1 + a2 + a3 + a5) > 420))) AND ((a1 + a2 + a3 + a5) > 420)) AS normal FROM data WHERE trialid=26) AS data WHERE normal;
// SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > 420) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < 420)) AS steps FROM data WHERE trialid=28) AS data WHERE steps;
