package receipt

import (
	"context"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/entity"
	"github.com/darcops/receipt-proccessor-challenge/util"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

/*
In order to make the code policy-less I defined a series of constants
to make the methods agnostic to the values of the rules. This way, if
the rules change, only the constants need to be changed and not the
methods themselves. Also these constants can be provided by a config
file or environment variables. In this way we can have different versions
of the rules for different environments.
*/
const (
	pointsForAlphanumericCharacter = 1
	pointsForTotalRounded          = 50
	pointsForTotalMultiple         = 25
	pointsForPurchaseTimeInBetween = 10
	pointsForItemPairs             = 5
	pointsForDayOdd                = 6
)

const (
	divisibilityFactorForTotalRounded     = 0.25
	multpliyingFactorForItemsDescriptions = 0.2
)

const (
	startTimeHourForTimeCheck = 14
	endTimeHourForTimeCheck   = 16
)

type receiptService struct{}

// NewReceiptService creates a new receipt service.
func NewReceiptService() *receiptService {
	return &receiptService{}
}

// CreateReceiptID creates an ID for receipt.
func (rs *receiptService) CreateReceiptID(ctx context.Context) string {
	return uuid.New().String()
}

// GetReceiptPoints gets the points of a receipt.
func (rs *receiptService) GetReceiptPoints(ctx context.Context, receipt entity.Receipt) (int64, error) {
	ruleFunctions := []func() (int64, error){
		func() (int64, error) { return rs.getPointsForRetailerName(receipt.Retailer), nil },
		func() (int64, error) { return rs.getPointsForTotalRounded(receipt.Total) },
		func() (int64, error) { return rs.getPointsForTotalMultiple(receipt.Total) },
		func() (int64, error) { return rs.getPointsForItemsCount(receipt.Items), nil },
		func() (int64, error) { return rs.getPointsForItemsDescriptions(receipt.Items), nil },
		func() (int64, error) { return rs.getPointsForPurchaseDate(receipt.PurchaseDate) },
		func() (int64, error) { return rs.getPointsForPurchaseHour(receipt.PurchaseTime) },
	}

	errGroup, _ := errgroup.WithContext(ctx)
	partialPoints := make([]int64, len(ruleFunctions))

	// Execute rule functions concurrently.
	for i, ruleFunc := range ruleFunctions {
		i, ruleFunc := i, ruleFunc
		errGroup.Go(func() error {
			points, err := ruleFunc()
			if err != nil {
				return err
			}
			partialPoints[i] = points
			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return 0, err
	}

	// Calculate total points.
	totalPoints := int64(0)
	for _, points := range partialPoints {
		totalPoints += points
	}

	return totalPoints, nil
}

func (rs *receiptService) getPointsForRetailerName(retailer string) int64 {
	var points int64

	for i := 0; i < len(retailer); i++ {
		if unicode.IsLetter(rune(retailer[i])) || unicode.IsNumber(rune(retailer[i])) {
			points += pointsForAlphanumericCharacter
		}
	}

	return points
}

func (rs *receiptService) getPointsForTotalRounded(total string) (int64, error) {
	value, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0, err
	}

	if value > 0 && value == math.Round(value) {
		return pointsForTotalRounded, nil
	}

	return 0, nil
}

func (rs *receiptService) getPointsForTotalMultiple(total string) (int64, error) {
	value, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0, err
	}

	if value > 0 && math.Mod(value, divisibilityFactorForTotalRounded) == 0 {
		return pointsForTotalMultiple, nil
	}

	return 0, nil
}

func (rs *receiptService) getPointsForItemsCount(items []entity.Item) int64 {
	return int64((len(items) / 2) * pointsForItemPairs)
}

func (rs *receiptService) getPointsForItemsDescriptions(items []entity.Item) int64 {
	var points int64

	for i := 0; i < len(items); i++ {
		description := strings.TrimSpace(items[i].ShortDescription)

		if len(description) > 0 && len(description)%3 == 0 {
			itemPrice, _ := strconv.ParseFloat(items[i].Price, 64)
			pricePoints := int64(math.Ceil(itemPrice * multpliyingFactorForItemsDescriptions))
			points += pricePoints
		}
	}

	return points
}

func (rs *receiptService) getPointsForPurchaseDate(purchaseDate string) (int64, error) {
	if isOdd, err := util.IsDayOdd(purchaseDate); err != nil || !isOdd {
		return 0, err
	}

	return pointsForDayOdd, nil
}

func (rs *receiptService) getPointsForPurchaseHour(purchaseTime string) (int64, error) {
	var isHourBetween bool
	var err error

	if isHourBetween, err = util.IsTimeBetween(
		purchaseTime,
		startTimeHourForTimeCheck,
		endTimeHourForTimeCheck,
	); err != nil || !isHourBetween {
		return 0, err
	}

	return pointsForPurchaseTimeInBetween, nil
}
