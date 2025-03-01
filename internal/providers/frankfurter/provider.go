package frankfurter

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

const amount = "1"

func (prv *provider) GetRate(ctx context.Context, toIso, fromIso string) ([]byte, error) {
	rqCtx, cancel := context.WithTimeout(ctx, prvTimeout)
	defer cancel()

	params := make(url.Values, 3)
	params.Add("amount", amount)
	params.Add("from", fromIso)
	params.Add("to", toIso)

	resp, err := prv.getRate.
		SetQueryParameters(params).
		DoWithoutBody(rqCtx)
	if err != nil {
		prv.logger.Error().Msg(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		_err := fmt.Errorf("unexpected status code from provider: %s", resp.Status)
		prv.logger.Error().Msg(_err.Error())
		return nil, _err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		prv.logger.Error().Msg(err.Error())
		return nil, err
	}

	return respBody, err
}

func (prv *provider) GetCurrencyList(ctx context.Context) ([]byte, error) {
	rqCtx, cancel := context.WithTimeout(ctx, prvTimeout)
	defer cancel()

	resp, err := prv.getCurrencyList.
		DoWithoutBody(rqCtx)
	if err != nil {
		prv.logger.Error().Msg(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		_err := fmt.Errorf("unexpected status code from provider: %s", resp.Status)
		prv.logger.Error().Msg(_err.Error())
		return nil, _err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		prv.logger.Error().Msg(err.Error())
		return nil, err
	}

	return respBody, err
}
