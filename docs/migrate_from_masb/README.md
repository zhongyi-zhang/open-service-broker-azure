This is a list of guidance for migrating Cloud Foundry service instances created by [Meta Azure Service Broker (MASB)](https://github.com/Azure/meta-azure-service-broker) to this advanced Open Service Broker of Azure (OSBA).

Now services supported migration, click to read the corresponding guidance:

  * [Azure SQL Database](./mssql.md)

All the guidance base on the scenario as below:

  * You have a Cloud Foundry cluster and installed CF CLI.

  * You installed MASB, created service instances, and bound them to your application.

  * You installed OSBA.

  * You want to switch to OSBA to take over those service instances and your application still work well.

***Note***: if you don't have any service instances created by MASB after migration, you can use `cf delete-service-broker <MASB-name>` to delete MASB. Also, the broker database of MASB could be deleted.
