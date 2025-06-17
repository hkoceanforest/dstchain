package client

import (
	"bytes"
	"errors"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/util"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	ttypes "github.com/tendermint/tendermint/types"
	"net/url"
	"strconv"
	"strings"
	"time"
	
	"io/ioutil"
	"net/http"
)

type StringEvent struct {
	Type       string      `json:"type,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type ABCIMessageLog struct {
	MsgIndex uint16 `json:"msg_index"`
	Log      string `json:"log"`

	Events []StringEvent `json:"events"`
}

type TxDetail struct {
	Height string `json:"height"`
	Status string `json:"status"`
	Txhash string `json:"txhash"`
	Error  string `json:"error"`
}

type ClientLatestHeight struct {
	ClientState struct {
		LatestHeight struct {
			RevisionNumber string `json:"revision_number"`
			RevisionHeight string `json:"revision_height"`
		} `json:"latest_height"`
	} `json:"client_state"`
}

type BlockResponse struct {
	Block struct {
		Header struct {
			Height string `json:"height"`
		} `json:"header"`
	} `json:"block"`
}

type AcknowledgementResponse struct {
	Acknowledgement string `json:"acknowledgement"`
}


func GetRequest(server, params string) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, 
			MaxIdleConnsPerHost: 512,  
		},
		Timeout: time.Second * 60, 
	}
	bodyReader := strings.NewReader("")
	req, err := http.NewRequest("GET", server+params, bodyReader)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	
	return string(body), err
}


func PostValuesRequest(server, url string, values url.Values) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, 
			MaxIdleConnsPerHost: 512,  
		},
		Timeout: time.Second * 60, 
	}
	req, err := http.NewRequest("POST", server+url, nil)
	req.PostForm = values
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	if resp.StatusCode != 200 {
		if len(body) == 0 {
			return "", errors.New("" + strconv.Itoa(resp.StatusCode))
		}
		return "", errors.New(string(body))
	}
	
	return string(body), err
}


func PostRequest(server, url string, params []byte) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, 
			MaxIdleConnsPerHost: 512,  
		},
		Timeout: time.Second * 60, 
	}
	req, err := http.NewRequest("POST", server+url, bytes.NewReader(params))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		if len(body) == 0 {
			return "", errors.New("" + strconv.Itoa(resp.StatusCode))
		}
		return "", errors.New(string(body))
	}

	return string(body), err
}

func PostRequestByTimeout(server, url string, params []byte, timeout time.Duration) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, 
			MaxIdleConnsPerHost: 512,  
		},
		Timeout: timeout, 
	}
	req, err := http.NewRequest("POST", server+url, bytes.NewReader(params))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	
	return string(body), err
}



func unmarshalMsg(msg sdk.Msg, obj interface{}) error {
	msgByte, err := util.Json.Marshal(msg)
	if err != nil {
		return err
	}
	return util.Json.Unmarshal(msgByte, &obj)
}

func convertTxToStdTx(cosmosTx sdk.Tx) (*legacytx.StdTx, error) {
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	signingTx, ok := cosmosTx.(xauthsigning.Tx)
	if !ok {
		return nil, errors.New("tx to stdtx error")
	}
	stdTx, err := clienttx.ConvertTxToStdTx(clientCtx.LegacyAmino, signingTx)
	if err != nil {
		return nil, err
	}
	return &stdTx, nil
}

func termintTx2CosmosTx(signTxs ttypes.Tx) (sdk.Tx, error) {
	return chainnet.ChainNetDst.GetClientCtx().TxConfig.TxDecoder()(signTxs)
}

func signTx2Bytes(signTxs xauthsigning.Tx) ([]byte, error) {
	return chainnet.ChainNetDst.GetClientCtx().TxConfig.TxEncoder()(signTxs)
}
