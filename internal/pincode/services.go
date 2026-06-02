package pincode

import "context"

func GetPincodesService(ctx context.Context) ([]Pincode, error) {
	return GetPincodesRepository(ctx)
}