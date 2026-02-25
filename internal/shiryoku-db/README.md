# Shiryoku-db

This layer interacts directly with the database*s*.

## Architecture

Multiple kinds of database can be used: OpenSearch/ElastricSearch, MySQL, etc. The idea is to provide an interface for the other layers.

> [!NOTE]
> We will mostly use Postgres databases.

It is better if we can use a similar interface between all engines, as it will make it simpler to work with.

To do that, we will use `SearchParams` (taken from [diracx](https://github.com/DIRACGrid/diracx/blob/main/diracx-core/src/diracx/core/models/search.py)). It provides a quite solid interface to the user as well as developper to build complex requests. For example:

```json
{
    "search": [
        {
            "parameter": "host",
            "operator": "eq",
            "value": "1.1.1.1"
        }
    ]
}
```

## Nmap storage

Nmap storage is divided in two parts: main storage, and the dashboard's. The first one is the result of every nmap scans, and the other one is dedicated to displaying scans on a dashboard (views calculated from the whole data).

```mermaid
classDiagram
    %% Main entities
    class NmapScan{
        +UUID ScanID
        +time ScanStart
        +string ScanArgs
        +string NmapVersion
        +time CreatedAt
    }
    
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

    class ScanResult {
        +UUID ScanResultID
        +UUID ScanID
        +UUID HostID
        +UUID ServiceID
        +uint16 Port
        +string PortState
        +time CreatedAt
    }
    
    class NmapScriptResult {
        +UUID NmapScriptResultID
        +UUID ScanResultID
        +string ScriptID
        +string ScriptOutput
        +time CreatedAt
    }

    %% Dashboard-specific view
    class WidgetDashboardScan {
        +string ScanID
        +string HostID
        +time ScanStart
        +string Host
        +int PortNumber
        +[]int Ports
        +[]string HostNames
    }

    %% Relationships
    NmapScan "1" --> "*" ScanResult: has
    NmapHost "1" --> "*" ScanResult: appears_in
    Service "1" --> "*" ScanResult: referenced_by
    ScanResult "1" --> "*" NmapScriptResult: has
    ScanResult "*" --> "*" WidgetDashboardScan: being_used_by
    NmapHost "*" --> "*" WidgetDashboardScan: being_used_by
```
