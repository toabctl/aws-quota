# aws-quota

A tool to list & create quota change requests for **all** regions in parallel.
It's useful if quotas need to be updated for **all** available regions.
The usual environment variables (eg. `AWS_PROFILE` are accepted).

## Usage

List the available services and the service codes (needed for other commands):
```
$ ./aws-quota services-list | head -10
+------------------------------+-----------------------------------------------------------------------+
| SERVICE CODE                 | SERVICE NAME                                                          |
+------------------------------+-----------------------------------------------------------------------+
| AWSCloudMap                  | AWS Cloud Map                                                         |
| a4b                          | Alexa for Business                                                    |
| access-analyzer              | Access Analyzer                                                       |
| account                      | AWS Account Management                                                |
| acm                          | AWS Certificate Manager (ACM)                                         |
| acm-pca                      | AWS Certificate Manager Private Certificate Authority (ACM PCA)       |
| airflow                      | Amazon Managed Workflows for Apache Airflow                           |
```

List the available quotas for a given service code (useful to get the quota code):
```
$ ./aws-quota service-quotas-list --servicecode vpc|head -10
+----------------+----------------------------------------------+------------+-------------+------------+
| REGION         | QUOTA NAME                                   | QUOTA CODE | QUOTA VALUE | ADJUSTABLE |
+----------------+----------------------------------------------+------------+-------------+------------+
| eu-west-3      | Active VPC peering connections per VPC       | L-7E9ECCDB |          50 | true       |
| eu-west-3      | Egress-only internet gateways per Region     | L-45FE3B85 |          50 | true       |
| eu-west-3      | Gateway VPC endpoints per Region             | L-1B52E74A |          20 | true       |
| eu-west-3      | IPv4 CIDR blocks per VPC                     | L-83CA0A9D |           5 | true       |
| eu-west-3      | Inbound or outbound rules per security group | L-0EA8095F |          60 | true       |
| eu-west-3      | Interface VPC endpoints per VPC              | L-29B6F2EB |          50 | true       |
| eu-west-3      | Internet gateways per Region                 | L-A4707A72 |          50 | true       |
```

List the currently available quota change requests:
```
$ ./aws-quota service-quota-history|head -10
+----------------+--------------+------------+---------------+-------------+-------------+
| REGION         | SERVICE CODE | QUOTA CODE | DESIRED VALUE | STATUS      | CASE ID     |
+----------------+--------------+------------+---------------+-------------+-------------+
| eu-west-2      | vpc          | L-F678F1CE |            50 | APPROVED    |             |
| eu-west-3      | vpc          | L-F678F1CE |            50 | APPROVED    |             |
| eu-north-1     | vpc          | L-F678F1CE |            50 | APPROVED    |             |
| eu-central-1   | vpc          | L-F678F1CE |            50 | APPROVED    |             |
| eu-south-1     | vpc          | L-F678F1CE |            50 | APPROVED    |             |
| eu-west-1      | vpc          | L-F678F1CE |            50 | APPROVED    |             |
| me-south-1     | vpc          | L-F678F1CE |            50 | APPROVED    |             |
```

Create a quota change request:
```
$ ./aws-quota service-quota-increase --servicecode ec2 --quotacode L-0263D0A3 --quotavalue 50
+----------------+--------------+------------+---------------+---------+
| REGION         | SERVICE CODE | QUOTA CODE | DESIRED VALUE | CASE ID |
+----------------+--------------+------------+---------------+---------+
| eu-central-1   | ec2          | L-0263D0A3 |            50 |         |
| eu-west-3      | ec2          | L-0263D0A3 |            50 |         |
| eu-south-1     | ec2          | L-0263D0A3 |            50 |         |
| eu-north-1     | ec2          | L-0263D0A3 |            50 |         |
| eu-west-2      | ec2          | L-0263D0A3 |            50 |         |
| eu-west-1      | ec2          | L-0263D0A3 |            50 |         |
| me-south-1     | ec2          | L-0263D0A3 |            50 |         |
| ca-central-1   | ec2          | L-0263D0A3 |            50 |         |
```
