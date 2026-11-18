package ashby

import "time"

type AshbyJob struct {
	ID                                     string    `json:"id"`
	Title                                  string    `json:"title"`
	Department                             string    `json:"department"`
	Team                                   string    `json:"team"`
	EmploymentType                         string    `json:"employmentType"`
	Location                               string    `json:"location"`
	ShouldDisplayCompensationOnJobPostings bool      `json:"shouldDisplayCompensationOnJobPostings"`
	SecondaryLocations                     []any     `json:"secondaryLocations"`
	PublishedAt                            time.Time `json:"publishedAt"`
	IsListed                               bool      `json:"isListed"`
	IsRemote                               *bool     `json:"isRemote"`
	WorkplaceType                          string    `json:"workplaceType"`
	Address                                struct {
		PostalAddress struct {
			AddressRegion   string `json:"addressRegion"`
			StreetAddress   string `json:"streetAddress"`
			AddressCountry  string `json:"addressCountry"`
			AddressLocality string `json:"addressLocality"`
		} `json:"postalAddress"`
	} `json:"address"`
	JobURL           string `json:"jobUrl"`
	ApplyURL         string `json:"applyUrl"`
	DescriptionHTML  string `json:"descriptionHtml"`
	DescriptionPlain string `json:"descriptionPlain"`
	Compensation     struct {
		CompensationTierSummary             any   `json:"compensationTierSummary"`
		ScrapeableCompensationSalarySummary any   `json:"scrapeableCompensationSalarySummary"`
		CompensationTiers                   []any `json:"compensationTiers"`
		SummaryComponents                   []any `json:"summaryComponents"`
	} `json:"compensation"`
}

type AshbyJobsList struct {
	Jobs       []AshbyJob `json:"jobs"`
	APIVersion string     `json:"apiVersion"`
}

type ghJobsPreData struct {
	Payload struct {
		CodeViewBlobRoute struct {
			Csv        any `json:"csv"`
			CsvError   any `json:"csvError"`
			HeaderInfo struct {
				Toc any `json:"toc"`
			} `json:"headerInfo"`
			IssueTemplate      any  `json:"issueTemplate"`
			DiscussionTemplate any  `json:"discussionTemplate"`
			RichText           any  `json:"richText"`
			RichTextTruncated  bool `json:"richTextTruncated"`
			RenderedFileInfo   any  `json:"renderedFileInfo"`
			Symbols            struct {
				TimedOut    bool  `json:"timed_out"`
				NotAnalyzed bool  `json:"not_analyzed"`
				Symbols     []any `json:"symbols"`
			} `json:"symbols"`
		} `json:"codeViewBlobRoute"`
		CodeViewLayoutRoute struct {
			Repo struct {
				ID                 int       `json:"id"`
				DefaultBranch      string    `json:"defaultBranch"`
				Name               string    `json:"name"`
				OwnerLogin         string    `json:"ownerLogin"`
				CurrentUserCanPush bool      `json:"currentUserCanPush"`
				IsFork             bool      `json:"isFork"`
				IsEmpty            bool      `json:"isEmpty"`
				CreatedAt          time.Time `json:"createdAt"`
				OwnerAvatar        string    `json:"ownerAvatar"`
				Public             bool      `json:"public"`
				Private            bool      `json:"private"`
				IsOrgOwned         bool      `json:"isOrgOwned"`
				IsArchived         bool      `json:"isArchived"`
			} `json:"repo"`
			CurrentUser         any    `json:"currentUser"`
			UploadToken         string `json:"uploadToken"`
			AllShortcutsEnabled bool   `json:"allShortcutsEnabled"`
			TreeExpanded        bool   `json:"treeExpanded"`
			Path                string `json:"path"`
			SymbolsExpanded     bool   `json:"symbolsExpanded"`
			RefInfo             struct {
				Name         string `json:"name"`
				ListCacheKey string `json:"listCacheKey"`
				CanEdit      bool   `json:"canEdit"`
				CurrentOid   string `json:"currentOid"`
			} `json:"refInfo"`
			HelpURL      string `json:"helpUrl"`
			GithubDevURL any    `json:"githubDevUrl"`
		} `json:"codeViewLayoutRoute"`
		CodeViewFileTreeLayoutRoute struct {
			FileTree struct {
				NAMING_FAILED struct {
					Items []struct {
						Name        string `json:"name"`
						Path        string `json:"path"`
						ContentType string `json:"contentType"`
					} `json:"items"`
					TotalCount int `json:"totalCount"`
				} `json:""`
			} `json:"fileTree"`
			FileTreeProcessingTime float64 `json:"fileTreeProcessingTime"`
			FoldersToFetch         []any   `json:"foldersToFetch"`
			RefName                string  `json:"refName"`
		} `json:"codeViewFileTreeLayoutRoute"`
		CodeViewBlobLayoutRoute struct {
			CodeLineWrapEnabled bool `json:"codeLineWrapEnabled"`
			RefInfo             struct {
				Name                   string `json:"name"`
				ListCacheKey           string `json:"listCacheKey"`
				CanEdit                bool   `json:"canEdit"`
				RefType                string `json:"refType"`
				CurrentOid             string `json:"currentOid"`
				CanEditOnDefaultBranch bool   `json:"canEditOnDefaultBranch"`
				FileExistsOnDefault    bool   `json:"fileExistsOnDefault"`
			} `json:"refInfo"`
			Path string `json:"path"`
			Blob struct {
				CopilotSWEAgentEnabled bool `json:"copilotSWEAgentEnabled"`
				DependabotInfo         struct {
					ShowConfigurationBanner        bool   `json:"showConfigurationBanner"`
					ConfigFilePath                 any    `json:"configFilePath"`
					NetworkDependabotPath          string `json:"networkDependabotPath"`
					DismissConfigurationNoticePath string `json:"dismissConfigurationNoticePath"`
					ConfigurationNoticeDismissed   any    `json:"configurationNoticeDismissed"`
				} `json:"dependabotInfo"`
				DisplayName string `json:"displayName"`
				DisplayURL  string `json:"displayUrl"`
				HeaderInfo  struct {
					BlobSize         string `json:"blobSize"`
					DeleteTooltip    string `json:"deleteTooltip"`
					EditTooltip      string `json:"editTooltip"`
					GhDesktopPath    string `json:"ghDesktopPath"`
					IsGitLfs         bool   `json:"isGitLfs"`
					OnBranch         bool   `json:"onBranch"`
					ShortPath        string `json:"shortPath"`
					SiteNavLoginPath string `json:"siteNavLoginPath"`
					IsCSV            bool   `json:"isCSV"`
					IsRichtext       bool   `json:"isRichtext"`
					LineInfo         struct {
						TruncatedLoc  string `json:"truncatedLoc"`
						TruncatedSloc string `json:"truncatedSloc"`
					} `json:"lineInfo"`
					Mode string `json:"mode"`
				} `json:"headerInfo"`
				Image                      bool   `json:"image"`
				IsCodeownersFile           any    `json:"isCodeownersFile"`
				IsPlain                    bool   `json:"isPlain"`
				IsValidLegacyIssueTemplate bool   `json:"isValidLegacyIssueTemplate"`
				IsIssueTemplate            bool   `json:"isIssueTemplate"`
				IsDiscussionTemplate       bool   `json:"isDiscussionTemplate"`
				Language                   string `json:"language"`
				LanguageID                 int    `json:"languageID"`
				Large                      bool   `json:"large"`
				PlanSupportInfo            struct {
					RepoIsFork                     any    `json:"repoIsFork"`
					RepoOwnedByCurrentUser         any    `json:"repoOwnedByCurrentUser"`
					RequestFullPath                string `json:"requestFullPath"`
					ShowFreeOrgGatedFeatureMessage any    `json:"showFreeOrgGatedFeatureMessage"`
					ShowPlanSupportBanner          any    `json:"showPlanSupportBanner"`
					UpgradeDataAttributes          any    `json:"upgradeDataAttributes"`
					UpgradePath                    any    `json:"upgradePath"`
				} `json:"planSupportInfo"`
				PublishBannersInfo struct {
					DismissActionNoticePath string `json:"dismissActionNoticePath"`
					ReleasePath             string `json:"releasePath"`
					ShowPublishActionBanner bool   `json:"showPublishActionBanner"`
				} `json:"publishBannersInfo"`
				RawBlobURL       string `json:"rawBlobUrl"`
				RenderImageOrRaw bool   `json:"renderImageOrRaw"`
				ShortPath        any    `json:"shortPath"`
				SymbolsEnabled   bool   `json:"symbolsEnabled"`
				TabSize          int    `json:"tabSize"`
				TopBannersInfo   struct {
					OverridingGlobalFundingFile bool   `json:"overridingGlobalFundingFile"`
					GlobalPreferredFundingPath  any    `json:"globalPreferredFundingPath"`
					ShowInvalidCitationWarning  bool   `json:"showInvalidCitationWarning"`
					CitationHelpURL             string `json:"citationHelpUrl"`
					ActionsOnboardingTip        any    `json:"actionsOnboardingTip"`
				} `json:"topBannersInfo"`
				Truncated           bool `json:"truncated"`
				Viewable            bool `json:"viewable"`
				WorkflowRedirectURL any  `json:"workflowRedirectUrl"`
			} `json:"blob"`
			CopilotInfo                  any  `json:"copilotInfo"`
			CopilotAccessAllowed         bool `json:"copilotAccessAllowed"`
			CopilotSpacesEnabled         bool `json:"copilotSpacesEnabled"`
			ModelsAccessAllowed          bool `json:"modelsAccessAllowed"`
			ModelsRepoIntegrationEnabled bool `json:"modelsRepoIntegrationEnabled"`
			IsMarketplaceEnabled         bool `json:"isMarketplaceEnabled"`
		} `json:"codeViewBlobLayoutRoute"`
		CodeViewBlobLayoutRouteStyledBlob struct {
			RawLines          []string `json:"rawLines"`
			StylingDirectives any      `json:"stylingDirectives"`
			ColorizedLines    []string `json:"colorizedLines"`
		} `json:"codeViewBlobLayoutRoute.StyledBlob"`
	} `json:"payload"`
	Title      string `json:"title"`
	AppPayload struct {
	} `json:"appPayload"`
	Meta struct {
		Title string `json:"title"`
	} `json:"meta"`
}

type ghRawdata struct {
	GeneratedAt   time.Time `json:"generated_at"`
	Ats           string    `json:"ats"`
	LiveCompanies int       `json:"live_companies"`
	Companies     []struct {
		Slug            string `json:"slug"`
		Company         string `json:"company"`
		Status          string `json:"status"`
		Jobs            int    `json:"jobs"`
		LastValidatedAt string `json:"last_validated_at"`
	} `json:"companies"`
}
