package utils

import validation "github.com/go-ozzo/ozzo-validation/v4"

func ValidateWhenNotNil(rules []validation.Rule) validation.RuleFunc {
	return func(value interface{}) error {
		if value == nil {
			return nil
		}

		return validation.Validate(value, rules...)
	}
}
