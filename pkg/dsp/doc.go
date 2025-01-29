/*
Package dsp provides an interface to interact with the Database Solutions Portal (DSP) API.

This package provides functionality to look up dataserver resources and handle their associated
metadata including support details, version information, and location data.

Example usage:

    dspClient := dsp.New()
    response, err := dspClient.LookupDataserverResource(context.Background(), "resourceName")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Resource details: %+v\n", response.Data)
*/
package dsp