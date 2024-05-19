// handler/person.go
package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goapi/model"
	"net/http"
)

var db *sql.DB

// SetDB sets the database connection
func SetDB(database *sql.DB) {
	db = database
}

// getPeople godoc
// @Summary Get all people
// @Description Get all people
// @Tags people
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Person
// @Router /people [get]
func GetPeople(c *gin.Context) {
	var people []model.Person
	rows, err := db.Query("SELECT id, name, email FROM people")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var person model.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		people = append(people, person)
	}

	c.JSON(http.StatusOK, people)
}

// getPerson godoc
// @Summary Get a person by ID
// @Description Get a person by ID
// @Tags people
// @Accept  json
// @Produce  json
// @Param id path int true "Person ID"
// @Success 200 {object} model.Person
// @Router /people/{id} [get]
func GetPerson(c *gin.Context) {
	id := c.Param("id")

	var person model.Person
	err := db.QueryRow("SELECT id, name, email FROM people WHERE id = @p1", id).Scan(&person.ID, &person.Name, &person.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Person not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, person)
}

// createPerson godoc
// @Summary Create a new person
// @Description Create a new person
// @Tags people
// @Accept  json
// @Produce  json
// @Param person body model.Person true "Person"
// @Success 201 {object} model.Person
// @Router /people [post]
func CreatePerson(c *gin.Context) {
	var person model.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow("INSERT INTO people (name, email) OUTPUT INSERTED.id VALUES (@p1, @p2)", person.Name, person.Email).Scan(&person.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, person)
}

// updatePerson godoc
// @Summary Update a person by ID
// @Description Update a person by ID
// @Tags people
// @Accept  json
// @Produce  json
// @Param id path int true "Person ID"
// @Param person body model.Person true "Person"
// @Success 200 {object} model.Person
// @Router /people/{id} [put]
func UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	var person model.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE people SET name = @p1, email = @p2 WHERE id = @p3", person.Name, person.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

// deletePerson godoc
// @Summary Delete a person by ID
// @Description Delete a person by ID
// @Tags people
// @Accept  json
// @Produce  json
// @Param id path int true "Person ID"
// @Success 204
// @Router /people/{id} [delete]
func DeletePerson(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM people WHERE id = @p1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
