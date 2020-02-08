# Volunteer Savvy Backend Test Plan

## Unit Tests

## Integration Tests

1. DescribeOrganization

    a. Accessible to public
    
    b. Schema includes org name, emergency contact info, geographical location

2. CreateOrganization
    
    a. Not accessible to public
    
    b. Only accessible to a user whose JWT claims include `SUPERADMIN`
    
    c. Input must include Org name
    
    d. Input may include emergency contact info, location data
    
    e. Response must include generated Org ID (should be the whole DescribeOrganization view)
    
3. UpdateOrganization

    a. Not accessible to public
    
    b. Accessible to a user whose JWT claims include `SUPERADMIN`
    
    c. Accessible to a user whose JWT claims include `ORG:{orgId}:ADMIN`
    
    d. Input may include new emergency contact info and/or location data
    
    e. Response should be the DescribeOrganization view

4. DestroyOrganization

    a. Not accessible to public
    
    b. Accessible only to a user whose JWT claims include `SUPERADMIN`

5. ListOrganizations

    a. Accessible to public
    
    b. Top level array, each member of which conforms to the DescribeOrganization schema
    
6. DescribeUser

    a. Accessible to public
    
    b. View changes based on login status:
      * Add work logs if the JWT claim includes `ORG:{orgId}:ADMIN` or `USER:{userId}:1`
      * Add signups if the JWT claim includes `ORG:{orgId}:ADMIN` or `USER:{userId}:1`
      
7. CreateUser
8. UpdateUser
9. DestroyUser
10. ListUsers 
    