package goinsta

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"time"
)

// Users is a struct that stores many user's returned by many different methods.
type Users struct {
	insta *Instagram

	// It's a bit confusing have the same structure
	// in the Instagram strucure and in the multiple users
	// calls

	err      error
	endpoint string

	Status    string          `json:"status"`
	BigList   bool            `json:"big_list"`
	Users     []*User         `json:"users"`
	PageSize  int             `json:"page_size"`
	RawNextID json.RawMessage `json:"next_max_id"`
	NextID    string          `json:"-"`
}

func newUsers(insta *Instagram) *Users {
	users := &Users{insta: insta}

	return users
}

// SetInstagram sets new instagram to user structure
func (users *Users) SetInstagram(insta *Instagram) {
	users.insta = insta
}

// Next allows to paginate after calling:
// Account.Follow* and User.Follow*
//
// New user list is stored inside Users
//
// returns false when list reach the end.
func (users *Users) Next() bool {
	if users.err != nil {
		return false
	}

	insta := users.insta
	endpoint := users.endpoint

	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			Query: map[string]string{
				"max_id":             users.NextID,
				"ig_sig_key_version": instaSigKeyVersion,
				"rank_token":         insta.rankToken,
			},
		},
	)
	if err == nil {
		usrs := Users{}
		err = json.Unmarshal(body, &usrs)
		if err == nil {
			if len(usrs.RawNextID) > 0 && usrs.RawNextID[0] == '"' && usrs.RawNextID[len(usrs.RawNextID)-1] == '"' {
				if err := json.Unmarshal(usrs.RawNextID, &usrs.NextID); err != nil {
					users.err = err
					return false
				}
			} else if usrs.RawNextID != nil {
				var nextID int64
				if err := json.Unmarshal(usrs.RawNextID, &nextID); err != nil {
					users.err = err
					return false
				}
				usrs.NextID = strconv.FormatInt(nextID, 10)
			}
			*users = usrs
			if !usrs.BigList || usrs.NextID == "" {
				users.err = ErrNoMore
			}
			users.insta = insta
			users.endpoint = endpoint
			users.setValues()
			return true
		}
	}
	users.err = err
	return false
}

// Error returns users error
func (users *Users) Error() error {
	return users.err
}

func (users *Users) setValues() {
	for i := range users.Users {
		users.Users[i].insta = users.insta
	}
}

type userResp struct {
	Status string `json:"status"`
	User   User   `json:"user"`
}

// User is the representation of instagram's user profile
type User struct {
	insta       *Instagram
	Collections *Collections

	ID                         int64         `json:"pk"`
	Username                   string        `json:"username"`
	FullName                   string        `json:"full_name"`
	Biography                  string        `json:"biography"`
	BestiesCount               int           `json:"besties_count"`
	ShowBestiesBadge           bool          `json:"show_besties_badge"`
	RecentlyBestiedByCount     int           `json:"recently_bestied_by_count"`
	ProfilePicURL              string        `json:"profile_pic_url"`
	Email                      string        `json:"email"`
	PhoneNumber                string        `json:"phone_number"`
	WhatsappNumber             string        `json:"whatsapp_number"`
	IsBusiness                 bool          `json:"is_business"`
	AccountType                int           `json:"account_type"`
	AccountBadges              []interface{} `json:"account_badges"`
	Gender                     int           `json:"gender"`
	ProfilePicID               string        `json:"profile_pic_id"`
	FbIdV2                     int64         `json:"fbid_v2"`
	HasAnonymousProfilePicture bool          `json:"has_anonymous_profile_picture"`
	IsPrivate                  bool          `json:"is_private"`
	IsUnpublished              bool          `json:"is_unpublished"`
	IsMutedWordsGlobalEnabled  bool          `json:"is_muted_words_global_enabled"`
	IsMutedWordsCustomEnabled  bool          `json:"is_muted_words_custom_enabled"`
	AllowedCommenterType       string        `json:"allowed_commenter_type"`
	UserTagsCount              int           `json:"usertags_count"`
	UserTagReviewEnabled       bool          `json:"usertag_review_enabled"`
	IsVerified                 bool          `json:"is_verified"`
	IsNeedy                    bool          `json:"is_needy"`
	IsInterestAccount          bool          `json:"is_interest_account"`
	IsVideoCreator             bool          `json:"is_video_creator"`
	MediaCount                 int           `json:"media_count"`
	IGTVCount                  int           `json:"total_igtv_videos"`
	HasIGTVSeries              bool          `json:"has_igtv_series"`
	TotalClipCount             int           `json:"total_clips_count"`
	TotalAREffects             int           `json:"total_ar_effects"`
	FollowerCount              int           `json:"follower_count"`
	FollowingCount             int           `json:"following_count"`
	FollowingTagCount          int           `json:"following_tag_count"`
	MutualFollowersID          []int64       `json:"profile_context_mutual_follow_ids"`
	FollowFrictionType         int           `json:"follow_friction_type"`
	ProfileContext             string        `json:"profile_context"`
	GeoMediaCount              int           `json:"geo_media_count"`
	ExternalURL                string        `json:"external_url"`
	HasBiographyTranslation    bool          `json:"has_biography_translation"`
	HasVideos                  bool          `json:"has_videos"`
	HasProfileVideoFeed        bool          `json:"has_profile_video_feed"`
	HasSavedItems              bool          `json:"has_saved_items"`
	ExternalLynxURL            string        `json:"external_lynx_url"`
	BiographyWithEntities      struct {
		RawText  string        `json:"raw_text"`
		Entities []interface{} `json:"entities"`
	} `json:"biography_with_entities"`
	Nametag                        Nametag `json:"nametag"`
	HasChaining                    bool    `json:"has_chaining"`
	HasPlacedOrders                bool    `json:"has_placed_orders"`
	IsFavorite                     bool    `json:"is_favorite"`
	IsFavoriteForStories           bool    `json:"is_favorite_for_stories"`
	IsFavoriteForHighlights        bool    `json:"is_favorite_for_highlights"`
	IsProfileActionNeeded          bool    `json:"is_profile_action_needed"`
	CanBeReportedAsFraud           bool    `json:"can_be_reported_as_fraud"`
	CanBoostPosts                  bool    `json:"can_boost_posts"`
	CanSeeOrganicInsights          bool    `json:"can_see_organic_insights"`
	CanConvertToBusiness           bool    `json:"can_convert_to_business"`
	CanCreateSponsorTags           bool    `json:"can_create_sponsor_tags"`
	CanCreateNewFundraiser         bool    `json:"can_create_new_standalone_fundraiser"`
	CanCreateNewPersonalFundraiser bool    `json:"can_create_new_standalone_personal_fundraiser"`
	CanBeTaggedAsSponsor           bool    `json:"can_be_tagged_as_sponsor"`
	CanSeeSupportInbox             bool    `json:"can_see_support_inbox"`
	CanSeeSupportInboxV1           bool    `json:"can_see_support_inbox_v1"`
	CanTagProductsFromMerchants    bool    `json:"can_tag_products_from_merchants"`
	CanSeePrimaryCountryInsettings bool    `json:"can_see_primary_country_in_settings"`
	CanFollowHashtag               bool    `json:"can_follow_hashtag"`
	PersonalAccountAdsPageName     string  `json:"personal_account_ads_page_name"`
	PersonalAccountAdsId           string  `json:"personal_account_ads_page_id"`
	ShowShoppableFeed              bool    `json:"show_shoppable_feed"`
	ShowInsightTerms               bool    `json:"show_insights_terms"`
	ShowConversionEditEntry        bool    `json:"show_conversion_edit_entry"`
	ShowPostsInsightEntryPoint     bool    `json:"show_post_insights_entry_point"`
	ShoppablePostsCount            int     `json:"shoppable_posts_count"`
	RequestContactEnabled          bool    `json:"request_contact_enabled"`
	FeedPostReshareDisabled        bool    `json:"feed_post_reshare_disabled"`
	CreatorShoppingInfo            struct {
		LinkedMerchantAccounts []interface{} `json:"linked_merchant_accounts"`
	} `json:"creator_shopping_info"`
	StandaloneFundraiserInfo struct {
		HasActiveFundraiser                 bool        `json:"has_active_fundraiser"`
		FundraiserId                        int64       `json:"fundraiser_id"`
		FundraiserTitle                     string      `json:"fundraiser_title"`
		FundraiserType                      interface{} `json:"fundraiser_type"`
		FormattedGoalAmount                 string      `json:"formatted_goal_amount"`
		BeneficiaryUsername                 string      `json:"beneficiary_username"`
		FormattedFundraiserProgressInfoText string      `json:"formatted_fundraiser_progress_info_text"`
		PercentRaised                       interface{} `json:"percent_raised"`
	} `json:"standalone_fundraiser_info"`
	AggregatePromoteEngagement   bool         `json:"aggregate_promote_engagement"`
	AllowMentionSetting          string       `json:"allow_mention_setting"`
	AllowTagSetting              string       `json:"allow_tag_setting"`
	LimitedInteractionsEnabled   bool         `json:"limited_interactions_enabled"`
	ReelAutoArchive              string       `json:"reel_auto_archive"`
	HasHighlightReels            bool         `json:"has_highlight_reels"`
	HightlightReshareDisabled    bool         `json:"highlight_reshare_disabled"`
	IsMemorialized               bool         `json:"is_memorialized"`
	HasGuides                    bool         `json:"has_guides"`
	PublicEmail                  string       `json:"public_email"`
	PublicPhoneNumber            string       `json:"public_phone_number"`
	PublicPhoneCountryCode       string       `json:"public_phone_country_code"`
	ContactPhoneNumber           string       `json:"contact_phone_number"`
	CityID                       int64        `json:"city_id"`
	CityName                     string       `json:"city_name"`
	AddressStreet                string       `json:"address_street"`
	DirectMessaging              string       `json:"direct_messaging"`
	Latitude                     float64      `json:"latitude"`
	Longitude                    float64      `json:"longitude"`
	Category                     string       `json:"category"`
	BusinessContactMethod        string       `json:"business_contact_method"`
	IncludeDirectBlacklistStatus bool         `json:"include_direct_blacklist_status"`
	HdProfilePicURLInfo          PicURLInfo   `json:"hd_profile_pic_url_info"`
	HdProfilePicVersions         []PicURLInfo `json:"hd_profile_pic_versions"`
	School                       School       `json:"school"`
	Byline                       string       `json:"byline"`
	SocialContext                string       `json:"social_context,omitempty"`
	SearchSocialContext          string       `json:"search_social_context,omitempty"`
	MutualFollowersCount         float64      `json:"mutual_followers_count"`
	LatestReelMedia              int64        `json:"latest_reel_media,omitempty"`
	IsCallToActionEnabled        bool         `json:"is_call_to_action_enabled"`
	IsPotentialBusiness          bool         `json:"is_potential_business"`
	FbPageCallToActionID         string       `json:"fb_page_call_to_action_id"`
	FbPayExperienceEnabled       bool         `json:"fbpay_experience_enabled"`
	Zip                          string       `json:"zip"`
	Friendship                   Friendship   `json:"friendship_status"`
	AutoExpandChaining           bool         `json:"auto_expand_chaining"`

	AllowedToCreateNonprofitFundraisers        bool          `json:"is_allowed_to_create_standalone_nonprofit_fundraisers"`
	AllowedToCreatePersonalFundraisers         bool          `json:"is_allowed_to_create_standalone_personal_fundraisers"`
	IsElegibleToShowFbCrossSharingNux          bool          `json:"is_eligible_to_show_fb_cross_sharing_nux"`
	PageIdForNewSumaBizAccount                 interface{}   `json:"page_id_for_new_suma_biz_account"`
	ElegibleShoppingSignupEntrypoints          []interface{} `json:"eligible_shopping_signup_entrypoints"`
	IsIgdProductPickerEnabled                  bool          `json:"is_igd_product_picker_enabled"`
	IsElegibleForAffiliateShopOnboarding       bool          `json:"is_eligible_for_affiliate_shop_onboarding"`
	ElegibleShoppingFormats                    []interface{} `json:"eligible_shopping_formats"`
	NeedsToAcceptShoppingSellerOnboardingTerms bool          `json:"needs_to_accept_shopping_seller_onboarding_terms"`
	IsShoppingCatalogSettingsEnabled           bool          `json:"is_shopping_settings_enabled"`
	IsShoppingCommunityContentEnabled          bool          `json:"is_shopping_community_content_enabled"`
	IsShoppingAutoHighlightEnabled             bool          `json:"is_shopping_auto_highlight_eligible"`
	IsShoppingCatalogSourceSelectionEnabled    bool          `json:"is_shopping_catalog_source_selection_enabled"`
	ProfessionalConversionSuggestedAccountType int           `json:"professional_conversion_suggested_account_type"`
	InteropMessagingUserfbid                   int64         `json:"interop_messaging_user_fbid"`
	LinkedFbInfo                               struct{}      `json:"linked_fb_info"`
	HasElegibleWhatsappLinkingCategory         struct{}      `json:"has_eligible_whatsapp_linking_category"`
	ExistingUserAgeCollectionEnabled           bool          `json:"existing_user_age_collection_enabled"`
	AboutYourAccountBloksEntrypointEnabled     bool          `json:"about_your_account_bloks_entrypoint_enabled"`
	OpenExternalUrlWithInAppBrowser            bool          `json:"open_external_url_with_in_app_browser"`

	// Profile picture as raw bytes, to populate call User.DownloadProfilePic()
	ProfilePic []byte
}

// SetInstagram will update instagram instance for selected User.
func (user *User) SetInstagram(insta *Instagram) {
	user.insta = insta
}

// NewUser returns prepared user to be used with his functions.
func (insta *Instagram) NewUser() *User {
	return &User{insta: insta}
}

// Info updates user info
// extra query arguments can be passes one after another as func(key, value).
// Only if an even number of string arguements will be passed, they will be
//   used in the query.
//
// See example: examples/user/friendship.go
func (user *User) Info(params ...interface{}) error {
	insta := user.insta
	query := map[string]string{}
	if len(params)%2 == 0 {
		for i := 0; i < len(params); i = i + 2 {
			query[params[i].(string)] = params[i+1].(string)
		}
	}

	body, _, err := insta.sendRequest(&reqOptions{
		Endpoint: fmt.Sprintf(urlUserInfo, user.ID),
		Query:    query,
	})
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, user)
	return err
}

// Sync wraps User.Info() 1:1
func (user *User) Sync(params ...interface{}) error {
	return user.Info(params...)
}

// Following returns a list of user following.
//
// Users.Next can be used to paginate
//
// See example: examples/user/following.go
func (user *User) Following() *Users {
	/*
		TODO: call friendships/show_many/

		returns:
		{
		  "friendship_statuses": {
			"xxxxxxx": {
			  "following": true,
			  "incoming_request": false,
			  "is_bestie": false,
			  "is_private": true,
			  "is_restricted": false,
			  "outgoing_request": false,
			  "is_feed_favorite": false
			},
			...
		}
	*/
	users := &Users{}
	users.insta = user.insta
	users.endpoint = fmt.Sprintf(urlFollowing, user.ID)
	return users
}

// Followers returns a list of user followers.
//
// Users.Next can be used to paginate
//
// See example: examples/user/followers.go
func (user *User) Followers() *Users {
	users := &Users{}
	users.insta = user.insta
	users.endpoint = fmt.Sprintf(urlFollowers, user.ID)
	return users
}

// Block blocks user
//
// This function updates current User.Friendship structure.
// Param: autoBlock - automatically block accounts registered on the same email/number
//
// See example: examples/user/block.go
func (user *User) Block(autoBlock bool) error {
	insta := user.insta
	data, err := json.Marshal(map[string]string{
		"surface":              "profile",
		"is_autoblock_enabled": strconv.FormatBool(autoBlock),
		"user_id":              strconv.Itoa(int(user.ID)),
		"_uid":                 strconv.Itoa(int(insta.Account.ID)),
		"_uuid":                insta.uuid,
	})
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserBlock, user.ID),
			IsPost:   true,
			Query: map[string]string{
				"signed_body": "SIGNATURE." + string(data),
			},
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

// Unblock unblocks user
//
// This function updates current User.Friendship structure.
//
// See example: examples/user/unblock.go
func (user *User) Unblock() error {
	insta := user.insta
	data, err := json.Marshal(
		map[string]interface{}{
			"user_id":          user.ID,
			"_uid":             insta.Account.ID,
			"_uuid":            insta.uuid,
			"container_module": "blended_search",
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserUnblock, user.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

// Mute mutes user from appearing in the feed or story reel
//
// Use one of the pre-defined constants to choose what exactly to mute:
// goinsta.MuteAll, goinsta.MuteStory, goinsta.MuteFeed
// This function updates current User.Friendship structure.
func (user *User) Mute(opt muteOption) error {
	if opt == MuteAll {
		err := user.muteOrUnmute(MuteStory, urlUserMute)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		return user.muteOrUnmute(MutePosts, urlUserMute)
	}
	return user.muteOrUnmute(opt, urlUserMute)
}

// Unmute unmutes user so it appears in the feed or story reel again
//
// Use one of the pre-defined constants to choose what exactly to unmute:
// goinsta.MuteAll, goinsta.MuteStory, goinsta.MuteFeed
// This function updates current User.Friendship structure.
func (user *User) Unmute(opt muteOption) error {
	if opt == MuteAll {
		err := user.muteOrUnmute(MuteStory, urlUserUnmute)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		return user.muteOrUnmute(MutePosts, urlUserUnmute)
	}
	return user.muteOrUnmute(opt, urlUserUnmute)
}

func (user *User) muteOrUnmute(opt muteOption, endpoint string) error {
	insta := user.insta
	data, err := json.Marshal(generateMuteData(user, opt))
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: endpoint,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

func generateMuteData(user *User, opt muteOption) map[string]string {
	insta := user.insta
	data := map[string]string{
		"_uid":             toString(insta.Account.ID),
		"_uuid":            insta.uuid,
		"container_module": "media_mute_sheet",
	}

	switch opt {
	case MuteStory:
		data["target_reel_author_id"] = toString(user.ID)
	case MutePosts:
		data["target_posts_author_id"] = toString(user.ID)
	}

	return data
}

// Follow started following some user
//
// This function performs a follow call. If user is private
// you have to wait until he/she accepts you.
//
// If the account is public User.Friendship will be updated
//
// See example: examples/user/follow.go
func (user *User) Follow() error {
	insta := user.insta
	data, err := json.Marshal(
		map[string]string{
			"user_id":    toString(user.ID),
			"radio_type": "wifi-none",
			"_uid":       toString(insta.Account.ID),
			"device_id":  insta.dID,
			"_uuid":      insta.uuid,
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserFollow, user.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

// Unfollow unfollows user
//
// User.Friendship will be updated
//
// See example: examples/user/unfollow.go
func (user *User) Unfollow() error {
	insta := user.insta
	data, err := json.Marshal(
		map[string]string{
			"user_id":          toString(user.ID),
			"radio_type":       "wifi-none",
			"_uid":             toString(insta.Account.ID),
			"device_id":        insta.dID,
			"_uuid":            insta.uuid,
			"container_module": "following_sheet",
		},
	)
	if err != nil {
		return err
	}
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserUnfollow, user.ID),
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)
	if err != nil {
		return err
	}
	resp := friendResp{}
	err = json.Unmarshal(body, &resp)
	user.Friendship = resp.Friendship
	if err != nil {
		return err
	}

	return nil
}

// GetFriendship allows user to get friend relationship.
//
// The result is stored in user.Friendship
func (user *User) GetFriendship() (fr *Friendship, err error) {
	insta := user.insta
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlFriendship, user.ID),
		},
	)
	if err == nil {
		fr = &user.Friendship
		err = json.Unmarshal(body, fr)
	}
	return
}

// GetFeaturedAccounts will call the featured accounts enpoint.
func (user *User) GetFeaturedAccounts() ([]*User, error) {
	body, _, err := user.insta.sendRequest(&reqOptions{
		Endpoint: urlFeaturedAccounts,
		Query: map[string]string{
			"target_user_id": strconv.Itoa(int(user.ID)),
		},
	})
	d := struct {
		Accounts []*User `json:"accounts"`
		Status   string  `json:"status"`
	}{}
	err = json.Unmarshal(body, &d)
	return d.Accounts, err
}

// Feed returns user feeds (media)
//
// 	params can be:
// 		string: timestamp of the minimum media timestamp.
//
// For pagination use FeedMedia.Next()
//
// See example: examples/user/feed.go
func (user *User) Feed(params ...interface{}) *FeedMedia {
	insta := user.insta

	media := &FeedMedia{
		insta:    insta,
		endpoint: urlUserFeed,
		uid:      user.ID,
	}

	for _, param := range params {
		switch s := param.(type) {
		case string:
			media.timestamp = s
		}
	}

	return media
}

// Tags returns media where user is tagged in
//
// For pagination use FeedMedia.Next()
//
// See example: examples/user/tags.go
func (user *User) Tags(minTimestamp []byte) (*FeedMedia, error) {
	insta := user.insta

	timestamp := string(minTimestamp)
	body, _, err := insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(urlUserTags, user.ID),
			Query: map[string]string{
				"max_id":         "",
				"rank_token":     user.insta.rankToken,
				"min_timestamp":  timestamp,
				"ranked_content": "true",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	media := &FeedMedia{
		insta:    insta,
		endpoint: urlUserTags,
		uid:      user.ID,
	}
	err = json.Unmarshal(body, media)
	if err != nil {
		return nil, err
	}
	media.setValues()
	return media, nil
}

// DownloadProfilePic will download a user's profile picture if available, and
//   return it as a byte slice.
func (user *User) DownloadProfilePic() ([]byte, error) {
	if user.ProfilePicURL == "" {
		return nil, ErrNoProfilePicUrl
	}
	insta := user.insta
	b, err := insta.download(user.ProfilePicURL)
	if err != nil {
		return nil, err
	}
	user.ProfilePic = b
	return b, nil
}

// DownloadProfilePicTo will download the user profile picture to the provided
//   path. If path does not include a file name, one will be extracted automatically.
// File extention does not need to be set, and will be set automatically.
func (user *User) DownloadProfilePicTo(dst string) error {
	folder, fn := path.Split(dst)
	b, err := user.DownloadProfilePic()
	if err != nil {
		return nil
	}
	fn, err = getDownloadName(user.ProfilePicURL, fn)
	if err != nil {
		return err
	}
	err = saveToFolder(folder, fn, b)
	return err
}
