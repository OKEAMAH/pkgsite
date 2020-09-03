// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"path"
	"strings"
	"time"

	"golang.org/x/mod/module"
	"golang.org/x/pkgsite/internal/licenses"
	"golang.org/x/pkgsite/internal/source"
	"golang.org/x/pkgsite/internal/stdlib"
)

const (
	// LatestVersion signifies the latest available version in requests to the
	// proxy client.
	LatestVersion = "latest"

	// MasterVersion signifies the version at master.
	MasterVersion = "master"

	// UnknownModulePath signifies that the module path for a given package
	// path is ambiguous or not known. This is because requests to the
	// frontend can come in the form of <import-path>[@<version>], and it is
	// not clear which part of the import-path is the module path.
	UnknownModulePath = "unknownModulePath"
)

// ModuleInfo holds metadata associated with a module.
type ModuleInfo struct {
	ModulePath        string
	Version           string
	CommitTime        time.Time
	IsRedistributable bool
	HasGoMod          bool // whether the module zip has a go.mod file
	SourceInfo        *source.Info
}

// VersionMap holds metadata associated with module queries for a version.
type VersionMap struct {
	ModulePath       string
	RequestedVersion string
	ResolvedVersion  string
	GoModPath        string
	Status           int
	Error            string
	UpdatedAt        time.Time
}

// SeriesPath returns the series path for the module.
//
// A series is a group of modules that share the same base path and are assumed
// to be major-version variants.
//
// The series path is the module path without the version. For most modules,
// this will be the module path for all module versions with major version 0 or
// 1. For gopkg.in modules, the series path does not correspond to any module
// version.
//
// Examples:
// The module paths "a/b" and "a/b/v2"  both have series path "a/b".
// The module paths "gopkg.in/yaml.v1" and "gopkg.in/yaml.v2" both have series
// path "gopkg.in/yaml".
func (v *ModuleInfo) SeriesPath() string {
	return SeriesPathForModule(v.ModulePath)
}

// SeriesPathForModule returns the series path for the provided modulePath.
func SeriesPathForModule(modulePath string) string {
	seriesPath, _, _ := module.SplitPathVersion(modulePath)
	return seriesPath
}

// Suffix returns the suffix of the fullPath. It assumes that basePath is a
// prefix of fullPath. If fullPath and basePath are the same, the empty string
// is returned.
func Suffix(fullPath, basePath string) string {
	return strings.TrimPrefix(strings.TrimPrefix(fullPath, basePath), "/")
}

// V1Path returns the path for version 1 of the package whose import path
// is fullPath. If modulePath is the standard library, then V1Path returns
// fullPath.
func V1Path(fullPath, modulePath string) string {
	if modulePath == stdlib.ModulePath {
		return fullPath
	}
	return path.Join(SeriesPathForModule(modulePath), Suffix(fullPath, modulePath))
}

// A Module is a specific, reproducible build of a module.
type Module struct {
	LegacyModuleInfo
	// Licenses holds all licenses within this module version, including those
	// that may be contained in nested subdirectories.
	Licenses []*licenses.License
	Units    []*Unit

	LegacyPackages []*LegacyPackage
}

// IndexVersion holds the version information returned by the module index.
type IndexVersion struct {
	Path      string
	Version   string
	Timestamp time.Time
}

// ModuleVersionState holds a worker module version state.
type ModuleVersionState struct {
	ModulePath string
	Version    string

	// IndexTimestamp is the timestamp received from the Index for this version,
	// which should correspond to the time this version was committed to the
	// Index.
	IndexTimestamp time.Time
	// CreatedAt is the time this version was originally inserted into the
	// module version state table.
	CreatedAt time.Time

	// Status is the most recent HTTP status code received from the Fetch service
	// for this version, or nil if no request to the fetch service has been made.
	Status int
	// Error is the most recent HTTP response body received from the Fetch
	// service, for a response with an unsuccessful status code. It is used for
	// debugging only, and has no semantic significance.
	Error string
	// TryCount is the number of times a fetch of this version has been
	// attempted.
	TryCount int
	// LastProcessedAt is the last time this version was updated with a result
	// from the fetch service.
	LastProcessedAt *time.Time
	// NextProcessedAfter is the next time a fetch for this version should be
	// attempted.
	NextProcessedAfter time.Time

	// AppVersion is the value of the GAE_VERSION environment variable, which is
	// set by app engine. It is a timestamp in the format 20190709t112655 that
	// is close to, but not the same as, the deployment time. For example, the
	// deployment time for the above timestamp might be Jul 9, 2019, 11:29:59 AM.
	AppVersion string

	// GoModPath is the path declared in the go.mod file.
	GoModPath string

	// NumPackages it the number of packages that were processed as part of the
	// module (regardless of whether the processing was successful).
	NumPackages *int
}

// PackageVersionState holds a worker package version state. It is associated
// with a given module version state.
type PackageVersionState struct {
	PackagePath string
	ModulePath  string
	Version     string
	Status      int
	Error       string
}

// SearchResult represents a single search result from SearchDocuments.
type SearchResult struct {
	Name        string
	PackagePath string
	ModulePath  string
	Version     string
	Synopsis    string
	Licenses    []string

	CommitTime time.Time
	// Score is used to sort items in an array of SearchResult.
	Score float64

	// NumImportedBy is the number of packages that import PackagePath.
	NumImportedBy uint64

	// NumResults is the total number of packages that were returned for this
	// search.
	NumResults uint64
	// Approximate reports whether NumResults is an approximate count. NumResults
	// can be approximate if search scanned only a subset of documents, and
	// result count is estimated using the hyperloglog algorithm.
	Approximate bool
}
