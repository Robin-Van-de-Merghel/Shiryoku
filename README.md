# Shiryoku

Shiryoku (*vision* in japanese) aims at having a wide vision of internet resources.

> [!NOTE]
> It is in development.

The idea is near what [IVRE](https://github.com/ivre/ivre): scan a range of IPs, domains, etc., and aggregate data to a dashboard.

> [!IMPORTANT]
> This is a lab-only project. Please use this tool cautiously on systems you own.

# Roadmap

What I want to achieve:

- [x] IP/Host list
- [x] Scan list
- [ ] Be able to "program" scans
- [ ] Import from : 
    - [x] `nmap` (XML only for now)
    - [ ] `nuclei`
    - [ ] `ffuf`
    - [ ] `wpscan`
    - [ ] `masscan`
    - [ ] `httpx`
    - [ ] More?
- [ ] User management (none for now)
- [ ] Agents (that collect data)
- [ ] Tasks Queue (targets to scan)

# Documentations

Some documentations will be written as time goes...:

- [`shiryoku-routers`](./internal/shiryoku-routers/README.md): handles requests
- [`shiryoku-logic`](./internal/shiryoku-logic/README.md): contains the business logic
- [`shiryoku-db`](./internal/shiryoku-db/README.md): handle db interactions
- [`shiryoku-core`](./internal/shiryoku-core/README.md): contains configurations as well as models

> [!NOTE]
> This design is (vastly) adapted from [diracx](https://github.com/diracgrid/diracx), a project I worked for.
