package services

// CarRequestBody holds the necessary fields for a car
type CarRequestBody struct {
	Make  string
	Model string
	Year  int
	Color string
}
