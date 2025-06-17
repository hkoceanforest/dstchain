package util

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	"google.golang.org/grpc/status"
	"regexp"
	"strings"
)

func QueryWithDataWithUnwrapErr(clientCtx client.Context, rpcUrc string, bz []byte) ([]byte, int64, error) {
	res, n, err := clientCtx.QueryWithData(rpcUrc, bz)

	err = GrpcErrFilter(err)

	return res, n, err
}

func GrpcErrFilter(err error) error {
	if grpcStatus, ok := status.FromError(err); ok && err != nil {
		err = errors.New(grpcStatus.Message())
	}

	return err
}

func ColonErrFilter(err error) error {
	sp := strings.Split(err.Error(), ": ")
	err = errors.New(sp[len(sp)-1])

	return err
}

func ErrFilter(err error) error {
	err = ColonErrFilter(GrpcErrFilter(err))

	return err
}


func ErrEegularFilter(err error) error {
	
	
	errContent := err.Error()
	regex1 := regexp.MustCompile(`codespace sdk code [0-9]+: `)

	errContent = regex1.ReplaceAllString(errContent, "")

	regex2 := regexp.MustCompile(`: failed to execute message; message index: [0-9]+`)

	errContent = regex2.ReplaceAllString(errContent, "")

	return errors.New(errContent)
}
