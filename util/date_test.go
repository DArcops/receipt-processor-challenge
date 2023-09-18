package util

import "testing"

func TestIsTimeBetween(t *testing.T) {
	testCases := []struct {
		name string

		timeToCheck   string
		startTimeHour int
		endTimeHour   int

		want    bool
		wantErr bool
	}{
		{
			name: "should return true for time between 2:00pm and 4:00pm",

			timeToCheck:   "15:13",
			startTimeHour: 14,
			endTimeHour:   16,

			want: true,
		},
		{
			name: "should return false for time after 4:00pm",

			timeToCheck:   "17:13",
			startTimeHour: 14,
			endTimeHour:   16,

			want: false,
		},
		{
			name: "should fail due invalid hour format",

			timeToCheck:   "20-01",
			startTimeHour: 14,
			endTimeHour:   16,

			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := IsTimeBetween(tc.timeToCheck, tc.startTimeHour, tc.endTimeHour)

			if got != tc.want {
				t.Errorf("IsTimeBetween() got = %v, want %v", got, tc.want)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("IsTimeBetween() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func TestIsDayOdd(t *testing.T) {
	testCases := []struct {
		name string

		date string

		want    bool
		wantErr bool
	}{
		{
			name: "should return true for odd day",

			date: "2021-01-01",

			want: true,
		},
		{
			name: "should return false for even day",

			date: "2021-01-02",

			want: false,
		},
		{
			name: "should fail due invalid date format",

			date: "01-01-2021",

			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := IsDayOdd(tc.date)

			if got != tc.want {
				t.Errorf("IsDayOdd() got = %v, want %v", got, tc.want)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("IsDayOdd() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
