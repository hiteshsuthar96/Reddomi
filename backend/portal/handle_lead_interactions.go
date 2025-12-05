package portal

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"fmt"
	"github.com/shank318/doota/datastore"
	"github.com/shank318/doota/models"
	pbcore "github.com/shank318/doota/pb/doota/core/v1"
	pbportal "github.com/shank318/doota/pb/doota/portal/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p *Portal) UpdateLeadInteraction(ctx context.Context, c *connect.Request[pbportal.UpdateLeadInteractionRequest]) (*connect.Response[emptypb.Empty], error) {
	actor, err := p.gethAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	// alteast one of comment or status should be provided
	if c.Msg.Comment == "" && c.Msg.Dm == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("at least one of comment or dm to edit should be provided"))
	}

	project, err := p.getProject(ctx, c.Header(), actor.OrganizationID)
	if err != nil {
		return nil, err
	}

	interaction, err := p.db.GetLeadInteractionByID(ctx, c.Msg.InteractionId)
	if err != nil {
		return nil, err
	}

	lead, err := p.db.GetLeadByID(ctx, project.ID, interaction.LeadID)
	if err != nil && !errors.Is(err, datastore.NotFound) {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to fetch lead: %w", err))
	}

	if lead == nil {
		return connect.NewResponse(&emptypb.Empty{}), nil
	}

	if c.Msg.Comment != "" {
		lead.LeadMetadata.SuggestedComment = c.Msg.Comment
	}
	if c.Msg.Dm != "" {
		lead.LeadMetadata.SuggestedDM = c.Msg.Dm
	}

	if err := p.db.UpdateLeadStatus(ctx, lead); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to update lead status: %w", err))
	}

	return connect.NewResponse(&emptypb.Empty{}), nil

}

func (p *Portal) UpdateLeadInteractionStatus(ctx context.Context, c *connect.Request[pbportal.UpdateLeadInteractionRequest]) (*connect.Response[emptypb.Empty], error) {
	actor, err := p.gethAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	status := c.Msg.Status

	// Only allow specific status updates
	if status != pbcore.LeadInteractionStatus_LEAD_INTERACTION_STATUS_REMOVED &&
		status != pbcore.LeadInteractionStatus_LEAD_INTERACTION_STATUS_CREATED {
		return connect.NewResponse(&emptypb.Empty{}), nil
	}

	interaction, err := p.db.GetLeadInteractionByID(ctx, c.Msg.InteractionId)
	if err != nil {
		return nil, err
	}

	// Handle CREATED status only if the current status is FAILED
	if status == pbcore.LeadInteractionStatus_LEAD_INTERACTION_STATUS_CREATED {
		if interaction.Status != models.LeadInteractionStatusFAILED {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("interaction status is not failed, cannot retry"))
		}

		interaction.Status = status.ToModel()
		if err := p.db.UpdateLeadInteraction(ctx, interaction); err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update interaction to CREATED: %w", err))
		}

		return connect.NewResponse(&emptypb.Empty{}), nil
	}

	// At this point, status must be REMOVED
	project, err := p.getProject(ctx, c.Header(), actor.OrganizationID)
	if err != nil {
		return nil, err
	}

	lead, err := p.db.GetLeadByID(ctx, project.ID, interaction.LeadID)
	if err != nil && !errors.Is(err, datastore.NotFound) {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to fetch lead: %w", err))
	}

	if lead == nil {
		return connect.NewResponse(&emptypb.Empty{}), nil
	}

	interaction.Status = status.ToModel()
	interaction.Reason = "Skipped, as user marked it as not relevant"

	// Clear scheduled timestamps
	lead.LeadMetadata.CommentScheduledAt = nil
	lead.LeadMetadata.DMScheduledAt = nil

	if err := p.db.UpdateLeadInteraction(ctx, interaction); err != nil {
		return nil, err
	}

	if err := p.db.UpdateLeadStatus(ctx, lead); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to update lead status: %w", err))
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (p *Portal) GetLeadInteractions(ctx context.Context, c *connect.Request[pbportal.GetLeadInteractionsRequest]) (*connect.Response[pbportal.GetLeadInteractionsResponse], error) {
	actor, err := p.gethAuthContext(ctx)
	if err != nil {
		return nil, err
	}
	project, err := p.getProject(ctx, c.Header(), actor.OrganizationID)
	if err != nil {
		return nil, err
	}

	interactions, err := p.db.GetAugmentedLeadInteractions(ctx, project.ID, c.Msg.GetDateRange())
	if err != nil {
		return nil, err
	}

	leadProtos := make([]*pbcore.LeadInteraction, 0, len(interactions))
	for _, interaction := range interactions {
		leadProtos = append(leadProtos, new(pbcore.LeadInteraction).FromModel(interaction))
	}

	return connect.NewResponse(&pbportal.GetLeadInteractionsResponse{Interactions: leadProtos}), nil
}
