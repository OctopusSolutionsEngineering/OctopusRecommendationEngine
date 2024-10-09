package performance

import (
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLongTaskQueue(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" || r.URL.Path == "/api/Spaces-1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
  "Application": "Octopus Deploy",
  "Version": "2023.1.9626",
  "ApiVersion": "3.0.0",
  "InstallationId": "79af4d4a-dc7f-479c-b0c5-8bf30cb0a827",
  "HasLongTermSupport": true,
  "Links": {
    "Self": "/api",
    "Accounts": "/api/Spaces-1/accounts{/id}{?skip,take,ids,partialName,accountType}",
    "ActionTemplateLogo": "/api/Spaces-1/actiontemplates/{typeOrId}/logo{?cb}",
    "ActionTemplates": "/api/Spaces-1/actiontemplates{/id}{?skip,take,ids,partialName}",
    "ActionTemplatesCategories": "/api/Spaces-1/actiontemplates/categories",
    "ActionTemplatesSearch": "/api/Spaces-1/actiontemplates/search{?type}",
    "ActionTemplateVersionedLogo": "/api/Spaces-1/actiontemplates/{typeOrId}/versions/{version}/logo{?cb}",
    "ArchivedEventFiles": "/api/events/archives{?skip,take}",
    "Artifacts": "/api/Spaces-1/artifacts{/id}{?skip,take,regarding,ids,partialName,order}",
    "AuditStreamConfiguration": "/api/audit-stream",
    "Authenticate_Azure AD": "/users/authenticate/AzureAD{?returnUrl}",
    "Authenticate_Octopus ID": "/users/authenticate/OctopusID{?returnUrl}",
    "Authentication": "/api/authentication",
    "AzureDevOpsConnectivityCheck": "/api/azuredevopsissuetracker/connectivitycheck",
    "AzureEnvironments": "/api/accounts/azureenvironments",
    "BuildInformation": "/api/Spaces-1/build-information{/id}{?packageId,filter,latest,skip,take,overwriteMode}",
    "BuildInformationBulk": "/api/Spaces-1/build-information/bulk{?ids}",
    "BuiltInFeedStats": "/api/feeds/stats",
    "CertificateConfiguration": "/api/configuration/certificates{/id}{?skip,take}",
    "Certificates": "/api/Spaces-1/certificates{/id}{?skip,take,search,archived,tenant,firstResult,orderBy,ids,partialName}",
    "Channels": "/api/Spaces-1/channels{/id}{?skip,take,ids,partialName}",
    "CloudTemplate": "/api/cloudtemplate/{id}/metadata{?packageId,feedId}",
    "CommunityActionTemplates": "/api/communityactiontemplates{/id}{?skip,take,ids}",
    "Configuration": "/api/configuration{/id}",
    "CurrentLicense": "/api/licenses/licenses-current",
    "CurrentLicenseStatus": "/api/licenses/licenses-current-status",
    "CurrentUser": "/api/users/me",
    "Dashboard": "/api/Spaces-1/dashboard{?projectId,releaseId,selectedTenants,selectedTags,showAll,highestLatestVersionPerProjectAndEnvironment}",
    "DashboardConfiguration": "/api/Spaces-1/dashboardconfiguration",
    "DashboardDynamic": "/api/Spaces-1/dashboard/dynamic{?projects,environments,includePrevious}",
    "DeploymentProcesses": "/api/Spaces-1/deploymentprocesses{/id}{?skip,take,ids}",
    "Deployments": "/api/Spaces-1/deployments{/id}{?skip,take,ids,projects,environments,tenants,channels,taskState,partialName}",
    "DiscoverMachine": "/api/Spaces-1/machines/discover{?host,port,type,proxyId}",
    "DiscoverWorker": "/api/Spaces-1/workers/discover{?host,port,type,proxyId}",
    "DynamicExtensionsFeaturesMetadata": "/api/dynamic-extensions/features/metadata",
    "DynamicExtensionsFeaturesValues": "/api/dynamic-extensions/features/values",
    "DynamicExtensionsScripts": "/api/dynamic-extensions/scripts",
    "EnabledFeatureToggles": "/api/configuration/enabled-feature-toggles",
    "Environments": "/api/Spaces-1/environments{/id}{?name,skip,ids,take,partialName}",
    "EnvironmentSortOrder": "/api/Spaces-1/environments/sortorder",
    "EnvironmentsSummary": "/api/Spaces-1/environments/summary{?ids,partialName,machinePartialName,roles,isDisabled,healthStatuses,commStyles,tenantIds,tenantTags,hideEmptyEnvironments,shellNames,deploymentTargetTypes}",
    "EventAgents": "/api/events/agents",
    "EventCategories": "/api/events/categories{?appliesTo}",
    "EventDocumentTypes": "/api/events/documenttypes",
    "EventGroups": "/api/events/groups{?appliesTo}",
    "Events": "/api/events{/id}{?skip,regarding,regardingAny,user,users,projects,projectGroups,environments,eventGroups,eventCategories,eventAgents,tags,tenants,from,to,internal,fromAutoId,toAutoId,documentTypes,asCsv,take,ids,spaces,includeSystem,excludeDifference}",
    "ExportProjects": "/api/Spaces-1/projects/import-export/export",
    "ExtensionStats": "/api/serverstatus/extensions",
    "ExternalSecurityGroupProviders": "/api/externalsecuritygroupproviders",
    "ExternalUserSearch": "/api/users/external-search{?partialName}",
    "FeaturesConfiguration": "/api/featuresconfiguration",
    "Feeds": "/api/feeds{/id}{?skip,take,ids,partialName,feedType,name}",
    "GitCredentials": "/api/Spaces-1/git-credentials{/id}{?skip,take,name}",
    "GitHubConnectivityCheck": "/api/githubissuetracker/connectivitycheck",
    "ImportProjects": "/api/Spaces-1/projects/import-export/import",
    "InsightsReports": "/api/Spaces-1/insights/reports{/id}{?skip,take}",
    "Interruptions": "/api/Spaces-1/interruptions{/id}{?skip,take,regarding,pendingOnly,ids}",
    "Invitations": "/api/users/invitations",
    "IssueTrackers": "/api/issuetrackers{?skip,take,ids,partialName}",
    "JiraConnectAppCredentialsTest": "/api/jiraintegration/connectivitycheck/connectapp",
    "JiraCredentialsTest": "/api/jiraintegration/connectivitycheck/jira",
    "LibraryVariables": "/api/Spaces-1/libraryvariablesets{/id}{?skip,contentType,take,ids,partialName}",
    "LifecyclePreviews": "/api/Spaces-1/lifecycles/previews{?ids}",
    "Lifecycles": "/api/Spaces-1/lifecycles{/id}{?skip,take,ids,partialName}",
    "LoginInitiated": "/api/authentication/checklogininitiated",
    "LogoIconCategories": "/api/icons/categories?cb=2023.1.9626",
    "LogoIcons": "/api/icons/all?cb=2023.1.9626",
    "MachineOperatingSystems": "/api/Spaces-1/machines/operatingsystem/names/all",
    "MachinePolicies": "/api/Spaces-1/machinepolicies{/id}{?skip,take,ids,partialName}",
    "MachinePolicyTemplate": "/api/Spaces-1/machinepolicies/template",
    "MachineRoles": "/api/Spaces-1/machineroles/all",
    "Machines": "/api/Spaces-1/machines{/id}{?skip,take,name,ids,partialName,roles,isDisabled,healthStatuses,commStyles,tenantIds,tenantTags,environmentIds,thumbprint,deploymentId,shellNames,deploymentTargetTypes}",
    "MachineShells": "/api/Spaces-1/machines/operatingsystem/shells/all",
    "MaintenanceConfiguration": "/api/maintenanceconfiguration",
    "MigrationsImport": "/api/migrations/import",
    "MigrationsPartialExport": "/api/migrations/partialexport",
    "OctopusServerClusterSummary": "/api/octopusservernodes/summary",
    "OctopusServerNodes": "/api/octopusservernodes{/id}{?skip,take,ids,partialName}",
    "PackageDeltaSignature": "/api/Spaces-1/packages/{packageId}/{version}/delta-signature",
    "PackageDeltaUpload": "/api/Spaces-1/packages/{packageId}/{baseVersion}/delta{?replace,overwriteMode}",
    "PackageNotesList": "/api/Spaces-1/packages/notes{?packageIds}",
    "Packages": "/api/Spaces-1/packages{/id}{?nuGetPackageId,filter,latest,skip,take,includeNotes}",
    "PackagesBulk": "/api/Spaces-1/packages/bulk{?ids}",
    "PackageUpload": "/api/Spaces-1/packages/raw{?replace,overwriteMode}",
    "PerformanceConfiguration": "/api/performanceconfiguration",
    "PermissionDescriptions": "/api/permissions/all",
    "ProjectGroups": "/api/Spaces-1/projectgroups{/id}{?skip,take,ids,partialName}",
    "ProjectImportFiles": "/api/Spaces-1/projects/import-export/import-files",
    "ProjectImportPreview": "/api/Spaces-1/projects/import-export/import/preview",
    "ProjectPulse": "/api/Spaces-1/projects/pulse{?projectIds}",
    "Projects": "/api/Spaces-1/projects{/id}{?name,skip,ids,clone,take,partialName,clonedFromProjectId}",
    "ProjectsExperimentalSummaries": "/api/Spaces-1/projects/experimental/summaries{?ids,isVersionControlled}",
    "ProjectTriggers": "/api/Spaces-1/projecttriggers{/id}{?skip,take,ids,runbooks}",
    "Proxies": "/api/Spaces-1/proxies{/id}{?skip,take,ids,partialName}",
    "Register": "/api/users/register",
    "Releases": "/api/Spaces-1/releases{/id}{?skip,ignoreChannelRules,take,ids}",
    "Reporting/DeploymentsCountedByWeek": "/api/Spaces-1/reporting/deployments-counted-by-week{?projectIds}",
    "RetentionDefaultConfiguration": "/api/configuration/retention-default",
    "RevokeUserSessions": "/api/users/{id}/revoke-sessions",
    "RunbookProcesses": "/api/Spaces-1/runbookProcesses{/id}{?skip,take,ids}",
    "RunbookRuns": "/api/Spaces-1/runbookRuns{/id}{?skip,take,ids,projects,environments,tenants,runbooks,taskState,partialName}",
    "Runbooks": "/api/Spaces-1/runbooks{/id}{?skip,take,ids,partialName,clone,projectIds}",
    "RunbookSnapshots": "/api/Spaces-1/runbookSnapshots{/id}{?skip,take,ids,publish}",
    "Scheduler": "/api/scheduler/{name}/logs{?verbose,tail}",
    "ScopedUserRoles": "/api/scopeduserroles{/id}{?skip,take,ids,partialName,spaces,includeSystem}",
    "ServerConfiguration": "/api/serverconfiguration",
    "ServerConfigurationSettings": "/api/serverconfiguration/settings",
    "ServerHealthStatus": "/api/serverstatus/health",
    "ServerStatus": "/api/serverstatus",
    "SignIn": "/api/users/login{?returnUrl}",
    "SignOut": "/api/users/logout",
    "SmtpConfiguration": "/api/smtpconfiguration",
    "SmtpIsConfigured": "/api/smtpconfiguration/isconfigured",
    "SpaceHome": "/api/{spaceId}",
    "Spaces": "/api/spaces{/id}{?skip,ids,take,partialName}",
    "SpaceSearch": "/api/spaces/{id}/search{?keyword}",
    "StepPackageDeploymentTargetTypes": "/api/steps/deploymenttargets",
    "Subscriptions": "/api/Spaces-1/subscriptions{/id}{?skip,take,ids,partialName,spaces}",
    "TagSets": "/api/Spaces-1/tagsets{/id}{?skip,take,ids,partialName}",
    "TagSetSortOrder": "/api/Spaces-1/tagsets/sortorder",
    "Tasks": "/api/tasks{/id}{?skip,active,environment,tenant,runbook,project,name,node,running,states,hasPendingInterruptions,hasWarningsOrErrors,take,ids,partialName,spaces,includeSystem,description,fromCompletedDate,toCompletedDate,fromQueueDate,toQueueDate,fromStartDate,toStartDate}",
    "TaskTypes": "/api/tasks/tasktypes",
    "TeamMembership": "/api/teammembership{?userId,spaces,includeSystem}",
    "TeamMembershipPreviewTeam": "/api/teammembership/previewteam",
    "Teams": "/api/teams{/id}{?skip,take,ids,partialName,spaces,includeSystem}",
    "TelemetryConfiguration": "/api/telemetryconfiguration",
    "TelemetryDownload": "/api/telemetry/download",
    "TelemetryLastTask": "/api/telemetry/lastTask",
    "TelemetrySend": "/api/telemetry/send",
    "Tenants": "/api/Spaces-1/tenants{/id}{?skip,projectId,name,tags,take,ids,clone,partialName,clonedFromTenantId}",
    "TenantsMissingVariables": "/api/Spaces-1/tenants/variables-missing{?tenantId,projectId,environmentId,includeDetails}",
    "TenantsStatus": "/api/Spaces-1/tenants/status",
    "TenantTagTest": "/api/Spaces-1/tenants/tag-test{?tenantIds,tags}",
    "TenantVariables": "/api/Spaces-1/tenantvariables/all{?projectId}",
    "Timezones": "/api/serverstatus/timezones",
    "UpgradeConfiguration": "/api/upgradeconfiguration",
    "UserAuthentication": "/api/users/authentication{/userId}",
    "UserIdentityMetadata": "/api/users/identity-metadata",
    "UserOnboarding": "/api/Spaces-1/useronboarding",
    "UserRoles": "/api/userroles{/id}{?skip,take,ids,partialName}",
    "Users": "/api/users{/id}{?skip,take,ids,filter}",
    "VariableNames": "/api/Spaces-1/variables/names{?project,runbook,projectEnvironmentsFilter,gitRef}",
    "VariablePreview": "/api/Spaces-1/variables/preview{?project,runbook,environment,channel,tenant,action,machine,role,gitRef}",
    "Variables": "/api/Spaces-1/variables{/id}{?ids}",
    "VersionControlClearCache": "/api/configuration/versioncontrol/clear-cache",
    "VersionRuleTest": "/api/Spaces-1/channels/rule-test{?version,versionRange,preReleaseTag,feetType}",
    "Web": "/app",
    "WorkerOperatingSystems": "/api/Spaces-1/workers/operatingsystem/names/all",
    "WorkerPools": "/api/Spaces-1/workerpools{/id}{?skip,ids,take,partialName}",
    "WorkerPoolsDynamicWorkerTypes": "/api/Spaces-1/workerpools/dynamicworkertypes",
    "WorkerPoolsSortOrder": "/api/Spaces-1/workerpools/sortorder",
    "WorkerPoolsSummary": "/api/Spaces-1/workerpools/summary{?ids,partialName,machinePartialName,isDisabled,healthStatuses,commStyles,hideEmptyWorkerPools,shellNames}",
    "WorkerPoolsSupportedTypes": "/api/Spaces-1/workerpools/supportedtypes",
    "Workers": "/api/Spaces-1/workers{/id}{?skip,take,name,ids,partialName,isDisabled,healthStatuses,commStyles,workerPoolIds,thumbprint,shellNames}",
    "WorkerShells": "/api/Spaces-1/workers/operatingsystem/shells/all",
    "WorkerToolsLatestImages": "/api/workertoolslatestimages"
  }
}`))
		}

		if r.URL.Path == "/api/Spaces-1/events" || r.URL.Path == "/api/events" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
  "ItemType": "Event",
  "TotalResults": 2095,
  "ItemsPerPage": 30,
  "NumberOfPages": 70,
  "LastPageNumber": 69,
  "Items": [
    {
      "Id": "Events-872850",
      "RelatedDocumentIds": [
        "Deployments-28317",
        "Projects-1492",
        "Releases-19095",
        "Environments-1022",
        "ServerTasks-268493",
        "Channels-1939"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-03T01:10:26.736+00:00",
      "Message": "Deploy to Development started for Random Failure release 0.0.1726 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28317'>Deploy to Development</a> started for <a href='#/projects/Projects-1492'>Random Failure</a> release <a href='#/releases/Releases-19095'>0.0.1726</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28317",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1492",
          "StartIndex": 34,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Releases-19095",
          "StartIndex": 57,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 69,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872850"
      }
    },
    {
      "Id": "Events-872834",
      "RelatedDocumentIds": [
        "Deployments-28317",
        "Projects-1492",
        "Releases-19095",
        "Environments-1022",
        "ServerTasks-268493",
        "Channels-1939"
      ],
      "Category": "DeploymentQueued",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-03T01:00:23.265+00:00",
      "Message": "Scheduled to deploy Random Failure release 0.0.1726 to Development at Friday, 03 March 2023 1:00 AM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1492'>Random Failure</a> release <a href='#/releases/Releases-19095'>0.0.1726</a> to <a href='#/environments/Environments-1022'>Development</a> at Friday, 03 March 2023 1:00 AM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1492",
          "StartIndex": 20,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Releases-19095",
          "StartIndex": 43,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 55,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872834"
      }
    },
    {
      "Id": "Events-872712",
      "RelatedDocumentIds": [
        "Deployments-28316",
        "Projects-1492",
        "Releases-19094",
        "Environments-1022",
        "ServerTasks-268476",
        "Channels-1939"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-03T00:10:24.849+00:00",
      "Message": "Deploy to Development started for Random Failure release 0.0.1725 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28316'>Deploy to Development</a> started for <a href='#/projects/Projects-1492'>Random Failure</a> release <a href='#/releases/Releases-19094'>0.0.1725</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28316",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1492",
          "StartIndex": 34,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Releases-19094",
          "StartIndex": 57,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 69,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872712"
      }
    },
    {
      "Id": "Events-872708",
      "RelatedDocumentIds": [
        "Deployments-28316",
        "Projects-1492",
        "Releases-19094",
        "Environments-1022",
        "ServerTasks-268476",
        "Channels-1939"
      ],
      "Category": "DeploymentQueued",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-03T00:00:23.502+00:00",
      "Message": "Scheduled to deploy Random Failure release 0.0.1725 to Development at Friday, 03 March 2023 12:00 AM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1492'>Random Failure</a> release <a href='#/releases/Releases-19094'>0.0.1725</a> to <a href='#/environments/Environments-1022'>Development</a> at Friday, 03 March 2023 12:00 AM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1492",
          "StartIndex": 20,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Releases-19094",
          "StartIndex": 43,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 55,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872708"
      }
    },
    {
      "Id": "Events-872534",
      "RelatedDocumentIds": [
        "Deployments-28315",
        "Projects-1492",
        "Releases-19093",
        "Environments-1022",
        "ServerTasks-268446",
        "Channels-1939"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T23:10:31.311+00:00",
      "Message": "Deploy to Development started for Random Failure release 0.0.1724 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28315'>Deploy to Development</a> started for <a href='#/projects/Projects-1492'>Random Failure</a> release <a href='#/releases/Releases-19093'>0.0.1724</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28315",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1492",
          "StartIndex": 34,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Releases-19093",
          "StartIndex": 57,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 69,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872534"
      }
    },
    {
      "Id": "Events-872530",
      "RelatedDocumentIds": [
        "Deployments-28315",
        "Projects-1492",
        "Releases-19093",
        "Environments-1022",
        "ServerTasks-268446",
        "Channels-1939"
      ],
      "Category": "DeploymentQueued",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T23:00:30.331+00:00",
      "Message": "Scheduled to deploy Random Failure release 0.0.1724 to Development at Thursday, 02 March 2023 11:00 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1492'>Random Failure</a> release <a href='#/releases/Releases-19093'>0.0.1724</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 11:00 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1492",
          "StartIndex": 20,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Releases-19093",
          "StartIndex": 43,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 55,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872530"
      }
    },
    {
      "Id": "Events-872448",
      "RelatedDocumentIds": [
        "Deployments-28313",
        "Projects-1492",
        "Releases-19092",
        "Environments-1022",
        "ServerTasks-268442",
        "Channels-1939"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T22:10:11.430+00:00",
      "Message": "Deploy to Development started for Random Failure release 0.0.1723 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28313'>Deploy to Development</a> started for <a href='#/projects/Projects-1492'>Random Failure</a> release <a href='#/releases/Releases-19092'>0.0.1723</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28313",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1492",
          "StartIndex": 34,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Releases-19092",
          "StartIndex": 57,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 69,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872448"
      }
    },
    {
      "Id": "Events-872444",
      "RelatedDocumentIds": [
        "Deployments-28313",
        "Projects-1492",
        "Releases-19092",
        "Environments-1022",
        "ServerTasks-268442",
        "Channels-1939"
      ],
      "Category": "DeploymentQueued",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T22:00:10.441+00:00",
      "Message": "Scheduled to deploy Random Failure release 0.0.1723 to Development at Thursday, 02 March 2023 10:00 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1492'>Random Failure</a> release <a href='#/releases/Releases-19092'>0.0.1723</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 10:00 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1492",
          "StartIndex": 20,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Releases-19092",
          "StartIndex": 43,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 55,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872444"
      }
    },
    {
      "Id": "Events-872410",
      "RelatedDocumentIds": [
        "Deployments-28312",
        "Projects-1108",
        "Releases-19090",
        "Environments-1022",
        "ServerTasks-268437",
        "Channels-1422"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:28:50.784+00:00",
      "Message": "Scheduled to deploy OctoFX release 2.10.257 to Development at Thursday, 02 March 2023 9:28 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1108'>OctoFX</a> release <a href='#/releases/Releases-19090'>2.10.257</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:28 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1108",
          "StartIndex": 20,
          "Length": 6
        },
        {
          "ReferencedDocumentId": "Releases-19090",
          "StartIndex": 35,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 47,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872410"
      }
    },
    {
      "Id": "Events-872364",
      "RelatedDocumentIds": [
        "Deployments-28311",
        "Projects-1110",
        "Releases-19088",
        "Environments-1022",
        "ServerTasks-268435",
        "Channels-1426"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:24:35.063+00:00",
      "Message": "Deploy to Development started for OctoPetShop release 2.1.0.187 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28311'>Deploy to Development</a> started for <a href='#/projects/Projects-1110'>OctoPetShop</a> release <a href='#/releases/Releases-19088'>2.1.0.187</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28311",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1110",
          "StartIndex": 34,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Releases-19088",
          "StartIndex": 54,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 67,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872364"
      }
    },
    {
      "Id": "Events-872360",
      "RelatedDocumentIds": [
        "Deployments-28311",
        "Projects-1110",
        "Releases-19088",
        "Environments-1022",
        "ServerTasks-268435",
        "Channels-1426"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:14:34.249+00:00",
      "Message": "Scheduled to deploy OctoPetShop release 2.1.0.187 to Development at Thursday, 02 March 2023 9:14 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1110'>OctoPetShop</a> release <a href='#/releases/Releases-19088'>2.1.0.187</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:14 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1110",
          "StartIndex": 20,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Releases-19088",
          "StartIndex": 40,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 53,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872360"
      }
    },
    {
      "Id": "Events-872356",
      "RelatedDocumentIds": [
        "Deployments-28310",
        "Projects-1110",
        "Releases-19085",
        "Environments-1023",
        "ServerTasks-268434",
        "Channels-1426"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:24:21.508+00:00",
      "Message": "Deploy to Test started for OctoPetShop release 2.1.0.186 to Test",
      "MessageHtml": "<a href='#/deployments/Deployments-28310'>Deploy to Test</a> started for <a href='#/projects/Projects-1110'>OctoPetShop</a> release <a href='#/releases/Releases-19085'>2.1.0.186</a> to <a href='#/environments/Environments-1023'>Test</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28310",
          "StartIndex": 0,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Projects-1110",
          "StartIndex": 27,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Releases-19085",
          "StartIndex": 47,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Environments-1023",
          "StartIndex": 60,
          "Length": 4
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872356"
      }
    },
    {
      "Id": "Events-872351",
      "RelatedDocumentIds": [
        "Deployments-28310",
        "Projects-1110",
        "Releases-19085",
        "Environments-1023",
        "ServerTasks-268434",
        "Channels-1426"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:14:20.669+00:00",
      "Message": "Scheduled to deploy OctoPetShop release 2.1.0.186 to Test at Thursday, 02 March 2023 9:14 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1110'>OctoPetShop</a> release <a href='#/releases/Releases-19085'>2.1.0.186</a> to <a href='#/environments/Environments-1023'>Test</a> at Thursday, 02 March 2023 9:14 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1110",
          "StartIndex": 20,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Releases-19085",
          "StartIndex": 40,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Environments-1023",
          "StartIndex": 53,
          "Length": 4
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872351"
      }
    },
    {
      "Id": "Events-872272",
      "RelatedDocumentIds": [
        "Deployments-28307",
        "Projects-1110",
        "Releases-19085",
        "Environments-1022",
        "ServerTasks-268430",
        "Channels-1426"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:22:08.247+00:00",
      "Message": "Deploy to Development started for OctoPetShop release 2.1.0.186 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28307'>Deploy to Development</a> started for <a href='#/projects/Projects-1110'>OctoPetShop</a> release <a href='#/releases/Releases-19085'>2.1.0.186</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28307",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1110",
          "StartIndex": 34,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Releases-19085",
          "StartIndex": 54,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 67,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872272"
      }
    },
    {
      "Id": "Events-872268",
      "RelatedDocumentIds": [
        "Deployments-28307",
        "Projects-1110",
        "Releases-19085",
        "Environments-1022",
        "ServerTasks-268430",
        "Channels-1426"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:12:07.450+00:00",
      "Message": "Scheduled to deploy OctoPetShop release 2.1.0.186 to Development at Thursday, 02 March 2023 9:12 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1110'>OctoPetShop</a> release <a href='#/releases/Releases-19085'>2.1.0.186</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:12 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1110",
          "StartIndex": 20,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Releases-19085",
          "StartIndex": 40,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 53,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872268"
      }
    },
    {
      "Id": "Events-872185",
      "RelatedDocumentIds": [
        "Deployments-28303",
        "Projects-1108",
        "Releases-19081",
        "Environments-1022",
        "ServerTasks-268425",
        "Channels-1422"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:18:23.160+00:00",
      "Message": "Deploy to Development started for OctoFX release 2.10.256 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28303'>Deploy to Development</a> started for <a href='#/projects/Projects-1108'>OctoFX</a> release <a href='#/releases/Releases-19081'>2.10.256</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28303",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1108",
          "StartIndex": 34,
          "Length": 6
        },
        {
          "ReferencedDocumentId": "Releases-19081",
          "StartIndex": 49,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 61,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872185"
      }
    },
    {
      "Id": "Events-872181",
      "RelatedDocumentIds": [
        "Deployments-28303",
        "Projects-1108",
        "Releases-19081",
        "Environments-1022",
        "ServerTasks-268425",
        "Channels-1422"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:08:22.283+00:00",
      "Message": "Scheduled to deploy OctoFX release 2.10.256 to Development at Thursday, 02 March 2023 9:08 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1108'>OctoFX</a> release <a href='#/releases/Releases-19081'>2.10.256</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:08 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1108",
          "StartIndex": 20,
          "Length": 6
        },
        {
          "ReferencedDocumentId": "Releases-19081",
          "StartIndex": 35,
          "Length": 8
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 47,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872181"
      }
    },
    {
      "Id": "Events-872154",
      "RelatedDocumentIds": [
        "Deployments-28300",
        "Projects-1109",
        "Releases-19077",
        "Environments-1022",
        "Tenants-830",
        "ServerTasks-268422",
        "Channels-1425"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:17:21.961+00:00",
      "Message": "Deploy to Development for _Internal started for Random Quotes release 2022.6.310 to Development for _Internal",
      "MessageHtml": "<a href='#/deployments/Deployments-28300'>Deploy to Development for _Internal</a> started for <a href='#/projects/Projects-1109'>Random Quotes</a> release <a href='#/releases/Releases-19077'>2022.6.310</a> to <a href='#/environments/Environments-1022'>Development</a> for <a href='#/tenants/Tenants-830'>_Internal</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28300",
          "StartIndex": 0,
          "Length": 35
        },
        {
          "ReferencedDocumentId": "Projects-1109",
          "StartIndex": 48,
          "Length": 13
        },
        {
          "ReferencedDocumentId": "Releases-19077",
          "StartIndex": 70,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 84,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Tenants-830",
          "StartIndex": 100,
          "Length": 9
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872154"
      }
    },
    {
      "Id": "Events-872150",
      "RelatedDocumentIds": [
        "Deployments-28300",
        "Projects-1109",
        "Releases-19077",
        "Environments-1022",
        "Tenants-830",
        "ServerTasks-268422",
        "Channels-1425"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:07:21.128+00:00",
      "Message": "Scheduled to deploy Random Quotes release 2022.6.310 to Development at Thursday, 02 March 2023 9:07 PM +00:00 for _Internal",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1109'>Random Quotes</a> release <a href='#/releases/Releases-19077'>2022.6.310</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:07 PM +00:00 for <a href='#/tenants/Tenants-830'>_Internal</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1109",
          "StartIndex": 20,
          "Length": 13
        },
        {
          "ReferencedDocumentId": "Releases-19077",
          "StartIndex": 42,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 56,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Tenants-830",
          "StartIndex": 114,
          "Length": 9
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872150"
      }
    },
    {
      "Id": "Events-872142",
      "RelatedDocumentIds": [
        "Deployments-28299",
        "Projects-1109",
        "Releases-19065",
        "Environments-1023",
        "Tenants-830",
        "ServerTasks-268421",
        "Channels-1425"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:17:08.760+00:00",
      "Message": "Deploy to Test for _Internal started for Random Quotes release 2022.6.309 to Test for _Internal",
      "MessageHtml": "<a href='#/deployments/Deployments-28299'>Deploy to Test for _Internal</a> started for <a href='#/projects/Projects-1109'>Random Quotes</a> release <a href='#/releases/Releases-19065'>2022.6.309</a> to <a href='#/environments/Environments-1023'>Test</a> for <a href='#/tenants/Tenants-830'>_Internal</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28299",
          "StartIndex": 0,
          "Length": 28
        },
        {
          "ReferencedDocumentId": "Projects-1109",
          "StartIndex": 41,
          "Length": 13
        },
        {
          "ReferencedDocumentId": "Releases-19065",
          "StartIndex": 63,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Environments-1023",
          "StartIndex": 77,
          "Length": 4
        },
        {
          "ReferencedDocumentId": "Tenants-830",
          "StartIndex": 86,
          "Length": 9
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872142"
      }
    },
    {
      "Id": "Events-872138",
      "RelatedDocumentIds": [
        "Deployments-28299",
        "Projects-1109",
        "Releases-19065",
        "Environments-1023",
        "Tenants-830",
        "ServerTasks-268421",
        "Channels-1425"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:07:07.843+00:00",
      "Message": "Scheduled to deploy Random Quotes release 2022.6.309 to Test at Thursday, 02 March 2023 9:07 PM +00:00 for _Internal",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1109'>Random Quotes</a> release <a href='#/releases/Releases-19065'>2022.6.309</a> to <a href='#/environments/Environments-1023'>Test</a> at Thursday, 02 March 2023 9:07 PM +00:00 for <a href='#/tenants/Tenants-830'>_Internal</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1109",
          "StartIndex": 20,
          "Length": 13
        },
        {
          "ReferencedDocumentId": "Releases-19065",
          "StartIndex": 42,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Environments-1023",
          "StartIndex": 56,
          "Length": 4
        },
        {
          "ReferencedDocumentId": "Tenants-830",
          "StartIndex": 107,
          "Length": 9
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872138"
      }
    },
    {
      "Id": "Events-872108",
      "RelatedDocumentIds": [
        "Deployments-28298",
        "Projects-1111",
        "Releases-19075",
        "Environments-1022",
        "ServerTasks-268419",
        "Channels-1427"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:16:06.981+00:00",
      "Message": "Deploy to Development started for PetClinic release 2023.03.02.1 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28298'>Deploy to Development</a> started for <a href='#/projects/Projects-1111'>PetClinic</a> release <a href='#/releases/Releases-19075'>2023.03.02.1</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28298",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1111",
          "StartIndex": 34,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Releases-19075",
          "StartIndex": 52,
          "Length": 12
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 68,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872108"
      }
    },
    {
      "Id": "Events-872103",
      "RelatedDocumentIds": [
        "Deployments-28298",
        "Projects-1111",
        "Releases-19075",
        "Environments-1022",
        "ServerTasks-268419",
        "Channels-1427"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:06:05.851+00:00",
      "Message": "Scheduled to deploy PetClinic release 2023.03.02.1 to Development at Thursday, 02 March 2023 9:06 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1111'>PetClinic</a> release <a href='#/releases/Releases-19075'>2023.03.02.1</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:06 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1111",
          "StartIndex": 20,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Releases-19075",
          "StartIndex": 38,
          "Length": 12
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 54,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872103"
      }
    },
    {
      "Id": "Events-872082",
      "RelatedDocumentIds": [
        "Deployments-28296",
        "Projects-1111",
        "Releases-19072",
        "Environments-1023",
        "ServerTasks-268416",
        "Channels-1427"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:15:52.104+00:00",
      "Message": "Deploy to Test started for PetClinic release 2023.03.02.0 to Test",
      "MessageHtml": "<a href='#/deployments/Deployments-28296'>Deploy to Test</a> started for <a href='#/projects/Projects-1111'>PetClinic</a> release <a href='#/releases/Releases-19072'>2023.03.02.0</a> to <a href='#/environments/Environments-1023'>Test</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28296",
          "StartIndex": 0,
          "Length": 14
        },
        {
          "ReferencedDocumentId": "Projects-1111",
          "StartIndex": 27,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Releases-19072",
          "StartIndex": 45,
          "Length": 12
        },
        {
          "ReferencedDocumentId": "Environments-1023",
          "StartIndex": 61,
          "Length": 4
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872082"
      }
    },
    {
      "Id": "Events-872075",
      "RelatedDocumentIds": [
        "Deployments-28296",
        "Projects-1111",
        "Releases-19072",
        "Environments-1023",
        "ServerTasks-268416",
        "Channels-1427"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:05:51.229+00:00",
      "Message": "Scheduled to deploy PetClinic release 2023.03.02.0 to Test at Thursday, 02 March 2023 9:05 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1111'>PetClinic</a> release <a href='#/releases/Releases-19072'>2023.03.02.0</a> to <a href='#/environments/Environments-1023'>Test</a> at Thursday, 02 March 2023 9:05 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1111",
          "StartIndex": 20,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Releases-19072",
          "StartIndex": 38,
          "Length": 12
        },
        {
          "ReferencedDocumentId": "Environments-1023",
          "StartIndex": 54,
          "Length": 4
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872075"
      }
    },
    {
      "Id": "Events-872046",
      "RelatedDocumentIds": [
        "Deployments-28294",
        "Projects-1111",
        "Releases-19072",
        "Environments-1022",
        "ServerTasks-268414",
        "Channels-1427"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:13:58.919+00:00",
      "Message": "Deploy to Development started for PetClinic release 2023.03.02.0 to Development",
      "MessageHtml": "<a href='#/deployments/Deployments-28294'>Deploy to Development</a> started for <a href='#/projects/Projects-1111'>PetClinic</a> release <a href='#/releases/Releases-19072'>2023.03.02.0</a> to <a href='#/environments/Environments-1022'>Development</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28294",
          "StartIndex": 0,
          "Length": 21
        },
        {
          "ReferencedDocumentId": "Projects-1111",
          "StartIndex": 34,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Releases-19072",
          "StartIndex": 52,
          "Length": 12
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 68,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872046"
      }
    },
    {
      "Id": "Events-872042",
      "RelatedDocumentIds": [
        "Deployments-28294",
        "Projects-1111",
        "Releases-19072",
        "Environments-1022",
        "ServerTasks-268414",
        "Channels-1427"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:03:58.132+00:00",
      "Message": "Scheduled to deploy PetClinic release 2023.03.02.0 to Development at Thursday, 02 March 2023 9:03 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1111'>PetClinic</a> release <a href='#/releases/Releases-19072'>2023.03.02.0</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:03 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1111",
          "StartIndex": 20,
          "Length": 9
        },
        {
          "ReferencedDocumentId": "Releases-19072",
          "StartIndex": 38,
          "Length": 12
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 54,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872042"
      }
    },
    {
      "Id": "Events-872036",
      "RelatedDocumentIds": [
        "Deployments-28293",
        "Projects-1262",
        "Releases-19070",
        "Environments-1022",
        "ServerTasks-268413",
        "Channels-1650"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:03:44.247+00:00",
      "Message": "Scheduled to deploy Snow Globe release 0.33.1.189 to Development at Thursday, 02 March 2023 9:03 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1262'>Snow Globe</a> release <a href='#/releases/Releases-19070'>0.33.1.189</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:03 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1262",
          "StartIndex": 20,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Releases-19070",
          "StartIndex": 39,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 53,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872036"
      }
    },
    {
      "Id": "Events-872008",
      "RelatedDocumentIds": [
        "Deployments-28289",
        "Projects-1262",
        "Releases-19067",
        "Environments-1022",
        "ServerTasks-268409",
        "Channels-1650"
      ],
      "Category": "DeploymentQueued",
      "UserId": "Users-261",
      "Username": "DemoSpaceManager",
      "IsService": true,
      "IdentityEstablishedWith": "API key 'DemoSpaceCreator.Octopus.APIKey' created Fri, 22 Jul 2022 17:46:35 GMT",
      "UserAgent": "OctopusClient-dotnet/11.6.3644 (Linux 5.4.0-1103-azure #109~18.04.1-Ubuntu SMP Wed Jan 25 20:53:00 UTC 2023; x64) NoneOrUnknown octo",
      "Occurred": "2023-03-02T21:03:09.656+00:00",
      "Message": "Scheduled to deploy Snow Globe release 0.33.1.188 to Development at Thursday, 02 March 2023 9:03 PM +00:00",
      "MessageHtml": "Scheduled to deploy <a href='#/projects/Projects-1262'>Snow Globe</a> release <a href='#/releases/Releases-19067'>0.33.1.188</a> to <a href='#/environments/Environments-1022'>Development</a> at Thursday, 02 March 2023 9:03 PM +00:00",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Projects-1262",
          "StartIndex": 20,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Releases-19067",
          "StartIndex": 39,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 53,
          "Length": 11
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": "52.143.76.163",
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-872008"
      }
    },
    {
      "Id": "Events-871995",
      "RelatedDocumentIds": [
        "Deployments-28287",
        "Projects-1109",
        "Releases-19065",
        "Environments-1022",
        "Tenants-830",
        "ServerTasks-268407",
        "Channels-1425"
      ],
      "Category": "DeploymentStarted",
      "UserId": "users-system",
      "Username": "system",
      "IsService": false,
      "IdentityEstablishedWith": "",
      "UserAgent": "Server",
      "Occurred": "2023-03-02T21:12:54.240+00:00",
      "Message": "Deploy to Development for _Internal started for Random Quotes release 2022.6.309 to Development for _Internal",
      "MessageHtml": "<a href='#/deployments/Deployments-28287'>Deploy to Development for _Internal</a> started for <a href='#/projects/Projects-1109'>Random Quotes</a> release <a href='#/releases/Releases-19065'>2022.6.309</a> to <a href='#/environments/Environments-1022'>Development</a> for <a href='#/tenants/Tenants-830'>_Internal</a>",
      "MessageReferences": [
        {
          "ReferencedDocumentId": "Deployments-28287",
          "StartIndex": 0,
          "Length": 35
        },
        {
          "ReferencedDocumentId": "Projects-1109",
          "StartIndex": 48,
          "Length": 13
        },
        {
          "ReferencedDocumentId": "Releases-19065",
          "StartIndex": 70,
          "Length": 10
        },
        {
          "ReferencedDocumentId": "Environments-1022",
          "StartIndex": 84,
          "Length": 11
        },
        {
          "ReferencedDocumentId": "Tenants-830",
          "StartIndex": 100,
          "Length": 9
        }
      ],
      "Comments": null,
      "Details": null,
      "ChangeDetails": {
        "DocumentContext": null,
        "Differences": null
      },
      "IpAddress": null,
      "SpaceId": "Spaces-792",
      "Links": {
        "Self": "/api/events/Events-871995"
      }
    }
  ],
  "Links": {
    "Self": "/api/events?skip=0&eventCategories=DeploymentQueued,DeploymentStarted&from=2%2F1%2F2023%2012%3A00%3A00%20AM%20%2B10%3A00&to=3%2F3%2F2023%2011%3A59%3A59%20PM%20%2B10%3A00&asCsv=false&take=30&excludeDifference=true",
    "Template": "/api/events{?skip,regarding,regardingAny,user,users,projects,projectGroups,environments,eventGroups,eventCategories,eventAgents,tags,tenants,from,to,internal,fromAutoId,toAutoId,documentTypes,asCsv,take,ids,spaces,includeSystem,excludeDifference}",
    "Page.All": "/api/events?skip=0&eventCategories=DeploymentQueued,DeploymentStarted&from=2%2F1%2F2023%2012%3A00%3A00%20AM%20%2B10%3A00&to=3%2F3%2F2023%2011%3A59%3A59%20PM%20%2B10%3A00&asCsv=false&take=2147483647&excludeDifference=true",
    "Page.Next": "/api/events?skip=30&eventCategories=DeploymentQueued,DeploymentStarted&from=2%2F1%2F2023%2012%3A00%3A00%20AM%20%2B10%3A00&to=3%2F3%2F2023%2011%3A59%3A59%20PM%20%2B10%3A00&asCsv=false&take=30&excludeDifference=true",
    "Page.Current": "/api/events?skip=0&eventCategories=DeploymentQueued,DeploymentStarted&from=2%2F1%2F2023%2012%3A00%3A00%20AM%20%2B10%3A00&to=3%2F3%2F2023%2011%3A59%3A59%20PM%20%2B10%3A00&asCsv=false&take=30&excludeDifference=true",
    "Page.Last": "/api/events?skip=2070&eventCategories=DeploymentQueued,DeploymentStarted&from=2%2F1%2F2023%2012%3A00%3A00%20AM%20%2B10%3A00&to=3%2F3%2F2023%2011%3A59%3A59%20PM%20%2B10%3A00&asCsv=false&take=30&excludeDifference=true"
  }
}`))
		}

	}))
	defer server.Close()

	// Act
	newSpaceClient, err := octoclient.CreateClient(server.URL, "Spaces-1", test.ApiKey)
	check := NewOctopusDeploymentQueuedTimeCheck(newSpaceClient, &config.OctolintConfig{}, "http://test.app", "Spaces-1", checks.OctopusClientPermissiveErrorHandler{})

	result, err := check.Execute(2)

	if err != nil {
		t.Fatal("Check produced an error")
	}

	// Assert
	if result.Severity() != checks.Warning {
		t.Fatal("Check should have returned a warning")
	}
}
