package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"sync"

	pkg "github.com/genmcp/grpc-demo/000_grpc_server/pkg"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// featureRecord matches the structure in the original HTTP server's main.go.
// Renamed from Feature to avoid conflict with generated protobuf message.
type featureRecord struct {
	ID          int
	Title       string
	Description string
	Details     string
	Upvotes     int
	Completed   bool
}

// server is used to implement FeatureServiceServer.
type server struct {
	pkg.UnimplementedFeatureServiceServer
	mu       sync.RWMutex
	features map[int]*featureRecord
	nextID   int
}

// newServer initializes the server with in-memory data.
func newServer() *server {
	s := &server{
		features: map[int]*featureRecord{
			1: {ID: 1, Title: "Dark Mode", Description: "Add dark theme support to the application", Details: "Implement a comprehensive dark mode that includes:\n\n- Automatic detection of system preference\n- Manual toggle in user settings\n- Dark variants for all UI components including buttons, forms, modals, and navigation\n- Proper contrast ratios for accessibility compliance\n- Smooth transitions between light and dark modes\n- Persistence of user preference across sessions\n- Support for custom accent colors in dark mode\n\nThis feature should integrate seamlessly with the existing design system and maintain consistency across all pages and components.", Upvotes: 142, Completed: false},
			2: {ID: 2, Title: "Mobile App", Description: "Native mobile application for iOS and Android", Details: "Develop native mobile applications for both iOS and Android platforms:\n\n**iOS App:**\n- Swift/SwiftUI implementation\n- iOS 14+ compatibility\n- App Store submission and compliance\n- Push notifications support\n- Offline functionality for core features\n\n**Android App:**\n- Kotlin implementation\n- Material Design 3 compliance\n- Android 8+ compatibility\n- Google Play Store submission\n- Background sync capabilities\n\n**Shared Features:**\n- Biometric authentication (Face ID, Touch ID, Fingerprint)\n- Deep linking support\n- Synchronized data across web and mobile\n- Performance optimization for battery life\n- Comprehensive testing on multiple devices", Upvotes: 98, Completed: false},
			3: {ID: 3, Title: "API Integration", Description: "Third-party API integrations for popular services", Details: "Build robust integrations with popular third-party services:\n\n**Communication APIs:**\n- Slack workspace integration\n- Microsoft Teams connector\n- Discord webhook support\n- Email service providers (SendGrid, Mailgun)\n\n**Productivity Tools:**\n- Google Workspace (Docs, Sheets, Calendar)\n- Microsoft Office 365\n- Trello and Asana project management\n- Notion database sync\n\n**Development Tools:**\n- GitHub repository integration\n- GitLab CI/CD webhooks\n- Jira issue tracking\n- Jenkins build notifications\n\n**Technical Requirements:**\n- OAuth 2.0 authentication flows\n- Rate limiting and retry mechanisms\n- Webhook validation and security\n- API key management interface\n- Real-time status monitoring\n- Comprehensive error handling and logging", Upvotes: 76, Completed: false},
			4: {ID: 4, Title: "Real-time Chat", Description: "Built-in real-time messaging system", Details: "Implement a comprehensive real-time messaging system:\n\n**Core Features:**\n- Instant messaging with WebSocket connections\n- Group chat rooms and private messaging\n- File sharing (images, documents, code snippets)\n- Message history and search functionality\n- Typing indicators and read receipts\n- Emoji reactions and custom emojis\n\n**Advanced Features:**\n- Message threading for organized discussions\n- Voice and video calling integration\n- Screen sharing capabilities\n- Message encryption for security\n- Customizable notifications\n- Message formatting (markdown support)\n\n**Technical Implementation:**\n- Scalable WebSocket infrastructure\n- Message persistence and backup\n- Real-time presence indicators\n- Mobile push notifications\n- Moderation tools and user management\n- Integration with existing user authentication", Upvotes: 54, Completed: false},
			5: {ID: 5, Title: "Advanced Analytics", Description: "Detailed analytics dashboard with custom metrics", Details: "Create a powerful analytics platform with comprehensive insights:\n\n**Dashboard Features:**\n- Customizable widget layout\n- Real-time data visualization\n- Interactive charts and graphs\n- Drill-down capabilities for detailed analysis\n- Export functionality (PDF, Excel, CSV)\n- Scheduled report generation\n\n**Metrics and KPIs:**\n- User engagement tracking\n- Performance monitoring\n- Conversion funnel analysis\n- A/B testing results\n- Custom event tracking\n- Revenue and growth metrics\n\n**Advanced Capabilities:**\n- Machine learning insights and predictions\n- Anomaly detection and alerts\n- Cohort analysis and user segmentation\n- Custom query builder\n- API for programmatic access\n- Integration with Google Analytics and other tools\n\n**Technical Features:**\n- High-performance data processing\n- Real-time data streaming\n- Historical data retention policies\n- GDPR compliance and data privacy controls", Upvotes: 31, Completed: false},
		},
		nextID: 6,
	}
	return s
}

// Helper to convert internal featureRecord struct to protobuf message
func (f *featureRecord) toProto() *pkg.Feature {
	return &pkg.Feature{
		Id:          int32(f.ID),
		Title:       f.Title,
		Description: f.Description,
		Details:     f.Details,
		Upvotes:     int32(f.Upvotes),
		Completed:   f.Completed,
	}
}

// Helper to convert internal featureRecord struct to protobuf summary message
func (f *featureRecord) toSummaryProto() *pkg.FeatureSummary {
	return &pkg.FeatureSummary{
		Id:        int32(f.ID),
		Title:     f.Title,
		Upvotes:   int32(f.Upvotes),
		Completed: f.Completed,
	}
}

// ListFeatures returns a list of all features.
func (s *server) ListFeatures(ctx context.Context, in *emptypb.Empty) (*pkg.ListFeaturesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	featureList := make([]*featureRecord, 0, len(s.features))
	for _, feature := range s.features {
		featureList = append(featureList, feature)
	}

	sort.Slice(featureList, func(i, j int) bool {
		return featureList[i].Upvotes > featureList[j].Upvotes
	})

	summaries := make([]*pkg.FeatureSummary, len(featureList))
	for i, feature := range featureList {
		summaries[i] = feature.toSummaryProto()
	}

	return &pkg.ListFeaturesResponse{Summaries: summaries}, nil
}

// GetTopFeature returns the feature with the most upvotes.
func (s *server) GetTopFeature(ctx context.Context, in *emptypb.Empty) (*pkg.FeatureSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var topFeature *featureRecord
	maxVotes := -1

	for _, feature := range s.features {
		if feature.Upvotes > maxVotes {
			maxVotes = feature.Upvotes
			topFeature = feature
		}
	}

	if topFeature == nil {
		return nil, status.Error(codes.NotFound, "no features found")
	}

	return topFeature.toSummaryProto(), nil
}

// GetFeature returns details for a specific feature.
func (s *server) GetFeature(ctx context.Context, in *pkg.GetFeatureRequest) (*pkg.Feature, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	feature, exists := s.features[int(in.GetId())]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "feature with ID %d not found", in.GetId())
	}

	return feature.toProto(), nil
}

// AddFeature creates a new feature request.
func (s *server) AddFeature(ctx context.Context, in *pkg.AddFeatureRequest) (*pkg.Feature, error) {
	if in.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	newFeature := &featureRecord{
		ID:          s.nextID,
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Details:     in.GetDetails(),
		Upvotes:     0,
		Completed:   false,
	}
	s.features[s.nextID] = newFeature
	s.nextID++

	return newFeature.toProto(), nil
}

// VoteFeature increments the upvote count for a feature.
func (s *server) VoteFeature(ctx context.Context, in *pkg.VoteFeatureRequest) (*pkg.Feature, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	feature, exists := s.features[int(in.GetId())]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "feature with ID %d not found", in.GetId())
	}

	feature.Upvotes++
	return feature.toProto(), nil
}

// CompleteFeature marks a feature as completed.
func (s *server) CompleteFeature(ctx context.Context, in *pkg.CompleteFeatureRequest) (*pkg.Feature, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	feature, exists := s.features[int(in.GetId())]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "feature with ID %d not found", in.GetId())
	}

	feature.Completed = true
	return feature.toProto(), nil
}

// DeleteFeature removes a feature.
func (s *server) DeleteFeature(ctx context.Context, in *pkg.DeleteFeatureRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.features[int(in.GetId())]; !exists {
		return nil, status.Errorf(codes.NotFound, "feature with ID %d not found", in.GetId())
	}

	delete(s.features, int(in.GetId()))
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pkg.RegisterFeatureServiceServer(grpcServer, newServer())

	// Enable reflection for gRPC server
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
