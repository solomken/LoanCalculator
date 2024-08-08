package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

const (
	paymentTypeDiff    = "diff"
	paymentTypeAnnuity = "annuity"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Incorrect parameters")

		return
	}

	paymentType := flag.String("type", "", "Payment type annuity or differentiated is required")
	payment := flag.Float64("payment", 0.0, "Payment amount is expected")
	principal := flag.Float64("principal", 0, "Principal is expected")
	periods := flag.Float64("periods", 0, "Expected number of months needed to repay the loan")
	interestRate := flag.Float64("interest", 0.0, "Interest rate is required")
	flag.Parse()

	if *interestRate <= 0 || *principal < 0 || *payment < 0 || *periods < 0 {
		fmt.Println("Incorrect parameters")

		return
	}

	if *paymentType != paymentTypeAnnuity && *paymentType != paymentTypeDiff {
		fmt.Println("Incorrect parameters")

		return
	}

	if *paymentType == paymentTypeDiff && (*payment > 0 || *periods <= 0) {
		fmt.Println("Incorrect parameters")

		return
	}

	switch *paymentType {
	case paymentTypeDiff:
		printDiffPayments(*interestRate, *periods, *principal)
		return
	case paymentTypeAnnuity:
		if *payment > 0 && *principal > 0 && *periods == 0 {
			numberOfPayments := calculateNumberOfPayments(*payment, *principal, *interestRate)
			printNumberOfMonths(numberOfPayments)
			totalPayment := *payment * float64(numberOfPayments)
			overpayment := math.Ceil(totalPayment - *principal)
			fmt.Printf("Overpayment = %.0f\n", overpayment)

			return
		}

		if *principal > 0 && *periods > 0 && *payment == 0 {
			monthlyPayment := calculateMonthlyPayment(*principal, *interestRate, *periods)
			totalPayment := float64(monthlyPayment) * *periods
			overpayment := math.Ceil(totalPayment - *principal)
			fmt.Printf("Your annuity payment = %d!\n", monthlyPayment)
			fmt.Printf("Overpayment = %.0f\n", overpayment)

			return
		}

		if *payment > 0 && *periods > 0 && *principal == 0 {
			loanPrincipal := calculateLoanPrincipal(*periods, *interestRate, *payment)
			totalPayment := *payment * *periods
			overpayment := math.Ceil(totalPayment - float64(loanPrincipal))
			fmt.Printf("Your loan principal = %d!\n", loanPrincipal)
			fmt.Printf("Overpayment = %.0f\n", overpayment)

			return
		}
	}
}

func calculateNumberOfPayments(payment, principal, interestRate float64) int {
	i := interestRate / (12 * 100)
	n := math.Log(payment/(payment-i*principal)) / math.Log(1+i)

	return int(math.Ceil(n))
}

func printNumberOfMonths(numberOfPayments int) {
	if numberOfPayments == 1 {
		fmt.Printf("It will take %d month to repay this loan!\n", numberOfPayments)
	} else if numberOfPayments > 1 && numberOfPayments < 12 {
		fmt.Printf("It will take %d months to repay this loan!\n", numberOfPayments)
	} else {
		years := numberOfPayments / 12
		months := numberOfPayments % 12

		if years == 1 && months == 0 {
			fmt.Printf("It will take 1 year to repay this loan!\n")
		} else if years == 1 {
			fmt.Printf("It will take 1 year and %d months to repay this loan!\n", months)
		} else if months == 0 {
			fmt.Printf("It will take %d years to repay this loan!\n", years)
		} else {
			fmt.Printf("It will take %d years and %d months to repay this loan!\n", years, months)
		}
	}
}

func calculateMonthlyPayment(principal, interestRate, periods float64) int {
	monthlyInterestRate := interestRate / 12 / 100
	monthlyPayment := principal * (monthlyInterestRate * math.Pow(1+monthlyInterestRate, periods)) /
		(math.Pow(1+monthlyInterestRate, periods) - 1)

	return int(math.Ceil(monthlyPayment))
}

func calculateLoanPrincipal(periods, interestRate, payment float64) int {
	monthlyInterestRate := interestRate / 12 / 100
	loanPrincipal := payment / ((monthlyInterestRate * math.Pow(1+monthlyInterestRate, periods)) /
		(math.Pow(1+monthlyInterestRate, periods) - 1))

	return int(loanPrincipal)
}

func printDiffPayments(interestRate, periods, principal float64) {
	var overpayment float64
	monthlyInterestRate := interestRate / 12 / 100

	for currentRepaymentMonth := 1; currentRepaymentMonth <= int(periods); currentRepaymentMonth++ {
		differentiatedPayment := math.Ceil((principal / periods) +
			monthlyInterestRate*(principal-(principal*float64(currentRepaymentMonth-1)/periods)))
		overpayment += differentiatedPayment
		fmt.Printf("Month %d: payment is %.0f\n", currentRepaymentMonth, differentiatedPayment)
	}
	overpayment = math.Ceil(overpayment - principal)
	fmt.Printf("\nOverpayment = %.0f\n", overpayment)
}
