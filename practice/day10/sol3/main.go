package main

/*
q3. Create a new custom type based on float64 for handling temperatures in Celsius.
    Implement the Following Methods (not functions):
    Method 1: ToFahrenheit
    Description: Converts the Celsius temperature to Fahrenheit.
    Signature: ToFahrenheit() float64
    Method 2: IsFreezing
    Description: Checks if the temperature is at or below the freezing point (0°C).
    Signature: IsFreezing() bool
*/
type celsius struct {
	temp float64
}

func (c *celsius) ToFahrenheit() float64 {
	return c.temp*(9/5) + 32
}

func (c *celsius) IsFreezing() bool {
	return c.temp <= 0
}

func main() {

}
