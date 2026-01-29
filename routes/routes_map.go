package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"rocket-map-example/geo"
	"rocket-map-example/templating"
	"rocket-map-example/types"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/starfederation/datastar-go/datastar"
)

var members []types.DiscordMember

func setupMapRoute(ctx context.Context, mux *http.ServeMux) error {

	if err := loadDiscordMembers("data/all-members.json"); err != nil {
		return fmt.Errorf("could not load discord members json: %w", err)
	}

	if err := geo.LoadLand("data/land.json"); err != nil {
		return fmt.Errorf("could not load land geometry: %w", err)
	}

	mux.HandleFunc("GET /markers", func(w http.ResponseWriter, r *http.Request) {
		ticker := time.NewTicker(5000 * time.Millisecond)
		defer ticker.Stop()

		numMembers := 20

		sse := datastar.NewSSE(w, r, datastar.WithCompression())
		markers := generateRandomMemberMarkers(numMembers)

		b, err := json.Marshal(&templating.MarkerSignals{Markers: markers})
		if err != nil {
			return
		}

		sse.PatchSignals(b)

		for {
			select {
			case <-ctx.Done():
				slog.Debug("server closed connection")
				return

			case <-r.Context().Done():
				slog.Debug("client closed connection")
				return

			case <-ticker.C:
				markers := generateRandomMemberMarkers(numMembers)

				b, err := json.Marshal(templating.MarkerSignals{Markers: markers})
				if err != nil {
					return
				}

				sse.PatchSignals(b)
			}
		}
	})

	return nil
}

func generateRandomMemberMarkers(n int) []types.MapMarker {
	markers := make([]types.MapMarker, 0, n)

	idxs := make([]int, 0, n)
	for range n {
		idxs = append(idxs, rand.IntN(len(members)))
	}

	for _, idx := range idxs {
		markers = append(markers, generateRandomMarker(idx))
	}

	return markers
}

func generateRandomMarker(idx int) types.MapMarker {
	return types.MapMarker{
		User:        members[idx],
		Coordinates: geo.RandomLandPoint(),
	}
}

func loadDiscordMembers(path string) error {
	rawJsonBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	parsed, err := gabs.ParseJSON(rawJsonBytes)
	if err != nil {
		return err
	}

	members = make([]types.DiscordMember, 0, len(parsed.Children()))
	for _, child := range parsed.Children() {
		id, ok := child.Path("member.user.id").Data().(string)
		if !ok {
			continue
		}
		username, ok := child.Path("member.user.username").Data().(string)
		if !ok {
			continue
		}
		avatar, ok := child.Path("member.user.avatar").Data().(string)
		if !ok {
			continue
		}

		member := types.DiscordMember{
			Id:       id,
			Username: username,
			Avatar:   avatar,
		}

		members = append(members, member)
	}

	return nil
}
