# Volunteer-Savvy.com Backend Design

Volunteer-Savvy is a SAAS approach to 'employee' scheduling, targeting volunteer organizations primarily.

Each Organization can have many Sites and many Volunteers, and those Volunteers may log hours worked at each Site. Managers may then approve the work logs for bookkeeping purposes.

## Users

Users belong to an Organization.
They have an email and a password used to log in.
They can subscribe to Sites to get notifications when it changes (hours, days open, location).
They can join "Teams", and opt-in to notifications for each of those Teams. If the Team name corresponds to a Site Feature, and the Site is updated, and the Volunteer has opted-in, then they will be notified of the change.

## Sites

Sites belong to an Organization.
They have location data, including address, lat/lon, and Google Place ID.
They have a default schedule, specifying which days they are usually open, and for what hours.
They can have 'calendar overrides', which are calendar entries specifying different hours or open status for a single specific day.

# Microservices

## Authentication

Auth will be handled via a standards-compliant OAuth2 server. It will handle password 
credential grants for users, and client credential grants for other microservices. 
We will use JWT where possible for the auth token.

The JWT will indicate which organizations the user is a member of. The APIs will 
require the caller to specify which organization they are operating on via header. 

## Site Management

The Site API service will handle CRUD for Sites. 

    # Create/Read/Update/Delete Sites
    GET  /sites/
    POST /sites/
    GET  /sites/{site-slug}
    PUT  /sites/{site-slug}
    DELETE /sites/{site-slug}
    
    # Add/remove features from Sites
    PUT    /sites/{site-slug}/feature/{feature-id}
    DELETE /sites/{site-slug}/feature/{feature-id}

## Volunteer Management

The Volunteer API service will handle CRUD for Users. Users will be 
identified by a GUID randomly assigned on user creation.

    # Create/Read/Update/Delete Users
    GET /users/
    GET /users/{user-guid}
    POST /users/
    PUT  /users/{user-guid}
    DELETE /users/{user-guid}
    
    # Add/remove roles from Users
    PUT /users/{user-guid}/role/{role-id}
    DELETE /users/{user-guid}/role/{role-id}
    
## Suggestions

The Suggestions API service will allow users or 
anonymous visitors to leave feedback for the volunteers.

    # Create/Read/Update/Delete Suggestions
    GET /suggestions/
    GET /suggestions/{suggestion-id}
    PUT /suggestions/{suggestion-id}
    POST /suggestions/
    DELETE /suggestions/{suggestion-id}
    
## Subscriptions

The Subscriptions API service will allow Users to request push 
notifications or email notifications on certain triggers.

    # List My Subscriptions
    GET /subscriptions/{user-guid}
    
    # On Site changes. Takes a JSON body that specifies message channel: push/email
    POST /subscriptions/site/{site-slug}
    DELETE /subscriptions/site/{site-slug}
    # Multi-select. Takes a JSON body that specifies all notification subscriptions + their channels.
    POST /subscriptions/site/ 
    
    # On new Suggestions
    POST /subscriptions/suggestions
    DELETE /subscriptions/suggestions
    
    # On new User creation
    POST /subscription/users/
    DELETE /subscriptions/users/
    
    # Device Registration for Push Notifications
    POST /subscription/register-device/
    
