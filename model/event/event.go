package eventmodel

import "time"

type Event struct {
	ID string `bson:"_id,omitempty" json:"id"`
	// event contains image
	WithImage bool `bson:"with_image,omitempty" json:"with_image"`
	// ID of the user who created the event
	CreatorID string `bson:"creator_id,omitempty" json:"creator_id"`
	// for which message users: worker, all
	UserType string `bson:"user_type,omitempty" json:"user_type"`
	// caption of message
	Caption string `bson:"caption,omitempty" json:"caption"`
	// The ID of the file that will be included in the message from the bot is not required
	FileID string `bson:"file_id,omitempty" json:"file_id"`
	// buttons leading to a specific application page is not required
	WebAppBtns []WebAppBtn `bson:"web_app_btn,omitempty" json:"web_app_btn"`
	// UTC format
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at"`
}

type WebAppBtn struct {
	// inscription on the button
	Text string `bson:"text,omitempty" json:"text"`
	// link for inline button on keyboard
	InlineUrl string `bson:"inline_url,omitempty" json:"inline_url"`
}
