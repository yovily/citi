package dsp

// DspTimeZone represents timezone information
type DspTimeZone struct {
    TimeZoneId   int    `json:"timeZoneId"`
    TimeZoneDesc string `json:"timeZoneDesc"`
}

// DspDBA represents DBA information
type DspDBA struct {
    RITSID string `json:"RITSID"`
    SOEID  string `json:"SOEID"`
    GEID   string `json:"GEID"`
    FName  string `json:"FName"`
    LName  string `json:"LName"`
    Email  string `json:"Email"`
}

// DspSupportDetails represents support information
type DspSupportDetails struct {
    PrimaryDBA        DspDBA  `json:"primaryDBA"`
    SecondaryDBA     DspDBA  `json:"secondaryDBA"`
    RotationId       string  `json:"rotationId"`
    RotationDesc     string  `json:"rotationDesc"`
    RotationEmail    string  `json:"rotationEmail"`
    RotationManager  DspDBA  `json:"rotationManager"`
    DBAGroupName     string  `json:"DBAGroupName"`
    DBAGroupEmail    string  `json:"DBAGroupEmail"`
    DBAGroupManager  DspDBA  `json:"DBAGroupManager"`
    DBASectorName    string  `json:"DBASectorName"`
    DBASectorDesc    string  `json:"DBASectorDesc"`
    DBASectorHead    DspDBA  `json:"DBASectorHead"`
}

// DspVersion represents version information
type DspVersion struct {
    RdbmsName       string `json:"RdbmsName"`
    Version         string `json:"version"`
    SubVersion      string `json:"subVersion"`
    CTCProductId    uint32 `json:"CTCProductId"`
    ProductName     string `json:"productName"`
    CTCVersionId    uint32 `json:"CTCVersionId"`
    CTCVersionName  string `json:"CTCVersionName"`
}

// DspCsidMapping represents CSID mapping information
type DspCsidMapping struct {
    CSIAppId         uint32 `json:"CSIAppId"`
    CSIAppName       string `json:"CSIAppName"`
    AppInstanceId    string `json:"appInstanceId"`
    AppInstanceName  string `json:"appInstanceName"`
}

// DspKeyValueAttribute represents key-value attributes
type DspKeyValueAttribute struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

// DspRelatedResources represents related resources
type DspRelatedResources struct {
    RelationType string `json:"relationType"`
    ResourceId   string `json:"resourceId"`
    ResourceName string `json:"resourceName"`
}

// DspLocation represents location information
type DspLocation struct {
    Region      string `json:"region"`
    Country     string `json:"country"`
    DataCenter  string `json:"dataCenter"`
}

// DspDedicatedHost represents dedicated host information
type DspDedicatedHost struct {
    ResourceId          uint32                 `json:"resourceId"`
    ResourceName        string                 `json:"resourceName"`
    EnvironmentType     string                 `json:"environmentType"`
    PortNumber         int                    `json:"portNumber"`
    DataServerType      string                 `json:"dataServerType"`
    ProvisionPlatform   string                 `json:"provisionPlatform"`
    ConnectionString    string                 `json:"connectionString"`
    CreateDate         string                 `json:"createDate"`
    ModifyDate         string                 `json:"modifyDate"`
    UpgradeDate        string                 `json:"upgradeDate"`
    LiveDate           string                 `json:"liveDate"`
    TimeZone           DspTimeZone            `json:"timeZone"`
    SupportDetails     DspSupportDetails      `json:"supportDetails"`
    ModifiedBy         DspDBA                 `json:"modifiedBy"`
    Version            DspVersion             `json:"version"`
    CSIMappings        []DspCsidMapping       `json:"csiMappings"`
    KeyValueAttributes []DspKeyValueAttribute `json:"keyValueAttributes"`
    RelatedResource    DspRelatedResources    `json:"relatedResources"`
    Location           []DspLocation          `json:"location"`
}

// DspResponse represents the response from DSP API
type DspResponse struct {
    Data DspDedicatedHost `json:"data"`
}

// DspStruct implements the DspInterface
type DspStruct struct{}