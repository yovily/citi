package dsp

import "context"

// DspInterface defines the interface for DSP API operations
type DspInterface interface {
    // LookupDataserverResource looks up a dataserver resource from the DSP API
    LookupDataserverResource(requestCtx context.Context, resourceName string) (DspResponse, error)
}

// New creates a new instance of DspInterface
func New() DspInterface {
    return &DspStruct{}
}