package profile

import "time"

// ProjectProfile represents the canonical project context stored by ForgeBE
type ProjectProfile struct {
	Version  string   `yaml:"version" json:"version"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Project  Project  `yaml:"project" json:"project"`
	Stack    Stack    `yaml:"stack" json:"stack"`
	Policy   Policy   `yaml:"policy" json:"policy"`
	Areas    Areas    `yaml:"areas" json:"areas"`
	AI       AI       `yaml:"ai" json:"ai"`
	Watch    Watch    `yaml:"watch" json:"watch"`
}

type Metadata struct {
	ProfileID  string    `yaml:"profile_id" json:"profile_id"`
	RepoPath   string    `yaml:"repo_path" json:"repo_path"`
	RepoRemote string    `yaml:"repo_remote" json:"repo_remote"`
	CreatedAt  time.Time `yaml:"created_at" json:"created_at"`
	UpdatedAt  time.Time `yaml:"updated_at" json:"updated_at"`
	ForgebeVer string    `yaml:"forgebe_version" json:"forgebe_version"`
	InitMethod string    `yaml:"init_method" json:"init_method"` // guided, scan, manual
}

type Project struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Maturity    string   `yaml:"maturity" json:"maturity"` // new, existing, legacy
	Type        string   `yaml:"type" json:"type"`         // monolith, service, library
	Languages   []string `yaml:"languages" json:"languages"`
}

type Stack struct {
	PrimaryLanguage string   `yaml:"primary_language" json:"primary_language"`
	Framework       string   `yaml:"framework,omitempty" json:"framework,omitempty"`
	Runtime         string   `yaml:"runtime,omitempty" json:"runtime,omitempty"`
	PackageManager  string   `yaml:"package_manager,omitempty" json:"package_manager,omitempty"`
	BuildTool       string   `yaml:"build_tool,omitempty" json:"build_tool,omitempty"`
	Architecture    string   `yaml:"architecture" json:"architecture"` // layered, clean, mvc, ddd, event-driven
	TestFramework   string   `yaml:"test_framework,omitempty" json:"test_framework,omitempty"`
	Linter          string   `yaml:"linter,omitempty" json:"linter,omitempty"`
	Formatter       string   `yaml:"formatter,omitempty" json:"formatter,omitempty"`
	CI              []string `yaml:"ci,omitempty" json:"ci,omitempty"`
}

type Policy struct {
	Testing      TestingPolicy    `yaml:"testing" json:"testing"`
	Dependencies DependencyPolicy `yaml:"dependencies" json:"dependencies"`
	Changes      ChangePolicy     `yaml:"changes" json:"changes"`
	Quality      QualityPolicy    `yaml:"quality" json:"quality"`
	Delivery     DeliveryPolicy   `yaml:"delivery" json:"delivery"`
}

type TestingPolicy struct {
	Required        bool   `yaml:"required" json:"required"`
	Coverage        int    `yaml:"coverage,omitempty" json:"coverage,omitempty"` // percentage
	Strategy        string `yaml:"strategy" json:"strategy"`                     // tdd, test-after, mixed
	UnitRequired    bool   `yaml:"unit_required" json:"unit_required"`
	IntegrationReqd bool   `yaml:"integration_required" json:"integration_required"`
}

type DependencyPolicy struct {
	AllowAddition   bool     `yaml:"allow_addition" json:"allow_addition"`
	RequireApproval bool     `yaml:"require_approval" json:"require_approval"`
	Forbidden       []string `yaml:"forbidden,omitempty" json:"forbidden,omitempty"`
}

type ChangePolicy struct {
	PreserveStructure bool     `yaml:"preserve_structure" json:"preserve_structure"`
	AllowRefactor     bool     `yaml:"allow_refactor" json:"allow_refactor"`
	AllowCrossFile    bool     `yaml:"allow_cross_file" json:"allow_cross_file"`
	ForbiddenPaths    []string `yaml:"forbidden_paths,omitempty" json:"forbidden_paths,omitempty"`
	RequireApproval   []string `yaml:"require_approval,omitempty" json:"require_approval,omitempty"` // patterns
}

type QualityPolicy struct {
	LintRequired   bool `yaml:"lint_required" json:"lint_required"`
	FormatRequired bool `yaml:"format_required" json:"format_required"`
	ReviewRequired bool `yaml:"review_required" json:"review_required"`
}

type DeliveryPolicy struct {
	Mode     string `yaml:"mode" json:"mode"`         // plan-first, direct, hybrid
	Priority string `yaml:"priority" json:"priority"` // speed, quality, safety, maintainability
}

type Areas struct {
	SourceRoots    []string `yaml:"source_roots,omitempty" json:"source_roots,omitempty"`
	TestRoots      []string `yaml:"test_roots,omitempty" json:"test_roots,omitempty"`
	SensitiveAreas []string `yaml:"sensitive_areas,omitempty" json:"sensitive_areas,omitempty"`
	PublicAPIs     []string `yaml:"public_apis,omitempty" json:"public_apis,omitempty"`
}

type AI struct {
	InteractionMode string   `yaml:"interaction_mode" json:"interaction_mode"` // guided, auto, hybrid
	ModelStrategy   string   `yaml:"model_strategy" json:"model_strategy"`     // single, multi, round-robin
	PreferredModels []string `yaml:"preferred_models,omitempty" json:"preferred_models,omitempty"`
}

// Watch configures the filesystem watcher behaviour for auto-sync.
type Watch struct {
	Recursive        bool          `yaml:"recursive" json:"recursive"`
	DebounceDuration time.Duration `yaml:"debounce_duration" json:"debounce_duration"`
	IgnorePatterns   []string      `yaml:"ignore_patterns" json:"ignore_patterns"`
	MatchPatterns    []string      `yaml:"match_patterns" json:"match_patterns"`
	FullResyncEvery  time.Duration `yaml:"full_resync_every" json:"full_resync_every"`
}

// DiscoveryReport represents the output of automatic codebase scanning
type DiscoveryReport struct {
	Timestamp   time.Time         `yaml:"timestamp" json:"timestamp"`
	RepoPath    string            `yaml:"repo_path" json:"repo_path"`
	Detections  Detections        `yaml:"detections" json:"detections"`
	Confidence  map[string]string `yaml:"confidence" json:"confidence"` // detector -> high/medium/low
	Assumptions []string          `yaml:"assumptions" json:"assumptions"`
	Suggestions []string          `yaml:"suggestions,omitempty" json:"suggestions,omitempty"`
}

type Detections struct {
	Languages    []LanguageDetection   `yaml:"languages" json:"languages"`
	Architecture ArchitectureDetection `yaml:"architecture" json:"architecture"`
	Testing      TestingDetection      `yaml:"testing" json:"testing"`
	Dependencies DependencyDetection   `yaml:"dependencies" json:"dependencies"`
	CI           CIDetection           `yaml:"ci" json:"ci"`
	Conventions  ConventionDetection   `yaml:"conventions" json:"conventions"`
}

type LanguageDetection struct {
	Language   string   `yaml:"language" json:"language"`
	Confidence string   `yaml:"confidence" json:"confidence"`
	Files      int      `yaml:"files" json:"files"`
	Extensions []string `yaml:"extensions" json:"extensions"`
	Manifest   string   `yaml:"manifest,omitempty" json:"manifest,omitempty"`
}

type ArchitectureDetection struct {
	Style      string   `yaml:"style" json:"style"`
	Confidence string   `yaml:"confidence" json:"confidence"`
	Layers     []string `yaml:"layers,omitempty" json:"layers,omitempty"`
	Patterns   []string `yaml:"patterns,omitempty" json:"patterns,omitempty"`
}

type TestingDetection struct {
	Framework  string   `yaml:"framework,omitempty" json:"framework,omitempty"`
	ConfigFile string   `yaml:"config_file,omitempty" json:"config_file,omitempty"`
	TestDirs   []string `yaml:"test_dirs,omitempty" json:"test_dirs,omitempty"`
	Coverage   bool     `yaml:"coverage" json:"coverage"`
}

type DependencyDetection struct {
	PackageManager string `yaml:"package_manager,omitempty" json:"package_manager,omitempty"`
	ManifestFile   string `yaml:"manifest_file,omitempty" json:"manifest_file,omitempty"`
	LockFile       string `yaml:"lock_file,omitempty" json:"lock_file,omitempty"`
	Dependencies   int    `yaml:"dependencies" json:"dependencies"`
}

type CIDetection struct {
	Providers  []string `yaml:"providers,omitempty" json:"providers,omitempty"`
	ConfigFile []string `yaml:"config_files,omitempty" json:"config_files,omitempty"`
}

type ConventionDetection struct {
	Linter    string `yaml:"linter,omitempty" json:"linter,omitempty"`
	Formatter string `yaml:"formatter,omitempty" json:"formatter,omitempty"`
	Style     string `yaml:"style,omitempty" json:"style,omitempty"`
}
