package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ValueCannotBeEmptyValidator struct {
	tfsdk.AttributeValidator
}

type ValueContainsInValidator struct {
	tfsdk.AttributeValidator
	values []string
}

type ValueRegexMatchValidator struct {
	tfsdk.AttributeValidator
	regex *regexp.Regexp
}

func (v ValueContainsInValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Value should be one of '%s'", strings.Join(v.values, ", "))
}

func (v ValueContainsInValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Value should be one of `'%s'`", strings.Join(v.values, ", "))
}

func (v ValueContainsInValidator) Validate(_ context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	value, ok := req.AttributeConfig.(types.String) // see also attr.ValueAs() proposal

	if !ok {

		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid value type",
			fmt.Sprintf("received incorrect value type (%T) at path: %s", req.AttributeConfig, req.AttributePath),
		)

		return
	}

	if value.Unknown || value.Null {
		return
	}

	for _, val := range v.values {
		if val == value.Value {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		fmt.Sprintf("Invalid value '%s'", value.Value),
		fmt.Sprintf("received incorrect value type (%T) at path: %s", req.AttributeConfig, req.AttributePath),
	)
}

func ValueMustBeOneOf(values ...string) ValueContainsInValidator {
	return ValueContainsInValidator{
		values: values,
	}
}

func (v ValueRegexMatchValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Value should match regex '%s'", v.regex.String())
}

func (v ValueRegexMatchValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Value should match regex `%s`", v.regex.String())
}

func (v ValueRegexMatchValidator) Validate(_ context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	value, ok := req.AttributeConfig.(types.String) // see also attr.ValueAs() proposal

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid value type",
			fmt.Sprintf("received incorrect value type (%T) at path: %s", req.AttributeConfig, req.AttributePath),
		)

		return
	}

	if value.Unknown || value.Null {
		return
	}

	if !v.regex.MatchString(value.Value) {

		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			fmt.Sprintf("Invalid value '%s'", value.Value),
			fmt.Sprintf("Value should match regex `%s`", v.regex.String()),
		)
	}
}

func ValueMustMatchRegex(regex string) ValueRegexMatchValidator {
	return ValueRegexMatchValidator{
		regex: regexp.MustCompile(regex),
	}
}

func (v ValueCannotBeEmptyValidator) Description(_ context.Context) string {
	return "Value Can not be empty"
}

func (v ValueCannotBeEmptyValidator) MarkdownDescription(_ context.Context) string {
	return "Value Can not be empty"
}

func (v ValueCannotBeEmptyValidator) Validate(_ context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	value, ok := req.AttributeConfig.(types.String) // see also attr.ValueAs() proposal

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid value type",
			fmt.Sprintf("received incorrect value type (%T) at path: %s", req.AttributeConfig, req.AttributePath),
		)

		return
	}

	if value.Unknown || value.Null {
		return
	}

	if strings.TrimSpace(value.Value) == "" {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Empty Value received",
			"",
		)

		return
	}
}

func ValueMustNotBeEmpty() ValueCannotBeEmptyValidator {
	return ValueCannotBeEmptyValidator{}
}
