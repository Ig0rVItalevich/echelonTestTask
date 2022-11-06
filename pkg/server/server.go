package server

import (
	"context"
	api "github.com/Ig0rVItalevich/echelon/pkg/api/proto"
	"github.com/Ig0rVItalevich/echelon/pkg/cache"
	"github.com/Ig0rVItalevich/echelon/pkg/thumbnail"
	"google.golang.org/api/youtube/v3"
)

type Server struct {
	api.UnimplementedThumbnailsServer
	YoutubeService *youtube.Service
	Cache          *cache.Cache
}

var _ api.ThumbnailsServer = (*Server)(nil)

func (s *Server) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	thumbnailBytes, err := s.Cache.GetThumbnail(request.GetVideoId())
	if err != nil {
		part := []string{"snippet"}
		videoInfo, err := s.YoutubeService.Videos.List(part).Fields("items/snippet/thumbnails").Id(request.GetVideoId()).Do()
		if err != nil {
			return nil, err
		}
		thumbnails := videoInfo.Items[0].Snippet.Thumbnails

		thumbnailBytes, err = thumbnail.ThumbnailHandler(thumbnails)
		if err != nil {
			return nil, err
		}

		if err := s.Cache.SetThumbnail(request.GetVideoId(), thumbnailBytes); err != nil {
			return nil, err
		}
	}

	thumbnailPath, err := thumbnail.SaveThumbnail(thumbnailBytes, request.GetVideoId())
	if err != nil {
		return nil, err
	}
	response := api.GetResponse{Thumbnail: thumbnailPath}

	return &response, nil
}

func NewServer(youtubeService *youtube.Service, cache *cache.Cache) *Server {
	return &Server{YoutubeService: youtubeService, Cache: cache}
}
