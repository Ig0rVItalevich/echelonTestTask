package server

import (
	"context"
	api "github.com/Ig0rVItalevich/echelon/pkg/api/proto"
	"github.com/Ig0rVItalevich/echelon/pkg/thumbnail"
	"google.golang.org/api/youtube/v3"
)

type Server struct {
	api.UnimplementedThumbnailsServer
	YoutubeService *youtube.Service
}

func NewServer(youtubeService *youtube.Service) *Server {
	return &Server{YoutubeService: youtubeService}
}

func (s *Server) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	part := []string{"snippet"}
	videoInfo, err := s.YoutubeService.Videos.List(part).Fields("items/snippet/thumbnails").Id(request.GetLink()).Do()
	if err != nil {
		return nil, err
	}
	thumbnails := videoInfo.Items[0].Snippet.Thumbnails

	thumbnailResult, err := thumbnail.ThumbnailHandler(thumbnails, request.GetLink())
	if err != nil {
		return nil, err
	}

	response := api.GetResponse{Thumbnail: thumbnailResult}

	return &response, nil
}
