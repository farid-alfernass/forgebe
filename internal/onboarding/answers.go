package onboarding

// Answers stores responses from guided onboarding.
type Answers struct {
	ProjectName      string   `survey:"projectName"`
	ProjectMaturity  string   `survey:"projectMaturity"`
	ProjectType      string   `survey:"projectType"`
	PrimaryLanguage  string   `survey:"primaryLanguage"`
	Framework        string   `survey:"framework"`
	PackageManager   string   `survey:"packageManager"`
	Architecture     string   `survey:"architecture"`
	TestingPolicy    string   `survey:"testingPolicy"`
	DependencyPolicy string   `survey:"dependencyPolicy"`
	ChangePolicy     string   `survey:"changePolicy"`
	QualityPriority  string   `survey:"qualityPriority"`
	DeliveryMode     string   `survey:"deliveryMode"`
	AIModelStrategy  string   `survey:"aiModelStrategy"`
	SensitiveAreas   []string `survey:"sensitiveAreas"`
	OtherSensitive   string   `survey:"otherSensitive"`
}
