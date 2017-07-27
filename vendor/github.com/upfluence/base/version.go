package base

import "github.com/upfluence/base/version"




var (
	semVersion = &version.SemanticVersion{0, 0, 21}
	gitVersion = &version.GitVersion{"33cc5b306c4a", "https://github.com/upfluence/upfluence-if", "master"}
	baseVersion = &version.Version{semVersion, gitVersion}
	)


type Base struct {}
func (p *Base) Name() string { return "base" }

func (p *Base) Version() *version.Version {
  return baseVersion
}
