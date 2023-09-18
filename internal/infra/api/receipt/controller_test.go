package receipt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darcops/receipt-proccessor-challenge/internal/pkg/entity"
	"github.com/darcops/receipt-proccessor-challenge/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func TestCreateReceipt(t *testing.T) {
	mockService := &mocks.ReceiptService{}

	testCases := []struct {
		name string

		service             *mocks.ReceiptService
		wantServiceResponse string

		request entity.Receipt

		wantStatusCode int
	}{
		{
			name: "should create a receipt",

			service:             mockService,
			wantServiceResponse: "1234567890",

			request: entity.Receipt{
				Retailer:     "$Walmart/   ",
				PurchaseDate: "2020-01-01",
				PurchaseTime: "15:00",
				Items: []entity.Item{
					{
						ShortDescription: "Item 1",
						Price:            "1.00",
					},
					{
						ShortDescription: "Item 1",
						Price:            "1.00",
					},
				},
				Total: "2.00",
			},

			wantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		// Set the memory storage for receipts and points.
		receiptByID = map[string]entity.Receipt{}
		receiptPointsByID = make(map[string]int64)

		// Create a new router for tests.
		router := gin.Default()
		gin.SetMode(gin.TestMode)

		controller := newReceiptController(tc.service)

		// Mock the desired response from the service.
		tc.service.On(
			"CreateReceiptID",
			mock.Anything, /* context.Context */
		).Return(tc.wantServiceResponse).Once()

		router.POST("/process", controller.createReceipt)

		server := httptest.NewServer(router)

		t.Run(tc.name, func(t *testing.T) {
			requestBody, err := json.Marshal(&tc.request)
			if err != nil {
				t.Errorf("CreateReceipt() = Marshaling error %v", err)
			}

			response, err := http.Post(
				fmt.Sprintf("%s/process", server.URL),
				"application/json",
				bytes.NewBuffer(requestBody),
			)
			if err != nil {
				t.Errorf("CreateReceipt() = error %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tc.wantStatusCode {
				t.Errorf("CreateReceipt() = %v, want %v", response.StatusCode, tc.wantStatusCode)
			}
		})
	}
}

func TestGetReceiptPoints(t *testing.T) {
	mockService := &mocks.ReceiptService{}

	testCases := []struct {
		name string

		service             *mocks.ReceiptService
		wantServiceResponse int64

		receiptID string

		wantStatusCode int
	}{
		{
			name: "should return points for receipt",

			service:             mockService,
			wantServiceResponse: 10,

			wantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		mockReceiptID := "1234567890"
		// Set the memory storage for receipts and points.
		receiptByID = map[string]entity.Receipt{
			mockReceiptID: {},
		}
		receiptPointsByID = make(map[string]int64)

		// Create a new router for tests.
		router := gin.Default()
		gin.SetMode(gin.TestMode)

		controller := newReceiptController(tc.service)

		// Mock the desired response from the service.
		tc.service.On(
			"GetReceiptPoints",
			mock.Anything, /* context.Context */
			mock.Anything, /* entity.Receipt */
		).Return(tc.wantServiceResponse, nil).Once()

		router.GET("/:receipt_id/points", controller.getReceiptPoints)

		server := httptest.NewServer(router)

		t.Run(tc.name, func(t *testing.T) {
			response, err := http.Get(
				fmt.Sprintf("%s/%s/points", server.URL, mockReceiptID),
			)
			if err != nil {
				t.Errorf("GetReceiptPoints() = error %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tc.wantStatusCode {
				t.Errorf("GetReceiptPoints() = %v, want %v", response.StatusCode, tc.wantStatusCode)
			}

			got := map[string]int64{}
			err = json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Errorf("GetReceiptPoints() = Unmarshaling response error %v", err)
			}

			if got["points"] != tc.wantServiceResponse {
				t.Errorf("GetReceiptPoints() = %v, want %v", got["points"], tc.wantServiceResponse)
			}
		})
	}
}
