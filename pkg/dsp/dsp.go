package dsp

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/citi/dsp/internal/auditlogger"
)

// getAuditValues safely extracts audit values from context
func getAuditValues(ctx context.Context) (auditlogger.Values, error) {
    if ctx == nil {
        return auditlogger.Values{}, fmt.Errorf("nil context")
    }
    
    val := ctx.Value("auditValues")
    if val == nil {
        return auditlogger.Values{}, fmt.Errorf("no audit values in context")
    }
    
    auditValues, ok := val.(auditlogger.Values)
    if !ok {
        return auditlogger.Values{}, fmt.Errorf("invalid audit values type in context")
    }
    
    return auditValues, nil
}

// LookupDataserverResource implements the DspInterface method
func (dspStruct *DspStruct) LookupDataserverResource(requestCtx context.Context, resourceName string) (DspResponse, error) {
    var dspResponse DspResponse
    
    // Safely get audit values
    auditValues, err := getAuditValues(requestCtx)
    if err != nil {
        return dspResponse, fmt.Errorf("failed to get audit values: %w", err)
    }

    var client *http.Client
    dspUrl := "https://databasesolutionsportal.citigroup.net/webservices/api/dsp/dataserver?resourceName="
    dataserverUrl := dspUrl + resourceName

    // set up the client
    transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
    client = &http.Client{Transport: transport}
    caCert := os.Getenv("CAROOT")
    ca, err := os.ReadFile(caCert)
    if err != nil {
        if auditValues.Logger != nil {
            auditValues.Logger.Error("Unable to set up secure client for DSP",
                "error", err.Error(),
                "uuid", auditValues.Uuid,
                "soeid", auditValues.Soeid)
        }
        return dspResponse, err
    }

    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(ca)

    req, err := http.NewRequest("GET", dataserverUrl, nil)
    if err != nil {
        if auditValues.Logger != nil {
            auditValues.Logger.Error("Failure creating dataserver request",
                "error", err.Error(),
                "uuid", auditValues.Uuid,
                "soeid", auditValues.Soeid)
        }
        return dspResponse, err
    }

    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
        if auditValues.Logger != nil {
            auditValues.Logger.Error("Failure on dataserver call",
                "error", err.Error(),
                "uuid", auditValues.Uuid,
                "soeid", auditValues.Soeid)
        }
        return dspResponse, err
    }

    defer res.Body.Close()
    var data DspDedicatedHost

    err = json.NewDecoder(res.Body).Decode(&data)
    if err != nil {
        if auditValues.Logger != nil {
            auditValues.Logger.Error("Failure decoding response",
                "error", err.Error(),
                "uuid", auditValues.Uuid,
                "soeid", auditValues.Soeid)
        }
        return dspResponse, err
    }

    dspResponse.Data = data
    return dspResponse, nil
}