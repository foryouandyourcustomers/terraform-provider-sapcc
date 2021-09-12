package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func (v ValueContainsInValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	value, ok := req.AttributeConfig.(types.String) // see also attr.ValueAs() proposal

	if !ok {
		resp.Diagnostics.AddError(
			"Invalid value type",
			fmt.Sprintf("received incorrect value type (%T) at path: %s", req.AttributeConfig, req.AttributePath),
		)

		return
	}

	if value.Unknown {
		resp.Diagnostics.AddError(
			"Unknown validation value",
			fmt.Sprintf("received unknown value at path: %s", req.AttributePath),
		)

		return
	}

	for _, val := range v.values {
		if val == value.Value {
			return
		}
	}

	resp.Diagnostics.AddError(
		fmt.Sprintf("Invalid value for '%s'", req.AttributePath.Steps()[0].(tftypes.AttributeName)),
		fmt.Sprintf("Expecting one of '%s'", strings.Join(v.values, ", ")),
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

func (v ValueRegexMatchValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	value, ok := req.AttributeConfig.(types.String) // see also attr.ValueAs() proposal

	if !ok {
		resp.Diagnostics.AddError(
			"Invalid value type",
			fmt.Sprintf("received incorrect value type (%T) at path: %s", req.AttributeConfig, req.AttributePath),
		)

		return
	}

	if value.Unknown {
		resp.Diagnostics.AddError(
			"Unknown validation value",
			fmt.Sprintf("received unknown value at path: %s", req.AttributePath),
		)

		return
	}

	if !v.regex.MatchString(value.Value) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid value for '%s'", req.AttributePath.Steps()[0].(tftypes.AttributeName)),
			fmt.Sprintf("Value should match regex `%s`", v.regex.String()),
		)
	}
}

func ValueMustMatchRegex(regex string) ValueRegexMatchValidator {
	return ValueRegexMatchValidator{
		regex: regexp.MustCompile(regex),
	}
}
