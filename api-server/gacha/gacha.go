package gacha

import (
	"math/rand"
	"time"
)



func DecideRank() string {
	rand.Seed(time.Now().UnixNano())

	switch percent := rand.Inin(100); percent {
	case percent < 5:

	case percent < 15:

	case percent:

	}
	
}

func SelectCharacter(string rank)  {

}