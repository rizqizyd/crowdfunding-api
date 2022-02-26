package campaign

import "strings"

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

// format detail campaign
type CampaignDetailFormatter struct {
	ID               int                      `json:"id"`
	Name             string                   `json:"name"`
	ShortDescription string                   `json:"short_description"`
	Description      string                   `json:"description"`
	ImageURL         string                   `json:"image_url"`
	GoalAmount       int                      `json:"goal_amount"`
	CurrentAmount    int                      `json:"current_amount"`
	UserID           int                      `json:"user_id"`
	Slug             string                   `json:"slug"`
	Perks            []string                 `json:"perks"`
	User             CampaignUserFormatter    `json:"user"`
	Images           []CampaignImageFormatter `json:"images`
}

// struct user yang akan berada di dalam CampaignDetailFormatter
type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// struct images
type CampaignImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	// campaignDetailFormatter adalah tipe dari CampaignDetailFormatter{}
	campaignDetailFormatter := CampaignDetailFormatter{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.GoalAmount = campaign.GoalAmount
	campaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	campaignDetailFormatter.UserID = campaign.UserID
	campaignDetailFormatter.Slug = campaign.Slug
	campaignDetailFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	// format perks (memecah string berdasarkan koma)
	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}

	// set perks berdasarkan looping yang telah dibuat
	campaignDetailFormatter.Perks = perks

	// format user formatter
	user := campaign.User
	campaignUserFormatter := CampaignUserFormatter{}
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageURL = user.AvatarFileName

	// set/isi field user ke dalam CampaignDetailFormatter
	campaignDetailFormatter.User = campaignUserFormatter

	// format images
	images := []CampaignImageFormatter{}
	// ambil campaign images
	for _, image := range campaign.CampaignImages {
		campaignImageFormatter := CampaignImageFormatter{}
		campaignImageFormatter.ImageURL = image.FileName

		// ubah isPrimary int ke boolean
		isPrimary := false
		if image.IsPrimary == 1 {
			isPrimary = true
		}
		campaignImageFormatter.IsPrimary = isPrimary

		// masukkan ke dalam slice/array images
		images = append(images, campaignImageFormatter)
	}

	// set field images dari campaignDetailFormatter
	campaignDetailFormatter.Images = images

	return campaignDetailFormatter
}
