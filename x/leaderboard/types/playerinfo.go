/**
 * Created by zhouwenzhe on 2023/6/1
 */

package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (info PlayerInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(info.Index)
	if err != nil {
		return err
	}
	return nil
}
