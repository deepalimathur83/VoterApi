package rest

import (
	"fmt"
	"strconv"
	"time"

	"drexel.edu/VoterApi/pkg/adding"
	"drexel.edu/VoterApi/pkg/changing"
	"drexel.edu/VoterApi/pkg/listing"
	"github.com/gofiber/fiber/v2"
)

var startTime time.Time

func Handler(add adding.Service, list listing.Service, change changing.Service) {

	startTime = time.Now()

	router := fiber.New()

	//GET /voters - Get all voter resources including all voter history for each voter (note we will discuss the concept of "paging" later, for now you can ignore)
	router.Get("/voters", func(c *fiber.Ctx) error {

		var voters []listing.Voter

		voters, err := list.GetAllVoters()
		if err != nil {
			return err
		}

		c.Status(fiber.StatusOK)
		return c.JSON(voters)
	})

	//GET&POST /voters/:id - Get a single voter resource with voterID=:id including their entire voting history.  POST version adds one to the "database"
	router.Get("/voters/:id", func(c *fiber.Ctx) error {

		var voter listing.Voter

		voterId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return err
		}

		voter, err = list.GetVoter(voterId)
		if err != nil {
			return err
		}

		c.Status(fiber.StatusOK)
		return c.JSON(voter)
	})

	router.Post("/voters", func(c *fiber.Ctx) error {

		var voter adding.Voter

		if err := c.BodyParser(&voter); err != nil {
			return err
		}

		voterId, err := add.RegisterVoter(voter)
		if err != nil {
			return err
		}

		c.Status(fiber.StatusCreated)

		returnMsg := fmt.Sprintf("Voter registration successful. received id: %d", voterId)

		return c.SendString(returnMsg)
	})

	//GET /voters/:id/polls - Gets the JUST the voter history for the voter with VoterID = :id
	router.Get("/voters/:id/polls", func(c *fiber.Ctx) error {
		var voter listing.Voter

		voterId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return err
		}

		voter, err = list.GetVoter(voterId)
		if err != nil {
			return err
		}

		c.Status(fiber.StatusOK)
		return c.JSON(voter.VoterHistory)

	})

	//GET&POST /voters/:id/polls/:pollid - Gets JUST the single voter poll data with PollID = :id and VoterID = :id.  POST version adds one to the "database"

	router.Get("/voters/:id/polls/:pollId", func(c *fiber.Ctx) error {
		var voter listing.VoterHistory

		voterId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return err
		}

		pollId, err := strconv.Atoi(c.Params("pollId"))
		if err != nil {
			return err
		}

		voter, err = list.GetVoterHistory(uint(voterId), uint(pollId))
		if err != nil {
			return err
		}

		c.Status(fiber.StatusOK)
		return c.JSON(voter)

	})

	router.Post("/voters/:id/polls/:pollId", func(c *fiber.Ctx) error {
		var voterHistory adding.VoterHistory

		voterId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return err
		}

		pollId, err := strconv.Atoi(c.Params("pollId"))
		if err != nil {
			return err
		}

		if err := c.BodyParser(&voterHistory); err != nil {
			return err
		}

		voterHistory.VoterId = uint(voterId)
		voterHistory.PollId = uint(pollId)

		err = add.AddVoterHistory(voterHistory)
		if err != nil {
			return err
		}

		c.Status(fiber.StatusCreated)

		return c.SendString("CREATED")

	})

	//GET /voters/health - Returns a "health" record indicating that the voter API is functioning properly and some metadata about the API.  Note the payload can be hard coded, we are mainly looking for a HTTP status code of 200, which means the API is functioning properly.
	router.Get("/health", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusOK)

		return c.SendString("OK!!!")
	})

	router.Listen(":3000")
}
