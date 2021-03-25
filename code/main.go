package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create the pet store.
	myPetStore := petStore{}

	r := gin.New()

	// List all available pets.
	r.GET("/pets", myPetStore.listPets)

	// Add one pet to the list of available pets.
	r.PUT("/pet", myPetStore.putPet)

	// Get one pet from the inventory of available pets.
	r.GET("/pet", myPetStore.getPet)

	r.Run() // Listen and serve requests on 0.0.0.0:8080.
}

func (p *petStore) listPets(c *gin.Context) {
	p.mu.Lock()
	defer p.mu.Unlock()
	c.JSON(http.StatusOK, p.inventory)
}

func (p *petStore) putPet(c *gin.Context) {
	var pet Animal

	if c.BindJSON(&pet) != nil {
		c.String(http.StatusBadRequest, "invalid json\n")
		return
	}

	if pet.Mythical {
		c.String(http.StatusBadRequest, "Sorry, this pet store does "+
			"not accept mythical beasts.\n")
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// customers must claim mythical creatures!
	if pet.Mythical {
		// Mythical beasts are dangerous.
		p.inventory = make([]Animal, 0)
	}
	p.inventory = append(p.inventory, pet)

	c.String(http.StatusOK, "Thank you for donating your pet\n")
}

func (p *petStore) getPet(c *gin.Context) {
	//Need to protect the store's inventory.
	p.mu.Lock()
	defer p.mu.Unlock()
	//return

	p.customers++
	defer func() { p.customers-- }()

	// Check to see if we have any pets in our inventory.
	if len(p.inventory) == 0 {
		c.String(http.StatusNotFound, "Sorry, we don't have any pets "+
			"available right now.\n")
		return
	}

	// Get the first pet in our inventory.
	petToReturn := p.inventory[0]

	// Walk the pet to the front of the store.
	time.Sleep(time.Second * 3)

	// Update the inventory to reflect that we gave the pet away.
	if len(p.inventory) > 0 {
		p.inventory = p.inventory[1:]
	}

	c.String(http.StatusOK, "We hope that you enjoy your new %s, %s\n",
		petToReturn.Species, petToReturn.Name)
}

type petStore struct {
	inventory []Animal
	customers int
	mu        sync.Mutex
}

type Animal struct {
	Name     string `json:"name" binding:"required"`
	Species  string `json:"species" binding:"required"`
	Mythical bool   `json:"mythical"`
}
