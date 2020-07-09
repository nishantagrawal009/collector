package profefe

import (
	"fmt"
	"testing"
)

func TestCache_writeAndGet(t *testing.T) {
	c:= newCache()
	c.PutProfilesIds("profile-push","abc","1")
	dashBoard,_ := c.GetProfileIds()
	fmt.Println(dashBoard)
}
func TestCache_writeAndGet(t *testing.T) {
	c:= newCache()
	c.PutProfilesIds("profile-push","abc","1")
	dashBoard,_ := c.GetProfileIds()
	fmt.Println(dashBoard)
}
