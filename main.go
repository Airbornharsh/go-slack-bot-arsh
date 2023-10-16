package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

// type User struct {
// 	ID                string               `json:"id"`
// 	TeamID            string               `json:"team_id"`
// 	Name              string               `json:"name"`
// 	Deleted           bool                 `json:"deleted"`
// 	Color             string               `json:"color"`
// 	RealName          string               `json:"real_name"`
// 	TZ                string               `json:"tz,omitempty"`
// 	TZLabel           string               `json:"tz_label"`
// 	TZOffset          int                  `json:"tz_offset"`
// 	Profile           slack.UserProfile    `json:"profile"`
// 	IsBot             bool                 `json:"is_bot"`
// 	IsAdmin           bool                 `json:"is_admin"`
// 	IsOwner           bool                 `json:"is_owner"`
// 	IsPrimaryOwner    bool                 `json:"is_primary_owner"`
// 	IsRestricted      bool                 `json:"is_restricted"`
// 	IsUltraRestricted bool                 `json:"is_ultra_restricted"`
// 	IsStranger        bool                 `json:"is_stranger"`
// 	IsAppUser         bool                 `json:"is_app_user"`
// 	IsInvitedUser     bool                 `json:"is_invited_user"`
// 	Has2FA            bool                 `json:"has_2fa"`
// 	TwoFactorType     *string              `json:"two_factor_type"`
// 	HasFiles          bool                 `json:"has_files"`
// 	Presence          string               `json:"presence"`
// 	Locale            string               `json:"locale"`
// 	Updated           slack.JSONTime       `json:"updated"`
// 	Enterprise        slack.EnterpriseUser `json:"enterprise_user,omitempty"`
// }

type UserPagination struct {
	Users []slack.User
}

func main() {
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	// channelId := os.Getenv("SLACK_CHANNEL_ID")

	client := slack.New(token, slack.OptionDebug(true))

	client.GetUsers(func(up *slack.UserPagination) {
		UserPagination := UserPagination{
			Users: up.Users,
		}

		fmt.Println(UserPagination)
	})

}
