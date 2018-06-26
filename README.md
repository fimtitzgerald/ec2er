## AWS Instance Describer (ec2er) 

Another excuse to play with golang. This command queries running EC2 instances in an account and returns the `project` tag value, `instance id`, and the `Public DNS` of the instance. 

#### Notes
* It's important to note that this command works against cross account sts roles - it does not query your root account or regular IAM. 
* The command expects to get an env var called `AWS_PROFILE_ARN` which is the full ARN of the cross account role to assume. 
* This command also requires the AWS Golang SDK to compile.

