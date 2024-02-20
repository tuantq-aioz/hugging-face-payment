package utils

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func Log(start time.Time, action string, err error) {
	responseTime := time.Since(start)
	if err != nil {
		log.Error().
			Err(err).
			Str("response_time", responseTime.String()).
			Str("err", fmt.Sprintf("%v", err)).
			Msg(action)
		return
	}

	log.Info().Str("response_time", responseTime.String()).Msg(action)
}
