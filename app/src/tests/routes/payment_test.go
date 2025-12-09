package route_tests

// func SetupPaymentTest(t *testing.T) *fiber.App {
// 	app, db, secretKey := tests.SetupTestApp()
// 	db.Create(&models.Play{
// 		Name: "Test Play",
// 		Author: "Test Author",
// 	})
// 	db.Create(&models.Presentation{
// 		PlayID: 1,
// 		DateTime: time.Now(),
// 		Location: "Test Location",
// 		Price: 100,
// 		SeatLimit: 100,
// 		AvailableSeats: 100,
// 	})
// 	return app
// }

// func TestStripeWebhook(t *testing.T) {
// 	app := SetupPaymentTest(t)
// 	stripeWebhook := services.StripeWebhook{
// 		AmountOfTickets: 1,
// 		PresentationID: 1,
// 		Email: "test@example.com",
// 	}
// }

// func TestSuccess(t *testing.T) {
// 	app := SetupPaymentTest(t)
// 	app.Get("/payment/success", func(c *fiber.Ctx) error {
// 		return c.SendStatus(200)
// 	})
// }

// func TestCancel(t *testing.T) {
// 	app := SetupPaymentTest(t)
// 	app.Get("/payment/cancel", func(c *fiber.Ctx) error {
// 		return c.SendStatus(200)
// 	})
// }