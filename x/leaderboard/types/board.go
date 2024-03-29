/**
 * Created by zhouwenzhe on 2023/6/1
 */

package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"sort"
	"time"
)

func ParseDateAddedAsTime(dateAdded string) (dateAddedParsed time.Time, err error) {
	dateAddedParsed, errDateAdded := time.Parse(TimeLayout, dateAdded)
	return dateAddedParsed, sdkerrors.Wrapf(errDateAdded, ErrInvalidDateAdded.Error(), dateAdded)
}

func SortPlayerInfo(playerInfoList []PlayerInfo) {
	sort.SliceStable(playerInfoList[:], func(i, j int) bool {
		if playerInfoList[i].WonCount > playerInfoList[j].WonCount {
			return true
		}
		if playerInfoList[i].WonCount < playerInfoList[j].WonCount {
			return false
		}
		firstPlayerTime, _ := ParseDateAddedAsTime(playerInfoList[i].DateUpdated)
		secondPlayerTime, _ := ParseDateAddedAsTime(playerInfoList[j].DateUpdated)

		return firstPlayerTime.After(secondPlayerTime)
	})
}
