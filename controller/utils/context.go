package controllerhelper

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

func getComponentIDFromContext(ctx context.Context, key string) (string, error) {
	ctxData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("room controller could not locate the Item metadata")
	}

	ctxEntry := ctxData.Get(key)
	if len(ctxEntry) == 0 {
		return "", errors.New("room controller could not get the Item number")
	}

	return ctxEntry[0], nil
}
