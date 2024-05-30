package payload

import (
	"testing"
)

func TestRedirects(t *testing.T) {
	t.Parallel()

	//var (
	//	fromURL   = "/test"
	//	redirects = []Redirect{
	//		{From: fromURL, To: "/new", Code: RedirectsCode301},
	//	}
	//)
	//
	//tt := map[string]struct {
	//	mock       func(cols *payloadfakes.MockCollectionService, store cache.Store)
	//	wantURL    string
	//	wantStatus int
	//}{
	//	"API error returns nil": {
	//		mock: func(cols *payloadfakes.MockCollectionService, store cache.Store) {
	//			cols.EXPECT().
	//				List(context.TODO(), CollectionRedirects, gomock.Any(), &payloadcms.ListResponse[Redirect]{}).
	//				Return(payloadcms.Response{}, errors.New("error"))
	//		},
	//		wantStatus: http.StatusOK,
	//	},
	//	"Invalid number defaults to 301": {
	//		mock: func(_ *payloadfakes.MockCollectionService, store cache.Store) {
	//			err := store.Set(context.TODO(), redirectCacheKey, []Redirect{
	//				{From: fromURL, To: "/new", Code: "wrong"},
	//			}, cache.Options{})
	//			require.NoError(t, err)
	//		},
	//		wantStatus: http.StatusMovedPermanently,
	//		wantURL:    "/new",
	//	},
	//	"No Matches": {
	//		mock: func(_ *payloadfakes.MockCollectionService, store cache.Store) {
	//			err := store.Set(context.TODO(), redirectCacheKey, []Redirect{
	//				{From: "/wrong", To: "/new", Code: RedirectsCode301},
	//			}, cache.Options{})
	//			require.NoError(t, err)
	//		},
	//		wantStatus: http.StatusOK,
	//	},
	//	"Redirects 301 from API": {
	//		mock: func(cols *payloadfakes.MockCollectionService, store cache.Store) {
	//			cols.EXPECT().
	//				List(context.TODO(), CollectionRedirects, gomock.Any(), &payloadcms.ListResponse[Redirect]{}).
	//				Do(func(_ context.Context, _ payloadcms.Collection, _ any, out any) error {
	//					*out.(*payloadcms.ListResponse[Redirect]) = payloadcms.ListResponse[Redirect]{
	//						Docs: redirects,
	//					}
	//					return nil
	//				})
	//		},
	//		wantStatus: http.StatusMovedPermanently,
	//		wantURL:    "/new",
	//	},
	//	"Redirects 301 from Cache": {
	//		mock: func(_ *payloadfakes.MockCollectionService, store cache.Store) {
	//			err := store.Set(context.TODO(), redirectCacheKey, redirects, cache.Options{})
	//			require.NoError(t, err)
	//		},
	//		wantStatus: http.StatusMovedPermanently,
	//		wantURL:    "/new",
	//	},
	//}
	//for name, test := range tt {
	//	t.Run(name, func(t *testing.T) {
	//		app := webkit.New()
	//		req := httptest.NewRequest(http.MethodGet, fromURL, nil)
	//		rr := httptest.NewRecorder()
	//
	//		ctrl := gomock.NewController(t)
	//		collections := payloadfakes.NewMockCollectionService(ctrl)
	//		store := cache.NewInMemory(time.Hour)
	//		payload := &payloadcms.Client{
	//			Collections: collections,
	//		}
	//
	//		test.mock(collections, store)
	//
	//		app.Plug(RedirectMiddleware(payload, store))
	//		app.Get(fromURL, func(c *webkit.Context) error {
	//			return c.String(http.StatusOK, "Middleware")
	//		})
	//		app.ServeHTTP(rr, req)
	//
	//		assert.Equal(t, test.wantStatus, rr.Code)
	//		assert.Equal(t, test.wantURL, rr.Header().Get("Location"))
	//	})
	//}
}
