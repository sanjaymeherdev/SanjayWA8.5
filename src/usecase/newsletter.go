package usecase

import (
	"context"

	domainNewsletter "whatsappbot/domains/newsletter"
	"whatsappbot/infrastructure/whatsapp"
	pkgError "whatsappbot/pkg/error"
	"whatsappbot/pkg/utils"
	"whatsappbot/validations"
)

type serviceNewsletter struct{}

func NewNewsletterService() domainNewsletter.INewsletterUsecase {
	return &serviceNewsletter{}
}

func (service serviceNewsletter) Unfollow(ctx context.Context, request domainNewsletter.UnfollowRequest) (err error) {
	if err = validations.ValidateUnfollowNewsletter(ctx, request); err != nil {
		return err
	}

	client := whatsapp.ClientFromContext(ctx)
	if client == nil {
		return pkgError.ErrWaCLI
	}

	JID, err := utils.ValidateJidWithLogin(client, request.NewsletterID)
	if err != nil {
		return err
	}

	return client.UnfollowNewsletter(ctx, JID)
}
