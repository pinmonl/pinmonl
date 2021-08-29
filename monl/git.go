package monl

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pinmonl/pinmonl/monlurl"
	"github.com/pinmonl/pinmonl/monlvars"
)

type Git struct {
	//
}

func NewGit() (*Git, error) {
	return &Git{}, nil
}

func (g *Git) MonlerName() string {
	return monlvars.Git
}

func (g *Git) Open(u *monlurl.Url) (Object, error) {
	return nil, nil
}

type GitRepo struct {
	SimpleObject
	url  string
	dir  string
	repo *git.Repository
}

func newGitRepo(url string) (*GitRepo, error) {
	return &GitRepo{
		url: url,
	}, nil
}

func (g *GitRepo) ObjectName() string {
	return monlvars.ObjectGitRepo
}

func (g *GitRepo) String() string {
	return monlvars.Git + monlurl.Separator + g.url
}

func (g *GitRepo) Populate() error {
	var err error
	if err = g.populateDir(); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			g.Close()
		}
	}()

	if err = g.populateRepo(); err != nil {
		return err
	}
	if err = g.populateReleases(); err != nil {
		return err
	}
	return nil
}

func (g *GitRepo) populateDir() error {
	dir, err := ioutil.TempDir("", "pml-git")
	if err != nil {
		return err
	}
	g.dir = dir
	return nil
}

func (g *GitRepo) populateRepo() error {
	repo, err := git.PlainClone(g.dir, false, &git.CloneOptions{
		URL:      g.url,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	g.repo = repo
	return nil
}

func (g *GitRepo) populateReleases() error {
	iter, err := g.repo.Tags()
	if err != nil {
		return err
	}

	stats := make([]Stat, 0)
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		// Parse annotated tag
		tag, err := g.repo.TagObject(ref.Hash())
		if err == nil {
			stat := SimpleStat{
				at:    tag.Tagger.When,
				value: tag.Name,
				hash:  tag.Hash.String(),
			}
			stats = append(stats, stat)
			return nil
		}
		// Parse lightweight tag
		commit, err := g.repo.CommitObject(ref.Hash())
		if err == nil {
			stat := SimpleStat{
				at:    commit.Committer.When,
				value: strings.TrimPrefix(ref.Name().String(), "refs/tags/"),
				hash:  commit.Hash.String(),
			}
			stats = append(stats, stat)
			return nil
		}
		return nil
	}); err != nil {
		return err
	}

	semverTags, err := NewSemverTags(stats, nil)
	if err != nil {
		return err
	}
	sort.Sort(BySemver(semverTags.Stable))

	statGroup := SimpleStatGroup{
		groupName: monlvars.StatRelease,
		isAppend:  false,
		stats:     semverTags.Stable,
	}
	if len(semverTags.Stable) > 0 {
		statGroup.latest = semverTags.Stable[:1]
	}
	g.stats = append(g.stats, statGroup)
	return nil
}

func (g *GitRepo) Close() error {
	if g.dir != "" {
		return nil
	}
	return os.RemoveAll(g.dir)
}

var _ Monler = &Git{}
var _ Object = &GitRepo{}
