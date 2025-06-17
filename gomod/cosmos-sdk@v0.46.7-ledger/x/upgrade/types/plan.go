package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"path"
	"strings"
)

// UpgradeInfoFileName file to store upgrade information
const UpgradeInfoFilename = "upgrade-info.json"

func (p Plan) String() string {
	due := p.DueAt()
	return fmt.Sprintf(`Upgrade Plan
  Name: %s
  %s
  Info: %s.`, p.Name, due, p.Info)
}

// ValidateBasic does basic validation of a Plan
func (p Plan) ValidateBasic() error {
	if !p.Time.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("time-based upgrades have been deprecated in the SDK")
	}
	if p.UpgradedClientState != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("upgrade logic for IBC has been moved to the IBC module")
	}
	if len(p.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}
	if p.Height <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "height must be greater than 0")
	}
	if !p.VerifyBinaries() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Unsupported download package")
	}

	return nil
}

// ShouldExecute returns true if the Plan is ready to execute given the current context
func (p Plan) ShouldExecute(ctx sdk.Context) bool {
	if p.Height > 0 {
		return p.Height <= ctx.BlockHeight()
	}
	return false
}

// DueAt is a string representation of when this plan is due to be executed
func (p Plan) DueAt() string {
	return fmt.Sprintf("height: %d", p.Height)
}

func (p Plan) UpgradeInfo() (UpgradeInfo, error) {
	upgradeInfo := UpgradeInfo{}
	err := json.Unmarshal([]byte(p.Info), &upgradeInfo)
	if err != nil {
		return upgradeInfo, err
	}
	return upgradeInfo, nil
}

func (p Plan) ChainUpgradeInfo(blockchainInfo UpgradeConfig) (string, error) {
	upgradeInfo := ChainUpgradeConfig{}
	binaries := make(map[string]string)
	for key, val := range blockchainInfo.Binaries {
		binaries[key] = val.Url
	}
	upgradeInfo.Binaries = binaries
	bytes, err := json.Marshal(upgradeInfo)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p Plan) VerifyBinaries() bool {
	upgradeInfo, err := p.UpgradeInfo()
	if err != nil {
		return false
	}

	if upgradeInfo.Gateway == nil && upgradeInfo.App == nil && upgradeInfo.Blockchain == nil {
		return false
	}

	if upgradeInfo.Gateway != nil {
		if len(upgradeInfo.Gateway.Binaries) == 0 {
			return false
		}

		for os, binaryInfo := range upgradeInfo.Gateway.Binaries {
			if os != "windows/amd64" && os != "linux/amd64" {
				return false
			}

			if binaryInfo.Url == "" || binaryInfo.Hash == "" || binaryInfo.Size == 0 || len(binaryInfo.Hash) != 128 {
				return false
			}
		}

		if _, ok := upgradeInfo.Gateway.Binaries["windows/amd64"]; ok {
			flag := verifySuffix(upgradeInfo.Gateway.Binaries["windows/amd64"].Url)
			if !flag {
				return flag
			}
		}
		if _, ok := upgradeInfo.Gateway.Binaries["linux/amd64"]; ok {
			flag := verifySuffix(upgradeInfo.Gateway.Binaries["linux/amd64"].Url)
			if !flag {
				return flag
			}
		}
	}
	if upgradeInfo.Blockchain != nil {
		if len(upgradeInfo.Blockchain.Binaries) == 0 {
			return false
		}

		for os, binaryInfo := range upgradeInfo.Blockchain.Binaries {
			if os != "windows/amd64" && os != "linux/amd64" {
				return false
			}

			if binaryInfo.Url == "" || binaryInfo.Hash == "" || binaryInfo.Size == 0 || len(binaryInfo.Hash) != 128 {
				return false
			}
		}

		if _, ok := upgradeInfo.Blockchain.Binaries["windows/amd64"]; ok {
			flag := verifySuffix(upgradeInfo.Blockchain.Binaries["windows/amd64"].Url)
			if !flag {
				return flag
			}
		}
		if _, ok := upgradeInfo.Blockchain.Binaries["linux/amd64"]; ok {
			flag := verifySuffix(upgradeInfo.Blockchain.Binaries["linux/amd64"].Url)
			if !flag {
				return flag
			}
		}
	}

	if upgradeInfo.App != nil {
		if len(upgradeInfo.App.Binaries) == 0 {
			return false
		}

		for os, binaryInfo := range upgradeInfo.App.Binaries {
			if os != "android" {
				return false
			}

			ok := verifySuffixApk(binaryInfo.Url)
			if !ok {
				return ok
			}

			if binaryInfo.Url == "" || binaryInfo.Hash == "" || binaryInfo.Size == 0 || len(binaryInfo.Hash) != 128 {
				return false
			}
		}
	}

	return true
}

func verifySuffix(url string) bool {
	suffix := strings.Split(path.Ext(url), "?")[0]
	if suffix == ".zip" || suffix == ".tar" || suffix == ".gz" {
		return true
	}
	return false
}

func verifySuffixApk(url string) bool {
	suffix := strings.Split(path.Ext(url), "?")[0]
	if suffix == ".apk" {
		return true
	}
	return false
}

type UpgradeConfig struct {
	Binaries      map[string]BinariesInfo `json:"binaries"`
	Version       string                  `json:"version"`
	UpgradeScript map[string]string       `json:"upgrade_script"`
}

type BinariesInfo struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
	Size uint64 `json:"size"`
}

type ChainUpgradeConfig struct {
	Binaries map[string]string `json:"binaries"`
}

type UpgradeInfo struct {
	Gateway    *UpgradeConfig `json:"gateway"`
	Blockchain *UpgradeConfig `json:"blockchain"`
	App        *UpgradeConfig `json:"app"`
}
