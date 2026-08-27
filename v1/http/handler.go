package http

// func (s *Server) RoutesHandler(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	path := c.Request.URL.Path
// 	route, err := s.RouteRepo.GetRouteByPath(ctx, path)
// 	if err != nil {
// 		responseJSON(c, http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
// 		return
// 	}
// }

// func (r *RouteRepository) GetRouteByPath(ctx context.Context, path string) (*utilsHeader.RouteConfig, error) {
// 	if path == "" {
// 		return nil, errors.New("path is empty")
// 	}
// 	key := "gateway:api:" + strings.ReplaceAll(strings.TrimPrefix(path, "/"), "/", ":")
// 	payload, err := r.client.Client.Get(ctx, key).Result()
// }
