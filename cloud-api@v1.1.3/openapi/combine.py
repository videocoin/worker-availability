import json
# result = ''

with open('./streams/v1/streams_service.swagger.json') as f:
    result = json.load(f)

with open('./profiles/v1/profiles_service.swagger.json') as f:
    profiles = json.load(f)
    result['paths'].update(profiles['paths'])
    result['definitions'].update(profiles['definitions'])
result['host'] = 'snb.videocoin.network'


# print(result)
with open('videocoin_service.json', 'w') as f:
    json.dump(result, f, indent=4)
