package campaign

// struct ini kita jadikan parameter dari function GetCampaignByID yang ada di dalam service
type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"`
}
