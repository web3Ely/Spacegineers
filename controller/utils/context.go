package controllerhelper

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc/metadata"
)

func getComponentIDFromContext(ctx context.Context, key string) (int32, error) {
	ctxData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return -1, errors.New("room controller could not locate the Item metadata")
	}

	ctxEntry := ctxData.Get(key)
	if len(ctxEntry) == 0 {
		return -1, errors.New("room controller could not get the Item number")
	}

	convert, _ := strconv.Atoi(ctxEntry[0])
	return int32(convert), nil
}
