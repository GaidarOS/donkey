# Donkey

Image storage and manipulation SaaS

## Why the name Donkey

Donkeys are pretty fantastic creatures! Here are a few reasons why:

Hardworking: Donkeys are known for their strong work ethic. They have been used for centuries as working animals, assisting with tasks like carrying loads and plowing fields.

Resilient: Donkeys are sturdy and resilient, able to thrive in challenging environments. Their adaptability makes them valuable in various climates and terrains.

Intelligent: Donkeys are intelligent animals. They can be trained easily and are known for their problem-solving abilities.

Affectionate: Donkeys can form strong bonds with their human caretakers. They are often affectionate and enjoy social interactions.

Surefooted: Donkeys are surefooted and have a good sense of balance. This makes them well-suited for navigating rough or hilly terrains.

Low Maintenance: Compared to some other animals, donkeys require relatively low maintenance. They can thrive on simple diets and are generally hardy.

Companionship: Donkeys are often kept as companions for other animals, providing a calming and protective presence. They are known to be good friends to horses and other livestock.

## Endpoints

| Path                  |     Method      | Header | Description        | Notes        | Status  |
| --------------------- | :-------------: | ------ | ------------------ | ------------ | ------- |
| `api/v1/login`        |      POST       | N/A    | Expects Basic Auth | Not Used atm | 200/401 |
| `api/v1/download/*`   |       GET       | Token  |                    |              |         |
| `api/v1/list/*`       |       GET       | Token  |                    |              |         |
| `api/v1/upload/*`     |      POST       | Token  |                    |              |         |
| `api/v1/delete/*`     |     DELETE      | Token  |                    |              |         |
| `api/v1/user/*`       |       GET       | Token  |                    |              |         |
| `api/v1/index/*`      |      POST       | Token  |                    |              |         |
| `api/v1/admin/users`  |       GET       | Token  |                    |              |         |
| `api/v1/admin/user`   | POST/PUT/DELETE | Token  |                    |              |         |
| `api/v1/admin/config` |    GET/POST     | Token  |                    |              |         |