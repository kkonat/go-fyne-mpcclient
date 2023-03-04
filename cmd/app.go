package main

import (
	"context"
)

type ControlCenterApp struct {
	mpdClient  *Client
	ctrlClient *Client
	ctx        context.Context
	state      *playerState
}

func NewControlCenterApp(ctx context.Context, updateStream chan any) *ControlCenterApp {
	app := &ControlCenterApp{
		mpdClient:  NewClient("192.168.0.95:6600"),
		ctrlClient: NewClient("192.168.0.95:1025"),
		ctx:        ctx,
	}
	app.state = NewState(app)
	return app
}
