package zerologgrpcprovider

// Function receiving key and value from the request and returns new value
// You can use if, for example, for hiding sensitive information from the output
// Example:
//
//	func(key, value string) (string, error) {
//	 if key == "password" {
//		return "<sensitive_data", nil
//	 }
//	 return value, nil
//	}
type RequestValueModifier func(key, value any) (newValue any, err error)
