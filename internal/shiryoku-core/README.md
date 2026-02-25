# Shiryoku-core

The core of the programs contains configurations, preferences, as well as models. It will be used by all the other parts of Shiryoku.

Its main goal is to hide from the workers as well as routers how it handles database requests orchestration (e.g. "query all scans, then merge with hosts").
