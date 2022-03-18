package cogs

type (
	CogDirectory map[string]CogFunc
)

func GetCog(cogName string) CogFunc {
	if cogName == "lost_ark" {
		return LostArkCog
	}

	return nil
}
