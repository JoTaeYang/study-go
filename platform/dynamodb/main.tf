provider "aws" {
    alias                   = "localstack"
    region                  = "us-east-1"  # You can set this to any AWS region
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true
    endpoints {
      dynamodb = "http://localhost:4566" # Use the LocalStack endpoint
    }
    access_key = var.access_key
    secret_key = var.secret_key
}

resource "aws_dynamodb_table" "project_b_table" {
    provider = aws.localstack  # Use the LocalStack-specific provider alias
    name           = "project-st"
    billing_mode   = "PROVISIONED"
    read_capacity = 100
    write_capacity = 100
    hash_key       = "PK"
    range_key      = "SK"

    dynamic attribute {
        for_each = var.attributes

        content {
            name =  attribute.value.name
            type =  attribute.value.type                 
        }        
    }

    dynamic global_secondary_index {
        for_each = var.global_secondary_indexs

        content {
            name =  global_secondary_index.value.name
            hash_key =  global_secondary_index.value.hash_key   
            range_key = global_secondary_index.value.range_key
            write_capacity = global_secondary_index.value.write_capacity
            read_capacity = global_secondary_index.value.read_capacity
            projection_type = global_secondary_index.value.projection_type         
            non_key_attributes = global_secondary_index.value.non_key_attributes == null ? null : global_secondary_index.value.non_key_attributes         
        }     
    }
}

resource "time_static" "table_items" {
    triggers = {
        items = "hello"
    }
}

resource "aws_dynamodb_table_item" "channels" {
    provider = aws.localstack
    table_name = aws_dynamodb_table.project_b_table.name
    hash_key = aws_dynamodb_table.project_b_table.hash_key
    range_key = aws_dynamodb_table.project_b_table.range_key
    

    for_each = {
        chan1 = {
            PK = "ChannelManage"
            SK = "CHAN#01"
        }
        chan2 = {
            PK = "ChannelManage"
            SK = "CHAN#02"
        }
        chan3 = {
            PK = "ChannelManage"
            SK = "CHAN#03"
        }
        chan4 = {
            PK = "ChannelManage"
            SK = "CHAN#04"                        
        }
    }

    item = <<ITEM
    {
        "PK" : {"S" : "${each.value.PK}"},
        "SK" : {"S" : "${each.value.SK}"},
        "cnt" : {"N" : "0"},
        "update_at": {"N" : "${time_static.table_items.unix}"},
        "create_at" : {"N" : "${time_static.table_items.unix}"}
    }
    ITEM
}