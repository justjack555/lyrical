package parser

import "os"

type ArgsError struct {
	err string
}

func (ae *ArgsError) Error() string {
	return ae.err
}

func Start() error {
	if numArgs := len(os.Args); numArgs < 2 {
		return &ArgsError{
			err: "Usage: ./parser <config_file>",
		}
	}


	return nil
}