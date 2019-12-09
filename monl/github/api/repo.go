package api

import "time"

// Repo defines the structure of repository
type Repo struct {
	ForkCount int
	UpdatedAt time.Time

	Stargazers struct {
		TotalCount int
	}
	Watchers struct {
		TotalCount int
	}
	OpenIssues struct {
		TotalCount int
	} `graphql:"issues(filterBy: {states: OPEN})"`
	DefaultBranchRef struct {
		Name string
	}
	Releases struct {
		PageInfo   *PageInfo
		TotalCount int
		Nodes      []RepoReleaseNode
	} `graphql:"releases(orderBy: {field: CREATED_AT, direction: DESC}, first: 1)"`
	Tags struct {
		TotalCount int
	} `graphql:"refs(refPrefix: \"refs/tags/\")"`
}

// RepoReleases handles releases of repository with pagination
type RepoReleases struct {
	Releases struct {
		PageInfo   *PageInfo
		TotalCount int
		Nodes      []RepoReleaseNode
	} `graphql:"releases(orderBy: {field: CREATED_AT, direction: DESC}, first: $first, after: $after, before: $before, last: $last)"`
}

// RepoTags handles tags of repository with pagination
type RepoTags struct {
	Tags struct {
		PageInfo   *PageInfo
		TotalCount int
		Nodes      []RepoTagNode
	} `graphql:"refs(refPrefix: \"refs/tags/\", orderBy: {field: TAG_COMMIT_DATE, direction: DESC}, first: $first, after: $after, before: $before, after: $after)"`
}

// RepoReleaseNode defines the structure of release
type RepoReleaseNode struct {
	CreatedAt time.Time
	Name      string
	TagName   string
}

// RepoTagNode defines the structure of tag
type RepoTagNode struct {
	Name string

	Target struct {
		Commit *struct {
			CommittedDate time.Time
		} `graphql:"... on Commit"`
		Tag *struct {
			Tagger struct {
				Date time.Time
			}
		} `graphql:"... on Tag"`
	}
}
