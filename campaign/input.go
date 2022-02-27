package campaign

import "api/user"

// struct ini kita jadikan parameter dari function GetCampaignByID yang ada di dalam service
type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateCampaignInput struct {
	// mapping inputan dari user ke dalam struct ini
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	// setelah itu kita mapping data user
	// data User dibutuhkan untuk mengetahui siapa yang membuat campaign
	// kemudian user id nya digunakan untuk pembuatan slug agar unique
	User user.User
	// isi field User ini dengan mengambil data user yang ada di Context yang sebelumnya sudah di set lewat middleware
}

// data yang dipakai untuk menangkap inputan dari user
type CreateCampaignImageInput struct {
	CampaignID int  `form:"campaign_id" binding:"required"`
	IsPrimary  bool `form:"is_primary"`
	User       user.User
}
