package campaign

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

// buat function untuk mengubah yang tadinya object struct campaign yang ada di entity menjadi struct campaign formatter
// parameternya adalah campaign yang ada di entity dan balikannya adalah CampaignFormatter
func FormatCampaign(campaign Campaign) CampaignFormatter {
	// mapping dari object campaign entity ke object campaign formatter
	// buat objectnya terlebih dahulu
	campaignFormatter := CampaignFormatter{}
	// kemudian mulai mappingnya
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.Slug = campaign.Slug
	// setting imageURL jika tidak memiliki gambar atau memiliki gambar
	campaignFormatter.ImageURL = ""
	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

// format apabila campaign lebih dari 1 maka lakukan perulangan untunk memformat tiap campaign
func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	// []CampaignFormatter{} memiliki nilai default jika campaign-nya tidak ada atau bernilai nol, maka kita kembalikan sebagai slice/array
	campaignsFormatter := []CampaignFormatter{}

	// di setiap perulangan kita dapatkan single object campaign
	// dari situ kita ubah menjadi struct campaign formatter menggunakkan fungsi FormatCampaign
	for _, campaign := range campaigns {
		// jika sudah mendapatkan campaign formatternya, kita masukkan ke dalam slice/array campaignsFormatter
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}
