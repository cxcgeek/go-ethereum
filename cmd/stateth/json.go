package main

var (
	jsonDashboard = `{
  "annotations": {
    "list": [
      {
        "$$hashKey": "object:448",
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "gnetId": null,
  "graphTooltip": 1,
  "id": 5,
  "iteration": 1526648602647,
  "links": [],
  "panels": [
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 40,
      "panels": [],
      "title": "LocalStore",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 0,
        "y": 1
      },
      "id": 42,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.localstore.get.cachehit.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LocalStore get cachehit",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 12,
        "y": 1
      },
      "id": 43,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.localstore.get.cachemiss.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LocalStore get cachemiss",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 0,
        "y": 7
      },
      "id": 44,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.localstore.getorcreaterequest.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "Total LocalStore.GetOrCreateRequest",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 12,
        "y": 7
      },
      "id": 47,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.localstore.getorcreaterequest.errfetching.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LocalStore GetOrCreateRequest ErrFetching",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 0,
        "y": 13
      },
      "id": 45,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.localstore.getorcreaterequest.hit.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LocalStore.GetOrCreateRequest hit",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 12,
        "y": 13
      },
      "id": 49,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.localstore.getorcreaterequest.miss.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LocalStore GetOrCreateRequest miss",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 0,
        "y": 19
      },
      "id": 48,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.localstore.get.error.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LocalStore get error",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 12,
        "y": 19
      },
      "id": 46,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.localstore.get.errfetching.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LocalStore get ErrFetching",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 25
      },
      "id": 27,
      "panels": [],
      "title": "LDBStore",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 0,
        "y": 26
      },
      "id": 29,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.ldbstore.get.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LDBStore get",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 12,
        "y": 26
      },
      "id": 30,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.ldbstore.put.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LDBStore put",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 0,
        "y": 32
      },
      "id": 31,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.ldbstore.synciterator.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LDBStore SyncIterator",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 12,
        "y": 32
      },
      "id": 32,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.ldbstore.synciterator.seek.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LDBStore SyncIterator Seek/Next",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 38
      },
      "id": 34,
      "panels": [],
      "title": "LDBDatabase",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 0,
        "y": 39
      },
      "id": 36,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.ldbdatabase.get.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LDBDatabase get",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 12,
        "y": 39
      },
      "id": 37,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.ldbdatabase.write.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LDBDatabase write",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "metrics",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 12,
        "x": 0,
        "y": 45
      },
      "id": 38,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "rightSide": true,
        "show": true,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "alias": "$tag_host",
          "groupBy": [
            {
              "params": [
                "$myinterval"
              ],
              "type": "time"
            },
            {
              "params": [
                "host"
              ],
              "type": "tag"
            },
            {
              "params": [
                "0"
              ],
              "type": "fill"
            }
          ],
          "measurement": "swarm.ldbdatabase.newiterator.count",
          "orderByTime": "ASC",
          "policy": "default",
          "refId": "A",
          "resultFormat": "time_series",
          "select": [
            [
              {
                "params": [
                  "value"
                ],
                "type": "field"
              },
              {
                "params": [],
                "type": "sum"
              }
            ]
          ],
          "tags": [
            {
              "key": "host",
              "operator": "=~",
              "value": "/^$host$/"
            }
          ]
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "LDBDatabase NewIterator",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ]
    }
  ],
  "refresh": "10s",
  "schemaVersion": 16,
  "style": "dark",
  "tags": [],
  "templating": {
    "list": [
      {
        "auto": false,
        "auto_count": 30,
        "auto_min": "10s",
        "current": {
          "text": "10s",
          "value": "10s"
        },
        "hide": 0,
        "label": "resolution",
        "name": "myinterval",
        "options": [
          {
            "selected": false,
            "text": "5s",
            "value": "5s"
          },
          {
            "selected": true,
            "text": "10s",
            "value": "10s"
          },
          {
            "selected": false,
            "text": "30s",
            "value": "30s"
          },
          {
            "selected": false,
            "text": "100s",
            "value": "100s"
          }
        ],
        "query": "5s,10s,30s,100s",
        "refresh": 2,
        "type": "interval"
      },
      {
        "allValue": null,
        "current": {
          "text": "swarm_30399 + swarm_30400 + swarm_30401",
          "value": [
            "swarm_30399",
            "swarm_30400",
            "swarm_30401"
          ]
        },
        "datasource": "metrics",
        "hide": 0,
        "includeAll": true,
        "label": null,
        "multi": true,
        "name": "host",
        "options": [],
        "query": "SHOW TAG VALUES WITH KEY = \"host\"",
        "refresh": 1,
        "regex": "",
        "sort": 1,
        "tagValuesQuery": "",
        "tags": [],
        "tagsQuery": "swarm.http.request.GET.time.span",
        "type": "query",
        "useTags": false
      }
    ]
  },
  "time": {
    "from": "now-15m",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": [
      "5s",
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ],
    "time_options": [
      "5m",
      "15m",
      "1h",
      "6h",
      "12h",
      "24h",
      "2d",
      "7d",
      "30d"
    ]
  },
  "timezone": "",
  "title": "LDBStore and LDBDatabase",
  "uid": "zS6beG7iz",
  "version": 26
}
`
)
