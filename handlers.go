package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// Handlers

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

func steps(c echo.Context) error {
	data := Trial{}
	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	trialID := strconv.Itoa(data.ID)
	steps := getSteps(trialID)

	if steps == "-1" {
		return echo.NewHTTPError(http.StatusBadRequest, steps)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
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

func trials(c echo.Context) error {
	temp := getTrials()
	if len(temp) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, temp)
}

func getAllData(c echo.Context) error {
	temp := getTrialsComplete()
	if len(temp) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, temp)
}
