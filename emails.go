package main

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

func SendSaleEmail(sale *SaleWithProducts, s *PostgresStore, customer *Customer, users []*User) error {
	// Recover the name and image of original products
	var products []*Product
	for _, pro := range sale.Products {
		product, err := s.GetProductByID(pro.ProductID)
		if err != nil {
			return err
		}
		products = append(products, product)
	}

	sdKey := os.Getenv("SENDGRID_API_KEY")
	sdSender := os.Getenv("SENDGRID_CUSTOM_SENDER")

	from := mail.NewEmail("Dashboard API", sdSender)
	subject := fmt.Sprintf("🛍️ Nueva venta para %v", customer.Name)
	plainTextContent := "Se ha generado una nueva venta"

	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>🛍️ Nueva venta para %s</title>
	</head>
	<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #111827;">
	
	<table cellpadding="20" cellspacing="0" width="100%%" style="border-collapse: collapse; max-width: 600px; margin: 0 auto; background-color: #111827; color: #fafafa;">
	
	<tr>
		<td colspan="2" style="padding-top: 20px;">
			<h2 style="margin: 0;">🦋 Se ha generado una nueva venta para %s 🍭</h2>
		</td>
		</tr>
		<tr>
		<td colspan="2" style="padding-top: 10px;">
			<hr style="border: 1px solid #444444;">
		</td>
	</tr>
	`,
		customer.Name,
		customer.Name,
	)

	var totalEarningsHTML int
	for _, pv := range sale.Products {
		totalEarningsHTML += pv.Price
		bgColor, textColor := ColorFromLocalConstants(pv.Color)
		htmlContent += fmt.Sprintf(`
		  <tr>
			<td colspan="2" style="padding-top: 20px; background-color: #1F2937; border-radius: 8px;">
			<table cellpadding="0" cellspacing="0" width="100%%">
			<tr>
			<td style="width: 10%%;">
			<img src="%s" alt="%s" width="42" height="42" style="border-radius: 50%%; display: block; margin-right: 6px;">
			</td>
			<td style="width: 70%%;">
			<p style="margin: 0 0 10px;">%s</p>
			<span style="background-color: %s; color: %s; padding: 5px 10px; border-radius: 5px; font-size: 14px;">%s</span>
			</td>
			<td style="width: 20%%;">
			<p>%s</p>
			</td>
			</tr>
			</table>
			</td>
		</tr>
		<tr style="height: 6px;">
			<td colspan="2"></td>
		</tr>`,
			findProductVariation(products, pv.ProductID).Image,
			findProductVariation(products, pv.ProductID).Name,
			findProductVariation(products, pv.ProductID).Name,
			bgColor,
			textColor,
			pv.Color,
			formatCurrency(pv.Price, "COP"),
		)
	}

	htmlContent += fmt.Sprintf(`
			<tr>
				<td style="width: 60%%;">
				</td>
				<td style="width: 40%%; text-align: right; color: #49DE80">
					<p style="margin: 0;">Ganancias: %s</p>
				</td>
				</tr>
				<tr>
					<td>
						<p style="margin: 0 0 4px;">%s</p>
						<p style="margin: 0 0 4px;">%s</p>
						<p style="margin: 0 0 4px">%s / %s</p>
						<a href="https://www.instagram.com/%s" style="color: #4f46e5; text-decoration: none; margin: 0 0 4px; display: block">@%s</a>
						<p style="margin: 0 0 4px;"><a href="tel:%v" style="color: #4f46e5; text-decoration: none;">%v</a></p>
						<p style="margin: 0 0 4px;">%s</p>
						<p style="margin: 0 0 4px;">cc %s</p>
					</td>
				</tr>
			</table>
		</body>
	</html>	
	`,
		formatCurrency(totalEarningsHTML, "COP"),
		customer.Name,
		customer.Address,
		customer.City,
		customer.Department,
		customer.InstagramAccount,
		customer.InstagramAccount,
		customer.Phone,
		customer.Phone,
		customer.Comments,
		customer.Cc,
	)

	for _, user := range users {
		to := mail.NewEmail(user.FirstName, user.Email)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(sdKey)
		_, err := client.Send(message)
		if err != nil {
			log.Printf("Error sending email to %s (%s): %s\n", user.FirstName, user.Email, err)
		}
	}

	return nil
}
