package validate

import (
	"context"
	"fmt"
	"strings"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/validation"
)

type validateOptions struct {
	AppID     string
	VersionID string
	Platform  string
	Strict    bool
	Output    string
	Pretty    bool
}

var clientFactory = shared.GetASCClient

func runValidate(ctx context.Context, opts validateOptions) error {
	client, err := clientFactory()
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	requestCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	versionResp, err := client.GetAppStoreVersion(requestCtx, opts.VersionID)
	if err != nil {
		return fmt.Errorf("validate: failed to fetch app store version: %w", err)
	}

	appResp, err := client.GetApp(requestCtx, opts.AppID)
	if err != nil {
		return fmt.Errorf("validate: failed to fetch app: %w", err)
	}

	versionLocsResp, err := client.GetAppStoreVersionLocalizations(requestCtx, opts.VersionID)
	if err != nil {
		return fmt.Errorf("validate: failed to fetch version localizations: %w", err)
	}

	appInfosResp, err := client.GetAppInfos(requestCtx, opts.AppID)
	if err != nil {
		return fmt.Errorf("validate: failed to fetch app info: %w", err)
	}

	appInfoID := shared.SelectBestAppInfoID(appInfosResp)
	if strings.TrimSpace(appInfoID) == "" {
		return fmt.Errorf("validate: failed to select app info for app")
	}

	appInfoLocsResp, err := client.GetAppInfoLocalizations(requestCtx, appInfoID)
	if err != nil {
		return fmt.Errorf("validate: failed to fetch app info localizations: %w", err)
	}

	var ageRatingDecl *validation.AgeRatingDeclaration
	ageRatingResp, err := client.GetAgeRatingDeclarationForAppStoreVersion(requestCtx, opts.VersionID)
	if err != nil {
		if !asc.IsNotFound(err) {
			return fmt.Errorf("validate: failed to fetch age rating declaration: %w", err)
		}
	} else {
		ageRatingDecl = mapAgeRatingDeclaration(ageRatingResp.Data.Attributes)
	}

	versionLocalizations := make([]validation.VersionLocalization, 0, len(versionLocsResp.Data))
	for _, loc := range versionLocsResp.Data {
		attrs := loc.Attributes
		versionLocalizations = append(versionLocalizations, validation.VersionLocalization{
			ID:              loc.ID,
			Locale:          attrs.Locale,
			Description:     attrs.Description,
			Keywords:        attrs.Keywords,
			WhatsNew:        attrs.WhatsNew,
			PromotionalText: attrs.PromotionalText,
			SupportURL:      attrs.SupportURL,
			MarketingURL:    attrs.MarketingURL,
		})
	}

	appInfoLocalizations := make([]validation.AppInfoLocalization, 0, len(appInfoLocsResp.Data))
	for _, loc := range appInfoLocsResp.Data {
		attrs := loc.Attributes
		appInfoLocalizations = append(appInfoLocalizations, validation.AppInfoLocalization{
			ID:       loc.ID,
			Locale:   attrs.Locale,
			Name:     attrs.Name,
			Subtitle: attrs.Subtitle,
		})
	}

	screenshotSets, err := fetchScreenshotSets(requestCtx, client, versionLocsResp.Data)
	if err != nil {
		return err
	}

	platform := opts.Platform
	if platform == "" {
		platform = string(versionResp.Data.Attributes.Platform)
	}

	report := validation.Validate(validation.Input{
		AppID:                opts.AppID,
		VersionID:            opts.VersionID,
		VersionString:        versionResp.Data.Attributes.VersionString,
		Platform:             platform,
		PrimaryLocale:        appResp.Data.Attributes.PrimaryLocale,
		VersionLocalizations: versionLocalizations,
		AppInfoLocalizations: appInfoLocalizations,
		ScreenshotSets:       screenshotSets,
		AgeRatingDeclaration: ageRatingDecl,
	}, opts.Strict)

	if err := shared.PrintOutput(&report, opts.Output, opts.Pretty); err != nil {
		return err
	}

	if report.Summary.Blocking > 0 {
		return shared.NewReportedError(fmt.Errorf("validate: found %d blocking issue(s)", report.Summary.Blocking))
	}

	return nil
}

func fetchScreenshotSets(ctx context.Context, client *asc.Client, localizations []asc.Resource[asc.AppStoreVersionLocalizationAttributes]) ([]validation.ScreenshotSet, error) {
	var sets []validation.ScreenshotSet
	for _, loc := range localizations {
		resp, err := client.GetAppStoreVersionLocalizationScreenshotSets(ctx, loc.ID)
		if err != nil {
			return nil, fmt.Errorf("validate: failed to fetch screenshot sets for %s: %w", loc.ID, err)
		}
		for _, set := range resp.Data {
			screenshotsResp, err := client.GetAppScreenshots(ctx, set.ID)
			if err != nil {
				return nil, fmt.Errorf("validate: failed to fetch screenshots for %s: %w", set.ID, err)
			}
			screenshots := make([]validation.Screenshot, 0, len(screenshotsResp.Data))
			for _, shot := range screenshotsResp.Data {
				width := 0
				height := 0
				if shot.Attributes.ImageAsset != nil {
					width = shot.Attributes.ImageAsset.Width
					height = shot.Attributes.ImageAsset.Height
				}
				screenshots = append(screenshots, validation.Screenshot{
					ID:       shot.ID,
					FileName: shot.Attributes.FileName,
					Width:    width,
					Height:   height,
				})
			}
			sets = append(sets, validation.ScreenshotSet{
				ID:             set.ID,
				DisplayType:    set.Attributes.ScreenshotDisplayType,
				Locale:         loc.Attributes.Locale,
				LocalizationID: loc.ID,
				Screenshots:    screenshots,
			})
		}
	}
	return sets, nil
}

func mapAgeRatingDeclaration(attrs asc.AgeRatingDeclarationAttributes) *validation.AgeRatingDeclaration {
	return &validation.AgeRatingDeclaration{
		Advertising:                                 attrs.Advertising,
		Gambling:                                    attrs.Gambling,
		HealthOrWellnessTopics:                      attrs.HealthOrWellnessTopics,
		LootBox:                                     attrs.LootBox,
		MessagingAndChat:                            attrs.MessagingAndChat,
		ParentalControls:                            attrs.ParentalControls,
		AgeAssurance:                                attrs.AgeAssurance,
		UnrestrictedWebAccess:                       attrs.UnrestrictedWebAccess,
		UserGeneratedContent:                        attrs.UserGeneratedContent,
		AlcoholTobaccoOrDrugUseOrReferences:         attrs.AlcoholTobaccoOrDrugUseOrReferences,
		Contests:                                    attrs.Contests,
		GamblingSimulated:                           attrs.GamblingSimulated,
		GunsOrOtherWeapons:                          attrs.GunsOrOtherWeapons,
		MedicalOrTreatmentInformation:               attrs.MedicalOrTreatmentInformation,
		ProfanityOrCrudeHumor:                       attrs.ProfanityOrCrudeHumor,
		SexualContentGraphicAndNudity:               attrs.SexualContentGraphicAndNudity,
		SexualContentOrNudity:                       attrs.SexualContentOrNudity,
		HorrorOrFearThemes:                          attrs.HorrorOrFearThemes,
		MatureOrSuggestiveThemes:                    attrs.MatureOrSuggestiveThemes,
		ViolenceCartoonOrFantasy:                    attrs.ViolenceCartoonOrFantasy,
		ViolenceRealistic:                           attrs.ViolenceRealistic,
		ViolenceRealisticProlongedGraphicOrSadistic: attrs.ViolenceRealisticProlongedGraphicOrSadistic,
		KidsAgeBand:                                 attrs.KidsAgeBand,
		AgeRatingOverride:                           attrs.AgeRatingOverride,
		AgeRatingOverrideV2:                         attrs.AgeRatingOverrideV2,
		KoreaAgeRatingOverride:                      attrs.KoreaAgeRatingOverride,
		DeveloperAgeRatingInfoURL:                   attrs.DeveloperAgeRatingInfoURL,
	}
}
