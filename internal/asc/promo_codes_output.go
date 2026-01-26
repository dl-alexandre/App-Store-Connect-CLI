package asc

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func promoCodeExpiresDate(attrs PromoCodeAttributes) string {
	if strings.TrimSpace(attrs.ExpiresDate) != "" {
		return strings.TrimSpace(attrs.ExpiresDate)
	}
	return strings.TrimSpace(attrs.ExpirationDate)
}

func promoCodeProductType(attrs PromoCodeAttributes) string {
	if strings.TrimSpace(string(attrs.ProductType)) != "" {
		return string(attrs.ProductType)
	}
	return string(attrs.Type)
}

func printPromoCodesTable(resp *PromoCodesResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Code\tExpires\tUsed\tExpired\tType")
	for _, item := range resp.Data {
		attrs := item.Attributes
		fmt.Fprintf(w, "%s\t%s\t%t\t%t\t%s\n",
			sanitizeTerminal(attrs.Code),
			sanitizeTerminal(promoCodeExpiresDate(attrs)),
			attrs.IsUsed,
			attrs.IsExpired,
			sanitizeTerminal(promoCodeProductType(attrs)),
		)
	}
	return w.Flush()
}

func printPromoCodesMarkdown(resp *PromoCodesResponse) error {
	fmt.Fprintln(os.Stdout, "| Code | Expires | Used | Expired | Type |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		attrs := item.Attributes
		fmt.Fprintf(os.Stdout, "| %s | %s | %t | %t | %s |\n",
			escapeMarkdown(attrs.Code),
			escapeMarkdown(promoCodeExpiresDate(attrs)),
			attrs.IsUsed,
			attrs.IsExpired,
			escapeMarkdown(promoCodeProductType(attrs)),
		)
	}
	return nil
}
