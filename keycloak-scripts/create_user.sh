#!/bin/bash
# Wait for Keycloak to start and then create a new user
echo "Waiting for Keycloak to start..."

# Give Keycloak some time to fully initialize
sleep 30

echo "Creating user admin with password..."

# Create user via Keycloak's REST API
curl -X POST "http://localhost:8080/admin/realms/peitho/users" \
-H "Content-Type: application/json" \
-d '{
      "username": "admin",
      "enabled": true,
      "firstName": "Admin",
      "lastName": "User",
      "email": "admin@peitho.com",
      "credentials": [
        {
          "type": "password",
          "value": "admin",
          "temporary": false
        }
      ],
      "roles": ["admin"]
    }'

echo "User created!"

