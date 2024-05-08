package jasmine

import (
	"context"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

type Backend struct {
}

func (b *Backend) CalendarHomeSetPath(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) ListCalendars(ctx context.Context) ([]caldav.Calendar, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) GetCalendar(ctx context.Context, path string) (*caldav.Calendar, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) GetCalendarObject(ctx context.Context, path string, req *caldav.CalendarCompRequest) (*caldav.CalendarObject, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) ListCalendarObjects(ctx context.Context, path string, req *caldav.CalendarCompRequest) ([]caldav.CalendarObject, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) QueryCalendarObjects(ctx context.Context, query *caldav.CalendarQuery) ([]caldav.CalendarObject, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) PutCalendarObject(ctx context.Context, path string, calendar *ical.Calendar, opts *caldav.PutCalendarObjectOptions) (loc string, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) DeleteCalendarObject(ctx context.Context, path string) error {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}
