variable "attributes" {
    type = list(object({
        name = string
        type = string
    }))

    default = [
        {
            name = "PK"
            type = "S"
        },
        {
            name = "SK"
            type = "S"
        },        
        {
            name = "channel"
            type = "S"
        },
        {
            name = "rank_point"
            type = "N"
        },
        {
            name = "session_key"
            type = "S"
        }
    ]
}

variable "global_secondary_indexs" {
    type = list(object({      
        name = string
        hash_key = string             
        range_key = string
        write_capacity = number
        read_capacity = number
        projection_type = string
        non_key_attributes = optional(list(string))
    }))
   

    default = [
        {                            
            name               = "ChannelSKIndex"
            hash_key           = "channel"
            range_key          = "SK"
            write_capacity     = 50
            read_capacity      = 50
            projection_type    = "KEYS_ONLY"                        
        },
        {                            
            name               = "ChannelRankPointIndex"
            hash_key           = "channel"
            range_key          = "rank_point"
            write_capacity     = 50
            read_capacity      = 50
            projection_type    = "KEYS_ONLY"                        
        },
        {                            
            name               = "session_key-SK-index"
            hash_key           = "session_key"
            range_key          = "SK"
            write_capacity     = 50
            read_capacity      = 50
            projection_type    = "INCLUDE"   
            non_key_attributes = ["grade"]               
        }
    ]
}

variable "access_key" {
  description = "My AWS access key"
  default = "SADASD"
}
variable "secret_key" {
  description = "My AWS secret key"
  default = "AASDF+si/"
}