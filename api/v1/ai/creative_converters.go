package ai

import (
	pb "Qingyu_backend/pkg/grpc/pb"
)

// ============ Proto → Model 转换 ============

// convertProtoOutlineToModel 将proto Outline转换为Model
func convertProtoOutlineToModel(proto *pb.OutlineData) *OutlineData {
	if proto == nil {
		return nil
	}

	chapters := make([]ChapterData, len(proto.Chapters))
	for i, ch := range proto.Chapters {
		chapters[i] = ChapterData{
			ChapterID:          ch.ChapterId,
			Title:              ch.Title,
			Summary:            ch.Summary,
			KeyEvents:          ch.KeyEvents,
			CharactersInvolved: ch.CharactersInvolved,
			ConflictType:       ch.ConflictType,
			EmotionalTone:      ch.EmotionalTone,
			EstimatedWordCount: ch.EstimatedWordCount,
			ChapterGoal:        ch.ChapterGoal,
			Cliffhanger:        ch.Cliffhanger,
		}
	}

	var storyArc *StoryArc
	if proto.StoryArc != nil {
		storyArc = &StoryArc{
			Setup:         proto.StoryArc.Setup,
			RisingAction:  proto.StoryArc.RisingAction,
			Climax:        proto.StoryArc.Climax,
			FallingAction: proto.StoryArc.FallingAction,
			Resolution:    proto.StoryArc.Resolution,
		}
	}

	return &OutlineData{
		Title:               proto.Title,
		Genre:               proto.Genre,
		CoreTheme:           proto.CoreTheme,
		TargetAudience:      proto.TargetAudience,
		EstimatedTotalWords: proto.EstimatedTotalWords,
		Chapters:            chapters,
		StoryArc:            storyArc,
	}
}

// convertProtoCharactersToModel 将proto Characters转换为Model
func convertProtoCharactersToModel(proto *pb.CharactersData) *CharactersData {
	if proto == nil {
		return nil
	}

	characters := make([]CharacterData, len(proto.Characters))
	for i, ch := range proto.Characters {
		var personality *PersonalityData
		if ch.Personality != nil {
			personality = &PersonalityData{
				Traits:     ch.Personality.Traits,
				Strengths:  ch.Personality.Strengths,
				Weaknesses: ch.Personality.Weaknesses,
				CoreValues: ch.Personality.CoreValues,
				Fears:      ch.Personality.Fears,
			}
		}

		var background *BackgroundData
		if ch.Background != nil {
			background = &BackgroundData{
				Summary:        ch.Background.Summary,
				Family:         ch.Background.Family,
				Education:      ch.Background.Education,
				KeyExperiences: ch.Background.KeyExperiences,
			}
		}

		relationships := make([]RelationshipData, len(ch.Relationships))
		for j, rel := range ch.Relationships {
			relationships[j] = RelationshipData{
				Character:    rel.Character,
				RelationType: rel.RelationType,
				Description:  rel.Description,
				Dynamics:     rel.Dynamics,
			}
		}

		var devArc *DevelopmentArc
		if ch.DevelopmentArc != nil {
			devArc = &DevelopmentArc{
				StartingPoint: ch.DevelopmentArc.StartingPoint,
				TurningPoints: ch.DevelopmentArc.TurningPoints,
				EndingPoint:   ch.DevelopmentArc.EndingPoint,
				GrowthTheme:   ch.DevelopmentArc.GrowthTheme,
			}
		}

		characters[i] = CharacterData{
			CharacterID:      ch.CharacterId,
			Name:             ch.Name,
			RoleType:         ch.RoleType,
			Importance:       ch.Importance,
			Age:              ch.Age,
			Gender:           ch.Gender,
			Appearance:       ch.Appearance,
			Personality:      personality,
			Background:       background,
			Motivation:       ch.Motivation,
			Relationships:    relationships,
			DevelopmentArc:   devArc,
			RoleInStory:      ch.RoleInStory,
			FirstAppearance:  ch.FirstAppearance,
			ChaptersInvolved: ch.ChaptersInvolved,
		}
	}

	var network *RelationshipNetwork
	if proto.RelationshipNetwork != nil {
		alliances := make([][]string, len(proto.RelationshipNetwork.Alliances))
		for i, alliance := range proto.RelationshipNetwork.Alliances {
			alliances[i] = alliance.Members
		}

		conflicts := make([][]string, len(proto.RelationshipNetwork.Conflicts))
		for i, conflict := range proto.RelationshipNetwork.Conflicts {
			conflicts[i] = conflict.Parties
		}

		mentorships := make([]MentorshipData, len(proto.RelationshipNetwork.Mentorships))
		for i, m := range proto.RelationshipNetwork.Mentorships {
			mentorships[i] = MentorshipData{
				Mentor:  m.Mentor,
				Student: m.Student,
			}
		}

		network = &RelationshipNetwork{
			Alliances:   alliances,
			Conflicts:   conflicts,
			Mentorships: mentorships,
		}
	}

	return &CharactersData{
		Characters:          characters,
		RelationshipNetwork: network,
	}
}

// convertProtoPlotToModel 将proto Plot转换为Model
func convertProtoPlotToModel(proto *pb.PlotData) *PlotData {
	if proto == nil {
		return nil
	}

	events := make([]TimelineEventData, len(proto.TimelineEvents))
	for i, ev := range proto.TimelineEvents {
		var impact *EventImpact
		if ev.Impact != nil {
			impact = &EventImpact{
				OnPlot:          ev.Impact.OnPlot,
				OnCharacters:    ev.Impact.OnCharacters,
				EmotionalImpact: ev.Impact.EmotionalImpact,
			}
		}

		events[i] = TimelineEventData{
			EventID:      ev.EventId,
			Timestamp:    ev.Timestamp,
			Location:     ev.Location,
			Title:        ev.Title,
			Description:  ev.Description,
			Participants: ev.Participants,
			EventType:    ev.EventType,
			Impact:       impact,
			Causes:       ev.Causes,
			Consequences: ev.Consequences,
			ChapterID:    ev.ChapterId,
		}
	}

	threads := make([]PlotThreadData, len(proto.PlotThreads))
	for i, th := range proto.PlotThreads {
		threads[i] = PlotThreadData{
			ThreadID:           th.ThreadId,
			Title:              th.Title,
			Description:        th.Description,
			Type:               th.Type,
			Events:             th.Events,
			StartingEvent:      th.StartingEvent,
			ClimaxEvent:        th.ClimaxEvent,
			ResolutionEvent:    th.ResolutionEvent,
			CharactersInvolved: th.CharactersInvolved,
		}
	}

	conflicts := make([]ConflictData, len(proto.Conflicts))
	for i, cf := range proto.Conflicts {
		conflicts[i] = ConflictData{
			ConflictID:       cf.ConflictId,
			Type:             cf.Type,
			Parties:          cf.Parties,
			Description:      cf.Description,
			EscalationEvents: cf.EscalationEvents,
			ResolutionEvent:  cf.ResolutionEvent,
		}
	}

	var keyPoints *KeyPlotPoints
	if proto.KeyPlotPoints != nil {
		keyPoints = &KeyPlotPoints{
			IncitingIncident: proto.KeyPlotPoints.IncitingIncident,
			PlotPoint1:       proto.KeyPlotPoints.PlotPoint_1,
			Midpoint:         proto.KeyPlotPoints.Midpoint,
			PlotPoint2:       proto.KeyPlotPoints.PlotPoint_2,
			Climax:           proto.KeyPlotPoints.Climax,
			Resolution:       proto.KeyPlotPoints.Resolution,
		}
	}

	return &PlotData{
		TimelineEvents: events,
		PlotThreads:    threads,
		Conflicts:      conflicts,
		KeyPlotPoints:  keyPoints,
	}
}

// convertProtoDiagnosticReportToModel 将proto DiagnosticReport转换为Model
func convertProtoDiagnosticReportToModel(proto *pb.DiagnosticReportData) *DiagnosticReportData {
	if proto == nil {
		return nil
	}

	issues := make([]DiagnosticIssue, len(proto.Issues))
	for i, issue := range proto.Issues {
		issues[i] = DiagnosticIssue{
			ID:               issue.Id,
			Severity:         issue.Severity,
			Category:         issue.Category,
			SubCategory:      issue.SubCategory,
			Title:            issue.Title,
			Description:      issue.Description,
			RootCause:        issue.RootCause,
			AffectedEntities: issue.AffectedEntities,
			Impact:           issue.Impact,
		}
	}

	return &DiagnosticReportData{
		Passed:             proto.Passed,
		QualityScore:       proto.QualityScore,
		Issues:             issues,
		CorrectionStrategy: proto.CorrectionStrategy,
		AffectedAgents:     proto.AffectedAgents,
		ReasoningChain:     proto.ReasoningChain,
	}
}

// ============ Model → Proto 转换 ============

// convertModelOutlineToProto 将Model Outline转换为Proto
func convertModelOutlineToProto(model *OutlineData) *pb.OutlineData {
	if model == nil {
		return nil
	}

	chapters := make([]*pb.ChapterData, len(model.Chapters))
	for i, ch := range model.Chapters {
		chapters[i] = &pb.ChapterData{
			ChapterId:          ch.ChapterID,
			Title:              ch.Title,
			Summary:            ch.Summary,
			KeyEvents:          ch.KeyEvents,
			CharactersInvolved: ch.CharactersInvolved,
			ConflictType:       ch.ConflictType,
			EmotionalTone:      ch.EmotionalTone,
			EstimatedWordCount: ch.EstimatedWordCount,
			ChapterGoal:        ch.ChapterGoal,
			Cliffhanger:        ch.Cliffhanger,
		}
	}

	var storyArc *pb.StoryArc
	if model.StoryArc != nil {
		storyArc = &pb.StoryArc{
			Setup:         model.StoryArc.Setup,
			RisingAction:  model.StoryArc.RisingAction,
			Climax:        model.StoryArc.Climax,
			FallingAction: model.StoryArc.FallingAction,
			Resolution:    model.StoryArc.Resolution,
		}
	}

	return &pb.OutlineData{
		Title:               model.Title,
		Genre:               model.Genre,
		CoreTheme:           model.CoreTheme,
		TargetAudience:      model.TargetAudience,
		EstimatedTotalWords: model.EstimatedTotalWords,
		Chapters:            chapters,
		StoryArc:            storyArc,
	}
}

// convertModelCharactersToProto 将Model Characters转换为Proto
func convertModelCharactersToProto(model *CharactersData) *pb.CharactersData {
	if model == nil {
		return nil
	}

	characters := make([]*pb.CharacterData, len(model.Characters))
	for i, ch := range model.Characters {
		var personality *pb.PersonalityData
		if ch.Personality != nil {
			personality = &pb.PersonalityData{
				Traits:     ch.Personality.Traits,
				Strengths:  ch.Personality.Strengths,
				Weaknesses: ch.Personality.Weaknesses,
				CoreValues: ch.Personality.CoreValues,
				Fears:      ch.Personality.Fears,
			}
		}

		var background *pb.BackgroundData
		if ch.Background != nil {
			background = &pb.BackgroundData{
				Summary:        ch.Background.Summary,
				Family:         ch.Background.Family,
				Education:      ch.Background.Education,
				KeyExperiences: ch.Background.KeyExperiences,
			}
		}

		relationships := make([]*pb.RelationshipData, len(ch.Relationships))
		for j, rel := range ch.Relationships {
			relationships[j] = &pb.RelationshipData{
				Character:    rel.Character,
				RelationType: rel.RelationType,
				Description:  rel.Description,
				Dynamics:     rel.Dynamics,
			}
		}

		var devArc *pb.DevelopmentArc
		if ch.DevelopmentArc != nil {
			devArc = &pb.DevelopmentArc{
				StartingPoint: ch.DevelopmentArc.StartingPoint,
				TurningPoints: ch.DevelopmentArc.TurningPoints,
				EndingPoint:   ch.DevelopmentArc.EndingPoint,
				GrowthTheme:   ch.DevelopmentArc.GrowthTheme,
			}
		}

		characters[i] = &pb.CharacterData{
			CharacterId:      ch.CharacterID,
			Name:             ch.Name,
			RoleType:         ch.RoleType,
			Importance:       ch.Importance,
			Age:              ch.Age,
			Gender:           ch.Gender,
			Appearance:       ch.Appearance,
			Personality:      personality,
			Background:       background,
			Motivation:       ch.Motivation,
			Relationships:    relationships,
			DevelopmentArc:   devArc,
			RoleInStory:      ch.RoleInStory,
			FirstAppearance:  ch.FirstAppearance,
			ChaptersInvolved: ch.ChaptersInvolved,
		}
	}

	var network *pb.RelationshipNetwork
	if model.RelationshipNetwork != nil {
		alliances := make([]*pb.Alliance, len(model.RelationshipNetwork.Alliances))
		for i, alliance := range model.RelationshipNetwork.Alliances {
			alliances[i] = &pb.Alliance{Members: alliance}
		}

		conflicts := make([]*pb.Conflict, len(model.RelationshipNetwork.Conflicts))
		for i, conflict := range model.RelationshipNetwork.Conflicts {
			conflicts[i] = &pb.Conflict{Parties: conflict}
		}

		mentorships := make([]*pb.Mentorship, len(model.RelationshipNetwork.Mentorships))
		for i, m := range model.RelationshipNetwork.Mentorships {
			mentorships[i] = &pb.Mentorship{
				Mentor:  m.Mentor,
				Student: m.Student,
			}
		}

		network = &pb.RelationshipNetwork{
			Alliances:   alliances,
			Conflicts:   conflicts,
			Mentorships: mentorships,
		}
	}

	return &pb.CharactersData{
		Characters:          characters,
		RelationshipNetwork: network,
	}
}
