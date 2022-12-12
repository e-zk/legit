package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitRepo struct {
	r *git.Repository
	h plumbing.Hash
}

func Open(path string, ref string) (*GitRepo, error) {
	var err error
	g := GitRepo{}
	g.r, err = git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}

	if ref == "" {
		head, err := g.r.Head()
		if err != nil {
			return nil, fmt.Errorf("getting head of %s: %w", path, err)
		}
		g.h = head.Hash()
	} else {
		hash, err := g.r.ResolveRevision(plumbing.Revision(ref))
		if err != nil {
			return nil, fmt.Errorf("resolving rev %s for %s: %w", ref, path, err)
		}
		g.h = *hash
	}
	return &g, nil
}

func (g *GitRepo) Commits() ([]*object.Commit, error) {
	ci, err := g.r.Log(&git.LogOptions{From: g.h})
	if err != nil {
		return nil, fmt.Errorf("commits from ref: %w", err)
	}

	commits := []*object.Commit{}
	ci.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	return commits, nil
}

func (g *GitRepo) LastCommit() (*object.Commit, error) {
	c, err := g.r.CommitObject(g.h)
	if err != nil {
		return nil, fmt.Errorf("last commit: %w", err)
	}
	return c, nil
}

func (g *GitRepo) FileContent(path string) (string, error) {
	c, err := g.r.CommitObject(g.h)
	if err != nil {
		return "", fmt.Errorf("commit object: %w", err)
	}

	tree, err := c.Tree()
	if err != nil {
		return "", fmt.Errorf("file tree: %w", err)
	}

	file, err := tree.File(path)
	if err != nil {
		return "", err
	}

	return file.Contents()
}

func (g *GitRepo) Tags() ([]*object.Tag, error) {
	ti, err := g.r.TagObjects()
	if err != nil {
		return nil, fmt.Errorf("tag objects: %w", err)
	}

	tags := []*object.Tag{}

	_ = ti.ForEach(func(t *object.Tag) error {
		tags = append(tags, t)
		return nil
	})

	return tags, nil
}

func (g *GitRepo) Branches() ([]*plumbing.Reference, error) {
	bi, err := g.r.Branches()
	if err != nil {
		return nil, fmt.Errorf("branchs: %w", err)
	}

	branches := []*plumbing.Reference{}

	_ = bi.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref)
		return nil
	})

	return branches, nil
}
