// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    packageJSON, err := UnmarshalPackageJSON(bytes)
//    bytes, err = packageJSON.Marshal()

package schemas

import "encoding/json"

func UnmarshalPackageJSON(data []byte) (PackageJSON, error) {
	var r PackageJSON
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *PackageJSON) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PackageJSON struct {
	Schema            string                       `json:"$schema"`
	Title             string                       `json:"title"`
	Definitions       Definitions                  `json:"definitions"`
	Type              TypeElement                  `json:"type"`
	PatternProperties PackageJSONPatternProperties `json:"patternProperties"`
	Properties        PackageJSONProperties        `json:"properties"`
	AnyOf             []AnyOf                      `json:"anyOf"`
	ID                string                       `json:"$id"`
}

type AnyOf struct {
	Type     TypeElement `json:"type"`
	Not      Not         `json:"not"`
	Required []string    `json:"required,omitempty"`
}

type Not struct {
	Required []string `json:"required"`
}

type Definitions struct {
	Person                        Person                  `json:"person"`
	Dependency                    Dependency              `json:"dependency"`
	DevDependency                 Dependency              `json:"devDependency"`
	OptionalDependency            Dependency              `json:"optionalDependency"`
	PeerDependency                Dependency              `json:"peerDependency"`
	PeerDependencyMeta            PeerDependencyMeta      `json:"peerDependencyMeta"`
	License                       DefinitionsLicenseClass `json:"license"`
	ScriptsInstallAfter           Scripts                 `json:"scriptsInstallAfter"`
	ScriptsPublishAfter           Scripts                 `json:"scriptsPublishAfter"`
	ScriptsRestart                Scripts                 `json:"scriptsRestart"`
	ScriptsStart                  Scripts                 `json:"scriptsStart"`
	ScriptsStop                   Scripts                 `json:"scriptsStop"`
	ScriptsTest                   Scripts                 `json:"scriptsTest"`
	ScriptsUninstallBefore        Scripts                 `json:"scriptsUninstallBefore"`
	ScriptsVersionBefore          Scripts                 `json:"scriptsVersionBefore"`
	PackageExportsEntryPath       PackageExportsEntryPath `json:"packageExportsEntryPath"`
	PackageExportsEntryObject     PackagePortsEntryObject `json:"packageExportsEntryObject"`
	PackageExportsEntry           PackageExportsEntry     `json:"packageExportsEntry"`
	PackageExportsFallback        PackageExportsFallback  `json:"packageExportsFallback"`
	PackageExportsEntryOrFallback PackageExportsEntry     `json:"packageExportsEntryOrFallback"`
	PackageImportsEntryPath       PackageImportsEntryPath `json:"packageImportsEntryPath"`
	PackageImportsEntryObject     PackagePortsEntryObject `json:"packageImportsEntryObject"`
	PackageImportsEntry           PackageExportsEntry     `json:"packageImportsEntry"`
	PackageImportsFallback        PackageExportsFallback  `json:"packageImportsFallback"`
	PackageImportsEntryOrFallback PackageExportsEntry     `json:"packageImportsEntryOrFallback"`
	FundingURL                    FundingURL              `json:"fundingUrl"`
	FundingWay                    FundingWay              `json:"fundingWay"`
	DevEngineDependency           DevEngineDependency     `json:"devEngineDependency"`
}

type Dependency struct {
	Description          string       `json:"description"`
	Type                 TypeElement  `json:"type"`
	AdditionalProperties EngineStrict `json:"additionalProperties"`
}

type EngineStrict struct {
	Type TypeElement `json:"type"`
}

type DevEngineDependency struct {
	Description string                        `json:"description"`
	Type        TypeElement                   `json:"type"`
	Required    []string                      `json:"required"`
	Properties  DevEngineDependencyProperties `json:"properties"`
}

type DevEngineDependencyProperties struct {
	Name    Description `json:"name"`
	Version Description `json:"version"`
	OnFail  TypeClass   `json:"onFail"`
}

type Description struct {
	Type        TypeElement `json:"type"`
	Description string      `json:"description"`
}

type TypeClass struct {
	Type        TypeElement `json:"type"`
	Enum        []string    `json:"enum"`
	Description string      `json:"description"`
	Default     *string     `json:"default,omitempty"`
}

type FundingURL struct {
	Type        TypeElement `json:"type"`
	Format      string      `json:"format"`
	Description string      `json:"description"`
}

type FundingWay struct {
	Type                 TypeElement          `json:"type"`
	Description          string               `json:"description"`
	Properties           FundingWayProperties `json:"properties"`
	AdditionalProperties bool                 `json:"additionalProperties"`
	Required             []string             `json:"required"`
}

type FundingWayProperties struct {
	URL  Author      `json:"url"`
	Type Description `json:"type"`
}

type Author struct {
	Ref string `json:"$ref"`
}

type DefinitionsLicenseClass struct {
	AnyOf []OneOfElement `json:"anyOf"`
}

type OneOfElement struct {
	Type *TypeElement `json:"type,omitempty"`
	Enum []string     `json:"enum,omitempty"`
}

type PackageExportsEntry struct {
	OneOf []Author `json:"oneOf"`
}

type PackagePortsEntryObject struct {
	Type                 TypeElement                         `json:"type"`
	Description          string                              `json:"description"`
	Properties           PackageExportsEntryObjectProperties `json:"properties"`
	PatternProperties    map[string]LicenseValue             `json:"patternProperties"`
	AdditionalProperties bool                                `json:"additionalProperties"`
}

type LicenseValue struct {
	Ref         Ref    `json:"$ref"`
	Description string `json:"description"`
}

type PackageExportsEntryObjectProperties struct {
	Require    LicenseValue `json:"require"`
	Import     LicenseValue `json:"import"`
	ModuleSync *ModuleSync  `json:"module-sync,omitempty"`
	Node       LicenseValue `json:"node"`
	Default    LicenseValue `json:"default"`
	Types      LicenseValue `json:"types"`
}

type ModuleSync struct {
	Ref         Ref    `json:"$ref"`
	Comment     string `json:"$comment"`
	Description string `json:"description"`
}

type PackageExportsEntryPath struct {
	Type        []string `json:"type"`
	Description string   `json:"description"`
	Pattern     string   `json:"pattern"`
}

type PackageExportsFallback struct {
	Type        PackageExportsFallbackType `json:"type"`
	Description string                     `json:"description"`
	Items       Author                     `json:"items"`
}

type PackageImportsEntryPath struct {
	Type        []string `json:"type"`
	Description string   `json:"description"`
}

type PeerDependencyMeta struct {
	Description          string                                 `json:"description"`
	Type                 TypeElement                            `json:"type"`
	AdditionalProperties PeerDependencyMetaAdditionalProperties `json:"additionalProperties"`
}

type PeerDependencyMetaAdditionalProperties struct {
	Type                 TypeElement      `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           PurpleProperties `json:"properties"`
}

type PurpleProperties struct {
	Optional Description `json:"optional"`
}

type Person struct {
	Description string           `json:"description"`
	Type        []TypeElement    `json:"type"`
	Required    []string         `json:"required"`
	Properties  PersonProperties `json:"properties"`
}

type PersonProperties struct {
	Name  EngineStrict `json:"name"`
	URL   Email        `json:"url"`
	Email Email        `json:"email"`
}

type Email struct {
	Type   TypeElement `json:"type"`
	Format string      `json:"format"`
}

type Scripts struct {
	Description                string      `json:"description"`
	Type                       TypeElement `json:"type"`
	XIntellijLanguageInjection string      `json:"x-intellij-language-injection"`
}

type PackageJSONPatternProperties struct {
	Empty Empty `json:"^_"`
}

type Empty struct {
	Description string `json:"description"`
	TsType      string `json:"tsType"`
}

type PackageJSONProperties struct {
	Name                 Name                   `json:"name"`
	Version              Description            `json:"version"`
	Description          Description            `json:"description"`
	Keywords             CPUElement             `json:"keywords"`
	Homepage             Description            `json:"homepage"`
	Bugs                 Bugs                   `json:"bugs"`
	License              LicenseValue           `json:"license"`
	Licenses             Licenses               `json:"licenses"`
	Author               Author                 `json:"author"`
	Contributors         PackageExportsFallback `json:"contributors"`
	Maintainers          PackageExportsFallback `json:"maintainers"`
	Files                CPUElement             `json:"files"`
	Main                 Description            `json:"main"`
	Exports              Exports                `json:"exports"`
	Imports              Config                 `json:"imports"`
	Bin                  Bin                    `json:"bin"`
	Type                 TypeClass              `json:"type"`
	Types                Description            `json:"types"`
	Typings              Description            `json:"typings"`
	TypesVersions        TypesVersions          `json:"typesVersions"`
	Man                  Man                    `json:"man"`
	Directories          Directories            `json:"directories"`
	Repository           Repository             `json:"repository"`
	Funding              Funding                `json:"funding"`
	Scripts              ScriptsClass           `json:"scripts"`
	Config               Config                 `json:"config"`
	Dependencies         Author                 `json:"dependencies"`
	DevDependencies      Author                 `json:"devDependencies"`
	OptionalDependencies Author                 `json:"optionalDependencies"`
	PeerDependencies     Author                 `json:"peerDependencies"`
	PeerDependenciesMeta Author                 `json:"peerDependenciesMeta"`
	BundleDependencies   BundleDependencies     `json:"bundleDependencies"`
	BundledDependencies  BundleDependencies     `json:"bundledDependencies"`
	Resolutions          Description            `json:"resolutions"`
	Overrides            Description            `json:"overrides"`
	PackageManager       PackageManager         `json:"packageManager"`
	Engines              Engines                `json:"engines"`
	Volta                Volta                  `json:"volta"`
	EngineStrict         EngineStrict           `json:"engineStrict"`
	OS                   CPUElement             `json:"os"`
	CPU                  CPUElement             `json:"cpu"`
	DevEngines           DevEngines             `json:"devEngines"`
	PreferGlobal         Description            `json:"preferGlobal"`
	Private              Private                `json:"private"`
	PublishConfig        PublishConfig          `json:"publishConfig"`
	Dist                 Dist                   `json:"dist"`
	Readme               EngineStrict           `json:"readme"`
	Module               Description            `json:"module"`
	Esnext               Esnext                 `json:"esnext"`
	Workspaces           Workspaces             `json:"workspaces"`
	Jspm                 Author                 `json:"jspm"`
	EslintConfig         Author                 `json:"eslintConfig"`
	Prettier             Author                 `json:"prettier"`
	Stylelint            Author                 `json:"stylelint"`
	Ava                  Author                 `json:"ava"`
	Release              Author                 `json:"release"`
	Jscpd                Author                 `json:"jscpd"`
	Pnpm                 Config                 `json:"pnpm"`
	Stackblitz           Config                 `json:"stackblitz"`
}

type Bin struct {
	Type                 []TypeElement `json:"type"`
	AdditionalProperties EngineStrict  `json:"additionalProperties"`
}

type Bugs struct {
	Description string         `json:"description"`
	Type        []TypeElement  `json:"type"`
	Properties  BugsProperties `json:"properties"`
}

type BugsProperties struct {
	URL   FundingURL `json:"url"`
	Email FundingURL `json:"email"`
}

type BundleDependencies struct {
	Description string  `json:"description"`
	OneOf       []OneOf `json:"oneOf"`
}

type OneOf struct {
	Type  string        `json:"type"`
	Items *EngineStrict `json:"items,omitempty"`
}

type CPUProperties struct {
	Packages CPUElement `json:"packages"`
	Nohoist  CPUElement `json:"nohoist"`
}

type CPUElement struct {
	Description *string                    `json:"description,omitempty"`
	Type        PackageExportsFallbackType `json:"type"`
	Items       *EngineStrict              `json:"items,omitempty"`
	Properties  *CPUProperties             `json:"properties,omitempty"`
}

type ConfigProperties struct {
	Overrides                   *Description             `json:"overrides,omitempty"`
	PackageExtensions           *Config                  `json:"packageExtensions,omitempty"`
	PeerDependencyRules         *PeerDependencyRules     `json:"peerDependencyRules,omitempty"`
	NeverBuiltDependencies      *CPUElement              `json:"neverBuiltDependencies,omitempty"`
	OnlyBuiltDependencies       *CPUElement              `json:"onlyBuiltDependencies,omitempty"`
	OnlyBuiltDependenciesFile   *Description             `json:"onlyBuiltDependenciesFile,omitempty"`
	IgnoredBuiltDependencies    *CPUElement              `json:"ignoredBuiltDependencies,omitempty"`
	AllowedDeprecatedVersions   *Description             `json:"allowedDeprecatedVersions,omitempty"`
	PatchedDependencies         *Description             `json:"patchedDependencies,omitempty"`
	AllowNonAppliedPatches      *Description             `json:"allowNonAppliedPatches,omitempty"`
	AllowUnusedPatches          *Description             `json:"allowUnusedPatches,omitempty"`
	UpdateConfig                *UpdateConfig            `json:"updateConfig,omitempty"`
	ConfigDependencies          *Description             `json:"configDependencies,omitempty"`
	AuditConfig                 *AuditConfig             `json:"auditConfig,omitempty"`
	RequiredScripts             *CPUElement              `json:"requiredScripts,omitempty"`
	SupportedArchitectures      *Config                  `json:"supportedArchitectures,omitempty"`
	IgnoredOptionalDependencies *CPUElement              `json:"ignoredOptionalDependencies,omitempty"`
	ExecutionEnv                *ExecutionEnv            `json:"executionEnv,omitempty"`
	OS                          *OneOf                   `json:"os,omitempty"`
	CPU                         *OneOf                   `json:"cpu,omitempty"`
	Libc                        *OneOf                   `json:"libc,omitempty"`
	InstallDependencies         *Description             `json:"installDependencies,omitempty"`
	StartCommand                *PackageImportsEntryPath `json:"startCommand,omitempty"`
	CompileTrigger              *Private                 `json:"compileTrigger,omitempty"`
	Env                         *Description             `json:"env,omitempty"`
}

type Config struct {
	Description          string                   `json:"description"`
	Type                 TypeElement              `json:"type"`
	AdditionalProperties bool                     `json:"additionalProperties"`
	PatternProperties    *ConfigPatternProperties `json:"patternProperties,omitempty"`
	Properties           *ConfigProperties        `json:"properties,omitempty"`
}

type AuditConfig struct {
	Type                 TypeElement           `json:"type"`
	Properties           AuditConfigProperties `json:"properties"`
	AdditionalProperties bool                  `json:"additionalProperties"`
}

type AuditConfigProperties struct {
	IgnoreCves  IgnoreCvesClass `json:"ignoreCves"`
	IgnoreGhsas IgnoreCvesClass `json:"ignoreGhsas"`
}

type IgnoreCvesClass struct {
	Description string                     `json:"description"`
	Type        PackageExportsFallbackType `json:"type"`
	Items       Items                      `json:"items"`
}

type Items struct {
	Type    TypeElement `json:"type"`
	Pattern string      `json:"pattern"`
}

type Private struct {
	Description string         `json:"description"`
	OneOf       []OneOfElement `json:"oneOf"`
}

type ExecutionEnv struct {
	Type                 TypeElement            `json:"type"`
	Properties           ExecutionEnvProperties `json:"properties"`
	AdditionalProperties bool                   `json:"additionalProperties"`
}

type ExecutionEnvProperties struct {
	NodeVersion Description `json:"nodeVersion"`
}

type PeerDependencyRules struct {
	Type                 TypeElement                   `json:"type"`
	Properties           PeerDependencyRulesProperties `json:"properties"`
	AdditionalProperties bool                          `json:"additionalProperties"`
}

type PeerDependencyRulesProperties struct {
	IgnoreMissing   CPUElement  `json:"ignoreMissing"`
	AllowedVersions Description `json:"allowedVersions"`
	AllowAny        CPUElement  `json:"allowAny"`
}

type UpdateConfig struct {
	Type                 TypeElement            `json:"type"`
	Properties           UpdateConfigProperties `json:"properties"`
	AdditionalProperties bool                   `json:"additionalProperties"`
}

type UpdateConfigProperties struct {
	IgnoreDependencies CPUElement `json:"ignoreDependencies"`
}

type ConfigPatternProperties struct {
	Empty             *LicenseValue `json:"^#.+$,omitempty"`
	PatternProperties *Class        `json:"^.+$,omitempty"`
}

type Class struct {
	Type                 TypeElement `json:"type"`
	Properties           Properties  `json:"properties"`
	AdditionalProperties bool        `json:"additionalProperties"`
}

type Properties struct {
	Dependencies         Author `json:"dependencies"`
	OptionalDependencies Author `json:"optionalDependencies"`
	PeerDependencies     Author `json:"peerDependencies"`
	PeerDependenciesMeta Author `json:"peerDependenciesMeta"`
}

type DevEngines struct {
	Description string               `json:"description"`
	Type        TypeElement          `json:"type"`
	Properties  DevEnginesProperties `json:"properties"`
}

type DevEnginesProperties struct {
	OS             LibcClass `json:"os"`
	CPU            LibcClass `json:"cpu"`
	Libc           LibcClass `json:"libc"`
	Runtime        LibcClass `json:"runtime"`
	PackageManager LibcClass `json:"packageManager"`
}

type LibcClass struct {
	OneOf       []CPUOneOf `json:"oneOf"`
	Description string     `json:"description"`
}

type CPUOneOf struct {
	Ref   *string                     `json:"$ref,omitempty"`
	Type  *PackageExportsFallbackType `json:"type,omitempty"`
	Items *Author                     `json:"items,omitempty"`
}

type Directories struct {
	Type       TypeElement           `json:"type"`
	Properties DirectoriesProperties `json:"properties"`
}

type DirectoriesProperties struct {
	Bin     Description  `json:"bin"`
	Doc     Description  `json:"doc"`
	Example Description  `json:"example"`
	LIB     Description  `json:"lib"`
	Man     Description  `json:"man"`
	Test    EngineStrict `json:"test"`
}

type Dist struct {
	Type       TypeElement    `json:"type"`
	Properties DistProperties `json:"properties"`
}

type DistProperties struct {
	Shasum  EngineStrict `json:"shasum"`
	Tarball EngineStrict `json:"tarball"`
}

type Engines struct {
	Type                 TypeElement       `json:"type"`
	Properties           EnginesProperties `json:"properties"`
	AdditionalProperties EngineStrict      `json:"additionalProperties"`
}

type EnginesProperties struct {
	Node EngineStrict `json:"node"`
}

type Esnext struct {
	Description          string           `json:"description"`
	Type                 []TypeElement    `json:"type"`
	Properties           EsnextProperties `json:"properties"`
	AdditionalProperties EngineStrict     `json:"additionalProperties"`
}

type EsnextProperties struct {
	Main    EngineStrict `json:"main"`
	Browser EngineStrict `json:"browser"`
}

type Exports struct {
	Description string         `json:"description"`
	OneOf       []ExportsOneOf `json:"oneOf"`
}

type ExportsOneOf struct {
	Ref                  *string                 `json:"$ref,omitempty"`
	Description          *string                 `json:"description,omitempty"`
	Type                 *TypeElement            `json:"type,omitempty"`
	Properties           *OneOfProperties        `json:"properties,omitempty"`
	PatternProperties    *OneOfPatternProperties `json:"patternProperties,omitempty"`
	AdditionalProperties *bool                   `json:"additionalProperties,omitempty"`
}

type OneOfPatternProperties struct {
	Empty LicenseValue `json:"^\\./.+"`
}

type OneOfProperties struct {
	Empty LicenseValue `json:"."`
}

type Funding struct {
	OneOf []FundingOneOf `json:"oneOf"`
}

type FundingOneOf struct {
	Ref         *string                     `json:"$ref,omitempty"`
	Type        *PackageExportsFallbackType `json:"type,omitempty"`
	Items       *PackageExportsEntry        `json:"items,omitempty"`
	MinItems    *int64                      `json:"minItems,omitempty"`
	UniqueItems *bool                       `json:"uniqueItems,omitempty"`
}

type Licenses struct {
	Description string                     `json:"description"`
	Type        PackageExportsFallbackType `json:"type"`
	Items       LicensesItems              `json:"items"`
}

type LicensesItems struct {
	Type       TypeElement     `json:"type"`
	Properties ItemsProperties `json:"properties"`
}

type ItemsProperties struct {
	Type Author `json:"type"`
	URL  Email  `json:"url"`
}

type Man struct {
	Type        []string     `json:"type"`
	Description string       `json:"description"`
	Items       EngineStrict `json:"items"`
}

type Name struct {
	Description string      `json:"description"`
	Type        TypeElement `json:"type"`
	MaxLength   int64       `json:"maxLength"`
	MinLength   int64       `json:"minLength"`
	Pattern     string      `json:"pattern"`
}

type PackageManager struct {
	Description string      `json:"description"`
	Type        TypeElement `json:"type"`
	Pattern     string      `json:"pattern"`
}

type PublishConfig struct {
	Type                 TypeElement             `json:"type"`
	Properties           PublishConfigProperties `json:"properties"`
	AdditionalProperties bool                    `json:"additionalProperties"`
}

type PublishConfigProperties struct {
	Access     OneOfElement `json:"access"`
	Tag        EngineStrict `json:"tag"`
	Registry   Email        `json:"registry"`
	Provenance EngineStrict `json:"provenance"`
}

type Repository struct {
	Description string               `json:"description"`
	Type        []TypeElement        `json:"type"`
	Properties  RepositoryProperties `json:"properties"`
}

type RepositoryProperties struct {
	Type      EngineStrict `json:"type"`
	URL       EngineStrict `json:"url"`
	Directory EngineStrict `json:"directory"`
}

type ScriptsClass struct {
	Description          string                      `json:"description"`
	Type                 TypeElement                 `json:"type"`
	Properties           ScriptsProperties           `json:"properties"`
	AdditionalProperties ScriptsAdditionalProperties `json:"additionalProperties"`
}

type ScriptsAdditionalProperties struct {
	Type                       TypeElement `json:"type"`
	TsType                     string      `json:"tsType"`
	XIntellijLanguageInjection string      `json:"x-intellij-language-injection"`
}

type ScriptsProperties struct {
	Lint           Description `json:"lint"`
	Prepublish     Description `json:"prepublish"`
	Prepare        Description `json:"prepare"`
	PrepublishOnly Description `json:"prepublishOnly"`
	Prepack        Description `json:"prepack"`
	Postpack       Description `json:"postpack"`
	Publish        Description `json:"publish"`
	Postpublish    Author      `json:"postpublish"`
	Preinstall     Description `json:"preinstall"`
	Install        Author      `json:"install"`
	Postinstall    Author      `json:"postinstall"`
	Preuninstall   Author      `json:"preuninstall"`
	Uninstall      Author      `json:"uninstall"`
	Postuninstall  Description `json:"postuninstall"`
	Preversion     Author      `json:"preversion"`
	Version        Author      `json:"version"`
	Postversion    Description `json:"postversion"`
	Pretest        Author      `json:"pretest"`
	Test           Author      `json:"test"`
	Posttest       Author      `json:"posttest"`
	Prestop        Author      `json:"prestop"`
	Stop           Author      `json:"stop"`
	Poststop       Author      `json:"poststop"`
	Prestart       Author      `json:"prestart"`
	Start          Author      `json:"start"`
	Poststart      Author      `json:"poststart"`
	Prerestart     Author      `json:"prerestart"`
	Restart        Author      `json:"restart"`
	Postrestart    Author      `json:"postrestart"`
	Serve          Description `json:"serve"`
}

type TypesVersions struct {
	Description          string                            `json:"description"`
	Type                 TypeElement                       `json:"type"`
	AdditionalProperties TypesVersionsAdditionalProperties `json:"additionalProperties"`
}

type TypesVersionsAdditionalProperties struct {
	Description          string                                `json:"description"`
	Type                 TypeElement                           `json:"type"`
	Properties           FluffyProperties                      `json:"properties"`
	PatternProperties    AdditionalPropertiesPatternProperties `json:"patternProperties"`
	AdditionalProperties bool                                  `json:"additionalProperties"`
}

type AdditionalPropertiesPatternProperties struct {
	PatternProperties CPUElement      `json:"^[^*]+$"`
	Empty             IgnoreCvesClass `json:"^[^*]*\\*[^*]*$"`
}

type FluffyProperties struct {
	Empty IgnoreCvesClass `json:"*"`
}

type Volta struct {
	Description       string                 `json:"description"`
	Type              TypeElement            `json:"type"`
	Properties        VoltaProperties        `json:"properties"`
	PatternProperties VoltaPatternProperties `json:"patternProperties"`
}

type VoltaPatternProperties struct {
	NodeNpmPnpmYarn EngineStrict `json:"(node|npm|pnpm|yarn)"`
}

type VoltaProperties struct {
	Extends Description `json:"extends"`
}

type Workspaces struct {
	Description string       `json:"description"`
	AnyOf       []CPUElement `json:"anyOf"`
}

type TypeElement string

const (
	Boolean      TypeElement = "boolean"
	PurpleObject TypeElement = "object"
	String       TypeElement = "string"
)

type Ref string

const (
	DefinitionsLicense                       Ref = "#/definitions/license"
	DefinitionsPackageExportsEntryOrFallback Ref = "#/definitions/packageExportsEntryOrFallback"
	DefinitionsPackageImportsEntryOrFallback Ref = "#/definitions/packageImportsEntryOrFallback"
)

type PackageExportsFallbackType string

const (
	Array        PackageExportsFallbackType = "array"
	FluffyObject PackageExportsFallbackType = "object"
)
