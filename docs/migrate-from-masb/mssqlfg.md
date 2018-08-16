# Migrate [Azure SQL Database Failover Group Service](https://github.com/Azure/meta-azure-service-broker/blob/master/docs/azure-sql-db-failover-group.md) Instances From MASB To OSBA

## Steps

### Create azure-sql-12-0-dr-dbms-pair-registered service instance by OSBA

OSBA doesn't support to provide existing servers' credentials from manifest and create database on them. You need to create a **azure-sql-12-0-dr-dbms-pair-registered** service instance for the two servers involved:

```
cf create-service azure-sql-12-0-dr-dbms-pair-registered dbms <sql-server-pair-instance-name> -c '{
  "primaryResourceGroup":"<primary-group-name>",
  "primaryServer":"<primary-server-name>",
  "primaryLocation":"<primary-server-location>",
  "primaryAdministratorLogin":"<primary-login>",
  "primaryAdministratorLoginPassword":"<primary-login-password>",
  "secondaryResourceGroup":"<secondary-group-name>",
  "secondaryServer":"<secondary-server-name>",
  "secondaryLocation":"<secondary-server-location>",
  "secondaryAdministratorLogin":"<secondary-login>",
  "secondaryAdministratorLoginPassword":"<secondary-login-password>",
  "alias":"<sql-server-pair-instance-name>"
}'
```

### Create azure-sql-12-0-dr-database-pair-from-existing service instance by OSBA

In MASB, the service for failover group creates the secondary database and failover group based on an existing primary database. In OSBA, the primary database, the secondary primary, and the failover group are managed by one service instance. To migrate, you need to create **azure-sql-12-0-dr-database-pair-from-existing** service instance.

```
cf env <your-app-name>
```

Find the database and failover group name in the field `VCAP_SERVICES` - `azure-sqldb-failover-group` - `sqldbName` / `sqlServerName`. (You know, the failover group is transparent to applications. So, its name is `sqlServerName` in the credentials.)

Then, you can run the following command to create a instance of **azure-sql-12-0-dr-database-pair-from-existing**:

```
cf create-service azure-sql-12-0-dr-database-pair-from-existing <plan-name> <sqldb-pair-instance-name> -c '{
  "parentAlias":"<the-alias-above>",
  "failoverGroup": "<sql-server-name>",
  "database":"<sqldb-name>"
}'
```

***Note***: OSBA provides plans by tier categories. For example, compared to that MASB provides plans `StandardS0` - `StandardS12`, OSBA provides only a plan `standard` for all the standard tiers. You should choose the right category. It is important. Though provisioning wouldn't change the tier, you wouldn't be able to update the tier as OSBA doesn't support changing plan for now. The plan name of existing service instances can be checked by `cf services`.

### Duplicate your application and update to adapt the SQL credentials delivered by OSBA

You should check the SQL credential differences between [OSBA](../modules/mssqlfg.md#credentials-1) and [MASB](https://github.com/Azure/meta-azure-service-broker/blob/master/docs/azure-sql-db-failover-group.md#format-of-credentials). Update how your application utilizes the credentials. Then `cf push` your updated application with another name and another route.

### Bind azure-sql-12-0-dr-database-pair-from-existing to your application

After successfully creating the service instance, bind it to your application:

```
cf bind-service <your-app-new-name> <sqldb-pair-instance-name>
```

Restage the application:

```
cf restage <your-app-new-name>
```

After you test the application and it works well, you can switch your application domain to the new route in your DNS.

### MASB service instance safely clean up

Unbind the old SQL database instance:

```
cf unbind-service <your-app-name> <masb-sqldb-pair-instance-name>
```

Purge the service instance (**IMPORTANT**):

```
cf purge-service-instance <masb-sqldb-instance-name>
```

Don't use `cf delete-service` here. Or, your database would be deleted in Azure!
