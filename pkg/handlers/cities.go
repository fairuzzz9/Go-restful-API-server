package handlers

import (
	"context"
	"fmt"
	"go-skeleton-rest-app/internal/db"
	"go-skeleton-rest-app/pkg/http/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CityData struct {
	CityName  string `json:"city_name" validate:"required"`
	CountryID int    `json:"country_id" validate:"required"`
}

func ListCities(c echo.Context) error {

	// serverTraceID := c.Response().Header().Get(echo.HeaderXRequestID)

	// uncomment these lines to pass context to child function
	//ctx := c.Request().Context()
	//ctx = contextkeys.SetContextValue(ctx, contextkeys.CONTEXT_KEY_SERVER_TRACE_ID, serverTraceID)
	//ctx = contextkeys.SetContextValue(ctx, contextkeys.CONTEXT_KEY_CLIENT_REQUEST_ID, clientRequestID)

	// reply with SuccessCode
	reply, err := responses.GetReponseMessageByCode(responses.SuccessCode)

	if err != nil {

		// override the reply message with the error message
		reply.Message = "GetReponseMessageByCode(responses.SuccessCode) error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)
	}

	// todo: get cities from database
	ctx := context.Background()
	conn, _ := db.DBPool().Acquire(ctx)

	defer conn.Release()

	var cities []CityData
	rows, err := conn.Query(context.Background(), db.GetSQLByName("GetAllFromCities"))
	if err != nil {

		reply.Message = "error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)

	}

	for rows.Next() {
		var city CityData
		err := rows.Scan(&city.CityName, &city.CountryID)
		if err != nil {
			fmt.Println(err)
		}
		cities = append(cities, city)
	}

	// override the success message
	reply.Message = reply.Message + ". List all the cities"
	reply.Data = cities

	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, reply)
}

func ListCitiesByCountryId(c echo.Context) error {

	id := c.Param("id")
	// return c.String(http.StatusOK, id)

	reply, err := responses.GetReponseMessageByCode(responses.SuccessCode)

	if err != nil {

		// override the reply message with the error message
		reply.Message = "GetReponseMessageByCode(responses.SuccessCode) error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)
	}

	// todo: get cities from database
	ctx := context.Background()
	conn, _ := db.DBPool().Acquire(ctx)

	defer conn.Release()

	var cities []CityData
	rows, err := conn.Query(context.Background(), db.GetSQLByName("GetCitiesByCountryID"), id)
	if err != nil {

		reply.Message = "error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)

	}

	for rows.Next() {
		var city CityData
		err := rows.Scan(&city.CityName, &city.CountryID)
		if err != nil {
			fmt.Println(err)
		}
		cities = append(cities, city)
	}

	// override the success message
	reply.Message = reply.Message + ". List the city by id"
	reply.Data = cities

	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, reply)
}

func CreateCity(c echo.Context) error {

	id := c.Param("id")
	city := c.Param("city")

	reply, err := responses.GetReponseMessageByCode(responses.SuccessCode)

	if err != nil {

		// override the reply message with the error message
		reply.Message = "GetReponseMessageByCode(responses.SuccessCode) error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)
	}

	// todo: create city to database
	ctx := context.Background()
	conn, _ := db.DBPool().Acquire(ctx)

	defer conn.Release()

	var cities []CityData
	rows, err := conn.Query(context.Background(), db.GetSQLByName("CreateCity"), city, id)
	if err != nil {

		reply.Message = "error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)

	}

	for rows.Next() {
		var city CityData
		err := rows.Scan(&city.CityName, &city.CountryID)
		if err != nil {
			fmt.Println(err)
		}
		cities = append(cities, city)
	}

	// override the success message
	reply.Message = reply.Message + ". The city has been added to database"
	reply.Data = cities

	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, reply)
}

func UpdateCity(c echo.Context) error {

	id := c.Param("id")
	city := c.Param("city")

	reply, err := responses.GetReponseMessageByCode(responses.SuccessCode)

	if err != nil {

		// override the reply message with the error message
		reply.Message = "GetReponseMessageByCode(responses.SuccessCode) error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)
	}

	// todo: create city to database
	ctx := context.Background()
	conn, _ := db.DBPool().Acquire(ctx)

	defer conn.Release()

	var cities []CityData
	rows, err := conn.Query(context.Background(), db.GetSQLByName("UpdateCity"), city, id)
	if err != nil {

		reply.Message = "error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)

	}

	for rows.Next() {
		var city CityData
		err := rows.Scan(&city.CityName, &city.CountryID)
		if err != nil {
			fmt.Println(err)
		}
		cities = append(cities, city)
	}

	// override the success message
	reply.Message = reply.Message + ". The city has been updated to database"
	reply.Data = cities

	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, reply)

}

func DeleteCity(c echo.Context) error {

	city := c.Param("city")

	reply, err := responses.GetReponseMessageByCode(responses.SuccessCode)

	if err != nil {

		// override the reply message with the error message
		reply.Message = "GetReponseMessageByCode(responses.SuccessCode) error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)
	}

	// todo: create city to database
	ctx := context.Background()
	conn, _ := db.DBPool().Acquire(ctx)

	defer conn.Release()

	var cities []CityData
	rows, err := conn.Query(context.Background(), db.GetSQLByName("DeleteCity"), city)
	if err != nil {

		reply.Message = "error : " + err.Error()

		c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, reply)

	}

	for rows.Next() {
		var city CityData
		err := rows.Scan(&city.CityName, &city.CountryID)
		if err != nil {
			fmt.Println(err)
		}
		cities = append(cities, city)
	}

	// override the success message
	reply.Message = reply.Message + ". The city has been deleted from database"
	reply.Data = cities

	c.Response().WriteHeader(http.StatusOK)
	return c.JSON(http.StatusOK, reply)

}
