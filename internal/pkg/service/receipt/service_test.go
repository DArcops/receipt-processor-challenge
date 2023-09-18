package receipt

import (
	"context"
	"testing"

	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/entity"
)

// func TestCreateReceiptID(t *testing.T) {
// 	testCases := []struct {
// 		name    string
// 		ctx     context.Context
// 		service *receiptService

// 		want string
// 	}{
// 		{
// 			name:    "CreateReceiptID",
// 			ctx:     context.Background(),
// 			service: NewReceiptService(),

// 			want: "1234567890",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			got := tc.service.CreateReceiptID(tc.ctx)

// 			if got != tc.want {
// 				t.Errorf("CreateReceiptID() = %v, want %v", got, tc.want)
// 			}
// 		})
// 	}
// }

func TestGetReceiptPoints(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     context.Context
		service *receiptService

		receipt entity.Receipt

		want    int64
		wantErr error
	}{
		{
			name:    "should return points for receipt",
			ctx:     context.Background(),
			service: NewReceiptService(),

			receipt: entity.Receipt{
				Retailer:     "$Walmart/   ",
				PurchaseDate: "2020-01-01",
				PurchaseTime: "15:00",
				Items: []entity.Item{
					{
						ShortDescription: "Item 1",
						Price:            "10.00",
					},
					{
						ShortDescription: "Item 2",
						Price:            "10.87",
					},
					{
						ShortDescription: "Item 3",
						Price:            "1.00",
					},
				},
				Total: "100.00",
			},

			want: 109,
		},
		{
			name:    "should return points only for retailer name",
			ctx:     context.Background(),
			service: NewReceiptService(),

			receipt: entity.Receipt{
				Retailer:     " retailer name example-2 ",
				PurchaseDate: "2020-01-02",
				PurchaseTime: "12:00",
				Total:        "100.01",
			},

			want: 20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.service.GetReceiptPoints(tc.ctx, tc.receipt)

			if got != tc.want {
				t.Errorf("GetReceiptPoints() = %v, want %v", got, tc.want)
			}

			if err != tc.wantErr {
				t.Errorf("GetReceiptPoints() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestGetPointsForRetailerName(t *testing.T) {
	testCases := []struct {
		name    string
		service *receiptService

		retailerName string

		want int64
	}{
		{
			name:    "should return zero points for retailer name",
			service: NewReceiptService(),

			retailerName: " ",

			want: 0,
		},
		{
			name:    "should return 0 for retailer name without alphanumeric characters",
			service: NewReceiptService(),

			retailerName: "%#&$", // retailer name without alphanumeric characters.

			want: 0,
		},
		{
			name:    "should return 1 for retailer name with one alphanumeric character",
			service: NewReceiptService(),

			retailerName: "a%#&$",

			want: 1,
		},
		{
			name:    "should ommit spaces",
			service: NewReceiptService(),

			retailerName: "  retailer name with spaces 1  ",

			want: 23,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.service.getPointsForRetailerName(tc.retailerName)

			if got != tc.want {
				t.Errorf("getPointsForRetailerName() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetPointsForTotalRounded(t *testing.T) {
	testCases := []struct {
		name    string
		service *receiptService

		total string

		want    int64
		wantErr bool
	}{
		{
			name:    "should return zero points for total",
			service: NewReceiptService(),

			total: "0.00",

			want: 0,
		},
		{
			name:    "should return zero points for total with cents",
			service: NewReceiptService(),

			total: "0.01",

			want: 0,
		},
		{
			name:    "should return 50 points for total with no cents",
			service: NewReceiptService(),

			total: "1.00",

			want: pointsForTotalRounded,
		},
		{
			name:    "should fail due invalid total",
			service: NewReceiptService(),

			total: "invalid total",

			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.service.getPointsForTotalRounded(tc.total)

			if got != tc.want {
				t.Errorf("getPointsForTotalRounded() = %v, want %v", got, tc.want)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("getPointsForTotalRounded() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestGetPointsForTotalMultiple(t *testing.T) {
	testCases := []struct {
		name    string
		service *receiptService

		total string

		want    int64
		wantErr bool
	}{
		{
			name:    "should return zero points for total",
			service: NewReceiptService(),

			total: "0.00",

			want: 0,
		},
		{
			name:    "should return zero points for total with cents",
			service: NewReceiptService(),

			total: "0.01",

			want: 0,
		},
		{
			name:    "should return 25 points for total if is multiple of 0.25",
			service: NewReceiptService(),

			total: "1.00",

			want: pointsForTotalMultiple,
		},
		{
			name:    "should fail due invalid total",
			service: NewReceiptService(),

			total: "invalid total",

			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.service.getPointsForTotalMultiple(tc.total)

			if got != tc.want {
				t.Errorf("getPointsForTotalMultiple() = %v, want %v", got, tc.want)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("getPointsForTotalMultiple() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestGetPointsForItemsCount(t *testing.T) {
	testCases := []struct {
		name    string
		service *receiptService

		items []entity.Item

		want int64
	}{
		{
			name:    "should return zero points for items count",
			service: NewReceiptService(),

			items: []entity.Item{},

			want: 0,
		},
		{
			name:    "should return 5 points for items count",
			service: NewReceiptService(),

			items: []entity.Item{
				{
					ShortDescription: "Item 1",
					Price:            "10.00",
				},
				{
					ShortDescription: "Item 2",
					Price:            "10.00",
				},
			},

			want: pointsForItemPairs,
		},
		{
			name:    "should return 5 points for 3 items count",
			service: NewReceiptService(),

			items: []entity.Item{
				{
					ShortDescription: "Item 1",
					Price:            "10.00",
				},
				{
					ShortDescription: "Item 2",
					Price:            "10.00",
				},
				{
					ShortDescription: "Item 3",
					Price:            "10.00",
				},
			},

			want: pointsForItemPairs,
		},
		{
			name:    "should return 15 points for 6 items count",
			service: NewReceiptService(),

			items: []entity.Item{
				{
					ShortDescription: "Item 1",
					Price:            "10.00",
				},
				{
					ShortDescription: "Item 2",
					Price:            "10.00",
				},
				{
					ShortDescription: "Item 3",
					Price:            "10.00",
				},
				{
					ShortDescription: "Item 4",
					Price:            "10.00",
				},
				{
					ShortDescription: "Item 5",
					Price:            "10.00",
				},
				{
					ShortDescription: "Item 6",
					Price:            "10.00",
				},
			},

			want: pointsForItemPairs * 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.service.getPointsForItemsCount(tc.items)

			if got != tc.want {
				t.Errorf("getPointsForItemsCount() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetPointsForPurchaseDate(t *testing.T) {
	testCases := []struct {
		name    string
		service *receiptService

		purchaseDate string

		want    int64
		wantErr bool
	}{
		{
			name:    "should return 6 points for odd day",
			service: NewReceiptService(),

			purchaseDate: "2020-01-01",

			want: pointsForDayOdd,
		},
		{
			name:    "should return zero points for even day",
			service: NewReceiptService(),

			purchaseDate: "2020-01-02",

			want: 0,
		},
		{
			name:    "should fail due invalid date",
			service: NewReceiptService(),

			purchaseDate: "invalid date",

			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.service.getPointsForPurchaseDate(tc.purchaseDate)

			if got != tc.want {
				t.Errorf("getPointsForPurchaseDate() = %v, want %v", got, tc.want)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("getPointsForTotalMultiple() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestGetPointsForPurchaseHour(t *testing.T) {
	testCases := []struct {
		name    string
		service *receiptService

		purchaseHour string

		want    int64
		wantErr bool
	}{
		{
			name:    "should return points for purchase hour at 2:00pm",
			service: NewReceiptService(),

			purchaseHour: "14:00",

			want: pointsForPurchaseTimeInBetween,
		},
		{
			name:    "should return points for purchase hour at 4:00pm",
			service: NewReceiptService(),

			purchaseHour: "16:00",

			want: pointsForPurchaseTimeInBetween,
		},
		{
			name:    "should return 10 points for purchase hour between 2:00pm and 4:00pm",
			service: NewReceiptService(),

			purchaseHour: "15:13",

			want: pointsForPurchaseTimeInBetween,
		},
		{
			name:    "should return zero points for purchase hour before 2:00pm",
			service: NewReceiptService(),

			purchaseHour: "13:59",

			want: 0,
		},
		{
			name:    "should return zero points for purchase hour after 4:00pm",
			service: NewReceiptService(),

			purchaseHour: "16:01",

			want: 0,
		},
		{
			name:    "should fail due invalid hour",
			service: NewReceiptService(),

			purchaseHour: "invalid hour",

			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.service.getPointsForPurchaseHour(tc.purchaseHour)

			if got != tc.want {
				t.Errorf("getPointsForPurchaseHour() = %v, want %v", got, tc.want)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("getPointsForPurchaseHour() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}
