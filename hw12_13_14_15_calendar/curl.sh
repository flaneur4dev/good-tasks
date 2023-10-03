curl -L -X POST -v '127.0.0.1:3000/event' \
-H 'Content-Type: application/json' \
--data-raw '{
    "id": "0c2e2081",
    "title": "event 1",
    "description": "Station National Rail",
    "owner_id": "230c866",
    "start_date": "2021-09-01T04:05:00Z",
    "fin_date": "2021-09-01T14:00:00Z",
    "notification_time": "2021-09-01T13:00:00Z"
}'

# curl -L -X GET 'localhost:3000/event?date=2019-01-08 04:05:06&period=day'

# curl -L -X PUT -v '127.0.0.1:3000/event' \
# -H 'Content-Type: application/json' \
# --data-raw '{
#     "id": "0c2e2081",
#     "title": "event 1_1",
#     "description": "",
#     "owner_id": "230c866-1",
#     "start_date": "2019-02-01T04:05:06Z",
#     "fin_date": "2019-02-09T14:35:00Z",
#     "notification_time": "2019-02-09T10:00:00Z"
# }'

# curl -L -X DELETE -v '127.0.0.1:3000/event?id=0c2e2081'
