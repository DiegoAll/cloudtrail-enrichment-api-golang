# Unit test Report


## Package Coverage

    root@h0stn4m3:/home/diegoall/Projects/cloudtrail-enrichment-api-golang# go test ./... -coverprofile=coverage.out
    ok      cloudtrail-enrichment-api-golang/cmd/api        0.003s  coverage: 0.0% of statements
    ok      cloudtrail-enrichment-api-golang/cmd/api/controllers    0.004s  coverage: 64.9% of statements
    ?       cloudtrail-enrichment-api-golang/models [no test files]
            cloudtrail-enrichment-api-golang/internal/config                coverage: 0.0% of statements
    ok      cloudtrail-enrichment-api-golang/database/mongo 1.217s  coverage: 87.3% of statements
    ok      cloudtrail-enrichment-api-golang/database/postgresql    0.004s  coverage: 88.6% of statements
    ok      cloudtrail-enrichment-api-golang/internal/middleware    0.004s  coverage: 100.0% of statements
    ok      cloudtrail-enrichment-api-golang/internal/pkg/logger    0.002s  coverage: 100.0% of statements
    ok      cloudtrail-enrichment-api-golang/internal/pkg/scopes    0.002s  coverage: 100.0% of statements
    ok      cloudtrail-enrichment-api-golang/internal/pkg/token     0.006s  coverage: 69.2% of statements
    ok      cloudtrail-enrichment-api-golang/internal/pkg/utils     0.004s  coverage: 72.5% of statements
    ok      cloudtrail-enrichment-api-golang/internal/repository    0.006s  coverage: 100.0% of statements
    ok      cloudtrail-enrichment-api-golang/services       1.216s  coverage: 59.2% of statements


##  Total Coverage  

    root@h0stn4m3:/home/diegoall/Projects/cloudtrail-enrichment-api-golang# go tool cover -func=coverage.out | grep total
    total:                                                                                  (statements)                    59.7%


## Findings

- Correct the data types for the time parameters in the config package so that they are standardized and do not affect the JWT token.

- It is possible to decouple the routes into a separate router package; this should be analyzed.