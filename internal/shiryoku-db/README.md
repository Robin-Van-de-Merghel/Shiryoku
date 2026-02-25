# Shiryoku-db

This layer interacts directly with the database*s*.

## Architecture

Multiple kinds of database can be used: OpenSearch/ElastricSearch, MySQL, etc. The idea is to provide an interface for the other layers.

It is better if we can use a similar interface between all engines, as it will make it simpler to work with.

### Nmap DB

```mermaid
classDiagram
    %% A scan discovered services
    %% NmapScan --> ScanResult
    %% ScanResult --> Service
    %% ScanResult --> NmapHost
    %% ScanResult --> NmapScriptResult

    class NmapScan{
        +UUID ScanID
        +time ScanStart
        +string ScanArgs
        +string NmapVersion
        +time CreatedAt
    }
    
    %% A host discovered during scans
    %% Unicity : Host (IP address, or first address if multiple)
    %% Can appear in multiple scans
    class NmapHost {
        +UUID HostID
        +string Host
        +[]string Addresses
        +[]string Hostnames
        +string HostStatus
        +string OSName
        +int OSAccuracy
        +string Comment
        +time CreatedAt
    }
    
    %% Usable by other scripts than nmap (nuclei, etc.)
    %% Unicity : (ServiceName + Product + Version + ServiceExtraInfo + Protocol + ServiceTunnel)
    %% Or else reuse existing one
    class Service {
        +UUID ServiceID
        +string ServiceName
        +string ServiceProduct
        +string ServiceVersion
        +string ServiceExtraInfo
        +string Protocol
        +string ServiceTunnel
        +time CreatedAt
    }

    %% A scan discovered "ScanResult"s
    %% One result = port + scan + host + service
    %% Composite key: (ScanID, HostID, ServiceID, Port)
    %% Represents: In this scan, this host had this service on this port
    class ScanResult {
        +UUID ScanResultID
        +UUID ScanID
        +UUID HostID
        +UUID ServiceID
        +uint16 Port
        +string PortState
        +time CreatedAt
    }
    
    %% NSE script results for a specific scan result
    %% (+ ScriptResults)
    class NmapScriptResult {
        +UUID NmapScriptResultID
        +UUID ScanResultID
        +string ScriptID
        +string ScriptOutput
        +time CreatedAt
    }

    NmapScan "1" --> "*" ScanResult: has
    NmapHost "1" --> "*" ScanResult: appears_in
    Service "1" --> "*" ScanResult: referenced_by
    ScanResult "1" --> "*" NmapScriptResult: has
```
