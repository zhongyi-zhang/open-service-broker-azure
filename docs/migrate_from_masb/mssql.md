# Migrate [Azure SQL Database](https://azure.microsoft.com/en-us/services/sql-database/) Service Instances From MASB To OSBA

***Note***: this guidance is specially for existing server scenario. Straighter, the SQL Databases created by MASB must be on a server whose credentials are provided from MASB manifest.

## Steps

### Create azure-sql-12-0-dbms-from-existing service instance by OSBA

OSBA doesn't support to provide existing servers' credentials from manifest and create database on them. Instead, it provides **azure-sql-12-0-dbms** service to create SQL servers and you can create database on the server by **azure-sql-12-0-database** service. The service **azure-sql-12-0-dbms-from-existing** is a service using your existing Azure SQL server. It does NOT CREATE new server in provisioning and does NOT DELETE the server in deprovisioning. You can run the following CF CLI command to create a instance of it:

```
cf create-service azure-sql-12-0-dbms-from-existing dbms <instance-name> -c '{"resourceGroup":"<group-name>", "location":"<server-location>", "server":"<server-name>", "administratorLogin":"<login>", "administratorLoginPassword":"<login-password>", "alias":"<instance-name>"}'
```

### Create azure-sql-12-0-database-from-existing

Similar to the server, the service **azure-sql-12-0-database-from-existing** is a service using your existing Azure SQL database. It does NOT CREATE new server in provisioning but DELETE the server in deprovisioning. Not like the server, it is to take over the database. First, you can run the following command to check the service credentials delivered to you application:

```
cf env <your-app-name>
```

Find the database name in field `VCAP_SERVICES` - `azure-sqldb` - `sqldbName`.

Then, you can run the following command to create a instance of **azure-sql-12-0-database-from-existing**:

```
cf create-service azure-sql-12-0-database-from-existing <plan-name> <sqldb-instance-name> -c '{"resourceGroup":"<group-name>", "location":"<server-location>", "server":"<server-name>", "administratorLogin":"<login>", "administratorLoginPassword":"<login-password>", "alias":"<sqldb-instance-name>"}'
```

***Note***: OSBA provides plans by tier categories. For example, compared to MASB provides plans `StandardS0` - `StandardS12`, OSBA provides only a plan `standard`. You should choose the right category. It is important. Though provisioning wouldn't change the tier, you wouldn't be able to update the tier as OSBA doesn't support changing plan for now. The plan name of existing service instances can be checked by `cf services`.

### Duplicate your application and update to adapt the SQL credentials delivered by OSBA

You should check the SQL credential differences between [OSBA](../modules/mssql.md#credentials-1) and [MASB](https://github.com/Azure/meta-azure-service-broker/blob/master/docs/azure-sql-db.md#format-of-credentials). Update how your application utilizes the credentials. Then `cf push` your updated application with another name and another route.

### Bind azure-sql-12-0-database-from-existing to your application

After successfully creating the service instance, bind it to your application:

```
cf bind-service <your-app-new-name> <instance-name> <sqldb-instance-name>
```

Restart the application:

```
cf restart <your-app-new-name>
```

After you test the application and it works well, you can switch your application domain to the new route in your DNS.

### MASB service instance safely clean up

Unbind the old SQL database instance:

```
cf unbind-service <your-app-name> <old-sqldb-instance-name>
```

Purge the service instance (**IMPORTANT**):

```
cf purge-service-instance <old-sqldb-instance-name>
```

Don't use `cf delete-service` here. Or, your database would be deleted in Azure!
